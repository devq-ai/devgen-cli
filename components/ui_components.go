package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/filepicker"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

// PlaybookLister shows available playbooks with interactive selection
type PlaybookLister struct {
	list   list.Model
	logger *log.Logger
	width  int
	height int
}

type playbookItem struct {
	name        string
	path        string
	description string
	lastModified time.Time
}

func (i playbookItem) FilterValue() string { return i.name }
func (i playbookItem) Title() string       { return i.name }
func (i playbookItem) Description() string {
	return fmt.Sprintf("%s • Modified: %s", i.description, i.lastModified.Format("Jan 2, 2006"))
}

func newPlaybookLister() *PlaybookLister {
	items := []list.Item{}

	// Scan for playbook files
	playbookFiles, err := filepath.Glob("*.yaml")
	if err == nil {
		for _, file := range playbookFiles {
			if strings.Contains(file, "playbook") || strings.Contains(file, "workflow") {
				info, err := os.Stat(file)
				if err == nil {
					item := playbookItem{
						name:         strings.TrimSuffix(file, ".yaml"),
						path:         file,
						description:  "Development playbook",
						lastModified: info.ModTime(),
					}
					items = append(items, item)
				}
			}
		}
	}

	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Available Playbooks"
	l.SetFilteringEnabled(true)
	l.SetShowStatusBar(true)
	l.SetShowPagination(true)

	return &PlaybookLister{
		list:   l,
		logger: log.New(nil),
	}
}

func (pl *PlaybookLister) Init() tea.Cmd {
	return nil
}

func (pl *PlaybookLister) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		pl.width = msg.Width
		pl.height = msg.Height
		pl.list.SetWidth(msg.Width)
		pl.list.SetHeight(msg.Height - 4)
		return pl, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return pl, tea.Quit
		case "enter":
			selected := pl.list.SelectedItem()
			if item, ok := selected.(playbookItem); ok {
				pl.logger.Info("Selected playbook", "name", item.name, "path", item.path)
				// Launch playbook execution
				return pl, tea.Quit
			}
		}
	}

	var cmd tea.Cmd
	pl.list, cmd = pl.list.Update(msg)
	return pl, cmd
}

func (pl *PlaybookLister) View() string {
	if pl.width == 0 {
		return "Loading..."
	}

	content := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(primaryColor).
		Padding(1).
		Render(pl.list.View())

	help := helpStyle.Render("↑/↓: navigate • enter: select • /: filter • q: quit")

	return lipgloss.JoinVertical(lipgloss.Left, content, help)
}

func (pl *PlaybookLister) run() error {
	p := tea.NewProgram(pl, tea.WithAltScreen())
	_, err := p.Run()
	return err
}

// TemplateLister shows available project templates
type TemplateLister struct {
	table  table.Model
	logger *log.Logger
	width  int
	height int
}

func newTemplateLister() *TemplateLister {
	columns := []table.Column{
		{Title: "Name", Width: 20},
		{Title: "Type", Width: 15},
		{Title: "Description", Width: 40},
		{Title: "Version", Width: 10},
		{Title: "Status", Width: 10},
	}

	rows := []table.Row{
		{"fastapi-basic", "Backend", "Basic FastAPI project with Logfire", "1.0.0", "✓ Available"},
		{"nextjs-app", "Frontend", "Next.js app with Tailwind and Shadcn", "1.2.0", "✓ Available"},
		{"cli-charm", "CLI", "CLI app with Charm components", "1.0.0", "✓ Available"},
		{"fullstack-saas", "Fullstack", "Complete SaaS with auth and billing", "2.1.0", "✓ Available"},
		{"microservice", "Backend", "Microservice with Docker and K8s", "1.5.0", "✓ Available"},
		{"data-pipeline", "Data", "ETL pipeline with Airflow", "1.3.0", "🔄 Installing"},
		{"ml-training", "ML", "Machine learning training pipeline", "1.1.0", "❌ Failed"},
		{"react-native", "Mobile", "React Native with Expo", "2.0.0", "✓ Available"},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(primaryColor).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(textColor).
		Background(primaryColor).
		Bold(false)
	t.SetStyles(s)

	return &TemplateLister{
		table:  t,
		logger: log.New(nil),
	}
}

func (tl *TemplateLister) Init() tea.Cmd {
	return nil
}

func (tl *TemplateLister) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		tl.width = msg.Width
		tl.height = msg.Height
		return tl, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return tl, tea.Quit
		case "enter":
			selected := tl.table.SelectedRow()
			if len(selected) > 0 {
				tl.logger.Info("Selected template", "name", selected[0], "type", selected[1])
				return tl, tea.Quit
			}
		}
	}

	tl.table, cmd = tl.table.Update(msg)
	return tl, cmd
}

