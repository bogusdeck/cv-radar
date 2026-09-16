package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/bogusdeck/ats-scanner/internal/ats"
	"github.com/bogusdeck/ats-scanner/internal/models"
	"github.com/bogusdeck/ats-scanner/internal/parser"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	state      int // 0:welcome, 1:input JD, 2:select ATS, 3:scanning, 4:result, 5:optimizing, 6:done
	jdText     string
	cvText     string
	platform   string
	platforms  []string
	cursor     int
	result     *models.ATSResult
	err        error
	width      int
	height     int
	spinner    spinner.Model
}

type optimizeResult struct {
	err error
}

func initialModel() model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return model{
		state:     0,
		platforms: []string{"Workday", "Taleo", "Greenhouse", "iCIMS"},
		cursor:    0,
		spinner:   s,
	}
}

func (m model) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}

		switch m.state {
		case 0:
			if msg.String() == "enter" {
				m.state = 1
			}
		case 1:
			if msg.String() == "enter" {
				b, _ := ioutil.ReadFile("jd_test.txt")
				m.jdText = string(b)
				if m.jdText == "" {
					m.jdText = "Backend Engineer\nRequirements: Kubernetes, Database, AWS, Python"
				}
				m.state = 2
			}
		case 2:
			switch msg.String() {
			case "up", "k":
				if m.cursor > 0 {
					m.cursor--
				}
			case "down", "j":
				if m.cursor < len(m.platforms)-1 {
					m.cursor++
				}
			case "enter":
				m.platform = m.platforms[m.cursor]
				m.state = 3
				return m, func() tea.Msg {
					b, err := ioutil.ReadFile("cv_test.txt")
					if err != nil {
						return err
					}
					m.cvText = string(b)
					parsed := parser.ParseText(m.cvText)
					engine := ats.Get(m.platform)
					var res models.ATSResult
					if engine != nil {
						res = engine.Analyze(parsed, m.jdText)
					}
					return res
				}
			}
		case 4:
			if msg.String() == "o" || msg.String() == "O" {
				m.state = 5
				return m, tea.Batch(m.spinner.Tick, runOptimizeCmd(m.cvText, m.jdText, m.platform))
			}
		case 6:
			if msg.String() == "enter" {
				return m, tea.Quit
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case models.ATSResult:
		m.result = &msg
		m.state = 4
	case optimizeResult:
		m.err = msg.err
		m.state = 6
	case error:
		m.err = msg
		m.state = 4
	
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	return m, nil
}

func detectAiTool() string {
	tools := []string{"opencode", "claude", "agy", "codex"}
	for _, t := range tools {
		if _, err := exec.LookPath(t); err == nil {
			return t
		}
	}
	return "claude"
}

func runOptimizeCmd(cv string, jd string, platform string) tea.Cmd {
	return func() tea.Msg {
		aiTool := detectAiTool()
		cmd := exec.Command("./optimize.sh", "cv_test.txt", "jd_test.txt", platform, aiTool)
		err := cmd.Run()
		return optimizeResult{err: err}
	}
}

var (
	titleStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FF75B5")).MarginBottom(1)
	boxStyle   = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#874BFD")).Padding(1, 2)
)

func (m model) View() string {
	if m.width == 0 {
		return "Initializing..."
	}

	var content string

	switch m.state {
	case 0:
		content = titleStyle.Render("CV Optimizer TUI 🚀") + "\n\nPress Enter to load Job Description..."
	case 1:
		content = "Loading Job Description..."
	case 2:
		content = titleStyle.Render("Select Target ATS Platform:") + "\n\n"
		for i, p := range m.platforms {
			cursor := " "
			if m.cursor == i {
				cursor = ">"
			}
			content += fmt.Sprintf("%s %s\n", cursor, p)
		}
		content += "\n(Press Enter to select)"
	case 3:
		content = "Scanning CV against Job Description...\nEvaluating layout, keywords, and structural integrity..."
	case 4:
		if m.err != nil {
			content = fmt.Sprintf("Error: %v", m.err)
		} else if m.result != nil {
			scoreColor := "#2EA043"
			if m.result.Score < 50 {
				scoreColor = "#F85149"
			} else if m.result.Score < 80 {
				scoreColor = "#D29922"
			}
			scoreStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(scoreColor))
			
			content = titleStyle.Render(fmt.Sprintf("Scan Results: %s", m.result.Platform)) + "\n"
			content += fmt.Sprintf("Total Score: %s\n", scoreStyle.Render(fmt.Sprintf("%.1f / 100", m.result.Score)))
			content += fmt.Sprintf("Grade: %s\n\n", scoreStyle.Render(m.result.Grade))
			
			content += lipgloss.NewStyle().Bold(true).Render("Missing Keywords:") + "\n"
			if len(m.result.MissingKeywords) > 0 {
				content += strings.Join(m.result.MissingKeywords, ", ") + "\n\n"
			} else {
				content += "None! Perfect match.\n\n"
			}
			
			content += lipgloss.NewStyle().Bold(true).Render("Actionable Recommendations:") + "\n"
			for _, rec := range m.result.Recommendations {
				content += "• " + rec + "\n"
			}
			
			content += "\nPress [O] to Optimize with AI (Headless)"
			content += "\nPress [Q] to Quit"
		}
	case 5:
		content = fmt.Sprintf("\n\n   %s Instructing LLM to rewrite CV headlessly...\n\n   Please wait, generating optimized LaTeX...", m.spinner.View())
	case 6:
		if m.err != nil {
			content = fmt.Sprintf("Optimization failed: %v", m.err)
		} else {
			content = titleStyle.Render("Success! 🎉") + "\n\nThe LLM has rewritten your CV and saved it to 'optimized_cv.md'!\n\nPress Enter to exit."
		}
	}

	return lipgloss.Place(m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		boxStyle.Render(content),
	)
}

func main() {
	// If launched with open flag or from a non-interactive tool call (AI agent like OpenCode / Claude Code / Antigravity)
	isOpenArg := len(os.Args) > 1 && (os.Args[1] == "open" || os.Args[1] == "-open" || os.Args[1] == "--open")
	if isOpenArg {
		if runtime.GOOS == "darwin" {
			if _, err := exec.LookPath("kitty"); err == nil {
				exec.Command("kitty", "-e", "cv-tui").Start()
				fmt.Println("🚀 Launched cv-tui in Kitty terminal window!")
				return
			}
			exec.Command("osascript", "-e", `tell application "Terminal" to activate`, "-e", `tell application "Terminal" to do script "cv-tui"`).Run()
			fmt.Println("🚀 Launched cv-tui in active Terminal window!")
			return
		}
	}

	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
