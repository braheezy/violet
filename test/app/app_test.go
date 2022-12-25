package app_test

import (
	"bytes"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

/*
Test Goals:
   1. Model: Test the data structures and their properties, such as the Environment and VM structs in the Violet model.

   2. View: Test the output of the View method to ensure that it is correct for different input data.

   3. Update: Test the Update method to ensure that it updates the model correctly in response to different messages.

   4. Interactions with the VagrantClient: Test the interactions with the VagrantClient, such as the RunCommand function, to ensure that they are working correctly and returning the expected results.

   5. User input handling: Test the handling of user input, such as keyboard and mouse events, to ensure that they are correctly processed and result in the expected changes to the model.
*/

type testModel struct {
	dummy string
}

func (m testModel) Init() tea.Cmd {
	return app.initGlobalStatus
}

func (m *testModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m *testModel) View() string {
	return ""
}

func TestEnvironment_Create(t *testing.T) {
	var out bytes.Buffer
	var in bytes.Buffer

	m := &testModel{}
	p := tea.NewProgram(m, tea.WithInput(&in), tea.WithOutput(&out))
	if _, err := p.Run(); err != nil {
		t.Error(err)
	}

}
