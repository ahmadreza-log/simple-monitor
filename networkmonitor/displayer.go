package networkmonitor

import (
	"fmt"
	"strings"
)

// NetworkMonitorDisplayer handles the display and formatting of network monitoring data
// This struct provides methods to format and display network data with graphical elements
type NetworkMonitorDisplayer struct {
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

// NewNetworkMonitorDisplayer creates a new instance of NetworkMonitorDisplayer
// with default configuration values
func NewNetworkMonitorDisplayer() *NetworkMonitorDisplayer {
	return &NetworkMonitorDisplayer{
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

// DisplayNetworkMonitorData displays comprehensive network monitoring data with graphics
func (displayer *NetworkMonitorDisplayer) DisplayNetworkMonitorData(data *NetworkMonitorData) {
	// Clear screen and move cursor to top
	fmt.Print("\033[2J\033[H")

	// Display header
	displayer.displayHeader(data)

	// Display overall network statistics
	displayer.displayOverallNetworkStats(data)

	// Display interface information
	if len(data.Interfaces) > 0 {
		displayer.displayInterfaceInfo(data)
	}

	// Display I/O statistics
	if len(data.InterfaceIO) > 0 {
		displayer.displayIOInfo(data)
	}

	// Display connection information
	if len(data.Connections) > 0 {
		displayer.displayConnectionInfo(data)
	}

	// Display latency information
	if len(data.LatencyInfo) > 0 {
		displayer.displayLatencyInfo(data)
	}

	// Display bandwidth information
	displayer.displayBandwidthInfo(data)

	// Display performance metrics
	displayer.displayPerformanceMetrics(data)

	// Display top processes
	if len(data.TopProcesses) > 0 {
		displayer.displayTopProcesses(data)
	}

	// Display network status and alerts
	displayer.displayNetworkStatus(data)

	// Display footer
	displayer.displayFooter(data)
}

// displayHeader displays the network monitor header
func (displayer *NetworkMonitorDisplayer) displayHeader(data *NetworkMonitorData) {
	fmt.Println(displayer.colorize("🌐 NETWORK MONITOR", displayer.ColorBold+displayer.ColorCyan))
	fmt.Println(strings.Repeat("=", 80))

	// Network summary
	fmt.Printf("%sTotal Sent: %s%s\n",
		displayer.colorize("", displayer.ColorBold),
		displayer.colorize(displayer.formatBytes(data.TotalBytesSent), displayer.ColorGreen),
		displayer.colorize("", displayer.ColorReset))

	fmt.Printf("%sTotal Received: %s%s\n",
		displayer.colorize("", displayer.ColorBold),
		displayer.colorize(displayer.formatBytes(data.TotalBytesRecv), displayer.ColorBlue),
		displayer.colorize("", displayer.ColorReset))

	fmt.Printf("%sTotal Throughput: %s%.2f Mbps%s\n",
		displayer.colorize("", displayer.ColorBold),
		displayer.colorize("", displayer.ColorYellow),
		data.TotalThroughput,
		displayer.colorize("", displayer.ColorReset))

	fmt.Printf("%sNetwork Status: %s%s%s\n",
		displayer.colorize("", displayer.ColorBold),
		displayer.getNetworkStatusColor(data.NetworkStatus),
		data.NetworkStatus,
		displayer.colorize("", displayer.ColorReset))

	fmt.Println(strings.Repeat("=", 80))
}

// displayOverallNetworkStats displays overall network statistics with graphical bars
func (displayer *NetworkMonitorDisplayer) displayOverallNetworkStats(data *NetworkMonitorData) {
	fmt.Println("\n📊 OVERALL NETWORK STATISTICS")
	fmt.Println(strings.Repeat("-", 50))

	// Send speed bar
	displayer.displayUsageBar("Send Speed", data.TotalSendSpeed, displayer.ColorGreen)

	// Receive speed bar
	displayer.displayUsageBar("Receive Speed", data.TotalRecvSpeed, displayer.ColorBlue)

	// Total throughput bar
	displayer.displayUsageBar("Total Throughput", data.TotalThroughput, displayer.ColorYellow)

	// Network utilization bar
	displayer.displayUsageBar("Network Utilization", data.NetworkUtilization, displayer.getUtilizationColor(data.NetworkUtilization))
}

// displayInterfaceInfo displays network interface information
func (displayer *NetworkMonitorDisplayer) displayInterfaceInfo(data *NetworkMonitorData) {
	fmt.Println("\n🔧 NETWORK INTERFACES")
	fmt.Println(strings.Repeat("-", 80))

	// Header
	fmt.Printf("%s%-15s %-10s %-15s %-15s %-8s %-6s %s\n",
		displayer.colorize("", displayer.ColorBold),
		"Interface",
		"Type",
		"IP Address",
		"MAC Address",
		"Status",
		"Speed",
		displayer.colorize("", displayer.ColorReset))

	fmt.Println(strings.Repeat("-", 80))

	// Display interfaces
	for _, iface := range data.Interfaces {
		// Color code based on status
		statusColor := displayer.getInterfaceStatusColor(iface.Status)

		fmt.Printf("%s%-15s %-10s %-15s %-15s %s%-8s %s%-6d %s\n",
			displayer.colorize("", displayer.ColorBold),
			iface.Name,
			iface.Type,
			iface.IPAddress,
			iface.MACAddress,
			statusColor,
			iface.Status,
			displayer.colorize("", displayer.ColorWhite),
			iface.Speed,
			displayer.colorize("", displayer.ColorReset))

		// Interface status indicator
		if iface.IsUp {
			fmt.Printf("%s  Status: %sUP%s\n",
				displayer.colorize("", displayer.ColorBold),
				displayer.colorize("", displayer.ColorGreen),
				displayer.colorize("", displayer.ColorReset))
		} else {
			fmt.Printf("%s  Status: %sDOWN%s\n",
				displayer.colorize("", displayer.ColorBold),
				displayer.colorize("", displayer.ColorRed),
				displayer.colorize("", displayer.ColorReset))
		}
	}
}

// displayIOInfo displays network I/O statistics
func (displayer *NetworkMonitorDisplayer) displayIOInfo(data *NetworkMonitorData) {
	fmt.Println("\n⚡ NETWORK I/O STATISTICS")
	fmt.Println(strings.Repeat("-", 80))

	// Header
	fmt.Printf("%s%-15s %-12s %-12s %-8s %-8s %-8s %s\n",
		displayer.colorize("", displayer.ColorBold),
		"Interface",
		"Send Speed",
		"Recv Speed",
		"Packets",
		"Errors",
		"Util%",
		displayer.colorize("", displayer.ColorReset))

	fmt.Println(strings.Repeat("-", 80))

	// Display I/O statistics
	for _, io := range data.InterfaceIO {
		// Color code based on utilization
		utilColor := displayer.getUtilizationColor(io.Utilization)

		fmt.Printf("%s%-15s %s%-12s %s%-12s %s%-8d %s%-8d %s%-8.2f %s\n",
			displayer.colorize("", displayer.ColorBold),
			io.InterfaceName,
			displayer.colorize("", displayer.ColorGreen),
			fmt.Sprintf("%.2f Mbps", io.SendSpeed),
			displayer.colorize("", displayer.ColorBlue),
			fmt.Sprintf("%.2f Mbps", io.RecvSpeed),
			displayer.colorize("", displayer.ColorWhite),
			io.PacketsSent + io.PacketsRecv,
			displayer.colorize("", displayer.ColorRed),
			io.SendErrors + io.RecvErrors,
			utilColor,
			io.Utilization,
			displayer.colorize("", displayer.ColorReset))

		// I/O utilization bar
		displayer.displayUsageBar("  "+io.InterfaceName, io.Utilization, utilColor)
	}

	// Overall I/O summary
	fmt.Printf("\n%sTotal Send Speed: %s%.2f Mbps%s\n",
		displayer.colorize("", displayer.ColorBold),
		displayer.colorize("", displayer.ColorGreen),
		data.TotalSendSpeed,
		displayer.colorize("", displayer.ColorReset))

	fmt.Printf("%sTotal Receive Speed: %s%.2f Mbps%s\n",
		displayer.colorize("", displayer.ColorBold),
		displayer.colorize("", displayer.ColorBlue),
		data.TotalRecvSpeed,
		displayer.colorize("", displayer.ColorReset))

	fmt.Printf("%sTotal Packets: %s%d%s\n",
		displayer.colorize("", displayer.ColorBold),
		displayer.colorize("", displayer.ColorWhite),
		data.TotalPacketsSent + data.TotalPacketsRecv,
		displayer.colorize("", displayer.ColorReset))
}

// displayConnectionInfo displays network connection information
func (displayer *NetworkMonitorDisplayer) displayConnectionInfo(data *NetworkMonitorData) {
	fmt.Println("\n🔗 NETWORK CONNECTIONS")
	fmt.Println(strings.Repeat("-", 80))

	// Header
	fmt.Printf("%s%-20s %-20s %-8s %-8s %-15s %s\n",
		displayer.colorize("", displayer.ColorBold),
		"Local Address",
		"Remote Address",
		"Type",
		"Status",
		"Process",
		displayer.colorize("", displayer.ColorReset))

	fmt.Println(strings.Repeat("-", 80))

	// Display connections
	for _, conn := range data.Connections {
		// Truncate long addresses
		localAddr := conn.LocalAddress
		remoteAddr := conn.RemoteAddress
		if len(localAddr) > 20 {
			localAddr = localAddr[:17] + "..."
		}
		if len(remoteAddr) > 20 {
			remoteAddr = remoteAddr[:17] + "..."
		}

		// Color code based on connection type
		typeColor := displayer.getConnectionTypeColor(conn.Type)

		fmt.Printf("%s%-20s %-20s %s%-8s %s%-8s %s%-15s %s\n",
			displayer.colorize("", displayer.ColorBold),
			localAddr,
			remoteAddr,
			typeColor,
			conn.Type,
			displayer.colorize("", displayer.ColorWhite),
			conn.Status,
			displayer.colorize("", displayer.ColorCyan),
			conn.ProcessName,
			displayer.colorize("", displayer.ColorReset))
	}
}

// displayLatencyInfo displays network latency information
func (displayer *NetworkMonitorDisplayer) displayLatencyInfo(data *NetworkMonitorData) {
	fmt.Println("\n⏱️  NETWORK LATENCY")
	fmt.Println(strings.Repeat("-", 50))

	for _, latency := range data.LatencyInfo {
		fmt.Printf("%sTarget: %s%s%s\n",
			displayer.colorize("", displayer.ColorBold),
			displayer.colorize("", displayer.ColorWhite),
			latency.Target,
			displayer.colorize("", displayer.ColorReset))

		fmt.Printf("%sLatency: %s%.2f ms%s\n",
			displayer.colorize("", displayer.ColorBold),
			displayer.getLatencyColor(latency.Latency),
			latency.Latency,
			displayer.colorize("", displayer.ColorReset))

		fmt.Printf("%sPacket Loss: %s%.2f%%%s\n",
			displayer.colorize("", displayer.ColorBold),
			displayer.getPacketLossColor(latency.PacketLoss),
			latency.PacketLoss,
			displayer.colorize("", displayer.ColorReset))

		// Latency bar
		latencyPercent := (latency.Latency / 200.0) * 100 // Normalize to 200ms
		if latencyPercent > 100 {
			latencyPercent = 100
		}
		displayer.displayUsageBar("Latency", latencyPercent, displayer.getLatencyColor(latency.Latency))

		// Status
		statusColor := displayer.getLatencyStatusColor(latency.Status)
		fmt.Printf("%sStatus: %s%s%s\n",
			displayer.colorize("", displayer.ColorBold),
			statusColor,
			latency.Status,
			displayer.colorize("", displayer.ColorReset))
	}
}

// displayBandwidthInfo displays bandwidth usage information
func (displayer *NetworkMonitorDisplayer) displayBandwidthInfo(data *NetworkMonitorData) {
	fmt.Println("\n📈 BANDWIDTH USAGE")
	fmt.Println(strings.Repeat("-", 50))

	// Bandwidth utilization bar
	displayer.displayUsageBar("Bandwidth Usage", data.BandwidthInfo.Utilization, displayer.getUtilizationColor(data.BandwidthInfo.Utilization))

	// Bandwidth details
	fmt.Printf("\n%sTotal Bandwidth: %s%.2f Mbps%s\n",
		displayer.colorize("", displayer.ColorBold),
		displayer.colorize("", displayer.ColorWhite),
		data.BandwidthInfo.TotalBandwidth,
		displayer.colorize("", displayer.ColorReset))

	fmt.Printf("%sUsed Bandwidth: %s%.2f Mbps%s\n",
		displayer.colorize("", displayer.ColorBold),
		displayer.colorize("", displayer.ColorYellow),
		data.BandwidthInfo.UsedBandwidth,
		displayer.colorize("", displayer.ColorReset))

	fmt.Printf("%sAvailable Bandwidth: %s%.2f Mbps%s\n",
		displayer.colorize("", displayer.ColorBold),
		displayer.colorize("", displayer.ColorGreen),
		data.BandwidthInfo.AvailableBandwidth,
		displayer.colorize("", displayer.ColorReset))

	fmt.Printf("%sPeak Usage: %s%.2f Mbps%s\n",
		displayer.colorize("", displayer.ColorBold),
		displayer.colorize("", displayer.ColorMagenta),
		data.BandwidthInfo.PeakUsage,
		displayer.colorize("", displayer.ColorReset))
}

// displayPerformanceMetrics displays network performance metrics
func (displayer *NetworkMonitorDisplayer) displayPerformanceMetrics(data *NetworkMonitorData) {
	fmt.Println("\n📊 PERFORMANCE METRICS")
	fmt.Println(strings.Repeat("-", 50))

	// Average latency
	displayer.displayUsageBar("Average Latency", data.AverageLatency, displayer.getLatencyColor(data.AverageLatency))

	// Packet loss rate
	displayer.displayUsageBar("Packet Loss Rate", data.PacketLossRate, displayer.getPacketLossColor(data.PacketLossRate))

	// Network utilization
	displayer.displayUsageBar("Network Utilization", data.NetworkUtilization, displayer.getUtilizationColor(data.NetworkUtilization))

	// Performance summary
	fmt.Printf("\n%sAverage Latency: %s%.2f ms%s\n",
		displayer.colorize("", displayer.ColorBold),
		displayer.colorize("", displayer.ColorCyan),
		data.AverageLatency,
		displayer.colorize("", displayer.ColorReset))

	fmt.Printf("%sPacket Loss Rate: %s%.2f%%%s\n",
		displayer.colorize("", displayer.ColorBold),
		displayer.colorize("", displayer.ColorRed),
		data.PacketLossRate,
		displayer.colorize("", displayer.ColorReset))

	fmt.Printf("%sNetwork Utilization: %s%.2f%%%s\n",
		displayer.colorize("", displayer.ColorBold),
		displayer.colorize("", displayer.ColorYellow),
		data.NetworkUtilization,
		displayer.colorize("", displayer.ColorReset))
}

// displayTopProcesses displays top network-consuming processes
func (displayer *NetworkMonitorDisplayer) displayTopProcesses(data *NetworkMonitorData) {
	fmt.Println("\n🔥 TOP NETWORK PROCESSES")
	fmt.Println(strings.Repeat("-", 80))

	// Header
	fmt.Printf("%s%-8s %-20s %-12s %-12s %-8s %-8s %s\n",
		displayer.colorize("", displayer.ColorBold),
		"PID",
		"Name",
		"Send Speed",
		"Recv Speed",
		"Total",
		"Connections",
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

		// Color code based on total speed
		speedColor := displayer.getNetworkSpeedColor(process.TotalSpeed)

		fmt.Printf("%s%-8d %-20s %s%-12s %s%-12s %s%-8.2f %s%-8d %s\n",
			displayer.colorize("", displayer.ColorBold),
			process.PID,
			name,
			displayer.colorize("", displayer.ColorGreen),
			fmt.Sprintf("%.2f Mbps", process.SendSpeed),
			displayer.colorize("", displayer.ColorBlue),
			fmt.Sprintf("%.2f Mbps", process.RecvSpeed),
			speedColor,
			process.TotalSpeed,
			displayer.colorize("", displayer.ColorWhite),
			process.Connections,
			displayer.colorize("", displayer.ColorReset))
	}
}

// displayNetworkStatus displays network status and alerts
func (displayer *NetworkMonitorDisplayer) displayNetworkStatus(data *NetworkMonitorData) {
	fmt.Println("\n🚨 NETWORK STATUS & ALERTS")
	fmt.Println(strings.Repeat("-", 50))

	// Network status
	statusColor := displayer.getNetworkStatusColor(data.NetworkStatus)
	fmt.Printf("%sNetwork Status: %s%s%s\n",
		displayer.colorize("", displayer.ColorBold),
		statusColor,
		data.NetworkStatus,
		displayer.colorize("", displayer.ColorReset))

	// High latency warning
	if data.HighLatencyWarning {
		fmt.Printf("%s⚠️  High Latency Warning: %sACTIVE%s\n",
			displayer.colorize("", displayer.ColorBold),
			displayer.colorize("", displayer.ColorRed),
			displayer.colorize("", displayer.ColorReset))
	} else {
		fmt.Printf("%s✅ Latency Status: %sNORMAL%s\n",
			displayer.colorize("", displayer.ColorBold),
			displayer.colorize("", displayer.ColorGreen),
			displayer.colorize("", displayer.ColorReset))
	}

	// Packet loss warning
	if data.PacketLossWarning {
		fmt.Printf("%s📦 Packet Loss Warning: %sACTIVE%s\n",
			displayer.colorize("", displayer.ColorBold),
			displayer.colorize("", displayer.ColorRed),
			displayer.colorize("", displayer.ColorReset))
	} else {
		fmt.Printf("%s✅ Packet Loss Status: %sNORMAL%s\n",
			displayer.colorize("", displayer.ColorBold),
			displayer.colorize("", displayer.ColorGreen),
			displayer.colorize("", displayer.ColorReset))
	}

	// Bandwidth warning
	if data.BandwidthWarning {
		fmt.Printf("%s📊 Bandwidth Warning: %sHIGH USAGE%s\n",
			displayer.colorize("", displayer.ColorBold),
			displayer.colorize("", displayer.ColorRed),
			displayer.colorize("", displayer.ColorReset))
	} else {
		fmt.Printf("%s✅ Bandwidth Status: %sNORMAL%s\n",
			displayer.colorize("", displayer.ColorBold),
			displayer.colorize("", displayer.ColorGreen),
			displayer.colorize("", displayer.ColorReset))
	}

	// Connection warning
	if data.ConnectionWarning {
		fmt.Printf("%s🔗 Connection Warning: %sISSUES DETECTED%s\n",
			displayer.colorize("", displayer.ColorBold),
			displayer.colorize("", displayer.ColorRed),
			displayer.colorize("", displayer.ColorReset))
	} else {
		fmt.Printf("%s✅ Connection Status: %sNORMAL%s\n",
			displayer.colorize("", displayer.ColorBold),
			displayer.colorize("", displayer.ColorGreen),
			displayer.colorize("", displayer.ColorReset))
	}
}

// displayUsageBar displays a graphical usage bar
func (displayer *NetworkMonitorDisplayer) displayUsageBar(label string, value float64, color string, customWidth ...int) {
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

// displayFooter displays the network monitor footer
func (displayer *NetworkMonitorDisplayer) displayFooter(data *NetworkMonitorData) {
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
func (displayer *NetworkMonitorDisplayer) formatBytes(bytes uint64) string {
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
func (displayer *NetworkMonitorDisplayer) colorize(text, color string) string {
	if !displayer.ShowColors {
		return text
	}
	return color + text + displayer.ColorReset
}

// getNetworkStatusColor returns the appropriate color for network status
func (displayer *NetworkMonitorDisplayer) getNetworkStatusColor(status string) string {
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

// getInterfaceStatusColor returns the appropriate color for interface status
func (displayer *NetworkMonitorDisplayer) getInterfaceStatusColor(status string) string {
	if !displayer.ShowColors {
		return ""
	}

	switch status {
	case "up":
		return displayer.ColorGreen
	case "down":
		return displayer.ColorRed
	default:
		return displayer.ColorYellow
	}
}

// getUtilizationColor returns the appropriate color for utilization percentage
func (displayer *NetworkMonitorDisplayer) getUtilizationColor(percentage float64) string {
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

// getLatencyColor returns the appropriate color for latency
func (displayer *NetworkMonitorDisplayer) getLatencyColor(latency float64) string {
	if !displayer.ShowColors {
		return ""
	}

	switch {
	case latency < 50:
		return displayer.ColorGreen
	case latency < 100:
		return displayer.ColorYellow
	case latency < 200:
		return displayer.ColorMagenta
	default:
		return displayer.ColorRed
	}
}

// getLatencyStatusColor returns the appropriate color for latency status
func (displayer *NetworkMonitorDisplayer) getLatencyStatusColor(status string) string {
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

// getPacketLossColor returns the appropriate color for packet loss
func (displayer *NetworkMonitorDisplayer) getPacketLossColor(packetLoss float64) string {
	if !displayer.ShowColors {
		return ""
	}

	switch {
	case packetLoss < 1:
		return displayer.ColorGreen
	case packetLoss < 5:
		return displayer.ColorYellow
	case packetLoss < 10:
		return displayer.ColorMagenta
	default:
		return displayer.ColorRed
	}
}

// getConnectionTypeColor returns the appropriate color for connection type
func (displayer *NetworkMonitorDisplayer) getConnectionTypeColor(connType string) string {
	if !displayer.ShowColors {
		return ""
	}

	switch connType {
	case "TCP":
		return displayer.ColorGreen
	case "UDP":
		return displayer.ColorBlue
	case "Unix":
		return displayer.ColorMagenta
	default:
		return displayer.ColorWhite
	}
}

// getNetworkSpeedColor returns the appropriate color for network speed
func (displayer *NetworkMonitorDisplayer) getNetworkSpeedColor(speed float64) string {
	if !displayer.ShowColors {
		return ""
	}

	switch {
	case speed < 10:
		return displayer.ColorGreen
	case speed < 50:
		return displayer.ColorYellow
	case speed < 100:
		return displayer.ColorMagenta
	default:
		return displayer.ColorRed
	}
}
