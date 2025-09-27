package main

import (
	"bufio"
	"fmt"
	"os"
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
	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println("           🖥️  Simple Monitor v1.0")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Println("1. Start Monitoring")
	fmt.Println("2. Settings")
	fmt.Println("3. Developer")
	fmt.Println("4. Quit")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Print("Please select an option (1-4): ")
}

// displayMonitoringMenu shows the monitoring submenu
func displayMonitoringMenu() {
	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println("           📊 Monitoring Options")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Println("1. System Information")
	fmt.Println("2. CPU Monitor")
	fmt.Println("3. Memory Monitor")
	fmt.Println("4. Disk Monitor")
	fmt.Println("5. Network Monitor")
	fmt.Println("6. Process Monitor")
	fmt.Println("7. Back to Main Menu")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Print("Please select a monitoring option (1-7): ")
}

// getUserChoice gets user input and validates it for main menu
func getUserChoice(maxOptions int) int {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		scanner.Scan()
		input := strings.TrimSpace(scanner.Text())

		choice, err := strconv.Atoi(input)
		if err != nil {
			fmt.Printf("❌ Invalid input! Please enter a number between 1 and %d: ", maxOptions)
			continue
		}

		if choice < 1 || choice > maxOptions {
			fmt.Printf("❌ Invalid option! Please enter a number between 1 and %d: ", maxOptions)
			continue
		}

		return choice
	}
}

// handleMainMenuChoice processes the user's main menu choice
func handleMainMenuChoice(choice int) {
	switch choice {
	case 1:
		fmt.Println("\n🚀 Starting monitoring...")
		startMonitoring()
	case 2:
		fmt.Println("\n⚙️  Opening settings...")
		showSettings()
	case 3:
		fmt.Println("\n👨‍💻 Developer options...")
		showDeveloper()
	case 4:
		fmt.Println("\n👋 Goodbye! Thank you for using Simple Monitor.")
		os.Exit(0)
	}
}

// startMonitoring handles the monitoring submenu
func startMonitoring() {
	for {
		displayMonitoringMenu()
		choice := getUserChoice(7)

		switch choice {
		case 1:
			fmt.Println("\n📊 Displaying system information...")
			showSystemInfo()
		case 2:
			fmt.Println("\n🖥️  Starting CPU monitoring...")
			monitorCPU()
		case 3:
			fmt.Println("\n💾 Starting memory monitoring...")
			monitorMemory()
		case 4:
			fmt.Println("\n💿 Starting disk monitoring...")
			monitorDisk()
		case 5:
			fmt.Println("\n🌐 Starting network monitoring...")
			monitorNetwork()
		case 6:
			fmt.Println("\n⚙️  Displaying processes...")
			showProcesses()
		case 7:
			fmt.Println("\n⬅️  Returning to main menu...")
			return
		}
	}
}

// showSettings displays settings menu
func showSettings() {
	fmt.Println("🔧 Settings will be implemented in the next version...")
	waitForEnter()
}

// showDeveloper displays developer options
func showDeveloper() {
	fmt.Println("👨‍💻 Developer options will be implemented in the next version...")
	waitForEnter()
}

// Placeholder functions for different monitoring features
func showSystemInfo() {
	fmt.Println("🔧 System information will be implemented in the next version...")
	waitForEnter()
}

func monitorCPU() {
	fmt.Println("🖥️  CPU monitoring will be implemented in the next version...")
	waitForEnter()
}

func monitorMemory() {
	fmt.Println("💾 Memory monitoring will be implemented in the next version...")
	waitForEnter()
}

func monitorDisk() {
	fmt.Println("💿 Disk monitoring will be implemented in the next version...")
	waitForEnter()
}

func monitorNetwork() {
	fmt.Println("🌐 Network monitoring will be implemented in the next version...")
	waitForEnter()
}

func showProcesses() {
	fmt.Println("⚙️  Process monitoring will be implemented in the next version...")
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
