package networkmonitor

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// NetworkMonitorExporter handles exporting network monitoring data to various formats
// This struct provides methods to save network monitoring data to files
type NetworkMonitorExporter struct {
	// Export configuration
	LogsDirectory string // Base directory for log files
	DateFormat    string // Date format for file naming
	CreateSubDirs bool   // Whether to create subdirectories for each module
	PrettyPrint   bool   // Whether to pretty print JSON output
}

// NewNetworkMonitorExporter creates a new instance of NetworkMonitorExporter
// with default configuration values
func NewNetworkMonitorExporter() *NetworkMonitorExporter {
	return &NetworkMonitorExporter{
		LogsDirectory: "logs",
		DateFormat:    "2006-01-02",
		CreateSubDirs: true,
		PrettyPrint:   true,
	}
}

// ExportToJSON exports network monitoring data to a JSON file
// The file will be saved with a timestamp-based filename in the appropriate subdirectory
func (exporter *NetworkMonitorExporter) ExportToJSON(data *NetworkMonitorData, moduleName string) (string, error) {
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

// ExportToCSV exports network monitoring data to a CSV file
// This creates a simplified CSV format with key metrics
func (exporter *NetworkMonitorExporter) ExportToCSV(data *NetworkMonitorData, moduleName string) (string, error) {
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

// ExportToTXT exports network monitoring data to a text file
// This creates a human-readable text format
func (exporter *NetworkMonitorExporter) ExportToTXT(data *NetworkMonitorData, moduleName string) (string, error) {
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

// generateCSVContent creates CSV content from network monitoring data
func (exporter *NetworkMonitorExporter) generateCSVContent(data *NetworkMonitorData) string {
	var content string

	// Header
	content += "Timestamp,Total Sent,Total Received,Total Throughput,Network Status,Average Latency,Packet Loss Rate,Network Utilization\n"

	// Data row
	content += fmt.Sprintf("%s,%d,%d,%.2f,%s,%.2f,%.2f,%.2f\n",
		data.Timestamp.Format("2006-01-02 15:04:05"),
		data.TotalBytesSent,
		data.TotalBytesRecv,
		data.TotalThroughput,
		data.NetworkStatus,
		data.AverageLatency,
		data.PacketLossRate,
		data.NetworkUtilization)

	// Interface data
	if len(data.Interfaces) > 0 {
		content += "\nInterface Data\n"
		content += "Name,Type,Status,IP Address,MAC Address,Speed,Is Up,Is Loopback,Is Virtual\n"
		for _, iface := range data.Interfaces {
			content += fmt.Sprintf("%s,%s,%s,%s,%s,%d,%t,%t,%t\n",
				iface.Name,
				iface.Type,
				iface.Status,
				iface.IPAddress,
				iface.MACAddress,
				iface.Speed,
				iface.IsUp,
				iface.IsLoopback,
				iface.IsVirtual)
		}
	}

	// I/O data
	if len(data.InterfaceIO) > 0 {
		content += "\nI/O Data\n"
		content += "Interface,Send Speed,Recv Speed,Total Speed,Utilization,Packets Sent,Packets Recv,Send Errors,Recv Errors\n"
		for _, io := range data.InterfaceIO {
			content += fmt.Sprintf("%s,%.2f,%.2f,%.2f,%.2f,%d,%d,%d,%d\n",
				io.InterfaceName,
				io.SendSpeed,
				io.RecvSpeed,
				io.TotalSpeed,
				io.Utilization,
				io.PacketsSent,
				io.PacketsRecv,
				io.SendErrors,
				io.RecvErrors)
		}
	}

	// Connection data
	if len(data.Connections) > 0 {
		content += "\nConnection Data\n"
		content += "Local Address,Remote Address,Type,Status,PID,Process Name,User,State,Family\n"
		for _, conn := range data.Connections {
			content += fmt.Sprintf("%s,%s,%s,%s,%d,%s,%s,%s,%s\n",
				conn.LocalAddress,
				conn.RemoteAddress,
				conn.Type,
				conn.Status,
				conn.PID,
				conn.ProcessName,
				conn.User,
				conn.State,
				conn.Family)
		}
	}

	// Process data
	if len(data.TopProcesses) > 0 {
		content += "\nProcess Data\n"
		content += "PID,Name,Send Speed,Recv Speed,Total Speed,Connections,Status,User\n"
		for _, process := range data.TopProcesses {
			content += fmt.Sprintf("%d,%s,%.2f,%.2f,%.2f,%d,%s,%s\n",
				process.PID,
				process.Name,
				process.SendSpeed,
				process.RecvSpeed,
				process.TotalSpeed,
				process.Connections,
				process.Status,
				process.User)
		}
	}

	// Latency data
	if len(data.LatencyInfo) > 0 {
		content += "\nLatency Data\n"
		content += "Target,Latency,Packet Loss,Status,Last Checked\n"
		for _, latency := range data.LatencyInfo {
			content += fmt.Sprintf("%s,%.2f,%.2f,%s,%s\n",
				latency.Target,
				latency.Latency,
				latency.PacketLoss,
				latency.Status,
				latency.LastChecked.Format("2006-01-02 15:04:05"))
		}
	}

	return content
}

// generateTXTContent creates text content from network monitoring data
func (exporter *NetworkMonitorExporter) generateTXTContent(data *NetworkMonitorData) string {
	var content string

	// Header
	content += "NETWORK MONITOR REPORT\n"
	content += "=====================\n\n"

	// Timestamp
	content += fmt.Sprintf("Generated: %s\n\n", data.Timestamp.Format("2006-01-02 15:04:05"))

	// Network summary
	content += "NETWORK SUMMARY\n"
	content += "---------------\n"
	content += fmt.Sprintf("Total Sent: %s\n", exporter.formatBytes(data.TotalBytesSent))
	content += fmt.Sprintf("Total Received: %s\n", exporter.formatBytes(data.TotalBytesRecv))
	content += fmt.Sprintf("Total Throughput: %.2f Mbps\n", data.TotalThroughput)
	content += fmt.Sprintf("Network Status: %s\n", data.NetworkStatus)
	content += fmt.Sprintf("Average Latency: %.2f ms\n", data.AverageLatency)
	content += fmt.Sprintf("Packet Loss Rate: %.2f%%\n", data.PacketLossRate)
	content += fmt.Sprintf("Network Utilization: %.2f%%\n\n", data.NetworkUtilization)

	// Interface information
	if len(data.Interfaces) > 0 {
		content += "INTERFACE INFORMATION\n"
		content += "--------------------\n"
		content += "Name\t\tType\t\tStatus\tIP Address\t\tMAC Address\t\tSpeed\n"
		content += "----\t\t----\t\t------\t----------\t\t------------\t\t-----\n"

		for _, iface := range data.Interfaces {
			content += fmt.Sprintf("%s\t\t%s\t\t%s\t%s\t\t%s\t\t%d\n",
				iface.Name,
				iface.Type,
				iface.Status,
				iface.IPAddress,
				iface.MACAddress,
				iface.Speed)
		}
		content += "\n"
	}

	// I/O statistics
	if len(data.InterfaceIO) > 0 {
		content += "I/O STATISTICS\n"
		content += "--------------\n"
		content += "Interface\t\tSend Speed\tRecv Speed\tTotal Speed\tUtilization\n"
		content += "--------\t\t----------\t----------\t-----------\t------------\n"

		for _, io := range data.InterfaceIO {
			content += fmt.Sprintf("%s\t\t%.2f Mbps\t%.2f Mbps\t%.2f Mbps\t%.2f%%\n",
				io.InterfaceName,
				io.SendSpeed,
				io.RecvSpeed,
				io.TotalSpeed,
				io.Utilization)
		}
		content += "\n"
	}

	// Connection information
	if len(data.Connections) > 0 {
		content += "CONNECTION INFORMATION\n"
		content += "---------------------\n"
		content += "Local Address\t\tRemote Address\t\tType\tStatus\tProcess\n"
		content += "-------------\t\t--------------\t\t----\t------\t-------\n"

		for _, conn := range data.Connections {
			content += fmt.Sprintf("%s\t\t%s\t\t%s\t%s\t%s\n",
				conn.LocalAddress,
				conn.RemoteAddress,
				conn.Type,
				conn.Status,
				conn.ProcessName)
		}
		content += "\n"
	}

	// Latency information
	if len(data.LatencyInfo) > 0 {
		content += "LATENCY INFORMATION\n"
		content += "------------------\n"
		content += "Target\t\tLatency\t\tPacket Loss\tStatus\n"
		content += "------\t\t-------\t\t-----------\t------\n"

		for _, latency := range data.LatencyInfo {
			content += fmt.Sprintf("%s\t\t%.2f ms\t\t%.2f%%\t\t%s\n",
				latency.Target,
				latency.Latency,
				latency.PacketLoss,
				latency.Status)
		}
		content += "\n"
	}

	// Bandwidth information
	content += "BANDWIDTH INFORMATION\n"
	content += "--------------------\n"
	content += fmt.Sprintf("Total Bandwidth: %.2f Mbps\n", data.BandwidthInfo.TotalBandwidth)
	content += fmt.Sprintf("Used Bandwidth: %.2f Mbps\n", data.BandwidthInfo.UsedBandwidth)
	content += fmt.Sprintf("Available Bandwidth: %.2f Mbps\n", data.BandwidthInfo.AvailableBandwidth)
	content += fmt.Sprintf("Utilization: %.2f%%\n", data.BandwidthInfo.Utilization)
	content += fmt.Sprintf("Peak Usage: %.2f Mbps\n", data.BandwidthInfo.PeakUsage)
	content += fmt.Sprintf("Average Usage: %.2f Mbps\n\n", data.BandwidthInfo.AverageUsage)

	// Top processes
	if len(data.TopProcesses) > 0 {
		content += "TOP NETWORK PROCESSES\n"
		content += "--------------------\n"
		content += "PID\tName\t\t\tSend Speed\tRecv Speed\tTotal Speed\tConnections\n"
		content += "---\t----\t\t\t----------\t----------\t-----------\t-----------\n"

		for _, process := range data.TopProcesses {
			content += fmt.Sprintf("%d\t%-20s\t%.2f Mbps\t%.2f Mbps\t%.2f Mbps\t%d\n",
				process.PID,
				process.Name,
				process.SendSpeed,
				process.RecvSpeed,
				process.TotalSpeed,
				process.Connections)
		}
		content += "\n"
	}

	// Alerts
	content += "ALERTS & WARNINGS\n"
	content += "-----------------\n"
	content += fmt.Sprintf("High Latency Warning: %t\n", data.HighLatencyWarning)
	content += fmt.Sprintf("Packet Loss Warning: %t\n", data.PacketLossWarning)
	content += fmt.Sprintf("Bandwidth Warning: %t\n", data.BandwidthWarning)
	content += fmt.Sprintf("Connection Warning: %t\n", data.ConnectionWarning)
	content += "\n"

	return content
}

// formatBytes formats bytes into human-readable format
func (exporter *NetworkMonitorExporter) formatBytes(bytes uint64) string {
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
func (exporter *NetworkMonitorExporter) SetLogsDirectory(dir string) {
	exporter.LogsDirectory = dir
}

// SetPrettyPrint sets whether to pretty print JSON output
func (exporter *NetworkMonitorExporter) SetPrettyPrint(pretty bool) {
	exporter.PrettyPrint = pretty
}

// SetCreateSubDirs sets whether to create subdirectories for each module
func (exporter *NetworkMonitorExporter) SetCreateSubDirs(create bool) {
	exporter.CreateSubDirs = create
}
