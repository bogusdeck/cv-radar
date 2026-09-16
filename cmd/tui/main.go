package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/bogusdeck/ats-scanner/internal/ats"
	"github.com/bogusdeck/ats-scanner/internal/models"
	"github.com/bogusdeck/ats-scanner/internal/parser"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type editorFinishedMsg struct {
	fieldIndex int
	tempPath   string
	err        error
}

type model struct {
	state        int // 0:welcome/input form, 1:select ATS, 2:scanning, 3:result, 4:optimizing, 5:done
	focusIndex   int // 0: JD input, 1: CV input, 2: AI Engine selector, 3: Let's Go button
	inputs       []textinput.Model
	jdText       string
	cvText       string
	platform     string
	platforms    []string
	cursor       int
	result       *models.ATSResult
	err          error
	width        int
	height       int
	spinner      spinner.Model
	availableAIs []string
	aiCursor     int
	selectedAI   string
}

type optimizeResult struct {
	err error
}

func initialModel() model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF75B5"))

	inputs := make([]textinput.Model, 2)
	inputs[0] = textinput.New()
	inputs[0].Placeholder = "jd.txt (or paste text, or press Ctrl+G for Vim)"
	inputs[0].Focus()
	inputs[0].CharLimit = 500
	inputs[0].Width = 54

	inputs[1] = textinput.New()
	inputs[1].Placeholder = "cv.tex (or file path, or press Ctrl+G for Vim)"
	inputs[1].CharLimit = 500
	inputs[1].Width = 54

	autoAI := detectAiTool()
	availableAIs := []string{"opencode", "claude", "agy", "codex"}
	aiCursor := 0
	for i, tool := range availableAIs {
		if tool == autoAI {
			aiCursor = i
			break
		}
	}

	return model{
		state:        0,
		focusIndex:   0,
		inputs:       inputs,
		platforms:    []string{"Workday", "Taleo", "Greenhouse", "iCIMS"},
		cursor:       0,
		spinner:      s,
		availableAIs: availableAIs,
		aiCursor:     aiCursor,
		selectedAI:   availableAIs[aiCursor],
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func openExternalEditor(fieldIndex int, initialContent string) tea.Cmd {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = os.Getenv("VISUAL")
	}
	if editor == "" {
		if _, err := exec.LookPath("nvim"); err == nil {
			editor = "nvim"
		} else if _, err := exec.LookPath("vim"); err == nil {
			editor = "vim"
		} else if _, err := exec.LookPath("nano"); err == nil {
			editor = "nano"
		} else {
			editor = "vi"
		}
	}

	ext := ".txt"
	if fieldIndex == 1 {
		ext = ".tex"
	}

	tmpFile, err := ioutil.TempFile("", "cv_radar_*"+ext)
	if err != nil {
		return func() tea.Msg {
			return editorFinishedMsg{fieldIndex: fieldIndex, tempPath: "", err: err}
		}
	}

	tmpPath := tmpFile.Name()
	if initialContent != "" {
		tmpFile.WriteString(initialContent)
	}
	tmpFile.Close()

	c := exec.Command(editor, tmpPath)
	return tea.ExecProcess(c, func(err error) tea.Msg {
		return editorFinishedMsg{
			fieldIndex: fieldIndex,
			tempPath:   tmpPath,
			err:        err,
		}
	})
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}

		switch m.state {
		case 0:
			if msg.String() == "ctrl+g" {
				if m.focusIndex < 2 {
					currentVal := m.inputs[m.focusIndex].Value()
					return m, openExternalEditor(m.focusIndex, currentVal)
				}
			}

			switch msg.String() {
			case "q":
				if m.focusIndex == 3 {
					return m, tea.Quit
				}
			case "tab", "down":
				m.focusIndex = (m.focusIndex + 1) % 4
				return m.updateFocus()
			case "shift+tab", "up":
				m.focusIndex = (m.focusIndex - 1 + 4) % 4
				return m.updateFocus()
			case "left", "h":
				if m.focusIndex == 2 {
					m.aiCursor = (m.aiCursor - 1 + len(m.availableAIs)) % len(m.availableAIs)
					m.selectedAI = m.availableAIs[m.aiCursor]
					return m, nil
				}
			case "right", "l", "space":
				if m.focusIndex == 2 {
					m.aiCursor = (m.aiCursor + 1) % len(m.availableAIs)
					m.selectedAI = m.availableAIs[m.aiCursor]
					return m, nil
				}
			case "enter":
				if m.focusIndex < 3 {
					m.focusIndex++
					return m.updateFocus()
				}
				// Button pressed ("Let's Go!")
				m.submitInputs()
				m.state = 1
				return m, nil
			}

			// Handle text input edits for fields 0 and 1
			var cmd tea.Cmd
			if m.focusIndex < 2 {
				m.inputs[m.focusIndex], cmd = m.inputs[m.focusIndex].Update(msg)
			}
			return m, cmd

		case 1:
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
				m.state = 2
				return m, func() tea.Msg {
					parsed := parser.ParseText(m.cvText)
					engine := ats.Get(m.platform)
					var res models.ATSResult
					if engine != nil {
						res = engine.Analyze(parsed, m.jdText)
					}
					return res
				}
			case "q":
				return m, tea.Quit
			}

		case 3:
			if msg.String() == "o" || msg.String() == "O" {
				m.state = 4
				return m, tea.Batch(m.spinner.Tick, runOptimizeCmd(m.cvText, m.jdText, m.platform, m.selectedAI))
			}
			if msg.String() == "q" || msg.String() == "Q" {
				return m, tea.Quit
			}

		case 5:
			if msg.String() == "enter" || msg.String() == "q" {
				return m, tea.Quit
			}
		}

	case editorFinishedMsg:
		if msg.err == nil && msg.tempPath != "" {
			if b, err := ioutil.ReadFile(msg.tempPath); err == nil {
				content := strings.TrimSpace(string(b))
				if content != "" {
					m.inputs[msg.fieldIndex].SetValue(content)
				}
			}
			os.Remove(msg.tempPath)
		}
		return m, nil

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case models.ATSResult:
		m.result = &msg
		m.state = 3

	case optimizeResult:
		m.err = msg.err
		m.state = 5

	case error:
		m.err = msg
		m.state = 3

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m *model) updateFocus() (tea.Model, tea.Cmd) {
	for i := 0; i < len(m.inputs); i++ {
		if i == m.focusIndex {
			m.inputs[i].Focus()
		} else {
			m.inputs[i].Blur()
		}
	}
	return *m, nil
}

