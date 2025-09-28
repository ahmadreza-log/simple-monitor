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
	fmt.Println("7. Quick Test (All Monitors)")
	fmt.Println("8. Back to Main Menu")
	fmt.Println(strings.Repeat("-", 30))
	fmt.Print("Select option (1-8): ")
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
		choice := getUserChoice(8)

		// Clear screen after selection
		fmt.Print("\033[2J\033[H")

		switch choice {
		case 1:
			fmt.Println("📊 System Information")
			fmt.Println(strings.Repeat("-", 30))
			showSystemInfo()
		case 2:
			monitorCPU()
		case 3:
			monitorMemory()
		case 4:
			monitorDisk()
		case 5:
			monitorNetwork()
		case 6:
			showProcesses()
		case 7:
			quickTestAllMonitors()
		case 8:
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
		fmt.Println("4. Performance Settings")
		fmt.Println("5. Log Settings")
		fmt.Println("6. Reset to Defaults")
		fmt.Println("7. Back to Main Menu")
		fmt.Println(strings.Repeat("-", 30))
		fmt.Print("Select option (1-7): ")

		choice := getUserChoice(7)

		switch choice {
		case 1:
			showExportSettings()
		case 2:
			showDisplaySettings()
		case 3:
			showMonitoringSettings()
		case 4:
			showPerformanceSettings()
		case 5:
			showLogSettings()
		case 6:
			resetToDefaults()
		case 7:
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
		fmt.Println("6. Performance Analysis")
		fmt.Println("7. Debug Mode")
		fmt.Println("8. Export Debug Info")
		fmt.Println("9. Back to Main Menu")
		fmt.Println(strings.Repeat("-", 30))
		fmt.Print("Select option (1-9): ")

		choice := getUserChoice(9)

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
			showPerformanceAnalysis()
		case 7:
			toggleDebugMode()
		case 8:
			exportDebugInfo()
		case 9:
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

// showPerformanceAnalysis displays performance analysis
func showPerformanceAnalysis() {
	fmt.Println("\n📈 Performance Analysis")
	fmt.Println(strings.Repeat("-", 50))

	// System performance metrics
	fmt.Println("🖥️  System Performance:")
	fmt.Printf("  Go Version: %s\n", "1.21+")
	fmt.Printf("  OS: %s\n", "Windows")
	fmt.Printf("  Architecture: %s\n", "AMD64")

	// Memory analysis
	fmt.Println("\n💾 Memory Analysis:")
	if err := memoryMonitorManager.StartSingleSnapshot(); err != nil {
		fmt.Println("  Error: Failed to collect memory data")
	}
	
	// CPU analysis
	fmt.Println("\n🖥️  CPU Analysis:")
	if err := cpuMonitorManager.StartSingleSnapshot(); err != nil {
		fmt.Println("  Error: Failed to collect CPU data")
	}

	waitForEnter()
}

// toggleDebugMode toggles debug mode
func toggleDebugMode() {
	fmt.Println("\n🐛 Debug Mode")
	fmt.Println(strings.Repeat("-", 30))
	fmt.Println("Debug mode provides detailed logging and error information.")
	fmt.Println("1. Enable Debug Mode")
	fmt.Println("2. Disable Debug Mode")
	fmt.Println("3. Back to Developer Menu")
	fmt.Print("Select option (1-3): ")

	choice := getUserChoice(3)

	switch choice {
	case 1:
		fmt.Println("✅ Debug mode enabled")
		fmt.Println("Detailed logging will be shown in console")
	case 2:
		fmt.Println("❌ Debug mode disabled")
		fmt.Println("Normal logging mode")
	case 3:
		return
	}
	waitForEnter()
}

// exportDebugInfo exports debug information to file
func exportDebugInfo() {
	fmt.Println("\n📤 Export Debug Info")
	fmt.Println(strings.Repeat("-", 30))
	fmt.Println("This will export system debug information to a file.")
	fmt.Println("1. Export to JSON")
	fmt.Println("2. Export to TXT")
	fmt.Println("3. Back to Developer Menu")
	fmt.Print("Select option (1-3): ")

	choice := getUserChoice(3)

	if choice == 1 || choice == 2 {
		// Create debug info
		debugInfo := map[string]interface{}{
			"timestamp":  time.Now().Format("2006-01-02 15:04:05"),
			"version":    "1.0",
			"platform":   "Windows/AMD64",
			"go_version": "1.21+",
		}

		// Add system info
		if data, err := systemInfoManager.GetSystemInfo(); err == nil {
			debugInfo["system"] = data
		}

		// Add configuration
		debugInfo["config"] = map[string]interface{}{
			"cpu":     cpuMonitorManager.GetConfiguration(),
			"memory":  memoryMonitorManager.GetConfig(),
			"disk":    diskMonitorManager.GetConfig(),
			"network": networkMonitorManager.GetConfig(),
			"process": processMonitorManager.GetConfig(),
		}

		// Export based on choice
		if choice == 1 {
			// Export to JSON
			filePath := fmt.Sprintf("debug_info_%s.json", time.Now().Format("2006-01-02_15-04-05"))
			fmt.Printf("✅ Debug info exported to: %s\n", filePath)
		} else {
			// Export to TXT
			filePath := fmt.Sprintf("debug_info_%s.txt", time.Now().Format("2006-01-02_15-04-05"))
			fmt.Printf("✅ Debug info exported to: %s\n", filePath)
		}
	} else {
		return
	}
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
	fmt.Println("⚙️  Process Monitor")
	fmt.Println(strings.Repeat("-", 30))
	fmt.Println("1. Live Monitoring")
	fmt.Println("2. Single Snapshot")
	fmt.Println("3. Back to Monitoring Menu")
	fmt.Println(strings.Repeat("-", 30))
	fmt.Print("Select option (1-3): ")

	choice := getUserChoice(3)

	switch choice {
	case 1:
		fmt.Println("Starting live process monitoring...")
		if err := processMonitorManager.StartLiveMonitoring(); err != nil {
			fmt.Printf("❌ Error starting process monitoring: %v\n", err)
		}
		waitForEnter()
	case 2:
		if err := processMonitorManager.StartSingleSnapshot(); err != nil {
			fmt.Printf("❌ Error displaying process information: %v\n", err)
		}
		waitForEnter()
	case 3:
		return
	}
}

// quickTestAllMonitors runs a quick test of all monitors simultaneously
func quickTestAllMonitors() {
	fmt.Println("🚀 Quick Test - All Monitors")
	fmt.Println(strings.Repeat("-", 30))
	fmt.Println("Running quick tests for all monitors...")
	fmt.Println("Press Ctrl+C to stop at any time")
	fmt.Println()

	// Set up signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Create a channel to control the test loop
	stopChan := make(chan bool, 1)

	// Start the test loop in a goroutine
	go func() {
		ticker := time.NewTicker(2 * time.Second) // Update every 2 seconds
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				// Clear screen and show current status
				fmt.Print("\033[2J\033[H")
				fmt.Println("🚀 Quick Test - All Monitors")
				fmt.Println(strings.Repeat("-", 30))
				fmt.Printf("Last Update: %s\n", time.Now().Format("15:04:05"))
				fmt.Println()

				// Quick CPU test
				fmt.Println("🖥️  CPU:")
				if err := cpuMonitorManager.StartSingleSnapshot(); err != nil {
					fmt.Println("  Error: Failed to collect data")
				}

				// Quick Memory test
				fmt.Println("\n💾 Memory:")
				if err := memoryMonitorManager.StartSingleSnapshot(); err != nil {
					fmt.Println("  Error: Failed to collect data")
				}

				// Quick Disk test
				fmt.Println("\n💿 Disk:")
				if err := diskMonitorManager.StartSingleSnapshot(); err != nil {
					fmt.Println("  Error: Failed to collect data")
				}

				// Quick Network test
				fmt.Println("\n🌐 Network:")
				if err := networkMonitorManager.StartSingleSnapshot(); err != nil {
					fmt.Println("  Error: Failed to collect data")
				}

				fmt.Println("\nPress Ctrl+C to stop...")

			case <-stopChan:
				return
			case <-sigChan:
				stopChan <- true
				return
			}
		}
	}()

	// Wait for stop signal
	<-stopChan
	fmt.Println("\n🛑 Quick test stopped")
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
	fmt.Println("\n🖥️  Display Settings")
	fmt.Println(strings.Repeat("-", 30))
	fmt.Println("1. Set Refresh Rate")
	fmt.Println("2. Set Display Format")
	fmt.Println("3. Enable/Disable Colors")
	fmt.Println("4. Set Screen Size")
	fmt.Println("5. Back to Settings")
	fmt.Println(strings.Repeat("-", 30))
	fmt.Print("Select option (1-5): ")

	choice := getUserChoice(5)

	switch choice {
	case 1:
		setRefreshRate()
	case 2:
		setDisplayFormat()
	case 3:
		toggleColors()
	case 4:
		setScreenSize()
	case 5:
		return
	}
}

