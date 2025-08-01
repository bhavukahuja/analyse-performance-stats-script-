#!/bin/zsh

echo "Server Performance Stats"
echo "======================="

# Total CPU usage
echo "CPU Usage:"
top -l 1 | grep "CPU usage"

# Total memory usage (Free vs Used including percentage)
echo
echo "Memory Usage:"
vm_stat | awk '
/page size of/ {psize=$8}
/Pages free/ {free=$3}
/Pages active/ {active=$3}
/Pages inactive/ {inactive=$3}
/Pages speculative/ {spec=$3}
/Pages wired down/ {wired=$4}
/Pages throttled/ {throttled=$3}
/Pages purgeable/ {purgeable=$3}
/Pages occupied by compressor/ {compressed=$5}
END {
  if (psize == "") psize=4096;
  free=free+0; active=active+0; inactive=inactive+0; spec=spec+0; wired=wired+0; throttled=throttled+0; purgeable=purgeable+0; compressed=compressed+0;
  total=free+active+inactive+spec+wired+throttled+purgeable+compressed;
  used=active+inactive+spec+wired+throttled+purgeable+compressed;
  total_bytes=total*psize;
  used_bytes=used*psize;
  free_bytes=free*psize;
  total_gb=total_bytes/1024/1024/1024;
  used_gb=used_bytes/1024/1024/1024;
  free_gb=free_bytes/1024/1024/1024;
  percent=(used_gb/total_gb)*100;
  printf "Used: %.2f GB, Free: %.2f GB, Total: %.2f GB (%.2f%% used)\n", used_gb, free_gb, total_gb, percent;
}'

# Total disk usage (Free vs Used including percentage)
echo
echo "Disk Usage:"
df -h / | awk 'NR==2 {printf "Used: %s, Free: %s, Usage: %s\n", $3, $4, $5}'

# OS version
echo
echo "OS Version:"
sw_vers

# Uptime
echo
echo "Uptime:"
uptime

# Load average
echo
echo "Load Average:"
uptime | awk -F'load averages?: ' '{print $2}'

# Logged in users
echo
echo "Logged in users:"
who | wc -l
who

# Failed login attempts (not supported on macOS by default)
echo
echo "Failed login attempts: Not supported on this OS"

# Top 5 processes by CPU usage
echo
echo "Top 5 Processes by CPU Usage:"
ps -axo pid,user,%cpu,%mem,comm | sort -k3 -nr | head -n 6

# Top 5 processes by Memory usage
echo
echo "Top 5 Processes by Memory Usage:"
ps -axo pid,user,%cpu,%mem,comm | sort -k4 -nr | head -n 6

# Go version and OS info (if Go is installed)
if command -v go >/dev/null 2>&1; then
    echo
    echo "Go version:"
    go version
    echo "OS info:"
    
