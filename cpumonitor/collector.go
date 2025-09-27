package cpumonitor

import (
	"fmt"
	"runtime"
	"sort"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/process"
)

// CPUMonitorCollector handles the collection of CPU monitoring data
// This struct provides methods to gather real-time CPU metrics and process information
type CPUMonitorCollector struct {
	// Configuration
	config        *CPUMonitorConfig
	lastCPUUsage  float64
	lastTimestamp time.Time

	// Process tracking
	processCache    map[int32]*CPUProcessInfo
	lastProcessTime map[int32]time.Time

	// History tracking
	history *CPUUsageHistory
}

// NewCPUMonitorCollector creates a new instance of CPUMonitorCollector
// with default configuration values
func NewCPUMonitorCollector() *CPUMonitorCollector {
	config := &CPUMonitorConfig{
		RefreshInterval:     1 * time.Second,
		MaxProcesses:        20,
		TemperatureWarning:  70.0,
		TemperatureCritical: 85.0,
		ShowCores:           true,
		ShowProcesses:       true,
		ShowTemperature:     true,
		ShowLoadAverage:     true,
		ExportToFile:        false,
		ExportInterval:      30 * time.Second,
		ExportFormat:        "json",
		MinCPUUsage:         1.0,
		ProcessNameFilter:   "",
	}

	return &CPUMonitorCollector{
		config:          config,
		lastCPUUsage:    0.0,
		lastTimestamp:   time.Now(),
		processCache:    make(map[int32]*CPUProcessInfo),
		lastProcessTime: make(map[int32]time.Time),
		history: &CPUUsageHistory{
			MaxDataPoints:  100,
			DataPointCount: 0,
		},
	}
}

// CollectCPUMonitorData gathers comprehensive CPU monitoring data
// This is the main method that collects all available CPU metrics
func (collector *CPUMonitorCollector) CollectCPUMonitorData() (*CPUMonitorData, error) {
	data := &CPUMonitorData{
		Timestamp:       time.Now(),
		RefreshInterval: collector.config.RefreshInterval,
		IsMonitoring:    true,
	}

	// Collect basic CPU information
	if err := collector.collectBasicCPUInfo(data); err != nil {
		return nil, fmt.Errorf("failed to collect basic CPU info: %w", err)
	}

	// Collect CPU usage statistics
	if err := collector.collectCPUUsageStats(data); err != nil {
		return nil, fmt.Errorf("failed to collect CPU usage stats: %w", err)
	}

	// Collect per-core information
	if collector.config.ShowCores {
		if err := collector.collectCoreInfo(data); err != nil {
			return nil, fmt.Errorf("failed to collect core info: %w", err)
		}
	}

	// Collect process information
	if collector.config.ShowProcesses {
		if err := collector.collectProcessInfo(data); err != nil {
			return nil, fmt.Errorf("failed to collect process info: %w", err)
		}
	}

	// Collect temperature information
	if collector.config.ShowTemperature {
		if err := collector.collectTemperatureInfo(data); err != nil {
			return nil, fmt.Errorf("failed to collect temperature info: %w", err)
		}
	}

	// Collect load average
	if collector.config.ShowLoadAverage {
		if err := collector.collectLoadAverage(data); err != nil {
			return nil, fmt.Errorf("failed to collect load average: %w", err)
		}
	}

	// Update history
	collector.updateHistory(data)

	return data, nil
}

// collectBasicCPUInfo gathers basic CPU identification information
func (collector *CPUMonitorCollector) collectBasicCPUInfo(data *CPUMonitorData) error {
	// Get basic information from runtime
	data.LogicalCores = runtime.NumCPU()
	data.Architecture = runtime.GOARCH

	// Get CPU info using gopsutil
	cpuInfo, err := cpu.Info()
	if err != nil {
		return fmt.Errorf("failed to get CPU info: %w", err)
	}

	if len(cpuInfo) > 0 {
		info := cpuInfo[0]
		data.ModelName = info.ModelName
		data.VendorID = info.VendorID
		data.PhysicalCores = int(info.Cores)
	}

	return nil
}

