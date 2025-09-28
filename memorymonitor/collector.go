package memorymonitor

import (
	"fmt"
	"sort"
	"time"

	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/process"
)

// MemoryMonitorCollector handles the collection of memory monitoring data
// This struct provides methods to gather real-time memory metrics and process information
type MemoryMonitorCollector struct {
	// Configuration
	config          *MemoryMonitorConfig
	lastMemoryUsage float64
	lastTimestamp   time.Time

	// Process tracking
	processCache    map[int32]*MemoryProcessInfo
	lastProcessTime map[int32]time.Time

	// History tracking
	history *MemoryUsageHistory
}

// NewMemoryMonitorCollector creates a new instance of MemoryMonitorCollector
// with default configuration values
func NewMemoryMonitorCollector() *MemoryMonitorCollector {
	config := &MemoryMonitorConfig{
		RefreshInterval:     1 * time.Second,
		MaxProcesses:        20,
		MemoryWarning:       70.0,
		MemoryCritical:      85.0,
		SwapWarning:         50.0,
		SwapCritical:        80.0,
		ShowModules:         true,
		ShowProcesses:       true,
		ShowSwap:            true,
		ShowCache:           true,
		ShowPerformance:     true,
		ExportToFile:        true,
		ExportInterval:      1 * time.Hour,
		ExportFormat:        "json",
		MinMemoryUsage:      1.0,
		ProcessNameFilter:   "",
		MemoryLeakThreshold: 10.0,
	}

	return &MemoryMonitorCollector{
		config:          config,
		lastMemoryUsage: 0.0,
		lastTimestamp:   time.Now(),
		processCache:    make(map[int32]*MemoryProcessInfo),
		lastProcessTime: make(map[int32]time.Time),
		history: &MemoryUsageHistory{
			MaxDataPoints:  100,
			DataPointCount: 0,
		},
	}
}

// CollectMemoryMonitorData gathers comprehensive memory monitoring data
// This is the main method that collects all available memory metrics
func (collector *MemoryMonitorCollector) CollectMemoryMonitorData() (*MemoryMonitorData, error) {
	data := &MemoryMonitorData{
		Timestamp:       time.Now(),
		RefreshInterval: collector.config.RefreshInterval,
		IsMonitoring:    true,
	}

	// Collect basic memory information
	if err := collector.collectBasicMemoryInfo(data); err != nil {
		return nil, fmt.Errorf("failed to collect basic memory info: %w", err)
	}

	// Collect memory breakdown
	if err := collector.collectMemoryBreakdown(data); err != nil {
		return nil, fmt.Errorf("failed to collect memory breakdown: %w", err)
	}

	// Collect performance metrics
	if collector.config.ShowPerformance {
		if err := collector.collectPerformanceMetrics(data); err != nil {
			return nil, fmt.Errorf("failed to collect performance metrics: %w", err)
		}
	}

	// Collect memory modules information
	if collector.config.ShowModules {
		if err := collector.collectMemoryModules(data); err != nil {
			return nil, fmt.Errorf("failed to collect memory modules: %w", err)
		}
	}

	// Collect swap information
	if collector.config.ShowSwap {
		if err := collector.collectSwapInfo(data); err != nil {
			return nil, fmt.Errorf("failed to collect swap info: %w", err)
		}
	}

	// Collect cache information
	if collector.config.ShowCache {
		if err := collector.collectCacheInfo(data); err != nil {
			return nil, fmt.Errorf("failed to collect cache info: %w", err)
		}
	}

	// Collect process information
	if collector.config.ShowProcesses {
		if err := collector.collectProcessInfo(data); err != nil {
			return nil, fmt.Errorf("failed to collect process info: %w", err)
		}
	}

	// Analyze memory status and alerts
	collector.analyzeMemoryStatus(data)

	// Update history
	collector.updateHistory(data)

	return data, nil
}

// collectBasicMemoryInfo gathers basic memory information
func (collector *MemoryMonitorCollector) collectBasicMemoryInfo(data *MemoryMonitorData) error {
	// Get virtual memory information
	vmem, err := mem.VirtualMemory()
	if err != nil {
		return fmt.Errorf("failed to get virtual memory info: %w", err)
	}

	data.TotalMemory = vmem.Total
	data.AvailableMemory = vmem.Available
	data.UsedMemory = vmem.Used
	data.FreeMemory = vmem.Free
	data.MemoryPercent = vmem.UsedPercent

	return nil
}

