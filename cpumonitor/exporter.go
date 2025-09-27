package cpumonitor

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// CPUMonitorExporter handles exporting CPU monitoring data to various formats
// This struct provides methods to save CPU monitoring data to files
type CPUMonitorExporter struct {
	// Export configuration
	LogsDirectory string // Base directory for log files
	DateFormat    string // Date format for file naming
	CreateSubDirs bool   // Whether to create subdirectories for each module
	PrettyPrint   bool   // Whether to pretty print JSON output
}

// NewCPUMonitorExporter creates a new instance of CPUMonitorExporter
// with default configuration values
func NewCPUMonitorExporter() *CPUMonitorExporter {
	return &CPUMonitorExporter{
		LogsDirectory: "logs",
		DateFormat:    "2006-01-02",
		CreateSubDirs: true,
		PrettyPrint:   true,
	}
}

// ExportToJSON exports CPU monitoring data to a JSON file
// The file will be saved with a timestamp-based filename in the appropriate subdirectory
func (exporter *CPUMonitorExporter) ExportToJSON(data *CPUMonitorData, moduleName string) (string, error) {
	// Generate filename with current date and time
	fileName := exporter.generateFileName(moduleName, "json")

	// Create full file path
	filePath, err := exporter.createFilePath(moduleName, fileName)
	if err != nil {
		return "", fmt.Errorf("failed to create file path: %w", err)
	}

	// Ensure directory exists
	if err := exporter.ensureDirectoryExists(filepath.Dir(filePath)); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	// Prepare data for export
	exportData := exporter.prepareExportData(data)

	// Convert to JSON
	var jsonData []byte
	if exporter.PrettyPrint {
		jsonData, err = json.MarshalIndent(exportData, "", "  ")
	} else {
		jsonData, err = json.Marshal(exportData)
	}

	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON: %w", err)
	}

	// Write to file
	if err := os.WriteFile(filePath, jsonData, 0644); err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	return filePath, nil
}

// ExportToCSV exports CPU monitoring data to a CSV file
// This method will be implemented when CSV export is needed
func (exporter *CPUMonitorExporter) ExportToCSV(data *CPUMonitorData, moduleName string) (string, error) {
	// TODO: Implement CSV export functionality
	return "", fmt.Errorf("CSV export not implemented yet")
}

// ExportToText exports CPU monitoring data to a text file
// This method will be implemented when text export is needed
func (exporter *CPUMonitorExporter) ExportToText(data *CPUMonitorData, moduleName string) (string, error) {
	// TODO: Implement text export functionality
	return "", fmt.Errorf("text export not implemented yet")
}

// ExportHistoryToJSON exports CPU usage history to a JSON file
func (exporter *CPUMonitorExporter) ExportHistoryToJSON(history *CPUUsageHistory, moduleName string) (string, error) {
	// Generate filename with current date and time
	fileName := exporter.generateHistoryFileName(moduleName, "json")

	// Create full file path
	filePath, err := exporter.createFilePath(moduleName, fileName)
	if err != nil {
		return "", fmt.Errorf("failed to create file path: %w", err)
	}

	// Ensure directory exists
	if err := exporter.ensureDirectoryExists(filepath.Dir(filePath)); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	// Prepare data for export
	exportData := map[string]interface{}{
		"export_info": map[string]interface{}{
			"exported_at":      time.Now().Format(time.RFC3339),
			"export_format":    "json",
			"exporter_version": "1.0.0",
			"data_type":        "cpu_usage_history",
		},
		"cpu_history": history,
	}

	// Convert to JSON
	var jsonData []byte
	if exporter.PrettyPrint {
		jsonData, err = json.MarshalIndent(exportData, "", "  ")
	} else {
		jsonData, err = json.Marshal(exportData)
	}

	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON: %w", err)
	}

	// Write to file
	if err := os.WriteFile(filePath, jsonData, 0644); err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	return filePath, nil
}

// generateFileName creates a filename with the specified module name and extension
// Format: {moduleName}-{date}-{time}.{extension}
func (exporter *CPUMonitorExporter) generateFileName(moduleName, extension string) string {
	now := time.Now()
	dateStr := now.Format(exporter.DateFormat)
	timeStr := now.Format("15-04-05")
	return fmt.Sprintf("%s-%s-%s.%s", moduleName, dateStr, timeStr, extension)
}

