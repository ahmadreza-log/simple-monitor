# Contributing to Simple Monitor

Thank you for your interest in contributing to Simple Monitor! This document provides guidelines and information for contributors.

## 🤝 How to Contribute

### Reporting Issues

Before creating an issue, please:
1. Check if the issue already exists
2. Use the latest version of the project
3. Provide detailed information about the problem

When creating an issue, include:
- **OS and version**: Windows 10, Ubuntu 20.04, macOS 12, etc.
- **Go version**: `go version`
- **Steps to reproduce**: Clear, numbered steps
- **Expected behavior**: What should happen
- **Actual behavior**: What actually happens
- **Screenshots**: If applicable
- **Logs**: Any error messages or logs

### Suggesting Features

We welcome feature suggestions! Please:
1. Check if the feature is already planned
2. Provide a clear description
3. Explain the use case
4. Consider implementation complexity

## 🛠️ Development Setup

### Prerequisites

- Go 1.21 or higher
- Git
- A code editor (VS Code, GoLand, Vim, etc.)

### Getting Started

1. **Fork the repository**
   ```bash
   git clone https://github.com/yourusername/simple-monitor.git
   cd simple-monitor
   ```

2. **Create a feature branch**
   ```bash
   git checkout -b feature/your-feature-name
   ```

3. **Install dependencies**
   ```bash
   go mod tidy
   ```

4. **Run tests**
   ```bash
   go test ./...
   ```

5. **Build the project**
   ```bash
   go build -o simple-monitor main.go
   ```

## 📝 Code Style Guidelines

### Go Code Style

