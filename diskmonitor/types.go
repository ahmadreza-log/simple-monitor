package diskmonitor

import "time"

// DiskPartitionInfo represents information about a disk partition
type DiskPartitionInfo struct {
	Device      string `json:"device"`       // Device name (e.g., /dev/sda1)
	Mountpoint  string `json:"mountpoint"`   // Mount point (e.g., /, /home)
	Fstype      string `json:"fstype"`       // File system type (e.g., ext4, ntfs)
	Total       uint64 `json:"total"`        // Total size in bytes
	Free        uint64 `json:"free"`         // Free space in bytes
	Used        uint64 `json:"used"`         // Used space in bytes
	UsagePercent float64 `json:"usage_percent"` // Usage percentage
	InodesTotal uint64 `json:"inodes_total"` // Total inodes
	InodesFree  uint64 `json:"inodes_free"`   // Free inodes
	InodesUsed  uint64 `json:"inodes_used"`  // Used inodes
}

// DiskIOInfo represents disk I/O statistics
type DiskIOInfo struct {
	DeviceName    string  `json:"device_name"`    // Device name
	ReadCount     uint64  `json:"read_count"`    // Number of reads
	WriteCount    uint64  `json:"write_count"`   // Number of writes
	ReadBytes     uint64  `json:"read_bytes"`    // Bytes read
	WriteBytes    uint64  `json:"write_bytes"`   // Bytes written
	ReadTime      uint64  `json:"read_time"`     // Time spent reading (ms)
	WriteTime     uint64  `json:"write_time"`    // Time spent writing (ms)
	ReadSpeed     float64 `json:"read_speed"`     // Read speed (MB/s)
	WriteSpeed    float64 `json:"write_speed"`   // Write speed (MB/s)
	IOPS          float64 `json:"iops"`           // I/O operations per second
	Utilization   float64 `json:"utilization"`   // Disk utilization percentage
}

// DiskTemperatureInfo represents disk temperature information
type DiskTemperatureInfo struct {
	DeviceName    string  `json:"device_name"`    // Device name
	Temperature   float64 `json:"temperature"`     // Current temperature in Celsius
	MaxTemperature float64 `json:"max_temperature"` // Maximum safe temperature
	Status        string  `json:"status"`         // Temperature status (Normal, Warning, Critical)
}

// DiskHealthInfo represents disk health and SMART data
type DiskHealthInfo struct {
	DeviceName     string `json:"device_name"`      // Device name
	HealthStatus   string `json:"health_status"`    // Overall health status
	PowerOnHours   uint64 `json:"power_on_hours"`   // Power-on hours
	PowerCycleCount uint64 `json:"power_cycle_count"` // Power cycle count
	ReallocatedSectors uint64 `json:"reallocated_sectors"` // Reallocated sectors count
	PendingSectors  uint64 `json:"pending_sectors"`  // Pending sectors count
	UncorrectableSectors uint64 `json:"uncorrectable_sectors"` // Uncorrectable sectors count
	Temperature     float64 `json:"temperature"`     // Current temperature
	WearLeveling    float64 `json:"wear_leveling"`   // Wear leveling percentage (SSD)
}

// DiskProcessInfo represents disk usage information for a specific process
type DiskProcessInfo struct {
	PID           int32   `json:"pid"`            // Process ID
	Name          string  `json:"name"`           // Process name
	ReadBytes     uint64  `json:"read_bytes"`    // Bytes read by process
	WriteBytes    uint64  `json:"write_bytes"`    // Bytes written by process
	ReadSpeed     float64 `json:"read_speed"`     // Read speed (MB/s)
	WriteSpeed    float64 `json:"write_speed"`   // Write speed (MB/s)
	TotalIO       uint64  `json:"total_io"`       // Total I/O operations
	IOPS          float64 `json:"iops"`           // I/O operations per second
	Status        string  `json:"status"`         // Process status
	User          string  `json:"user"`           // Process owner
}

