package cpumonitor

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// CPUMonitorManager is the main interface for CPU monitoring operations
// This struct coordinates between the collector, displayer, and exporter to provide
// a complete CPU monitoring solution with live updates
type CPUMonitorManager struct {
	collector *CPUMonitorCollector
	displayer *CPUMonitorDisplayer
	exporter  *CPUMonitorExporter

	// Monitoring state
	isRunning     bool
	stopChannel   chan bool
	refreshTicker *time.Ticker
}

// NewCPUMonitorManager creates a new instance of CPUMonitorManager
// with default collector, displayer, and exporter configurations
func NewCPUMonitorManager() *CPUMonitorManager {
	return &CPUMonitorManager{
		collector:   NewCPUMonitorCollector(),
		displayer:   NewCPUMonitorDisplayer(),
		exporter:    NewCPUMonitorExporter(),
		isRunning:   false,
		stopChannel: make(chan bool, 1),
	}
}

// StartLiveMonitoring starts live CPU monitoring with real-time updates
// This method runs continuously until stopped by the user
func (manager *CPUMonitorManager) StartLiveMonitoring() error {
	if manager.isRunning {
		return fmt.Errorf("CPU monitoring is already running")
	}

	manager.isRunning = true
	manager.refreshTicker = time.NewTicker(manager.collector.config.RefreshInterval)

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("🚀 Starting live CPU monitoring...")
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

// StartSingleSnapshot displays a single snapshot of CPU information
func (manager *CPUMonitorManager) StartSingleSnapshot() error {
	fmt.Println("📊 Collecting CPU information...")

	// Collect CPU data
	data, err := manager.collector.CollectCPUMonitorData()
	if err != nil {
		return fmt.Errorf("failed to collect CPU data: %w", err)
	}

	// Display CPU data
	manager.displayer.DisplayCPUMonitorData(data)

	// Always export to file for CPU monitor
	filePath, err := manager.exporter.ExportToJSON(data, "cpumonitor")
	if err != nil {
		fmt.Printf("⚠️  Warning: Failed to export data: %v\n", err)
	} else {
		fmt.Printf("\n💾 CPU data saved to: %s\n", filePath)
	}

	return nil
}

// StopMonitoring stops the live monitoring
func (manager *CPUMonitorManager) StopMonitoring() {
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

	fmt.Println("\n🛑 CPU monitoring stopped")
}

// updateAndDisplay collects new data and updates the display
func (manager *CPUMonitorManager) updateAndDisplay() {
	// Collect new CPU data
	data, err := manager.collector.CollectCPUMonitorData()
	if err != nil {
		// Display error but continue monitoring
		fmt.Printf("\n❌ Error collecting CPU data: %v\n", err)
		return
	}
	
	// Display updated data
	manager.displayer.DisplayCPUMonitorData(data)
	
	// Always export data for live monitoring
	manager.exportData(data)
}

// exportData exports CPU data to file
func (manager *CPUMonitorManager) exportData(data *CPUMonitorData) {
	filePath, err := manager.exporter.ExportToJSON(data, "cpumonitor")
	if err != nil {
		// Don't display error for every export to avoid cluttering the display
		return
	}
	
	// Show export message every 10 seconds to avoid spam
	now := time.Now()
	if now.Second()%10 == 0 {
		fmt.Printf("\n💾 Data exported to: %s\n", filePath)
	}
}

// GetCPUUsageHistory returns the current CPU usage history
func (manager *CPUMonitorManager) GetCPUUsageHistory() *CPUUsageHistory {
	return manager.collector.GetCPUUsageHistory()
}

// SetRefreshInterval sets the refresh interval for live monitoring
func (manager *CPUMonitorManager) SetRefreshInterval(interval time.Duration) {
	manager.collector.config.RefreshInterval = interval
	if manager.refreshTicker != nil {
		manager.refreshTicker.Stop()
		manager.refreshTicker = time.NewTicker(interval)
	}
}

// SetMaxProcesses sets the maximum number of processes to display
func (manager *CPUMonitorManager) SetMaxProcesses(max int) {
	manager.collector.config.MaxProcesses = max
	manager.displayer.SetMaxProcesses(max)
}

// SetTemperatureThresholds sets the temperature warning and critical thresholds
func (manager *CPUMonitorManager) SetTemperatureThresholds(warning, critical float64) {
	manager.collector.config.TemperatureWarning = warning
	manager.collector.config.TemperatureCritical = critical
}

// SetDisplayOptions configures the display options
func (manager *CPUMonitorManager) SetDisplayOptions(showGraphics, showColors bool, barWidth int) {
	manager.displayer.SetGraphicsEnabled(showGraphics)
	manager.displayer.SetColorsEnabled(showColors)
	manager.displayer.SetBarWidth(barWidth)
}

// SetExportOptions configures the export options
func (manager *CPUMonitorManager) SetExportOptions(exportToFile bool, exportInterval time.Duration, exportFormat string) {
	manager.collector.config.ExportToFile = exportToFile
	manager.collector.config.ExportInterval = exportInterval
	manager.collector.config.ExportFormat = exportFormat
}

// GetCurrentData returns the current CPU monitoring data
func (manager *CPUMonitorManager) GetCurrentData() (*CPUMonitorData, error) {
	return manager.collector.CollectCPUMonitorData()
}

// IsRunning returns whether the monitoring is currently running
func (manager *CPUMonitorManager) IsRunning() bool {
	return manager.isRunning
}

// ResetHistory clears the CPU usage history
func (manager *CPUMonitorManager) ResetHistory() {
	manager.collector.ResetHistory()
}

// GetConfiguration returns the current monitoring configuration
func (manager *CPUMonitorManager) GetConfiguration() *CPUMonitorConfig {
	return manager.collector.GetConfig()
}

// SetConfiguration updates the monitoring configuration
func (manager *CPUMonitorManager) SetConfiguration(config *CPUMonitorConfig) {
	manager.collector.SetConfig(config)
}

// StartContinuousExport starts continuous export of CPU data
func (manager *CPUMonitorManager) StartContinuousExport() error {
	if !manager.collector.config.ExportToFile {
		return fmt.Errorf("export is not enabled")
	}

	exportTicker := time.NewTicker(manager.collector.config.ExportInterval)

	go func() {
		for {
			select {
			case <-exportTicker.C:
				data, err := manager.collector.CollectCPUMonitorData()
				if err == nil {
					manager.exportData(data)
				}
			case <-manager.stopChannel:
				exportTicker.Stop()
				return
			}
		}
	}()

	return nil
}

// GetProcessCPUUsage returns CPU usage for a specific process
func (manager *CPUMonitorManager) GetProcessCPUUsage(pid int32) (float64, error) {
	return manager.collector.GetProcessCPUUsage(pid)
}

// GetCoreCPUUsage returns CPU usage for a specific core
func (manager *CPUMonitorManager) GetCoreCPUUsage(coreID int) (float64, error) {
	return manager.collector.GetCoreCPUUsage(coreID)
}
