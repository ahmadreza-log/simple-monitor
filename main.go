package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"simple-monitor/cpumonitor"
	"simple-monitor/diskmonitor"
	"simple-monitor/memorymonitor"
	"simple-monitor/networkmonitor"
	"simple-monitor/processmonitor"
	"simple-monitor/systeminfo"
	"strconv"
	"strings"
	"syscall"
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

	// Set up signal handling for this function
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	for {
		// Use a goroutine to handle input and signals concurrently
		inputChan := make(chan string, 1)
		
		go func() {
			scanner.Scan()
			inputChan <- scanner.Text()
		}()
		
		select {
		case input := <-inputChan:
			input = strings.TrimSpace(input)
			
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
			
		case <-sigChan:
			fmt.Println("\n\n👋 Goodbye! Thank you for using Simple Monitor.")
			os.Exit(0)
		}
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
	for {
		fmt.Println("\n👨‍💻 Developer Section")
		fmt.Println(strings.Repeat("-", 30))
		fmt.Println("1. View System Information")
		fmt.Println("2. View Log Files")
		fmt.Println("3. Clear Log Files")
		fmt.Println("4. View Configuration")
		fmt.Println("5. Test All Monitors")
		fmt.Println("6. Back to Main Menu")
		fmt.Println(strings.Repeat("-", 30))
		fmt.Print("Select option (1-6): ")

		choice := getUserChoice(6)

		switch choice {
		case 1:
			showDeveloperSystemInfo()
		case 2:
			showLogFiles()
		case 3:
			clearLogFiles()
		case 4:
			showConfiguration()
		case 5:
			testAllMonitors()
		case 6:
			return
		}
	}
}

// showDeveloperSystemInfo displays detailed system information for developers
func showDeveloperSystemInfo() {
	fmt.Println("\n🔧 Developer System Information")
	fmt.Println(strings.Repeat("-", 50))
	
	// Show system info
	if err := systemInfoManager.ShowSystemInfo(); err != nil {
		fmt.Printf("❌ Error displaying system information: %v\n", err)
	}
	
	// Show Go version and build info
	fmt.Println("\n📋 Build Information:")
	fmt.Printf("Go Version: %s\n", "1.21+")
	fmt.Printf("Build Time: %s\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Printf("Platform: %s\n", "Windows/AMD64")
	
	waitForEnter()
}

// showLogFiles displays available log files
func showLogFiles() {
	fmt.Println("\n📁 Log Files")
	fmt.Println(strings.Repeat("-", 30))
	
	// Check logs directory
	logDir := "logs"
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		fmt.Println("No log files found.")
		waitForEnter()
		return
	}
	
	// List log files
	files, err := os.ReadDir(logDir)
	if err != nil {
		fmt.Printf("❌ Error reading log directory: %v\n", err)
		waitForEnter()
		return
	}
	
	if len(files) == 0 {
		fmt.Println("No log files found.")
		waitForEnter()
		return
	}
	
	fmt.Println("Available log files:")
	for i, file := range files {
		if !file.IsDir() {
			fmt.Printf("%d. %s\n", i+1, file.Name())
		}
	}
	
	waitForEnter()
}

// clearLogFiles clears all log files
func clearLogFiles() {
	fmt.Println("\n🗑️  Clear Log Files")
	fmt.Println(strings.Repeat("-", 30))
	fmt.Println("This will delete all log files. Are you sure?")
	fmt.Println("1. Yes, delete all logs")
	fmt.Println("2. No, keep logs")
	fmt.Print("Select option (1-2): ")
	
	choice := getUserChoice(2)
	
	if choice == 1 {
		logDir := "logs"
		if _, err := os.Stat(logDir); os.IsNotExist(err) {
			fmt.Println("No log directory found.")
			waitForEnter()
			return
		}
		
		files, err := os.ReadDir(logDir)
		if err != nil {
			fmt.Printf("❌ Error reading log directory: %v\n", err)
			waitForEnter()
			return
		}
		
		deletedCount := 0
		for _, file := range files {
			if !file.IsDir() {
				if err := os.Remove(logDir + "/" + file.Name()); err == nil {
					deletedCount++
				}
			}
		}
		
		fmt.Printf("✅ Deleted %d log files.\n", deletedCount)
	} else {
		fmt.Println("Log files kept.")
	}
	
	waitForEnter()
}

