package main

import (
	"bufio"
	"fmt"
	"os"
	"simple-monitor/cpumonitor"
	"simple-monitor/diskmonitor"
	"simple-monitor/memorymonitor"
	"simple-monitor/networkmonitor"
	"simple-monitor/processmonitor"
	"simple-monitor/systeminfo"
	"strconv"
	"strings"
	"time"
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
	for {
		fmt.Println("\n⚙️  Settings")
		fmt.Println(strings.Repeat("-", 30))
		fmt.Println("1. Export Settings")
		fmt.Println("2. Display Settings")
		fmt.Println("3. Monitoring Settings")
		fmt.Println("4. Back to Main Menu")
		fmt.Println(strings.Repeat("-", 30))
		fmt.Print("Select option (1-4): ")

		choice := getUserChoice(4)

		switch choice {
		case 1:
			showExportSettings()
		case 2:
			showDisplaySettings()
		case 3:
			showMonitoringSettings()
		case 4:
			return
		}
	}
}

// showDeveloper displays developer options
func showDeveloper() {
	fmt.Println("Developer options - Coming soon...")
	waitForEnter()
}

// System information manager instance
var systemInfoManager = systeminfo.NewSystemInfoManager()

// CPU monitor manager instance
var cpuMonitorManager = cpumonitor.NewCPUMonitorManager()

// Memory monitor manager instance
var memoryMonitorManager = memorymonitor.NewMemoryMonitorManager()

// Disk monitor manager instance
var diskMonitorManager = diskmonitor.NewDiskMonitorManager()

// Network monitor manager instance
var networkMonitorManager = networkmonitor.NewNetworkMonitorManager()
var processMonitorManager = processmonitor.NewProcessMonitorManager()

// showSystemInfo displays comprehensive system information
func showSystemInfo() {
	if err := systemInfoManager.ShowSystemInfo(); err != nil {
		fmt.Printf("❌ Error displaying system information: %v\n", err)
	}
	waitForEnter()
}

func monitorCPU() {
	fmt.Println("🖥️  CPU Monitor")
	fmt.Println(strings.Repeat("-", 30))
	fmt.Println("1. Live Monitoring")
	fmt.Println("2. Single Snapshot")
	fmt.Println("3. Back to Monitoring Menu")
	fmt.Println(strings.Repeat("-", 30))
	fmt.Print("Select option (1-3): ")

	choice := getUserChoice(3)

	switch choice {
	case 1:
		fmt.Println("Starting live CPU monitoring...")
		if err := cpuMonitorManager.StartLiveMonitoring(); err != nil {
			fmt.Printf("❌ Error starting CPU monitoring: %v\n", err)
		}
		waitForEnter()
	case 2:
		if err := cpuMonitorManager.StartSingleSnapshot(); err != nil {
			fmt.Printf("❌ Error displaying CPU information: %v\n", err)
		}
		waitForEnter()
	case 3:
		return
	}
}

func monitorMemory() {
	fmt.Println("💾 Memory Monitor")
	fmt.Println(strings.Repeat("-", 30))
	fmt.Println("1. Live Monitoring")
	fmt.Println("2. Single Snapshot")
	fmt.Println("3. Back to Monitoring Menu")
	fmt.Println(strings.Repeat("-", 30))
	fmt.Print("Select option (1-3): ")

	choice := getUserChoice(3)

	switch choice {
	case 1:
		fmt.Println("Starting live memory monitoring...")
		if err := memoryMonitorManager.StartLiveMonitoring(); err != nil {
			fmt.Printf("❌ Error starting memory monitoring: %v\n", err)
		}
		waitForEnter()
	case 2:
		if err := memoryMonitorManager.StartSingleSnapshot(); err != nil {
			fmt.Printf("❌ Error displaying memory information: %v\n", err)
		}
		waitForEnter()
	case 3:
		return
	}
}

