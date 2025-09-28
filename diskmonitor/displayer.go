package diskmonitor

import (
	"fmt"
	"strings"
)

// DiskMonitorDisplayer handles the display and formatting of disk monitoring data
// This struct provides methods to format and display disk data with graphical elements
type DiskMonitorDisplayer struct {
	// Display configuration
	ShowGraphics bool // Whether to show graphical elements
	ShowColors   bool // Whether to use colored output
	BarWidth     int  // Width of progress bars
	MaxProcesses int  // Maximum number of processes to display

	// Color codes for different elements
	ColorReset   string
	ColorRed     string
	ColorGreen   string
	ColorYellow  string
	ColorBlue    string
	ColorCyan    string
	ColorMagenta string
	ColorWhite   string
	ColorBold    string
}

// NewDiskMonitorDisplayer creates a new instance of DiskMonitorDisplayer
// with default configuration values
func NewDiskMonitorDisplayer() *DiskMonitorDisplayer {
	return &DiskMonitorDisplayer{
		ShowGraphics: true,
		ShowColors:   true,
		BarWidth:     50,
		MaxProcesses: 10,
		ColorReset:   "\033[0m",
		ColorRed:     "\033[31m",
		ColorGreen:   "\033[32m",
		ColorYellow:  "\033[33m",
		ColorBlue:    "\033[34m",
		ColorCyan:    "\033[36m",
		ColorMagenta: "\033[35m",
		ColorWhite:   "\033[37m",
		ColorBold:    "\033[1m",
	}
}

// DisplayDiskMonitorData displays comprehensive disk monitoring data with graphics
func (displayer *DiskMonitorDisplayer) DisplayDiskMonitorData(data *DiskMonitorData) {
	// Clear screen and move cursor to top
	fmt.Print("\033[2J\033[H")

	// Display header
	displayer.displayHeader(data)

	// Display overall disk usage with graphics
	displayer.displayOverallDiskUsage(data)

	// Display partition information
	if len(data.Partitions) > 0 {
		displayer.displayPartitionInfo(data)
	}

	// Display I/O statistics
	if len(data.DiskIO) > 0 {
		displayer.displayIOInfo(data)
	}

	// Display temperature information
	if len(data.DiskTemperatures) > 0 {
		displayer.displayTemperatureInfo(data)
	}

	// Display health information
	if len(data.DiskHealth) > 0 {
		displayer.displayHealthInfo(data)
	}

	// Display performance metrics
	displayer.displayPerformanceMetrics(data)

	// Display top processes
	if len(data.TopProcesses) > 0 {
		displayer.displayTopProcesses(data)
	}

	// Display disk status and alerts
	displayer.displayDiskStatus(data)

	// Display footer
	displayer.displayFooter(data)
}

// displayHeader displays the disk monitor header
func (displayer *DiskMonitorDisplayer) displayHeader(data *DiskMonitorData) {
	fmt.Println(displayer.colorize("💿 DISK MONITOR", displayer.ColorBold+displayer.ColorCyan))
	fmt.Println(strings.Repeat("=", 80))

	// Disk summary
	fmt.Printf("%sTotal Space: %s%s\n",
		displayer.colorize("", displayer.ColorBold),
		displayer.colorize(displayer.formatBytes(data.TotalSpace), displayer.ColorWhite),
		displayer.colorize("", displayer.ColorReset))

	fmt.Printf("%sUsed Space: %s%s\n",
		displayer.colorize("", displayer.ColorBold),
		displayer.colorize(displayer.formatBytes(data.UsedSpace), displayer.ColorRed),
		displayer.colorize("", displayer.ColorReset))

	fmt.Printf("%sFree Space: %s%s\n",
		displayer.colorize("", displayer.ColorBold),
		displayer.colorize(displayer.formatBytes(data.FreeSpace), displayer.ColorGreen),
		displayer.colorize("", displayer.ColorReset))

	fmt.Printf("%sUsage: %s%.2f%%%s\n",
		displayer.colorize("", displayer.ColorBold),
		displayer.colorize("", displayer.ColorYellow),
		data.UsagePercent,
		displayer.colorize("", displayer.ColorReset))

	fmt.Println(strings.Repeat("=", 80))
}

