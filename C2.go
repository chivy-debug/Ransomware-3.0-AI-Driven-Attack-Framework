package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/tealeg/xlsx"
)

// ================= TYPES =================

type LogEntry struct {
	Time        string `json:"time"`
	Phase       string `json:"phase"`
	LLMResponse string `json:"llm_response"`
	Code        string `json:"code"`
	Status      string `json:"status"`
}

// ================= GLOBAL =================

var (
	logs   []LogEntry
	logMux sync.Mutex

	clients   = make(map[*websocket.Conn]bool)
	clientsMu sync.Mutex

	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	uploadDir   string
	downloadDir string
)

// ================= MAIN =================

func main() {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	uploadDir = filepath.Join(wd, "upload")
	downloadDir = filepath.Join(wd, "download")

	os.MkdirAll(uploadDir, 0755)
	os.MkdirAll(downloadDir, 0755)

	http.HandleFunc("/log", logHandler)
	http.HandleFunc("/logs", logsHandler)
	http.HandleFunc("/ws", wsHandler)

	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/download/", downloadHandler)

	http.HandleFunc("/", dashboardHandler)

	log.Println("C2 running at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// ================= LOG HANDLER =================

func logHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "only POST", 405)
		return
	}

	var entry LogEntry
	if err := json.NewDecoder(r.Body).Decode(&entry); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	logMux.Lock()

	logs = append(logs, entry)

	// limit logs (anti RAM leak)
	if len(logs) > 2000 {
		logs = logs[len(logs)-1000:]
	}

	// copy để save tránh race
	logsCopy := make([]LogEntry, len(logs))
	copy(logsCopy, logs)

	logMux.Unlock()

	// save file
	go saveJSON(logsCopy)
	go saveXLSX(logsCopy)

	// realtime push
	broadcast(entry)

	log.Println("LOG:", entry.Phase, entry.Status)

	w.WriteHeader(200)
}

// ================= GET LOGS =================

func logsHandler(w http.ResponseWriter, r *http.Request) {
	logMux.Lock()
	data := make([]LogEntry, len(logs))
	copy(data, logs)
	logMux.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// ================= WEBSOCKET =================

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	clientsMu.Lock()
	clients[conn] = true
	clientsMu.Unlock()

	// send history
	logMux.Lock()
	for _, l := range logs {
		conn.WriteJSON(l)
	}
	logMux.Unlock()
}

func broadcast(entry LogEntry) {
	clientsMu.Lock()
	defer clientsMu.Unlock()

	for c := range clients {
		go func(conn *websocket.Conn) {
			conn.SetWriteDeadline(time.Now().Add(2 * time.Second))

			if err := conn.WriteJSON(entry); err != nil {
				conn.Close()
				delete(clients, conn)
			}
		}(c)
	}
}

// ================= FILE UPLOAD =================

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "only POST", 405)
		return
	}

	filename := filepath.Base(r.Header.Get("X-Filename"))
	if filename == "" {
		http.Error(w, "missing filename", 400)
		return
	}

	uploadPath := filepath.Join(uploadDir, filename)
	downloadPath := filepath.Join(downloadDir, filename)

	os.MkdirAll(filepath.Dir(uploadPath), 0755)

	dst, err := os.Create(uploadPath)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, r.Body); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	copyFile(uploadPath, downloadPath)

	log.Println("Uploaded:", filename)

	w.WriteHeader(201)
}

// ================= DOWNLOAD =================

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	filename := filepath.Base(r.URL.Path)
	path := filepath.Join(downloadDir, filename)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		http.NotFound(w, r)
		return
	}

	http.ServeFile(w, r, path)
}

func copyFile(src, dst string) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer out.Close()

	io.Copy(out, in)
}

// ================= SAVE =================

func saveJSON(data []LogEntry) {
	f, err := os.Create("logs.json")
	if err != nil {
		log.Println("JSON save error:", err)
		return
	}
	defer f.Close()

	json.NewEncoder(f).Encode(data)
}

func saveXLSX(data []LogEntry) {
	file := xlsx.NewFile()
	sheet, _ := file.AddSheet("Logs")

	header := sheet.AddRow()
	for _, h := range []string{"Time", "Phase", "Status", "LLM", "Code"} {
		header.AddCell().Value = h
	}

	for _, e := range data {
		row := sheet.AddRow()
		row.AddCell().Value = e.Time
		row.AddCell().Value = e.Phase
		row.AddCell().Value = e.Status
		row.AddCell().Value = e.LLMResponse
		row.AddCell().Value = e.Code
	}

	if err := file.Save("logs.xlsx"); err != nil {
		log.Println("XLSX save error:", err)
	}
}

// ================= DASHBOARD =================

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	html := `
<!doctype html>
<html>
<head>
<title>C2 Dashboard</title>
<style>
body{font-family:Arial;margin:20px}
table{border-collapse:collapse;width:100%}
th,td{border:1px solid #ccc;padding:8px}
th{background:#f0f0f0}
ul{list-style:none;padding:0}
.success { color:#22c55e }
.fail { color:#ef4444 }
</style>
</head>
<body>

<h1>C2 Realtime Dashboard</h1>

<table id="log">
<thead>
<tr>
<th>Time</th>
<th>Phase</th>
<th>Status</th>
<th>LLM</th>
<th>Code</th>
</tr>
</thead>
<tbody></tbody>
</table>

<script>
const ws = new WebSocket("ws://" + location.host + "/ws");

ws.onmessage = function(event){
	const e = JSON.parse(event.data);

	const tr = document.createElement("tr");

	tr.innerHTML = 
	"<td>"+e.time+"</td>" +
	"<td>"+e.phase+"</td>" +
	"<td class='"+e.status+"'>"+e.status+"</td>" +
	"<td>"+escapeHtml(e.llm_response)+"</td>" +
	"<td>"+escapeHtml(e.code)+"</td>";

	document.querySelector("#log tbody").prepend(tr);
};

function escapeHtml(text){
	return text.replace(/</g,"&lt;");
}
</script>

</body>
</html>
`
	w.Write([]byte(html))
}
