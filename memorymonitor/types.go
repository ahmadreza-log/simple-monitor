package memorymonitor

import "time"

// MemoryProcessInfo represents memory usage information for a specific process
type MemoryProcessInfo struct {
	PID           int32   `json:"pid"`            // Process ID
	Name          string  `json:"name"`           // Process name
	MemoryUsage   uint64  `json:"memory_usage"`   // Memory usage in bytes
	MemoryPercent float64 `json:"memory_percent"` // Memory usage percentage
	RSS           uint64  `json:"rss"`            // Resident Set Size
	VMS           uint64  `json:"vms"`            // Virtual Memory Size
	Status        string  `json:"status"`         // Process status
	User          string  `json:"user"`           // Process owner
	CreateTime    int64   `json:"create_time"`    // Process creation time
}

// MemoryModuleInfo represents memory information for a specific memory module
type MemoryModuleInfo struct {
	ModuleID     int     `json:"module_id"`     // Module identifier
	Type         string  `json:"type"`          // Memory type (RAM, Swap, etc.)
	TotalSize    uint64  `json:"total_size"`    // Total size in bytes
	UsedSize     uint64  `json:"used_size"`     // Used size in bytes
	FreeSize     uint64  `json:"free_size"`     // Free size in bytes
	UsagePercent float64 `json:"usage_percent"` // Usage percentage
	Speed        uint64  `json:"speed"`         // Memory speed in MHz
	Manufacturer string  `json:"manufacturer"`  // Memory manufacturer
	Model        string  `json:"model"`         // Memory model
	SerialNumber string  `json:"serial_number"` // Memory serial number
}

// MemorySwapInfo represents swap memory information
type MemorySwapInfo struct {
	TotalSwap   uint64  `json:"total_swap"`   // Total swap space in bytes
	UsedSwap    uint64  `json:"used_swap"`    // Used swap space in bytes
	FreeSwap    uint64  `json:"free_swap"`    // Free swap space in bytes
	SwapPercent float64 `json:"swap_percent"` // Swap usage percentage
	SwapIn      uint64  `json:"swap_in"`      // Pages swapped in
	SwapOut     uint64  `json:"swap_out"`     // Pages swapped out
	SwapStatus  string  `json:"swap_status"`  // Swap status (Normal, Warning, Critical)
}

// MemoryCacheInfo represents system cache information
type MemoryCacheInfo struct {
	BufferCache  uint64  `json:"buffer_cache"`  // Buffer cache size
	PageCache    uint64  `json:"page_cache"`    // Page cache size
	SlabCache    uint64  `json:"slab_cache"`    // Slab cache size
	TotalCache   uint64  `json:"total_cache"`   // Total cache size
	CachePercent float64 `json:"cache_percent"` // Cache percentage of total memory
}

