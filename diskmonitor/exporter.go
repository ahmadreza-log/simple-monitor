package diskmonitor

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// DiskMonitorExporter handles exporting disk monitoring data to various formats
// This struct provides methods to save disk monitoring data to files
type DiskMonitorExporter struct {
	// Export configuration
	LogsDirectory string // Base directory for log files
	DateFormat    string // Date format for file naming
	CreateSubDirs bool   // Whether to create subdirectories for each module
	PrettyPrint   bool   // Whether to pretty print JSON output
}

// NewDiskMonitorExporter creates a new instance of DiskMonitorExporter
// with default configuration values
func NewDiskMonitorExporter() *DiskMonitorExporter {
	return &DiskMonitorExporter{
		LogsDirectory: "logs",
		DateFormat:    "2006-01-02",
		CreateSubDirs: true,
		PrettyPrint:   true,
	}
}

// ExportToJSON exports disk monitoring data to a JSON file
// The file will be saved with a timestamp-based filename in the appropriate subdirectory
func (exporter *DiskMonitorExporter) ExportToJSON(data *DiskMonitorData, moduleName string) (string, error) {
	// Create logs directory if it doesn't exist
	if err := os.MkdirAll(exporter.LogsDirectory, 0755); err != nil {
		return "", fmt.Errorf("failed to create logs directory: %w", err)
	}

	// Create subdirectory for the module if enabled
	var targetDir string
	if exporter.CreateSubDirs {
		targetDir = filepath.Join(exporter.LogsDirectory, moduleName)
		if err := os.MkdirAll(targetDir, 0755); err != nil {
			return "", fmt.Errorf("failed to create module directory: %w", err)
		}
	} else {
		targetDir = exporter.LogsDirectory
	}

	// Generate filename with timestamp
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	filename := fmt.Sprintf("%s_%s.json", moduleName, timestamp)
	filePath := filepath.Join(targetDir, filename)

	// Prepare JSON data
	var jsonData []byte
	var err error

	if exporter.PrettyPrint {
		jsonData, err = json.MarshalIndent(data, "", "  ")
	} else {
		jsonData, err = json.Marshal(data)
	}

	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON data: %w", err)
	}

	// Write to file
	if err := os.WriteFile(filePath, jsonData, 0644); err != nil {
		return "", fmt.Errorf("failed to write JSON file: %w", err)
	}

	return filePath, nil
}

// ExportToCSV exports disk monitoring data to a CSV file
// This creates a simplified CSV format with key metrics
func (exporter *DiskMonitorExporter) ExportToCSV(data *DiskMonitorData, moduleName string) (string, error) {
	// Create logs directory if it doesn't exist
	if err := os.MkdirAll(exporter.LogsDirectory, 0755); err != nil {
		return "", fmt.Errorf("failed to create logs directory: %w", err)
	}

	// Create subdirectory for the module if enabled
	var targetDir string
	if exporter.CreateSubDirs {
		targetDir = filepath.Join(exporter.LogsDirectory, moduleName)
		if err := os.MkdirAll(targetDir, 0755); err != nil {
			return "", fmt.Errorf("failed to create module directory: %w", err)
		}
	} else {
		targetDir = exporter.LogsDirectory
	}

	// Generate filename with timestamp
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	filename := fmt.Sprintf("%s_%s.csv", moduleName, timestamp)
	filePath := filepath.Join(targetDir, filename)

	// Create CSV content
	csvContent := exporter.generateCSVContent(data)

	// Write to file
	if err := os.WriteFile(filePath, []byte(csvContent), 0644); err != nil {
		return "", fmt.Errorf("failed to write CSV file: %w", err)
	}

	return filePath, nil
}

// ExportToTXT exports disk monitoring data to a text file
// This creates a human-readable text format
func (exporter *DiskMonitorExporter) ExportToTXT(data *DiskMonitorData, moduleName string) (string, error) {
	// Create logs directory if it doesn't exist
	if err := os.MkdirAll(exporter.LogsDirectory, 0755); err != nil {
		return "", fmt.Errorf("failed to create logs directory: %w", err)
	}

	// Create subdirectory for the module if enabled
	var targetDir string
	if exporter.CreateSubDirs {
		targetDir = filepath.Join(exporter.LogsDirectory, moduleName)
		if err := os.MkdirAll(targetDir, 0755); err != nil {
			return "", fmt.Errorf("failed to create module directory: %w", err)
		}
	} else {
		targetDir = exporter.LogsDirectory
	}

	// Generate filename with timestamp
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	filename := fmt.Sprintf("%s_%s.txt", moduleName, timestamp)
	filePath := filepath.Join(targetDir, filename)

	// Create text content
	txtContent := exporter.generateTXTContent(data)

	// Write to file
	if err := os.WriteFile(filePath, []byte(txtContent), 0644); err != nil {
		return "", fmt.Errorf("failed to write TXT file: %w", err)
	}

	return filePath, nil
}

