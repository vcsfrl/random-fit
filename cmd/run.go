package cmd

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"github.com/vcsfrl/random-fit/internal/service"
)

const (
	colorAccent  = "170"
	colorSubtle  = "241"
	colorSuccess = "42"
	colorError   = "196"

	cursorSymbol  = "▸ "
	paddingPrefix = "  "
	inputCursor   = "█"
	inputPrefix   = "> "

	executeChoice = "Execute"
	cancelChoice  = "Cancel"

	commandGenerate   = "generate"
	commandDefinition = "definition"
	commandCode       = "code"

	subCommandCombination = "combination"
	subCommandPlan        = "plan"

	operationList   = "list"
	operationNew    = "new"
	operationEdit   = "edit"
	operationDelete = "delete"

	selectionGenerate       = 0
	selectionSubCommand     = 1
	selectionOperation      = 2
	selectionItem           = 3
	minGenerateConfirmDepth = 4
	minDefinitionItemDepth  = 4

	confirmChoiceCount = 2
)

type runModel struct {
	selections []string
	choices    []string
	cursor     int
	textInput  string
	inputMode  bool
	title      string
	conf       *Config
	quitting   bool
	execute    bool
	errorMsg   string
}

func newRunModel(conf *Config) runModel {
	model := runModel{
		conf:       conf,
		selections: make([]string, 0),
	}

	return model.refreshedStep()
}

func (m runModel) Init() tea.Cmd {
	return nil
}

func (m runModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) { //nolint:ireturn
	keyMsg, isKeyMsg := msg.(tea.KeyMsg)
	if !isKeyMsg {
		return m, nil
	}

	if m.inputMode {
		return m.handleTextInput(keyMsg)
	}

	return m.handleListInput(keyMsg)
}

func (m runModel) View() string {
	if m.quitting {
		return "Cancelled.\n"
	}

	if m.execute {
		args := m.buildCommandArgs()

		return "\nExecuting: random-fit " + strings.Join(args, " ") + "\n\n"
	}

	return m.renderWizard()
}

func (m runModel) handleTextInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) { //nolint:ireturn
	switch msg.Type { //nolint:exhaustive
	case tea.KeyCtrlC:
		m.quitting = true

		return m, tea.Quit
	case tea.KeyEsc:
		return m.goBack()
	case tea.KeyBackspace:
		if len(m.textInput) > 0 {
			m.textInput = m.textInput[:len(m.textInput)-1]
		}

		return m, nil
	case tea.KeyEnter:
		return m.submitTextInput()
	case tea.KeyRunes:
		m.textInput += string(msg.Runes)

		return m, nil
	}

	return m, nil
}

func (m runModel) handleListInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) { //nolint:ireturn
	switch msg.String() {
	case "ctrl+c", "q":
		m.quitting = true

		return m, tea.Quit
	case "esc":
		return m.goBack()
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down", "j":
		if m.cursor < len(m.choices)-1 {
			m.cursor++
		}
	case "enter":
		return m.selectChoice()
	}

	return m, nil
}

func (m runModel) selectChoice() (tea.Model, tea.Cmd) { //nolint:ireturn
	if len(m.choices) == 0 {
		return m, nil
	}

	selected := m.choices[m.cursor]

	if m.isConfirmStep() {
		return m.handleConfirm(selected)
	}

	m.selections = append(m.selections, selected)
	m.cursor = 0
	m.errorMsg = ""
	m = m.refreshedStep()

	return m, nil
}

func (m runModel) handleConfirm(selected string) (tea.Model, tea.Cmd) { //nolint:ireturn
	if selected == executeChoice {
		m.execute = true

		return m, tea.Quit
	}

	m.quitting = true

	return m, tea.Quit
}

func (m runModel) submitTextInput() (tea.Model, tea.Cmd) { //nolint:ireturn
	name := strings.TrimSpace(m.textInput)

	if name == "" {
		m.errorMsg = "Name cannot be empty."

		return m, nil
	}

	m.selections = append(m.selections, name)
	m.textInput = ""
	m.inputMode = false
	m.cursor = 0
	m.errorMsg = ""
	m = m.refreshedStep()

	return m, nil
}