func monitorDisk() {
	fmt.Println("💿 Disk Monitor")
	fmt.Println(strings.Repeat("-", 30))
	fmt.Println("1. Live Monitoring")
	fmt.Println("2. Single Snapshot")
	fmt.Println("3. Back to Monitoring Menu")
	fmt.Println(strings.Repeat("-", 30))
	fmt.Print("Select option (1-3): ")

	choice := getUserChoice(3)

	switch choice {
	case 1:
		fmt.Println("Starting live disk monitoring...")
		if err := diskMonitorManager.StartLiveMonitoring(); err != nil {
			fmt.Printf("❌ Error starting disk monitoring: %v\n", err)
		}
		waitForEnter()
	case 2:
		if err := diskMonitorManager.StartSingleSnapshot(); err != nil {
			fmt.Printf("❌ Error displaying disk information: %v\n", err)
		}
		waitForEnter()
	case 3:
		return
	}
}

func monitorNetwork() {
	fmt.Println("🌐 Network Monitor")
	fmt.Println(strings.Repeat("-", 30))
	fmt.Println("1. Live Monitoring")
	fmt.Println("2. Single Snapshot")
	fmt.Println("3. Back to Monitoring Menu")
	fmt.Println(strings.Repeat("-", 30))
	fmt.Print("Select option (1-3): ")

	choice := getUserChoice(3)

	switch choice {
	case 1:
		fmt.Println("Starting live network monitoring...")
		if err := networkMonitorManager.StartLiveMonitoring(); err != nil {
			fmt.Printf("❌ Error starting network monitoring: %v\n", err)
		}
		waitForEnter()
	case 2:
		if err := networkMonitorManager.StartSingleSnapshot(); err != nil {
			fmt.Printf("❌ Error displaying network information: %v\n", err)
		}
		waitForEnter()
	case 3:
		return
	}
}

func showProcesses() {
	fmt.Println("Process Monitor - Coming soon...")
	waitForEnter()
}

// showExportSettings displays export settings menu
func showExportSettings() {
	for {
		fmt.Println("\n📁 Export Settings")
		fmt.Println(strings.Repeat("-", 30))
		fmt.Println("1. Set Export Interval")
		fmt.Println("2. Set Export Format")
		fmt.Println("3. Enable/Disable Export")
		fmt.Println("4. Back to Settings")
		fmt.Println(strings.Repeat("-", 30))
		fmt.Print("Select option (1-4): ")

		choice := getUserChoice(4)

		switch choice {
		case 1:
			setExportInterval()
		case 2:
			setExportFormat()
		case 3:
			toggleExport()
		case 4:
			return
		}
	}
}