// displayOverallDiskUsage displays overall disk usage with graphical bars
func (displayer *DiskMonitorDisplayer) displayOverallDiskUsage(data *DiskMonitorData) {
	fmt.Println("\n📊 OVERALL DISK USAGE")
	fmt.Println(strings.Repeat("-", 50))

	// Overall usage bar
	displayer.displayUsageBar("Disk Usage", data.UsagePercent, displayer.getDiskUsageColor(data.UsagePercent))

	// Disk status indicator
	statusColor := displayer.getDiskStatusColor(data.DiskStatus)
	fmt.Printf("\n%sStatus: %s%s%s\n",
		displayer.colorize("", displayer.ColorBold),
		statusColor,
		data.DiskStatus,
		displayer.colorize("", displayer.ColorReset))
}

// displayPartitionInfo displays partition information
func (displayer *DiskMonitorDisplayer) displayPartitionInfo(data *DiskMonitorData) {
	fmt.Println("\n🔧 DISK PARTITIONS")
	fmt.Println(strings.Repeat("-", 80))

	// Header
	fmt.Printf("%s%-15s %-20s %-8s %-12s %-12s %-8s %s\n",
		displayer.colorize("", displayer.ColorBold),
		"Device",
		"Mountpoint",
		"Type",
		"Total",
		"Used",
		"Usage%",
		displayer.colorize("", displayer.ColorReset))

	fmt.Println(strings.Repeat("-", 80))

	// Display partitions
	for _, partition := range data.Partitions {
		// Truncate long mountpoints
		mountpoint := partition.Mountpoint
		if len(mountpoint) > 20 {
			mountpoint = mountpoint[:17] + "..."
		}

		// Color code based on usage
		usageColor := displayer.getDiskUsageColor(partition.UsagePercent)

		fmt.Printf("%s%-15s %-20s %-8s %s%-12s %s%-12s %s%-8.2f %s\n",
			displayer.colorize("", displayer.ColorBold),
			partition.Device,
			mountpoint,
			partition.Fstype,
			displayer.colorize("", displayer.ColorWhite),
			displayer.formatBytes(partition.Total),
			usageColor,
			displayer.formatBytes(partition.Used),
			usageColor,
			partition.UsagePercent,
			displayer.colorize("", displayer.ColorReset))

		// Partition usage bar
		displayer.displayUsageBar("  "+partition.Device, partition.UsagePercent, usageColor)
	}
}

// displayIOInfo displays disk I/O statistics
func (displayer *DiskMonitorDisplayer) displayIOInfo(data *DiskMonitorData) {
	fmt.Println("\n⚡ DISK I/O STATISTICS")
	fmt.Println(strings.Repeat("-", 80))

	// Header
	fmt.Printf("%s%-15s %-12s %-12s %-8s %-8s %-8s %s\n",
		displayer.colorize("", displayer.ColorBold),
		"Device",
		"Read Speed",
		"Write Speed",
		"IOPS",
		"Util%",
		"Reads",
		displayer.colorize("", displayer.ColorReset))

	fmt.Println(strings.Repeat("-", 80))

	// Display I/O statistics
	for _, io := range data.DiskIO {
		// Color code based on utilization
		utilColor := displayer.getUtilizationColor(io.Utilization)

		fmt.Printf("%s%-15s %s%-12s %s%-12s %s%-8.2f %s%-8.2f %s%-8d %s\n",
			displayer.colorize("", displayer.ColorBold),
			io.DeviceName,
			displayer.colorize("", displayer.ColorGreen),
			fmt.Sprintf("%.2f MB/s", io.ReadSpeed),
			displayer.colorize("", displayer.ColorYellow),
			fmt.Sprintf("%.2f MB/s", io.WriteSpeed),
			displayer.colorize("", displayer.ColorCyan),
			io.IOPS,
			utilColor,
			io.Utilization,
			displayer.colorize("", displayer.ColorWhite),
			io.ReadCount,
			displayer.colorize("", displayer.ColorReset))

		// I/O utilization bar
		displayer.displayUsageBar("  "+io.DeviceName, io.Utilization, utilColor)
	}

	// Overall I/O summary
	fmt.Printf("\n%sTotal Read Speed: %s%.2f MB/s%s\n",
		displayer.colorize("", displayer.ColorBold),
		displayer.colorize("", displayer.ColorGreen),
		data.TotalReadSpeed,
		displayer.colorize("", displayer.ColorReset))

	fmt.Printf("%sTotal Write Speed: %s%.2f MB/s%s\n",
		displayer.colorize("", displayer.ColorBold),
		displayer.colorize("", displayer.ColorYellow),
		data.TotalWriteSpeed,
		displayer.colorize("", displayer.ColorReset))

	fmt.Printf("%sAverage IOPS: %s%.2f%s\n",
		displayer.colorize("", displayer.ColorBold),
		displayer.colorize("", displayer.ColorCyan),
		data.AverageIOPS,
		displayer.colorize("", displayer.ColorReset))
}

