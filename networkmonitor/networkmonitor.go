package networkmonitor

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// NetworkMonitorManager is the main interface for network monitoring operations
// This struct coordinates between the collector, displayer, and exporter to provide
// a complete network monitoring solution with live updates
type NetworkMonitorManager struct {
	collector *NetworkMonitorCollector
	displayer *NetworkMonitorDisplayer
	exporter  *NetworkMonitorExporter

	// Monitoring state
	isRunning     bool
	stopChannel   chan bool
	refreshTicker *time.Ticker
	lastExportTime time.Time
}

// NewNetworkMonitorManager creates a new instance of NetworkMonitorManager
// with default collector, displayer, and exporter configurations
func NewNetworkMonitorManager() *NetworkMonitorManager {
	return &NetworkMonitorManager{
		collector:   NewNetworkMonitorCollector(),
		displayer:   NewNetworkMonitorDisplayer(),
		exporter:    NewNetworkMonitorExporter(),
		isRunning:   false,
		stopChannel: make(chan bool, 1),
	}
}

// StartLiveMonitoring starts live network monitoring with real-time updates
// This method runs continuously until stopped by the user
func (manager *NetworkMonitorManager) StartLiveMonitoring() error {
	if manager.isRunning {
		return fmt.Errorf("network monitoring is already running")
	}

	manager.isRunning = true
	manager.refreshTicker = time.NewTicker(manager.collector.config.RefreshInterval)

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("🚀 Starting live network monitoring...")
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

// StartSingleSnapshot displays a single snapshot of network information
func (manager *NetworkMonitorManager) StartSingleSnapshot() error {
	fmt.Println("📊 Collecting network information...")

	// Collect network data
	data, err := manager.collector.CollectNetworkMonitorData()
	if err != nil {
		return fmt.Errorf("failed to collect network data: %w", err)
	}

	// Display network data
	manager.displayer.DisplayNetworkMonitorData(data)

	// Always export to file for network monitor
	filePath, err := manager.exporter.ExportToJSON(data, "networkmonitor")
	if err != nil {
		fmt.Printf("⚠️  Warning: Failed to export data: %v\n", err)
	} else {
		fmt.Printf("\n💾 Network data saved to: %s\n", filePath)
	}

	return nil
}

// StopMonitoring stops the live monitoring
func (manager *NetworkMonitorManager) StopMonitoring() {
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

	fmt.Println("\n🛑 Network monitoring stopped")
}

// updateAndDisplay collects new data and updates the display
func (manager *NetworkMonitorManager) updateAndDisplay() {
	// Collect new network data
	data, err := manager.collector.CollectNetworkMonitorData()
	if err != nil {
		// Display error but continue monitoring
		fmt.Printf("\n❌ Error collecting network data: %v\n", err)
		return
	}
	
	// Display updated data
	manager.displayer.DisplayNetworkMonitorData(data)
	
	// Export data based on export interval
	manager.exportDataIfNeeded(data)
}

// exportDataIfNeeded exports network data to file based on export interval
func (manager *NetworkMonitorManager) exportDataIfNeeded(data *NetworkMonitorData) {
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

// exportData exports network data to file
func (manager *NetworkMonitorManager) exportData(data *NetworkMonitorData) {
	filePath, err := manager.exporter.ExportToJSON(data, "networkmonitor")
	if err != nil {
		// Don't display error for every export to avoid cluttering the display
		return
	}
	
	fmt.Printf("\n💾 Data exported to: %s\n", filePath)
}

// GetNetworkUsageHistory returns the current network usage history
func (manager *NetworkMonitorManager) GetNetworkUsageHistory() *NetworkUsageHistory {
	return manager.collector.GetNetworkUsageHistory()
}

// GetConfig returns the current configuration
func (manager *NetworkMonitorManager) GetConfig() *NetworkMonitorConfig {
	return manager.collector.GetConfig()
}

// UpdateConfig updates the network monitor configuration
func (manager *NetworkMonitorManager) UpdateConfig(config *NetworkMonitorConfig) {
	manager.collector.UpdateConfig(config)
}

// SetDisplayOptions configures the displayer options
func (manager *NetworkMonitorManager) SetDisplayOptions(showGraphics, showColors bool, barWidth, maxProcesses int) {
	manager.displayer.ShowGraphics = showGraphics
	manager.displayer.ShowColors = showColors
	manager.displayer.BarWidth = barWidth
	manager.displayer.MaxProcesses = maxProcesses
}

// SetExportOptions configures the exporter options
func (manager *NetworkMonitorManager) SetExportOptions(logsDir string, prettyPrint, createSubDirs bool) {
	manager.exporter.SetLogsDirectory(logsDir)
	manager.exporter.SetPrettyPrint(prettyPrint)
	manager.exporter.SetCreateSubDirs(createSubDirs)
}

// ExportToFile exports current network data to a file
func (manager *NetworkMonitorManager) ExportToFile(format string) error {
	// Collect current data
	data, err := manager.collector.CollectNetworkMonitorData()
	if err != nil {
		return fmt.Errorf("failed to collect network data: %w", err)
	}

	// Export based on format
	var filePath string
	switch format {
	case "json":
		filePath, err = manager.exporter.ExportToJSON(data, "networkmonitor")
	case "csv":
		filePath, err = manager.exporter.ExportToCSV(data, "networkmonitor")
	case "txt":
		filePath, err = manager.exporter.ExportToTXT(data, "networkmonitor")
	default:
		return fmt.Errorf("unsupported export format: %s", format)
	}

	if err != nil {
		return fmt.Errorf("failed to export data: %w", err)
	}

	fmt.Printf("💾 Network data exported to: %s\n", filePath)
	return nil
}

// GetNetworkStatus returns the current network status
func (manager *NetworkMonitorManager) GetNetworkStatus() (string, error) {
	data, err := manager.collector.CollectNetworkMonitorData()
	if err != nil {
		return "", fmt.Errorf("failed to collect network data: %w", err)
	}

	return data.NetworkStatus, nil
}

// GetNetworkAlerts returns current network alerts
func (manager *NetworkMonitorManager) GetNetworkAlerts() (bool, bool, bool, bool, error) {
	data, err := manager.collector.CollectNetworkMonitorData()
	if err != nil {
		return false, false, false, false, fmt.Errorf("failed to collect network data: %w", err)
	}

	return data.HighLatencyWarning, data.PacketLossWarning, data.BandwidthWarning, data.ConnectionWarning, nil
}

// IsRunning returns whether the network monitor is currently running
func (manager *NetworkMonitorManager) IsRunning() bool {
	return manager.isRunning
}
