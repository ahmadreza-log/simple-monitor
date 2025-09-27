package systeminfo

import (
	"fmt"
	"time"
)

// SystemInfoManager is the main interface for system information operations
// This struct coordinates between the collector and displayer to provide
// a complete system information management solution
type SystemInfoManager struct {
	collector *SystemInfoCollector
	displayer *SystemInfoDisplayer
}

// NewSystemInfoManager creates a new instance of SystemInfoManager
// with default collector and displayer configurations
func NewSystemInfoManager() *SystemInfoManager {
	return &SystemInfoManager{
		collector: NewSystemInfoCollector(),
		displayer: NewSystemInfoDisplayer(),
	}
}

// ShowSystemInfo displays comprehensive system information
// This is the main public method that collects and displays all system data
func (manager *SystemInfoManager) ShowSystemInfo() error {
	fmt.Println("🔍 Collecting system information...")

	// Collect system information
	systemInfo, err := manager.collector.CollectSystemInfo()
	if err != nil {
		return fmt.Errorf("failed to collect system information: %w", err)
	}

	// Display system information
	manager.displayer.DisplaySystemInfo(systemInfo)

	return nil
}

// ShowBasicInfo displays only basic system information
// This method shows essential system details without detailed metrics
func (manager *SystemInfoManager) ShowBasicInfo() error {
	fmt.Println("🔍 Collecting basic system information...")

	// Collect system information
	systemInfo, err := manager.collector.CollectSystemInfo()
	if err != nil {
		return fmt.Errorf("failed to collect system information: %w", err)
	}

	// Display only basic information
	manager.displayer.displayBasicInfo(systemInfo)

	return nil
}

// ShowCPUInfo displays only CPU information and usage statistics
func (manager *SystemInfoManager) ShowCPUInfo() error {
	fmt.Println("🔍 Collecting CPU information...")

	// Collect system information
	systemInfo, err := manager.collector.CollectSystemInfo()
	if err != nil {
		return fmt.Errorf("failed to collect system information: %w", err)
	}

	// Display only CPU information
	manager.displayer.displayCPUInfo(&systemInfo.CPUInfo)

	return nil
}

// ShowMemoryInfo displays only memory usage and statistics
func (manager *SystemInfoManager) ShowMemoryInfo() error {
	fmt.Println("🔍 Collecting memory information...")

	// Collect system information
	systemInfo, err := manager.collector.CollectSystemInfo()
	if err != nil {
		return fmt.Errorf("failed to collect system information: %w", err)
	}

	// Display only memory information
	manager.displayer.displayMemoryInfo(&systemInfo.MemoryInfo)

	return nil
}

// ShowDiskInfo displays only disk usage information for all drives
func (manager *SystemInfoManager) ShowDiskInfo() error {
	fmt.Println("🔍 Collecting disk information...")

	// Collect system information
	systemInfo, err := manager.collector.CollectSystemInfo()
	if err != nil {
		return fmt.Errorf("failed to collect system information: %w", err)
	}

	// Display only disk information
	manager.displayer.displayDiskInfo(systemInfo.DiskInfo)

	return nil
}

// ShowNetworkInfo displays only network interface information
func (manager *SystemInfoManager) ShowNetworkInfo() error {
	fmt.Println("🔍 Collecting network information...")

	// Collect system information
	systemInfo, err := manager.collector.CollectSystemInfo()
	if err != nil {
		return fmt.Errorf("failed to collect system information: %w", err)
	}

	// Display only network information
	manager.displayer.displayNetworkInfo(systemInfo.NetworkInfo)

	return nil
}

// StartContinuousMonitoring starts continuous monitoring of system information
// This method runs in a loop, refreshing system information at regular intervals
func (manager *SystemInfoManager) StartContinuousMonitoring(interval time.Duration) error {
	fmt.Printf("🔄 Starting continuous monitoring (refresh every %v)...\n", interval)
	fmt.Println("Press Ctrl+C to stop monitoring")

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Clear screen (platform-specific)
			fmt.Print("\033[2J\033[H")

			// Collect and display system information
			if err := manager.ShowSystemInfo(); err != nil {
				fmt.Printf("❌ Error collecting system information: %v\n", err)
			}

			// Show next refresh time
			fmt.Printf("\n⏰ Next refresh in %v\n", interval)
		}
	}
}

// GetSystemInfo returns the current system information as a struct
// This method is useful for programmatic access to system data
func (manager *SystemInfoManager) GetSystemInfo() (*SystemInfo, error) {
	return manager.collector.CollectSystemInfo()
}

// GetCPUUsage returns the current CPU usage percentage
func (manager *SystemInfoManager) GetCPUUsage() (float64, error) {
	return manager.collector.GetCPUUsage()
}

// GetMemoryUsage returns the current memory usage information
func (manager *SystemInfoManager) GetMemoryUsage() (*MemoryInfo, error) {
	return manager.collector.GetMemoryUsage()
}

// GetDiskUsage returns the current disk usage for all drives
func (manager *SystemInfoManager) GetDiskUsage() ([]DiskInfo, error) {
	return manager.collector.GetDiskUsage()
}

// GetNetworkStats returns the current network interface statistics
func (manager *SystemInfoManager) GetNetworkStats() ([]NetworkInfo, error) {
	return manager.collector.GetNetworkStats()
}

// SetDisplayOptions configures the display options for the system information
func (manager *SystemInfoManager) SetDisplayOptions(showDetailed bool, useColors bool) {
	manager.displayer.ShowDetailedInfo = showDetailed
	manager.displayer.UseColors = useColors
}

// SetCollectionOptions configures the collection options for system information
func (manager *SystemInfoManager) SetCollectionOptions(includeTemperature bool, includeNetworkStats bool) {
	manager.collector.IncludeTemperature = includeTemperature
	manager.collector.IncludeNetworkStats = includeNetworkStats
}
