package memorymonitor

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// MemoryMonitorExporter handles exporting memory monitoring data to various formats
// This struct provides methods to save memory monitoring data to files
type MemoryMonitorExporter struct {
	// Export configuration
	LogsDirectory string // Base directory for log files
	DateFormat    string // Date format for file naming
	CreateSubDirs bool   // Whether to create subdirectories for each module
	PrettyPrint   bool   // Whether to pretty print JSON output
}

// NewMemoryMonitorExporter creates a new instance of MemoryMonitorExporter
// with default configuration values
func NewMemoryMonitorExporter() *MemoryMonitorExporter {
	return &MemoryMonitorExporter{
		LogsDirectory: "logs",
		DateFormat:    "2006-01-02",
		CreateSubDirs: true,
		PrettyPrint:   true,
	}
}

// ExportToJSON exports memory monitoring data to a JSON file
// The file will be saved with a timestamp-based filename in the appropriate subdirectory
func (exporter *MemoryMonitorExporter) ExportToJSON(data *MemoryMonitorData, moduleName string) (string, error) {
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

// ExportToCSV exports memory monitoring data to a CSV file
// This creates a simplified CSV format with key metrics
func (exporter *MemoryMonitorExporter) ExportToCSV(data *MemoryMonitorData, moduleName string) (string, error) {
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

// ExportToTXT exports memory monitoring data to a text file
// This creates a human-readable text format
func (exporter *MemoryMonitorExporter) ExportToTXT(data *MemoryMonitorData, moduleName string) (string, error) {
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

// generateCSVContent creates CSV content from memory monitoring data
func (exporter *MemoryMonitorExporter) generateCSVContent(data *MemoryMonitorData) string {
	var content string

	// Header
	content += "Timestamp,Total Memory,Used Memory,Free Memory,Memory Percent,Swap Total,Swap Used,Swap Percent,Memory Status\n"

	// Data row
	content += fmt.Sprintf("%s,%d,%d,%d,%.2f,%d,%d,%.2f,%s\n",
		data.Timestamp.Format("2006-01-02 15:04:05"),
		data.TotalMemory,
		data.UsedMemory,
		data.FreeMemory,
		data.MemoryPercent,
		data.SwapInfo.TotalSwap,
		data.SwapInfo.UsedSwap,
		data.SwapInfo.SwapPercent,
		data.MemoryStatus)

	// Process data
	if len(data.TopProcesses) > 0 {
		content += "\nProcess Data\n"
		content += "PID,Name,Memory Usage,Memory Percent,RSS,Status\n"
		for _, process := range data.TopProcesses {
			content += fmt.Sprintf("%d,%s,%d,%.2f,%d,%s\n",
				process.PID,
				process.Name,
				process.MemoryUsage,
				process.MemoryPercent,
				process.RSS,
				process.Status)
		}
	}

	return content
}

// generateTXTContent creates text content from memory monitoring data
func (exporter *MemoryMonitorExporter) generateTXTContent(data *MemoryMonitorData) string {
	var content string

	// Header
	content += "MEMORY MONITOR REPORT\n"
	content += "====================\n\n"

	// Timestamp
	content += fmt.Sprintf("Generated: %s\n\n", data.Timestamp.Format("2006-01-02 15:04:05"))

	// Memory summary
	content += "MEMORY SUMMARY\n"
	content += "--------------\n"
	content += fmt.Sprintf("Total Memory: %s\n", exporter.formatBytes(data.TotalMemory))
	content += fmt.Sprintf("Used Memory: %s\n", exporter.formatBytes(data.UsedMemory))
	content += fmt.Sprintf("Free Memory: %s\n", exporter.formatBytes(data.FreeMemory))
	content += fmt.Sprintf("Memory Usage: %.2f%%\n", data.MemoryPercent)
	content += fmt.Sprintf("Memory Status: %s\n\n", data.MemoryStatus)

	// Memory breakdown
	content += "MEMORY BREAKDOWN\n"
	content += "----------------\n"
	content += fmt.Sprintf("User Memory: %s\n", exporter.formatBytes(data.UserMemory))
	content += fmt.Sprintf("System Memory: %s\n", exporter.formatBytes(data.SystemMemory))
	content += fmt.Sprintf("Buffer Memory: %s\n", exporter.formatBytes(data.BufferMemory))
	content += fmt.Sprintf("Cache Memory: %s\n", exporter.formatBytes(data.CacheMemory))
	content += fmt.Sprintf("Shared Memory: %s\n\n", exporter.formatBytes(data.SharedMemory))

	// Swap information
	if data.SwapInfo.TotalSwap > 0 {
		content += "SWAP INFORMATION\n"
		content += "----------------\n"
		content += fmt.Sprintf("Total Swap: %s\n", exporter.formatBytes(data.SwapInfo.TotalSwap))
		content += fmt.Sprintf("Used Swap: %s\n", exporter.formatBytes(data.SwapInfo.UsedSwap))
		content += fmt.Sprintf("Free Swap: %s\n", exporter.formatBytes(data.SwapInfo.FreeSwap))
		content += fmt.Sprintf("Swap Usage: %.2f%%\n", data.SwapInfo.SwapPercent)
		content += fmt.Sprintf("Swap Status: %s\n\n", data.SwapInfo.SwapStatus)
	}

	// Cache information
	content += "CACHE INFORMATION\n"
	content += "-----------------\n"
	content += fmt.Sprintf("Buffer Cache: %s\n", exporter.formatBytes(data.CacheInfo.BufferCache))
	content += fmt.Sprintf("Page Cache: %s\n", exporter.formatBytes(data.CacheInfo.PageCache))
	content += fmt.Sprintf("Slab Cache: %s\n", exporter.formatBytes(data.CacheInfo.SlabCache))
	content += fmt.Sprintf("Total Cache: %s\n", exporter.formatBytes(data.CacheInfo.TotalCache))
	content += fmt.Sprintf("Cache Usage: %.2f%%\n\n", data.CacheInfo.CachePercent)

	// Performance metrics
	content += "PERFORMANCE METRICS\n"
	content += "-------------------\n"
	content += fmt.Sprintf("Memory Pressure: %.2f%%\n", data.MemoryPressure)
	content += fmt.Sprintf("Memory Fragmentation: %.2f%%\n", data.MemoryFragmentation)
	content += fmt.Sprintf("Page Faults: %d/sec\n", data.PageFaults)
	content += fmt.Sprintf("Page Ins: %d/sec\n", data.PageIns)
	content += fmt.Sprintf("Page Outs: %d/sec\n\n", data.PageOuts)

	// Top processes
	if len(data.TopProcesses) > 0 {
		content += "TOP MEMORY PROCESSES\n"
		content += "--------------------\n"
		content += "PID\tName\t\t\tMemory\t\tPercent\tStatus\n"
		content += "---\t----\t\t\t------\t\t-------\t------\n"

		for _, process := range data.TopProcesses {
			content += fmt.Sprintf("%d\t%-20s\t%s\t%.2f%%\t%s\n",
				process.PID,
				process.Name,
				exporter.formatBytes(process.MemoryUsage),
				process.MemoryPercent,
				process.Status)
		}
		content += "\n"
	}

	// Alerts
	content += "ALERTS & WARNINGS\n"
	content += "-----------------\n"
	content += fmt.Sprintf("Low Memory Warning: %t\n", data.LowMemoryWarning)
	content += fmt.Sprintf("Memory Leak Alert: %t\n", data.MemoryLeakAlert)
	content += "\n"

	return content
}

// formatBytes formats bytes into human-readable format
func (exporter *MemoryMonitorExporter) formatBytes(bytes uint64) string {
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
func (exporter *MemoryMonitorExporter) SetLogsDirectory(dir string) {
	exporter.LogsDirectory = dir
}

// SetPrettyPrint sets whether to pretty print JSON output
func (exporter *MemoryMonitorExporter) SetPrettyPrint(pretty bool) {
	exporter.PrettyPrint = pretty
}

// SetCreateSubDirs sets whether to create subdirectories for each module
func (exporter *MemoryMonitorExporter) SetCreateSubDirs(create bool) {
	exporter.CreateSubDirs = create
}
