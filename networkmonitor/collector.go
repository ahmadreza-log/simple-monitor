package networkmonitor

import (
	"fmt"
	"net"
	"sort"
	"time"

	netutil "github.com/shirou/gopsutil/v3/net"
	"github.com/shirou/gopsutil/v3/process"
)

// NetworkMonitorCollector handles the collection of network monitoring data
// This struct provides methods to gather real-time network metrics and process information
type NetworkMonitorCollector struct {
	// Configuration
	config        *NetworkMonitorConfig
	lastTimestamp time.Time

	// Process tracking
	processCache    map[int32]*NetworkProcessInfo
	lastProcessTime map[int32]time.Time

	// History tracking
	history *NetworkUsageHistory
}

// NewNetworkMonitorCollector creates a new instance of NetworkMonitorCollector
// with default configuration values
func NewNetworkMonitorCollector() *NetworkMonitorCollector {
	config := &NetworkMonitorConfig{
		RefreshInterval:     1 * time.Second,
		MaxProcesses:        20,
		MaxConnections:      100,
		LatencyWarning:      100.0,
		LatencyCritical:     200.0,
		PacketLossWarning:   5.0,
		BandwidthWarning:    80.0,
		ConnectionTimeout:    5 * time.Second,
		ShowInterfaces:      true,
		ShowIO:              true,
		ShowConnections:     true,
		ShowProcesses:       true,
		ShowLatency:         true,
		ShowBandwidth:       true,
		ShowPerformance:     true,
		ExportToFile:        true,
		ExportInterval:      1 * time.Hour,
		ExportFormat:        "json",
		MinNetworkUsage:     1.0,
		ProcessNameFilter:   "",
		InterfaceFilter:     "",
		ConnectionTypeFilter: "",
		LatencyTargets:      []string{"8.8.8.8", "1.1.1.1", "google.com"},
	}

	return &NetworkMonitorCollector{
		config:          config,
		lastTimestamp:   time.Now(),
		processCache:    make(map[int32]*NetworkProcessInfo),
		lastProcessTime: make(map[int32]time.Time),
		history: &NetworkUsageHistory{
			MaxDataPoints:  100,
			DataPointCount: 0,
		},
	}
}

// CollectNetworkMonitorData gathers comprehensive network monitoring data
// This is the main method that collects all available network metrics
func (collector *NetworkMonitorCollector) CollectNetworkMonitorData() (*NetworkMonitorData, error) {
	data := &NetworkMonitorData{
		Timestamp:       time.Now(),
		RefreshInterval: collector.config.RefreshInterval,
		IsMonitoring:    true,
	}

	// Collect interface information
	if collector.config.ShowInterfaces {
		if err := collector.collectInterfaceInfo(data); err != nil {
			return nil, fmt.Errorf("failed to collect interface info: %w", err)
		}
	}

	// Collect I/O statistics
	if collector.config.ShowIO {
		if err := collector.collectIOInfo(data); err != nil {
			return nil, fmt.Errorf("failed to collect I/O info: %w", err)
		}
	}

	// Collect connection information
	if collector.config.ShowConnections {
		if err := collector.collectConnectionInfo(data); err != nil {
			return nil, fmt.Errorf("failed to collect connection info: %w", err)
		}
	}

	// Collect process information
	if collector.config.ShowProcesses {
		if err := collector.collectProcessInfo(data); err != nil {
			return nil, fmt.Errorf("failed to collect process info: %w", err)
		}
	}

	// Collect latency information
	if collector.config.ShowLatency {
		if err := collector.collectLatencyInfo(data); err != nil {
			return nil, fmt.Errorf("failed to collect latency info: %w", err)
		}
	}

	// Collect bandwidth information
	if collector.config.ShowBandwidth {
		collector.collectBandwidthInfo(data)
	}

	// Calculate performance metrics
	if collector.config.ShowPerformance {
		collector.calculatePerformanceMetrics(data)
	}

	// Analyze network status and alerts
	collector.analyzeNetworkStatus(data)

	// Update history
	collector.updateHistory(data)

	return data, nil
}