func (m runModel) goBack() (tea.Model, tea.Cmd) { //nolint:ireturn
	if len(m.selections) == 0 {
		m.quitting = true

		return m, tea.Quit
	}

	m.selections = m.selections[:len(m.selections)-1]
	m.cursor = 0
	m.textInput = ""
	m.inputMode = false
	m.errorMsg = ""
	m = m.refreshedStep()

	return m, nil
}

func (m runModel) refreshedStep() runModel {
	title, choices, isInput := m.computeStepInfo()
	m.title = title
	m.choices = choices
	m.inputMode = isInput

	return m
}

func (m runModel) computeStepInfo() (string, []string, bool) {
	depth := len(m.selections)

	if depth == 0 {
		return "What would you like to do?",
			[]string{commandGenerate, commandDefinition, commandCode}, false
	}

	switch m.selections[selectionGenerate] {
	case commandGenerate:
		return m.generateStepInfo()
	case commandDefinition:
		return m.definitionStepInfo()
	case commandCode:
		return m.codeStepInfo()
	}

	return "", nil, false
}

func (m runModel) generateStepInfo() (string, []string, bool) {
	depth := len(m.selections)

	switch depth {
	case selectionSubCommand:
		return "Generate what?", []string{subCommandCombination}, false
	case selectionOperation:
		return m.selectCombinationDefinitionStep()
	case selectionItem:
		return m.selectPlanDefinitionStep()
	case minGenerateConfirmDepth:
		return m.buildConfirmTitle(), []string{executeChoice, cancelChoice}, false
	}

	return "", nil, false
}

func (m runModel) definitionStepInfo() (string, []string, bool) {
	depth := len(m.selections)

	switch depth {
	case selectionSubCommand:
		return "Manage which type?",
			[]string{subCommandCombination, subCommandPlan}, false
	case selectionOperation:
		return "Select operation:",
			[]string{operationList, operationNew, operationEdit, operationDelete}, false
	case selectionItem:
		return m.definitionOperationStepInfo()
	case minDefinitionItemDepth:
		return m.buildConfirmTitle(), []string{executeChoice, cancelChoice}, false
	}

	return "", nil, false
}

func (m runModel) definitionOperationStepInfo() (string, []string, bool) {
	operation := m.selections[selectionOperation]

	switch operation {
	case operationList:
		return m.buildConfirmTitle(), []string{executeChoice, cancelChoice}, false
	case operationNew:
		return "Enter name for new definition:", nil, true
	case operationEdit, operationDelete:
		return m.definitionSelectItemStep()
	}

	return "", nil, false
}

func (m runModel) definitionSelectItemStep() (string, []string, bool) {
	defType := m.selections[selectionSubCommand]

	if defType == subCommandCombination {
		return m.selectCombinationDefinitionStep()
	}

	return m.selectPlanDefinitionStep()
}

func (m runModel) selectCombinationDefinitionStep() (string, []string, bool) {
	definitions := m.listCombinationDefinitions()

	if len(definitions) == 0 {
		return "No combination definitions found. Press esc to go back.", nil, false
	}

	return "Select combination definition:", definitions, false
}

func (m runModel) selectPlanDefinitionStep() (string, []string, bool) {
	plans := m.listPlanDefinitions()

	if len(plans) == 0 {
		return "No plan definitions found. Press esc to go back.", nil, false
	}

	return "Select plan definition:", plans, false
}

func (m runModel) codeStepInfo() (string, []string, bool) {
	depth := len(m.selections)

	switch depth {
	case selectionSubCommand:
		return "Code tools:", []string{commandGenerate}, false
	case selectionOperation:
		return m.buildConfirmTitle(), []string{executeChoice, cancelChoice}, false
	}

	return "", nil, false
}

func (m runModel) isConfirmStep() bool {
	return len(m.choices) == confirmChoiceCount &&
		m.choices[0] == executeChoice &&
		m.choices[1] == cancelChoice
}

func (m runModel) buildCommandArgs() []string {
	if len(m.selections) == 0 {
		return nil
	}

	switch m.selections[selectionGenerate] {
	case commandGenerate:
		return m.buildGenerateArgs()
	case commandDefinition:
		return m.buildDefinitionArgs()
	case commandCode:
		return []string{commandCode, commandGenerate}
	}

	return nil
}