func (m *model) submitInputs() {
	jdInput := strings.TrimSpace(m.inputs[0].Value())
	if jdInput == "" {
		jdInput = "jd.txt"
	}
	if b, err := ioutil.ReadFile(jdInput); err == nil {
		m.jdText = string(b)
	} else if b, err := ioutil.ReadFile("jd_test.txt"); err == nil && jdInput == "jd.txt" {
		m.jdText = string(b)
	} else {
		m.jdText = jdInput
	}

	cvInput := strings.TrimSpace(m.inputs[1].Value())
	if cvInput == "" {
		cvInput = "cv.tex"
	}
	if b, err := ioutil.ReadFile(cvInput); err == nil {
		m.cvText = string(b)
	} else if b, err := ioutil.ReadFile("cv_test.txt"); err == nil && cvInput == "cv.tex" {
		m.cvText = string(b)
	} else {
		m.cvText = cvInput
	}
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

func findOptimizeScript() string {
	// 1. Check current working directory
	if _, err := os.Stat("./optimize.sh"); err == nil {
		return "./optimize.sh"
	}

	// 2. Check directory of cv-tui binary (resolving symlinks)
	if execPath, err := os.Executable(); err == nil {
		if evalPath, err := filepath.EvalSymlinks(execPath); err == nil {
			scriptPath := filepath.Join(filepath.Dir(evalPath), "optimize.sh")
			if _, err := os.Stat(scriptPath); err == nil {
				return scriptPath
			}
		}
	}

	// 3. Check ~/.cv-radar/optimize.sh fallback
	if home, err := os.UserHomeDir(); err == nil {
		globalPath := filepath.Join(home, ".cv-radar", "optimize.sh")
		if _, err := os.Stat(globalPath); err == nil {
			return globalPath
		}
	}

	return "./optimize.sh"
}

func runOptimizeCmd(cvText string, jdText string, platform string, selectedAI string) tea.Cmd {
	return func() tea.Msg {
		scriptPath := findOptimizeScript()

		// Write CV text to temp file
		cvTmp, err := ioutil.TempFile("", "cv_input_*.tex")
		if err != nil {
			return optimizeResult{err: fmt.Errorf("failed to create temp CV file: %v", err)}
		}
		defer os.Remove(cvTmp.Name())
		if strings.TrimSpace(cvText) == "" {
			cvText = "Backend Engineer CV"
		}
		cvTmp.WriteString(cvText)
		cvTmp.Close()

		// Write JD text to temp file
		jdTmp, err := ioutil.TempFile("", "jd_input_*.txt")
		if err != nil {
			return optimizeResult{err: fmt.Errorf("failed to create temp JD file: %v", err)}
		}
		defer os.Remove(jdTmp.Name())
		if strings.TrimSpace(jdText) == "" {
			jdText = "Backend Engineer"
		}
		jdTmp.WriteString(jdText)
		jdTmp.Close()

		cmd := exec.Command(scriptPath, cvTmp.Name(), jdTmp.Name(), platform, selectedAI)
		out, err := cmd.CombinedOutput()
		if err != nil {
			lines := strings.Split(strings.TrimSpace(string(out)), "\n")
			lastLine := err.Error()
			for i := len(lines) - 1; i >= 0; i-- {
				if l := strings.TrimSpace(lines[i]); l != "" {
					lastLine = l
					break
				}
			}
			return optimizeResult{err: fmt.Errorf("%s", lastLine)}
		}
		return optimizeResult{err: nil}
	}
}

var (
	bannerStyle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FFCC00"))
	badgeStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#00FFFF")).Background(lipgloss.Color("#1A1A2E")).Padding(0, 1)
	titleStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FF75B5")).MarginBottom(1)
	boxStyle     = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#874BFD")).Padding(1, 2)
	focusedStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FFCC00"))
	blurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#888888"))
	btnFocused   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#000000")).Background(lipgloss.Color("#FFCC00")).Padding(0, 2)
	btnBlurred   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FFFFFF")).Background(lipgloss.Color("#333333")).Padding(0, 2)
)

func (m model) View() string {
	if m.width == 0 {
		return "Initializing..."
	}

	var content string

	switch m.state {
	case 0:
		asciiBanner := bannerStyle.Render(`
  ██████╗██╗   ██╗     ██████╗  █████╗ ██████╗  █████╗ ██████╗ 
 ██╔════╝██║   ██║     ██╔══██╗██╔══██╗██╔══██╗██╔══██╗██╔══██╗
 ██║     ██║   ██║     ██████╔╝███████║██║  ██║███████║██████╔╝
 ██║     ╚██╗ ██╔╝     ██╔══██╗██╔══██║██║  ██║██╔══██║██╔══██╗
 ╚██████╗ ╚████╔╝      ██║  ██║██║  ██║██████╔╝██║  ██║██║  ██║
  ╚═════╝  ╚═══╝       ╚═╝  ╚═╝╚═╝  ╚═╝╚═════╝ ╚═╝  ╚═╝╚═════╝ `)

		subtitle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#888888")).Render("   ATS RESUME SCANNER & AI CV OPTIMIZER")
		aiBadge := badgeStyle.Render(fmt.Sprintf("🤖 AI Engine Selected: %s", strings.ToUpper(m.selectedAI)))

		var jdLabel, cvLabel, aiLabel string
		if m.focusIndex == 0 {
			jdLabel = focusedStyle.Render("► 1. Job Description File (path, text, or Ctrl+G for Vim):")
		} else {
			jdLabel = blurredStyle.Render("  1. Job Description File (path, text, or Ctrl+G for Vim):")
		}

		if m.focusIndex == 1 {
			cvLabel = focusedStyle.Render("► 2. CV Upload Path or LaTeX (file, paste, or Ctrl+G for Vim):")
		} else {
			cvLabel = blurredStyle.Render("  2. CV Upload Path or LaTeX (file, paste, or Ctrl+G for Vim):")
		}

		if m.focusIndex == 2 {
			aiLabel = focusedStyle.Render(fmt.Sprintf("► 3. AI Backend Engine:  ◄ [ %s ] ►  (Use ← / → arrows to toggle)", strings.ToUpper(m.selectedAI)))
		} else {
			aiLabel = blurredStyle.Render(fmt.Sprintf("  3. AI Backend Engine:    [ %s ]", strings.ToUpper(m.selectedAI)))
		}

		var btn string
		if m.focusIndex == 3 {
			btn = btnFocused.Render("[ 🚀 LET'S GO! ]")
		} else {
			btn = btnBlurred.Render("[   LET'S GO!   ]")
		}

		helpText := lipgloss.NewStyle().Foreground(lipgloss.Color("#666666")).Render("(Press Tab to navigate, ←/→ for AI, Ctrl+G for Vim editor, Enter to submit)")

		content = asciiBanner + "\n" + subtitle + "\n\n" + aiBadge + "\n\n" +
			jdLabel + "\n" + m.inputs[0].View() + "\n\n" +
			cvLabel + "\n" + m.inputs[1].View() + "\n\n" +
			aiLabel + "\n\n" +
			btn + "\n\n" + helpText

	case 1:
		content = titleStyle.Render("Select Target ATS Platform:") + "\n\n"
		for i, p := range m.platforms {
			cursor := " "
			if m.cursor == i {
				cursor = ">"
			}
			content += fmt.Sprintf("%s %s\n", cursor, p)
		}
		content += "\n(Use Up/Down arrows, Press Enter to select)"

	case 2:
		content = fmt.Sprintf("\n\n   %s Scanning CV against Job Description...\n\n   Evaluating layout, exact keywords, and structural integrity...", m.spinner.View())

	case 3:
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

			content += fmt.Sprintf("\nAI Backend Engine: %s\n", badgeStyle.Render(strings.ToUpper(m.selectedAI)))
			content += "\nPress [O] to Optimize with AI (Headless)"
			content += "\nPress [Q] to Quit"
		}

	case 4:
		content = fmt.Sprintf("\n\n   %s Rewriting CV using %s LLM...\n\n   Please wait, generating optimized LaTeX and PDF...", m.spinner.View(), strings.ToUpper(m.selectedAI))

	case 5:
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