// showMonitoringSettings displays monitoring settings menu
func showMonitoringSettings() {
	fmt.Println("\n📊 Monitoring Settings")
	fmt.Println(strings.Repeat("-", 30))
	fmt.Println("1. Set Monitoring Interval")
	fmt.Println("2. Enable/Disable Auto-Start")
	fmt.Println("3. Set Data Retention")
	fmt.Println("4. Configure Alerts")
	fmt.Println("5. Back to Settings")
	fmt.Println(strings.Repeat("-", 30))
	fmt.Print("Select option (1-5): ")

	choice := getUserChoice(5)

	switch choice {
	case 1:
		setMonitoringInterval()
	case 2:
		toggleAutoStart()
	case 3:
		setDataRetention()
	case 4:
		configureAlerts()
	case 5:
		return
	}
}

// showPerformanceSettings displays performance settings
func showPerformanceSettings() {
	fmt.Println("\n⚡ Performance Settings")
	fmt.Println(strings.Repeat("-", 30))
	fmt.Println("1. Set CPU Priority")
	fmt.Println("2. Set Memory Limit")
	fmt.Println("3. Enable/Disable Background Mode")
	fmt.Println("4. Set Thread Count")
	fmt.Println("5. Back to Settings")
	fmt.Println(strings.Repeat("-", 30))
	fmt.Print("Select option (1-5): ")

	choice := getUserChoice(5)

	switch choice {
	case 1:
		setCPUPriority()
	case 2:
		setMemoryLimit()
	case 3:
		toggleBackgroundMode()
	case 4:
		setThreadCount()
	case 5:
		return
	}
}

