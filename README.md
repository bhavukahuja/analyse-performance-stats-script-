🖥️ Server Performance Monitor
This project contains:

A Go script that collects and prints key system performance metrics.

A Zsh script that performs a similar task using macOS built-in CLI tools.


📁 Contents
Script.go: Go-based server performance tool.

monitor.zsh: Zsh shell script for performance monitoring.

README.md: This documentation.


⚙️ Features
Both scripts report the following system statistics:

✅ CPU Usage

✅ Memory Usage (Free/Used with % calculation)

✅ Disk Usage

✅ OS Version

✅ System Uptime

✅ Load Average

✅ Logged-in Users

✅ Top 5 Processes by CPU and Memory Usage

✅ Go Version (if available)

❌ Failed Login Attempts (macOS not supported)


🏁 How to Run
Go Script
1. Prerequisites:
Go installed (go version)

MacOS (the Go script uses top, vm_stat, and other macOS-specific commands)
 
2. Run:
1: bash
go run main.go

2: build and run:
bash
go build -o perfmonitor Script.go
./perfmonitor



Zsh Script
1. Make the script executable:
bash
chmod +x monitor.zsh

2. Run the script:
bash
./monitor.zsh

Or run directly with Zsh:
bash
zsh monitor.zsh

📌 Notes
These scripts are macOS specific due to the use of tools like top -l, vm_stat, and sw_vers.

Linux users would need equivalent commands (top, free, df, etc.).

🤝 Contributing
Feel free to fork the repo and submit PRs to add:

1:Linux compatibility
2:Network stats
3:JSON/CSV output formats
4:Web dashboard integration (e.g., using Go + HTML)
