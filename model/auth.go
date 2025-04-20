package model

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type AuthState int

const (
	AuthUsername AuthState = iota
	AuthPassword
)

type AuthAction int

const (
	AuthLogin AuthAction = iota
	AuthRegister
)

type AuthModel struct {
	state     AuthState
	action    AuthAction
	username  string
	password  string
	input     textinput.Model
	errMsg    string
	userData  UserActivity
	session   string
	onSuccess func(user UserActivity, path string, username string) tea.Cmd
}

// Constructor
func NewAuthModel(action AuthAction, onSuccess func(UserActivity, string, string) tea.Cmd) *AuthModel {
	ti := textinput.New()
	ti.Placeholder = "Username"
	ti.Focus()
	ti.CharLimit = 32
	ti.Width = 30

	return &AuthModel{
		state:     AuthUsername,
		action:    action,
		input:     ti,
		onSuccess: onSuccess,
	}
}

func (m AuthModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m *AuthModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.state {
		case AuthUsername:
			if msg.Type == tea.KeyEnter {
				m.username = m.input.Value()
				m.state = AuthPassword
				m.input.Placeholder = "Password"
				m.input.EchoMode = textinput.EchoPassword
				m.input.Reset()
			}
		case AuthPassword:
			if msg.Type == tea.KeyEnter {
				m.password = m.input.Value()
				var user UserActivity
				var err error

				if m.action == AuthLogin {
					user, err = login(m.username, m.password)
				} else {
					user, err = register(m.username, m.password)
				}

				if err != nil {
					m.errMsg = err.Error()
					m.state = AuthUsername
					m.input.Placeholder = "Username"
					m.input.EchoMode = textinput.EchoNormal
					m.input.Reset()
					return m, nil
				}

				m.userData = user
				m.session = "./data/" + m.username + ".json"
				return m, m.onSuccess(m.userData, m.session, m.username)
			}
		}
	case tea.WindowSizeMsg:
		m.input.Width = msg.Width
	}

	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m AuthModel) View() string {
	var builder strings.Builder

	switch m.state {
	case AuthUsername, AuthPassword:
		builder.WriteString(fmt.Sprintf("\n%s:\n%s\n\n(Press Enter to continue)", strings.Title(m.input.Placeholder), m.input.View()))
		if m.errMsg != "" {
			builder.WriteString(fmt.Sprintf("\n\n❌ Error: %s", m.errMsg))
		}
	}

	return builder.String()
}
