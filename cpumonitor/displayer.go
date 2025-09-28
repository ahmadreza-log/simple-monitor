package cpumonitor

import (
	"fmt"
	"strings"
)

// CPUMonitorDisplayer handles the display and formatting of CPU monitoring data
// This struct provides methods to format and display CPU data with graphical elements
type CPUMonitorDisplayer struct {
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

// NewCPUMonitorDisplayer creates a new instance of CPUMonitorDisplayer
// with default configuration values
func NewCPUMonitorDisplayer() *CPUMonitorDisplayer {
	return &CPUMonitorDisplayer{
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

// DisplayCPUMonitorData displays comprehensive CPU monitoring data with graphics
func (displayer *CPUMonitorDisplayer) DisplayCPUMonitorData(data *CPUMonitorData) {
	// Clear screen and move cursor to top
	fmt.Print("\033[2J\033[H")

	// Display header
	displayer.displayHeader(data)

	// Display overall CPU usage with graphics
	displayer.displayOverallUsage(data)

	// Display per-core information
	if len(data.Cores) > 0 {
		displayer.displayCoreInfo(data)
	}

	// Display temperature information
	if data.Temperature > 0 {
		displayer.displayTemperatureInfo(data)
	}

	// Display load average
	if data.LoadAverage1Min > 0 || data.LoadAverage5Min > 0 || data.LoadAverage15Min > 0 {
		displayer.displayLoadAverage(data)
	}

	// Display top processes
	if len(data.TopProcesses) > 0 {
		displayer.displayTopProcesses(data)
	}

	// Display footer
	displayer.displayFooter(data)
}

// displayHeader displays the CPU monitor header
func (displayer *CPUMonitorDisplayer) displayHeader(data *CPUMonitorData) {
	fmt.Println(displayer.colorize("🖥️  CPU MONITOR", displayer.ColorBold+displayer.ColorCyan))
	fmt.Println(strings.Repeat("=", 80))

	// CPU model and basic info
	fmt.Printf("%sCPU Model: %s%s\n",
		displayer.colorize("", displayer.ColorBold),
		displayer.colorize(data.ModelName, displayer.ColorWhite),
		displayer.colorize("", displayer.ColorReset))

	fmt.Printf("%sArchitecture: %s%s\n",
		displayer.colorize("", displayer.ColorBold),
		displayer.colorize(data.Architecture, displayer.ColorWhite),
		displayer.colorize("", displayer.ColorReset))

	fmt.Printf("%sCores: %s%d Physical, %d Logical%s\n",
		displayer.colorize("", displayer.ColorBold),
		displayer.colorize("", displayer.ColorWhite),
		data.PhysicalCores,
		data.LogicalCores,
		displayer.colorize("", displayer.ColorReset))

	fmt.Println(strings.Repeat("=", 80))
}

// displayOverallUsage displays overall CPU usage with graphical bars
func (displayer *CPUMonitorDisplayer) displayOverallUsage(data *CPUMonitorData) {
	fmt.Println("\n📊 OVERALL CPU USAGE")
	fmt.Println(strings.Repeat("-", 50))

	// Overall usage bar
	displayer.displayUsageBar("Overall", data.OverallUsage, displayer.getUsageColor(data.OverallUsage))

	// Detailed breakdown
	fmt.Printf("\n%sUser Processes: %s%.2f%%%s\n",
		displayer.colorize("", displayer.ColorBold),
		displayer.colorize("", displayer.ColorGreen),
		data.UserUsage,
		displayer.colorize("", displayer.ColorReset))

	fmt.Printf("%sSystem Processes: %s%.2f%%%s\n",
		displayer.colorize("", displayer.ColorBold),
		displayer.colorize("", displayer.ColorYellow),
		data.SystemUsage,
		displayer.colorize("", displayer.ColorReset))

	fmt.Printf("%sIdle: %s%.2f%%%s\n",
		displayer.colorize("", displayer.ColorBold),
		displayer.colorize("", displayer.ColorBlue),
		data.IdleUsage,
		displayer.colorize("", displayer.ColorReset))

	if data.IOWaitUsage > 0 {
		fmt.Printf("%sI/O Wait: %s%.2f%%%s\n",
			displayer.colorize("", displayer.ColorBold),
			displayer.colorize("", displayer.ColorMagenta),
			data.IOWaitUsage,
			displayer.colorize("", displayer.ColorReset))
	}
}

// displayCoreInfo displays per-core CPU usage information
func (displayer *CPUMonitorDisplayer) displayCoreInfo(data *CPUMonitorData) {
	fmt.Println("\n🔧 PER-CORE USAGE")
	fmt.Println(strings.Repeat("-", 50))

	// Display cores in a grid layout
	coresPerRow := 4
	for i := 0; i < len(data.Cores); i += coresPerRow {
		fmt.Printf("\n")
		for j := 0; j < coresPerRow && i+j < len(data.Cores); j++ {
			core := data.Cores[i+j]
			coreLabel := fmt.Sprintf("Core %d", core.CoreID)
			if core.IsHyperthreaded {
				coreLabel += " (HT)"
			}

			// Shortened bar for grid display
			shortBarWidth := displayer.BarWidth / 2
			displayer.displayUsageBar(coreLabel, core.UsagePercent, displayer.getUsageColor(core.UsagePercent), shortBarWidth)
		}
	}
}

// displayTemperatureInfo displays CPU temperature information
func (displayer *CPUMonitorDisplayer) displayTemperatureInfo(data *CPUMonitorData) {
	fmt.Println("\n🌡️  TEMPERATURE")
	fmt.Println(strings.Repeat("-", 50))

	// Temperature bar
	tempPercent := (data.Temperature / data.MaxTemperature) * 100
	if tempPercent > 100 {
		tempPercent = 100
	}

	fmt.Printf("%sCPU Temperature: %s%.1f°C%s / %s%.1f°C%s\n",
		displayer.colorize("", displayer.ColorBold),
		displayer.getTemperatureColor(data.Temperature),
		data.Temperature,
		displayer.colorize("", displayer.ColorReset),
		displayer.colorize("", displayer.ColorWhite),
		data.MaxTemperature,
		displayer.colorize("", displayer.ColorReset))

	// Temperature bar
	displayer.displayUsageBar("Temperature", tempPercent, displayer.getTemperatureColor(data.Temperature))

	// Temperature status
	statusColor := displayer.getTemperatureStatusColor(data.TemperatureStatus)
	fmt.Printf("\n%sStatus: %s%s%s\n",
		displayer.colorize("", displayer.ColorBold),
		statusColor,
		data.TemperatureStatus,
		displayer.colorize("", displayer.ColorReset))
}

// displayLoadAverage displays system load average
func (displayer *CPUMonitorDisplayer) displayLoadAverage(data *CPUMonitorData) {
	fmt.Println("\n📈 LOAD AVERAGE")
	fmt.Println(strings.Repeat("-", 50))

	fmt.Printf("%s1 minute:  %s%.2f%s\n",
		displayer.colorize("", displayer.ColorBold),
		displayer.colorize("", displayer.ColorGreen),
		data.LoadAverage1Min,
		displayer.colorize("", displayer.ColorReset))

	fmt.Printf("%s5 minutes: %s%.2f%s\n",
		displayer.colorize("", displayer.ColorBold),
		displayer.colorize("", displayer.ColorYellow),
		data.LoadAverage5Min,
		displayer.colorize("", displayer.ColorReset))

	fmt.Printf("%s15 minutes:%s%.2f%s\n",
		displayer.colorize("", displayer.ColorBold),
		displayer.colorize("", displayer.ColorRed),
		data.LoadAverage15Min,
		displayer.colorize("", displayer.ColorReset))
}

// displayTopProcesses displays top CPU-consuming processes
func (displayer *CPUMonitorDisplayer) displayTopProcesses(data *CPUMonitorData) {
	fmt.Println("\n⚙️  TOP PROCESSES")
	fmt.Println(strings.Repeat("-", 50))

	// Limit number of processes to display
	maxProcesses := displayer.MaxProcesses
	if len(data.TopProcesses) < maxProcesses {
		maxProcesses = len(data.TopProcesses)
	}

	// Display header
	fmt.Printf("%s%-8s %-20s %-8s %-10s %s\n",
		displayer.colorize("", displayer.ColorBold),
		"PID",
		"Process",
		"CPU%",
		"Status",
		displayer.colorize("", displayer.ColorReset))

	fmt.Println(strings.Repeat("-", 50))

	// Display processes
	for i := 0; i < maxProcesses; i++ {
		process := data.TopProcesses[i]
		displayer.displayProcessInfo(process)
	}
}

// displayProcessInfo displays information about a single process
func (displayer *CPUMonitorDisplayer) displayProcessInfo(process CPUProcessInfo) {
	// Truncate process name if too long
	processName := process.Name
	if len(processName) > 20 {
		processName = processName[:17] + "..."
	}

	// Get color based on CPU usage
	cpuColor := displayer.getUsageColor(process.CPUUsagePercent)

	fmt.Printf("%s%-8d %-20s %s%-8.2f%s %-10.2f %s\n",
		displayer.colorize("", displayer.ColorWhite),
		process.PID,
		processName,
		cpuColor,
		process.CPUUsagePercent,
		displayer.colorize("", displayer.ColorReset),
		process.MemoryPercent,
		process.Status)
}

// displayUsageBar displays a graphical usage bar
func (displayer *CPUMonitorDisplayer) displayUsageBar(label string, percentage float64, color string, customWidth ...int) {
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

	fmt.Printf("%s%-15s %s[%s]%s %s%.2f%%%s\n",
		displayer.colorize("", displayer.ColorBold),
		label,
		displayer.colorize("", displayer.ColorWhite),
		displayer.colorize(bar, color),
		displayer.colorize("", displayer.ColorWhite),
		color,
		percentage,
		displayer.colorize("", displayer.ColorReset))
}

// displayFooter displays the CPU monitor footer
func (displayer *CPUMonitorDisplayer) displayFooter(data *CPUMonitorData) {
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

// getUsageColor returns the appropriate color for a given usage percentage
func (displayer *CPUMonitorDisplayer) getUsageColor(percentage float64) string {
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

// getTemperatureColor returns the appropriate color for a given temperature
func (displayer *CPUMonitorDisplayer) getTemperatureColor(temperature float64) string {
	if !displayer.ShowColors {
		return ""
	}

	switch {
	case temperature < 50:
		return displayer.ColorGreen
	case temperature < 70:
		return displayer.ColorYellow
	case temperature < 85:
		return displayer.ColorMagenta
	default:
		return displayer.ColorRed
	}
}

// getTemperatureStatusColor returns the appropriate color for temperature status
func (displayer *CPUMonitorDisplayer) getTemperatureStatusColor(status string) string {
	if !displayer.ShowColors {
		return ""
	}

	switch strings.ToLower(status) {
	case "normal":
		return displayer.ColorGreen
	case "warning":
		return displayer.ColorYellow
	case "critical":
		return displayer.ColorRed
	default:
		return displayer.ColorWhite
	}
}

// colorize applies color to text if colors are enabled
func (displayer *CPUMonitorDisplayer) colorize(text, color string) string {
	if !displayer.ShowColors {
		return text
	}
	return color + text + displayer.ColorReset
}

// SetGraphicsEnabled enables or disables graphical elements
func (displayer *CPUMonitorDisplayer) SetGraphicsEnabled(enabled bool) {
	displayer.ShowGraphics = enabled
}

// SetColorsEnabled enables or disables colored output
func (displayer *CPUMonitorDisplayer) SetColorsEnabled(enabled bool) {
	displayer.ShowColors = enabled
}

// SetBarWidth sets the width of progress bars
func (displayer *CPUMonitorDisplayer) SetBarWidth(width int) {
	displayer.BarWidth = width
}

// SetMaxProcesses sets the maximum number of processes to display
func (displayer *CPUMonitorDisplayer) SetMaxProcesses(max int) {
	displayer.MaxProcesses = max
}