// MemoryMonitorData represents comprehensive memory monitoring data
type MemoryMonitorData struct {
	// Basic memory information
	TotalMemory     uint64  `json:"total_memory"`     // Total system memory in bytes
	AvailableMemory uint64  `json:"available_memory"` // Available memory in bytes
	UsedMemory      uint64  `json:"used_memory"`      // Used memory in bytes
	FreeMemory      uint64  `json:"free_memory"`      // Free memory in bytes
	MemoryPercent   float64 `json:"memory_percent"`   // Memory usage percentage

	// Memory breakdown
	UserMemory   uint64 `json:"user_memory"`   // User process memory
	SystemMemory uint64 `json:"system_memory"` // System process memory
	BufferMemory uint64 `json:"buffer_memory"` // Buffer memory
	CacheMemory  uint64 `json:"cache_memory"`  // Cache memory
	SharedMemory uint64 `json:"shared_memory"` // Shared memory

	// Memory performance metrics
	MemoryPressure      float64 `json:"memory_pressure"`      // Memory pressure indicator
	MemoryFragmentation float64 `json:"memory_fragmentation"` // Memory fragmentation percentage
	PageFaults          uint64  `json:"page_faults"`          // Page faults per second
	PageIns             uint64  `json:"page_ins"`             // Pages swapped in per second
	PageOuts            uint64  `json:"page_outs"`            // Pages swapped out per second

	// Memory modules information
	MemoryModules []MemoryModuleInfo `json:"memory_modules"` // Information about memory modules

	// Swap information
	SwapInfo MemorySwapInfo `json:"swap_info"` // Swap memory information

	// Cache information
	CacheInfo MemoryCacheInfo `json:"cache_info"` // System cache information

	// Top processes by memory usage
	TopProcesses []MemoryProcessInfo `json:"top_processes"` // Top memory-consuming processes

	// Memory alerts and warnings
	MemoryStatus     string `json:"memory_status"`      // Memory status (Normal, Warning, Critical)
	LowMemoryWarning bool   `json:"low_memory_warning"` // Low memory warning flag
	MemoryLeakAlert  bool   `json:"memory_leak_alert"`  // Potential memory leak alert

	// Monitoring configuration
	RefreshInterval time.Duration `json:"refresh_interval"` // How often data is refreshed
	IsMonitoring    bool          `json:"is_monitoring"`    // Whether monitoring is active

	// Timestamps
	Timestamp time.Time     `json:"timestamp"` // When this data was collected
	Uptime    time.Duration `json:"uptime"`    // System uptime
}

// MemoryMonitorConfig represents configuration options for memory monitoring
type MemoryMonitorConfig struct {
	// Monitoring settings
	RefreshInterval time.Duration `json:"refresh_interval"` // How often to refresh data
	MaxProcesses    int           `json:"max_processes"`    // Maximum number of processes to track
	MemoryWarning   float64       `json:"memory_warning"`   // Memory warning threshold (percentage)
	MemoryCritical  float64       `json:"memory_critical"`  // Memory critical threshold (percentage)
	SwapWarning     float64       `json:"swap_warning"`     // Swap warning threshold (percentage)
	SwapCritical    float64       `json:"swap_critical"`    // Swap critical threshold (percentage)

	// Display settings
	ShowModules     bool `json:"show_modules"`     // Whether to show memory modules
	ShowProcesses   bool `json:"show_processes"`   // Whether to show process information
	ShowSwap        bool `json:"show_swap"`        // Whether to show swap information
	ShowCache       bool `json:"show_cache"`       // Whether to show cache information
	ShowPerformance bool `json:"show_performance"` // Whether to show performance metrics

	// Export settings
	ExportToFile   bool          `json:"export_to_file"`  // Whether to export data to file
	ExportInterval time.Duration `json:"export_interval"` // How often to export data
	ExportFormat   string        `json:"export_format"`   // Export format (json, csv, txt)

	// Filter settings
	MinMemoryUsage      float64 `json:"min_memory_usage"`      // Minimum memory usage to show process
	ProcessNameFilter   string  `json:"process_name_filter"`   // Filter processes by name
	MemoryLeakThreshold float64 `json:"memory_leak_threshold"` // Memory leak detection threshold
}

// MemoryUsageHistory represents historical memory usage data for graphing
type MemoryUsageHistory struct {
	// Time series data
	Timestamps  []time.Time `json:"timestamps"`   // Time points for data
	TotalUsage  []float64   `json:"total_usage"`  // Total memory usage over time
	UserUsage   []float64   `json:"user_usage"`   // User memory usage over time
	SystemUsage []float64   `json:"system_usage"` // System memory usage over time
	CacheUsage  []float64   `json:"cache_usage"`  // Cache usage over time
	SwapUsage   []float64   `json:"swap_usage"`   // Swap usage over time

	// Configuration
	MaxDataPoints  int `json:"max_data_points"`  // Maximum number of data points to store
	DataPointCount int `json:"data_point_count"` // Current number of data points
}