// displayTemperatureInfo displays disk temperature information
func (displayer *DiskMonitorDisplayer) displayTemperatureInfo(data *DiskMonitorData) {
	fmt.Println("\n🌡️  DISK TEMPERATURE")
	fmt.Println(strings.Repeat("-", 50))

	for _, temp := range data.DiskTemperatures {
		fmt.Printf("%sDevice: %s%s%s\n",
			displayer.colorize("", displayer.ColorBold),
			displayer.colorize("", displayer.ColorWhite),
			temp.DeviceName,
			displayer.colorize("", displayer.ColorReset))

		fmt.Printf("%sTemperature: %s%.1f°C%s / %s%.1f°C%s\n",
			displayer.colorize("", displayer.ColorBold),
			displayer.getTemperatureColor(temp.Temperature),
			temp.Temperature,
			displayer.colorize("", displayer.ColorReset),
			displayer.colorize("", displayer.ColorWhite),
			temp.MaxTemperature,
			displayer.colorize("", displayer.ColorReset))

		// Temperature bar
		tempPercent := (temp.Temperature / temp.MaxTemperature) * 100
		displayer.displayUsageBar("Temperature", tempPercent, displayer.getTemperatureColor(temp.Temperature))

		// Temperature status
		statusColor := displayer.getTemperatureStatusColor(temp.Status)
		fmt.Printf("%sStatus: %s%s%s\n",
			displayer.colorize("", displayer.ColorBold),
			statusColor,
			temp.Status,
			displayer.colorize("", displayer.ColorReset))
	}
}

// displayHealthInfo displays disk health information
func (displayer *DiskMonitorDisplayer) displayHealthInfo(data *DiskMonitorData) {
	fmt.Println("\n💚 DISK HEALTH")
	fmt.Println(strings.Repeat("-", 80))

	// Header
	fmt.Printf("%s%-15s %-10s %-12s %-12s %-8s %s\n",
		displayer.colorize("", displayer.ColorBold),
		"Device",
		"Health",
		"Power On",
		"Cycles",
		"Wear%",
		displayer.colorize("", displayer.ColorReset))

	fmt.Println(strings.Repeat("-", 80))

	// Display health information
	for _, health := range data.DiskHealth {
		// Color code based on health status
		healthColor := displayer.getHealthStatusColor(health.HealthStatus)

		fmt.Printf("%s%-15s %s%-10s %s%-12d %s%-12d %s%-8.1f %s\n",
			displayer.colorize("", displayer.ColorBold),
			health.DeviceName,
			healthColor,
			health.HealthStatus,
			displayer.colorize("", displayer.ColorWhite),
			health.PowerOnHours,
			displayer.colorize("", displayer.ColorCyan),
			health.PowerCycleCount,
			displayer.colorize("", displayer.ColorYellow),
			health.WearLeveling,
			displayer.colorize("", displayer.ColorReset))

		// Health details
		if health.ReallocatedSectors > 0 || health.PendingSectors > 0 || health.UncorrectableSectors > 0 {
			fmt.Printf("%s  Reallocated: %s%d%s, Pending: %s%d%s, Uncorrectable: %s%d%s\n",
				displayer.colorize("", displayer.ColorBold),
				displayer.colorize("", displayer.ColorRed),
				health.ReallocatedSectors,
				displayer.colorize("", displayer.ColorReset),
				displayer.colorize("", displayer.ColorYellow),
				health.PendingSectors,
				displayer.colorize("", displayer.ColorReset),
				displayer.colorize("", displayer.ColorRed),
				health.UncorrectableSectors,
				displayer.colorize("", displayer.ColorReset))
		}
	}
}

