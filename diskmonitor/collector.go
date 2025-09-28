package diskmonitor

import (
	"fmt"
	"sort"
	"time"

	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/process"
)

// DiskMonitorCollector handles the collection of disk monitoring data
// This struct provides methods to gather real-time disk metrics and process information
type DiskMonitorCollector struct {
	// Configuration
	config        *DiskMonitorConfig
	lastTimestamp time.Time

	// Process tracking
	processCache    map[int32]*DiskProcessInfo
	lastProcessTime map[int32]time.Time

	// History tracking
	history *DiskUsageHistory
}

// NewDiskMonitorCollector creates a new instance of DiskMonitorCollector
// with default configuration values
func NewDiskMonitorCollector() *DiskMonitorCollector {
	config := &DiskMonitorConfig{
		RefreshInterval:     1 * time.Second,
		MaxProcesses:        20,
		LowSpaceWarning:     80.0,
		LowSpaceCritical:    90.0,
		TempWarning:         50.0,
		TempCritical:        60.0,
		IOBottleneckThreshold: 80.0,
		ShowPartitions:      true,
		ShowIO:              true,
		ShowTemperature:     true,
		ShowHealth:          true,
		ShowProcesses:       true,
		ShowPerformance:     true,
		ExportToFile:        true,
		ExportInterval:      1 * time.Hour,
		ExportFormat:        "json",
		MinIOUsage:          1.0,
		ProcessNameFilter:   "",
		DeviceFilter:        "",
		MountpointFilter:    "",
	}

	return &DiskMonitorCollector{
		config:          config,
		lastTimestamp:   time.Now(),
		processCache:    make(map[int32]*DiskProcessInfo),
		lastProcessTime: make(map[int32]time.Time),
		history: &DiskUsageHistory{
			MaxDataPoints:  100,
			DataPointCount: 0,
		},
	}
}

// CollectDiskMonitorData gathers comprehensive disk monitoring data
// This is the main method that collects all available disk metrics
func (collector *DiskMonitorCollector) CollectDiskMonitorData() (*DiskMonitorData, error) {
	data := &DiskMonitorData{
		Timestamp:       time.Now(),
		RefreshInterval: collector.config.RefreshInterval,
		IsMonitoring:    true,
	}

	// Collect partition information
	if collector.config.ShowPartitions {
		if err := collector.collectPartitionInfo(data); err != nil {
			return nil, fmt.Errorf("failed to collect partition info: %w", err)
		}
	}

	// Collect I/O statistics
	if collector.config.ShowIO {
		if err := collector.collectIOInfo(data); err != nil {
			return nil, fmt.Errorf("failed to collect I/O info: %w", err)
		}
	}

	// Collect temperature information
	if collector.config.ShowTemperature {
		if err := collector.collectTemperatureInfo(data); err != nil {
			return nil, fmt.Errorf("failed to collect temperature info: %w", err)
		}
	}

	// Collect health information
	if collector.config.ShowHealth {
		if err := collector.collectHealthInfo(data); err != nil {
			return nil, fmt.Errorf("failed to collect health info: %w", err)
		}
	}

	// Collect process information
	if collector.config.ShowProcesses {
		if err := collector.collectProcessInfo(data); err != nil {
			return nil, fmt.Errorf("failed to collect process info: %w", err)
		}
	}

	// Calculate performance metrics
	if collector.config.ShowPerformance {
		collector.calculatePerformanceMetrics(data)
	}

	// Analyze disk status and alerts
	collector.analyzeDiskStatus(data)

	// Update history
	collector.updateHistory(data)

	return data, nil
}

// collectPartitionInfo gathers disk partition information
func (collector *DiskMonitorCollector) collectPartitionInfo(data *DiskMonitorData) error {
	// Get all partitions
	allPartitions, err := disk.Partitions(false)
	if err != nil {
		return fmt.Errorf("failed to get partitions: %w", err)
	}

	var partitionInfos []DiskPartitionInfo
	var totalSpace, usedSpace, freeSpace uint64

	// Process each partition
	for _, partition := range allPartitions {
		// Skip if device filter is specified and doesn't match
		if collector.config.DeviceFilter != "" && partition.Device != collector.config.DeviceFilter {
			continue
		}

		// Skip if mountpoint filter is specified and doesn't match
		if collector.config.MountpointFilter != "" && partition.Mountpoint != collector.config.MountpointFilter {
			continue
		}

		// Get usage for this partition
		usage, err := disk.Usage(partition.Mountpoint)
		if err != nil {
			continue // Skip partitions we can't access
		}

		partitionInfo := DiskPartitionInfo{
			Device:       partition.Device,
			Mountpoint:   partition.Mountpoint,
			Fstype:       partition.Fstype,
			Total:        usage.Total,
			Free:         usage.Free,
			Used:         usage.Used,
			UsagePercent: usage.UsedPercent,
			InodesTotal:  usage.InodesTotal,
			InodesFree:   usage.InodesFree,
			InodesUsed:   usage.InodesUsed,
		}

		partitionInfos = append(partitionInfos, partitionInfo)

		// Add to totals
		totalSpace += usage.Total
		usedSpace += usage.Used
		freeSpace += usage.Free
	}

	data.Partitions = partitionInfos
	data.TotalSpace = totalSpace
	data.UsedSpace = usedSpace
	data.FreeSpace = freeSpace

	if totalSpace > 0 {
		data.UsagePercent = (float64(usedSpace) / float64(totalSpace)) * 100
	}

	return nil
}

