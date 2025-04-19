package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type TransactionType string

const (
	Income  TransactionType = "income"
	Expense TransactionType = "expense"
)

type Transaction struct {
	ID          int             `json:"id"`
	Type        TransactionType `json:"type"`
	Amount      float64         `json:"amount"`
	Category    string          `json:"category"`
	Description string          `json:"description"`
	Date        time.Time       `json:"date"`
}

type UserActivity struct {
	Balance      float64       `json:"balance"`
	Password     string        `json:"password"`
	Transactions []Transaction `json:"transactions"`
}

// Estados de la app
type state int

const (
	stateMenu state = iota
	stateInputUsername
	stateInputPassword
	stateSuccess
	stateError
)

type action int

const (
	actionLogin action = iota
	actionRegister
)

type model struct {
	state    state
	action   action
	username string
	password string
	input    textinput.Model
	errorMsg string
	userData UserActivity
	session  string
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "Enter a number"
	ti.Focus()
	ti.CharLimit = 32
	ti.Width = 30

	return model{
		state: stateMenu,
		input: ti,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.state {

		case stateMenu:
			switch msg.String() {
			case "1":
				m.action = actionRegister
				m.state = stateInputUsername
				m.input.Placeholder = "Username"
				m.input.Reset()
			case "2":
				m.action = actionLogin
				m.state = stateInputUsername
				m.input.Placeholder = "Username"
				m.input.Reset()
			case "q", "ctrl+c":
				return m, tea.Quit
			}

		case stateInputUsername:
			if msg.Type == tea.KeyEnter {
				m.username = m.input.Value()
				m.state = stateInputPassword
				m.input.Placeholder = "Password"
				m.input.Reset()
				return m, nil
			}

		case stateInputPassword:
			if msg.Type == tea.KeyEnter {
				m.password = m.input.Value()
				var user UserActivity
				var err error

				if m.action == actionLogin {
					user, err = login(m.username, m.password)
				} else {
					user, err = register(m.username, m.password)
				}

				if err != nil {
					m.state = stateError
					m.errorMsg = err.Error()
					return m, nil
				}

				m.userData = user
				m.session = "./data/" + m.username + ".json"
				m.state = stateSuccess
				return m, nil
			}
		case stateError, stateSuccess:
			if msg.String() == "enter" || msg.String() == "q" {
				return m, tea.Quit
			}
		}

	case tea.WindowSizeMsg:
		m.input.Width = msg.Width
	}

	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m model) View() string {
	switch m.state {
	case stateMenu:
		return "\nWelcome to the Expense Tracker\n\n1) Register\n2) Login\n\nPress 'q' to quit.\n"
	case stateInputUsername, stateInputPassword:
		return fmt.Sprintf("\n%s:\n%s\n\n(Press Enter to continue)", strings.Title(m.input.Placeholder), m.input.View())
	case stateSuccess:
		return fmt.Sprintf("\n✅ Welcome %s! Your balance is $%.2f\nSession file: %s\n\nPress Enter to quit.\n", m.username, m.userData.Balance, m.session)
	case stateError:
		return fmt.Sprintf("\n❌ Error: %s\n\nPress Enter to quit.\n", m.errorMsg)
	default:
		return ""
	}
}

// ========== Lógica de login y register ==========
func login(user string, password string) (UserActivity, error) {
	filePath := "./data/" + user + ".json"
	data, err := os.ReadFile(filePath)
	if err != nil {
		return UserActivity{}, fmt.Errorf("user not found")
	}
	var ua UserActivity
	if err := json.Unmarshal(data, &ua); err != nil {
		return UserActivity{}, fmt.Errorf("error parsing file")
	}
	if password != ua.Password {
		return UserActivity{}, fmt.Errorf("incorrect password")
	}
	return ua, nil
}

func register(user string, password string) (UserActivity, error) {
	filePath := "./data/" + user + ".json"
	if _, err := os.Stat(filePath); err == nil {
		return UserActivity{}, fmt.Errorf("user already exists")
	}

	file, err := os.Create(filePath)
	if err != nil {
		return UserActivity{}, fmt.Errorf("error creating file")
	}
	defer file.Close()

	ua := UserActivity{
		Balance:      0,
		Password:     password,
		Transactions: []Transaction{},
	}

	if err := json.NewEncoder(file).Encode(ua); err != nil {
		return UserActivity{}, fmt.Errorf("error writing file")
	}
	return ua, nil
}

// ========== Main ==========
func main() {
	if _, err := tea.NewProgram(initialModel()).Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
