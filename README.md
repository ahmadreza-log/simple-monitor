# 🖥️ Simple Monitor

A comprehensive system monitoring tool built with Go that provides real-time insights into your system's performance, CPU usage, memory consumption, disk usage, network activity, and process monitoring.

## ✨ Features

### 🔧 System Information
- **Basic System Info**: Hostname, OS, Architecture, Kernel version
- **Uptime Tracking**: System uptime and boot time
- **Hardware Details**: CPU model, cores, memory specifications

### 🖥️ CPU Monitoring
- **Live CPU Monitoring**: Real-time CPU usage with graphical display
- **Per-Core Analysis**: Individual core usage tracking
- **Process Monitoring**: Top CPU-consuming processes
- **Temperature Monitoring**: CPU temperature tracking with alerts
- **Load Average**: 1-minute, 5-minute, and 15-minute load averages
- **Graphical Display**: Color-coded progress bars and charts

### 💾 Memory Monitoring
- **RAM Usage**: Total, used, available, and free memory
- **Swap Space**: Swap usage and statistics
- **Memory Details**: Cache and buffer information

### 💿 Disk Monitoring
- **Disk Usage**: Capacity and usage for all drives
- **Performance Metrics**: Read/write speeds
- **Disk Health**: SSD/HDD detection, removable drive support

### 🌐 Network Monitoring
- **Interface Status**: Network interface information
- **Traffic Statistics**: Bytes sent/received, packet counts
- **IP Configuration**: IP addresses, subnet masks, gateways

### ⚙️ Process Monitoring
- **Process List**: Running processes with CPU and memory usage
- **Process Details**: PID, name, status, priority
- **Thread Information**: Thread count per process

### 🚀 Quick Test Feature
- **Simultaneous Monitoring**: Monitor all systems at once
- **Real-time Updates**: Live data refresh every 2 seconds
- **Compact Display**: Essential information in one view
- **Easy Exit**: Press Ctrl+C to stop anytime

### ⚙️ Advanced Settings
- **Display Settings**: Refresh rate, format, colors, screen size
- **Monitoring Settings**: Intervals, auto-start, data retention, alerts
- **Performance Settings**: CPU priority, memory limits, background mode
- **Log Settings**: Log level, rotation, directory management
- **Reset to Defaults**: Restore all settings to factory defaults

### 👨‍💻 Developer Tools
- **Performance Analysis**: Detailed system performance metrics
- **Debug Mode**: Enhanced logging and error information
- **Export Debug Info**: Export system information for troubleshooting
- **Log Management**: View, clear, and manage log files
- **Configuration Viewer**: Display all monitor configurations
- **Test All Monitors**: Comprehensive testing of all components

## 🚀 Quick Start

### Prerequisites
- Go 1.21 or higher
- Windows, Linux, or macOS

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/ahmadreza-log/simple-monitor.git
   cd simple-monitor
   ```

2. **Build the application**
   ```bash
   go build -o simple-monitor main.go
   ```

3. **Run the application**
   ```bash
   ./simple-monitor
   ```

### Development Setup

1. **Install dependencies**
   ```bash
   go mod tidy
   ```

2. **Run in development mode**
   ```bash
   go run main.go
   ```

## 📖 Usage

### Main Menu
```
🖥️  Simple Monitor v1.0
------------------------------
1. Start Monitoring
2. Settings
3. Developer
4. Quit
------------------------------
```

### Monitoring Options
```
📊 Monitoring Options
------------------------------
1. System Information
2. CPU Monitor
3. Memory Monitor
4. Disk Monitor
5. Network Monitor
6. Process Monitor
7. Quick Test (All Monitors)
8. Back to Main Menu
------------------------------
```

### Quick Test Feature
```
🚀 Quick Test - All Monitors
------------------------------
Last Update: 15:30:45

🖥️  CPU:
  Usage: 45.2%

💾 Memory:
  Used: 8.5 GB / 16.0 GB (53.1%)

