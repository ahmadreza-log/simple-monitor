package systeminfo

import (
	"fmt"
	"strings"
	"time"
)

// SystemInfoDisplayer handles the display and formatting of system information
// This struct provides methods to format and display system data in a user-friendly way
type SystemInfoDisplayer struct {
	// Display configuration options
	ShowDetailedInfo bool   // Whether to show detailed information
	UseColors        bool   // Whether to use colored output
	DateFormat       string // Date format for timestamps
}

// NewSystemInfoDisplayer creates a new instance of SystemInfoDisplayer
// with default configuration values
func NewSystemInfoDisplayer() *SystemInfoDisplayer {
	return &SystemInfoDisplayer{
		ShowDetailedInfo: true,
		UseColors:        true,
		DateFormat:       "2006-01-02 15:04:05",
	}
}

// DisplaySystemInfo displays comprehensive system information in a formatted way
// This is the main method that formats and displays all system information
func (displayer *SystemInfoDisplayer) DisplaySystemInfo(systemInfo *SystemInfo) {
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("                    🖥️  SYSTEM INFORMATION")
	fmt.Println(strings.Repeat("=", 80))

	// Display basic system information
	displayer.displayBasicInfo(systemInfo)

	// Display CPU information
	displayer.displayCPUInfo(&systemInfo.CPUInfo)

	// Display memory information
	displayer.displayMemoryInfo(&systemInfo.MemoryInfo)

	// Display disk information
	displayer.displayDiskInfo(systemInfo.DiskInfo)

	// Display network information
	displayer.displayNetworkInfo(systemInfo.NetworkInfo)

	// Display performance metrics
	displayer.displayPerformanceMetrics(systemInfo)

	fmt.Println(strings.Repeat("=", 80))
	fmt.Printf("📅 Last Updated: %s\n", systemInfo.Timestamp.Format(displayer.DateFormat))
	fmt.Println(strings.Repeat("=", 80))
}

// displayBasicInfo displays basic system identification information
func (displayer *SystemInfoDisplayer) displayBasicInfo(systemInfo *SystemInfo) {
	fmt.Println("\n🔧 BASIC SYSTEM INFORMATION")
	fmt.Println(strings.Repeat("-", 50))

	fmt.Printf("Hostname:        %s\n", displayer.formatValue(systemInfo.HostName, "Unknown"))
	fmt.Printf("Operating System: %s\n", displayer.formatValue(systemInfo.OperatingSystem, "Unknown"))
	fmt.Printf("Architecture:    %s\n", displayer.formatValue(systemInfo.Architecture, "Unknown"))
	fmt.Printf("Kernel Version:  %s\n", displayer.formatValue(systemInfo.KernelVersion, "Unknown"))

	// Display uptime information
	if systemInfo.Uptime > 0 {
		fmt.Printf("System Uptime:   %s\n", displayer.formatDuration(systemInfo.Uptime))
	}

	// Display boot time
	if !systemInfo.BootTime.IsZero() {
		fmt.Printf("Boot Time:       %s\n", systemInfo.BootTime.Format(displayer.DateFormat))
	}
}

// displayCPUInfo displays detailed CPU information and usage statistics
func (displayer *SystemInfoDisplayer) displayCPUInfo(cpuInfo *CPUInfo) {
	fmt.Println("\n🖥️  CPU INFORMATION")
	fmt.Println(strings.Repeat("-", 50))

	fmt.Printf("Model:           %s\n", displayer.formatValue(cpuInfo.ModelName, "Unknown"))
	fmt.Printf("Vendor:          %s\n", displayer.formatValue(cpuInfo.VendorID, "Unknown"))
	fmt.Printf("Family:          %s\n", displayer.formatValue(cpuInfo.CPUFamily, "Unknown"))
	fmt.Printf("Frequency:       %.2f MHz\n", cpuInfo.CPUMHz)

	fmt.Printf("Physical Cores:  %d\n", cpuInfo.PhysicalCores)
	fmt.Printf("Logical Cores:   %d\n", cpuInfo.LogicalCores)

	fmt.Println("\n📊 CPU USAGE")
	fmt.Println(strings.Repeat("-", 30))
	fmt.Printf("Overall Usage:   %.2f%%\n", cpuInfo.UsagePercent)
	fmt.Printf("User Processes:  %.2f%%\n", cpuInfo.UserPercent)
	fmt.Printf("System Processes: %.2f%%\n", cpuInfo.SystemPercent)
	fmt.Printf("Idle:            %.2f%%\n", cpuInfo.IdlePercent)

	// Display temperature if available
	if cpuInfo.Temperature > 0 {
		fmt.Printf("Temperature:     %.1f°C\n", cpuInfo.Temperature)
	}
}