func (tl *TemplateLister) View() string {
	content := baseStyle.Render(tl.table.View()) + "\n"
	help := helpStyle.Render("↑/↓: navigate • enter: install • q: quit")
	return content + help
}

func (tl *TemplateLister) run() error {
	p := tea.NewProgram(tl, tea.WithAltScreen())
	_, err := p.Run()
	return err
}

// TemplateInstaller handles template installation with progress
type TemplateInstaller struct {
	templateName string
	spinner      spinner.Model
	installing   bool
	completed    bool
	error        error
	logger       *log.Logger
	steps        []string
	currentStep  int
}

func newTemplateInstaller(templateName string) *TemplateInstaller {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(primaryColor)

	steps := []string{
		"Downloading template...",
		"Extracting files...",
		"Installing dependencies...",
		"Configuring project...",
		"Running setup scripts...",
		"Finalizing installation...",
	}

	return &TemplateInstaller{
		templateName: templateName,
		spinner:      s,
		installing:   false,
		completed:    false,
		logger:       log.New(nil),
		steps:        steps,
		currentStep:  0,
	}
}

func (ti *TemplateInstaller) Init() tea.Cmd {
	return tea.Batch(
		ti.spinner.Tick,
		ti.startInstallation(),
	)
}

func (ti *TemplateInstaller) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			if ti.completed {
				return ti, tea.Quit
			}
		}

	case spinner.TickMsg:
		ti.spinner, cmd = ti.spinner.Update(msg)
		return ti, cmd

	case installationStepMsg:
		ti.currentStep = int(msg)
		if ti.currentStep >= len(ti.steps) {
			ti.installing = false
			ti.completed = true
			ti.logger.Info("Template installation completed", "template", ti.templateName)
		}
		return ti, ti.nextStep()

	case installationErrorMsg:
		ti.installing = false
		ti.error = msg.error
		ti.logger.Error("Template installation failed", "template", ti.templateName, "error", msg.error)
	}

	return ti, cmd
}

func (ti *TemplateInstaller) View() string {
	var content strings.Builder

	title := fmt.Sprintf("Installing Template: %s", ti.templateName)
	content.WriteString(titleStyle.Render(title))
	content.WriteString("\n\n")

	if ti.installing {
		if ti.currentStep < len(ti.steps) {
			currentStepText := fmt.Sprintf("%s %s", ti.spinner.View(), ti.steps[ti.currentStep])
			content.WriteString(statusRunningStyle.Render(currentStepText))
			content.WriteString("\n\n")
		}

		// Show progress
		content.WriteString("Progress:\n")
		for i, step := range ti.steps {
			if i < ti.currentStep {
				content.WriteString(statusCompletedStyle.Render(fmt.Sprintf("✓ %s", step)))
			} else if i == ti.currentStep {
				content.WriteString(statusRunningStyle.Render(fmt.Sprintf("⟳ %s", step)))
			} else {
				content.WriteString(statusPendingStyle.Render(fmt.Sprintf("○ %s", step)))
			}
			content.WriteString("\n")
		}
	} else if ti.completed {
		content.WriteString(statusCompletedStyle.Render("✓ Installation completed successfully!"))
		content.WriteString("\n\n")
		content.WriteString(helpStyle.Render("Press 'q' to exit"))
	} else if ti.error != nil {
		content.WriteString(statusFailedStyle.Render(fmt.Sprintf("✗ Installation failed: %s", ti.error.Error())))
		content.WriteString("\n\n")
		content.WriteString(helpStyle.Render("Press 'q' to exit"))
	}

	return baseStyle.Render(content.String())
}