// setExportInterval allows user to set export interval
func setExportInterval() {
	fmt.Println("\n⏰ Export Interval Settings")
	fmt.Println(strings.Repeat("-", 30))
	fmt.Println("1. 1 minute")
	fmt.Println("2. 1 hour")
	fmt.Println("3. 1 day")
	fmt.Println("4. Custom interval")
	fmt.Println("5. Back to Export Settings")
	fmt.Println(strings.Repeat("-", 30))
	fmt.Print("Select option (1-5): ")

	choice := getUserChoice(5)

	var interval time.Duration
	switch choice {
	case 1:
		interval = 1 * time.Minute
	case 2:
		interval = 1 * time.Hour
	case 3:
		interval = 24 * time.Hour
	case 4:
		fmt.Print("Enter custom interval in minutes: ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		input := strings.TrimSpace(scanner.Text())
		if minutes, err := strconv.Atoi(input); err == nil {
			interval = time.Duration(minutes) * time.Minute
		} else {
			fmt.Println("❌ Invalid input! Using default 1 hour.")
			interval = 1 * time.Hour
		}
	case 5:
		return
	}

	// Update all monitor configurations
	cpuConfig := cpuMonitorManager.GetConfiguration()
	cpuConfig.ExportInterval = interval
	cpuMonitorManager.SetConfiguration(cpuConfig)

	memoryConfig := memoryMonitorManager.GetConfig()
	memoryConfig.ExportInterval = interval
	memoryMonitorManager.UpdateConfig(memoryConfig)

	diskConfig := diskMonitorManager.GetConfig()
	diskConfig.ExportInterval = interval
	diskMonitorManager.UpdateConfig(diskConfig)

	networkConfig := networkMonitorManager.GetConfig()
	networkConfig.ExportInterval = interval
	networkMonitorManager.UpdateConfig(networkConfig)

	processConfig := processMonitorManager.GetConfig()
	processConfig.ExportInterval = interval
	processMonitorManager.UpdateConfig(processConfig)

	fmt.Printf("✅ Export interval set to: %v\n", interval)
	waitForEnter()
}

// setExportFormat allows user to set export format
func setExportFormat() {
	fmt.Println("\n📄 Export Format Settings")
	fmt.Println(strings.Repeat("-", 30))
	fmt.Println("1. JSON")
	fmt.Println("2. CSV")
	fmt.Println("3. TXT")
	fmt.Println("4. Back to Export Settings")
	fmt.Println(strings.Repeat("-", 30))
	fmt.Print("Select option (1-4): ")

	choice := getUserChoice(4)

	var format string
	switch choice {
	case 1:
		format = "json"
	case 2:
		format = "csv"
	case 3:
		format = "txt"
	case 4:
		return
	}

	// Update all monitor configurations
	cpuConfig := cpuMonitorManager.GetConfiguration()
	cpuConfig.ExportFormat = format
	cpuMonitorManager.SetConfiguration(cpuConfig)

	memoryConfig := memoryMonitorManager.GetConfig()
	memoryConfig.ExportFormat = format
	memoryMonitorManager.UpdateConfig(memoryConfig)

	diskConfig := diskMonitorManager.GetConfig()
	diskConfig.ExportFormat = format
	diskMonitorManager.UpdateConfig(diskConfig)

	networkConfig := networkMonitorManager.GetConfig()
	networkConfig.ExportFormat = format
	networkMonitorManager.UpdateConfig(networkConfig)

	processConfig := processMonitorManager.GetConfig()
	processConfig.ExportFormat = format
	processMonitorManager.UpdateConfig(processConfig)

	fmt.Printf("✅ Export format set to: %s\n", strings.ToUpper(format))
	waitForEnter()
}

// toggleExport allows user to enable/disable export
func toggleExport() {
	fmt.Println("\n🔄 Export Toggle")
	fmt.Println(strings.Repeat("-", 30))
	fmt.Println("1. Enable Export")
	fmt.Println("2. Disable Export")
	fmt.Println("3. Back to Export Settings")
	fmt.Println(strings.Repeat("-", 30))
	fmt.Print("Select option (1-3): ")

	choice := getUserChoice(3)

	var enable bool
	switch choice {
	case 1:
		enable = true
	case 2:
		enable = false
	case 3:
		return
	}

	// Update all monitor configurations
	cpuConfig := cpuMonitorManager.GetConfiguration()
	cpuConfig.ExportToFile = enable
	cpuMonitorManager.SetConfiguration(cpuConfig)

	memoryConfig := memoryMonitorManager.GetConfig()
	memoryConfig.ExportToFile = enable
	memoryMonitorManager.UpdateConfig(memoryConfig)

	diskConfig := diskMonitorManager.GetConfig()
	diskConfig.ExportToFile = enable
	diskMonitorManager.UpdateConfig(diskConfig)

	networkConfig := networkMonitorManager.GetConfig()
	networkConfig.ExportToFile = enable
	networkMonitorManager.UpdateConfig(networkConfig)

	processConfig := processMonitorManager.GetConfig()
	processConfig.ExportToFile = enable
	processMonitorManager.UpdateConfig(processConfig)

	status := "disabled"
	if enable {
		status = "enabled"
	}
	fmt.Printf("✅ Export %s\n", status)
	waitForEnter()
}

// showDisplaySettings displays display settings menu
func showDisplaySettings() {
	fmt.Println("Display Settings - Coming soon...")
	waitForEnter()
}

// showMonitoringSettings displays monitoring settings menu
func showMonitoringSettings() {
	fmt.Println("Monitoring Settings - Coming soon...")
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
