package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ochanuco/prr/internal/github"
)

type prListMsg []github.PR
type errMsg struct{ err error }

func (e errMsg) Error() string { return e.err.Error() }

type Model struct {
	repo    string
	state   string
	table   table.Model
	err     error
	loading bool
}

func New(repo, state string) Model {
	columns := []table.Column{
		{Title: "#", Width: 5},
		{Title: "Title", Width: 50},
		{Title: "Author", Width: 16},
		{Title: "Updated", Width: 16},
		{Title: "State", Width: 10},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithHeight(12),
	)

	headerStyle := lipgloss.NewStyle().Bold(true)
	selectedStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("229")).Background(lipgloss.Color("57")).Bold(true)
	styles := table.DefaultStyles()
	styles.Header = styles.Header.Inherit(headerStyle)
	styles.Selected = styles.Selected.Inherit(selectedStyle)
	t.SetStyles(styles)

	return Model{
		repo:    repo,
		state:   state,
		table:   t,
		loading: true,
	}
}

func (m Model) Init() tea.Cmd {
	return fetchPRs(m.repo, m.state)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	case prListMsg:
		m.loading = false
		m.table.SetRows(prsToRows([]github.PR(msg)))
		return m, nil
	case errMsg:
		m.loading = false
		m.err = msg.err
		return m, nil
	}

	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	if m.loading {
		return "Loading PRs...\n"
	}
	if m.err != nil {
		return fmt.Sprintf("Error: %v\n", m.err)

	}

	header := lipgloss.NewStyle().Bold(true).Render(fmt.Sprintf("PRs: %s (%s)", m.repo, m.state))
	return header + "\n\n" + m.table.View() + "\n\n(q to quit)\n"
}

func fetchPRs(repo, state string) tea.Cmd {
	return func() tea.Msg {
		prs, err := github.FetchPRs(repo, state)
		if err != nil {
			return errMsg{err: err}
		}
		return prListMsg(prs)
	}
}

func prsToRows(prs []github.PR) []table.Row {
	rows := make([]table.Row, 0, len(prs))
	for _, pr := range prs {
		state := pr.State
		if pr.Draft {
			state = "draft"
		}
		rows = append(rows, table.Row{
			fmt.Sprintf("%d", pr.Number),
			truncate(pr.Title, 48),
			pr.User.Login,
			pr.UpdatedAt.Local().Format("2006-01-02"),
			state,
		})
	}
	return rows
}

func truncate(s string, max int) string {
	runes := []rune(s)
	if len(runes) <= max {
		return s
	}
	if max <= 3 {
		return string(runes[:max])
	}
	return string(runes[:max-3]) + "..."
}