func (ti *TemplateInstaller) startInstallation() tea.Cmd {
	ti.installing = true
	return ti.nextStep()
}

func (ti *TemplateInstaller) nextStep() tea.Cmd {
	return tea.Tick(time.Millisecond*1500, func(t time.Time) tea.Msg {
		return installationStepMsg(ti.currentStep + 1)
	})
}

func (ti *TemplateInstaller) run() error {
	p := tea.NewProgram(ti, tea.WithAltScreen())
	_, err := p.Run()
	return err
}

// ProjectInitializer handles interactive project setup
type ProjectInitializer struct {
	projectName string
	step        int
	form        *huh.Form
	completed   bool
	logger      *log.Logger
}

func newProjectInitializer(projectName string) *ProjectInitializer {
	return &ProjectInitializer{
		projectName: projectName,
		step:        0,
		completed:   false,
		logger:      log.New(nil),
	}
}

func (pi *ProjectInitializer) Init() tea.Cmd {
	return pi.createForm()
}

func (pi *ProjectInitializer) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return pi, tea.Quit
		}
	}

	return pi, nil
}

func (pi *ProjectInitializer) View() string {
	if pi.completed {
		return statusCompletedStyle.Render("✓ Project initialized successfully!")
	}

	title := "Project Initialization"
	if pi.projectName != "" {
		title = fmt.Sprintf("Initializing: %s", pi.projectName)
	}

	content := titleStyle.Render(title) + "\n\n"
	content += "Setting up your new project with DevQ.ai standards...\n\n"

	return baseStyle.Render(content)
}

func (pi *ProjectInitializer) createForm() tea.Cmd {
	var name, description, template string
	var includeCI, includeDocker bool

	pi.form = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Project Name").
				Value(&name).
				Validate(func(s string) error {
					if len(s) == 0 {
						return fmt.Errorf("project name cannot be empty")
					}
					return nil
				}),

			huh.NewText().
				Title("Description").
				Value(&description).
				Lines(3),

			huh.NewSelect[string]().
				Title("Project Template").
				Options(
					huh.NewOption("FastAPI Backend", "fastapi"),
					huh.NewOption("Next.js Frontend", "nextjs"),
					huh.NewOption("CLI Application", "cli"),
					huh.NewOption("Fullstack SaaS", "fullstack"),
				).
				Value(&template),
		),

		huh.NewGroup(
			huh.NewConfirm().
				Title("Include CI/CD pipeline?").
				Value(&includeCI),

			huh.NewConfirm().
				Title("Include Docker configuration?").
				Value(&includeDocker),
		),
	)

	return nil
}

func (pi *ProjectInitializer) run() error {
	if pi.form != nil {
		err := pi.form.Run()
		if err != nil {
			return fmt.Errorf("form error: %w", err)
		}
	}

	pi.completed = true
	pi.logger.Info("Project initialization completed", "name", pi.projectName)
	return nil
}

// ProjectStatusViewer shows current project status
type ProjectStatusViewer struct {
	table  table.Model
	logger *log.Logger
}

func newProjectStatusViewer() *ProjectStatusViewer {
	columns := []table.Column{
		{Title: "Component", Width: 20},
		{Title: "Status", Width: 15},
		{Title: "Version", Width: 12},
		{Title: "Last Updated", Width: 15},
		{Title: "Health", Width: 10},
	}

	rows := []table.Row{
		{"FastAPI", "✓ Running", "0.104.1", "2 hours ago", "🟢 Healthy"},
		{"Database", "✓ Connected", "15.4", "1 day ago", "🟢 Healthy"},
		{"Redis", "✓ Active", "7.2.0", "3 hours ago", "🟢 Healthy"},
		{"Logfire", "✓ Monitoring", "0.28.0", "30 min ago", "🟢 Healthy"},
		{"Tests", "✓ Passing", "90% cov", "15 min ago", "🟢 Healthy"},
		{"CI/CD", "🔄 Building", "v1.2.3", "5 min ago", "🟡 Building"},
		{"Docker", "✓ Ready", "24.0.6", "1 hour ago", "🟢 Healthy"},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(primaryColor).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(textColor).
		Background(primaryColor).
		Bold(false)
	t.SetStyles(s)

	return &ProjectStatusViewer{
		table:  t,
		logger: log.New(nil),
	}
}

func (psv *ProjectStatusViewer) Init() tea.Cmd {
	return nil
}

func (psv *ProjectStatusViewer) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return psv, tea.Quit
		case "r":
			// Refresh status
			psv.logger.Info("Refreshing project status...")
			return psv, nil
		}
	}

	psv.table, cmd = psv.table.Update(msg)
	return psv, cmd
}

