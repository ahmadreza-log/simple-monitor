package processmonitor

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// ProcessMonitorExporter handles exporting process monitoring data to various formats
// This struct provides methods to save process monitoring data to files
type ProcessMonitorExporter struct {
	// Export configuration
	LogsDirectory string // Base directory for log files
	DateFormat    string // Date format for file naming
	CreateSubDirs bool   // Whether to create subdirectories for each module
	PrettyPrint   bool   // Whether to pretty print JSON output
}

// NewProcessMonitorExporter creates a new instance of ProcessMonitorExporter
// with default configuration values
func NewProcessMonitorExporter() *ProcessMonitorExporter {
	return &ProcessMonitorExporter{
		LogsDirectory: "logs",
		DateFormat:    "2006-01-02",
		CreateSubDirs: true,
		PrettyPrint:   true,
	}
}

// ExportToJSON exports process monitoring data to a JSON file
// The file will be saved with a timestamp-based filename in the appropriate subdirectory
func (exporter *ProcessMonitorExporter) ExportToJSON(data *ProcessMonitorData, moduleName string) (string, error) {
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

// ExportToCSV exports process monitoring data to a CSV file
// This creates a simplified CSV format with key metrics
func (exporter *ProcessMonitorExporter) ExportToCSV(data *ProcessMonitorData, moduleName string) (string, error) {
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

// ExportToTXT exports process monitoring data to a text file
// This creates a human-readable text format
func (exporter *ProcessMonitorExporter) ExportToTXT(data *ProcessMonitorData, moduleName string) (string, error) {
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

// generateCSVContent creates CSV content from process monitoring data
func (exporter *ProcessMonitorExporter) generateCSVContent(data *ProcessMonitorData) string {
	var content string

	// Header
	content += "Timestamp,Total Processes,Running,Sleeping,Zombie,Stopped,Total CPU,Total Memory,Total Threads,Process Status\n"

	// Data row
	content += fmt.Sprintf("%s,%d,%d,%d,%d,%d,%.2f,%.2f,%d,%s\n",
		data.Timestamp.Format("2006-01-02 15:04:05"),
		data.TotalProcesses,
		data.RunningProcesses,
		data.SleepingProcesses,
		data.ZombieProcesses,
		data.StoppedProcesses,
		data.TotalCPUUsage,
		data.TotalMemoryUsage,
		data.TotalThreads,
		data.ProcessStatus)

	// Process data
	if len(data.ProcessInfos) > 0 {
		content += "\nProcess Data\n"
		content += "PID,Name,Status,User,CPU%,Memory%,Threads,Open Files,Priority,Parent PID,Command Line\n"
		for _, proc := range data.ProcessInfos {
			content += fmt.Sprintf("%d,%s,%s,%s,%.2f,%.2f,%d,%d,%d,%d,%s\n",
				proc.PID,
				proc.Name,
				proc.Status,
				proc.User,
				proc.CPUUsage,
				proc.MemoryUsage,
				proc.Threads,
				proc.OpenFiles,
				proc.Priority,
				proc.ParentPID,
				proc.CommandLine)
		}
	}

	// Top CPU processes
	if len(data.TopCPUProcesses) > 0 {
		content += "\nTop CPU Processes\n"
		content += "PID,Name,CPU%,Memory%,Threads,Status\n"
		for _, proc := range data.TopCPUProcesses {
			content += fmt.Sprintf("%d,%s,%.2f,%.2f,%d,%s\n",
				proc.PID,
				proc.Name,
				proc.CPUUsage,
				proc.MemoryUsage,
				proc.Threads,
				proc.Status)
		}
	}

	// Top memory processes
	if len(data.TopMemoryProcesses) > 0 {
		content += "\nTop Memory Processes\n"
		content += "PID,Name,CPU%,Memory%,Threads,Status\n"
		for _, proc := range data.TopMemoryProcesses {
			content += fmt.Sprintf("%d,%s,%.2f,%.2f,%d,%s\n",
				proc.PID,
				proc.Name,
				proc.CPUUsage,
				proc.MemoryUsage,
				proc.Threads,
				proc.Status)
		}
	}

	// Process alerts
	if len(data.ProcessAlerts) > 0 {
		content += "\nProcess Alerts\n"
		content += "PID,Name,Alert Type,Severity,Value,Threshold,Timestamp\n"
		for _, alert := range data.ProcessAlerts {
			content += fmt.Sprintf("%d,%s,%s,%s,%.2f,%.2f,%s\n",
				alert.PID,
				alert.Name,
				alert.AlertType,
				alert.Severity,
				alert.Value,
				alert.Threshold,
				alert.Timestamp.Format("2006-01-02 15:04:05"))
		}
	}

	return content
}

// generateTXTContent creates text content from process monitoring data
func (exporter *ProcessMonitorExporter) generateTXTContent(data *ProcessMonitorData) string {
	var content string

	// Header
	content += "PROCESS MONITOR REPORT\n"
	content += "=====================\n\n"

	// Timestamp
	content += fmt.Sprintf("Generated: %s\n\n", data.Timestamp.Format("2006-01-02 15:04:05"))

	// Process summary
	content += "PROCESS SUMMARY\n"
	content += "---------------\n"
	content += fmt.Sprintf("Total Processes: %d\n", data.TotalProcesses)
	content += fmt.Sprintf("Running: %d\n", data.RunningProcesses)
	content += fmt.Sprintf("Sleeping: %d\n", data.SleepingProcesses)
	content += fmt.Sprintf("Zombie: %d\n", data.ZombieProcesses)
	content += fmt.Sprintf("Stopped: %d\n", data.StoppedProcesses)
	content += fmt.Sprintf("Total CPU Usage: %.2f%%\n", data.TotalCPUUsage)
	content += fmt.Sprintf("Total Memory Usage: %.2f%%\n", data.TotalMemoryUsage)
	content += fmt.Sprintf("Total Threads: %d\n", data.TotalThreads)
	content += fmt.Sprintf("Total Open Files: %d\n", data.TotalOpenFiles)
	content += fmt.Sprintf("Process Status: %s\n\n", data.ProcessStatus)

	// Top CPU processes
	if len(data.TopCPUProcesses) > 0 {
		content += "TOP CPU PROCESSES\n"
		content += "-----------------\n"
		content += "PID\tName\t\t\tCPU%\tMemory%\tThreads\tStatus\n"
		content += "---\t----\t\t\t----\t-------\t-------\t------\n"

		for _, proc := range data.TopCPUProcesses {
			content += fmt.Sprintf("%d\t%-20s\t%.2f\t%.2f\t%d\t%s\n",
				proc.PID,
				proc.Name,
				proc.CPUUsage,
				proc.MemoryUsage,
				proc.Threads,
				proc.Status)
		}
		content += "\n"
	}

	// Top memory processes
	if len(data.TopMemoryProcesses) > 0 {
		content += "TOP MEMORY PROCESSES\n"
		content += "--------------------\n"
		content += "PID\tName\t\t\tCPU%\tMemory%\tThreads\tStatus\n"
		content += "---\t----\t\t\t----\t-------\t-------\t------\n"

		for _, proc := range data.TopMemoryProcesses {
			content += fmt.Sprintf("%d\t%-20s\t%.2f\t%.2f\t%d\t%s\n",
				proc.PID,
				proc.Name,
				proc.CPUUsage,
				proc.MemoryUsage,
				proc.Threads,
				proc.Status)
		}
		content += "\n"
	}

	// Top I/O processes
	if len(data.TopIOProcesses) > 0 {
		content += "TOP I/O PROCESSES\n"
		content += "-----------------\n"
		content += "PID\tName\t\t\tCPU%\tMemory%\tThreads\tStatus\n"
		content += "---\t----\t\t\t----\t-------\t-------\t------\n"

		for _, proc := range data.TopIOProcesses {
			content += fmt.Sprintf("%d\t%-20s\t%.2f\t%.2f\t%d\t%s\n",
				proc.PID,
				proc.Name,
				proc.CPUUsage,
				proc.MemoryUsage,
				proc.Threads,
				proc.Status)
		}
		content += "\n"
	}

	// Top thread processes
	if len(data.TopThreadProcesses) > 0 {
		content += "TOP THREAD PROCESSES\n"
		content += "--------------------\n"
		content += "PID\tName\t\t\tCPU%\tMemory%\tThreads\tStatus\n"
		content += "---\t----\t\t\t----\t-------\t-------\t------\n"

		for _, proc := range data.TopThreadProcesses {
			content += fmt.Sprintf("%d\t%-20s\t%.2f\t%.2f\t%d\t%s\n",
				proc.PID,
				proc.Name,
				proc.CPUUsage,
				proc.MemoryUsage,
				proc.Threads,
				proc.Status)
		}
		content += "\n"
	}

	// Process tree
	if len(data.ProcessTree) > 0 {
		content += "PROCESS TREE\n"
		content += "------------\n"
		exporter.addTreeToContent(data.ProcessTree, &content, 0)
		content += "\n"
	}

	// Process alerts
	if len(data.ProcessAlerts) > 0 {
		content += "PROCESS ALERTS\n"
		content += "--------------\n"
		content += "PID\tName\t\t\tAlert Type\t\tSeverity\tValue\tThreshold\n"
		content += "---\t----\t\t\t----------\t\t--------\t-----\t---------\n"

		for _, alert := range data.ProcessAlerts {
			content += fmt.Sprintf("%d\t%-20s\t%-20s\t%s\t%.2f\t%.2f\n",
				alert.PID,
				alert.Name,
				alert.AlertType,
				alert.Severity,
				alert.Value,
				alert.Threshold)
		}
		content += "\n"
	}

	// Alerts
	content += "ALERTS & WARNINGS\n"
	content += "-----------------\n"
	content += fmt.Sprintf("High CPU Warning: %t\n", data.HighCPUWarning)
	content += fmt.Sprintf("High Memory Warning: %t\n", data.HighMemoryWarning)
	content += fmt.Sprintf("High I/O Warning: %t\n", data.HighIOWarning)
	content += fmt.Sprintf("Zombie Warning: %t\n", data.ZombieWarning)
	content += fmt.Sprintf("Thread Warning: %t\n", data.ThreadWarning)
	content += "\n"

	return content
}

// addTreeToContent recursively adds process tree to content
func (exporter *ProcessMonitorExporter) addTreeToContent(tree []ProcessTreeInfo, content *string, level int) {
	for _, node := range tree {
		// Indent based on level
		indent := strings.Repeat("  ", level)

		// Tree character
		treeChar := "├─"
		if node.IsLeaf {
			treeChar = "└─"
		}

		*content += fmt.Sprintf("%s%s%s (PID: %d)\n",
			indent,
			treeChar,
			node.Name,
			node.PID)

		// Add children
		if len(node.Children) > 0 {
			exporter.addTreeToContent(node.Children, content, level+1)
		}
	}
}

// SetLogsDirectory sets the logs directory
func (exporter *ProcessMonitorExporter) SetLogsDirectory(dir string) {
	exporter.LogsDirectory = dir
}

// SetPrettyPrint sets whether to pretty print JSON output
func (exporter *ProcessMonitorExporter) SetPrettyPrint(pretty bool) {
	exporter.PrettyPrint = pretty
}

// SetCreateSubDirs sets whether to create subdirectories for each module
func (exporter *ProcessMonitorExporter) SetCreateSubDirs(create bool) {
	exporter.CreateSubDirs = create
}
