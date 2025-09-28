package processmonitor

import (
	"fmt"
	"sort"
	"time"

	"github.com/shirou/gopsutil/v3/process"
)

// ProcessMonitorCollector handles the collection of process monitoring data
// This struct provides methods to gather real-time process metrics and information
type ProcessMonitorCollector struct {
	// Configuration
	config        *ProcessMonitorConfig
	lastTimestamp time.Time

	// Process tracking
	processCache    map[int32]*ProcessInfo
	lastProcessTime map[int32]time.Time

	// History tracking
	history *ProcessUsageHistory
}

// NewProcessMonitorCollector creates a new instance of ProcessMonitorCollector
// with default configuration values
func NewProcessMonitorCollector() *ProcessMonitorCollector {
	config := &ProcessMonitorConfig{
		RefreshInterval:     1 * time.Second,
		MaxProcesses:        50,
		MaxTreeDepth:        5,
		HighCPUThreshold:    80.0,
		HighMemoryThreshold: 80.0,
		HighIOThreshold:     100 * 1024 * 1024, // 100MB
		HighThreadThreshold: 100,
		ZombieThreshold:     5,
		ShowProcessTree:     true,
		ShowResourceUsage:   true,
		ShowTopProcesses:    true,
		ShowAlerts:          true,
		ShowPerformance:     true,
		ExportToFile:        true,
		ExportInterval:      1 * time.Hour,
		ExportFormat:        "json",
		MinCPUUsage:         1.0,
		MinMemoryUsage:      1.0,
		ProcessNameFilter:   "",
		UserFilter:          "",
		StatusFilter:        "",
	}

	return &ProcessMonitorCollector{
		config:          config,
		lastTimestamp:   time.Now(),
		processCache:    make(map[int32]*ProcessInfo),
		lastProcessTime: make(map[int32]time.Time),
		history: &ProcessUsageHistory{
			MaxDataPoints:  100,
			DataPointCount: 0,
		},
	}
}

// CollectProcessMonitorData gathers comprehensive process monitoring data
// This is the main method that collects all available process metrics
func (collector *ProcessMonitorCollector) CollectProcessMonitorData() (*ProcessMonitorData, error) {
	data := &ProcessMonitorData{
		Timestamp:       time.Now(),
		RefreshInterval: collector.config.RefreshInterval,
		IsMonitoring:    true,
	}

	// Collect all processes
	if err := collector.collectAllProcesses(data); err != nil {
		return nil, fmt.Errorf("failed to collect processes: %w", err)
	}

	// Collect process tree
	if collector.config.ShowProcessTree {
		if err := collector.collectProcessTree(data); err != nil {
			return nil, fmt.Errorf("failed to collect process tree: %w", err)
		}
	}

	// Collect resource usage
	if collector.config.ShowResourceUsage {
		if err := collector.collectResourceUsage(data); err != nil {
			return nil, fmt.Errorf("failed to collect resource usage: %w", err)
		}
	}

	// Collect top processes
	if collector.config.ShowTopProcesses {
		collector.collectTopProcesses(data)
	}

	// Collect process alerts
	if collector.config.ShowAlerts {
		collector.collectProcessAlerts(data)
	}

	// Calculate performance metrics
	if collector.config.ShowPerformance {
		collector.calculatePerformanceMetrics(data)
	}

	// Analyze process status and alerts
	collector.analyzeProcessStatus(data)

	// Update history
	collector.updateHistory(data)

	return data, nil
}

// collectAllProcesses gathers information about all processes
func (collector *ProcessMonitorCollector) collectAllProcesses(data *ProcessMonitorData) error {
	// Get all processes
	processes, err := process.Processes()
	if err != nil {
		return fmt.Errorf("failed to get processes: %w", err)
	}

	var processInfos []ProcessInfo
	var totalCPU, totalMemory float64
	var totalIORead, totalIOWrite uint64
	var totalThreads, totalOpenFiles int32

	// Process each process
	for _, p := range processes {
		// Get basic process information
		processInfo, err := collector.getProcessInfo(p)
		if err != nil {
			continue // Skip processes we can't access
		}

		// Apply filters
		if !collector.passesFilters(processInfo) {
			continue
		}

		processInfos = append(processInfos, processInfo)

		// Add to totals
		totalCPU += processInfo.CPUUsage
		totalMemory += processInfo.MemoryUsage
		totalIORead += processInfo.IOReadBytes
		totalIOWrite += processInfo.IOWriteBytes
		totalThreads += processInfo.Threads
		totalOpenFiles += processInfo.OpenFiles

		// Count processes by status
		switch processInfo.Status {
		case "R":
			data.RunningProcesses++
		case "S":
			data.SleepingProcesses++
		case "Z":
			data.ZombieProcesses++
		case "T":
			data.StoppedProcesses++
		}
	}

	data.ProcessInfos = processInfos
	data.TotalProcesses = len(processInfos)
	data.TotalCPUUsage = totalCPU
	data.TotalMemoryUsage = totalMemory
	data.TotalIORead = totalIORead
	data.TotalIOWrite = totalIOWrite
	data.TotalThreads = totalThreads
	data.TotalOpenFiles = totalOpenFiles

	return nil
}