// collectInterfaceInfo gathers network interface information
func (collector *NetworkMonitorCollector) collectInterfaceInfo(data *NetworkMonitorData) error {
	// Get network interfaces
	interfaces, err := netutil.Interfaces()
	if err != nil {
		return fmt.Errorf("failed to get network interfaces: %w", err)
	}

	var interfaceInfos []NetworkInterfaceInfo

	for _, iface := range interfaces {
		// Skip if interface filter is specified and doesn't match
		if collector.config.InterfaceFilter != "" && iface.Name != collector.config.InterfaceFilter {
			continue
		}

		// Get interface addresses
		addrs := iface.Addrs
		if len(addrs) == 0 {
			continue // Skip interfaces with no addresses
		}

		var ipAddress, subnetMask string
		if len(addrs) > 0 {
			ipAddress = addrs[0].Addr
			subnetMask = addrs[0].Addr // Simplified
		}

		interfaceInfo := NetworkInterfaceInfo{
			Name:        iface.Name,
			DisplayName: iface.Name,
			Type:        collector.getInterfaceType(iface.Name),
			Status:      collector.getInterfaceStatus(iface.Flags),
			MTU:         iface.MTU,
			Speed:       1000, // Default speed in Mbps (would need platform-specific implementation)
			MACAddress:  iface.HardwareAddr,
			IPAddress:   ipAddress,
			SubnetMask:  subnetMask,
			Gateway:     "", // Would need platform-specific implementation
			DNSServers:  []string{}, // Would need platform-specific implementation
			IsUp:        collector.isInterfaceUp(iface.Flags),
			IsLoopback:  collector.isLoopbackInterface(iface.Name),
			IsVirtual:   collector.isVirtualInterface(iface.Name),
		}

		interfaceInfos = append(interfaceInfos, interfaceInfo)
	}

	data.Interfaces = interfaceInfos
	return nil
}

// collectIOInfo gathers network I/O statistics
func (collector *NetworkMonitorCollector) collectIOInfo(data *NetworkMonitorData) error {
	// Get I/O counters
	ioCounters, err := netutil.IOCounters(true)
	if err != nil {
		return fmt.Errorf("failed to get I/O counters: %w", err)
	}

	var interfaceIOs []NetworkIOInfo
	var totalBytesSent, totalBytesRecv, totalPacketsSent, totalPacketsRecv uint64
	var totalSendSpeed, totalRecvSpeed float64

	for _, counter := range ioCounters {
		// Skip if interface filter is specified and doesn't match
		if collector.config.InterfaceFilter != "" && counter.Name != collector.config.InterfaceFilter {
			continue
		}

		// Calculate speeds (simplified - would need time-based calculation for accurate speeds)
		sendSpeed := float64(counter.BytesSent) / (1024 * 1024) // Convert to Mbps (simplified)
		recvSpeed := float64(counter.BytesRecv) / (1024 * 1024) // Convert to Mbps (simplified)
		totalSpeed := sendSpeed + recvSpeed
		utilization := (totalSpeed / 1000.0) * 100 // Simplified utilization calculation

		interfaceIO := NetworkIOInfo{
			InterfaceName: counter.Name,
			BytesSent:     counter.BytesSent,
			BytesRecv:     counter.BytesRecv,
			PacketsSent:   counter.PacketsSent,
			PacketsRecv:   counter.PacketsRecv,
			SendSpeed:     sendSpeed,
			RecvSpeed:     recvSpeed,
			TotalSpeed:    totalSpeed,
			SendErrors:    counter.Errin + counter.Errout,
			RecvErrors:    counter.Errin,
			DropIn:        counter.Dropin,
			DropOut:       counter.Dropout,
			Utilization:   utilization,
		}

		interfaceIOs = append(interfaceIOs, interfaceIO)

		// Add to totals
		totalBytesSent += counter.BytesSent
		totalBytesRecv += counter.BytesRecv
		totalPacketsSent += counter.PacketsSent
		totalPacketsRecv += counter.PacketsRecv
		totalSendSpeed += sendSpeed
		totalRecvSpeed += recvSpeed
	}

	data.InterfaceIO = interfaceIOs
	data.TotalBytesSent = totalBytesSent
	data.TotalBytesRecv = totalBytesRecv
	data.TotalPacketsSent = totalPacketsSent
	data.TotalPacketsRecv = totalPacketsRecv
	data.TotalSendSpeed = totalSendSpeed
	data.TotalRecvSpeed = totalRecvSpeed
	data.TotalThroughput = totalSendSpeed + totalRecvSpeed

	return nil
}