// showConfiguration displays current configuration
func showConfiguration() {
	fmt.Println("\n⚙️  Current Configuration")
	fmt.Println(strings.Repeat("-", 50))
	
	// CPU Monitor Config
	cpuConfig := cpuMonitorManager.GetConfiguration()
	fmt.Println("🖥️  CPU Monitor:")
	fmt.Printf("  - Refresh Interval: %v\n", cpuConfig.RefreshInterval)
	fmt.Printf("  - Export Format: %s\n", cpuConfig.ExportFormat)
	fmt.Printf("  - Export Enabled: %t\n", cpuConfig.ExportToFile)
	
	// Memory Monitor Config
	memoryConfig := memoryMonitorManager.GetConfig()
	fmt.Println("\n💾 Memory Monitor:")
	fmt.Printf("  - Refresh Interval: %v\n", memoryConfig.RefreshInterval)
	fmt.Printf("  - Export Format: %s\n", memoryConfig.ExportFormat)
	fmt.Printf("  - Export Enabled: %t\n", memoryConfig.ExportToFile)
	
	// Disk Monitor Config
	diskConfig := diskMonitorManager.GetConfig()
	fmt.Println("\n💿 Disk Monitor:")
	fmt.Printf("  - Refresh Interval: %v\n", diskConfig.RefreshInterval)
	fmt.Printf("  - Export Format: %s\n", diskConfig.ExportFormat)
	fmt.Printf("  - Export Enabled: %t\n", diskConfig.ExportToFile)
	
	// Network Monitor Config
	networkConfig := networkMonitorManager.GetConfig()
	fmt.Println("\n🌐 Network Monitor:")
	fmt.Printf("  - Refresh Interval: %v\n", networkConfig.RefreshInterval)
	fmt.Printf("  - Export Format: %s\n", networkConfig.ExportFormat)
	fmt.Printf("  - Export Enabled: %t\n", networkConfig.ExportToFile)
	
	// Process Monitor Config
	processConfig := processMonitorManager.GetConfig()
	fmt.Println("\n⚙️  Process Monitor:")
	fmt.Printf("  - Refresh Interval: %v\n", processConfig.RefreshInterval)
	fmt.Printf("  - Export Format: %s\n", processConfig.ExportFormat)
	fmt.Printf("  - Export Enabled: %t\n", processConfig.ExportToFile)
	
	waitForEnter()
}

// testAllMonitors tests all monitor components
func testAllMonitors() {
	fmt.Println("\n🧪 Testing All Monitors")
	fmt.Println(strings.Repeat("-", 30))
	
	// Test System Info
	fmt.Println("Testing System Info...")
	if err := systemInfoManager.ShowSystemInfo(); err != nil {
		fmt.Printf("❌ System Info Error: %v\n", err)
	} else {
		fmt.Println("✅ System Info: OK")
	}
	
	// Test CPU Monitor
	fmt.Println("\nTesting CPU Monitor...")
	if err := cpuMonitorManager.StartSingleSnapshot(); err != nil {
		fmt.Printf("❌ CPU Monitor Error: %v\n", err)
	} else {
		fmt.Println("✅ CPU Monitor: OK")
	}
	
	// Test Memory Monitor
	fmt.Println("\nTesting Memory Monitor...")
	if err := memoryMonitorManager.StartSingleSnapshot(); err != nil {
		fmt.Printf("❌ Memory Monitor Error: %v\n", err)
	} else {
		fmt.Println("✅ Memory Monitor: OK")
	}
	
	// Test Disk Monitor
	fmt.Println("\nTesting Disk Monitor...")
	if err := diskMonitorManager.StartSingleSnapshot(); err != nil {
		fmt.Printf("❌ Disk Monitor Error: %v\n", err)
	} else {
		fmt.Println("✅ Disk Monitor: OK")
	}
	
	// Test Network Monitor
	fmt.Println("\nTesting Network Monitor...")
	if err := networkMonitorManager.StartSingleSnapshot(); err != nil {
		fmt.Printf("❌ Network Monitor Error: %v\n", err)
	} else {
		fmt.Println("✅ Network Monitor: OK")
	}
	
	fmt.Println("\n🎉 All tests completed!")
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