// getProcessInfo gathers detailed information about a specific process
func (collector *ProcessMonitorCollector) getProcessInfo(p *process.Process) (ProcessInfo, error) {
	var processInfo ProcessInfo

	// Basic information
	processInfo.PID = p.Pid

	// Get process name
	if name, err := p.Name(); err == nil {
		processInfo.Name = name
	}

	// Get process status
	if status, err := p.Status(); err == nil && len(status) > 0 {
		processInfo.Status = status[0]
	}

	// Get process user
	if user, err := p.Username(); err == nil {
		processInfo.User = user
	}

	// Get CPU usage
	if cpu, err := p.CPUPercent(); err == nil {
		processInfo.CPUUsage = cpu
	}

	// Get memory information
	if memInfo, err := p.MemoryInfo(); err == nil {
		processInfo.MemoryRSS = memInfo.RSS
		processInfo.MemoryVMS = memInfo.VMS
	}

	// Get memory percentage
	if memPercent, err := p.MemoryPercent(); err == nil {
		processInfo.MemoryUsage = float64(memPercent)
	}

	// Get thread count
	if threads, err := p.NumThreads(); err == nil {
		processInfo.Threads = threads
	}

	// Get open files count
	if openFiles, err := p.NumFDs(); err == nil {
		processInfo.OpenFiles = openFiles
	}

	// Get creation time
	if createTime, err := p.CreateTime(); err == nil {
		processInfo.CreateTime = createTime
		processInfo.Uptime = time.Now().Unix() - createTime/1000
	}

	// Get parent PID
	if parentPID, err := p.Parent(); err == nil {
		processInfo.ParentPID = parentPID.Pid
	}

	// Get command line
	if cmdline, err := p.Cmdline(); err == nil {
		processInfo.CommandLine = cmdline
	}

	// Get working directory
	if cwd, err := p.Cwd(); err == nil {
		processInfo.WorkingDir = cwd
	}

	// Get executable path
	if exe, err := p.Exe(); err == nil {
		processInfo.Executable = exe
	}

	// Get priority
	if priority, err := p.Nice(); err == nil {
		processInfo.Priority = priority
		processInfo.Nice = priority
	}

	// Get I/O information
	if ioInfo, err := p.IOCounters(); err == nil {
		processInfo.IOReadBytes = ioInfo.ReadBytes
		processInfo.IOWriteBytes = ioInfo.WriteBytes
		processInfo.IOReadCount = ioInfo.ReadCount
		processInfo.IOWriteCount = ioInfo.WriteCount
	}

	// Get context switches
	if ctxSwitches, err := p.NumCtxSwitches(); err == nil {
		processInfo.ContextSwitches = uint64(ctxSwitches.Voluntary + ctxSwitches.Involuntary)
	}

	// Get page faults
	if pageFaults, err := p.PageFaults(); err == nil {
		processInfo.PageFaults = pageFaults.MinorFaults + pageFaults.MajorFaults
	}

	// Get children count
	if children, err := p.Children(); err == nil {
		processInfo.Children = int32(len(children))
	}

	return processInfo, nil
}

// collectProcessTree builds the process tree structure
func (collector *ProcessMonitorCollector) collectProcessTree(data *ProcessMonitorData) error {
	// Build process tree
	processTree := collector.buildProcessTree(data.ProcessInfos, 0, collector.config.MaxTreeDepth)
	data.ProcessTree = processTree
	return nil
}

