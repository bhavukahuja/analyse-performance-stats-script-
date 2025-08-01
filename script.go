package main

import (
	"fmt"
	// "os"
	"os/exec"
	"runtime"
	"sort"
	"strconv"
	"strings"
)

type ProcessStat struct {
	PID     string
	User    string
	CPU     float64
	MEM     float64
	Command string
}

func getCPUUsage() (string, error) {
	out, err := exec.Command("sh", "-c", "top -l 1 | grep 'CPU usage'").Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func getMemUsage() (string, error) {
	out, err := exec.Command("sh", "-c", "vm_stat").Output()
	if err != nil {
		return "", err
	}
	lines := strings.Split(string(out), "\n")
	pageSize := 4096 // default
	totalPages := 0
	freePages := 0
	for _, line := range lines {
		if strings.Contains(line, "page size of") {
			parts := strings.Fields(line)
			if len(parts) > 3 {
				pageSize, _ = strconv.Atoi(parts[3])
			}
		}
		if strings.Contains(line, "Pages free") {
			parts := strings.Fields(line)
			freePages, _ = strconv.Atoi(strings.Trim(parts[2], "."))
		}
		if strings.Contains(line, "Pages active") ||
			strings.Contains(line, "Pages inactive") ||
			strings.Contains(line, "Pages speculative") ||
			strings.Contains(line, "Pages wired down") ||
			strings.Contains(line, "Pages throttled") ||
			strings.Contains(line, "Pages purgeable") ||
			strings.Contains(line, "Pages occupied by compressor") {
			parts := strings.Fields(line)
			totalPages += func() int {
				v, _ := strconv.Atoi(strings.Trim(parts[2], "."))
				return v
			}()
		}
	}
	totalMem := float64((totalPages+freePages)*pageSize) / (1024 * 1024 * 1024)
	usedMem := float64(totalPages*pageSize) / (1024 * 1024 * 1024)
	freeMem := float64(freePages*pageSize) / (1024 * 1024 * 1024)
	percentUsed := (usedMem / totalMem) * 100
	return fmt.Sprintf("Memory Usage: Used: %.2f GB, Free: %.2f GB, Total: %.2f GB (%.2f%% used)", usedMem, freeMem, totalMem, percentUsed), nil
}

func getDiskUsage() (string, error) {
	out, err := exec.Command("sh", "-c", "df -h /").Output()
	if err != nil {
		return "", err
	}
	lines := strings.Split(string(out), "\n")
	if len(lines) < 2 {
		return "", fmt.Errorf("unexpected df output")
	}
	fields := strings.Fields(lines[1])
	if len(fields) < 5 {
		return "", fmt.Errorf("unexpected df output fields")
	}
	used := fields[2]
	free := fields[3]
	percent := fields[4]
	return fmt.Sprintf("Disk Usage: Used: %s, Free: %s, Usage: %s", used, free, percent), nil
}

func getTopProcesses(sortBy string) ([]ProcessStat, error) {

	out, err := exec.Command("sh", "-c", "ps -axo pid,user,%cpu,%mem,comm").Output()
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(out), "\n")
	var procs []ProcessStat
	for i, line := range lines {
		if i == 0 || strings.TrimSpace(line) == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 5 {
			continue
		}
		cpu, _ := strconv.ParseFloat(fields[2], 64)
		mem, _ := strconv.ParseFloat(fields[3], 64)
		procs = append(procs, ProcessStat{
			PID:     fields[0],
			User:    fields[1],
			CPU:     cpu,
			MEM:     mem,
			Command: strings.Join(fields[4:], " "),
		})
	}
	if sortBy == "cpu" {
		sort.Slice(procs, func(i, j int) bool { return procs[i].CPU > procs[j].CPU })
	} else {
		sort.Slice(procs, func(i, j int) bool { return procs[i].MEM > procs[j].MEM })
	}
	if len(procs) > 5 {
		procs = procs[:5]
	}
	return procs, nil
}

func getOSVersion() (string, error) {
	out, err := exec.Command("sw_vers").Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func getUptime() (string, error) {
	out, err := exec.Command("uptime").Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func getLoadAverage() (string, error) {
	out, err := exec.Command("sh", "-c", "uptime | awk -F'load averages?: ' '{print $2}'").Output()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Load Average: %s", strings.TrimSpace(string(out))), nil
}

func getLoggedInUsers() (string, error) {
	out, err := exec.Command("who").Output()
	if err != nil {
		return "", err
	}
	users := strings.Split(strings.TrimSpace(string(out)), "\n")
	if len(users) == 1 && users[0] == "" {
		return "Logged in users: 0", nil
	}
	return fmt.Sprintf("Logged in users: %d\n%s", len(users), string(out)), nil
}

func getFailedLoginAttempts() (string, error) {
	// lastb is not available on macOS by default, so we return not supported
	return "Failed login attempts: Not supported on this OS", nil
}

func main() {
	fmt.Println("Server Performance Stats")
	fmt.Println("=======================")

	cpu, err := getCPUUsage()
	if err != nil {
		fmt.Println("CPU Usage: Error:", err)
	} else {
		fmt.Println(cpu)
	}

	mem, err := getMemUsage()
	if err != nil {
		fmt.Println("Memory Usage: Error:", err)
	} else {
		fmt.Println(mem)
	}

	disk, err := getDiskUsage()
	if err != nil {
		fmt.Println("Disk Usage: Error:", err)
	} else {
		fmt.Println(disk)
	}

	osver, err := getOSVersion()
	if err != nil {
		fmt.Println("OS Version: Error:", err)
	} else {
		fmt.Println("\nOS Version:\n" + osver)
	}

	uptime, err := getUptime()
	if err != nil {
		fmt.Println("Uptime: Error:", err)
	} else {
		fmt.Println("\nUptime:", uptime)
	}

	load, err := getLoadAverage()
	if err != nil {
		fmt.Println("Load Average: Error:", err)
	} else {
		fmt.Println(load)
	}

	users, err := getLoggedInUsers()
	if err != nil {
		fmt.Println("Logged in users: Error:", err)
	} else {
		fmt.Println("\n" + users)
	}

	failed, err := getFailedLoginAttempts()
	if err != nil {
		fmt.Println("Failed login attempts: Error:", err)
	} else {
		fmt.Println("\n" + failed)
	}

	fmt.Println("\nTop 5 Processes by CPU Usage:")
	procs, err := getTopProcesses("cpu")
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		for _, p := range procs {
			fmt.Printf("PID: %s, User: %s, CPU: %.2f%%, MEM: %.2f%%, CMD: %s\n", p.PID, p.User, p.CPU, p.MEM, p.Command)
		}
	}

	fmt.Println("\nTop 5 Processes by Memory Usage:")
	procs, err = getTopProcesses("mem")
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		for _, p := range procs {
			fmt.Printf("PID: %s, User: %s, CPU: %.2f%%, MEM: %.2f%%, CMD: %s\n", p.PID, p.User, p.CPU, p.MEM, p.Command)
		}
	}

	fmt.Printf("\nGo version: %s\n", runtime.Version())
	fmt.Printf("OS: %s/%s\n", runtime.GOOS, runtime.GOARCH)
}