// collectCPUUsageStats gathers overall CPU usage statistics
func (collector *CPUMonitorCollector) collectCPUUsageStats(data *CPUMonitorData) error {
	// Get CPU usage percentages
	percentages, err := cpu.Percent(time.Second, false)
	if err != nil {
		return fmt.Errorf("failed to get CPU usage: %w", err)
	}

	if len(percentages) > 0 {
		data.OverallUsage = percentages[0]
		data.IdleUsage = 100.0 - data.OverallUsage
		data.UserUsage = data.OverallUsage * 0.7   // Approximate user usage
		data.SystemUsage = data.OverallUsage * 0.3 // Approximate system usage
	}

	// Get detailed CPU times
	times, err := cpu.Times(false)
	if err == nil && len(times) > 0 {
		time := times[0]
		total := time.User + time.System + time.Idle + time.Iowait
		if total > 0 {
			data.UserUsage = (time.User / total) * 100
			data.SystemUsage = (time.System / total) * 100
			data.IdleUsage = (time.Idle / total) * 100
			data.IOWaitUsage = (time.Iowait / total) * 100
		}
	}

	return nil
}

// collectCoreInfo gathers per-core CPU information
func (collector *CPUMonitorCollector) collectCoreInfo(data *CPUMonitorData) error {
	// Get per-core CPU usage
	percentages, err := cpu.Percent(time.Second, true)
	if err != nil {
		return fmt.Errorf("failed to get per-core CPU usage: %w", err)
	}

	// Get per-core CPU times
	times, err := cpu.Times(true)
	if err != nil {
		return fmt.Errorf("failed to get per-core CPU times: %w", err)
	}

	// Create core info for each logical core
	data.Cores = make([]CPUCoreInfo, data.LogicalCores)
	for i := 0; i < data.LogicalCores; i++ {
		coreInfo := CPUCoreInfo{
			CoreID:          i,
			PhysicalID:      i / 2, // Assuming hyperthreading
			IsOnline:        true,
			IsHyperthreaded: i%2 == 1,
			LastUpdated:     time.Now(),
		}

		// Set usage percentage
		if i < len(percentages) {
			coreInfo.UsagePercent = percentages[i]
		}

		// Set detailed times if available
		if i < len(times) {
			time := times[i]
			total := time.User + time.System + time.Idle + time.Iowait
			if total > 0 {
				coreInfo.UserPercent = (time.User / total) * 100
				coreInfo.SystemPercent = (time.System / total) * 100
				coreInfo.IdlePercent = (time.Idle / total) * 100
			}
		}

		data.Cores[i] = coreInfo
	}

	return nil
}

// collectProcessInfo gathers information about CPU-consuming processes
func (collector *CPUMonitorCollector) collectProcessInfo(data *CPUMonitorData) error {
	// Get all processes
	processes, err := process.Processes()
	if err != nil {
		return fmt.Errorf("failed to get processes: %w", err)
	}

	var processInfos []CPUProcessInfo

	// Collect process information
	for _, proc := range processes {
		// Get process info
		name, _ := proc.Name()
		cpuPercent, _ := proc.CPUPercent()
		memInfo, _ := proc.MemoryInfo()
		status, _ := proc.Status()
		threads, _ := proc.NumThreads()

		// Skip processes with very low CPU usage
		if cpuPercent < collector.config.MinCPUUsage {
			continue
		}

		// Apply name filter if set
		if collector.config.ProcessNameFilter != "" &&
			!strings.Contains(strings.ToLower(name), strings.ToLower(collector.config.ProcessNameFilter)) {
			continue
		}

		processInfo := CPUProcessInfo{
			PID:             proc.Pid,
			Name:            name,
			CPUUsagePercent: cpuPercent,
			Status:          strings.Join(status, ","),
			ThreadCount:     threads,
			LastUpdated:     time.Now(),
		}

		// Add memory info if available
		if memInfo != nil {
			processInfo.MemoryUsage = memInfo.RSS
		}

		processInfos = append(processInfos, processInfo)
	}

	// Sort by CPU usage
	sort.Slice(processInfos, func(i, j int) bool {
		return processInfos[i].CPUUsagePercent > processInfos[j].CPUUsagePercent
	})

	// Limit to max processes
	maxProcesses := collector.config.MaxProcesses
	if len(processInfos) > maxProcesses {
		processInfos = processInfos[:maxProcesses]
	}

	data.TopProcesses = processInfos
	return nil
}