func (m runModel) buildGenerateArgs() []string {
	if len(m.selections) < minGenerateConfirmDepth {
		return nil
	}

	return []string{
		commandGenerate, subCommandCombination,
		"--combination", m.selections[selectionOperation],
		"--plan", m.selections[selectionItem],
	}
}

func (m runModel) buildDefinitionArgs() []string {
	if len(m.selections) < selectionItem {
		return nil
	}

	defType := m.selections[selectionSubCommand]
	operation := m.selections[selectionOperation]

	if operation == operationList {
		return []string{commandDefinition, defType}
	}

	if len(m.selections) < minDefinitionItemDepth {
		return nil
	}

	return []string{commandDefinition, defType, operation, "--name", m.selections[selectionItem]}
}

func (m runModel) buildConfirmTitle() string {
	args := m.buildCommandArgs()

	return "Execute: random-fit " + strings.Join(args, " ")
}

func (m runModel) listCombinationDefinitions() []string {
	manager := service.NewCombinationStarDefinitionManager(m.conf.DefinitionFolder())

	definitions, err := manager.List()
	if err != nil {
		return nil
	}

	return definitions
}

func (m runModel) listPlanDefinitions() []string {
	manager := service.NewPlanDefinitionManager(m.conf.PlanFolder())

	plans, err := manager.List()
	if err != nil {
		return nil
	}

	return plans
}

func (m runModel) renderWizard() string {
	var builder strings.Builder

	accentStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(colorAccent))
	subtleStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(colorSubtle))
	successStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(colorSuccess))

	builder.WriteString(accentStyle.Render("Random Fit Interactive Mode"))
	builder.WriteString("\n\n")

	if len(m.selections) > 0 {
		breadcrumb := "Path: " + strings.Join(m.selections, " → ")
		builder.WriteString(subtleStyle.Render(breadcrumb))
		builder.WriteString("\n\n")
	}

	builder.WriteString(successStyle.Render(m.title))
	builder.WriteString("\n\n")

	if m.errorMsg != "" {
		errStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(colorError))
		builder.WriteString(errStyle.Render(m.errorMsg))
		builder.WriteString("\n\n")
	}

	if m.inputMode {
		fmt.Fprintf(&builder, "%s%s%s\n", inputPrefix, m.textInput, inputCursor)
	} else {
		m.renderChoices(&builder, accentStyle)
	}

	builder.WriteString("\n")
	builder.WriteString(subtleStyle.Render("↑/↓: navigate • enter: select • esc: back • q: quit"))
	builder.WriteString("\n")

	return builder.String()
}

func (m runModel) renderChoices(builder *strings.Builder, accentStyle lipgloss.Style) {
	for index, choice := range m.choices {
		if index == m.cursor {
			builder.WriteString(accentStyle.Render(paddingPrefix + cursorSymbol + choice))
		} else {
			builder.WriteString(paddingPrefix + paddingPrefix + choice)
		}

		builder.WriteString("\n")
	}
}

func runInteractiveCmd(rootCmd *cobra.Command) {
	var run = &cobra.Command{
		Use:   "run",
		Short: "Interactive mode for building and executing commands.",
		Long:  `Interactive mode that guides you through building CLI commands step by step.`,
		Run: func(cmd *cobra.Command, _ []string) {
			executeInteractiveRun(cmd)
		},
	}

	rootCmd.AddCommand(run)
}

func executeInteractiveRun(cmd *cobra.Command) {
	conf := NewConfig()
	program := tea.NewProgram(newRunModel(conf))

	finalModel, err := program.Run()
	if err != nil {
		cmd.PrintErrln("Error running interactive mode:", err)

		return
	}

	result, isRunModel := finalModel.(runModel)
	if !isRunModel {
		cmd.PrintErrln("Error: unexpected model type")

		return
	}

	if !result.execute {
		return
	}

	args := result.buildCommandArgs()
	cmd.Println("Executing:", "random-fit", strings.Join(args, " "))

	root := cmd.Root()
	root.SetArgs(args)

	if execErr := root.ExecuteContext(cmd.Context()); execErr != nil {
		cmd.PrintErrln("Error executing command:", execErr)
	}
}