// collectIOInfo gathers disk I/O statistics
func (collector *DiskMonitorCollector) collectIOInfo(data *DiskMonitorData) error {
	// Get I/O counters
	ioCounters, err := disk.IOCounters()
	if err != nil {
		return fmt.Errorf("failed to get I/O counters: %w", err)
	}

	var diskIOs []DiskIOInfo
	var totalReadSpeed, totalWriteSpeed, totalIOPS float64

	for device, counter := range ioCounters {
		// Calculate speeds (simplified - would need time-based calculation for accurate speeds)
		readSpeed := float64(counter.ReadBytes) / (1024 * 1024) // Convert to MB/s (simplified)
		writeSpeed := float64(counter.WriteBytes) / (1024 * 1024) // Convert to MB/s (simplified)
		iops := float64(counter.ReadCount + counter.WriteCount) // Simplified IOPS
		utilization := float64(counter.ReadTime + counter.WriteTime) / 1000 // Simplified utilization

		diskIO := DiskIOInfo{
			DeviceName:   device,
			ReadCount:    counter.ReadCount,
			WriteCount:   counter.WriteCount,
			ReadBytes:    counter.ReadBytes,
			WriteBytes:   counter.WriteBytes,
			ReadTime:     counter.ReadTime,
			WriteTime:    counter.WriteTime,
			ReadSpeed:    readSpeed,
			WriteSpeed:   writeSpeed,
			IOPS:         iops,
			Utilization:  utilization,
		}

		diskIOs = append(diskIOs, diskIO)

		// Add to totals
		totalReadSpeed += readSpeed
		totalWriteSpeed += writeSpeed
		totalIOPS += iops
	}

	data.DiskIO = diskIOs
	data.TotalReadSpeed = totalReadSpeed
	data.TotalWriteSpeed = totalWriteSpeed
	data.AverageIOPS = totalIOPS / float64(len(diskIOs))

	return nil
}

// collectTemperatureInfo gathers disk temperature information
func (collector *DiskMonitorCollector) collectTemperatureInfo(data *DiskMonitorData) error {
	// Note: Temperature monitoring requires platform-specific implementation
	// For now, we'll create placeholder data
	var temperatures []DiskTemperatureInfo

	// Get list of disk devices
	partitions, err := disk.Partitions(false)
	if err != nil {
		return fmt.Errorf("failed to get partitions: %w", err)
	}

	for _, partition := range partitions {
		// Create placeholder temperature data
		// In a real implementation, this would query hardware sensors
		temperature := DiskTemperatureInfo{
			DeviceName:     partition.Device,
			Temperature:    35.0, // Placeholder temperature
			MaxTemperature: 60.0, // Placeholder max temperature
			Status:        "Normal",
		}

		temperatures = append(temperatures, temperature)
	}

	data.DiskTemperatures = temperatures
	return nil
}

// collectHealthInfo gathers disk health information
func (collector *DiskMonitorCollector) collectHealthInfo(data *DiskMonitorData) error {
	// Note: Health monitoring requires platform-specific implementation
	// For now, we'll create placeholder data
	var healthInfos []DiskHealthInfo

	// Get list of disk devices
	partitions, err := disk.Partitions(false)
	if err != nil {
		return fmt.Errorf("failed to get partitions: %w", err)
	}

	for _, partition := range partitions {
		// Create placeholder health data
		// In a real implementation, this would query SMART data
		health := DiskHealthInfo{
			DeviceName:         partition.Device,
			HealthStatus:       "Good",
			PowerOnHours:       8760, // Placeholder: 1 year
			PowerCycleCount:    100,  // Placeholder
			ReallocatedSectors: 0,
			PendingSectors:     0,
			UncorrectableSectors: 0,
			Temperature:        35.0,
			WearLeveling:       10.0, // Placeholder: 10% wear
		}

		healthInfos = append(healthInfos, health)
	}

	data.DiskHealth = healthInfos
	return nil
}