// buildProcessTree recursively builds the process tree
func (collector *ProcessMonitorCollector) buildProcessTree(processes []ProcessInfo, parentPID int32, maxDepth int) []ProcessTreeInfo {
	if maxDepth <= 0 {
		return []ProcessTreeInfo{}
	}

	var tree []ProcessTreeInfo
	for _, proc := range processes {
		if proc.ParentPID == parentPID {
			treeInfo := ProcessTreeInfo{
				PID:    proc.PID,
				Name:   proc.Name,
				Level:  maxDepth,
				IsLeaf: true,
			}

			// Find children
			children := collector.buildProcessTree(processes, proc.PID, maxDepth-1)
			if len(children) > 0 {
				treeInfo.Children = children
				treeInfo.IsLeaf = false
			}

			tree = append(tree, treeInfo)
		}
	}

	return tree
}

// collectResourceUsage gathers resource usage information
func (collector *ProcessMonitorCollector) collectResourceUsage(data *ProcessMonitorData) error {
	var resourceUsage []ProcessResourceInfo

	for _, proc := range data.ProcessInfos {
		resourceInfo := ProcessResourceInfo{
			PID:             proc.PID,
			Name:            proc.Name,
			CPUUsage:        proc.CPUUsage,
			MemoryUsage:     proc.MemoryUsage,
			MemoryRSS:       proc.MemoryRSS,
			MemoryVMS:       proc.MemoryVMS,
			Threads:         proc.Threads,
			OpenFiles:       proc.OpenFiles,
			IOReadBytes:     proc.IOReadBytes,
			IOWriteBytes:    proc.IOWriteBytes,
			ContextSwitches: proc.ContextSwitches,
			PageFaults:      proc.PageFaults,
			Priority:        proc.Priority,
			Nice:            proc.Nice,
		}

		resourceUsage = append(resourceUsage, resourceInfo)
	}

	data.ResourceUsage = resourceUsage
	return nil
}

// collectTopProcesses identifies top processes by different metrics
func (collector *ProcessMonitorCollector) collectTopProcesses(data *ProcessMonitorData) {
	// Sort by CPU usage
	cpuProcesses := make([]ProcessInfo, len(data.ProcessInfos))
	copy(cpuProcesses, data.ProcessInfos)
	sort.Slice(cpuProcesses, func(i, j int) bool {
		return cpuProcesses[i].CPUUsage > cpuProcesses[j].CPUUsage
	})
	if len(cpuProcesses) > collector.config.MaxProcesses {
		cpuProcesses = cpuProcesses[:collector.config.MaxProcesses]
	}
	data.TopCPUProcesses = cpuProcesses

	// Sort by memory usage
	memoryProcesses := make([]ProcessInfo, len(data.ProcessInfos))
	copy(memoryProcesses, data.ProcessInfos)
	sort.Slice(memoryProcesses, func(i, j int) bool {
		return memoryProcesses[i].MemoryUsage > memoryProcesses[j].MemoryUsage
	})
	if len(memoryProcesses) > collector.config.MaxProcesses {
		memoryProcesses = memoryProcesses[:collector.config.MaxProcesses]
	}
	data.TopMemoryProcesses = memoryProcesses

	// Sort by I/O usage
	ioProcesses := make([]ProcessInfo, len(data.ProcessInfos))
	copy(ioProcesses, data.ProcessInfos)
	sort.Slice(ioProcesses, func(i, j int) bool {
		return (ioProcesses[i].IOReadBytes + ioProcesses[i].IOWriteBytes) > (ioProcesses[j].IOReadBytes + ioProcesses[j].IOWriteBytes)
	})
	if len(ioProcesses) > collector.config.MaxProcesses {
		ioProcesses = ioProcesses[:collector.config.MaxProcesses]
	}
	data.TopIOProcesses = ioProcesses

	// Sort by thread count
	threadProcesses := make([]ProcessInfo, len(data.ProcessInfos))
	copy(threadProcesses, data.ProcessInfos)
	sort.Slice(threadProcesses, func(i, j int) bool {
		return threadProcesses[i].Threads > threadProcesses[j].Threads
	})
	if len(threadProcesses) > collector.config.MaxProcesses {
		threadProcesses = threadProcesses[:collector.config.MaxProcesses]
	}
	data.TopThreadProcesses = threadProcesses
}

