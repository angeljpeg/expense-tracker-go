package main

import (
	"fmt"
	"os"

	"github.com/angeljpeg/expense-tracker-go/model"
	tea "github.com/charmbracelet/bubbletea"
)

type AppState int

const (
	AppMenu AppState = iota
	AppAuth
	AppDashboard
)

type MainModel struct {
	state       AppState
	currentView tea.Model
	user        model.UserActivity
	username    string
	sessionPath string
}

func NewMainModel() *MainModel {
	return &MainModel{
		state: AppMenu,
	}
}

func (m MainModel) Init() tea.Cmd {
	return nil
}

func (m *MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Delegamos a UpdateWithDashboard para detectar el mensaje personalizado
	if msg, ok := msg.(dashboardInitMsg); ok {
		return m.UpdateWithDashboard(msg)
	}

	switch m.state {

	case AppMenu:
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			switch keyMsg.String() {
			case "1":
				m.state = AppAuth
				m.currentView = model.NewAuthModel(model.AuthRegister, m.authSuccess)
				return m, m.currentView.Init()
			case "2":
				m.state = AppAuth
				m.currentView = model.NewAuthModel(model.AuthLogin, m.authSuccess)
				return m, m.currentView.Init()
			case "q", "ctrl+c":
				return m, tea.Quit
			}
		}

	case AppAuth, AppDashboard:
		var cmd tea.Cmd
		m.currentView, cmd = m.currentView.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m MainModel) View() string {
	switch m.state {
	case AppMenu:
		return "\nWelcome to the Expense Tracker\n\n1) Register\n2) Login\n\nPress 'q' to quit.\n"
	case AppAuth, AppDashboard:
		return m.currentView.View()
	default:
		return ""
	}
}

func (m *MainModel) authSuccess(user model.UserActivity, path string, username string) tea.Cmd {
	return func() tea.Msg {
		return dashboardInitMsg{
			user:     user,
			path:     path,
			username: username,
		}
	}
}

type dashboardInitMsg struct {
	user     model.UserActivity
	path     string
	username string
}

// Este método maneja el mensaje custom del dashboard
func (m *MainModel) UpdateWithDashboard(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case dashboardInitMsg:
		fmt.Println("Recibido dashboardInitMsg, cambiando a Dashboard")

		m.state = AppDashboard
		m.user = msg.user
		m.username = msg.username
		m.sessionPath = msg.path
		m.currentView = model.NewDashboard(m.username)
		cmd := m.currentView.Init()
		return m, cmd
	}
	return m.Update(msg)
}

// ===== MAIN =====
func main() {
	program := tea.NewProgram(NewMainModel(), tea.WithAltScreen())

	if err := program.Start(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