// collectTemperatureInfo gathers CPU temperature information
func (collector *CPUMonitorCollector) collectTemperatureInfo(data *CPUMonitorData) error {
	// TODO: Implement platform-specific temperature collection
	// This should collect:
	// - Overall CPU temperature
	// - Per-core temperatures
	// - Temperature status (Normal/Warning/Critical)
	// - Maximum safe temperature

	// For now, return placeholder data
	data.Temperature = 0.0
	data.MaxTemperature = 100.0
	data.TemperatureStatus = "Unknown"

	return nil
}

// collectLoadAverage gathers system load average information
func (collector *CPUMonitorCollector) collectLoadAverage(data *CPUMonitorData) error {
	// TODO: Implement platform-specific load average collection
	// This should collect:
	// - 1-minute load average
	// - 5-minute load average
	// - 15-minute load average

	// For now, return placeholder data
	data.LoadAverage1Min = 0.0
	data.LoadAverage5Min = 0.0
	data.LoadAverage15Min = 0.0

	return nil
}

// updateHistory updates the CPU usage history with new data
func (collector *CPUMonitorCollector) updateHistory(data *CPUMonitorData) {
	now := time.Now()

	// Add new data point
	collector.history.Timestamps = append(collector.history.Timestamps, now)
	collector.history.OverallUsage = append(collector.history.OverallUsage, data.OverallUsage)
	collector.history.UserUsage = append(collector.history.UserUsage, data.UserUsage)
	collector.history.SystemUsage = append(collector.history.SystemUsage, data.SystemUsage)
	collector.history.IdleUsage = append(collector.history.IdleUsage, data.IdleUsage)
	collector.history.Temperature = append(collector.history.Temperature, data.Temperature)

	// Add per-core data
	if len(data.Cores) > 0 {
		coreUsage := make([]float64, len(data.Cores))
		for i, core := range data.Cores {
			coreUsage[i] = core.UsagePercent
		}
		collector.history.CoreUsage = append(collector.history.CoreUsage, coreUsage)
	}

	// Maintain maximum data points
	if len(collector.history.Timestamps) > collector.history.MaxDataPoints {
		collector.history.Timestamps = collector.history.Timestamps[1:]
		collector.history.OverallUsage = collector.history.OverallUsage[1:]
		collector.history.UserUsage = collector.history.UserUsage[1:]
		collector.history.SystemUsage = collector.history.SystemUsage[1:]
		collector.history.IdleUsage = collector.history.IdleUsage[1:]
		collector.history.Temperature = collector.history.Temperature[1:]
		if len(collector.history.CoreUsage) > 0 {
			collector.history.CoreUsage = collector.history.CoreUsage[1:]
		}
	}

	collector.history.DataPointCount = len(collector.history.Timestamps)
}

// GetCPUUsageHistory returns the current CPU usage history
func (collector *CPUMonitorCollector) GetCPUUsageHistory() *CPUUsageHistory {
	return collector.history
}

// SetConfig updates the collector configuration
func (collector *CPUMonitorCollector) SetConfig(config *CPUMonitorConfig) {
	collector.config = config
}

// GetConfig returns the current collector configuration
func (collector *CPUMonitorCollector) GetConfig() *CPUMonitorConfig {
	return collector.config
}

// ResetHistory clears the CPU usage history
func (collector *CPUMonitorCollector) ResetHistory() {
	collector.history = &CPUUsageHistory{
		MaxDataPoints:  collector.history.MaxDataPoints,
		DataPointCount: 0,
	}
}

// GetProcessCPUUsage returns CPU usage for a specific process
func (collector *CPUMonitorCollector) GetProcessCPUUsage(pid int32) (float64, error) {
	// TODO: Implement process-specific CPU usage collection
	return 0.0, fmt.Errorf("process CPU usage not implemented yet")
}

// GetCoreCPUUsage returns CPU usage for a specific core
func (collector *CPUMonitorCollector) GetCoreCPUUsage(coreID int) (float64, error) {
	// TODO: Implement core-specific CPU usage collection
	return 0.0, fmt.Errorf("core CPU usage not implemented yet")
}