// DiskMonitorData represents comprehensive disk monitoring data
type DiskMonitorData struct {
	// Overall disk statistics
	TotalSpace     uint64 `json:"total_space"`      // Total disk space across all partitions
	UsedSpace      uint64 `json:"used_space"`       // Used disk space across all partitions
	FreeSpace      uint64 `json:"free_space"`       // Free disk space across all partitions
	UsagePercent   float64 `json:"usage_percent"`  // Overall disk usage percentage

	// Disk partitions information
	Partitions []DiskPartitionInfo `json:"partitions"` // Information about each partition

	// Disk I/O statistics
	DiskIO []DiskIOInfo `json:"disk_io"` // I/O statistics for each disk

	// Disk temperature information
	DiskTemperatures []DiskTemperatureInfo `json:"disk_temperatures"` // Temperature for each disk

	// Disk health information
	DiskHealth []DiskHealthInfo `json:"disk_health"` // Health information for each disk

	// Top processes by disk usage
	TopProcesses []DiskProcessInfo `json:"top_processes"` // Top disk-consuming processes

	// Performance metrics
	TotalReadSpeed  float64 `json:"total_read_speed"`  // Total read speed across all disks
	TotalWriteSpeed float64 `json:"total_write_speed"` // Total write speed across all disks
	AverageIOPS     float64 `json:"average_iops"`      // Average I/O operations per second
	DiskUtilization float64 `json:"disk_utilization"` // Overall disk utilization

	// Disk alerts and warnings
	DiskStatus       string `json:"disk_status"`        // Overall disk status (Normal, Warning, Critical)
	LowSpaceWarning  bool   `json:"low_space_warning"`   // Low disk space warning
	HighTempWarning  bool   `json:"high_temp_warning"`   // High temperature warning
	HealthWarning    bool   `json:"health_warning"`     // Disk health warning
	IOBottleneck     bool   `json:"io_bottleneck"`       // I/O bottleneck detection

	// Monitoring configuration
	RefreshInterval time.Duration `json:"refresh_interval"` // How often data is refreshed
	IsMonitoring    bool          `json:"is_monitoring"`    // Whether monitoring is active

	// Timestamps
	Timestamp time.Time     `json:"timestamp"` // When this data was collected
	Uptime    time.Duration `json:"uptime"`    // System uptime
}

// DiskMonitorConfig represents configuration options for disk monitoring
type DiskMonitorConfig struct {
	// Monitoring settings
	RefreshInterval     time.Duration `json:"refresh_interval"`     // How often to refresh data
	MaxProcesses        int           `json:"max_processes"`        // Maximum number of processes to track
	LowSpaceWarning     float64       `json:"low_space_warning"`   // Low space warning threshold (percentage)
	LowSpaceCritical    float64       `json:"low_space_critical"`  // Low space critical threshold (percentage)
	TempWarning         float64       `json:"temp_warning"`        // Temperature warning threshold (Celsius)
	TempCritical       float64       `json:"temp_critical"`        // Temperature critical threshold (Celsius)
	IOBottleneckThreshold float64    `json:"io_bottleneck_threshold"` // I/O bottleneck threshold (percentage)

	// Display settings
	ShowPartitions    bool `json:"show_partitions"`     // Whether to show partition information
	ShowIO           bool `json:"show_io"`            // Whether to show I/O statistics
	ShowTemperature  bool `json:"show_temperature"`   // Whether to show temperature information
	ShowHealth       bool `json:"show_health"`       // Whether to show health information
	ShowProcesses    bool `json:"show_processes"`    // Whether to show process information
	ShowPerformance  bool `json:"show_performance"`   // Whether to show performance metrics

	// Export settings
	ExportToFile   bool          `json:"export_to_file"`  // Whether to export data to file
	ExportInterval time.Duration `json:"export_interval"` // How often to export data
	ExportFormat   string        `json:"export_format"`   // Export format (json, csv, txt)

	// Filter settings
	MinIOUsage         float64 `json:"min_io_usage"`         // Minimum I/O usage to show process
	ProcessNameFilter  string  `json:"process_name_filter"`  // Filter processes by name
	DeviceFilter       string  `json:"device_filter"`        // Filter specific devices
	MountpointFilter   string  `json:"mountpoint_filter"`   // Filter specific mountpoints
}

// DiskUsageHistory represents historical disk usage data for graphing
type DiskUsageHistory struct {
	// Time series data
	Timestamps     []time.Time `json:"timestamps"`      // Time points for data
	TotalUsage     []float64  `json:"total_usage"`     // Total disk usage over time
	ReadSpeed      []float64  `json:"read_speed"`      // Read speed over time
	WriteSpeed     []float64  `json:"write_speed"`     // Write speed over time
	IOPS           []float64  `json:"iops"`            // I/O operations per second over time
	Utilization    []float64  `json:"utilization"`     // Disk utilization over time

	// Configuration
	MaxDataPoints int `json:"max_data_points"` // Maximum number of data points to store
	DataPointCount int `json:"data_point_count"` // Current number of data points
}