// collectConnectionInfo gathers network connection information
func (collector *NetworkMonitorCollector) collectConnectionInfo(data *NetworkMonitorData) error {
	// Get network connections
	connections, err := netutil.Connections("all")
	if err != nil {
		return fmt.Errorf("failed to get network connections: %w", err)
	}

	var connectionInfos []NetworkConnectionInfo
	connectionCount := 0

	for _, conn := range connections {
		// Limit number of connections
		if connectionCount >= collector.config.MaxConnections {
			break
		}

		// Skip if connection type filter is specified and doesn't match
		connType := collector.getConnectionType(conn.Type)
		if collector.config.ConnectionTypeFilter != "" && connType != collector.config.ConnectionTypeFilter {
			continue
		}

		// Get process information
		processName := "Unknown"
		user := "Unknown"
		if conn.Pid > 0 {
			if proc, err := process.NewProcess(conn.Pid); err == nil {
				if name, err := proc.Name(); err == nil {
					processName = name
				}
				if username, err := proc.Username(); err == nil {
					user = username
				}
			}
		}

		connectionInfo := NetworkConnectionInfo{
			LocalAddress:  fmt.Sprintf("%s:%d", conn.Laddr.IP, conn.Laddr.Port),
			RemoteAddress: fmt.Sprintf("%s:%d", conn.Raddr.IP, conn.Raddr.Port),
			Status:        conn.Status,
			Type:          connType,
			PID:           conn.Pid,
			ProcessName:     processName,
			User:          user,
			State:         conn.Status,
			Family:        collector.getConnectionFamily(conn.Family),
		}

		connectionInfos = append(connectionInfos, connectionInfo)
		connectionCount++
	}

	data.Connections = connectionInfos
	return nil
}

// collectProcessInfo gathers top network-consuming processes
func (collector *NetworkMonitorCollector) collectProcessInfo(data *NetworkMonitorData) error {
	// Get all processes
	processes, err := process.Processes()
	if err != nil {
		return fmt.Errorf("failed to get processes: %w", err)
	}

	var networkProcesses []NetworkProcessInfo

	// Collect network I/O information for each process
	for _, p := range processes {
		// Get process I/O info
		ioInfo, err := p.IOCounters()
		if err != nil {
			continue // Skip processes we can't access
		}

		// Get process name
		name, err := p.Name()
		if err != nil {
			name = "Unknown"
		}

		// Get process status
		status, err := p.Status()
		if err != nil {
			status = []string{"Unknown"}
		}

		// Get process user
		user, err := p.Username()
		if err != nil {
			user = "Unknown"
		}

		// Calculate network metrics
		sendSpeed := float64(ioInfo.WriteBytes) / (1024 * 1024) // Convert to Mbps (simplified)
		recvSpeed := float64(ioInfo.ReadBytes) / (1024 * 1024) // Convert to Mbps (simplified)
		totalSpeed := sendSpeed + recvSpeed

		// Filter by minimum network usage
		if totalSpeed < collector.config.MinNetworkUsage {
			continue
		}

		// Filter by process name if specified
		if collector.config.ProcessNameFilter != "" && name != collector.config.ProcessNameFilter {
			continue
		}

		// Count connections for this process
		connections := 0
		for _, conn := range data.Connections {
			if conn.PID == p.Pid {
				connections++
			}
		}

		processInfo := NetworkProcessInfo{
			PID:         p.Pid,
			Name:        name,
			BytesSent:   ioInfo.WriteBytes,
			BytesRecv:   ioInfo.ReadBytes,
			SendSpeed:   sendSpeed,
			RecvSpeed:   recvSpeed,
			TotalSpeed:  totalSpeed,
			Connections: connections,
			Status:      status[0],
			User:        user,
		}

		networkProcesses = append(networkProcesses, processInfo)
	}

	// Sort by total network usage
	sort.Slice(networkProcesses, func(i, j int) bool {
		return networkProcesses[i].TotalSpeed > networkProcesses[j].TotalSpeed
	})

	// Limit to max processes
	if len(networkProcesses) > collector.config.MaxProcesses {
		networkProcesses = networkProcesses[:collector.config.MaxProcesses]
	}

	data.TopProcesses = networkProcesses
	return nil
}

