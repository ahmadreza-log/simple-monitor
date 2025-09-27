package cpumonitor

import (
	"time"
)

// CPUProcessInfo represents information about a single process and its CPU usage
type CPUProcessInfo struct {
	// Process identification
	PID            int32  `json:"pid"`             // Process ID
	Name           string `json:"name"`            // Process name
	ExecutablePath string `json:"executable_path"` // Full path to executable

	// CPU usage information
	CPUUsagePercent float64 `json:"cpu_usage_percent"` // CPU usage percentage
	CPUUsageTime    uint64  `json:"cpu_usage_time"`    // Total CPU time used (in milliseconds)

	// Process status
	Status   string `json:"status"`   // Process status (Running, Sleeping, etc.)
	Priority int32  `json:"priority"` // Process priority

	// Memory usage (for context)
	MemoryUsage   uint64  `json:"memory_usage"`   // Memory usage in bytes
	MemoryPercent float64 `json:"memory_percent"` // Memory usage percentage

	// Thread information
	ThreadCount int32 `json:"thread_count"` // Number of threads

	// Timestamp
	LastUpdated time.Time `json:"last_updated"` // When this data was last updated
}

// CPUCoreInfo represents information about a single CPU core
type CPUCoreInfo struct {
	// Core identification
	CoreID     int `json:"core_id"`     // Core identifier
	PhysicalID int `json:"physical_id"` // Physical processor ID

	// Core usage statistics
	UsagePercent  float64 `json:"usage_percent"`  // Core usage percentage
	UserPercent   float64 `json:"user_percent"`   // User process usage
	SystemPercent float64 `json:"system_percent"` // System process usage
	IdlePercent   float64 `json:"idle_percent"`   // Idle percentage

	// Core performance
	Frequency   float64 `json:"frequency"`   // Core frequency in MHz
	Temperature float64 `json:"temperature"` // Core temperature in Celsius

	// Core status
	IsOnline        bool `json:"is_online"`        // Whether core is online
	IsHyperthreaded bool `json:"is_hyperthreaded"` // Whether this is a hyperthread

	// Timestamp
	LastUpdated time.Time `json:"last_updated"` // When this data was last updated
}

// CPUMonitorData represents comprehensive CPU monitoring data
type CPUMonitorData struct {
	// Basic CPU information
	ModelName     string `json:"model_name"`     // CPU model name
	VendorID      string `json:"vendor_id"`      // CPU vendor
	Architecture  string `json:"architecture"`   // CPU architecture
	PhysicalCores int    `json:"physical_cores"` // Number of physical cores
	LogicalCores  int    `json:"logical_cores"`  // Number of logical cores

	// Overall CPU statistics
	OverallUsage float64 `json:"overall_usage"` // Overall CPU usage percentage
	UserUsage    float64 `json:"user_usage"`    // User process usage
	SystemUsage  float64 `json:"system_usage"`  // System process usage
	IdleUsage    float64 `json:"idle_usage"`    // Idle usage
	IOWaitUsage  float64 `json:"io_wait_usage"` // I/O wait usage

	// CPU performance metrics
	LoadAverage1Min  float64 `json:"load_1_min"`  // Load average over 1 minute
	LoadAverage5Min  float64 `json:"load_5_min"`  // Load average over 5 minutes
	LoadAverage15Min float64 `json:"load_15_min"` // Load average over 15 minutes

	// Temperature information
	Temperature       float64 `json:"temperature"`        // Overall CPU temperature
	MaxTemperature    float64 `json:"max_temperature"`    // Maximum safe temperature
	TemperatureStatus string  `json:"temperature_status"` // Temperature status (Normal, Warning, Critical)

	// Per-core information
	Cores []CPUCoreInfo `json:"cores"` // Information about each core

	// Top processes by CPU usage
	TopProcesses []CPUProcessInfo `json:"top_processes"` // Top CPU-consuming processes

	// Monitoring configuration
	RefreshInterval time.Duration `json:"refresh_interval"` // How often data is refreshed
	IsMonitoring    bool          `json:"is_monitoring"`    // Whether monitoring is active

	// Timestamps
	Timestamp time.Time     `json:"timestamp"` // When this data was collected
	Uptime    time.Duration `json:"uptime"`    // System uptime
}

// CPUMonitorConfig represents configuration options for CPU monitoring
type CPUMonitorConfig struct {
	// Monitoring settings
	RefreshInterval     time.Duration `json:"refresh_interval"`     // How often to refresh data
	MaxProcesses        int           `json:"max_processes"`        // Maximum number of processes to track
	TemperatureWarning  float64       `json:"temperature_warning"`  // Temperature warning threshold
	TemperatureCritical float64       `json:"temperature_critical"` // Temperature critical threshold

	// Display settings
	ShowCores       bool `json:"show_cores"`        // Whether to show per-core information
	ShowProcesses   bool `json:"show_processes"`    // Whether to show process information
	ShowTemperature bool `json:"show_temperature"`  // Whether to show temperature
	ShowLoadAverage bool `json:"show_load_average"` // Whether to show load average

	// Export settings
	ExportToFile   bool          `json:"export_to_file"`  // Whether to export data to file
	ExportInterval time.Duration `json:"export_interval"` // How often to export data
	ExportFormat   string        `json:"export_format"`   // Export format (json, csv, txt)

	// Filter settings
	MinCPUUsage       float64 `json:"min_cpu_usage"`       // Minimum CPU usage to show process
	ProcessNameFilter string  `json:"process_name_filter"` // Filter processes by name
}

// CPUUsageHistory represents historical CPU usage data for graphing
type CPUUsageHistory struct {
	// Time series data
	Timestamps   []time.Time `json:"timestamps"`    // Time points
	OverallUsage []float64   `json:"overall_usage"` // Overall usage over time
	UserUsage    []float64   `json:"user_usage"`    // User usage over time
	SystemUsage  []float64   `json:"system_usage"`  // System usage over time
	IdleUsage    []float64   `json:"idle_usage"`    // Idle usage over time

	// Core usage over time
	CoreUsage [][]float64 `json:"core_usage"` // Usage for each core over time

	// Temperature history
	Temperature []float64 `json:"temperature"` // Temperature over time

	// Configuration
	MaxDataPoints  int `json:"max_data_points"`  // Maximum number of data points to keep
	DataPointCount int `json:"data_point_count"` // Current number of data points
}

// CPUMonitorAlert represents an alert condition for CPU monitoring
type CPUMonitorAlert struct {
	// Alert identification
	ID       string `json:"id"`       // Unique alert ID
	Type     string `json:"type"`     // Alert type (CPU_USAGE, TEMPERATURE, etc.)
	Severity string `json:"severity"` // Alert severity (INFO, WARNING, CRITICAL)

	// Alert conditions
	Threshold    float64       `json:"threshold"`     // Threshold value
	CurrentValue float64       `json:"current_value"` // Current value that triggered alert
	Duration     time.Duration `json:"duration"`      // How long condition has been met

	// Alert details
	Message     string `json:"message"`     // Alert message
	Description string `json:"description"` // Detailed description

	// Alert status
	IsActive       bool `json:"is_active"`       // Whether alert is currently active
	IsAcknowledged bool `json:"is_acknowledged"` // Whether alert has been acknowledged

	// Timestamps
	TriggeredAt    time.Time `json:"triggered_at"`    // When alert was triggered
	LastUpdated    time.Time `json:"last_updated"`    // When alert was last updated
	AcknowledgedAt time.Time `json:"acknowledged_at"` // When alert was acknowledged
}
