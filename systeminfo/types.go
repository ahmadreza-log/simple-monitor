package systeminfo

import (
	"time"
)

// SystemInfo represents the complete system information structure
// This struct contains all the essential system details that can be monitored
type SystemInfo struct {
	// Basic system identification
	HostName        string `json:"hostname"`       // Computer hostname
	OperatingSystem string `json:"os"`             // Operating system name and version
	Architecture    string `json:"architecture"`   // System architecture (x86, x64, ARM, etc.)
	KernelVersion   string `json:"kernel_version"` // Kernel version information

	// System uptime and boot information
	Uptime   time.Duration `json:"uptime"`    // System uptime since last boot
	BootTime time.Time     `json:"boot_time"` // When the system was last booted

	// Hardware information
	CPUInfo     CPUInfo       `json:"cpu_info"`     // CPU details and statistics
	MemoryInfo  MemoryInfo    `json:"memory_info"`  // Memory usage and statistics
	DiskInfo    []DiskInfo    `json:"disk_info"`    // Disk drives information
	NetworkInfo []NetworkInfo `json:"network_info"` // Network interfaces information

	// System performance metrics
	LoadAverage  LoadAverage `json:"load_average"`  // System load average (1m, 5m, 15m)
	ProcessCount int         `json:"process_count"` // Number of running processes

	// Timestamp when this information was collected
	Timestamp time.Time `json:"timestamp"` // When this data was collected
}

// CPUInfo contains detailed CPU information and statistics
type CPUInfo struct {
	// Basic CPU identification
	ModelName string  `json:"model_name"` // CPU model name (e.g., "Intel Core i7-8700K")
	VendorID  string  `json:"vendor_id"`  // CPU vendor (e.g., "GenuineIntel")
	CPUFamily string  `json:"cpu_family"` // CPU family identifier
	CPUMHz    float64 `json:"cpu_mhz"`    // CPU frequency in MHz

	// CPU core information
	PhysicalCores int `json:"physical_cores"` // Number of physical CPU cores
	LogicalCores  int `json:"logical_cores"`  // Number of logical CPU cores (including hyperthreading)

	// CPU usage statistics
	UsagePercent  float64 `json:"usage_percent"`  // Current CPU usage percentage
	UserPercent   float64 `json:"user_percent"`   // CPU usage by user processes
	SystemPercent float64 `json:"system_percent"` // CPU usage by system processes
	IdlePercent   float64 `json:"idle_percent"`   // CPU idle percentage

	// CPU temperature (if available)
	Temperature float64 `json:"temperature"` // CPU temperature in Celsius
}

// MemoryInfo contains detailed memory usage information
type MemoryInfo struct {
	// Physical memory (RAM) information
	TotalMemory     uint64 `json:"total_memory"`     // Total physical memory in bytes
	AvailableMemory uint64 `json:"available_memory"` // Available physical memory in bytes
	UsedMemory      uint64 `json:"used_memory"`      // Used physical memory in bytes
	FreeMemory      uint64 `json:"free_memory"`      // Free physical memory in bytes

	// Memory usage percentages
	MemoryUsagePercent float64 `json:"memory_usage_percent"` // Memory usage percentage

	// Virtual memory information
	TotalSwap uint64 `json:"total_swap"` // Total swap space in bytes
	UsedSwap  uint64 `json:"used_swap"`  // Used swap space in bytes
	FreeSwap  uint64 `json:"free_swap"`  // Free swap space in bytes

	// Memory statistics
	CacheSize  uint64 `json:"cache_size"`  // Memory used for caching
	BufferSize uint64 `json:"buffer_size"` // Memory used for buffers
}

// DiskInfo contains information about a single disk drive
type DiskInfo struct {
	// Basic disk identification
	DeviceName string `json:"device_name"` // Disk device name (e.g., "/dev/sda1", "C:")
	MountPoint string `json:"mount_point"` // Mount point or drive letter
	FileSystem string `json:"file_system"` // File system type (e.g., "NTFS", "ext4")

	// Disk capacity information
	TotalSize uint64 `json:"total_size"` // Total disk size in bytes
	UsedSize  uint64 `json:"used_size"`  // Used disk space in bytes
	FreeSize  uint64 `json:"free_size"`  // Free disk space in bytes

	// Disk usage percentages
	UsagePercent float64 `json:"usage_percent"` // Disk usage percentage

	// Disk performance metrics
	ReadSpeed  uint64 `json:"read_speed"`  // Disk read speed in bytes/second
	WriteSpeed uint64 `json:"write_speed"` // Disk write speed in bytes/second

	// Disk health information
	IsRemovable bool `json:"is_removable"` // Whether the disk is removable
	IsSSD       bool `json:"is_ssd"`       // Whether the disk is an SSD
}

// NetworkInfo contains information about a network interface
type NetworkInfo struct {
	// Interface identification
	InterfaceName string `json:"interface_name"` // Network interface name (e.g., "eth0", "Wi-Fi")
	InterfaceType string `json:"interface_type"` // Interface type (e.g., "Ethernet", "WiFi", "Loopback")

	// IP address information
	IPAddress  string `json:"ip_address"`  // Primary IP address
	SubnetMask string `json:"subnet_mask"` // Subnet mask
	Gateway    string `json:"gateway"`     // Default gateway

	// Network statistics
	BytesReceived   uint64 `json:"bytes_received"`   // Total bytes received
	BytesSent       uint64 `json:"bytes_sent"`       // Total bytes sent
	PacketsReceived uint64 `json:"packets_received"` // Total packets received
	PacketsSent     uint64 `json:"packets_sent"`     // Total packets sent

	// Network status
	IsUp       bool `json:"is_up"`       // Whether the interface is up
	IsLoopback bool `json:"is_loopback"` // Whether this is a loopback interface
}

// LoadAverage represents system load average over different time periods
type LoadAverage struct {
	Load1Minute   float64 `json:"load_1_minute"`   // Load average over 1 minute
	Load5Minutes  float64 `json:"load_5_minutes"`  // Load average over 5 minutes
	Load15Minutes float64 `json:"load_15_minutes"` // Load average over 15 minutes
}
