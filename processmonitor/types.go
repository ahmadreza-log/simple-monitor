package processmonitor

import "time"

// ProcessInfo represents comprehensive information about a process
type ProcessInfo struct {
	PID             int32   `json:"pid"`              // Process ID
	Name            string  `json:"name"`             // Process name
	Status          string  `json:"status"`           // Process status
	User            string  `json:"user"`             // Process owner
	CPUUsage        float64 `json:"cpu_usage"`        // CPU usage percentage
	MemoryUsage     float64 `json:"memory_usage"`     // Memory usage percentage
	MemoryRSS       uint64  `json:"memory_rss"`       // Resident Set Size in bytes
	MemoryVMS       uint64  `json:"memory_vms"`       // Virtual Memory Size in bytes
	Threads         int32   `json:"threads"`          // Number of threads
	OpenFiles       int32   `json:"open_files"`       // Number of open files
	CreateTime      int64   `json:"create_time"`      // Process creation time
	Uptime          int64   `json:"uptime"`           // Process uptime in seconds
	ParentPID       int32   `json:"parent_pid"`       // Parent process ID
	CommandLine     string  `json:"command_line"`     // Full command line
	WorkingDir      string  `json:"working_dir"`      // Working directory
	Executable      string  `json:"executable"`       // Executable path
	Priority        int32   `json:"priority"`         // Process priority
	Nice            int32   `json:"nice"`             // Nice value
	IOReadBytes     uint64  `json:"io_read_bytes"`    // I/O read bytes
	IOWriteBytes    uint64  `json:"io_write_bytes"`   // I/O write bytes
	IOReadCount     uint64  `json:"io_read_count"`    // I/O read count
	IOWriteCount    uint64  `json:"io_write_count"`   // I/O write count
	ContextSwitches uint64  `json:"context_switches"` // Context switches
	PageFaults      uint64  `json:"page_faults"`      // Page faults
	Children        int32   `json:"children"`         // Number of child processes
}

// ProcessTreeInfo represents process tree information
type ProcessTreeInfo struct {
	PID      int32             `json:"pid"`      // Process ID
	Name     string            `json:"name"`     // Process name
	Children []ProcessTreeInfo `json:"children"` // Child processes
	Level    int               `json:"level"`    // Tree level
	IsLeaf   bool              `json:"is_leaf"`  // Whether this is a leaf process
}

// ProcessResourceInfo represents resource usage information for a process
type ProcessResourceInfo struct {
	PID             int32   `json:"pid"`              // Process ID
	Name            string  `json:"name"`             // Process name
	CPUUsage        float64 `json:"cpu_usage"`        // CPU usage percentage
	MemoryUsage     float64 `json:"memory_usage"`     // Memory usage percentage
	MemoryRSS       uint64  `json:"memory_rss"`       // Resident Set Size in bytes
	MemoryVMS       uint64  `json:"memory_vms"`       // Virtual Memory Size in bytes
	Threads         int32   `json:"threads"`          // Number of threads
	OpenFiles       int32   `json:"open_files"`       // Number of open files
	IOReadBytes     uint64  `json:"io_read_bytes"`    // I/O read bytes
	IOWriteBytes    uint64  `json:"io_write_bytes"`   // I/O write bytes
	ContextSwitches uint64  `json:"context_switches"` // Context switches
	PageFaults      uint64  `json:"page_faults"`      // Page faults
	Priority        int32   `json:"priority"`         // Process priority
	Nice            int32   `json:"nice"`             // Nice value
}

// ProcessAlertInfo represents process alert information
type ProcessAlertInfo struct {
	PID          int32     `json:"pid"`           // Process ID
	Name         string    `json:"name"`          // Process name
	AlertType    string    `json:"alert_type"`    // Type of alert
	AlertMessage string    `json:"alert_message"` // Alert message
	Severity     string    `json:"severity"`      // Alert severity (Low, Medium, High, Critical)
	Timestamp    time.Time `json:"timestamp"`     // When the alert was generated
	Value        float64   `json:"value"`         // Alert value
	Threshold    float64   `json:"threshold"`     // Alert threshold
}

