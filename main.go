package main

import (
	"bufio"
	"fmt"
	"os"
	"simple-monitor/systeminfo"
	"strconv"
	"strings"
)

// MainMenuOptions represents the main menu options
type MainMenuOptions struct {
	StartMonitoring bool
	Settings        bool
	Developer       bool
	Quit            bool
}

// MonitoringOptions represents the monitoring submenu options
type MonitoringOptions struct {
	SystemInfo bool
	CPU        bool
	Memory     bool
	Disk       bool
	Network    bool
	Processes  bool
	Back       bool
}

// displayMainMenu shows the main menu to the user
func displayMainMenu() {
	fmt.Println("\n🖥️  Simple Monitor v1.0")
	fmt.Println(strings.Repeat("-", 30))
	fmt.Println("1. Start Monitoring")
	fmt.Println("2. Settings")
	fmt.Println("3. Developer")
	fmt.Println("4. Quit")
	fmt.Println(strings.Repeat("-", 30))
	fmt.Print("Select option (1-4): ")
}

// displayMonitoringMenu shows the monitoring submenu
func displayMonitoringMenu() {
	fmt.Println("\n📊 Monitoring Options")
	fmt.Println(strings.Repeat("-", 30))
	fmt.Println("1. System Information")
	fmt.Println("2. CPU Monitor")
	fmt.Println("3. Memory Monitor")
	fmt.Println("4. Disk Monitor")
	fmt.Println("5. Network Monitor")
	fmt.Println("6. Process Monitor")
	fmt.Println("7. Back to Main Menu")
	fmt.Println(strings.Repeat("-", 30))
	fmt.Print("Select option (1-7): ")
}

// getUserChoice gets user input and validates it for main menu
func getUserChoice(maxOptions int) int {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		scanner.Scan()
		input := strings.TrimSpace(scanner.Text())

		choice, err := strconv.Atoi(input)
		if err != nil {
			fmt.Printf("Invalid input! Enter 1-%d: ", maxOptions)
			continue
		}

		if choice < 1 || choice > maxOptions {
			fmt.Printf("Invalid option! Enter 1-%d: ", maxOptions)
			continue
		}

		return choice
	}
}

// handleMainMenuChoice processes the user's main menu choice
func handleMainMenuChoice(choice int) {
	// Clear screen after selection
	fmt.Print("\033[2J\033[H")

	switch choice {
	case 1:
		fmt.Println("🚀 Start Monitoring")
		fmt.Println(strings.Repeat("-", 30))
		startMonitoring()
	case 2:
		fmt.Println("⚙️  Settings")
		fmt.Println(strings.Repeat("-", 30))
		showSettings()
	case 3:
		fmt.Println("👨‍💻 Developer")
		fmt.Println(strings.Repeat("-", 30))
		showDeveloper()
	case 4:
		fmt.Println("👋 Goodbye! Thank you for using Simple Monitor.")
		os.Exit(0)
	}
}

// startMonitoring handles the monitoring submenu
func startMonitoring() {
	for {
		displayMonitoringMenu()
		choice := getUserChoice(7)

		// Clear screen after selection
		fmt.Print("\033[2J\033[H")

		switch choice {
		case 1:
			fmt.Println("📊 System Information")
			fmt.Println(strings.Repeat("-", 30))
			showSystemInfo()
		case 2:
			fmt.Println("🖥️  CPU Monitor")
			fmt.Println(strings.Repeat("-", 30))
			monitorCPU()
		case 3:
			fmt.Println("💾 Memory Monitor")
			fmt.Println(strings.Repeat("-", 30))
			monitorMemory()
		case 4:
			fmt.Println("💿 Disk Monitor")
			fmt.Println(strings.Repeat("-", 30))
			monitorDisk()
		case 5:
			fmt.Println("🌐 Network Monitor")
			fmt.Println(strings.Repeat("-", 30))
			monitorNetwork()
		case 6:
			fmt.Println("⚙️  Process Monitor")
			fmt.Println(strings.Repeat("-", 30))
			showProcesses()
		case 7:
			fmt.Println("⬅️  Returning to main menu...")
			return
		}
	}
}

// showSettings displays settings menu
func showSettings() {
	fmt.Println("Settings - Coming soon...")
	waitForEnter()
}

// showDeveloper displays developer options
func showDeveloper() {
	fmt.Println("Developer options - Coming soon...")
	waitForEnter()
}

// System information manager instance
var systemInfoManager = systeminfo.NewSystemInfoManager()

// showSystemInfo displays comprehensive system information
func showSystemInfo() {
	if err := systemInfoManager.ShowSystemInfo(); err != nil {
		fmt.Printf("❌ Error displaying system information: %v\n", err)
	}
	waitForEnter()
}

func monitorCPU() {
	fmt.Println("CPU Monitor - Coming soon...")
	waitForEnter()
}

func monitorMemory() {
	fmt.Println("Memory Monitor - Coming soon...")
	waitForEnter()
}

func monitorDisk() {
	fmt.Println("Disk Monitor - Coming soon...")
	waitForEnter()
}

func monitorNetwork() {
	fmt.Println("Network Monitor - Coming soon...")
	waitForEnter()
}

func showProcesses() {
	fmt.Println("Process Monitor - Coming soon...")
	waitForEnter()
}

// waitForEnter waits for user to press Enter
func waitForEnter() {
	fmt.Print("\nPress Enter to continue...")
	bufio.NewScanner(os.Stdin).Scan()
}

func main() {
	fmt.Println("🚀 Simple Monitor started!")

	for {
		displayMainMenu()
		choice := getUserChoice(4)
		handleMainMenuChoice(choice)
	}
}