// collectLatencyInfo gathers network latency information
func (collector *NetworkMonitorCollector) collectLatencyInfo(data *NetworkMonitorData) error {
	var latencyInfos []NetworkLatencyInfo

	for _, target := range collector.config.LatencyTargets {
		// Measure latency (simplified implementation)
		latency, packetLoss, status := collector.measureLatency(target)

		latencyInfo := NetworkLatencyInfo{
			Target:      target,
			Latency:     latency,
			PacketLoss:  packetLoss,
			Jitter:      0.0, // Would need multiple measurements
			Status:      status,
			LastChecked: time.Now(),
		}

		latencyInfos = append(latencyInfos, latencyInfo)
	}

	data.LatencyInfo = latencyInfos
	return nil
}

// collectBandwidthInfo calculates bandwidth usage information
func (collector *NetworkMonitorCollector) collectBandwidthInfo(data *NetworkMonitorData) {
	// Calculate bandwidth information
	totalBandwidth := float64(len(data.Interfaces)) * 1000.0 // Assume 1Gbps per interface
	usedBandwidth := data.TotalThroughput
	availableBandwidth := totalBandwidth - usedBandwidth
	utilization := (usedBandwidth / totalBandwidth) * 100

	data.BandwidthInfo = NetworkBandwidthInfo{
		TotalBandwidth:     totalBandwidth,
		UsedBandwidth:      usedBandwidth,
		AvailableBandwidth: availableBandwidth,
		Utilization:        utilization,
		PeakUsage:          usedBandwidth, // Simplified
		AverageUsage:       usedBandwidth, // Simplified
	}
}

// calculatePerformanceMetrics calculates overall performance metrics
func (collector *NetworkMonitorCollector) calculatePerformanceMetrics(data *NetworkMonitorData) {
	// Calculate average latency
	if len(data.LatencyInfo) > 0 {
		var totalLatency float64
		for _, latency := range data.LatencyInfo {
			totalLatency += latency.Latency
		}
		data.AverageLatency = totalLatency / float64(len(data.LatencyInfo))
	}

	// Calculate packet loss rate
	if len(data.LatencyInfo) > 0 {
		var totalPacketLoss float64
		for _, latency := range data.LatencyInfo {
			totalPacketLoss += latency.PacketLoss
		}
		data.PacketLossRate = totalPacketLoss / float64(len(data.LatencyInfo))
	}

	// Calculate network utilization
	data.NetworkUtilization = data.BandwidthInfo.Utilization
}

// analyzeNetworkStatus analyzes network status and sets alerts
func (collector *NetworkMonitorCollector) analyzeNetworkStatus(data *NetworkMonitorData) {
	// Analyze latency
	for _, latency := range data.LatencyInfo {
		if latency.Latency >= collector.config.LatencyCritical {
			data.NetworkStatus = "Critical"
			data.HighLatencyWarning = true
			break
		} else if latency.Latency >= collector.config.LatencyWarning {
			data.HighLatencyWarning = true
			if data.NetworkStatus == "Normal" {
				data.NetworkStatus = "Warning"
			}
		}
	}

	// Analyze packet loss
	for _, latency := range data.LatencyInfo {
		if latency.PacketLoss >= collector.config.PacketLossWarning {
			data.PacketLossWarning = true
			if data.NetworkStatus == "Normal" {
				data.NetworkStatus = "Warning"
			}
			break
		}
	}

	// Analyze bandwidth usage
	if data.BandwidthInfo.Utilization >= collector.config.BandwidthWarning {
		data.BandwidthWarning = true
		if data.NetworkStatus == "Normal" {
			data.NetworkStatus = "Warning"
		}
	}

	// Analyze connection issues
	if len(data.Connections) == 0 {
		data.ConnectionWarning = true
		if data.NetworkStatus == "Normal" {
			data.NetworkStatus = "Warning"
		}
	}

	// Set default status if no issues
	if data.NetworkStatus == "" {
		data.NetworkStatus = "Normal"
	}
}