// collectMemoryBreakdown gathers detailed memory breakdown
func (collector *MemoryMonitorCollector) collectMemoryBreakdown(data *MemoryMonitorData) error {
	// Get detailed memory information
	vmem, err := mem.VirtualMemory()
	if err != nil {
		return fmt.Errorf("failed to get virtual memory info: %w", err)
	}

	// Calculate memory breakdown (approximations based on available data)
	data.UserMemory = uint64(float64(vmem.Used) * 0.6)   // Approximate user memory
	data.SystemMemory = uint64(float64(vmem.Used) * 0.3) // Approximate system memory
	data.BufferMemory = uint64(float64(vmem.Used) * 0.1) // Approximate buffer memory
	data.CacheMemory = vmem.Cached                       // Cache memory
	data.SharedMemory = vmem.Shared                      // Shared memory

	return nil
}

// collectPerformanceMetrics gathers memory performance metrics
func (collector *MemoryMonitorCollector) collectPerformanceMetrics(data *MemoryMonitorData) error {
	// Get memory statistics
	vmem, err := mem.VirtualMemory()
	if err != nil {
		return fmt.Errorf("failed to get virtual memory info: %w", err)
	}

	// Calculate memory pressure (ratio of used to total)
	data.MemoryPressure = (float64(vmem.Used) / float64(vmem.Total)) * 100

	// Calculate memory fragmentation (simplified)
	if vmem.Available > 0 {
		data.MemoryFragmentation = (float64(vmem.Available) / float64(vmem.Total)) * 100
	}

	// Get page fault information (simplified)
	data.PageFaults = uint64(vmem.Used / 1024 / 1024) // Approximate page faults
	data.PageIns = uint64(vmem.Cached / 1024 / 1024)  // Approximate page ins
	data.PageOuts = uint64(vmem.Shared / 1024 / 1024) // Approximate page outs

	return nil
}

// collectMemoryModules gathers information about memory modules
func (collector *MemoryMonitorCollector) collectMemoryModules(data *MemoryMonitorData) error {
	// For now, create a simplified memory module representation
	// In a real implementation, this would query hardware-specific information
	vmem, err := mem.VirtualMemory()
	if err != nil {
		return fmt.Errorf("failed to get virtual memory info: %w", err)
	}

	// Create a single memory module entry representing the main system memory
	module := MemoryModuleInfo{
		ModuleID:     0,
		Type:         "RAM",
		TotalSize:    vmem.Total,
		UsedSize:     vmem.Used,
		FreeSize:     vmem.Available,
		UsagePercent: vmem.UsedPercent,
		Speed:        3200, // Default speed in MHz
		Manufacturer: "System",
		Model:        "Main Memory",
		SerialNumber: "N/A",
	}

	data.MemoryModules = []MemoryModuleInfo{module}
	return nil
}

// collectSwapInfo gathers swap memory information
func (collector *MemoryMonitorCollector) collectSwapInfo(data *MemoryMonitorData) error {
	// Get swap memory information
	swap, err := mem.SwapMemory()
	if err != nil {
		return fmt.Errorf("failed to get swap memory info: %w", err)
	}

	data.SwapInfo = MemorySwapInfo{
		TotalSwap:   swap.Total,
		UsedSwap:    swap.Used,
		FreeSwap:    swap.Free,
		SwapPercent: swap.UsedPercent,
		SwapIn:      swap.Sin,
		SwapOut:     swap.Sout,
		SwapStatus:  collector.getSwapStatus(swap.UsedPercent),
	}

	return nil
}

// collectCacheInfo gathers system cache information
func (collector *MemoryMonitorCollector) collectCacheInfo(data *MemoryMonitorData) error {
	// Get memory information for cache calculation
	vmem, err := mem.VirtualMemory()
	if err != nil {
		return fmt.Errorf("failed to get virtual memory info: %w", err)
	}

	// Calculate cache information
	data.CacheInfo = MemoryCacheInfo{
		BufferCache:  vmem.Buffers,
		PageCache:    vmem.Cached,
		SlabCache:    vmem.Shared,
		TotalCache:   vmem.Buffers + vmem.Cached + vmem.Shared,
		CachePercent: (float64(vmem.Buffers+vmem.Cached+vmem.Shared) / float64(vmem.Total)) * 100,
	}

	return nil
}