// displayPerformanceMetrics displays disk performance metrics
func (displayer *DiskMonitorDisplayer) displayPerformanceMetrics(data *DiskMonitorData) {
	fmt.Println("\n📈 PERFORMANCE METRICS")
	fmt.Println(strings.Repeat("-", 50))

	// Overall utilization
	displayer.displayUsageBar("Disk Utilization", data.DiskUtilization, displayer.getUtilizationColor(data.DiskUtilization))

	// Performance summary
	fmt.Printf("\n%sTotal Read Speed: %s%.2f MB/s%s\n",
		displayer.colorize("", displayer.ColorBold),
		displayer.colorize("", displayer.ColorGreen),
		data.TotalReadSpeed,
		displayer.colorize("", displayer.ColorReset))

	fmt.Printf("%sTotal Write Speed: %s%.2f MB/s%s\n",
		displayer.colorize("", displayer.ColorBold),
		displayer.colorize("", displayer.ColorYellow),
		data.TotalWriteSpeed,
		displayer.colorize("", displayer.ColorReset))

	fmt.Printf("%sAverage IOPS: %s%.2f%s\n",
		displayer.colorize("", displayer.ColorBold),
		displayer.colorize("", displayer.ColorCyan),
		data.AverageIOPS,
		displayer.colorize("", displayer.ColorReset))
}

// displayTopProcesses displays top disk-consuming processes
func (displayer *DiskMonitorDisplayer) displayTopProcesses(data *DiskMonitorData) {
	fmt.Println("\n🔥 TOP DISK PROCESSES")
	fmt.Println(strings.Repeat("-", 80))

	// Header
	fmt.Printf("%s%-8s %-20s %-12s %-12s %-8s %-8s %s\n",
		displayer.colorize("", displayer.ColorBold),
		"PID",
		"Name",
		"Read Speed",
		"Write Speed",
		"IOPS",
		"Total IO",
		displayer.colorize("", displayer.ColorReset))

	fmt.Println(strings.Repeat("-", 80))

	// Display processes
	for i, process := range data.TopProcesses {
		if i >= displayer.MaxProcesses {
			break
		}

		// Truncate long process names
		name := process.Name
		if len(name) > 20 {
			name = name[:17] + "..."
		}

		// Color code based on I/O usage
		ioColor := displayer.getIOUsageColor(process.IOPS)

		fmt.Printf("%s%-8d %-20s %s%-12s %s%-12s %s%-8.2f %s%-8d %s\n",
			displayer.colorize("", displayer.ColorBold),
			process.PID,
			name,
			displayer.colorize("", displayer.ColorGreen),
			fmt.Sprintf("%.2f MB/s", process.ReadSpeed),
			displayer.colorize("", displayer.ColorYellow),
			fmt.Sprintf("%.2f MB/s", process.WriteSpeed),
			ioColor,
			process.IOPS,
			displayer.colorize("", displayer.ColorWhite),
			process.TotalIO,
			displayer.colorize("", displayer.ColorReset))
	}
}