// generateCSVContent creates CSV content from disk monitoring data
func (exporter *DiskMonitorExporter) generateCSVContent(data *DiskMonitorData) string {
	var content string

	// Header
	content += "Timestamp,Total Space,Used Space,Free Space,Usage Percent,Disk Status,Total Read Speed,Total Write Speed,Average IOPS,Disk Utilization\n"

	// Data row
	content += fmt.Sprintf("%s,%d,%d,%d,%.2f,%s,%.2f,%.2f,%.2f,%.2f\n",
		data.Timestamp.Format("2006-01-02 15:04:05"),
		data.TotalSpace,
		data.UsedSpace,
		data.FreeSpace,
		data.UsagePercent,
		data.DiskStatus,
		data.TotalReadSpeed,
		data.TotalWriteSpeed,
		data.AverageIOPS,
		data.DiskUtilization)

	// Partition data
	if len(data.Partitions) > 0 {
		content += "\nPartition Data\n"
		content += "Device,Mountpoint,Type,Total,Used,Free,Usage Percent\n"
		for _, partition := range data.Partitions {
			content += fmt.Sprintf("%s,%s,%s,%d,%d,%d,%.2f\n",
				partition.Device,
				partition.Mountpoint,
				partition.Fstype,
				partition.Total,
				partition.Used,
				partition.Free,
				partition.UsagePercent)
		}
	}

	// I/O data
	if len(data.DiskIO) > 0 {
		content += "\nI/O Data\n"
		content += "Device,Read Speed,Write Speed,IOPS,Utilization,Read Count,Write Count\n"
		for _, io := range data.DiskIO {
			content += fmt.Sprintf("%s,%.2f,%.2f,%.2f,%.2f,%d,%d\n",
				io.DeviceName,
				io.ReadSpeed,
				io.WriteSpeed,
				io.IOPS,
				io.Utilization,
				io.ReadCount,
				io.WriteCount)
		}
	}

	// Process data
	if len(data.TopProcesses) > 0 {
		content += "\nProcess Data\n"
		content += "PID,Name,Read Speed,Write Speed,IOPS,Total IO,Status\n"
		for _, process := range data.TopProcesses {
			content += fmt.Sprintf("%d,%s,%.2f,%.2f,%.2f,%d,%s\n",
				process.PID,
				process.Name,
				process.ReadSpeed,
				process.WriteSpeed,
				process.IOPS,
				process.TotalIO,
				process.Status)
		}
	}

	return content
}

