// model/dashboard.go
package model

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type DashboardModel struct {
	username   string
	options    []string
	cursor     int
	errMsg     string
	balance    int
	balanceMsg string // Para mostrar el balance correctamente
}

func NewDashboard(username string) DashboardModel {
	return DashboardModel{
		username: username,
		options: []string{
			"Ver balance",
			"Agregar transacción",
			"Ver historial",
			"Cerrar sesión",
		},
		cursor: 0,
	}
}

func (m DashboardModel) Init() tea.Cmd {
	return nil
}

func (m DashboardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.options)-1 {
				m.cursor++
			}
		case "enter":
			return m.handleSelection()
		case "q":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m DashboardModel) View() string {
	s := fmt.Sprintf("\nHola, %s\n\n", m.username)
	s += "Selecciona una opción:\n\n"

	for i, option := range m.options {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		s += fmt.Sprintf(" %s %s\n", cursor, option)
	}

	if m.balanceMsg != "" {
		s += fmt.Sprintf("\n%s\n", m.balanceMsg)
	}

	if m.errMsg != "" {
		s += fmt.Sprintf("\n[Error]: %s\n", m.errMsg)
	}

	s += "\nPresiona 'q' para salir.\n"
	return s
}

// Puedes luego reemplazar este método con eventos reales según lo que necesites.
func (m DashboardModel) handleSelection() (tea.Model, tea.Cmd) {
	switch m.cursor {
	case 0:
		balance, err := getBalance(m.username)
		if err != nil {
			m.errMsg = err.Error()
			m.balanceMsg = ""
		} else {
			m.balance = balance
			m.balanceMsg = fmt.Sprintf("Balance: $%d", m.balance)
			m.errMsg = ""
		}

	case 1:
		m.errMsg = "➕ Agregar transacción aún no implementado"
	case 2:
		m.errMsg = "📜 Ver historial aún no implementado"
	case 3:
		return m, tea.Quit
	}
	return m, nil
}