// collectProcessAlerts identifies process alerts and warnings
func (collector *ProcessMonitorCollector) collectProcessAlerts(data *ProcessMonitorData) {
	var alerts []ProcessAlertInfo

	for _, proc := range data.ProcessInfos {
		// High CPU usage alert
		if proc.CPUUsage >= collector.config.HighCPUThreshold {
			alerts = append(alerts, ProcessAlertInfo{
				PID:          proc.PID,
				Name:         proc.Name,
				AlertType:    "High CPU Usage",
				AlertMessage: fmt.Sprintf("Process %s (PID %d) is using %.2f%% CPU", proc.Name, proc.PID, proc.CPUUsage),
				Severity:     collector.getSeverity(proc.CPUUsage, collector.config.HighCPUThreshold),
				Timestamp:    time.Now(),
				Value:        proc.CPUUsage,
				Threshold:    collector.config.HighCPUThreshold,
			})
		}

		// High memory usage alert
		if proc.MemoryUsage >= collector.config.HighMemoryThreshold {
			alerts = append(alerts, ProcessAlertInfo{
				PID:          proc.PID,
				Name:         proc.Name,
				AlertType:    "High Memory Usage",
				AlertMessage: fmt.Sprintf("Process %s (PID %d) is using %.2f%% memory", proc.Name, proc.PID, proc.MemoryUsage),
				Severity:     collector.getSeverity(proc.MemoryUsage, collector.config.HighMemoryThreshold),
				Timestamp:    time.Now(),
				Value:        proc.MemoryUsage,
				Threshold:    collector.config.HighMemoryThreshold,
			})
		}

		// High I/O usage alert
		totalIO := proc.IOReadBytes + proc.IOWriteBytes
		if totalIO >= collector.config.HighIOThreshold {
			alerts = append(alerts, ProcessAlertInfo{
				PID:          proc.PID,
				Name:         proc.Name,
				AlertType:    "High I/O Usage",
				AlertMessage: fmt.Sprintf("Process %s (PID %d) is using %d bytes I/O", proc.Name, proc.PID, totalIO),
				Severity:     collector.getSeverity(float64(totalIO), float64(collector.config.HighIOThreshold)),
				Timestamp:    time.Now(),
				Value:        float64(totalIO),
				Threshold:    float64(collector.config.HighIOThreshold),
			})
		}

		// High thread count alert
		if proc.Threads >= collector.config.HighThreadThreshold {
			alerts = append(alerts, ProcessAlertInfo{
				PID:          proc.PID,
				Name:         proc.Name,
				AlertType:    "High Thread Count",
				AlertMessage: fmt.Sprintf("Process %s (PID %d) has %d threads", proc.Name, proc.PID, proc.Threads),
				Severity:     collector.getSeverity(float64(proc.Threads), float64(collector.config.HighThreadThreshold)),
				Timestamp:    time.Now(),
				Value:        float64(proc.Threads),
				Threshold:    float64(collector.config.HighThreadThreshold),
			})
		}
	}

	data.ProcessAlerts = alerts
}

// calculatePerformanceMetrics calculates overall performance metrics
func (collector *ProcessMonitorCollector) calculatePerformanceMetrics(data *ProcessMonitorData) {
	// Calculate overall system metrics
	data.TotalCPUUsage = 0
	data.TotalMemoryUsage = 0
	data.TotalIORead = 0
	data.TotalIOWrite = 0
	data.TotalThreads = 0
	data.TotalOpenFiles = 0

	for _, proc := range data.ProcessInfos {
		data.TotalCPUUsage += proc.CPUUsage
		data.TotalMemoryUsage += proc.MemoryUsage
		data.TotalIORead += proc.IOReadBytes
		data.TotalIOWrite += proc.IOWriteBytes
		data.TotalThreads += proc.Threads
		data.TotalOpenFiles += proc.OpenFiles
	}
}