// displayMemoryInfo displays memory usage and statistics
func (displayer *SystemInfoDisplayer) displayMemoryInfo(memoryInfo *MemoryInfo) {
	fmt.Println("\n💾 MEMORY INFORMATION")
	fmt.Println(strings.Repeat("-", 50))

	// Display physical memory
	fmt.Printf("Total Memory:    %s\n", displayer.formatBytes(memoryInfo.TotalMemory))
	fmt.Printf("Used Memory:     %s (%.2f%%)\n",
		displayer.formatBytes(memoryInfo.UsedMemory),
		memoryInfo.MemoryUsagePercent)
	fmt.Printf("Available Memory: %s\n", displayer.formatBytes(memoryInfo.AvailableMemory))
	fmt.Printf("Free Memory:     %s\n", displayer.formatBytes(memoryInfo.FreeMemory))

	// Display swap information
	if memoryInfo.TotalSwap > 0 {
		fmt.Println("\n🔄 SWAP INFORMATION")
		fmt.Println(strings.Repeat("-", 30))
		fmt.Printf("Total Swap:      %s\n", displayer.formatBytes(memoryInfo.TotalSwap))
		fmt.Printf("Used Swap:       %s\n", displayer.formatBytes(memoryInfo.UsedSwap))
		fmt.Printf("Free Swap:       %s\n", displayer.formatBytes(memoryInfo.FreeSwap))
	}

	// Display cache and buffer information
	if displayer.ShowDetailedInfo {
		fmt.Println("\n📋 MEMORY DETAILS")
		fmt.Println(strings.Repeat("-", 30))
		fmt.Printf("Cache Size:      %s\n", displayer.formatBytes(memoryInfo.CacheSize))
		fmt.Printf("Buffer Size:     %s\n", displayer.formatBytes(memoryInfo.BufferSize))
	}
}

// displayDiskInfo displays information about all disk drives
func (displayer *SystemInfoDisplayer) displayDiskInfo(diskInfo []DiskInfo) {
	if len(diskInfo) == 0 {
		fmt.Println("\n💿 DISK INFORMATION")
		fmt.Println(strings.Repeat("-", 50))
		fmt.Println("No disk information available")
		return
	}

	fmt.Println("\n💿 DISK INFORMATION")
	fmt.Println(strings.Repeat("-", 50))

	for i, disk := range diskInfo {
		fmt.Printf("\n📀 Disk %d: %s\n", i+1, disk.DeviceName)
		fmt.Println(strings.Repeat("-", 30))
		fmt.Printf("Mount Point:     %s\n", disk.MountPoint)
		fmt.Printf("File System:     %s\n", disk.FileSystem)
		fmt.Printf("Total Size:      %s\n", displayer.formatBytes(disk.TotalSize))
		fmt.Printf("Used Size:       %s (%.2f%%)\n",
			displayer.formatBytes(disk.UsedSize),
			disk.UsagePercent)
		fmt.Printf("Free Size:       %s\n", displayer.formatBytes(disk.FreeSize))

		// Display disk type and properties
		diskType := "HDD"
		if disk.IsSSD {
			diskType = "SSD"
		}
		fmt.Printf("Disk Type:       %s\n", diskType)

		if disk.IsRemovable {
			fmt.Printf("Removable:       Yes\n")
		}

		// Display performance metrics if available
		if disk.ReadSpeed > 0 || disk.WriteSpeed > 0 {
			fmt.Printf("Read Speed:      %s/s\n", displayer.formatBytes(disk.ReadSpeed))
			fmt.Printf("Write Speed:     %s/s\n", displayer.formatBytes(disk.WriteSpeed))
		}
	}
}

