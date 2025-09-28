package memorymonitor

import (
	"fmt"
	"strings"
)

// MemoryMonitorDisplayer handles the display and formatting of memory monitoring data
// This struct provides methods to format and display memory data with graphical elements
type MemoryMonitorDisplayer struct {
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

// NewMemoryMonitorDisplayer creates a new instance of MemoryMonitorDisplayer
// with default configuration values
func NewMemoryMonitorDisplayer() *MemoryMonitorDisplayer {
	return &MemoryMonitorDisplayer{
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

// DisplayMemoryMonitorData displays comprehensive memory monitoring data with graphics
func (displayer *MemoryMonitorDisplayer) DisplayMemoryMonitorData(data *MemoryMonitorData) {
	// Clear screen and move cursor to top
	fmt.Print("\033[2J\033[H")

	// Display header
	displayer.displayHeader(data)

	// Display overall memory usage with graphics
	displayer.displayOverallMemoryUsage(data)

	// Display memory breakdown
	displayer.displayMemoryBreakdown(data)

	// Display memory modules
	if len(data.MemoryModules) > 0 {
		displayer.displayMemoryModules(data)
	}

	// Display swap information
	if data.SwapInfo.TotalSwap > 0 {
		displayer.displaySwapInfo(data)
	}

	// Display cache information
	displayer.displayCacheInfo(data)

	// Display performance metrics
	displayer.displayPerformanceMetrics(data)

	// Display top processes
	if len(data.TopProcesses) > 0 {
		displayer.displayTopProcesses(data)
	}

	// Display memory status and alerts
	displayer.displayMemoryStatus(data)

	// Display footer
	displayer.displayFooter(data)
}

// displayHeader displays the memory monitor header
func (displayer *MemoryMonitorDisplayer) displayHeader(data *MemoryMonitorData) {
	fmt.Println(displayer.colorize("💾 MEMORY MONITOR", displayer.ColorBold+displayer.ColorCyan))
	fmt.Println(strings.Repeat("=", 80))

	// Memory summary
	fmt.Printf("%sTotal Memory: %s%s\n",
		displayer.colorize("", displayer.ColorBold),
		displayer.colorize(displayer.formatBytes(data.TotalMemory), displayer.ColorWhite),
		displayer.colorize("", displayer.ColorReset))

	fmt.Printf("%sAvailable: %s%s\n",
		displayer.colorize("", displayer.ColorBold),
		displayer.colorize(displayer.formatBytes(data.AvailableMemory), displayer.ColorGreen),
		displayer.colorize("", displayer.ColorReset))

	fmt.Printf("%sUsed: %s%s\n",
		displayer.colorize("", displayer.ColorBold),
		displayer.colorize(displayer.formatBytes(data.UsedMemory), displayer.ColorRed),
		displayer.colorize("", displayer.ColorReset))

	fmt.Printf("%sFree: %s%s\n",
		displayer.colorize("", displayer.ColorBold),
		displayer.colorize(displayer.formatBytes(data.FreeMemory), displayer.ColorBlue),
		displayer.colorize("", displayer.ColorReset))

	fmt.Println(strings.Repeat("=", 80))
}

// displayOverallMemoryUsage displays overall memory usage with graphical bars
func (displayer *MemoryMonitorDisplayer) displayOverallMemoryUsage(data *MemoryMonitorData) {
	fmt.Println("\n📊 OVERALL MEMORY USAGE")
	fmt.Println(strings.Repeat("-", 50))

	// Overall usage bar
	displayer.displayUsageBar("Memory Usage", data.MemoryPercent, displayer.getMemoryUsageColor(data.MemoryPercent))

	// Memory status indicator
	statusColor := displayer.getMemoryStatusColor(data.MemoryStatus)
	fmt.Printf("\n%sStatus: %s%s%s\n",
		displayer.colorize("", displayer.ColorBold),
		statusColor,
		data.MemoryStatus,
		displayer.colorize("", displayer.ColorReset))

	// Memory pressure indicator
	if data.MemoryPressure > 0 {
		displayer.displayUsageBar("Memory Pressure", data.MemoryPressure, displayer.getMemoryUsageColor(data.MemoryPressure))
	}
}

// displayMemoryBreakdown displays detailed memory breakdown
func (displayer *MemoryMonitorDisplayer) displayMemoryBreakdown(data *MemoryMonitorData) {
	fmt.Println("\n🔧 MEMORY BREAKDOWN")
	fmt.Println(strings.Repeat("-", 50))

	// User memory
	userPercent := (float64(data.UserMemory) / float64(data.TotalMemory)) * 100
	displayer.displayUsageBar("User Processes", userPercent, displayer.ColorGreen)

	// System memory
	systemPercent := (float64(data.SystemMemory) / float64(data.TotalMemory)) * 100
	displayer.displayUsageBar("System Processes", systemPercent, displayer.ColorYellow)

	// Buffer memory
	bufferPercent := (float64(data.BufferMemory) / float64(data.TotalMemory)) * 100
	displayer.displayUsageBar("Buffer Memory", bufferPercent, displayer.ColorBlue)

	// Cache memory
	cachePercent := (float64(data.CacheMemory) / float64(data.TotalMemory)) * 100
	displayer.displayUsageBar("Cache Memory", cachePercent, displayer.ColorMagenta)

	// Shared memory
	sharedPercent := (float64(data.SharedMemory) / float64(data.TotalMemory)) * 100
	displayer.displayUsageBar("Shared Memory", sharedPercent, displayer.ColorCyan)
}

// displayMemoryModules displays memory modules information
func (displayer *MemoryMonitorDisplayer) displayMemoryModules(data *MemoryMonitorData) {
	fmt.Println("\n🔧 MEMORY MODULES")
	fmt.Println(strings.Repeat("-", 50))

	for i, module := range data.MemoryModules {
		fmt.Printf("\n%sModule %d: %s%s\n",
			displayer.colorize("", displayer.ColorBold),
			i+1,
			displayer.colorize(module.Type, displayer.ColorWhite),
			displayer.colorize("", displayer.ColorReset))

		fmt.Printf("%s  Total: %s%s\n",
			displayer.colorize("", displayer.ColorBold),
			displayer.colorize(displayer.formatBytes(module.TotalSize), displayer.ColorWhite),
			displayer.colorize("", displayer.ColorReset))

		fmt.Printf("%s  Used: %s%s\n",
			displayer.colorize("", displayer.ColorBold),
			displayer.colorize(displayer.formatBytes(module.UsedSize), displayer.ColorRed),
			displayer.colorize("", displayer.ColorReset))

		fmt.Printf("%s  Free: %s%s\n",
			displayer.colorize("", displayer.ColorBold),
			displayer.colorize(displayer.formatBytes(module.FreeSize), displayer.ColorGreen),
			displayer.colorize("", displayer.ColorReset))

		// Module usage bar
		displayer.displayUsageBar("Module Usage", module.UsagePercent, displayer.getMemoryUsageColor(module.UsagePercent))

		if module.Speed > 0 {
			fmt.Printf("%s  Speed: %s%d MHz%s\n",
				displayer.colorize("", displayer.ColorBold),
				displayer.colorize("", displayer.ColorCyan),
				module.Speed,
				displayer.colorize("", displayer.ColorReset))
		}
	}
}

// displaySwapInfo displays swap memory information
func (displayer *MemoryMonitorDisplayer) displaySwapInfo(data *MemoryMonitorData) {
	fmt.Println("\n🔄 SWAP MEMORY")
	fmt.Println(strings.Repeat("-", 50))

	swapInfo := data.SwapInfo

	fmt.Printf("%sTotal Swap: %s%s\n",
		displayer.colorize("", displayer.ColorBold),
		displayer.colorize(displayer.formatBytes(swapInfo.TotalSwap), displayer.ColorWhite),
		displayer.colorize("", displayer.ColorReset))

	fmt.Printf("%sUsed Swap: %s%s\n",
		displayer.colorize("", displayer.ColorBold),
		displayer.colorize(displayer.formatBytes(swapInfo.UsedSwap), displayer.ColorRed),
		displayer.colorize("", displayer.ColorReset))

	fmt.Printf("%sFree Swap: %s%s\n",
		displayer.colorize("", displayer.ColorBold),
		displayer.colorize(displayer.formatBytes(swapInfo.FreeSwap), displayer.ColorGreen),
		displayer.colorize("", displayer.ColorReset))

	// Swap usage bar
	displayer.displayUsageBar("Swap Usage", swapInfo.SwapPercent, displayer.getSwapUsageColor(swapInfo.SwapPercent))

	// Swap status
	statusColor := displayer.getSwapStatusColor(swapInfo.SwapStatus)
	fmt.Printf("\n%sSwap Status: %s%s%s\n",
		displayer.colorize("", displayer.ColorBold),
		statusColor,
		swapInfo.SwapStatus,
		displayer.colorize("", displayer.ColorReset))

	// Swap activity
	if swapInfo.SwapIn > 0 || swapInfo.SwapOut > 0 {
		fmt.Printf("%sSwap In: %s%d pages%s\n",
			displayer.colorize("", displayer.ColorBold),
			displayer.colorize("", displayer.ColorGreen),
			swapInfo.SwapIn,
			displayer.colorize("", displayer.ColorReset))

		fmt.Printf("%sSwap Out: %s%d pages%s\n",
			displayer.colorize("", displayer.ColorBold),
			displayer.colorize("", displayer.ColorRed),
			swapInfo.SwapOut,
			displayer.colorize("", displayer.ColorReset))
	}
}

// displayCacheInfo displays system cache information
func (displayer *MemoryMonitorDisplayer) displayCacheInfo(data *MemoryMonitorData) {
	fmt.Println("\n💾 CACHE INFORMATION")
	fmt.Println(strings.Repeat("-", 50))

	cacheInfo := data.CacheInfo

	fmt.Printf("%sBuffer Cache: %s%s\n",
		displayer.colorize("", displayer.ColorBold),
		displayer.colorize(displayer.formatBytes(cacheInfo.BufferCache), displayer.ColorBlue),
		displayer.colorize("", displayer.ColorReset))

	fmt.Printf("%sPage Cache: %s%s\n",
		displayer.colorize("", displayer.ColorBold),
		displayer.colorize(displayer.formatBytes(cacheInfo.PageCache), displayer.ColorYellow),
		displayer.colorize("", displayer.ColorReset))

	fmt.Printf("%sSlab Cache: %s%s\n",
		displayer.colorize("", displayer.ColorBold),
		displayer.colorize(displayer.formatBytes(cacheInfo.SlabCache), displayer.ColorMagenta),
		displayer.colorize("", displayer.ColorReset))

	fmt.Printf("%sTotal Cache: %s%s\n",
		displayer.colorize("", displayer.ColorBold),
		displayer.colorize(displayer.formatBytes(cacheInfo.TotalCache), displayer.ColorCyan),
		displayer.colorize("", displayer.ColorReset))

	// Cache usage bar
	displayer.displayUsageBar("Cache Usage", cacheInfo.CachePercent, displayer.getCacheUsageColor(cacheInfo.CachePercent))
}

// displayPerformanceMetrics displays memory performance metrics
func (displayer *MemoryMonitorDisplayer) displayPerformanceMetrics(data *MemoryMonitorData) {
	fmt.Println("\n⚡ PERFORMANCE METRICS")
	fmt.Println(strings.Repeat("-", 50))

	// Memory fragmentation
	displayer.displayUsageBar("Memory Fragmentation", data.MemoryFragmentation, displayer.getFragmentationColor(data.MemoryFragmentation))

	// Page faults
	fmt.Printf("\n%sPage Faults: %s%d/sec%s\n",
		displayer.colorize("", displayer.ColorBold),
		displayer.colorize("", displayer.ColorYellow),
		data.PageFaults,
		displayer.colorize("", displayer.ColorReset))

	// Page ins
	fmt.Printf("%sPage Ins: %s%d/sec%s\n",
		displayer.colorize("", displayer.ColorBold),
		displayer.colorize("", displayer.ColorGreen),
		data.PageIns,
		displayer.colorize("", displayer.ColorReset))

	// Page outs
	fmt.Printf("%sPage Outs: %s%d/sec%s\n",
		displayer.colorize("", displayer.ColorBold),
		displayer.colorize("", displayer.ColorRed),
		data.PageOuts,
		displayer.colorize("", displayer.ColorReset))
}

// displayTopProcesses displays top memory-consuming processes
func (displayer *MemoryMonitorDisplayer) displayTopProcesses(data *MemoryMonitorData) {
	fmt.Println("\n🔥 TOP MEMORY PROCESSES")
	fmt.Println(strings.Repeat("-", 80))

	// Header
	fmt.Printf("%s%-8s %-20s %-12s %-8s %-10s %-8s%s\n",
		displayer.colorize("", displayer.ColorBold),
		"PID",
		"Name",
		"Memory",
		"Percent",
		"RSS",
		"Status",
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

		// Color code based on memory usage
		memColor := displayer.getMemoryUsageColor(process.MemoryPercent)
		statusColor := displayer.getProcessStatusColor(process.Status)

		fmt.Printf("%s%-8d %-20s %s%-12s %s%-8.2f %s%-10s %s%-8s%s\n",
			displayer.colorize("", displayer.ColorBold),
			process.PID,
			name,
			memColor,
			displayer.formatBytes(process.MemoryUsage),
			memColor,
			process.MemoryPercent,
			displayer.colorize("", displayer.ColorWhite),
			displayer.formatBytes(process.RSS),
			statusColor,
			process.Status,
			displayer.colorize("", displayer.ColorReset))
	}
}

// displayMemoryStatus displays memory status and alerts
func (displayer *MemoryMonitorDisplayer) displayMemoryStatus(data *MemoryMonitorData) {
	fmt.Println("\n🚨 MEMORY STATUS & ALERTS")
	fmt.Println(strings.Repeat("-", 50))

	// Memory status
	statusColor := displayer.getMemoryStatusColor(data.MemoryStatus)
	fmt.Printf("%sMemory Status: %s%s%s\n",
		displayer.colorize("", displayer.ColorBold),
		statusColor,
		data.MemoryStatus,
		displayer.colorize("", displayer.ColorReset))

	// Low memory warning
	if data.LowMemoryWarning {
		fmt.Printf("%s⚠️  Low Memory Warning: %sACTIVE%s\n",
			displayer.colorize("", displayer.ColorBold),
			displayer.colorize("", displayer.ColorRed),
			displayer.colorize("", displayer.ColorReset))
	} else {
		fmt.Printf("%s✅ Low Memory Warning: %sINACTIVE%s\n",
			displayer.colorize("", displayer.ColorBold),
			displayer.colorize("", displayer.ColorGreen),
			displayer.colorize("", displayer.ColorReset))
	}

	// Memory leak alert
	if data.MemoryLeakAlert {
		fmt.Printf("%s🔍 Memory Leak Alert: %sPOTENTIAL LEAK DETECTED%s\n",
			displayer.colorize("", displayer.ColorBold),
			displayer.colorize("", displayer.ColorRed),
			displayer.colorize("", displayer.ColorReset))
	} else {
		fmt.Printf("%s✅ Memory Leak Alert: %sNORMAL%s\n",
			displayer.colorize("", displayer.ColorBold),
			displayer.colorize("", displayer.ColorGreen),
			displayer.colorize("", displayer.ColorReset))
	}
}

// displayUsageBar displays a graphical usage bar
func (displayer *MemoryMonitorDisplayer) displayUsageBar(label string, percentage float64, color string, customWidth ...int) {
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

// displayFooter displays the memory monitor footer
func (displayer *MemoryMonitorDisplayer) displayFooter(data *MemoryMonitorData) {
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
func (displayer *MemoryMonitorDisplayer) formatBytes(bytes uint64) string {
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
func (displayer *MemoryMonitorDisplayer) colorize(text, color string) string {
	if !displayer.ShowColors {
		return text
	}
	return color + text + displayer.ColorReset
}

// getMemoryUsageColor returns the appropriate color for memory usage percentage
func (displayer *MemoryMonitorDisplayer) getMemoryUsageColor(percentage float64) string {
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

// getMemoryStatusColor returns the appropriate color for memory status
func (displayer *MemoryMonitorDisplayer) getMemoryStatusColor(status string) string {
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

// getSwapUsageColor returns the appropriate color for swap usage
func (displayer *MemoryMonitorDisplayer) getSwapUsageColor(percentage float64) string {
	if !displayer.ShowColors {
		return ""
	}

	switch {
	case percentage < 20:
		return displayer.ColorGreen
	case percentage < 50:
		return displayer.ColorYellow
	default:
		return displayer.ColorRed
	}
}

// getSwapStatusColor returns the appropriate color for swap status
func (displayer *MemoryMonitorDisplayer) getSwapStatusColor(status string) string {
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

// getCacheUsageColor returns the appropriate color for cache usage
func (displayer *MemoryMonitorDisplayer) getCacheUsageColor(percentage float64) string {
	if !displayer.ShowColors {
		return ""
	}

	switch {
	case percentage < 20:
		return displayer.ColorBlue
	case percentage < 40:
		return displayer.ColorCyan
	default:
		return displayer.ColorMagenta
	}
}

// getFragmentationColor returns the appropriate color for fragmentation
func (displayer *MemoryMonitorDisplayer) getFragmentationColor(percentage float64) string {
	if !displayer.ShowColors {
		return ""
	}

	switch {
	case percentage < 20:
		return displayer.ColorGreen
	case percentage < 40:
		return displayer.ColorYellow
	default:
		return displayer.ColorRed
	}
}

// getProcessStatusColor returns the appropriate color for process status
func (displayer *MemoryMonitorDisplayer) getProcessStatusColor(status string) string {
	if !displayer.ShowColors {
		return ""
	}

	switch status {
	case "running":
		return displayer.ColorGreen
	case "sleeping":
		return displayer.ColorBlue
	case "stopped":
		return displayer.ColorYellow
	case "zombie":
		return displayer.ColorRed
	default:
		return displayer.ColorWhite
	}
}
