package processmonitor

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// ProcessMonitorManager is the main interface for process monitoring operations
// This struct coordinates between the collector, displayer, and exporter to provide
// a complete process monitoring solution with live updates
type ProcessMonitorManager struct {
	collector *ProcessMonitorCollector
	displayer *ProcessMonitorDisplayer
	exporter  *ProcessMonitorExporter

	// Monitoring state
	isRunning      bool
	stopChannel    chan bool
	refreshTicker  *time.Ticker
	lastExportTime time.Time
}

// NewProcessMonitorManager creates a new instance of ProcessMonitorManager
// with default collector, displayer, and exporter configurations
func NewProcessMonitorManager() *ProcessMonitorManager {
	return &ProcessMonitorManager{
		collector:   NewProcessMonitorCollector(),
		displayer:   NewProcessMonitorDisplayer(),
		exporter:    NewProcessMonitorExporter(),
		isRunning:   false,
		stopChannel: make(chan bool, 1),
	}
}

// StartLiveMonitoring starts live process monitoring with real-time updates
// This method runs continuously until stopped by the user
func (manager *ProcessMonitorManager) StartLiveMonitoring() error {
	if manager.isRunning {
		return fmt.Errorf("process monitoring is already running")
	}

	manager.isRunning = true
	manager.refreshTicker = time.NewTicker(manager.collector.config.RefreshInterval)

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("🚀 Starting live process monitoring...")
	fmt.Println("Press Ctrl+C to stop monitoring")

	// Start monitoring loop
	go func() {
		for {
			select {
			case <-manager.refreshTicker.C:
				manager.updateAndDisplay()
			case <-manager.stopChannel:
				return
			case <-sigChan:
				manager.StopMonitoring()
				return
			}
		}
	}()

	// Wait for stop signal
	<-manager.stopChannel
	return nil
}

// StartSingleSnapshot displays a single snapshot of process information
func (manager *ProcessMonitorManager) StartSingleSnapshot() error {
	fmt.Println("📊 Collecting process information...")

	// Collect process data
	data, err := manager.collector.CollectProcessMonitorData()
	if err != nil {
		return fmt.Errorf("failed to collect process data: %w", err)
	}

	// Display process data
	manager.displayer.DisplayProcessMonitorData(data)

	// Always export to file for process monitor
	filePath, err := manager.exporter.ExportToJSON(data, "processmonitor")
	if err != nil {
		fmt.Printf("⚠️  Warning: Failed to export data: %v\n", err)
	} else {
		fmt.Printf("\n💾 Process data saved to: %s\n", filePath)
	}

	return nil
}

// StopMonitoring stops the live monitoring
func (manager *ProcessMonitorManager) StopMonitoring() {
	if !manager.isRunning {
		return
	}

	manager.isRunning = false

	if manager.refreshTicker != nil {
		manager.refreshTicker.Stop()
	}

	select {
	case manager.stopChannel <- true:
	default:
	}

	fmt.Println("\n🛑 Process monitoring stopped")
}

// updateAndDisplay collects new data and updates the display
func (manager *ProcessMonitorManager) updateAndDisplay() {
	// Collect new process data
	data, err := manager.collector.CollectProcessMonitorData()
	if err != nil {
		// Display error but continue monitoring
		fmt.Printf("\n❌ Error collecting process data: %v\n", err)
		return
	}

	// Display updated data
	manager.displayer.DisplayProcessMonitorData(data)

	// Export data based on export interval
	manager.exportDataIfNeeded(data)
}

// exportDataIfNeeded exports process data to file based on export interval
func (manager *ProcessMonitorManager) exportDataIfNeeded(data *ProcessMonitorData) {
	// Check if export is enabled
	if !manager.collector.config.ExportToFile {
		return
	}

	// Check if it's time to export based on export interval
	now := time.Now()
	if manager.lastExportTime.IsZero() || now.Sub(manager.lastExportTime) >= manager.collector.config.ExportInterval {
		manager.exportData(data)
		manager.lastExportTime = now
	}
}

// exportData exports process data to file
func (manager *ProcessMonitorManager) exportData(data *ProcessMonitorData) {
	filePath, err := manager.exporter.ExportToJSON(data, "processmonitor")
	if err != nil {
		// Don't display error for every export to avoid cluttering the display
		return
	}

	fmt.Printf("\n💾 Data exported to: %s\n", filePath)
}

// GetProcessUsageHistory returns the current process usage history
func (manager *ProcessMonitorManager) GetProcessUsageHistory() *ProcessUsageHistory {
	return manager.collector.GetProcessUsageHistory()
}

// GetConfig returns the current configuration
func (manager *ProcessMonitorManager) GetConfig() *ProcessMonitorConfig {
	return manager.collector.GetConfig()
}

// UpdateConfig updates the process monitor configuration
func (manager *ProcessMonitorManager) UpdateConfig(config *ProcessMonitorConfig) {
	manager.collector.UpdateConfig(config)
}

// SetDisplayOptions configures the displayer options
func (manager *ProcessMonitorManager) SetDisplayOptions(showGraphics, showColors bool, barWidth, maxProcesses int) {
	manager.displayer.ShowGraphics = showGraphics
	manager.displayer.ShowColors = showColors
	manager.displayer.BarWidth = barWidth
	manager.displayer.MaxProcesses = maxProcesses
}

// SetExportOptions configures the exporter options
func (manager *ProcessMonitorManager) SetExportOptions(logsDir string, prettyPrint, createSubDirs bool) {
	manager.exporter.SetLogsDirectory(logsDir)
	manager.exporter.SetPrettyPrint(prettyPrint)
	manager.exporter.SetCreateSubDirs(createSubDirs)
}

// ExportToFile exports current process data to a file
func (manager *ProcessMonitorManager) ExportToFile(format string) error {
	// Collect current data
	data, err := manager.collector.CollectProcessMonitorData()
	if err != nil {
		return fmt.Errorf("failed to collect process data: %w", err)
	}

	// Export based on format
	var filePath string
	switch format {
	case "json":
		filePath, err = manager.exporter.ExportToJSON(data, "processmonitor")
	case "csv":
		filePath, err = manager.exporter.ExportToCSV(data, "processmonitor")
	case "txt":
		filePath, err = manager.exporter.ExportToTXT(data, "processmonitor")
	default:
		return fmt.Errorf("unsupported export format: %s", format)
	}

	if err != nil {
		return fmt.Errorf("failed to export data: %w", err)
	}

	fmt.Printf("💾 Process data exported to: %s\n", filePath)
	return nil
}

// GetProcessStatus returns the current process status
func (manager *ProcessMonitorManager) GetProcessStatus() (string, error) {
	data, err := manager.collector.CollectProcessMonitorData()
	if err != nil {
		return "", fmt.Errorf("failed to collect process data: %w", err)
	}

	return data.ProcessStatus, nil
}

// GetProcessAlerts returns current process alerts
func (manager *ProcessMonitorManager) GetProcessAlerts() (bool, bool, bool, bool, bool, error) {
	data, err := manager.collector.CollectProcessMonitorData()
	if err != nil {
		return false, false, false, false, false, fmt.Errorf("failed to collect process data: %w", err)
	}

	return data.HighCPUWarning, data.HighMemoryWarning, data.HighIOWarning, data.ZombieWarning, data.ThreadWarning, nil
}

// IsRunning returns whether the process monitor is currently running
func (manager *ProcessMonitorManager) IsRunning() bool {
	return manager.isRunning
}