// generateTXTContent creates text content from disk monitoring data
func (exporter *DiskMonitorExporter) generateTXTContent(data *DiskMonitorData) string {
	var content string

	// Header
	content += "DISK MONITOR REPORT\n"
	content += "==================\n\n"

	// Timestamp
	content += fmt.Sprintf("Generated: %s\n\n", data.Timestamp.Format("2006-01-02 15:04:05"))

	// Disk summary
	content += "DISK SUMMARY\n"
	content += "------------\n"
	content += fmt.Sprintf("Total Space: %s\n", exporter.formatBytes(data.TotalSpace))
	content += fmt.Sprintf("Used Space: %s\n", exporter.formatBytes(data.UsedSpace))
	content += fmt.Sprintf("Free Space: %s\n", exporter.formatBytes(data.FreeSpace))
	content += fmt.Sprintf("Usage: %.2f%%\n", data.UsagePercent)
	content += fmt.Sprintf("Status: %s\n\n", data.DiskStatus)

	// Partition information
	if len(data.Partitions) > 0 {
		content += "PARTITION INFORMATION\n"
		content += "--------------------\n"
		content += "Device\t\tMountpoint\t\tType\tTotal\t\tUsed\t\tFree\t\tUsage%\n"
		content += "------\t\t-----------\t\t----\t-----\t\t----\t\t----\t\t------\n"

		for _, partition := range data.Partitions {
			content += fmt.Sprintf("%s\t\t%s\t\t%s\t%s\t\t%s\t\t%s\t\t%.2f%%\n",
				partition.Device,
				partition.Mountpoint,
				partition.Fstype,
				exporter.formatBytes(partition.Total),
				exporter.formatBytes(partition.Used),
				exporter.formatBytes(partition.Free),
				partition.UsagePercent)
		}
		content += "\n"
	}

	// I/O statistics
	if len(data.DiskIO) > 0 {
		content += "I/O STATISTICS\n"
		content += "--------------\n"
		content += "Device\t\tRead Speed\tWrite Speed\tIOPS\t\tUtilization\n"
		content += "------\t\t----------\t-----------\t----\t\t------------\n"

		for _, io := range data.DiskIO {
			content += fmt.Sprintf("%s\t\t%.2f MB/s\t%.2f MB/s\t%.2f\t\t%.2f%%\n",
				io.DeviceName,
				io.ReadSpeed,
				io.WriteSpeed,
				io.IOPS,
				io.Utilization)
		}
		content += "\n"
	}

	// Temperature information
	if len(data.DiskTemperatures) > 0 {
		content += "TEMPERATURE INFORMATION\n"
		content += "-----------------------\n"
		content += "Device\t\tTemperature\tMax Temperature\tStatus\n"
		content += "------\t\t-----------\t---------------\t------\n"

		for _, temp := range data.DiskTemperatures {
			content += fmt.Sprintf("%s\t\t%.1f°C\t\t%.1f°C\t\t%s\n",
				temp.DeviceName,
				temp.Temperature,
				temp.MaxTemperature,
				temp.Status)
		}
		content += "\n"
	}

	// Health information
	if len(data.DiskHealth) > 0 {
		content += "HEALTH INFORMATION\n"
		content += "------------------\n"
		content += "Device\t\tHealth\t\tPower On Hours\tCycles\t\tWear%\n"
		content += "------\t\t------\t\t---------------\t------\t\t-----\n"

		for _, health := range data.DiskHealth {
			content += fmt.Sprintf("%s\t\t%s\t\t%d\t\t%d\t\t%.1f%%\n",
				health.DeviceName,
				health.HealthStatus,
				health.PowerOnHours,
				health.PowerCycleCount,
				health.WearLeveling)
		}
		content += "\n"
	}

	// Performance metrics
	content += "PERFORMANCE METRICS\n"
	content += "------------------\n"
	content += fmt.Sprintf("Total Read Speed: %.2f MB/s\n", data.TotalReadSpeed)
	content += fmt.Sprintf("Total Write Speed: %.2f MB/s\n", data.TotalWriteSpeed)
	content += fmt.Sprintf("Average IOPS: %.2f\n", data.AverageIOPS)
	content += fmt.Sprintf("Disk Utilization: %.2f%%\n\n", data.DiskUtilization)

	// Top processes
	if len(data.TopProcesses) > 0 {
		content += "TOP DISK PROCESSES\n"
		content += "------------------\n"
		content += "PID\tName\t\t\tRead Speed\tWrite Speed\tIOPS\tTotal IO\n"
		content += "---\t----\t\t\t----------\t-----------\t----\t--------\n"

		for _, process := range data.TopProcesses {
			content += fmt.Sprintf("%d\t%-20s\t%.2f MB/s\t%.2f MB/s\t%.2f\t%d\n",
				process.PID,
				process.Name,
				process.ReadSpeed,
				process.WriteSpeed,
				process.IOPS,
				process.TotalIO)
		}
		content += "\n"
	}

	// Alerts
	content += "ALERTS & WARNINGS\n"
	content += "-----------------\n"
	content += fmt.Sprintf("Low Space Warning: %t\n", data.LowSpaceWarning)
	content += fmt.Sprintf("High Temperature Warning: %t\n", data.HighTempWarning)
	content += fmt.Sprintf("Health Warning: %t\n", data.HealthWarning)
	content += fmt.Sprintf("I/O Bottleneck: %t\n", data.IOBottleneck)
	content += "\n"

	return content
}

// formatBytes formats bytes into human-readable format
func (exporter *DiskMonitorExporter) formatBytes(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// SetLogsDirectory sets the logs directory
func (exporter *DiskMonitorExporter) SetLogsDirectory(dir string) {
	exporter.LogsDirectory = dir
}

// SetPrettyPrint sets whether to pretty print JSON output
func (exporter *DiskMonitorExporter) SetPrettyPrint(pretty bool) {
	exporter.PrettyPrint = pretty
}

// SetCreateSubDirs sets whether to create subdirectories for each module
func (exporter *DiskMonitorExporter) SetCreateSubDirs(create bool) {
	exporter.CreateSubDirs = create
}