💿 Disk:
  C: 120.5 GB / 500.0 GB (24.1%)

🌐 Network:
  Interfaces: 3
  Ethernet: RX 125.3 MB, TX 89.7 MB

Press Ctrl+C to stop...
```

### Settings Menu
```
⚙️  Settings
------------------------------
1. Export Settings
2. Display Settings
3. Monitoring Settings
4. Performance Settings
5. Log Settings
6. Reset to Defaults
7. Back to Main Menu
------------------------------
```

### Developer Menu
```
👨‍💻 Developer Section
------------------------------
1. View System Information
2. View Log Files
3. Clear Log Files
4. View Configuration
5. Test All Monitors
6. Performance Analysis
7. Debug Mode
8. Export Debug Info
9. Back to Main Menu
------------------------------
```

## 🏗️ Architecture

### Project Structure
```
simple-monitor/
├── main.go                 # Main application entry point
├── go.mod                  # Go module definition
├── go.sum                  # Go module checksums
├── .gitignore             # Git ignore rules
├── README.md              # Project documentation
├── CHANGELOG.md           # Version history
├── CONTRIBUTING.md        # Contribution guidelines
├── LICENSE                # MIT License
├── Makefile              # Build automation
├── build.ps1             # PowerShell build script
├── logs/                 # Log files directory
│   ├── systeminfo/       # System info exports
│   ├── cpumonitor/       # CPU monitor exports
│   ├── memorymonitor/    # Memory monitor exports
│   ├── diskmonitor/      # Disk monitor exports
│   ├── networkmonitor/   # Network monitor exports
│   └── processmonitor/    # Process monitor exports
├── systeminfo/           # System information module
│   ├── types.go          # Data structures
│   ├── collector.go      # Data collection
│   ├── displayer.go      # Data display
│   ├── exporter.go       # Data export
│   └── systeminfo.go     # Main interface
├── cpumonitor/           # CPU monitoring module
│   ├── types.go         # Data structures
│   ├── collector.go     # Data collection
│   ├── displayer.go     # Data display
│   ├── exporter.go      # Data export
│   └── cpumonitor.go    # Main interface
├── memorymonitor/       # Memory monitoring module
├── diskmonitor/         # Disk monitoring module
├── networkmonitor/      # Network monitoring module
└── processmonitor/      # Process monitoring module
```

### Design Patterns

- **Modular Architecture**: Each monitoring feature is a separate module
- **Separation of Concerns**: Collector, Displayer, and Exporter are separate
- **Interface-Based Design**: Clean interfaces for easy testing and extension
- **Configuration-Driven**: Configurable options for all features
- **Signal Handling**: Proper Ctrl+C handling for graceful shutdown

## 🔧 Configuration

### CPU Monitor Settings
```go
config := &CPUMonitorConfig{
    RefreshInterval:     1 * time.Second,
    MaxProcesses:        20,
    TemperatureWarning:  70.0,
    TemperatureCritical: 85.0,
    ShowCores:           true,
    ShowProcesses:       true,
    ShowTemperature:     true,
    ShowLoadAverage:     true,
    ExportToFile:        false,
    ExportInterval:      30 * time.Second,
    ExportFormat:        "json",
    MinCPUUsage:         1.0,
    ProcessNameFilter:   "",
}
```

### Export Settings
```go
exporter.SetLogsDirectory("logs")
exporter.SetPrettyPrint(true)
exporter.SetCreateSubDirs(true)
```

## 📊 Data Export

### Supported Formats
- **JSON**: Structured data export with metadata
- **CSV**: Tabular data export
- **Text**: Human-readable format

### Export Structure
```json
{
  "export_info": {
    "exported_at": "2025-09-27T15:12:19+03:30",
    "export_format": "json",
    "exporter_version": "1.0.0",
    "data_type": "cpu_monitoring"
  },
  "cpu_data": {
    "model_name": "Intel Core i7-8700K",
    "architecture": "amd64",
    "physical_cores": 6,
    "logical_cores": 12,
    "overall_usage": 75.5,
    "cores": [...],
    "top_processes": [...],
    "timestamp": "2025-09-27T15:12:19.6308837+03:30"
  }
}
```

## 🎨 Display Features

### Color Coding
- 🟢 **Green**: Low usage (< 30%)
- 🟡 **Yellow**: Medium usage (30-60%)
- 🟣 **Purple**: High usage (60-80%)
- 🔴 **Red**: Critical usage (> 80%)

### Graphical Elements
- **Progress Bars**: Visual representation of usage percentages
- **Grid Layout**: Organized display of multiple cores
- **Real-time Updates**: Live refreshing of data
- **Responsive Design**: Adapts to different terminal sizes

## 🛠️ Development

### Adding New Modules

1. **Create module directory**
   ```bash
   mkdir newmodule
   ```

2. **Implement required files**
   - `types.go`: Data structures
   - `collector.go`: Data collection logic
   - `displayer.go`: Display formatting
   - `exporter.go`: Export functionality
   - `newmodule.go`: Main interface

3. **Update main.go**
   - Add import
   - Create manager instance
   - Add menu options

### Code Style

- **PascalCase**: For all exported types and functions
- **camelCase**: For internal variables and functions
- **Comprehensive Comments**: All public functions documented
- **Error Handling**: Proper error propagation and logging

## 🚀 Key Features

### ✅ Completed Features
- [x] System Information module
- [x] CPU Monitor with live updates
- [x] Memory Monitor with real-time tracking
- [x] Disk Monitor with usage statistics
- [x] Network Monitor with traffic analysis
- [x] Process Monitor with detailed information
- [x] JSON export functionality
- [x] Modular architecture
- [x] Quick Test feature for simultaneous monitoring
- [x] Advanced Settings with comprehensive options
- [x] Developer Tools for debugging and analysis
- [x] Proper Ctrl+C signal handling
- [x] Configuration management
- [x] Log file management

### 🔄 Future Enhancements
- [ ] Web dashboard interface
- [ ] Historical data analysis
- [ ] Alert system with notifications
- [ ] Plugin system for custom monitors
- [ ] Cross-platform optimizations
- [ ] Performance profiling tools
- [ ] Data visualization charts
- [ ] Remote monitoring capabilities

## 🤝 Contributing

We welcome contributions! Please follow these steps:

1. **Fork the repository**
2. **Create a feature branch**
   ```bash
   git checkout -b feature/amazing-feature
   ```
3. **Commit your changes**
   ```bash
   git commit -m 'Add some amazing feature'
   ```
4. **Push to the branch**
   ```bash
   git push origin feature/amazing-feature
   ```
5. **Open a Pull Request**

### Contribution Guidelines

- Follow the existing code style
- Add tests for new features
- Update documentation
- Ensure all tests pass
- Add appropriate comments

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Go community for excellent documentation
- Open source monitoring tools for inspiration
- Contributors and testers

## 📞 Support

- **Issues**: [GitHub Issues](https://github.com/ahmadreza-log/simple-monitor/issues)
- **Discussions**: [GitHub Discussions](https://github.com/ahmadreza-log/simple-monitor/discussions)
- **Email**: ahmadreza.log@gmail.com

## 🔗 Links

- **Documentation**: [Wiki](https://github.com/ahmadreza-log/simple-monitor/wiki)
- **Changelog**: [CHANGELOG.md](CHANGELOG.md)
- **Contributing**: [CONTRIBUTING.md](CONTRIBUTING.md)

---

<div align="center">

**Made with ❤️ and Go by [@ahmadreza-log](https://github.com/ahmadreza-log)**

[⭐ Star this repo](https://github.com/ahmadreza-log/simple-monitor) | [🐛 Report Bug](https://github.com/ahmadreza-log/simple-monitor/issues) | [💡 Request Feature](https://github.com/ahmadreza-log/simple-monitor/issues)

</div>