# CLI Testing Strategies for DevGen with Charm UI

A comprehensive guide to testing terminal user interfaces, CLI applications, and interactive components built with the Charm ecosystem.

## Table of Contents

- [Testing Philosophy](#testing-philosophy)
- [Test Types and Structure](#test-types-and-structure)
- [Bubble Tea Testing](#bubble-tea-testing)
- [UI Component Testing](#ui-component-testing)
- [Integration Testing](#integration-testing)
- [Performance Testing](#performance-testing)
- [Accessibility Testing](#accessibility-testing)
- [CI/CD Integration](#cicd-integration)
- [Testing Tools and Libraries](#testing-tools-and-libraries)

---

## Testing Philosophy

### Core Principles

1. **Test User Interactions**: Focus on how users interact with the CLI, not just code coverage
2. **Visual Regression Testing**: Ensure UI layouts remain consistent across changes
3. **Terminal Compatibility**: Test across different terminal emulators and capabilities
4. **State Management**: Verify application state transitions and data flow
5. **Error Scenarios**: Test edge cases, network failures, and invalid inputs

### Testing Pyramid for CLI Applications

```
    ┌─────────────────────┐
    │   E2E Tests (5%)    │  Full application workflows
    │                     │  Real terminal interactions
    └─────────────────────┘
    ┌─────────────────────┐
    │ Integration (25%)   │  Component interactions
    │                     │  Message handling
    └─────────────────────┘
    ┌─────────────────────┐
    │  Unit Tests (70%)   │  Individual functions
    │                     │  Business logic
    └─────────────────────┘
```

---

## Test Types and Structure

### Unit Tests

Test individual functions and business logic without UI components.

```go
// config_test.go
package main

import (
    "testing"
    "os"
    "path/filepath"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestLoadPlaybook(t *testing.T) {
    tests := []struct {
        name        string
        yamlContent string
        expectError bool
        validate    func(*testing.T, *Playbook)
    }{
        {
            name: "valid playbook",
            yamlContent: `
name: "Test Playbook"
version: "1.0.0"
branches:
  - name: "test-branch"
    steps:
      - agent: "test-agent"
        action: "test action"
        condition: "start"
`,
            expectError: false,
            validate: func(t *testing.T, p *Playbook) {
                assert.Equal(t, "Test Playbook", p.Name)
                assert.Len(t, p.Branches, 1)
                assert.Equal(t, "test-branch", p.Branches[0].Name)
            },
        },
        {
            name: "invalid yaml",
            yamlContent: `
name: "Test Playbook"
branches:
  - name: "test-branch"
    steps:
      - agent: # missing value
`,
            expectError: true,
        },
        {
            name: "missing required fields",
            yamlContent: `
version: "1.0.0"
# missing name
branches: []
`,
            expectError: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Create temporary file
            tmpDir := t.TempDir()
            tmpFile := filepath.Join(tmpDir, "test.yaml")

            err := os.WriteFile(tmpFile, []byte(tt.yamlContent), 0644)
            require.NoError(t, err)

            // Test loading
            playbook, err := loadPlaybook(tmpFile)

            if tt.expectError {
                assert.Error(t, err)
                assert.Nil(t, playbook)
            } else {
                assert.NoError(t, err)
                require.NotNil(t, playbook)
                if tt.validate != nil {
                    tt.validate(t, playbook)
                }
            }
        })
    }
}

func TestEngineExecution(t *testing.T) {
    playbook := &Playbook{
        Name: "Test Execution",
        Branches: []Branch{
            {
                Name: "sequential",
                Steps: []Step{
                    {Agent: "agent1", Action: "action1", Condition: "start"},
                    {Agent: "agent2", Action: "action2", Condition: "agent1-completed"},
                },
            },
        },
    }

    engine := newEngine(playbook)

    t.Run("initial state", func(t *testing.T) {
        assert.False(t, engine.isPlaybookComplete())

        steps := engine.getNextExecutableSteps()
        assert.Len(t, steps, 1)
        assert.Equal(t, 0, steps[0].BranchIndex)
        assert.Equal(t, 0, steps[0].StepIndex)
    })

    t.Run("execute first step", func(t *testing.T) {
        err := engine.executeStep(0, 0)
        assert.NoError(t, err)

        step := &playbook.Branches[0].Steps[0]
        assert.Equal(t, StatusCompleted, step.Status)

        // Second step should now be executable
        steps := engine.getNextExecutableSteps()
        assert.Len(t, steps, 1)
        assert.Equal(t, 1, steps[0].StepIndex)
    })

    t.Run("complete execution", func(t *testing.T) {
        err := engine.executeStep(0, 1)
        assert.NoError(t, err)

        assert.True(t, engine.isPlaybookComplete())

        steps := engine.getNextExecutableSteps()
        assert.Len(t, steps, 0)
    })
}
```

### Component Tests

Test UI components in isolation with mocked dependencies.

```go
// ui_test.go
package main

import (
    "testing"
    "time"

    tea "github.com/charmbracelet/bubbletea"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestPlaybookUI(t *testing.T) {
    playbook := &Playbook{
        Name: "Test UI Playbook",
        Branches: []Branch{
            {
                Name: "test-branch",
                Steps: []Step{
                    {Agent: "test-agent", Action: "test", Condition: "start"},
                },
            },
        },
    }

    ui := newPlaybookUI(playbook)

    t.Run("initialization", func(t *testing.T) {
        assert.Equal(t, "overview", ui.currentView)
        assert.False(t, ui.running)
        assert.False(t, ui.paused)
        assert.Len(t, ui.logs, 0)
    })

    t.Run("view switching", func(t *testing.T) {
        ui.switchView()
        assert.Equal(t, "branches", ui.currentView)

        ui.switchView()
        assert.Equal(t, "logs", ui.currentView)

        ui.switchView()
        assert.Equal(t, "progress", ui.currentView)

        ui.switchView()
        assert.Equal(t, "overview", ui.currentView)
    })

    t.Run("keyboard handling", func(t *testing.T) {
        // Test quit key
        model, cmd := ui.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
        assert.Equal(t, ui, model)
        assert.Equal(t, tea.Quit, cmd())

        // Test help toggle
        ui.help.ShowAll = false
        model, _ = ui.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
        updatedUI := model.(*PlaybookUI)
        assert.True(t, updatedUI.help.ShowAll)
    })

    t.Run("window resize", func(t *testing.T) {
        msg := tea.WindowSizeMsg{Width: 120, Height: 40}
        model, _ := ui.Update(msg)
        updatedUI := model.(*PlaybookUI)

        assert.Equal(t, 120, updatedUI.width)
        assert.Equal(t, 40, updatedUI.height)
    })
}

func TestPlaybookUIExecution(t *testing.T) {
    playbook := &Playbook{
        Name: "Execution Test",
        Branches: []Branch{
            {
                Name: "quick-test",
                Steps: []Step{
                    {Agent: "fast-agent", Action: "quick-action", Condition: "start"},
                },
            },
        },
    }

    ui := newPlaybookUI(playbook)

    t.Run("start execution", func(t *testing.T) {
        // Simulate execute key press
        keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}}
        model, cmd := ui.Update(keyMsg)
        updatedUI := model.(*PlaybookUI)

        assert.True(t, updatedUI.running)
        assert.False(t, updatedUI.paused)
        assert.NotNil(t, cmd)
    })

    t.Run("pause execution", func(t *testing.T) {
        ui.running = true
        ui.paused = false

        keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'p'}}
        model, _ := ui.Update(keyMsg)
        updatedUI := model.(*PlaybookUI)

        assert.True(t, updatedUI.paused)
    })

    t.Run("reset execution", func(t *testing.T) {
        ui.running = true
        ui.logs = []string{"test log"}

        keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}}
        model, _ := ui.Update(keyMsg)
        updatedUI := model.(*PlaybookUI)

        assert.False(t, updatedUI.running)
        assert.Len(t, updatedUI.logs, 1) // Should include reset log
    })
}
```

---

## Bubble Tea Testing

### Model Testing Pattern

Test Bubble Tea models with a structured approach:

```go
// bubble_tea_test.go
package main

import (
    "testing"
    tea "github.com/charmbracelet/bubbletea"
)

// TestModel is a helper for testing Bubble Tea models
type TestModel struct {
    model tea.Model
    t     *testing.T
}

func NewTestModel(t *testing.T, model tea.Model) *TestModel {
    return &TestModel{model: model, t: t}
}

func (tm *TestModel) Init() *TestModel {
    cmd := tm.model.Init()
    if cmd != nil {
        // Execute initialization command
        msg := cmd()
        tm.model, _ = tm.model.Update(msg)
    }
    return tm
}

func (tm *TestModel) SendKey(key string) *TestModel {
    var msg tea.Msg

    switch key {
    case "enter":
        msg = tea.KeyMsg{Type: tea.KeyEnter}
    case "esc":
        msg = tea.KeyMsg{Type: tea.KeyEsc}
    case "ctrl+c":
        msg = tea.KeyMsg{Type: tea.KeyCtrlC}
    case "tab":
        msg = tea.KeyMsg{Type: tea.KeyTab}
    case "up":
        msg = tea.KeyMsg{Type: tea.KeyUp}
    case "down":
        msg = tea.KeyMsg{Type: tea.KeyDown}
    default:
        msg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(key)}
    }

    tm.model, _ = tm.model.Update(msg)
    return tm
}

func (tm *TestModel) SendMessage(msg tea.Msg) *TestModel {
    tm.model, _ = tm.model.Update(msg)
    return tm
}

func (tm *TestModel) Resize(width, height int) *TestModel {
    msg := tea.WindowSizeMsg{Width: width, Height: height}
    tm.model, _ = tm.model.Update(msg)
    return tm
}

func (tm *TestModel) AssertContains(text string) *TestModel {
    view := tm.model.View()
    if !strings.Contains(view, text) {
        tm.t.Errorf("Expected view to contain %q, but got:\n%s", text, view)
    }
    return tm
}

func (tm *TestModel) AssertNotContains(text string) *TestModel {
    view := tm.model.View()
    if strings.Contains(view, text) {
        tm.t.Errorf("Expected view NOT to contain %q, but got:\n%s", text, view)
    }
    return tm
}

func (tm *TestModel) GetModel() tea.Model {
    return tm.model
}

// Example usage
func TestPlaybookUIWithHelper(t *testing.T) {
    playbook := &Playbook{
        Name: "Test Playbook",
        Branches: []Branch{{Name: "test", Steps: []Step{}}},
    }

    ui := newPlaybookUI(playbook)

    NewTestModel(t, ui).
        Init().
        Resize(100, 30).
        AssertContains("Test Playbook").
        SendKey("tab").
        AssertContains("branches").
        SendKey("?").
        AssertContains("help")
}
```

### Message Testing

Test specific message handling:

```go
func TestMessageHandling(t *testing.T) {
    ui := newPlaybookUI(&Playbook{Name: "Test", Branches: []Branch{}})

    tests := []struct {
        name        string
        message     tea.Msg
        expectView  string
        expectState func(*PlaybookUI) bool
    }{
        {
            name:       "tick message",
            message:    tickMsg(time.Now()),
            expectView: "overview",
            expectState: func(ui *PlaybookUI) bool {
                return !ui.running
            },
        },
        {
            name:       "step complete",
            message:    stepCompleteMsg{branchIdx: 0, stepIdx: 0},
            expectView: "overview",
            expectState: func(ui *PlaybookUI) bool {
                return len(ui.logs) > 0
            },
        },
        {
            name:       "window resize",
            message:    tea.WindowSizeMsg{Width: 80, Height: 24},
            expectView: "overview",
            expectState: func(ui *PlaybookUI) bool {
                return ui.width == 80 && ui.height == 24
            },
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            model, _ := ui.Update(tt.message)
            updatedUI := model.(*PlaybookUI)

            assert.Equal(t, tt.expectView, updatedUI.currentView)
            if tt.expectState != nil {
                assert.True(t, tt.expectState(updatedUI))
            }
        })
    }
}
```

---

## UI Component Testing

### Style Testing

Verify that styles are applied correctly:

```go
func TestStyles(t *testing.T) {
    t.Run("color values", func(t *testing.T) {
        assert.Equal(t, "#00ffff", string(primaryColor))
        assert.Equal(t, "#ff0080", string(secondaryColor))
        assert.Equal(t, "#00ff00", string(successColor))
    })

    t.Run("style application", func(t *testing.T) {
        text := "Test Text"
        styled := titleStyle.Render(text)

        // Verify style contains expected elements
        assert.Contains(t, styled, text)
        // In a real test, you'd verify ANSI codes
    })
}

func TestStatusRendering(t *testing.T) {
    tests := []struct {
        status   StepStatus
        expected string
    }{
        {StatusPending, "○"},
        {StatusRunning, "⟳"},
        {StatusCompleted, "✓"},
        {StatusFailed, "✗"},
    }

    ui := newPlaybookUI(&Playbook{})

    for _, tt := range tests {
        t.Run(tt.status.String(), func(t *testing.T) {
            icon := ui.getStatusIcon(tt.status)
            assert.Equal(t, tt.expected, icon)
        })
    }
}
```

### Layout Testing

Test component layouts and positioning:

```go
func TestLayoutRendering(t *testing.T) {
    playbook := &Playbook{
        Name: "Layout Test",
        Branches: []Branch{
            {
                Name: "branch1",
                Steps: []Step{
                    {Agent: "agent1", Action: "action1", Condition: "start"},
                    {Agent: "agent2", Action: "action2", Condition: "agent1-completed"},
                },
            },
        },
    }

    ui := newPlaybookUI(playbook)
    ui.width = 100
    ui.height = 40

    view := ui.View()

    t.Run("contains required elements", func(t *testing.T) {
        assert.Contains(t, view, "Layout Test")
        assert.Contains(t, view, "branch1")
        assert.Contains(t, view, "agent1")
        assert.Contains(t, view, "agent2")
    })

    t.Run("view switching", func(t *testing.T) {
        ui.currentView = "branches"
        view := ui.View()
        assert.Contains(t, view, "🌿 Branches")

        ui.currentView = "logs"
        view = ui.View()
        assert.Contains(t, view, "📋 Activity Logs")

        ui.currentView = "progress"
        view = ui.View()
        assert.Contains(t, view, "📈 Progress Details")
    })
}
```

---

## Integration Testing

### Full Workflow Testing

Test complete user workflows:

```go
// integration_test.go
package main

import (
    "os"
    "path/filepath"
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestPlaybookExecutionWorkflow(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test in short mode")
    }

    // Create temporary directory
    tmpDir := t.TempDir()
    playbookPath := filepath.Join(tmpDir, "test-playbook.yaml")

    // Create test playbook
    playbookContent := `
name: "Integration Test Playbook"
version: "1.0.0"
branches:
  - name: "setup"
    steps:
      - agent: "setup-agent"
        action: "initialize environment"
        condition: "start"
      - agent: "config-agent"
        action: "configure settings"
        condition: "setup-completed"
  - name: "execution"
    steps:
      - agent: "work-agent"
        action: "do work"
        condition: "config-completed"
`

    err := os.WriteFile(playbookPath, []byte(playbookContent), 0644)
    require.NoError(t, err)

    t.Run("load and validate playbook", func(t *testing.T) {
        playbook, err := loadPlaybook(playbookPath)
        require.NoError(t, err)

        assert.Equal(t, "Integration Test Playbook", playbook.Name)
        assert.Len(t, playbook.Branches, 2)
    })

    t.Run("execute playbook workflow", func(t *testing.T) {
        playbook, err := loadPlaybook(playbookPath)
        require.NoError(t, err)

        engine := newEngine(playbook)

        // Verify initial state
        assert.False(t, engine.isPlaybookComplete())

        steps := engine.getNextExecutableSteps()
        assert.Len(t, steps, 1)
        assert.Equal(t, "setup-agent", playbook.Branches[steps[0].BranchIndex].Steps[steps[0].StepIndex].Agent)

        // Execute first step
        err = engine.executeStep(steps[0].BranchIndex, steps[0].StepIndex)
        assert.NoError(t, err)

        // Check progression
        steps = engine.getNextExecutableSteps()
        assert.Len(t, steps, 1)
        assert.Equal(t, "config-agent", playbook.Branches[steps[0].BranchIndex].Steps[steps[0].StepIndex].Agent)

        // Complete workflow
        for !engine.isPlaybookComplete() {
            steps = engine.getNextExecutableSteps()
            if len(steps) == 0 {
                break
            }

            for _, step := range steps {
                err = engine.executeStep(step.BranchIndex, step.StepIndex)
                assert.NoError(t, err)
            }
        }

        assert.True(t, engine.isPlaybookComplete())
    })
}

func TestConfigurationWorkflow(t *testing.T) {
    tmpDir := t.TempDir()
    configPath := filepath.Join(tmpDir, "config.yaml")

    t.Run("create default config", func(t *testing.T) {
        config := CreateDefaultConfig()
        err := SaveConfig(config, configPath)
        require.NoError(t, err)

        // Verify file exists
        _, err = os.Stat(configPath)
        assert.NoError(t, err)
    })

    t.Run("load and modify config", func(t *testing.T) {
        config, err := LoadConfig(configPath)
        require.NoError(t, err)

        // Modify config
        config.DevGen.DefaultOutputDir = "/custom/output"
        config.Logging.Level = "debug"

        err = SaveConfig(config, configPath)
        require.NoError(t, err)

        // Reload and verify changes
        reloaded, err := LoadConfig(configPath)
        require.NoError(t, err)

        assert.Equal(t, "/custom/output", reloaded.DevGen.DefaultOutputDir)
        assert.Equal(t, "debug", reloaded.Logging.Level)
    })
}
```

### CLI Command Testing

Test CLI commands end-to-end:

```go
func TestCLICommands(t *testing.T) {
    // Helper to run CLI commands
    runCLI := func(args ...string) (string, error) {
        // In a real implementation, you'd capture stdout/stderr
        // and run the CLI with the given arguments
        return "", nil
    }

    t.Run("help command", func(t *testing.T) {
        output, err := runCLI("--help")
        assert.NoError(t, err)
        assert.Contains(t, output, "DevGen")
        assert.Contains(t, output, "Usage:")
    })

    t.Run("version command", func(t *testing.T) {
        output, err := runCLI("--version")
        assert.NoError(t, err)
        assert.Contains(t, output, "1.0.0")
    })

    t.Run("config init", func(t *testing.T) {
        output, err := runCLI("config", "init")
        assert.NoError(t, err)
        assert.Contains(t, output, "Configuration initialized")
    })
}
```

---

## Performance Testing

### Benchmark Tests

Test performance of critical operations:

```go
// benchmark_test.go
package main

import (
    "testing"
    "time"
)

func BenchmarkPlaybookLoading(b *testing.B) {
    // Create a large playbook for testing
    playbook := createLargePlaybook(100, 10) // 100 branches, 10 steps each

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        engine := newEngine(playbook)
        _ = engine.getNextExecutableSteps()
    }
}

func BenchmarkUIRendering(b *testing.B) {
    playbook := createLargePlaybook(50, 5)
    ui := newPlaybookUI(playbook)
    ui.width = 120
    ui.height = 40

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _ = ui.View()
    }
}

func BenchmarkStepExecution(b *testing.B) {
    playbook := &Playbook{
        Name: "Benchmark Test",
        Branches: []Branch{
            {
                Name: "bench",
                Steps: []Step{
                    {Agent: "fast-agent", Action: "quick-action", Condition: "start"},
                },
            },
        },
    }

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        engine := newEngine(playbook)
        _ = engine.executeStep(0, 0)
    }
}

func createLargePlaybook(numBranches, stepsPerBranch int) *Playbook {
    playbook := &Playbook{
        Name:     "Large Test Playbook",
        Branches: make([]Branch, numBranches),
    }

    for i := 0; i < numBranches; i++ {
        branch := Branch{
            Name:  fmt.Sprintf("branch-%d", i),
            Steps: make([]Step, stepsPerBranch),
        }

        for j := 0; j < stepsPerBranch; j++ {
            condition := "start"
            if j > 0 {
                condition = fmt.Sprintf("step-%d-%d-completed", i, j-1)
            }

            branch.Steps[j] = Step{
                Agent:     fmt.Sprintf("agent-%d-%d", i, j),
                Action:    fmt.Sprintf("action-%d-%d", i, j),
                Condition: condition,
            }
        }

        playbook.Branches[i] = branch
    }

    return playbook
}
```

### Memory Usage Testing

Monitor memory usage during execution:

```go
func TestMemoryUsage(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping memory test in short mode")
    }

    var m1, m2 runtime.MemStats

    // Measure initial memory
    runtime.GC()
    runtime.ReadMemStats(&m1)

    // Create and run large playbook
    playbook := createLargePlaybook(1000, 20)
    ui := newPlaybookUI(playbook)

    // Simulate UI operations
    for i := 0; i < 1000; i++ {
        ui.addLog(fmt.Sprintf("Log entry %d", i))
        _ = ui.View()
        ui.switchView()
    }

    // Measure final memory
    runtime.GC()
    runtime.ReadMemStats(&m2)

    // Calculate memory usage
    memUsed := m2.TotalAlloc - m1.TotalAlloc
    t.Logf("Memory used: %d KB", memUsed/1024)

    // Assert reasonable memory usage (adjust threshold as needed)
    assert.Less(t, memUsed, uint64(100*1024*1024)) // 100MB threshold
}
```

---

## Accessibility Testing

### Color Contrast Testing

Ensure adequate color contrast for accessibility:

```go
func TestColorContrast(t *testing.T) {
    tests := []struct {
        name       string
        background string
        foreground string
        minRatio   float64
    }{
        {"primary on background", "#000000", "#00ffff", 4.5},
        {"success on background", "#000000", "#00ff00", 4.5},
        {"error on background", "#000000", "#ff0080", 4.5},
        {"text on background", "#000000", "#ffffff", 7.0},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            ratio := calculateContrastRatio(tt.background, tt.foreground)
            assert.GreaterOrEqual(t, ratio, tt.minRatio,
                "Insufficient contrast ratio for %s", tt.name)
        })
    }
}

func calculateContrastRatio(bg, fg string) float64 {
    // Simplified contrast calculation
    // In practice, you'd use a proper color library
    return 7.0 // Placeholder
}
```

### Keyboard Navigation Testing

Verify keyboard accessibility:

```go
func TestKeyboardNavigation(t *testing.T) {
    ui := newPlaybookUI(&Playbook{
        Name: "Keyboard Test",
        Branches: []Branch{
            {Name: "test", Steps: []Step{{Agent: "test", Action: "test", Condition: "start"}}},
        },
    })

    t.Run("tab navigation", func(t *testing.T) {
        initialView := ui.currentView

        // Tab should cycle through views
        ui.switchView()
        assert.NotEqual(t, initialView, ui.currentView)
    })

    t.Run("help accessibility", func(t *testing.T) {
        // Help should be toggleable
        initialState := ui.help.ShowAll

        model, _ := ui.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
        updatedUI := model.(*PlaybookUI)

        assert.NotEqual(t, initialState, updatedUI.help.ShowAll)
    })

    t.Run("quit commands", func(t *testing.T) {
        // Multiple quit methods should work
        quitKeys := []tea.KeyMsg{
            {Type: tea.KeyCtrlC},
            {Type: tea.KeyRunes, Runes: []rune{'q'}},
        }

        for _, key := range quitKeys {
            _, cmd := ui.Update(key)
            assert.Equal(t, tea.Quit, cmd())
        }
    })
}
```

---

## CI/CD Integration

### GitHub Actions Workflow

```yaml
# .github/workflows/test.yml
name: Test DevGen CLI

on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main]

jobs:
  test:
    strategy:
      matrix:
        os: [ubuntu-latest, macos-latest, windows-latest]
        go-version: [1.21.x]

    runs-on: ${{ matrix.os }}

    steps:
      - uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: ${{ matrix.go-version }}

      - name: Cache dependencies
        uses: actions/cache@v3
        with:
          path: ~/go/pkg/mod
          key: ${{
```