// displayDiskStatus displays disk status and alerts
func (displayer *DiskMonitorDisplayer) displayDiskStatus(data *DiskMonitorData) {
	fmt.Println("\n🚨 DISK STATUS & ALERTS")
	fmt.Println(strings.Repeat("-", 50))

	// Disk status
	statusColor := displayer.getDiskStatusColor(data.DiskStatus)
	fmt.Printf("%sDisk Status: %s%s%s\n",
		displayer.colorize("", displayer.ColorBold),
		statusColor,
		data.DiskStatus,
		displayer.colorize("", displayer.ColorReset))

	// Low space warning
	if data.LowSpaceWarning {
		fmt.Printf("%s⚠️  Low Space Warning: %sACTIVE%s\n",
			displayer.colorize("", displayer.ColorBold),
			displayer.colorize("", displayer.ColorRed),
			displayer.colorize("", displayer.ColorReset))
	} else {
		fmt.Printf("%s✅ Low Space Warning: %sINACTIVE%s\n",
			displayer.colorize("", displayer.ColorBold),
			displayer.colorize("", displayer.ColorGreen),
			displayer.colorize("", displayer.ColorReset))
	}

	// High temperature warning
	if data.HighTempWarning {
		fmt.Printf("%s🌡️  High Temperature Warning: %sACTIVE%s\n",
			displayer.colorize("", displayer.ColorBold),
			displayer.colorize("", displayer.ColorRed),
			displayer.colorize("", displayer.ColorReset))
	} else {
		fmt.Printf("%s✅ Temperature Warning: %sNORMAL%s\n",
			displayer.colorize("", displayer.ColorBold),
			displayer.colorize("", displayer.ColorGreen),
			displayer.colorize("", displayer.ColorReset))
	}

	// Health warning
	if data.HealthWarning {
		fmt.Printf("%s💚 Health Warning: %sISSUES DETECTED%s\n",
			displayer.colorize("", displayer.ColorBold),
			displayer.colorize("", displayer.ColorRed),
			displayer.colorize("", displayer.ColorReset))
	} else {
		fmt.Printf("%s✅ Health Status: %sGOOD%s\n",
			displayer.colorize("", displayer.ColorBold),
			displayer.colorize("", displayer.ColorGreen),
			displayer.colorize("", displayer.ColorReset))
	}

	// I/O bottleneck
	if data.IOBottleneck {
		fmt.Printf("%s⚡ I/O Bottleneck: %sDETECTED%s\n",
			displayer.colorize("", displayer.ColorBold),
			displayer.colorize("", displayer.ColorRed),
			displayer.colorize("", displayer.ColorReset))
	} else {
		fmt.Printf("%s✅ I/O Performance: %sNORMAL%s\n",
			displayer.colorize("", displayer.ColorBold),
			displayer.colorize("", displayer.ColorGreen),
			displayer.colorize("", displayer.ColorReset))
	}
}

// displayUsageBar displays a graphical usage bar
func (displayer *DiskMonitorDisplayer) displayUsageBar(label string, percentage float64, color string, customWidth ...int) {
	width := displayer.BarWidth
	if len(customWidth) > 0 {
		width = customWidth[0]
	}

	// Calculate filled width
	filledWidth := int((percentage / 100.0) * float64(width))
	if filledWidth > width {
		filledWidth = width
	}

	// Create bar
	bar := strings.Repeat("█", filledWidth) + strings.Repeat("░", width-filledWidth)

	fmt.Printf("%s%-20s %s[%s]%s %s%.2f%%%s\n",
		displayer.colorize("", displayer.ColorBold),
		label,
		displayer.colorize("", displayer.ColorWhite),
		displayer.colorize(bar, color),
		displayer.colorize("", displayer.ColorWhite),
		color,
		percentage,
		displayer.colorize("", displayer.ColorReset))
}