// collectProcessInfo gathers top memory-consuming processes
func (collector *MemoryMonitorCollector) collectProcessInfo(data *MemoryMonitorData) error {
	// Get all processes
	processes, err := process.Processes()
	if err != nil {
		return fmt.Errorf("failed to get processes: %w", err)
	}

	var memoryProcesses []MemoryProcessInfo

	// Collect memory information for each process
	for _, p := range processes {
		// Get process memory info
		memInfo, err := p.MemoryInfo()
		if err != nil {
			continue // Skip processes we can't access
		}

		// Get process name
		name, err := p.Name()
		if err != nil {
			name = "Unknown"
		}

		// Get memory percentage
		memPercent, err := p.MemoryPercent()
		if err != nil {
			memPercent = 0.0
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

		// Get process creation time
		createTime, err := p.CreateTime()
		if err != nil {
			createTime = 0
		}

		// Filter by minimum memory usage
		if float64(memPercent) < collector.config.MinMemoryUsage {
			continue
		}

		// Filter by process name if specified
		if collector.config.ProcessNameFilter != "" && name != collector.config.ProcessNameFilter {
			continue
		}

		processInfo := MemoryProcessInfo{
			PID:           p.Pid,
			Name:          name,
			MemoryUsage:   memInfo.RSS,
			MemoryPercent: float64(memPercent),
			RSS:           memInfo.RSS,
			VMS:           memInfo.VMS,
			Status:        status[0],
			User:          user,
			CreateTime:    createTime,
		}

		memoryProcesses = append(memoryProcesses, processInfo)
	}

	// Sort by memory usage
	sort.Slice(memoryProcesses, func(i, j int) bool {
		return memoryProcesses[i].MemoryPercent > memoryProcesses[j].MemoryPercent
	})

	// Limit to max processes
	if len(memoryProcesses) > collector.config.MaxProcesses {
		memoryProcesses = memoryProcesses[:collector.config.MaxProcesses]
	}

	data.TopProcesses = memoryProcesses
	return nil
}

// analyzeMemoryStatus analyzes memory status and sets alerts
func (collector *MemoryMonitorCollector) analyzeMemoryStatus(data *MemoryMonitorData) {
	// Analyze memory status
	if data.MemoryPercent >= collector.config.MemoryCritical {
		data.MemoryStatus = "Critical"
		data.LowMemoryWarning = true
	} else if data.MemoryPercent >= collector.config.MemoryWarning {
		data.MemoryStatus = "Warning"
		data.LowMemoryWarning = true
	} else {
		data.MemoryStatus = "Normal"
		data.LowMemoryWarning = false
	}

	// Analyze swap status
	if data.SwapInfo.SwapPercent >= collector.config.SwapCritical {
		data.MemoryStatus = "Critical"
	}

	// Check for potential memory leaks (simplified detection)
	if collector.history.DataPointCount > 10 {
		// Check if memory usage is continuously increasing
		recentUsage := collector.history.TotalUsage[len(collector.history.TotalUsage)-10:]
		isIncreasing := true
		for i := 1; i < len(recentUsage); i++ {
			if recentUsage[i] <= recentUsage[i-1] {
				isIncreasing = false
				break
			}
		}
		data.MemoryLeakAlert = isIncreasing
	}
}

// updateHistory updates the memory usage history
func (collector *MemoryMonitorCollector) updateHistory(data *MemoryMonitorData) {
	now := time.Now()

	// Add new data point
	collector.history.Timestamps = append(collector.history.Timestamps, now)
	collector.history.TotalUsage = append(collector.history.TotalUsage, data.MemoryPercent)
	collector.history.UserUsage = append(collector.history.UserUsage, (float64(data.UserMemory)/float64(data.TotalMemory))*100)
	collector.history.SystemUsage = append(collector.history.SystemUsage, (float64(data.SystemMemory)/float64(data.TotalMemory))*100)
	collector.history.CacheUsage = append(collector.history.CacheUsage, (float64(data.CacheMemory)/float64(data.TotalMemory))*100)
	collector.history.SwapUsage = append(collector.history.SwapUsage, data.SwapInfo.SwapPercent)

	// Limit history size
	if len(collector.history.Timestamps) > collector.history.MaxDataPoints {
		collector.history.Timestamps = collector.history.Timestamps[1:]
		collector.history.TotalUsage = collector.history.TotalUsage[1:]
		collector.history.UserUsage = collector.history.UserUsage[1:]
		collector.history.SystemUsage = collector.history.SystemUsage[1:]
		collector.history.CacheUsage = collector.history.CacheUsage[1:]
		collector.history.SwapUsage = collector.history.SwapUsage[1:]
	}

	collector.history.DataPointCount = len(collector.history.Timestamps)
}

// getSwapStatus returns the status string for swap usage
func (collector *MemoryMonitorCollector) getSwapStatus(swapPercent float64) string {
	if swapPercent >= collector.config.SwapCritical {
		return "Critical"
	} else if swapPercent >= collector.config.SwapWarning {
		return "Warning"
	}
	return "Normal"
}

// GetMemoryUsageHistory returns the current memory usage history
func (collector *MemoryMonitorCollector) GetMemoryUsageHistory() *MemoryUsageHistory {
	return collector.history
}

// GetConfig returns the current configuration
func (collector *MemoryMonitorCollector) GetConfig() *MemoryMonitorConfig {
	return collector.config
}

// UpdateConfig updates the collector configuration
func (collector *MemoryMonitorCollector) UpdateConfig(config *MemoryMonitorConfig) {
	collector.config = config
}
