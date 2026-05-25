# 🤖 Ransomware 3.0: AI-Driven Attack Framework

> A cutting-edge, AI-driven framework for advanced ransomware attack simulation and cybersecurity research.

---

## ⚠️ Disclaimer

This project is intended **strictly for educational purposes**, cybersecurity research, and ethical red teaming exercises only. It should **never** be used to conduct illegal activities, cause harm, or compromise systems without explicit, informed consent. Misuse of this software for malicious purposes is strictly prohibited and the author does not endorse or take responsibility for any such actions.

---

## 📖 Overview

The **"Ransomware 3.0 AI-Driven Attack Framework"** represents a sophisticated approach to understanding and simulating next-generation ransomware threats. This framework explores the integration of Large Language Models (LLMs) and artificial intelligence into the ransomware attack lifecycle, enabling more adaptive, autonomous, and evasive operations compared to traditional methods.

The architecture and conceptual model of this framework are heavily inspired by cutting-edge cybersecurity research, including:

- **LLM-Orchestrated Malware:** Concepts modeled after the paradigm shift detailed in the research paper *"Ransomware 3.0: Self-Composing and LLM-Orchestrated"* ([arXiv:2508.20444](https://arxiv.org/abs/2508.20444)), which explores how autonomous agents can dynamically compose attack chains.

- **Real-world AI Threats:** Insights gained from emerging real-world threats such as **Promptlock**, the first documented AI-powered ransomware discovered by ESET, which utilizes LLMs to dynamically bypass security controls and optimize execution. Read the full analysis on the [ESET Research Newsroom](https://www.eset.com/int/about/newsroom/).

Designed for security researchers, penetration testers, and academic institutions, this project provides a foundation for:

- Analyzing potential future threats
- Developing advanced detection mechanisms
- Strengthening enterprise defensive strategies against intelligent, adaptive malware

It comprises two core Go-based components: a **Command and Control (C2) server** for managing compromised systems and an **Orchestrator** responsible for coordinating the overall AI-driven attack flow via local LLM instances.

---

## ✨ Features

| Feature | Description |
|---|---|
| 🧠 **AI-Driven Attack Logic** | Adaptive targeting, autonomous decision-making, and dynamic evasion throughout the attack chain |
| 🖥️ **Modular C2 Server** | Robust backend (`C2.go`) for secure implant communication, command issuance, data exfiltration, and key delivery |
| ⚙️ **Attack Orchestration Engine** | `Orchestrator.go` coordinates multi-stage attacks: initial compromise, persistence, lateral movement, encryption, and key management |
| 🌐 **Cross-Platform Potential** | Built with Go for inherent cross-platform compatibility across various operating systems |

---

## 🛠️ Tech Stack

- **Language:** Go (Golang)
- **AI Backend:** [Ollama](https://ollama.com/) (local LLM runtime)
- **Supported Models:** llama3, mistral, qwen2.5, and more

---

## 🦙 Ollama — Local LLM Runtime

[Ollama](https://ollama.com/) is an open-source tool that allows you to run large language models (LLMs) **locally on your own machine**, without relying on external cloud APIs. It acts as the AI backbone of this framework, enabling the Orchestrator to make intelligent, adaptive decisions in real time during the attack simulation.

### Why Ollama?

- **Privacy & Isolation:** All LLM inference happens locally — no data leaves your machine.
- **Model Flexibility:** Supports a wide range of open-source models (Llama 3, Mistral, Qwen 2.5, Gemma, etc.).
- **Simple REST API:** Exposes a lightweight HTTP API (`/api/generate`) that the Orchestrator communicates with directly.
- **Low Latency:** Significantly faster response times compared to remote API calls, critical for real-time attack orchestration.

### Installing Ollama

```bash
# Linux / macOS
curl -fsSL https://ollama.com/install.sh | sh

# Windows
# Download the installer from: https://ollama.com/download
```

### Pulling a Model

After installation, pull the LLM model you intend to use:

```bash
ollama pull llama3       # Meta Llama 3 (recommended)
ollama pull mistral      # Mistral 7B
ollama pull qwen2.5      # Qwen 2.5
```

### Starting the Ollama Server

> ⚠️ **Ollama must be running before you execute any Go components.**

```bash
ollama serve
```

By default, Ollama listens on `http://localhost:11434`. Verify it is running:

```bash
curl http://localhost:11434/api/tags
```

You should see a JSON response listing your available local models.

---

## 🚀 Quick Start

### Prerequisites

- [Go (Golang)](https://go.dev/) version **1.16 or higher**
- [Ollama](https://ollama.com/) installed and **running** with at least one model pulled

### Installation & Configuration

**1. Clone the repository**

```bash
git clone https://github.com/chivy-debug/Ransomware-3.0-AI-Driven-Attack-Framework.git
cd Ransomware-3.0-AI-Driven-Attack-Framework
```

**2. Configure Environment Variables**

Before running the framework, open the source files and configure the LLM endpoints, C2 API server addresses, and the model you intend to use:

```go
OllamaAPI = "http://<IP>/api/generate"
C2API     = "http://<IP>"
Model     = "<MODEL LLM>" // e.g., llama3, mistral, qwen2.5, etc.
```

**3. Initialize Go Modules**

```bash
go mod init Ransomware-3.0-AI-Driven-Attack-Framework  # If go.mod is not present
go mod tidy                                             # To download dependencies
```

### Running Components

The framework consists of two main executable components. These are typically run as **separate processes**.

> ⚠️ **Important:** You must start the Ollama server **before** running any Go component. The Orchestrator depends on a live Ollama instance to function.

```bash
# Step 0 — Start Ollama first (keep this terminal open)
ollama serve
```

#### 1. Running the Command & Control (C2) Server

```bash
# Run directly
go run C2.go

# Or build and run the executable
go build -o c2 C2.go
./c2
```

#### 2. Running the Orchestrator

```bash
# Run directly
go run Orchestrator.go

# Or build and run the executable
go build -o orchestrator Orchestrator.go
./orchestrator
```

---

## 📁 Project Structure

```
Ransomware-3.0-AI-Driven-Attack-Framework/
├── C2.go            # Command and Control (C2) server implementation
├── Orchestrator.go  # Core attack orchestration logic, including AI/Ollama integration
└── README.md        # Project documentation
```

---

## 🔧 Development

### Building Executables

To create standalone binaries for deployment:

```bash
# Build C2 server
go build -o c2 C2.go

# Build Orchestrator
go build -o orchestrator Orchestrator.go
```

### Cleaning Up

```bash
# Linux/macOS
rm c2 orchestrator

# Windows
del c2.exe orchestrator.exe
```

---

## 🤝 Contributing

We welcome contributions to enhance the framework, especially in areas like:

- AI integration
- Defense evasion research
- Behavioral detection signatures

Please ensure any contributions align with the **ethical use guidelines** outlined in the disclaimer.

1. Fork the repository
2. Create a new branch: `git checkout -b feature/amazing-feature`
3. Commit your changes: `git commit -m 'Add amazing feature'`
4. Push to the branch: `git push origin feature/amazing-feature`
5. Open a **Pull Request**

---

## 🎬 Video Demo

Watch the framework in action — the demo videos cover end-to-end attack simulation, C2 communication, and AI-driven orchestration behavior.

> 📁 **[View Demo Videos on Google Drive](https://drive.google.com/drive/folders/14krAfItgylRvVqQdn6-hku3I1BxKDWWM?usp=sharing)**

The demo folder includes:

- Full attack lifecycle walkthrough
- C2 server & Orchestrator interaction
- AI decision-making in real time

---

## 📄 License

This project is licensed under the LICENSE_NAME — see the `LICENSE` file for details.

> **TODO:** Add a LICENSE file.

---

> ⭐ **Star this repo** if you find it helpful for your research!
>
> Made with ❤️ by [chivy-debug](https://github.com/chivy-debug)
