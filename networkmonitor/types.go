package networkmonitor

import "time"

// NetworkInterfaceInfo represents information about a network interface
type NetworkInterfaceInfo struct {
	Name         string `json:"name"`          // Interface name (e.g., eth0, wlan0)
	DisplayName  string `json:"display_name"`  // Human-readable name
	Type         string `json:"type"`           // Interface type (Ethernet, WiFi, etc.)
	Status       string `json:"status"`        // Interface status (up, down, unknown)
	MTU          int    `json:"mtu"`           // Maximum Transmission Unit
	Speed        uint64 `json:"speed"`         // Interface speed in Mbps
	MACAddress   string `json:"mac_address"`   // MAC address
	IPAddress    string `json:"ip_address"`   // Primary IP address
	SubnetMask   string `json:"subnet_mask"`   // Subnet mask
	Gateway      string `json:"gateway"`       // Default gateway
	DNSServers   []string `json:"dns_servers"` // DNS servers
	IsUp         bool   `json:"is_up"`         // Whether interface is up
	IsLoopback   bool   `json:"is_loopback"`   // Whether interface is loopback
	IsVirtual    bool   `json:"is_virtual"`    // Whether interface is virtual
}

// NetworkIOInfo represents network I/O statistics for an interface
type NetworkIOInfo struct {
	InterfaceName string  `json:"interface_name"` // Interface name
	BytesSent     uint64  `json:"bytes_sent"`     // Total bytes sent
	BytesRecv     uint64  `json:"bytes_recv"`     // Total bytes received
	PacketsSent   uint64  `json:"packets_sent"`   // Total packets sent
	PacketsRecv   uint64  `json:"packets_recv"`   // Total packets received
	SendSpeed     float64 `json:"send_speed"`     // Current send speed (Mbps)
	RecvSpeed     float64 `json:"recv_speed"`     // Current receive speed (Mbps)
	TotalSpeed    float64 `json:"total_speed"`    // Total throughput (Mbps)
	SendErrors    uint64  `json:"send_errors"`   // Send errors
	RecvErrors    uint64  `json:"recv_errors"`   // Receive errors
	DropIn        uint64  `json:"drop_in"`       // Incoming packets dropped
	DropOut       uint64  `json:"drop_out"`      // Outgoing packets dropped
	Utilization   float64 `json:"utilization"`   // Interface utilization percentage
}

// NetworkConnectionInfo represents information about network connections
type NetworkConnectionInfo struct {
	LocalAddress  string `json:"local_address"`   // Local address:port
	RemoteAddress string `json:"remote_address"`   // Remote address:port
	Status        string `json:"status"`          // Connection status
	Type          string `json:"type"`            // Connection type (TCP, UDP)
	PID           int32  `json:"pid"`             // Process ID
	ProcessName   string `json:"process_name"`   // Process name
	User          string `json:"user"`            // Process owner
	State         string `json:"state"`          // Connection state
	Family        string `json:"family"`          // Address family (IPv4, IPv6)
}

// NetworkProcessInfo represents network usage information for a specific process
type NetworkProcessInfo struct {
	PID           int32   `json:"pid"`            // Process ID
	Name          string  `json:"name"`           // Process name
	BytesSent     uint64  `json:"bytes_sent"`     // Bytes sent by process
	BytesRecv     uint64  `json:"bytes_recv"`     // Bytes received by process
	SendSpeed     float64 `json:"send_speed"`     // Send speed (Mbps)
	RecvSpeed     float64 `json:"recv_speed"`     // Receive speed (Mbps)
	TotalSpeed    float64 `json:"total_speed"`    // Total throughput (Mbps)
	Connections   int     `json:"connections"`    // Number of connections
	Status        string  `json:"status"`        // Process status
	User          string  `json:"user"`           // Process owner
}

// NetworkLatencyInfo represents network latency information
type NetworkLatencyInfo struct {
	Target        string  `json:"target"`         // Target host/IP
	Latency       float64 `json:"latency"`        // Latency in milliseconds
	PacketLoss    float64 `json:"packet_loss"`    // Packet loss percentage
	Jitter        float64 `json:"jitter"`         // Jitter in milliseconds
	Status        string  `json:"status"`         // Connection status
	LastChecked   time.Time `json:"last_checked"` // Last check time
}

// NetworkBandwidthInfo represents bandwidth usage information
type NetworkBandwidthInfo struct {
	TotalBandwidth    float64 `json:"total_bandwidth"`     // Total available bandwidth (Mbps)
	UsedBandwidth     float64 `json:"used_bandwidth"`      // Currently used bandwidth (Mbps)
	AvailableBandwidth float64 `json:"available_bandwidth"` // Available bandwidth (Mbps)
	Utilization       float64 `json:"utilization"`         // Bandwidth utilization percentage
	PeakUsage         float64 `json:"peak_usage"`          // Peak usage (Mbps)
	AverageUsage      float64 `json:"average_usage"`       // Average usage (Mbps)
}