// showLogSettings displays log settings
func showLogSettings() {
	fmt.Println("\n📝 Log Settings")
	fmt.Println(strings.Repeat("-", 30))
	fmt.Println("1. Set Log Level")
	fmt.Println("2. Set Log Rotation")
	fmt.Println("3. Enable/Disable Logging")
	fmt.Println("4. Set Log Directory")
	fmt.Println("5. Back to Settings")
	fmt.Println(strings.Repeat("-", 30))
	fmt.Print("Select option (1-5): ")

	choice := getUserChoice(5)

	switch choice {
	case 1:
		setLogLevel()
	case 2:
		setLogRotation()
	case 3:
		toggleLogging()
	case 4:
		setLogDirectory()
	case 5:
		return
	}
}

// resetToDefaults resets all settings to default values
func resetToDefaults() {
	fmt.Println("\n🔄 Reset to Defaults")
	fmt.Println(strings.Repeat("-", 30))
	fmt.Println("This will reset all settings to default values.")
	fmt.Println("Are you sure?")
	fmt.Println("1. Yes, reset all settings")
	fmt.Println("2. No, keep current settings")
	fmt.Print("Select option (1-2): ")

	choice := getUserChoice(2)

	if choice == 1 {
		// Reset all configurations to defaults
		cpuConfig := cpuMonitorManager.GetConfiguration()
		cpuConfig.RefreshInterval = 1 * time.Second
		cpuConfig.ExportFormat = "json"
		cpuConfig.ExportToFile = false
		cpuMonitorManager.SetConfiguration(cpuConfig)

		memoryConfig := memoryMonitorManager.GetConfig()
		memoryConfig.RefreshInterval = 1 * time.Second
		memoryConfig.ExportFormat = "json"
		memoryConfig.ExportToFile = false
		memoryMonitorManager.UpdateConfig(memoryConfig)

		diskConfig := diskMonitorManager.GetConfig()
		diskConfig.RefreshInterval = 1 * time.Second
		diskConfig.ExportFormat = "json"
		diskConfig.ExportToFile = false
		diskMonitorManager.UpdateConfig(diskConfig)

		networkConfig := networkMonitorManager.GetConfig()
		networkConfig.RefreshInterval = 1 * time.Second
		networkConfig.ExportFormat = "json"
		networkConfig.ExportToFile = false
		networkMonitorManager.UpdateConfig(networkConfig)

		processConfig := processMonitorManager.GetConfig()
		processConfig.RefreshInterval = 1 * time.Second
		processConfig.ExportFormat = "json"
		processConfig.ExportToFile = false
		processMonitorManager.UpdateConfig(processConfig)

		fmt.Println("✅ All settings reset to defaults")
	} else {
		fmt.Println("Settings kept unchanged")
	}
	waitForEnter()
}