// generateHistoryFileName creates a filename for history exports
// Format: {moduleName}-history-{date}-{time}.{extension}
func (exporter *CPUMonitorExporter) generateHistoryFileName(moduleName, extension string) string {
	now := time.Now()
	dateStr := now.Format(exporter.DateFormat)
	timeStr := now.Format("15-04-05")
	return fmt.Sprintf("%s-history-%s-%s.%s", moduleName, dateStr, timeStr, extension)
}

// createFilePath creates the full file path for the export
// If CreateSubDirs is true, it will create a subdirectory for the module
func (exporter *CPUMonitorExporter) createFilePath(moduleName, fileName string) (string, error) {
	if exporter.CreateSubDirs {
		return filepath.Join(exporter.LogsDirectory, moduleName, fileName), nil
	}
	return filepath.Join(exporter.LogsDirectory, fileName), nil
}

// ensureDirectoryExists creates the directory if it doesn't exist
func (exporter *CPUMonitorExporter) ensureDirectoryExists(dirPath string) error {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		return os.MkdirAll(dirPath, 0755)
	}
	return nil
}

// prepareExportData prepares the CPU monitoring data for export
// This method can be extended to add metadata or modify the data structure
func (exporter *CPUMonitorExporter) prepareExportData(data *CPUMonitorData) map[string]interface{} {
	// Convert struct to map for more flexible JSON structure
	exportData := map[string]interface{}{
		"export_info": map[string]interface{}{
			"exported_at":      time.Now().Format(time.RFC3339),
			"export_format":    "json",
			"exporter_version": "1.0.0",
			"data_type":        "cpu_monitoring",
		},
		"cpu_data": data,
	}

	return exportData
}

// SetLogsDirectory sets the base directory for log files
func (exporter *CPUMonitorExporter) SetLogsDirectory(directory string) {
	exporter.LogsDirectory = directory
}

// SetDateFormat sets the date format for file naming
func (exporter *CPUMonitorExporter) SetDateFormat(format string) {
	exporter.DateFormat = format
}

// SetPrettyPrint sets whether to pretty print JSON output
func (exporter *CPUMonitorExporter) SetPrettyPrint(pretty bool) {
	exporter.PrettyPrint = pretty
}

// SetCreateSubDirs sets whether to create subdirectories for each module
func (exporter *CPUMonitorExporter) SetCreateSubDirs(create bool) {
	exporter.CreateSubDirs = create
}

// GetExportPath returns the full path where files for a specific module are stored
func (exporter *CPUMonitorExporter) GetExportPath(moduleName string) string {
	if exporter.CreateSubDirs {
		return filepath.Join(exporter.LogsDirectory, moduleName)
	}
	return exporter.LogsDirectory
}

// ListExportedFiles returns a list of exported files for a specific module
func (exporter *CPUMonitorExporter) ListExportedFiles(moduleName string) ([]string, error) {
	exportPath := exporter.GetExportPath(moduleName)

	// Check if directory exists
	if _, err := os.Stat(exportPath); os.IsNotExist(err) {
		return []string{}, nil
	}

	// Read directory contents
	files, err := os.ReadDir(exportPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	var fileList []string
	for _, file := range files {
		if !file.IsDir() {
			fileList = append(fileList, file.Name())
		}
	}

	return fileList, nil
}

// CleanOldFiles removes files older than the specified number of days
func (exporter *CPUMonitorExporter) CleanOldFiles(moduleName string, daysToKeep int) error {
	exportPath := exporter.GetExportPath(moduleName)

	// Check if directory exists
	if _, err := os.Stat(exportPath); os.IsNotExist(err) {
		return nil
	}

	// Read directory contents
	files, err := os.ReadDir(exportPath)
	if err != nil {
		return fmt.Errorf("failed to read directory: %w", err)
	}

	cutoffTime := time.Now().AddDate(0, 0, -daysToKeep)

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		fileInfo, err := file.Info()
		if err != nil {
			continue
		}

		if fileInfo.ModTime().Before(cutoffTime) {
			filePath := filepath.Join(exportPath, file.Name())
			if err := os.Remove(filePath); err != nil {
				// Log error but continue with other files
				fmt.Printf("Warning: Failed to remove old file %s: %v\n", filePath, err)
			}
		}
	}

	return nil
}