// collectProcessInfo gathers top disk-consuming processes
func (collector *DiskMonitorCollector) collectProcessInfo(data *DiskMonitorData) error {
	// Get all processes
	processes, err := process.Processes()
	if err != nil {
		return fmt.Errorf("failed to get processes: %w", err)
	}

	var diskProcesses []DiskProcessInfo

	// Collect disk I/O information for each process
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

		// Calculate I/O metrics
		readSpeed := float64(ioInfo.ReadBytes) / (1024 * 1024) // Convert to MB/s (simplified)
		writeSpeed := float64(ioInfo.WriteBytes) / (1024 * 1024) // Convert to MB/s (simplified)
		totalIO := ioInfo.ReadCount + ioInfo.WriteCount
		iops := float64(totalIO) // Simplified IOPS

		// Filter by minimum I/O usage
		if iops < collector.config.MinIOUsage {
			continue
		}

		// Filter by process name if specified
		if collector.config.ProcessNameFilter != "" && name != collector.config.ProcessNameFilter {
			continue
		}

		processInfo := DiskProcessInfo{
			PID:        p.Pid,
			Name:       name,
			ReadBytes:  ioInfo.ReadBytes,
			WriteBytes: ioInfo.WriteBytes,
			ReadSpeed:  readSpeed,
			WriteSpeed: writeSpeed,
			TotalIO:    totalIO,
			IOPS:       iops,
			Status:     status[0],
			User:       user,
		}

		diskProcesses = append(diskProcesses, processInfo)
	}

	// Sort by total I/O
	sort.Slice(diskProcesses, func(i, j int) bool {
		return diskProcesses[i].TotalIO > diskProcesses[j].TotalIO
	})

	// Limit to max processes
	if len(diskProcesses) > collector.config.MaxProcesses {
		diskProcesses = diskProcesses[:collector.config.MaxProcesses]
	}

	data.TopProcesses = diskProcesses
	return nil
}

// calculatePerformanceMetrics calculates overall performance metrics
func (collector *DiskMonitorCollector) calculatePerformanceMetrics(data *DiskMonitorData) {
	// Calculate overall disk utilization
	if len(data.DiskIO) > 0 {
		var totalUtilization float64
		for _, io := range data.DiskIO {
			totalUtilization += io.Utilization
		}
		data.DiskUtilization = totalUtilization / float64(len(data.DiskIO))
	}
}

// analyzeDiskStatus analyzes disk status and sets alerts
func (collector *DiskMonitorCollector) analyzeDiskStatus(data *DiskMonitorData) {
	// Analyze disk space status
	if data.UsagePercent >= collector.config.LowSpaceCritical {
		data.DiskStatus = "Critical"
		data.LowSpaceWarning = true
	} else if data.UsagePercent >= collector.config.LowSpaceWarning {
		data.DiskStatus = "Warning"
		data.LowSpaceWarning = true
	} else {
		data.DiskStatus = "Normal"
		data.LowSpaceWarning = false
	}

	// Analyze temperature status
	for _, temp := range data.DiskTemperatures {
		if temp.Temperature >= collector.config.TempCritical {
			data.HighTempWarning = true
			data.DiskStatus = "Critical"
			break
		} else if temp.Temperature >= collector.config.TempWarning {
			data.HighTempWarning = true
			if data.DiskStatus == "Normal" {
				data.DiskStatus = "Warning"
			}
		}
	}

	// Analyze I/O bottleneck
	if data.DiskUtilization >= collector.config.IOBottleneckThreshold {
		data.IOBottleneck = true
		if data.DiskStatus == "Normal" {
			data.DiskStatus = "Warning"
		}
	}

	// Analyze health status
	for _, health := range data.DiskHealth {
		if health.HealthStatus != "Good" {
			data.HealthWarning = true
			if data.DiskStatus == "Normal" {
				data.DiskStatus = "Warning"
			}
			break
		}
	}
}

// updateHistory updates the disk usage history
func (collector *DiskMonitorCollector) updateHistory(data *DiskMonitorData) {
	now := time.Now()
	
	// Add new data point
	collector.history.Timestamps = append(collector.history.Timestamps, now)
	collector.history.TotalUsage = append(collector.history.TotalUsage, data.UsagePercent)
	collector.history.ReadSpeed = append(collector.history.ReadSpeed, data.TotalReadSpeed)
	collector.history.WriteSpeed = append(collector.history.WriteSpeed, data.TotalWriteSpeed)
	collector.history.IOPS = append(collector.history.IOPS, data.AverageIOPS)
	collector.history.Utilization = append(collector.history.Utilization, data.DiskUtilization)

	// Limit history size
	if len(collector.history.Timestamps) > collector.history.MaxDataPoints {
		collector.history.Timestamps = collector.history.Timestamps[1:]
		collector.history.TotalUsage = collector.history.TotalUsage[1:]
		collector.history.ReadSpeed = collector.history.ReadSpeed[1:]
		collector.history.WriteSpeed = collector.history.WriteSpeed[1:]
		collector.history.IOPS = collector.history.IOPS[1:]
		collector.history.Utilization = collector.history.Utilization[1:]
	}

	collector.history.DataPointCount = len(collector.history.Timestamps)
}

// GetDiskUsageHistory returns the current disk usage history
func (collector *DiskMonitorCollector) GetDiskUsageHistory() *DiskUsageHistory {
	return collector.history
}

// GetConfig returns the current configuration
func (collector *DiskMonitorCollector) GetConfig() *DiskMonitorConfig {
	return collector.config
}

// UpdateConfig updates the collector configuration
func (collector *DiskMonitorCollector) UpdateConfig(config *DiskMonitorConfig) {
	collector.config = config
}