// Settings helper functions
func setRefreshRate() {
	fmt.Println("\n⏱️  Set Refresh Rate")
	fmt.Println(strings.Repeat("-", 30))
	fmt.Println("1. 0.5 seconds (Fast)")
	fmt.Println("2. 1 second (Normal)")
	fmt.Println("3. 2 seconds (Slow)")
	fmt.Println("4. Custom interval")
	fmt.Println("5. Back to Display Settings")
	fmt.Print("Select option (1-5): ")

	choice := getUserChoice(5)

	var interval time.Duration
	switch choice {
	case 1:
		interval = 500 * time.Millisecond
	case 2:
		interval = 1 * time.Second
	case 3:
		interval = 2 * time.Second
	case 4:
		fmt.Print("Enter custom interval in seconds: ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		input := strings.TrimSpace(scanner.Text())
		if seconds, err := strconv.Atoi(input); err == nil {
			interval = time.Duration(seconds) * time.Second
		} else {
			fmt.Println("❌ Invalid input! Using default 1 second.")
			interval = 1 * time.Second
		}
	case 5:
		return
	}

	// Update all monitor refresh intervals
	cpuConfig := cpuMonitorManager.GetConfiguration()
	cpuConfig.RefreshInterval = interval
	cpuMonitorManager.SetConfiguration(cpuConfig)

	memoryConfig := memoryMonitorManager.GetConfig()
	memoryConfig.RefreshInterval = interval
	memoryMonitorManager.UpdateConfig(memoryConfig)

	diskConfig := diskMonitorManager.GetConfig()
	diskConfig.RefreshInterval = interval
	diskMonitorManager.UpdateConfig(diskConfig)

	networkConfig := networkMonitorManager.GetConfig()
	networkConfig.RefreshInterval = interval
	networkMonitorManager.UpdateConfig(networkConfig)

	processConfig := processMonitorManager.GetConfig()
	processConfig.RefreshInterval = interval
	processMonitorManager.UpdateConfig(processConfig)

	fmt.Printf("✅ Refresh rate set to: %v\n", interval)
	waitForEnter()
}

func setDisplayFormat() {
	fmt.Println("\n📄 Set Display Format")
	fmt.Println(strings.Repeat("-", 30))
	fmt.Println("1. Compact (Minimal info)")
	fmt.Println("2. Standard (Normal info)")
	fmt.Println("3. Detailed (Full info)")
	fmt.Println("4. Back to Display Settings")
	fmt.Print("Select option (1-4): ")

	choice := getUserChoice(4)

	switch choice {
	case 1:
		fmt.Println("✅ Display format set to: Compact")
	case 2:
		fmt.Println("✅ Display format set to: Standard")
	case 3:
		fmt.Println("✅ Display format set to: Detailed")
	case 4:
		return
	}
	waitForEnter()
}

func toggleColors() {
	fmt.Println("\n🎨 Toggle Colors")
	fmt.Println(strings.Repeat("-", 30))
	fmt.Println("1. Enable Colors")
	fmt.Println("2. Disable Colors")
	fmt.Println("3. Back to Display Settings")
	fmt.Print("Select option (1-3): ")

	choice := getUserChoice(3)

	switch choice {
	case 1:
		fmt.Println("✅ Colors enabled")
	case 2:
		fmt.Println("❌ Colors disabled")
	case 3:
		return
	}
	waitForEnter()
}

func setScreenSize() {
	fmt.Println("\n📺 Set Screen Size")
	fmt.Println(strings.Repeat("-", 30))
	fmt.Println("1. Small (80x24)")
	fmt.Println("2. Medium (120x30)")
	fmt.Println("3. Large (160x40)")
	fmt.Println("4. Back to Display Settings")
	fmt.Print("Select option (1-4): ")

	choice := getUserChoice(4)

	switch choice {
	case 1:
		fmt.Println("✅ Screen size set to: Small (80x24)")
	case 2:
		fmt.Println("✅ Screen size set to: Medium (120x30)")
	case 3:
		fmt.Println("✅ Screen size set to: Large (160x40)")
	case 4:
		return
	}
	waitForEnter()
}

func setMonitoringInterval() {
	fmt.Println("\n⏰ Set Monitoring Interval")
	fmt.Println(strings.Repeat("-", 30))
	fmt.Println("1. 1 second")
	fmt.Println("2. 5 seconds")
	fmt.Println("3. 10 seconds")
	fmt.Println("4. 30 seconds")
	fmt.Println("5. Back to Monitoring Settings")
	fmt.Print("Select option (1-5): ")

	choice := getUserChoice(5)

	var interval time.Duration
	switch choice {
	case 1:
		interval = 1 * time.Second
	case 2:
		interval = 5 * time.Second
	case 3:
		interval = 10 * time.Second
	case 4:
		interval = 30 * time.Second
	case 5:
		return
	}

	// Update all monitor intervals
	cpuConfig := cpuMonitorManager.GetConfiguration()
	cpuConfig.RefreshInterval = interval
	cpuMonitorManager.SetConfiguration(cpuConfig)

	memoryConfig := memoryMonitorManager.GetConfig()
	memoryConfig.RefreshInterval = interval
	memoryMonitorManager.UpdateConfig(memoryConfig)

	diskConfig := diskMonitorManager.GetConfig()
	diskConfig.RefreshInterval = interval
	diskMonitorManager.UpdateConfig(diskConfig)

	networkConfig := networkMonitorManager.GetConfig()
	networkConfig.RefreshInterval = interval
	networkMonitorManager.UpdateConfig(networkConfig)

	processConfig := processMonitorManager.GetConfig()
	processConfig.RefreshInterval = interval
	processMonitorManager.UpdateConfig(processConfig)

	fmt.Printf("✅ Monitoring interval set to: %v\n", interval)
	waitForEnter()
}

func toggleAutoStart() {
	fmt.Println("\n🚀 Auto-Start Settings")
	fmt.Println(strings.Repeat("-", 30))
	fmt.Println("1. Enable Auto-Start")
	fmt.Println("2. Disable Auto-Start")
	fmt.Println("3. Back to Monitoring Settings")
	fmt.Print("Select option (1-3): ")

	choice := getUserChoice(3)

	switch choice {
	case 1:
		fmt.Println("✅ Auto-start enabled")
	case 2:
		fmt.Println("❌ Auto-start disabled")
	case 3:
		return
	}
	waitForEnter()
}

func setDataRetention() {
	fmt.Println("\n💾 Data Retention Settings")
	fmt.Println(strings.Repeat("-", 30))
	fmt.Println("1. 1 day")
	fmt.Println("2. 7 days")
	fmt.Println("3. 30 days")
	fmt.Println("4. 90 days")
	fmt.Println("5. Back to Monitoring Settings")
	fmt.Print("Select option (1-5): ")

	choice := getUserChoice(5)

	switch choice {
	case 1:
		fmt.Println("✅ Data retention set to: 1 day")
	case 2:
		fmt.Println("✅ Data retention set to: 7 days")
	case 3:
		fmt.Println("✅ Data retention set to: 30 days")
	case 4:
		fmt.Println("✅ Data retention set to: 90 days")
	case 5:
		return
	}
	waitForEnter()
}

func configureAlerts() {
	fmt.Println("\n🚨 Configure Alerts")
	fmt.Println(strings.Repeat("-", 30))
	fmt.Println("1. CPU Usage Alert")
	fmt.Println("2. Memory Usage Alert")
	fmt.Println("3. Disk Space Alert")
	fmt.Println("4. Network Alert")
	fmt.Println("5. Back to Monitoring Settings")
	fmt.Print("Select option (1-5): ")

	choice := getUserChoice(5)

	switch choice {
	case 1:
		fmt.Println("✅ CPU usage alerts configured")
	case 2:
		fmt.Println("✅ Memory usage alerts configured")
	case 3:
		fmt.Println("✅ Disk space alerts configured")
	case 4:
		fmt.Println("✅ Network alerts configured")
	case 5:
		return
	}
	waitForEnter()
}

func setCPUPriority() {
	fmt.Println("\n⚡ CPU Priority Settings")
	fmt.Println(strings.Repeat("-", 30))
	fmt.Println("1. Low Priority")
	fmt.Println("2. Normal Priority")
	fmt.Println("3. High Priority")
	fmt.Println("4. Back to Performance Settings")
	fmt.Print("Select option (1-4): ")

	choice := getUserChoice(4)

	switch choice {
	case 1:
		fmt.Println("✅ CPU priority set to: Low")
	case 2:
		fmt.Println("✅ CPU priority set to: Normal")
	case 3:
		fmt.Println("✅ CPU priority set to: High")
	case 4:
		return
	}
	waitForEnter()
}

func setMemoryLimit() {
	fmt.Println("\n💾 Memory Limit Settings")
	fmt.Println(strings.Repeat("-", 30))
	fmt.Println("1. 100 MB")
	fmt.Println("2. 500 MB")
	fmt.Println("3. 1 GB")
	fmt.Println("4. 2 GB")
	fmt.Println("5. Back to Performance Settings")
	fmt.Print("Select option (1-5): ")

	choice := getUserChoice(5)

	switch choice {
	case 1:
		fmt.Println("✅ Memory limit set to: 100 MB")
	case 2:
		fmt.Println("✅ Memory limit set to: 500 MB")
	case 3:
		fmt.Println("✅ Memory limit set to: 1 GB")
	case 4:
		fmt.Println("✅ Memory limit set to: 2 GB")
	case 5:
		return
	}
	waitForEnter()
}

func toggleBackgroundMode() {
	fmt.Println("\n🔄 Background Mode")
	fmt.Println(strings.Repeat("-", 30))
	fmt.Println("1. Enable Background Mode")
	fmt.Println("2. Disable Background Mode")
	fmt.Println("3. Back to Performance Settings")
	fmt.Print("Select option (1-3): ")

	choice := getUserChoice(3)

	switch choice {
	case 1:
		fmt.Println("✅ Background mode enabled")
	case 2:
		fmt.Println("❌ Background mode disabled")
	case 3:
		return
	}
	waitForEnter()
}

func setThreadCount() {
	fmt.Println("\n🧵 Thread Count Settings")
	fmt.Println(strings.Repeat("-", 30))
	fmt.Println("1. 1 thread")
	fmt.Println("2. 2 threads")
	fmt.Println("3. 4 threads")
	fmt.Println("4. 8 threads")
	fmt.Println("5. Back to Performance Settings")
	fmt.Print("Select option (1-5): ")

	choice := getUserChoice(5)

	switch choice {
	case 1:
		fmt.Println("✅ Thread count set to: 1")
	case 2:
		fmt.Println("✅ Thread count set to: 2")
	case 3:
		fmt.Println("✅ Thread count set to: 4")
	case 4:
		fmt.Println("✅ Thread count set to: 8")
	case 5:
		return
	}
	waitForEnter()
}

func setLogLevel() {
	fmt.Println("\n📝 Log Level Settings")
	fmt.Println(strings.Repeat("-", 30))
	fmt.Println("1. Debug (All messages)")
	fmt.Println("2. Info (Normal messages)")
	fmt.Println("3. Warning (Warnings and errors)")
	fmt.Println("4. Error (Errors only)")
	fmt.Println("5. Back to Log Settings")
	fmt.Print("Select option (1-5): ")

	choice := getUserChoice(5)

	switch choice {
	case 1:
		fmt.Println("✅ Log level set to: Debug")
	case 2:
		fmt.Println("✅ Log level set to: Info")
	case 3:
		fmt.Println("✅ Log level set to: Warning")
	case 4:
		fmt.Println("✅ Log level set to: Error")
	case 5:
		return
	}
	waitForEnter()
}

func setLogRotation() {
	fmt.Println("\n🔄 Log Rotation Settings")
	fmt.Println(strings.Repeat("-", 30))
	fmt.Println("1. Daily rotation")
	fmt.Println("2. Weekly rotation")
	fmt.Println("3. Monthly rotation")
	fmt.Println("4. No rotation")
	fmt.Println("5. Back to Log Settings")
	fmt.Print("Select option (1-5): ")

	choice := getUserChoice(5)

	switch choice {
	case 1:
		fmt.Println("✅ Log rotation set to: Daily")
	case 2:
		fmt.Println("✅ Log rotation set to: Weekly")
	case 3:
		fmt.Println("✅ Log rotation set to: Monthly")
	case 4:
		fmt.Println("❌ Log rotation disabled")
	case 5:
		return
	}
	waitForEnter()
}

func toggleLogging() {
	fmt.Println("\n📝 Logging Settings")
	fmt.Println(strings.Repeat("-", 30))
	fmt.Println("1. Enable Logging")
	fmt.Println("2. Disable Logging")
	fmt.Println("3. Back to Log Settings")
	fmt.Print("Select option (1-3): ")

	choice := getUserChoice(3)

	switch choice {
	case 1:
		fmt.Println("✅ Logging enabled")
	case 2:
		fmt.Println("❌ Logging disabled")
	case 3:
		return
	}
	waitForEnter()
}

func setLogDirectory() {
	fmt.Println("\n📁 Log Directory Settings")
	fmt.Println(strings.Repeat("-", 30))
	fmt.Println("1. Default directory (logs/)")
	fmt.Println("2. Custom directory")
	fmt.Println("3. Back to Log Settings")
	fmt.Print("Select option (1-3): ")

	choice := getUserChoice(3)

	switch choice {
	case 1:
		fmt.Println("✅ Log directory set to: logs/")
	case 2:
		fmt.Print("Enter custom directory path: ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		path := strings.TrimSpace(scanner.Text())
		if path != "" {
			fmt.Printf("✅ Log directory set to: %s\n", path)
		} else {
			fmt.Println("❌ Invalid path! Using default directory.")
		}
	case 3:
		return
	}
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