// displayFooter displays the disk monitor footer
func (displayer *DiskMonitorDisplayer) displayFooter(data *DiskMonitorData) {
	fmt.Println(strings.Repeat("=", 80))
	fmt.Printf("%sLast Updated: %s%s%s\n",
		displayer.colorize("", displayer.ColorBold),
		displayer.colorize("", displayer.ColorCyan),
		data.Timestamp.Format("2006-01-02 15:04:05"),
		displayer.colorize("", displayer.ColorReset))

	fmt.Printf("%sRefresh Rate: %s%.1fs%s\n",
		displayer.colorize("", displayer.ColorBold),
		displayer.colorize("", displayer.ColorCyan),
		data.RefreshInterval.Seconds(),
		displayer.colorize("", displayer.ColorReset))

	fmt.Println(strings.Repeat("=", 80))
}

// formatBytes formats bytes into human-readable format
func (displayer *DiskMonitorDisplayer) formatBytes(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// colorize applies color to text if colors are enabled
func (displayer *DiskMonitorDisplayer) colorize(text, color string) string {
	if !displayer.ShowColors {
		return text
	}
	return color + text + displayer.ColorReset
}

// getDiskUsageColor returns the appropriate color for disk usage percentage
func (displayer *DiskMonitorDisplayer) getDiskUsageColor(percentage float64) string {
	if !displayer.ShowColors {
		return ""
	}

	switch {
	case percentage < 50:
		return displayer.ColorGreen
	case percentage < 80:
		return displayer.ColorYellow
	case percentage < 90:
		return displayer.ColorMagenta
	default:
		return displayer.ColorRed
	}
}

// getDiskStatusColor returns the appropriate color for disk status
func (displayer *DiskMonitorDisplayer) getDiskStatusColor(status string) string {
	if !displayer.ShowColors {
		return ""
	}

	switch status {
	case "Normal":
		return displayer.ColorGreen
	case "Warning":
		return displayer.ColorYellow
	case "Critical":
		return displayer.ColorRed
	default:
		return displayer.ColorWhite
	}
}

// getUtilizationColor returns the appropriate color for utilization percentage
func (displayer *DiskMonitorDisplayer) getUtilizationColor(percentage float64) string {
	if !displayer.ShowColors {
		return ""
	}

	switch {
	case percentage < 30:
		return displayer.ColorGreen
	case percentage < 60:
		return displayer.ColorYellow
	case percentage < 80:
		return displayer.ColorMagenta
	default:
		return displayer.ColorRed
	}
}

// getTemperatureColor returns the appropriate color for temperature
func (displayer *DiskMonitorDisplayer) getTemperatureColor(temperature float64) string {
	if !displayer.ShowColors {
		return ""
	}

	switch {
	case temperature < 40:
		return displayer.ColorGreen
	case temperature < 50:
		return displayer.ColorYellow
	case temperature < 60:
		return displayer.ColorMagenta
	default:
		return displayer.ColorRed
	}
}

// getTemperatureStatusColor returns the appropriate color for temperature status
func (displayer *DiskMonitorDisplayer) getTemperatureStatusColor(status string) string {
	if !displayer.ShowColors {
		return ""
	}

	switch status {
	case "Normal":
		return displayer.ColorGreen
	case "Warning":
		return displayer.ColorYellow
	case "Critical":
		return displayer.ColorRed
	default:
		return displayer.ColorWhite
	}
}

// getHealthStatusColor returns the appropriate color for health status
func (displayer *DiskMonitorDisplayer) getHealthStatusColor(status string) string {
	if !displayer.ShowColors {
		return ""
	}

	switch status {
	case "Good":
		return displayer.ColorGreen
	case "Warning":
		return displayer.ColorYellow
	case "Critical":
		return displayer.ColorRed
	default:
		return displayer.ColorWhite
	}
}

// getIOUsageColor returns the appropriate color for I/O usage
func (displayer *DiskMonitorDisplayer) getIOUsageColor(iops float64) string {
	if !displayer.ShowColors {
		return ""
	}

	switch {
	case iops < 100:
		return displayer.ColorGreen
	case iops < 500:
		return displayer.ColorYellow
	case iops < 1000:
		return displayer.ColorMagenta
	default:
		return displayer.ColorRed
	}
}
