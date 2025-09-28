package memorymonitor

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// MemoryMonitorManager is the main interface for memory monitoring operations
// This struct coordinates between the collector, displayer, and exporter to provide
// a complete memory monitoring solution with live updates
type MemoryMonitorManager struct {
	collector *MemoryMonitorCollector
	displayer *MemoryMonitorDisplayer
	exporter  *MemoryMonitorExporter

	// Monitoring state
	isRunning      bool
	stopChannel    chan bool
	refreshTicker  *time.Ticker
	lastExportTime time.Time
}

// NewMemoryMonitorManager creates a new instance of MemoryMonitorManager
// with default collector, displayer, and exporter configurations
func NewMemoryMonitorManager() *MemoryMonitorManager {
	return &MemoryMonitorManager{
		collector:   NewMemoryMonitorCollector(),
		displayer:   NewMemoryMonitorDisplayer(),
		exporter:    NewMemoryMonitorExporter(),
		isRunning:   false,
		stopChannel: make(chan bool, 1),
	}
}

// StartLiveMonitoring starts live memory monitoring with real-time updates
// This method runs continuously until stopped by the user
func (manager *MemoryMonitorManager) StartLiveMonitoring() error {
	if manager.isRunning {
		return fmt.Errorf("memory monitoring is already running")
	}

	manager.isRunning = true
	manager.refreshTicker = time.NewTicker(manager.collector.config.RefreshInterval)

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("🚀 Starting live memory monitoring...")
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

// StartSingleSnapshot displays a single snapshot of memory information
func (manager *MemoryMonitorManager) StartSingleSnapshot() error {
	fmt.Println("📊 Collecting memory information...")

	// Collect memory data
	data, err := manager.collector.CollectMemoryMonitorData()
	if err != nil {
		return fmt.Errorf("failed to collect memory data: %w", err)
	}

	// Display memory data
	manager.displayer.DisplayMemoryMonitorData(data)

	// Always export to file for memory monitor
	filePath, err := manager.exporter.ExportToJSON(data, "memorymonitor")
	if err != nil {
		fmt.Printf("⚠️  Warning: Failed to export data: %v\n", err)
	} else {
		fmt.Printf("\n💾 Memory data saved to: %s\n", filePath)
	}

	return nil
}

// StopMonitoring stops the live monitoring
func (manager *MemoryMonitorManager) StopMonitoring() {
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

	fmt.Println("\n🛑 Memory monitoring stopped")
}

// updateAndDisplay collects new data and updates the display
func (manager *MemoryMonitorManager) updateAndDisplay() {
	// Collect new memory data
	data, err := manager.collector.CollectMemoryMonitorData()
	if err != nil {
		// Display error but continue monitoring
		fmt.Printf("\n❌ Error collecting memory data: %v\n", err)
		return
	}

	// Display updated data
	manager.displayer.DisplayMemoryMonitorData(data)

	// Export data based on export interval
	manager.exportDataIfNeeded(data)
}

// exportDataIfNeeded exports memory data to file based on export interval
func (manager *MemoryMonitorManager) exportDataIfNeeded(data *MemoryMonitorData) {
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

// exportData exports memory data to file
func (manager *MemoryMonitorManager) exportData(data *MemoryMonitorData) {
	filePath, err := manager.exporter.ExportToJSON(data, "memorymonitor")
	if err != nil {
		// Don't display error for every export to avoid cluttering the display
		return
	}

	fmt.Printf("\n💾 Data exported to: %s\n", filePath)
}

// GetMemoryUsageHistory returns the current memory usage history
func (manager *MemoryMonitorManager) GetMemoryUsageHistory() *MemoryUsageHistory {
	return manager.collector.GetMemoryUsageHistory()
}

// GetConfig returns the current configuration
func (manager *MemoryMonitorManager) GetConfig() *MemoryMonitorConfig {
	return manager.collector.GetConfig()
}

// UpdateConfig updates the memory monitor configuration
func (manager *MemoryMonitorManager) UpdateConfig(config *MemoryMonitorConfig) {
	manager.collector.UpdateConfig(config)
}

// SetDisplayOptions configures the displayer options
func (manager *MemoryMonitorManager) SetDisplayOptions(showGraphics, showColors bool, barWidth, maxProcesses int) {
	manager.displayer.ShowGraphics = showGraphics
	manager.displayer.ShowColors = showColors
	manager.displayer.BarWidth = barWidth
	manager.displayer.MaxProcesses = maxProcesses
}

// SetExportOptions configures the exporter options
func (manager *MemoryMonitorManager) SetExportOptions(logsDir string, prettyPrint, createSubDirs bool) {
	manager.exporter.SetLogsDirectory(logsDir)
	manager.exporter.SetPrettyPrint(prettyPrint)
	manager.exporter.SetCreateSubDirs(createSubDirs)
}

// ExportToFile exports current memory data to a file
func (manager *MemoryMonitorManager) ExportToFile(format string) error {
	// Collect current data
	data, err := manager.collector.CollectMemoryMonitorData()
	if err != nil {
		return fmt.Errorf("failed to collect memory data: %w", err)
	}

	// Export based on format
	var filePath string
	switch format {
	case "json":
		filePath, err = manager.exporter.ExportToJSON(data, "memorymonitor")
	case "csv":
		filePath, err = manager.exporter.ExportToCSV(data, "memorymonitor")
	case "txt":
		filePath, err = manager.exporter.ExportToTXT(data, "memorymonitor")
	default:
		return fmt.Errorf("unsupported export format: %s", format)
	}

	if err != nil {
		return fmt.Errorf("failed to export data: %w", err)
	}

	fmt.Printf("💾 Memory data exported to: %s\n", filePath)
	return nil
}

// GetMemoryStatus returns the current memory status
func (manager *MemoryMonitorManager) GetMemoryStatus() (string, error) {
	data, err := manager.collector.CollectMemoryMonitorData()
	if err != nil {
		return "", fmt.Errorf("failed to collect memory data: %w", err)
	}

	return data.MemoryStatus, nil
}

// GetMemoryAlerts returns current memory alerts
func (manager *MemoryMonitorManager) GetMemoryAlerts() (bool, bool, error) {
	data, err := manager.collector.CollectMemoryMonitorData()
	if err != nil {
		return false, false, fmt.Errorf("failed to collect memory data: %w", err)
	}

	return data.LowMemoryWarning, data.MemoryLeakAlert, nil
}

// IsRunning returns whether the memory monitor is currently running
func (manager *MemoryMonitorManager) IsRunning() bool {
	return manager.isRunning
}
