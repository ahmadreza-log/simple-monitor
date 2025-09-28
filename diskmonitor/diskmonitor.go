package diskmonitor

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// DiskMonitorManager is the main interface for disk monitoring operations
// This struct coordinates between the collector, displayer, and exporter to provide
// a complete disk monitoring solution with live updates
type DiskMonitorManager struct {
	collector *DiskMonitorCollector
	displayer *DiskMonitorDisplayer
	exporter  *DiskMonitorExporter

	// Monitoring state
	isRunning     bool
	stopChannel   chan bool
	refreshTicker *time.Ticker
	lastExportTime time.Time
}

// NewDiskMonitorManager creates a new instance of DiskMonitorManager
// with default collector, displayer, and exporter configurations
func NewDiskMonitorManager() *DiskMonitorManager {
	return &DiskMonitorManager{
		collector:   NewDiskMonitorCollector(),
		displayer:   NewDiskMonitorDisplayer(),
		exporter:    NewDiskMonitorExporter(),
		isRunning:   false,
		stopChannel: make(chan bool, 1),
	}
}

// StartLiveMonitoring starts live disk monitoring with real-time updates
// This method runs continuously until stopped by the user
func (manager *DiskMonitorManager) StartLiveMonitoring() error {
	if manager.isRunning {
		return fmt.Errorf("disk monitoring is already running")
	}

	manager.isRunning = true
	manager.refreshTicker = time.NewTicker(manager.collector.config.RefreshInterval)

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("🚀 Starting live disk monitoring...")
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

// StartSingleSnapshot displays a single snapshot of disk information
func (manager *DiskMonitorManager) StartSingleSnapshot() error {
	fmt.Println("📊 Collecting disk information...")

	// Collect disk data
	data, err := manager.collector.CollectDiskMonitorData()
	if err != nil {
		return fmt.Errorf("failed to collect disk data: %w", err)
	}

	// Display disk data
	manager.displayer.DisplayDiskMonitorData(data)

	// Always export to file for disk monitor
	filePath, err := manager.exporter.ExportToJSON(data, "diskmonitor")
	if err != nil {
		fmt.Printf("⚠️  Warning: Failed to export data: %v\n", err)
	} else {
		fmt.Printf("\n💾 Disk data saved to: %s\n", filePath)
	}

	return nil
}

// StopMonitoring stops the live monitoring
func (manager *DiskMonitorManager) StopMonitoring() {
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

	fmt.Println("\n🛑 Disk monitoring stopped")
}

// updateAndDisplay collects new data and updates the display
func (manager *DiskMonitorManager) updateAndDisplay() {
	// Collect new disk data
	data, err := manager.collector.CollectDiskMonitorData()
	if err != nil {
		// Display error but continue monitoring
		fmt.Printf("\n❌ Error collecting disk data: %v\n", err)
		return
	}
	
	// Display updated data
	manager.displayer.DisplayDiskMonitorData(data)
	
	// Export data based on export interval
	manager.exportDataIfNeeded(data)
}

// exportDataIfNeeded exports disk data to file based on export interval
func (manager *DiskMonitorManager) exportDataIfNeeded(data *DiskMonitorData) {
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

// exportData exports disk data to file
func (manager *DiskMonitorManager) exportData(data *DiskMonitorData) {
	filePath, err := manager.exporter.ExportToJSON(data, "diskmonitor")
	if err != nil {
		// Don't display error for every export to avoid cluttering the display
		return
	}
	
	fmt.Printf("\n💾 Data exported to: %s\n", filePath)
}

// GetDiskUsageHistory returns the current disk usage history
func (manager *DiskMonitorManager) GetDiskUsageHistory() *DiskUsageHistory {
	return manager.collector.GetDiskUsageHistory()
}

// GetConfig returns the current configuration
func (manager *DiskMonitorManager) GetConfig() *DiskMonitorConfig {
	return manager.collector.GetConfig()
}

// UpdateConfig updates the disk monitor configuration
func (manager *DiskMonitorManager) UpdateConfig(config *DiskMonitorConfig) {
	manager.collector.UpdateConfig(config)
}

// SetDisplayOptions configures the displayer options
func (manager *DiskMonitorManager) SetDisplayOptions(showGraphics, showColors bool, barWidth, maxProcesses int) {
	manager.displayer.ShowGraphics = showGraphics
	manager.displayer.ShowColors = showColors
	manager.displayer.BarWidth = barWidth
	manager.displayer.MaxProcesses = maxProcesses
}

// SetExportOptions configures the exporter options
func (manager *DiskMonitorManager) SetExportOptions(logsDir string, prettyPrint, createSubDirs bool) {
	manager.exporter.SetLogsDirectory(logsDir)
	manager.exporter.SetPrettyPrint(prettyPrint)
	manager.exporter.SetCreateSubDirs(createSubDirs)
}

// ExportToFile exports current disk data to a file
func (manager *DiskMonitorManager) ExportToFile(format string) error {
	// Collect current data
	data, err := manager.collector.CollectDiskMonitorData()
	if err != nil {
		return fmt.Errorf("failed to collect disk data: %w", err)
	}

	// Export based on format
	var filePath string
	switch format {
	case "json":
		filePath, err = manager.exporter.ExportToJSON(data, "diskmonitor")
	case "csv":
		filePath, err = manager.exporter.ExportToCSV(data, "diskmonitor")
	case "txt":
		filePath, err = manager.exporter.ExportToTXT(data, "diskmonitor")
	default:
		return fmt.Errorf("unsupported export format: %s", format)
	}

	if err != nil {
		return fmt.Errorf("failed to export data: %w", err)
	}

	fmt.Printf("💾 Disk data exported to: %s\n", filePath)
	return nil
}

// GetDiskStatus returns the current disk status
func (manager *DiskMonitorManager) GetDiskStatus() (string, error) {
	data, err := manager.collector.CollectDiskMonitorData()
	if err != nil {
		return "", fmt.Errorf("failed to collect disk data: %w", err)
	}

	return data.DiskStatus, nil
}

// GetDiskAlerts returns current disk alerts
func (manager *DiskMonitorManager) GetDiskAlerts() (bool, bool, bool, bool, error) {
	data, err := manager.collector.CollectDiskMonitorData()
	if err != nil {
		return false, false, false, false, fmt.Errorf("failed to collect disk data: %w", err)
	}

	return data.LowSpaceWarning, data.HighTempWarning, data.HealthWarning, data.IOBottleneck, nil
}

// IsRunning returns whether the disk monitor is currently running
func (manager *DiskMonitorManager) IsRunning() bool {
	return manager.isRunning
}