func (psv *ProjectStatusViewer) View() string {
	title := titleStyle.Render("Project Status Dashboard")
	content := title + "\n\n" + baseStyle.Render(psv.table.View()) + "\n"
	help := helpStyle.Render("r: refresh • q: quit")
	return content + help
}

func (psv *ProjectStatusViewer) run() error {
	p := tea.NewProgram(psv, tea.WithAltScreen())
	_, err := p.Run()
	return err
}

// ArtifactGenerator handles code and file generation
type ArtifactGenerator struct {
	artifactType string
	filepicker   filepicker.Model
	textinput    textinput.Model
	step         int
	generated    bool
	logger       *log.Logger
}

func newArtifactGenerator(artifactType string) *ArtifactGenerator {
	fp := filepicker.New()
	fp.AllowedTypes = []string{".go", ".py", ".js", ".ts", ".yaml", ".json"}
	fp.CurrentDirectory, _ = os.Getwd()

	ti := textinput.New()
	ti.Placeholder = "Enter artifact name..."
	ti.Focus()

	return &ArtifactGenerator{
		artifactType: artifactType,
		filepicker:   fp,
		textinput:    ti,
		step:         0,
		generated:    false,
		logger:       log.New(nil),
	}
}

func (ag *ArtifactGenerator) Init() tea.Cmd {
	return tea.Batch(
		textinput.Blink,
		ag.filepicker.Init(),
	)
}

func (ag *ArtifactGenerator) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return ag, tea.Quit
		case "enter":
			if ag.step == 0 && ag.textinput.Value() != "" {
				ag.step = 1
				return ag, nil
			} else if ag.step == 1 {
				ag.generated = true
				ag.logger.Info("Generated artifact",
					"type", ag.artifactType,
					"name", ag.textinput.Value(),
					"path", ag.filepicker.CurrentDirectory)
				return ag, nil
			}
		}
	}

	if ag.step == 0 {
		ag.textinput, cmd = ag.textinput.Update(msg)
	} else if ag.step == 1 {
		ag.filepicker, cmd = ag.filepicker.Update(msg)
	}

	return ag, cmd
}

func (ag *ArtifactGenerator) View() string {
	var content strings.Builder

	title := fmt.Sprintf("Generate %s Artifact", strings.Title(ag.artifactType))
	content.WriteString(titleStyle.Render(title))
	content.WriteString("\n\n")

	if ag.generated {
		content.WriteString(statusCompletedStyle.Render("✓ Artifact generated successfully!"))
		content.WriteString("\n\n")
		content.WriteString(helpStyle.Render("Press 'q' to exit"))
	} else if ag.step == 0 {
		content.WriteString("Artifact Name:\n")
		content.WriteString(ag.textinput.View())
		content.WriteString("\n\n")
		content.WriteString(helpStyle.Render("Enter name and press enter to continue"))
	} else if ag.step == 1 {
		content.WriteString(fmt.Sprintf("Name: %s\n", ag.textinput.Value()))
		content.WriteString("Choose output directory:\n\n")
		content.WriteString(ag.filepicker.View())
		content.WriteString("\n\n")
		content.WriteString(helpStyle.Render("Navigate and press enter to select directory"))
	}

	return baseStyle.Render(content.String())
}

func (ag *ArtifactGenerator) run() error {
	p := tea.NewProgram(ag, tea.WithAltScreen())
	_, err := p.Run()
	return err
}

// Message types for components
type installationStepMsg int
type installationErrorMsg struct {
	error error
}

// Additional supporting components would include:
// - DevServer for development server management
// - ConfigEditor for interactive configuration
// - ServerStatusViewer for monitoring
// - TemplateCreator for creating new templates
// - LogViewer for real-time log streaming
