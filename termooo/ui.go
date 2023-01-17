package termooo

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/exp/slices"

	tea "github.com/charmbracelet/bubbletea"
)

var (
	alphabet = "A B C D E F G H I J K L M N O P Q R S T U V W X Y Z"

	err        = createColoredStyle("#FE6F6D", false)
	wrongGuess = createColoredStyle("#ECB225", false)
	rightGuess = createColoredStyle("#94E07A", false)

	words, _ = readFile("words.txt")
)

type model struct {
	word     *word
	target   string
	results  []string
	alphabet *string

	textarea textarea.Model

	viewport viewport.Model
	err      viewport.Model
	result   viewport.Model
	end      viewport.Model
}

func createColoredStyle(color string, bold bool) lipgloss.Style {
	defaultStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(color))

	if !bold {
		return defaultStyle
	} else {
		defaultStyle = defaultStyle.Bold(true)
		return defaultStyle
	}
}

func (m model) colorizeAlphabet(word string) {
	var (
		green  = createColoredStyle("#54D854", false)
		yellow = createColoredStyle("#EDED66", false)
		gray   = createColoredStyle("#525552", false)
	)

	m.word.applyAll(word)

	for _, alphaLet := range *m.alphabet {
		for _, value := range m.word.status {
			upper := strings.ToUpper(value.letter)

			if upper == string(alphaLet) {
				switch value.color {
				case "Green":
					colored := green.Render(upper)
					*m.alphabet = strings.Replace(*m.alphabet, upper, colored, 1)

				case "Yellow":
					if strings.Contains(*m.alphabet, green.Render(upper)) {
						continue
					}

					colored := yellow.Render(upper)
					*m.alphabet = strings.Replace(*m.alphabet, upper, colored, 1)

				case "Gray":
					if strings.Contains(*m.alphabet, green.Render(upper)) ||
						strings.Contains(*m.alphabet, yellow.Render(upper)) {

						continue
					}

					colored := gray.Render(upper)
					*m.alphabet = strings.Replace(*m.alphabet, upper, colored, 1)
				}
			}
		}
	}
}

func (m model) colorizeWordStatus() string {
	var colorfulWord string

	var (
		green  = createColoredStyle("#5CD201", true)
		yellow = createColoredStyle("#E8DF2E", true)
		gray   = createColoredStyle("#808079", true)
	)

	m.word.applyAll(m.target)

	for i := 0; i < 5; i++ {
		letterColor := m.word.status[i].color
		letter := strings.ToUpper(string(m.word.status[i].letter))

		switch letterColor {
		case "Green":
			colorfulWord += green.Render(letter + "  ")

		case "Yellow":
			colorfulWord += yellow.Render(letter + "  ")

		case "Gray":
			colorfulWord += gray.Render(letter + "  ")
		}
	}

	return colorfulWord
}

func InitialModel() model {

	textArea := textarea.New()
	textArea.Placeholder = "Digite uma palavra..."
	textArea.Focus()

	textArea.CharLimit = 5

	textArea.SetWidth(30)
	textArea.SetHeight(1)

	textArea.ShowLineNumbers = false

	textArea.FocusedStyle.CursorLine = lipgloss.NewStyle()
	textArea.KeyMap.InsertNewline.SetEnabled(false)

	// Target
	target, _ := generateTarget()

	// Viewport
	vp := viewport.New(100, 6)

	// Error Viewport
	errVp := viewport.New(100, 1)

	// Result Viewport
	resVp := viewport.New(100, 1)

	// endVp
	endVp := viewport.New(100, 1)
	endVp.SetContent("Feito por isaacszf ⚬ Sair (Ctrl+C)")

	return model{
		textarea: textArea,

		viewport: vp,
		err:      errVp,
		result:   resVp,
		end:      endVp,

		word:     &word{status: make(map[int]letterStatus)},
		target:   target,
		results:  []string{},
		alphabet: &alphabet,
	}
}

func (m model) Init() tea.Cmd {
	return textarea.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd    tea.Cmd
		vpCmd    tea.Cmd
		errVpCmd tea.Cmd
	)

	m.textarea, tiCmd = m.textarea.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)

	m.err, errVpCmd = m.err.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit

		case tea.KeyEnter:
			guess := m.textarea.Value()

			if utf8.RuneCountInString(guess) != 5 {
				msg := err.Render("Só palavras com 5 letras!")

				m.err.SetContent(msg)
				return m, tea.Batch(tiCmd, errVpCmd)
			}

			if !isAlpha(guess) {
				msg := err.Render("Apenas letras são permitidas!")

				m.err.SetContent(msg)
				return m, tea.Batch(tiCmd, errVpCmd)
			}

			if m.err.View() != "" {
				m.err.SetContent("")
			}

			parsedGuess := strings.ToLower(normalizeWord(guess))
			if !slices.Contains(words, parsedGuess) {
				msg := err.Render("Essa palavra é inválida!")

				m.err.SetContent(msg)
				return m, tea.Batch(tiCmd, errVpCmd)
			}

			m.word.guess = parsedGuess

			result := m.colorizeWordStatus()
			m.colorizeAlphabet(m.target)

			m.results = append(m.results, result)
			m.viewport.SetContent(strings.Join(m.results, "\n"))

			// Resetting fields
			m.word.guess = ""
			m.word.status = make(map[int]letterStatus)

			m.textarea.Reset()

			if m.target == parsedGuess {
				msg := rightGuess.Render("Você ganhou!")

				m.result.SetContent(msg)
				return m, tea.Quit
			}

			if len(m.results) == 6 {
				msg := wrongGuess.Render("Você perdeu! A palavra era: " + m.target)

				m.result.SetContent(msg)
				return m, tea.Quit
			}
		}
	}

	return m, tea.Batch(tiCmd, vpCmd)
}

func (m model) View() string {
	var alphabet = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#333333")).
		Padding(0, 3).
		MarginLeft(2)

	var style = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#333333")).
		Align(lipgloss.Center).
		MarginTop(1).
		MarginLeft(13).
		Width(35)

	var end = lipgloss.NewStyle().
		MarginTop(2).
		MarginLeft(15).
		MarginBottom(1).
		Foreground(lipgloss.Color("#414141")).
		Align(lipgloss.Center)

	template := fmt.Sprintf(
		"%s\n%s\n\n%s\n\n%s",
		m.err.View(),
		m.textarea.View(),
		m.viewport.View(),
		m.result.View(),
	) + "\n\n"

	return "\n" +
		alphabet.Render(*m.alphabet) +
		style.Render(template) + end.Render(m.end.View()+"\n")
}
