package model

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type DashboardState int

const (
	DashMenu DashboardState = iota
)

type DashboardModel struct {
	state    DashboardState
	username string
}

// Constructor
func NewDashboard(username string) DashboardModel {
	return DashboardModel{
		state:    DashMenu,
		username: username,
	}
}

func (m DashboardModel) Init() tea.Cmd {
	return nil
}

func (m DashboardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.state {
	case DashMenu:
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			switch keyMsg.String() {
			case "1":
				// Aquí iría lógica para ver el balance
				return m, tea.Println("\n[🔍] Ver balance (por implementar)")
			case "2":
				// Aquí iría lógica para agregar una transacción
				return m, tea.Println("\n[➕] Agregar transacción (por implementar)")
			case "3":
				// Aquí iría lógica para ver el historial
				return m, tea.Println("\n[📜] Ver historial (por implementar)")
			case "q":
				// Cierra sesión
				return m, tea.Quit
			}
		}
	}
	return m, nil
}

func (m DashboardModel) View() string {
	return fmt.Sprintf(`
👤 Usuario: %s

📊 Dashboard:
1) Ver Balance
2) Agregar Transacción
3) Ver Historial
q) Cerrar Sesión
`, m.username)
}