// NetworkMonitorData represents comprehensive network monitoring data
type NetworkMonitorData struct {
	// Network interfaces
	Interfaces []NetworkInterfaceInfo `json:"interfaces"` // Available network interfaces

	// I/O statistics for each interface
	InterfaceIO []NetworkIOInfo `json:"interface_io"` // I/O statistics per interface

	// Network connections
	Connections []NetworkConnectionInfo `json:"connections"` // Active network connections

	// Top processes by network usage
	TopProcesses []NetworkProcessInfo `json:"top_processes"` // Top network-consuming processes

	// Network latency information
	LatencyInfo []NetworkLatencyInfo `json:"latency_info"` // Latency to various targets

	// Bandwidth information
	BandwidthInfo NetworkBandwidthInfo `json:"bandwidth_info"` // Bandwidth usage information

	// Overall network statistics
	TotalBytesSent    uint64  `json:"total_bytes_sent"`     // Total bytes sent across all interfaces
	TotalBytesRecv    uint64  `json:"total_bytes_recv"`     // Total bytes received across all interfaces
	TotalPacketsSent  uint64  `json:"total_packets_sent"`   // Total packets sent across all interfaces
	TotalPacketsRecv  uint64  `json:"total_packets_recv"`   // Total packets received across all interfaces
	TotalSendSpeed    float64 `json:"total_send_speed"`     // Total send speed across all interfaces
	TotalRecvSpeed    float64 `json:"total_recv_speed"`     // Total receive speed across all interfaces
	TotalThroughput  float64 `json:"total_throughput"`     // Total network throughput

	// Network performance metrics
	AverageLatency    float64 `json:"average_latency"`      // Average latency across all targets
	PacketLossRate    float64 `json:"packet_loss_rate"`     // Overall packet loss rate
	NetworkUtilization float64 `json:"network_utilization"` // Overall network utilization

	// Network alerts and warnings
	NetworkStatus     string `json:"network_status"`      // Overall network status (Normal, Warning, Critical)
	HighLatencyWarning bool  `json:"high_latency_warning"` // High latency warning
	PacketLossWarning  bool  `json:"packet_loss_warning"`  // Packet loss warning
	BandwidthWarning   bool  `json:"bandwidth_warning"`    // Bandwidth usage warning
	ConnectionWarning  bool  `json:"connection_warning"`    // Connection issues warning

	// Monitoring configuration
	RefreshInterval time.Duration `json:"refresh_interval"` // How often data is refreshed
	IsMonitoring    bool          `json:"is_monitoring"`    // Whether monitoring is active

	// Timestamps
	Timestamp time.Time     `json:"timestamp"` // When this data was collected
	Uptime    time.Duration `json:"uptime"`    // System uptime
}

// NetworkMonitorConfig represents configuration options for network monitoring
type NetworkMonitorConfig struct {
	// Monitoring settings
	RefreshInterval     time.Duration `json:"refresh_interval"`     // How often to refresh data
	MaxProcesses        int           `json:"max_processes"`        // Maximum number of processes to track
	MaxConnections      int           `json:"max_connections"`      // Maximum number of connections to track
	LatencyWarning      float64       `json:"latency_warning"`      // Latency warning threshold (ms)
	LatencyCritical     float64       `json:"latency_critical"`     // Latency critical threshold (ms)
	PacketLossWarning   float64       `json:"packet_loss_warning"`  // Packet loss warning threshold (%)
	BandwidthWarning    float64       `json:"bandwidth_warning"`   // Bandwidth warning threshold (%)
	ConnectionTimeout   time.Duration `json:"connection_timeout"`   // Connection timeout

	// Display settings
	ShowInterfaces    bool `json:"show_interfaces"`     // Whether to show interface information
	ShowIO           bool `json:"show_io"`            // Whether to show I/O statistics
	ShowConnections  bool `json:"show_connections"`   // Whether to show connection information
	ShowProcesses    bool `json:"show_processes"`     // Whether to show process information
	ShowLatency      bool `json:"show_latency"`       // Whether to show latency information
	ShowBandwidth    bool `json:"show_bandwidth"`     // Whether to show bandwidth information
	ShowPerformance  bool `json:"show_performance"`   // Whether to show performance metrics

	// Export settings
	ExportToFile   bool          `json:"export_to_file"`  // Whether to export data to file
	ExportInterval time.Duration `json:"export_interval"` // How often to export data
	ExportFormat   string        `json:"export_format"`   // Export format (json, csv, txt)

	// Filter settings
	MinNetworkUsage     float64 `json:"min_network_usage"`     // Minimum network usage to show process
	ProcessNameFilter   string  `json:"process_name_filter"`   // Filter processes by name
	InterfaceFilter     string  `json:"interface_filter"`      // Filter specific interfaces
	ConnectionTypeFilter string  `json:"connection_type_filter"` // Filter connection types
	LatencyTargets      []string `json:"latency_targets"`      // Targets for latency monitoring
}

// NetworkUsageHistory represents historical network usage data for graphing
type NetworkUsageHistory struct {
	// Time series data
	Timestamps     []time.Time `json:"timestamps"`      // Time points for data
	TotalSent      []float64  `json:"total_sent"`      // Total bytes sent over time
	TotalRecv      []float64  `json:"total_recv"`      // Total bytes received over time
	SendSpeed      []float64  `json:"send_speed"`      // Send speed over time
	RecvSpeed      []float64  `json:"recv_speed"`      // Receive speed over time
	Throughput     []float64  `json:"throughput"`      // Total throughput over time
	Latency        []float64  `json:"latency"`         // Average latency over time
	Utilization    []float64  `json:"utilization"`     // Network utilization over time

	// Configuration
	MaxDataPoints int `json:"max_data_points"` // Maximum number of data points to store
	DataPointCount int `json:"data_point_count"` // Current number of data points
}