// analyzeProcessStatus analyzes process status and sets alerts
func (collector *ProcessMonitorCollector) analyzeProcessStatus(data *ProcessMonitorData) {
	// Analyze zombie processes
	if data.ZombieProcesses >= collector.config.ZombieThreshold {
		data.ZombieWarning = true
		data.ProcessStatus = "Warning"
	} else {
		data.ZombieWarning = false
	}

	// Analyze high CPU usage
	if data.TotalCPUUsage >= collector.config.HighCPUThreshold {
		data.HighCPUWarning = true
		if data.ProcessStatus == "" {
			data.ProcessStatus = "Warning"
		}
	} else {
		data.HighCPUWarning = false
	}

	// Analyze high memory usage
	if data.TotalMemoryUsage >= collector.config.HighMemoryThreshold {
		data.HighMemoryWarning = true
		if data.ProcessStatus == "" {
			data.ProcessStatus = "Warning"
		}
	} else {
		data.HighMemoryWarning = false
	}

	// Analyze high I/O usage
	totalIO := data.TotalIORead + data.TotalIOWrite
	if totalIO >= collector.config.HighIOThreshold {
		data.HighIOWarning = true
		if data.ProcessStatus == "" {
			data.ProcessStatus = "Warning"
		}
	} else {
		data.HighIOWarning = false
	}

	// Analyze high thread count
	if data.TotalThreads >= int32(collector.config.HighThreadThreshold) {
		data.ThreadWarning = true
		if data.ProcessStatus == "" {
			data.ProcessStatus = "Warning"
		}
	} else {
		data.ThreadWarning = false
	}

	// Set default status if no issues
	if data.ProcessStatus == "" {
		data.ProcessStatus = "Normal"
	}
}

// updateHistory updates the process usage history
func (collector *ProcessMonitorCollector) updateHistory(data *ProcessMonitorData) {
	now := time.Now()

	// Add new data point
	collector.history.Timestamps = append(collector.history.Timestamps, now)
	collector.history.TotalCPUUsage = append(collector.history.TotalCPUUsage, data.TotalCPUUsage)
	collector.history.TotalMemoryUsage = append(collector.history.TotalMemoryUsage, data.TotalMemoryUsage)
	collector.history.TotalIORead = append(collector.history.TotalIORead, float64(data.TotalIORead))
	collector.history.TotalIOWrite = append(collector.history.TotalIOWrite, float64(data.TotalIOWrite))
	collector.history.TotalThreads = append(collector.history.TotalThreads, float64(data.TotalThreads))
	collector.history.ProcessCount = append(collector.history.ProcessCount, float64(data.TotalProcesses))

	// Limit history size
	if len(collector.history.Timestamps) > collector.history.MaxDataPoints {
		collector.history.Timestamps = collector.history.Timestamps[1:]
		collector.history.TotalCPUUsage = collector.history.TotalCPUUsage[1:]
		collector.history.TotalMemoryUsage = collector.history.TotalMemoryUsage[1:]
		collector.history.TotalIORead = collector.history.TotalIORead[1:]
		collector.history.TotalIOWrite = collector.history.TotalIOWrite[1:]
		collector.history.TotalThreads = collector.history.TotalThreads[1:]
		collector.history.ProcessCount = collector.history.ProcessCount[1:]
	}

	collector.history.DataPointCount = len(collector.history.Timestamps)
}

// Helper methods

// passesFilters checks if a process passes all configured filters
func (collector *ProcessMonitorCollector) passesFilters(proc ProcessInfo) bool {
	// CPU usage filter
	if proc.CPUUsage < collector.config.MinCPUUsage {
		return false
	}

	// Memory usage filter
	if proc.MemoryUsage < collector.config.MinMemoryUsage {
		return false
	}

	// Process name filter
	if collector.config.ProcessNameFilter != "" && proc.Name != collector.config.ProcessNameFilter {
		return false
	}

	// User filter
	if collector.config.UserFilter != "" && proc.User != collector.config.UserFilter {
		return false
	}

	// Status filter
	if collector.config.StatusFilter != "" && proc.Status != collector.config.StatusFilter {
		return false
	}

	return true
}

// getSeverity determines the severity level based on value and threshold
func (collector *ProcessMonitorCollector) getSeverity(value, threshold float64) string {
	ratio := value / threshold
	if ratio >= 2.0 {
		return "Critical"
	} else if ratio >= 1.5 {
		return "High"
	} else if ratio >= 1.0 {
		return "Medium"
	} else {
		return "Low"
	}
}

// GetProcessUsageHistory returns the current process usage history
func (collector *ProcessMonitorCollector) GetProcessUsageHistory() *ProcessUsageHistory {
	return collector.history
}

// GetConfig returns the current configuration
func (collector *ProcessMonitorCollector) GetConfig() *ProcessMonitorConfig {
	return collector.config
}

// UpdateConfig updates the collector configuration
func (collector *ProcessMonitorCollector) UpdateConfig(config *ProcessMonitorConfig) {
	collector.config = config
}