- Follow the official [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Use `gofmt` to format your code
- Use `golint` to check for style issues
- Use `go vet` to check for potential issues

### Naming Conventions

- **PascalCase**: For exported types, functions, and methods
- **camelCase**: For internal variables and functions
- **snake_case**: For file names (when appropriate)
- **UPPER_CASE**: For constants

### Documentation

- Document all exported functions and types
- Use clear, concise comments
- Include examples for complex functions
- Update README.md for new features

### Example Code Style

```go
// CPUMonitorManager manages CPU monitoring operations
type CPUMonitorManager struct {
    collector *CPUMonitorCollector
    displayer *CPUMonitorDisplayer
    exporter  *CPUMonitorExporter
}

// NewCPUMonitorManager creates a new instance of CPUMonitorManager
// with default collector, displayer, and exporter configurations
func NewCPUMonitorManager() *CPUMonitorManager {
    return &CPUMonitorManager{
        collector: NewCPUMonitorCollector(),
        displayer: NewCPUMonitorDisplayer(),
        exporter:  NewCPUMonitorExporter(),
    }
}
```

## 🏗️ Architecture Guidelines

### Module Structure

Each monitoring module should follow this structure:
```
module/
├── types.go           # Data structures and types
├── collector.go       # Data collection logic
├── displayer.go       # Display and formatting
├── exporter.go        # Data export functionality
└── module.go          # Main interface and manager
```

### Design Principles

1. **Separation of Concerns**: Each component has a single responsibility
2. **Interface-Based Design**: Use interfaces for better testability
3. **Configuration-Driven**: Make features configurable
4. **Error Handling**: Always handle errors appropriately
5. **Logging**: Use consistent logging throughout

### Adding New Modules

1. **Create module directory**
   ```bash
   mkdir newmodule
   ```

2. **Implement required files**
   - `types.go`: Define data structures
   - `collector.go`: Implement data collection
   - `displayer.go`: Implement display logic
   - `exporter.go`: Implement export functionality
   - `newmodule.go`: Create main interface

3. **Update main.go**
   - Add import
   - Create manager instance
   - Add menu options

4. **Add tests**
   - Unit tests for each component
   - Integration tests for the module

## 🧪 Testing

### Writing Tests

- Write tests for all new functionality
- Aim for high test coverage
- Use table-driven tests when appropriate
- Test both success and error cases

### Test Structure

```go
func TestCPUMonitorCollector_CollectData(t *testing.T) {
    tests := []struct {
        name    string
        want    *CPUMonitorData
        wantErr bool
    }{
        {
            name:    "successful collection",
            want:    &CPUMonitorData{},
            wantErr: false,
        },
        {
            name:    "collection error",
            want:    nil,
            wantErr: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            collector := NewCPUMonitorCollector()
            got, err := collector.CollectCPUMonitorData()
            
            if (err != nil) != tt.wantErr {
                t.Errorf("CollectCPUMonitorData() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            
            if !reflect.DeepEqual(got, tt.want) {
                t.Errorf("CollectCPUMonitorData() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with verbose output
go test -v ./...

# Run specific test
go test -run TestCPUMonitorCollector ./cpumonitor
```

## 📋 Pull Request Process

### Before Submitting

1. **Ensure tests pass**
   ```bash
   go test ./...
   ```

2. **Check code style**
   ```bash
   gofmt -w .
   golint ./...
   go vet ./...
   ```

3. **Update documentation**
   - Update README.md if needed
   - Add/update code comments
   - Update CHANGELOG.md

4. **Commit changes**
   ```bash
   git add .
   git commit -m "Add feature: brief description"
   ```

### Pull Request Template

When creating a PR, include:

- **Description**: What changes were made and why
- **Type**: Bug fix, feature, documentation, etc.
- **Testing**: How the changes were tested
- **Screenshots**: If applicable
- **Breaking Changes**: Any breaking changes
- **Related Issues**: Link to related issues

### PR Review Process

1. **Automated Checks**: CI/CD pipeline runs tests
2. **Code Review**: Maintainers review the code
3. **Testing**: Manual testing if needed
4. **Approval**: At least one maintainer approval required
5. **Merge**: Changes are merged to main branch

## 🐛 Bug Reports

### Before Reporting

1. Check if the bug is already reported
2. Try to reproduce the issue
3. Check the latest version
4. Look at recent changes

### Bug Report Template

```markdown
**Bug Description**
A clear description of what the bug is.

**To Reproduce**
Steps to reproduce the behavior:
1. Go to '...'
2. Click on '....'
3. Scroll down to '....'
4. See error

**Expected Behavior**
What you expected to happen.

**Actual Behavior**
What actually happened.

**Environment**
- OS: [e.g. Windows 10, Ubuntu 20.04]
- Go Version: [e.g. 1.21.0]
- Simple Monitor Version: [e.g. 0.2.0]

**Additional Context**
Any other context about the problem.
```

## 💡 Feature Requests

### Before Requesting

1. Check if the feature is already planned
2. Consider if it fits the project scope
3. Think about implementation complexity
4. Consider alternatives

### Feature Request Template

```markdown
**Feature Description**
A clear description of the feature you'd like to see.

**Use Case**
Describe the use case and why this feature would be useful.

**Proposed Solution**
Describe your proposed solution.

**Alternatives**
Describe any alternative solutions you've considered.

**Additional Context**
Any other context about the feature request.
```

## 📚 Documentation

### Code Documentation

- Document all exported functions and types
- Use clear, concise comments
- Include examples for complex functions
- Follow Go documentation conventions

### README Updates

- Update feature list for new features
- Add usage examples
- Update installation instructions
- Add configuration examples

### API Documentation

- Document all public APIs
- Include parameter descriptions
- Provide usage examples
- Document return values and errors

## 🏷️ Release Process

### Version Numbering

We follow [Semantic Versioning](https://semver.org/):
- **MAJOR**: Incompatible API changes
- **MINOR**: New functionality (backward compatible)
- **PATCH**: Bug fixes (backward compatible)

### Release Checklist

- [ ] All tests pass
- [ ] Documentation updated
- [ ] CHANGELOG.md updated
- [ ] Version numbers updated
- [ ] Release notes prepared
- [ ] Tag created
- [ ] Release published

## 🤔 Questions?

If you have questions about contributing:

- **GitHub Discussions**: For general questions
- **GitHub Issues**: For specific problems
- **Email**: contributors@simplemonitor.dev

## 🙏 Recognition

Contributors will be recognized in:
- README.md contributors section
- Release notes
- Project documentation

Thank you for contributing to Simple Monitor! 🎉
