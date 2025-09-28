package processmonitor

import (
	"fmt"
	"strings"
)

// ProcessMonitorDisplayer handles the display and formatting of process monitoring data
// This struct provides methods to format and display process data with graphical elements
type ProcessMonitorDisplayer struct {
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

// NewProcessMonitorDisplayer creates a new instance of ProcessMonitorDisplayer
// with default configuration values
func NewProcessMonitorDisplayer() *ProcessMonitorDisplayer {
	return &ProcessMonitorDisplayer{
		ShowGraphics: true,
		ShowColors:   true,
		BarWidth:     50,
		MaxProcesses: 20,
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

// DisplayProcessMonitorData displays comprehensive process monitoring data with graphics
func (displayer *ProcessMonitorDisplayer) DisplayProcessMonitorData(data *ProcessMonitorData) {
	// Clear screen and move cursor to top
	fmt.Print("\033[2J\033[H")

	// Display header
	displayer.displayHeader(data)

	// Display overall process statistics
	displayer.displayOverallProcessStats(data)

	// Display top processes by CPU
	if len(data.TopCPUProcesses) > 0 {
		displayer.displayTopProcesses(data.TopCPUProcesses, "CPU", "🔥 TOP CPU PROCESSES")
	}

	// Display top processes by memory
	if len(data.TopMemoryProcesses) > 0 {
		displayer.displayTopProcesses(data.TopMemoryProcesses, "Memory", "💾 TOP MEMORY PROCESSES")
	}

	// Display top processes by I/O
	if len(data.TopIOProcesses) > 0 {
		displayer.displayTopProcesses(data.TopIOProcesses, "I/O", "⚡ TOP I/O PROCESSES")
	}

	// Display top processes by threads
	if len(data.TopThreadProcesses) > 0 {
		displayer.displayTopProcesses(data.TopThreadProcesses, "Threads", "🧵 TOP THREAD PROCESSES")
	}

	// Display process tree
	if len(data.ProcessTree) > 0 {
		displayer.displayProcessTree(data.ProcessTree)
	}

	// Display process alerts
	if len(data.ProcessAlerts) > 0 {
		displayer.displayProcessAlerts(data.ProcessAlerts)
	}

	// Display process status and alerts
	displayer.displayProcessStatus(data)

	// Display footer
	displayer.displayFooter(data)
}

// displayHeader displays the process monitor header
func (displayer *ProcessMonitorDisplayer) displayHeader(data *ProcessMonitorData) {
	fmt.Println(displayer.colorize("⚙️  PROCESS MONITOR", displayer.ColorBold+displayer.ColorCyan))
	fmt.Println(strings.Repeat("=", 80))

	// Process summary
	fmt.Printf("%sTotal Processes: %s%d%s\n",
		displayer.colorize("", displayer.ColorBold),
		displayer.colorize("", displayer.ColorWhite),
		data.TotalProcesses,
		displayer.colorize("", displayer.ColorReset))

	fmt.Printf("%sRunning: %s%d%s, Sleeping: %s%d%s, Zombie: %s%d%s, Stopped: %s%d%s\n",
		displayer.colorize("", displayer.ColorBold),
		displayer.colorize("", displayer.ColorGreen),
		data.RunningProcesses,
		displayer.colorize("", displayer.ColorReset),
		displayer.colorize("", displayer.ColorYellow),
		data.SleepingProcesses,
		displayer.colorize("", displayer.ColorReset),
		displayer.colorize("", displayer.ColorRed),
		data.ZombieProcesses,
		displayer.colorize("", displayer.ColorReset),
		displayer.colorize("", displayer.ColorMagenta),
		data.StoppedProcesses,
		displayer.colorize("", displayer.ColorReset))

	fmt.Printf("%sTotal CPU: %s%.2f%%%s, Total Memory: %s%.2f%%%s\n",
		displayer.colorize("", displayer.ColorBold),
		displayer.colorize("", displayer.ColorYellow),
		data.TotalCPUUsage,
		displayer.colorize("", displayer.ColorReset),
		displayer.colorize("", displayer.ColorBlue),
		data.TotalMemoryUsage,
		displayer.colorize("", displayer.ColorReset))

	fmt.Printf("%sTotal Threads: %s%d%s, Total Open Files: %s%d%s\n",
		displayer.colorize("", displayer.ColorBold),
		displayer.colorize("", displayer.ColorCyan),
		data.TotalThreads,
		displayer.colorize("", displayer.ColorReset),
		displayer.colorize("", displayer.ColorWhite),
		data.TotalOpenFiles,
		displayer.colorize("", displayer.ColorReset))

	fmt.Println(strings.Repeat("=", 80))
}

// displayOverallProcessStats displays overall process statistics with graphical bars
func (displayer *ProcessMonitorDisplayer) displayOverallProcessStats(data *ProcessMonitorData) {
	fmt.Println("\n📊 OVERALL PROCESS STATISTICS")
	fmt.Println(strings.Repeat("-", 50))

	// CPU usage bar
	displayer.displayUsageBar("Total CPU Usage", data.TotalCPUUsage, displayer.getCPUUsageColor(data.TotalCPUUsage))

	// Memory usage bar
	displayer.displayUsageBar("Total Memory Usage", data.TotalMemoryUsage, displayer.getMemoryUsageColor(data.TotalMemoryUsage))

	// Process count bar
	processPercent := float64(data.TotalProcesses) / 1000.0 * 100 // Normalize to 1000 processes
	if processPercent > 100 {
		processPercent = 100
	}
	displayer.displayUsageBar("Process Count", processPercent, displayer.getProcessCountColor(data.TotalProcesses))

	// Thread count bar
	threadPercent := float64(data.TotalThreads) / 1000.0 * 100 // Normalize to 1000 threads
	if threadPercent > 100 {
		threadPercent = 100
	}
	displayer.displayUsageBar("Thread Count", threadPercent, displayer.getThreadCountColor(data.TotalThreads))
}

// displayTopProcesses displays top processes by a specific metric
func (displayer *ProcessMonitorDisplayer) displayTopProcesses(processes []ProcessInfo, metric, title string) {
	fmt.Printf("\n%s\n", title)
	fmt.Println(strings.Repeat("-", 80))

	// Header
	fmt.Printf("%s%-8s %-20s %-8s %-8s %-8s %-8s %-8s %s\n",
		displayer.colorize("", displayer.ColorBold),
		"PID",
		"Name",
		"CPU%",
		"Memory%",
		"Threads",
		"Status",
		"User",
		displayer.colorize("", displayer.ColorReset))

	fmt.Println(strings.Repeat("-", 80))

	// Display processes
	for i, proc := range processes {
		if i >= displayer.MaxProcesses {
			break
		}

		// Truncate long process names
		name := proc.Name
		if len(name) > 20 {
			name = name[:17] + "..."
		}

		// Color code based on metric
		var metricColor string
		var metricValue float64
		switch metric {
		case "CPU":
			metricColor = displayer.getCPUUsageColor(proc.CPUUsage)
			metricValue = proc.CPUUsage
		case "Memory":
			metricColor = displayer.getMemoryUsageColor(proc.MemoryUsage)
			metricValue = proc.MemoryUsage
		case "I/O":
			metricColor = displayer.getIOUsageColor(proc.IOReadBytes + proc.IOWriteBytes)
			metricValue = float64(proc.IOReadBytes+proc.IOWriteBytes) / (1024 * 1024) // Convert to MB
		case "Threads":
			metricColor = displayer.getThreadCountColor(proc.Threads)
			metricValue = float64(proc.Threads)
		}

		// Status color
		statusColor := displayer.getProcessStatusColor(proc.Status)

		fmt.Printf("%s%-8d %-20s %s%-8.2f %s%-8.2f %s%-8d %s%-8s %s%-8s %s\n",
			displayer.colorize("", displayer.ColorBold),
			proc.PID,
			name,
			metricColor,
			proc.CPUUsage,
			displayer.colorize("", displayer.ColorBlue),
			proc.MemoryUsage,
			displayer.colorize("", displayer.ColorCyan),
			proc.Threads,
			statusColor,
			proc.Status,
			displayer.colorize("", displayer.ColorWhite),
			proc.User,
			displayer.colorize("", displayer.ColorReset))

		// Metric bar
		displayer.displayUsageBar("  "+name, metricValue, metricColor)
	}
}

// displayProcessTree displays the process tree
func (displayer *ProcessMonitorDisplayer) displayProcessTree(tree []ProcessTreeInfo) {
	fmt.Println("\n🌳 PROCESS TREE")
	fmt.Println(strings.Repeat("-", 50))

	displayer.displayTreeLevel(tree, 0)
}

// displayTreeLevel recursively displays a level of the process tree
func (displayer *ProcessMonitorDisplayer) displayTreeLevel(tree []ProcessTreeInfo, level int) {
	for _, node := range tree {
		// Indent based on level
		indent := strings.Repeat("  ", level)

		// Tree character
		treeChar := "├─"
		if node.IsLeaf {
			treeChar = "└─"
		}

		// Color code based on level
		levelColor := displayer.getTreeLevelColor(level)

		fmt.Printf("%s%s%s%s %s%s%s\n",
			indent,
			treeChar,
			levelColor,
			node.Name,
			displayer.colorize("", displayer.ColorBold),
			fmt.Sprintf("(PID: %d)", node.PID),
			displayer.colorize("", displayer.ColorReset))

		// Display children
		if len(node.Children) > 0 {
			displayer.displayTreeLevel(node.Children, level+1)
		}
	}
}

// displayProcessAlerts displays process alerts
func (displayer *ProcessMonitorDisplayer) displayProcessAlerts(alerts []ProcessAlertInfo) {
	fmt.Println("\n🚨 PROCESS ALERTS")
	fmt.Println(strings.Repeat("-", 80))

	// Header
	fmt.Printf("%s%-8s %-20s %-20s %-10s %-15s %s\n",
		displayer.colorize("", displayer.ColorBold),
		"PID",
		"Name",
		"Alert Type",
		"Severity",
		"Value",
		displayer.colorize("", displayer.ColorReset))

	fmt.Println(strings.Repeat("-", 80))

	// Display alerts
	for _, alert := range alerts {
		// Truncate long names
		name := alert.Name
		if len(name) > 20 {
			name = name[:17] + "..."
		}

		// Severity color
		severityColor := displayer.getSeverityColor(alert.Severity)

		fmt.Printf("%s%-8d %-20s %-20s %s%-10s %s%-15.2f %s\n",
			displayer.colorize("", displayer.ColorBold),
			alert.PID,
			name,
			alert.AlertType,
			severityColor,
			alert.Severity,
			displayer.colorize("", displayer.ColorWhite),
			alert.Value,
			displayer.colorize("", displayer.ColorReset))
	}
}

// displayProcessStatus displays process status and alerts
func (displayer *ProcessMonitorDisplayer) displayProcessStatus(data *ProcessMonitorData) {
	fmt.Println("\n🚨 PROCESS STATUS & ALERTS")
	fmt.Println(strings.Repeat("-", 50))

	// Process status
	statusColor := displayer.getProcessStatusColor(data.ProcessStatus)
	fmt.Printf("%sProcess Status: %s%s%s\n",
		displayer.colorize("", displayer.ColorBold),
		statusColor,
		data.ProcessStatus,
		displayer.colorize("", displayer.ColorReset))

	// High CPU warning
	if data.HighCPUWarning {
		fmt.Printf("%s⚠️  High CPU Warning: %sACTIVE%s\n",
			displayer.colorize("", displayer.ColorBold),
			displayer.colorize("", displayer.ColorRed),
			displayer.colorize("", displayer.ColorReset))
	} else {
		fmt.Printf("%s✅ CPU Usage: %sNORMAL%s\n",
			displayer.colorize("", displayer.ColorBold),
			displayer.colorize("", displayer.ColorGreen),
			displayer.colorize("", displayer.ColorReset))
	}

	// High memory warning
	if data.HighMemoryWarning {
		fmt.Printf("%s⚠️  High Memory Warning: %sACTIVE%s\n",
			displayer.colorize("", displayer.ColorBold),
			displayer.colorize("", displayer.ColorRed),
			displayer.colorize("", displayer.ColorReset))
	} else {
		fmt.Printf("%s✅ Memory Usage: %sNORMAL%s\n",
			displayer.colorize("", displayer.ColorBold),
			displayer.colorize("", displayer.ColorGreen),
			displayer.colorize("", displayer.ColorReset))
	}

	// High I/O warning
	if data.HighIOWarning {
		fmt.Printf("%s⚠️  High I/O Warning: %sACTIVE%s\n",
			displayer.colorize("", displayer.ColorBold),
			displayer.colorize("", displayer.ColorRed),
			displayer.colorize("", displayer.ColorReset))
	} else {
		fmt.Printf("%s✅ I/O Usage: %sNORMAL%s\n",
			displayer.colorize("", displayer.ColorBold),
			displayer.colorize("", displayer.ColorGreen),
			displayer.colorize("", displayer.ColorReset))
	}

	// Zombie warning
	if data.ZombieWarning {
		fmt.Printf("%s⚠️  Zombie Process Warning: %sACTIVE%s\n",
			displayer.colorize("", displayer.ColorBold),
			displayer.colorize("", displayer.ColorRed),
			displayer.colorize("", displayer.ColorReset))
	} else {
		fmt.Printf("%s✅ Zombie Processes: %sNORMAL%s\n",
			displayer.colorize("", displayer.ColorBold),
			displayer.colorize("", displayer.ColorGreen),
			displayer.colorize("", displayer.ColorReset))
	}

	// Thread warning
	if data.ThreadWarning {
		fmt.Printf("%s⚠️  High Thread Count Warning: %sACTIVE%s\n",
			displayer.colorize("", displayer.ColorBold),
			displayer.colorize("", displayer.ColorRed),
			displayer.colorize("", displayer.ColorReset))
	} else {
		fmt.Printf("%s✅ Thread Count: %sNORMAL%s\n",
			displayer.colorize("", displayer.ColorBold),
			displayer.colorize("", displayer.ColorGreen),
			displayer.colorize("", displayer.ColorReset))
	}
}

// displayUsageBar displays a graphical usage bar
func (displayer *ProcessMonitorDisplayer) displayUsageBar(label string, value float64, color string, customWidth ...int) {
	width := displayer.BarWidth
	if len(customWidth) > 0 {
		width = customWidth[0]
	}

	// Calculate filled width
	filledWidth := int((value / 100.0) * float64(width))
	if filledWidth > width {
		filledWidth = width
	}

	// Create bar
	bar := strings.Repeat("█", filledWidth) + strings.Repeat("░", width-filledWidth)

	fmt.Printf("%s%-20s %s[%s]%s %s%.2f%s\n",
		displayer.colorize("", displayer.ColorBold),
		label,
		displayer.colorize("", displayer.ColorWhite),
		displayer.colorize(bar, color),
		displayer.colorize("", displayer.ColorWhite),
		color,
		value,
		displayer.colorize("", displayer.ColorReset))
}

// displayFooter displays the process monitor footer
func (displayer *ProcessMonitorDisplayer) displayFooter(data *ProcessMonitorData) {
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

// colorize applies color to text if colors are enabled
func (displayer *ProcessMonitorDisplayer) colorize(text, color string) string {
	if !displayer.ShowColors {
		return text
	}
	return color + text + displayer.ColorReset
}

// getCPUUsageColor returns the appropriate color for CPU usage
func (displayer *ProcessMonitorDisplayer) getCPUUsageColor(usage float64) string {
	if !displayer.ShowColors {
		return ""
	}

	switch {
	case usage < 20:
		return displayer.ColorGreen
	case usage < 50:
		return displayer.ColorYellow
	case usage < 80:
		return displayer.ColorMagenta
	default:
		return displayer.ColorRed
	}
}

// getMemoryUsageColor returns the appropriate color for memory usage
func (displayer *ProcessMonitorDisplayer) getMemoryUsageColor(usage float64) string {
	if !displayer.ShowColors {
		return ""
	}

	switch {
	case usage < 20:
		return displayer.ColorGreen
	case usage < 50:
		return displayer.ColorYellow
	case usage < 80:
		return displayer.ColorMagenta
	default:
		return displayer.ColorRed
	}
}

// getProcessCountColor returns the appropriate color for process count
func (displayer *ProcessMonitorDisplayer) getProcessCountColor(count int) string {
	if !displayer.ShowColors {
		return ""
	}

	switch {
	case count < 100:
		return displayer.ColorGreen
	case count < 500:
		return displayer.ColorYellow
	case count < 1000:
		return displayer.ColorMagenta
	default:
		return displayer.ColorRed
	}
}

// getThreadCountColor returns the appropriate color for thread count
func (displayer *ProcessMonitorDisplayer) getThreadCountColor(count int32) string {
	if !displayer.ShowColors {
		return ""
	}

	switch {
	case count < 50:
		return displayer.ColorGreen
	case count < 100:
		return displayer.ColorYellow
	case count < 200:
		return displayer.ColorMagenta
	default:
		return displayer.ColorRed
	}
}

// getIOUsageColor returns the appropriate color for I/O usage
func (displayer *ProcessMonitorDisplayer) getIOUsageColor(usage uint64) string {
	if !displayer.ShowColors {
		return ""
	}

	usageMB := float64(usage) / (1024 * 1024)
	switch {
	case usageMB < 10:
		return displayer.ColorGreen
	case usageMB < 50:
		return displayer.ColorYellow
	case usageMB < 100:
		return displayer.ColorMagenta
	default:
		return displayer.ColorRed
	}
}

// getProcessStatusColor returns the appropriate color for process status
func (displayer *ProcessMonitorDisplayer) getProcessStatusColor(status string) string {
	if !displayer.ShowColors {
		return ""
	}

	switch status {
	case "R":
		return displayer.ColorGreen
	case "S":
		return displayer.ColorYellow
	case "Z":
		return displayer.ColorRed
	case "T":
		return displayer.ColorMagenta
	default:
		return displayer.ColorWhite
	}
}

// getTreeLevelColor returns the appropriate color for tree level
func (displayer *ProcessMonitorDisplayer) getTreeLevelColor(level int) string {
	if !displayer.ShowColors {
		return ""
	}

	switch level {
	case 0:
		return displayer.ColorBold
	case 1:
		return displayer.ColorGreen
	case 2:
		return displayer.ColorBlue
	case 3:
		return displayer.ColorYellow
	default:
		return displayer.ColorCyan
	}
}

// getSeverityColor returns the appropriate color for alert severity
func (displayer *ProcessMonitorDisplayer) getSeverityColor(severity string) string {
	if !displayer.ShowColors {
		return ""
	}

	switch severity {
	case "Critical":
		return displayer.ColorRed
	case "High":
		return displayer.ColorMagenta
	case "Medium":
		return displayer.ColorYellow
	case "Low":
		return displayer.ColorGreen
	default:
		return displayer.ColorWhite
	}
}