// ProcessMonitorData represents comprehensive process monitoring data
type ProcessMonitorData struct {
	// All processes
	ProcessInfos []ProcessInfo `json:"process_infos"` // All process information

	// Overall process statistics
	TotalProcesses    int `json:"total_processes"`    // Total number of processes
	RunningProcesses  int `json:"running_processes"`  // Number of running processes
	SleepingProcesses int `json:"sleeping_processes"` // Number of sleeping processes
	ZombieProcesses   int `json:"zombie_processes"`   // Number of zombie processes
	StoppedProcesses  int `json:"stopped_processes"`  // Number of stopped processes

	// Top processes by different metrics
	TopCPUProcesses    []ProcessInfo `json:"top_cpu_processes"`    // Top processes by CPU usage
	TopMemoryProcesses []ProcessInfo `json:"top_memory_processes"` // Top processes by memory usage
	TopIOProcesses     []ProcessInfo `json:"top_io_processes"`     // Top processes by I/O usage
	TopThreadProcesses []ProcessInfo `json:"top_thread_processes"` // Top processes by thread count

	// Process tree information
	ProcessTree []ProcessTreeInfo `json:"process_tree"` // Process tree structure

	// Resource usage information
	ResourceUsage []ProcessResourceInfo `json:"resource_usage"` // Resource usage for all processes

	// Process alerts
	ProcessAlerts []ProcessAlertInfo `json:"process_alerts"` // Process alerts and warnings

	// Overall system metrics
	TotalCPUUsage    float64 `json:"total_cpu_usage"`    // Total CPU usage across all processes
	TotalMemoryUsage float64 `json:"total_memory_usage"` // Total memory usage across all processes
	TotalIORead      uint64  `json:"total_io_read"`      // Total I/O read across all processes
	TotalIOWrite     uint64  `json:"total_io_write"`     // Total I/O write across all processes
	TotalThreads     int32   `json:"total_threads"`      // Total number of threads
	TotalOpenFiles   int32   `json:"total_open_files"`   // Total number of open files

	// Process alerts and warnings
	ProcessStatus     string `json:"process_status"`      // Overall process status (Normal, Warning, Critical)
	HighCPUWarning    bool   `json:"high_cpu_warning"`    // High CPU usage warning
	HighMemoryWarning bool   `json:"high_memory_warning"` // High memory usage warning
	HighIOWarning     bool   `json:"high_io_warning"`     // High I/O usage warning
	ZombieWarning     bool   `json:"zombie_warning"`      // Zombie process warning
	ThreadWarning     bool   `json:"thread_warning"`      // High thread count warning

	// Monitoring configuration
	RefreshInterval time.Duration `json:"refresh_interval"` // How often data is refreshed
	IsMonitoring    bool          `json:"is_monitoring"`    // Whether monitoring is active

	// Timestamps
	Timestamp time.Time     `json:"timestamp"` // When this data was collected
	Uptime    time.Duration `json:"uptime"`    // System uptime
}

// ProcessMonitorConfig represents configuration options for process monitoring
type ProcessMonitorConfig struct {
	// Monitoring settings
	RefreshInterval     time.Duration `json:"refresh_interval"`      // How often to refresh data
	MaxProcesses        int           `json:"max_processes"`         // Maximum number of processes to track
	MaxTreeDepth        int           `json:"max_tree_depth"`        // Maximum depth for process tree
	HighCPUThreshold    float64       `json:"high_cpu_threshold"`    // High CPU usage threshold (percentage)
	HighMemoryThreshold float64       `json:"high_memory_threshold"` // High memory usage threshold (percentage)
	HighIOThreshold     uint64        `json:"high_io_threshold"`     // High I/O usage threshold (bytes)
	HighThreadThreshold int32         `json:"high_thread_threshold"` // High thread count threshold
	ZombieThreshold     int           `json:"zombie_threshold"`      // Zombie process threshold

	// Display settings
	ShowProcessTree   bool `json:"show_process_tree"`   // Whether to show process tree
	ShowResourceUsage bool `json:"show_resource_usage"` // Whether to show resource usage
	ShowTopProcesses  bool `json:"show_top_processes"`  // Whether to show top processes
	ShowAlerts        bool `json:"show_alerts"`         // Whether to show process alerts
	ShowPerformance   bool `json:"show_performance"`    // Whether to show performance metrics

	// Export settings
	ExportToFile   bool          `json:"export_to_file"`  // Whether to export data to file
	ExportInterval time.Duration `json:"export_interval"` // How often to export data
	ExportFormat   string        `json:"export_format"`   // Export format (json, csv, txt)

	// Filter settings
	MinCPUUsage       float64 `json:"min_cpu_usage"`       // Minimum CPU usage to show process
	MinMemoryUsage    float64 `json:"min_memory_usage"`    // Minimum memory usage to show process
	ProcessNameFilter string  `json:"process_name_filter"` // Filter processes by name
	UserFilter        string  `json:"user_filter"`         // Filter processes by user
	StatusFilter      string  `json:"status_filter"`       // Filter processes by status
}

// ProcessUsageHistory represents historical process usage data for graphing
type ProcessUsageHistory struct {
	// Time series data
	Timestamps       []time.Time `json:"timestamps"`         // Time points for data
	TotalCPUUsage    []float64   `json:"total_cpu_usage"`    // Total CPU usage over time
	TotalMemoryUsage []float64   `json:"total_memory_usage"` // Total memory usage over time
	TotalIORead      []float64   `json:"total_io_read"`      // Total I/O read over time
	TotalIOWrite     []float64   `json:"total_io_write"`     // Total I/O write over time
	TotalThreads     []float64   `json:"total_threads"`      // Total threads over time
	ProcessCount     []float64   `json:"process_count"`      // Process count over time

	// Configuration
	MaxDataPoints  int `json:"max_data_points"`  // Maximum number of data points to store
	DataPointCount int `json:"data_point_count"` // Current number of data points
}