// updateHistory updates the network usage history
func (collector *NetworkMonitorCollector) updateHistory(data *NetworkMonitorData) {
	now := time.Now()
	
	// Add new data point
	collector.history.Timestamps = append(collector.history.Timestamps, now)
	collector.history.TotalSent = append(collector.history.TotalSent, float64(data.TotalBytesSent))
	collector.history.TotalRecv = append(collector.history.TotalRecv, float64(data.TotalBytesRecv))
	collector.history.SendSpeed = append(collector.history.SendSpeed, data.TotalSendSpeed)
	collector.history.RecvSpeed = append(collector.history.RecvSpeed, data.TotalRecvSpeed)
	collector.history.Throughput = append(collector.history.Throughput, data.TotalThroughput)
	collector.history.Latency = append(collector.history.Latency, data.AverageLatency)
	collector.history.Utilization = append(collector.history.Utilization, data.NetworkUtilization)

	// Limit history size
	if len(collector.history.Timestamps) > collector.history.MaxDataPoints {
		collector.history.Timestamps = collector.history.Timestamps[1:]
		collector.history.TotalSent = collector.history.TotalSent[1:]
		collector.history.TotalRecv = collector.history.TotalRecv[1:]
		collector.history.SendSpeed = collector.history.SendSpeed[1:]
		collector.history.RecvSpeed = collector.history.RecvSpeed[1:]
		collector.history.Throughput = collector.history.Throughput[1:]
		collector.history.Latency = collector.history.Latency[1:]
		collector.history.Utilization = collector.history.Utilization[1:]
	}

	collector.history.DataPointCount = len(collector.history.Timestamps)
}

// Helper methods

// getInterfaceType returns the interface type based on name
func (collector *NetworkMonitorCollector) getInterfaceType(name string) string {
	if name[:2] == "wl" || name[:4] == "wlan" {
		return "WiFi"
	} else if name[:2] == "et" || name[:4] == "eth" {
		return "Ethernet"
	} else if name[:2] == "lo" {
		return "Loopback"
	} else if name[:2] == "vm" || name[:2] == "vb" {
		return "Virtual"
	}
	return "Unknown"
}

// getInterfaceStatus returns the interface status based on flags
func (collector *NetworkMonitorCollector) getInterfaceStatus(flags []string) string {
	for _, flag := range flags {
		if flag == "up" {
			return "up"
		}
	}
	return "down"
}

// isInterfaceUp checks if interface is up
func (collector *NetworkMonitorCollector) isInterfaceUp(flags []string) bool {
	for _, flag := range flags {
		if flag == "up" {
			return true
		}
	}
	return false
}

// isLoopbackInterface checks if interface is loopback
func (collector *NetworkMonitorCollector) isLoopbackInterface(name string) bool {
	return name == "lo" || name[:2] == "lo"
}

// getConnectionType returns the connection type as string
func (collector *NetworkMonitorCollector) getConnectionType(connType uint32) string {
	switch connType {
	case 1:
		return "TCP"
	case 2:
		return "UDP"
	case 3:
		return "Unix"
	default:
		return "Unknown"
	}
}

// getConnectionFamily returns the connection family as string
func (collector *NetworkMonitorCollector) getConnectionFamily(family uint32) string {
	switch family {
	case 2:
		return "IPv4"
	case 10:
		return "IPv6"
	case 1:
		return "Unix"
	default:
		return "Unknown"
	}
}

// isVirtualInterface checks if an interface is virtual
func (collector *NetworkMonitorCollector) isVirtualInterface(name string) bool {
	virtualPrefixes := []string{"vm", "vb", "veth", "docker", "br-", "virbr"}
	for _, prefix := range virtualPrefixes {
		if len(name) >= len(prefix) && name[:len(prefix)] == prefix {
			return true
		}
	}
	return false
}

// measureLatency measures latency to a target (simplified implementation)
func (collector *NetworkMonitorCollector) measureLatency(target string) (float64, float64, string) {
	// Simplified latency measurement
	// In a real implementation, this would use ping or similar
	start := time.Now()
	conn, err := net.DialTimeout("tcp", target+":80", collector.config.ConnectionTimeout)
	if err != nil {
		return 0.0, 100.0, "Failed"
	}
	conn.Close()
	latency := float64(time.Since(start).Nanoseconds()) / 1000000.0 // Convert to milliseconds
	
	if latency < collector.config.LatencyWarning {
		return latency, 0.0, "Good"
	} else if latency < collector.config.LatencyCritical {
		return latency, 0.0, "Warning"
	} else {
		return latency, 0.0, "Critical"
	}
}

// GetNetworkUsageHistory returns the current network usage history
func (collector *NetworkMonitorCollector) GetNetworkUsageHistory() *NetworkUsageHistory {
	return collector.history
}

// GetConfig returns the current configuration
func (collector *NetworkMonitorCollector) GetConfig() *NetworkMonitorConfig {
	return collector.config
}

// UpdateConfig updates the collector configuration
func (collector *NetworkMonitorCollector) UpdateConfig(config *NetworkMonitorConfig) {
	collector.config = config
}