// displayNetworkInfo displays information about network interfaces
func (displayer *SystemInfoDisplayer) displayNetworkInfo(networkInfo []NetworkInfo) {
	if len(networkInfo) == 0 {
		fmt.Println("\n🌐 NETWORK INFORMATION")
		fmt.Println(strings.Repeat("-", 50))
		fmt.Println("No network information available")
		return
	}

	fmt.Println("\n🌐 NETWORK INFORMATION")
	fmt.Println(strings.Repeat("-", 50))

	for i, network := range networkInfo {
		fmt.Printf("\n🔌 Interface %d: %s\n", i+1, network.InterfaceName)
		fmt.Println(strings.Repeat("-", 30))
		fmt.Printf("Type:            %s\n", network.InterfaceType)
		fmt.Printf("IP Address:      %s\n", displayer.formatValue(network.IPAddress, "Not assigned"))
		fmt.Printf("Subnet Mask:     %s\n", displayer.formatValue(network.SubnetMask, "Not assigned"))
		fmt.Printf("Gateway:         %s\n", displayer.formatValue(network.Gateway, "Not assigned"))

		// Display status
		status := "Down"
		if network.IsUp {
			status = "Up"
		}
		fmt.Printf("Status:          %s\n", status)

		if network.IsLoopback {
			fmt.Printf("Loopback:        Yes\n")
		}

		// Display statistics if available
		if network.BytesReceived > 0 || network.BytesSent > 0 {
			fmt.Println("\n📊 Network Statistics:")
			fmt.Printf("Bytes Received:  %s\n", displayer.formatBytes(network.BytesReceived))
			fmt.Printf("Bytes Sent:      %s\n", displayer.formatBytes(network.BytesSent))
			fmt.Printf("Packets Received: %d\n", network.PacketsReceived)
			fmt.Printf("Packets Sent:    %d\n", network.PacketsSent)
		}
	}
}

// displayPerformanceMetrics displays system performance metrics
func (displayer *SystemInfoDisplayer) displayPerformanceMetrics(systemInfo *SystemInfo) {
	fmt.Println("\n📈 PERFORMANCE METRICS")
	fmt.Println(strings.Repeat("-", 50))

	// Display load average
	if systemInfo.LoadAverage.Load1Minute > 0 ||
		systemInfo.LoadAverage.Load5Minutes > 0 ||
		systemInfo.LoadAverage.Load15Minutes > 0 {
		fmt.Println("Load Average:")
		fmt.Printf("  1 minute:      %.2f\n", systemInfo.LoadAverage.Load1Minute)
		fmt.Printf("  5 minutes:     %.2f\n", systemInfo.LoadAverage.Load5Minutes)
		fmt.Printf("  15 minutes:    %.2f\n", systemInfo.LoadAverage.Load15Minutes)
	}

	// Display process count
	if systemInfo.ProcessCount > 0 {
		fmt.Printf("Running Processes: %d\n", systemInfo.ProcessCount)
	}
}

// formatValue formats a string value with a default fallback
func (displayer *SystemInfoDisplayer) formatValue(value, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}

// formatBytes formats a byte count into human-readable format
func (displayer *SystemInfoDisplayer) formatBytes(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}

	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	units := []string{"KB", "MB", "GB", "TB", "PB"}
	if exp >= len(units) {
		exp = len(units) - 1
	}

	return fmt.Sprintf("%.2f %s", float64(bytes)/float64(div), units[exp])
}

// formatDuration formats a duration into human-readable format
func (displayer *SystemInfoDisplayer) formatDuration(duration time.Duration) string {
	days := int(duration.Hours() / 24)
	hours := int(duration.Hours()) % 24
	minutes := int(duration.Minutes()) % 60
	seconds := int(duration.Seconds()) % 60

	if days > 0 {
		return fmt.Sprintf("%d days, %d hours, %d minutes", days, hours, minutes)
	} else if hours > 0 {
		return fmt.Sprintf("%d hours, %d minutes", hours, minutes)
	} else if minutes > 0 {
		return fmt.Sprintf("%d minutes, %d seconds", minutes, seconds)
	} else {
		return fmt.Sprintf("%d seconds", seconds)
	}
}
