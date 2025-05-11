package main

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"slices"
	"strings"
	"sync"

	"github.com/c-bata/go-prompt"
)

type CommandSet struct {
	mu       sync.RWMutex
	commands map[string]string
}

var commandSet = &CommandSet{
	commands: make(map[string]string),
}

func (cs *CommandSet) setCommand(name, command string) {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	cs.commands[name] = command
}

func (cs *CommandSet) getCommand(name string) (string, bool) {
	cs.mu.RLock()
	defer cs.mu.RUnlock()
	command, exists := cs.commands[name]
	return command, exists
}

func (cs *CommandSet) listCommands() []string {
	cs.mu.RLock()
	defer cs.mu.RUnlock()

	commandList := make([]string, 0, len(cs.commands))
	for name, command := range cs.commands {
		commandList = append(commandList, fmt.Sprintf("%s: %s", name, command))
	}
	return commandList
}

type ValueMemory struct {
	mu     sync.RWMutex
	values map[string][]string
}

var valueMemory = &ValueMemory{
	values: make(map[string][]string),
}

var flagTypes = map[string]string{
	"-p":              "port",
	"--exclude-ports": "port",
	"--top-ports":     "port_count",
	"-PS":             "port",
	"-PA":             "port",
	"-PU":             "port",
	"-PY":             "port",

	"-iL":           "file",
	"-oN":           "file",
	"-oX":           "file",
	"-oS":           "file",
	"-oG":           "file",
	"--excludefile": "file",
	"--resume":      "file",

	"-S":        "ip",
	"--exclude": "ip_network",
	"-D":        "ip_list",
	"--proxies": "url_list",

	"--min-rtt-timeout": "time",
	"--max-rtt-timeout": "time",
	"--script-timeout":  "time",
	"--scan-delay":      "time",
	"--host-timeout":    "time",
}

func (vm *ValueMemory) addValue(flag, value string) {
	if value == "" {
		return
	}

	vm.mu.Lock()
	defer vm.mu.Unlock()

	value = strings.TrimSpace(value)

	if _, exists := vm.values[flag]; !exists {
		vm.values[flag] = []string{}
	}

	for i, v := range vm.values[flag] {
		if v == value {
			vm.values[flag] = append([]string{value}, slices.Delete(vm.values[flag], i, i+1)...)
			return
		}
	}

	vm.values[flag] = append([]string{value}, vm.values[flag]...)

	if flagType, exists := flagTypes[flag]; exists {
		for otherFlag, otherType := range flagTypes {
			if otherFlag != flag && otherType == flagType {
				if _, exists := vm.values[otherFlag]; !exists {
					vm.values[otherFlag] = []string{}
				}

				duplicate := false
				for i, v := range vm.values[otherFlag] {
					if v == value {
						vm.values[otherFlag] = append([]string{value}, slices.Delete(vm.values[otherFlag], i, i+1)...)
						duplicate = true
						break
					}
				}

				if !duplicate {
					vm.values[otherFlag] = append([]string{value}, vm.values[otherFlag]...)
				}
			}
		}
	}
}

func (vm *ValueMemory) getLastValue(flag string) string {
	vm.mu.RLock()
	defer vm.mu.RUnlock()

	values, exists := vm.values[flag]
	if !exists || len(values) == 0 {
		return ""
	}

	return values[0]
}

func createEnhancedSuggestions() []prompt.Suggest {
	enhanced := make([]prompt.Suggest, len(Suggestions))

	for i, s := range Suggestions {
		flag := strings.Fields(s.Text)[0]

		lastValue := valueMemory.getLastValue(flag)

		if lastValue != "" && strings.Contains(s.Text, "<") && strings.Contains(s.Text, ">") {
			re := regexp.MustCompile(`<[^>]+>`)

			enhanced[i] = prompt.Suggest{
				Text:        re.ReplaceAllString(s.Text, lastValue),
				Description: s.Description + " [last used: " + lastValue + "]",
			}
		} else if lastValue != "" && strings.HasSuffix(s.Text, flag) {
			enhanced[i] = prompt.Suggest{
				Text:        s.Text + " " + lastValue,
				Description: s.Description + " [last used: " + lastValue + "]",
			}
		} else {
			enhanced[i] = s
		}
	}

	commandList := commandSet.listCommands()
	for _, cmd := range commandList {
		parts := strings.SplitN(cmd, ":", 2)
		if len(parts) == 2 {
			enhanced = append(enhanced, prompt.Suggest{
				Text:        parts[0],
				Description: "Saved command set: " + strings.TrimSpace(parts[1]),
			})
		}
	}

	return enhanced
}

func parseCommand(cmd string) {
	parts := strings.Fields(cmd)
	var currentFlag string

	for _, part := range parts {
		if strings.HasPrefix(part, "-") {
			for flag := range flagTypes {
				if strings.HasPrefix(part, flag) && len(part) > len(flag) {
					value := part[len(flag):]
					valueMemory.addValue(flag, value)
					break
				}
			}

			currentFlag = part

			if strings.Contains(part, "=") {
				eqIndex := strings.Index(part, "=")
				currentFlag = part[:eqIndex]
				value := part[eqIndex+1:]
				valueMemory.addValue(currentFlag, value)
				currentFlag = ""
			}
		} else if currentFlag != "" {
			valueMemory.addValue(currentFlag, part)
			currentFlag = ""
		}
	}
}

func completer(d prompt.Document) []prompt.Suggest {
	word := strings.Trim(d.GetWordBeforeCursor(), " ")
	line := d.TextBeforeCursor()

	if strings.HasPrefix(line, "help") {
		return HelpCategories
	}

	if strings.HasPrefix(line, "set ") {
		// Empty suggestions
		return []prompt.Suggest{}
	}

	enhanced := createEnhancedSuggestions()

	if word == "--script" {
		enhanced = append(enhanced, ScriptCategories...)
	}

	return CustomFilter(enhanced, word)
}

func checkSudo(cmd string) string {
	if runtime.GOOS == "windows" {
		return cmd
	}
	for _, v := range SudoRequiredFlags {
		if strings.Contains(cmd, v) {
			fmt.Printf("The flag %s requires sudo, you may need to enter your password\n", v)
			return "sudo " + cmd
		}
	}
	return cmd
}

func execute(t string) {
	if t == "exit" || t == "quit" || t == "q" {
		fmt.Println("Exiting...")
		os.Exit(0)
	}

	if strings.HasPrefix(t, "help") {
		parts := strings.Fields(t)
		if len(parts) == 1 {
			fmt.Println("Please pick a help category")
			return
		}
		GetHelp(parts[1])
		return
	}

	if strings.HasPrefix(t, "set ") {
		parts := strings.Fields(t)
		if len(parts) < 3 {
			fmt.Println("Usage: set <name> <command>")
			return
		}
		name := parts[1]
		command := strings.Join(parts[2:], " ")
		commandSet.setCommand(name, command)
		fmt.Printf("Saved command set '%s': %s\n", name, command)
		return
	}

	if t == "list" {
		commands := commandSet.listCommands()
		if len(commands) == 0 {
			fmt.Println("No saved command sets")
			return
		}
		fmt.Println("Saved command sets:")
		for _, cmd := range commands {
			fmt.Println("  " + cmd)
		}
		return
	}

	parts := strings.Fields(t)
	newParts := make([]string, 0, len(parts))

	for _, part := range parts {
		if cmd, exists := commandSet.getCommand(part); exists {
			cmdParts := strings.Fields(cmd)
			newParts = append(newParts, cmdParts...)
		} else {
			newParts = append(newParts, part)
		}
	}

	// Actual command
	t = strings.Join(newParts, " ")

	t = "nmap " + t

	parseCommand(t)

	t = checkSudo(t)
	execParts := strings.Fields(t)
	var cmd *exec.Cmd
	if len(execParts) == 0 {
		fmt.Println("No command entered")
		return
	} else if len(execParts) == 1 && execParts[0] == "nmap" {
		cmd = exec.Command("nmap", "-h")
	} else {
		cmd = exec.Command(execParts[0], execParts[1:]...)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(output))
}

func main() {
	fmt.Println("NMAP Interactive CLI. Press 'Ctrl-D' to exit.")
	fmt.Println("Press tab for options")
	fmt.Println("Save variables/commands with 'set <name> <command>' and use them with 'nmap <name>'")
	fmt.Println("Use 'list' to see all saved variables/commands")
	fmt.Println("Type `help` for tips and more info.")

	p := prompt.New(
		func(input string) {
			execute(input)
		},
		completer,
		prompt.OptionPrefix("> nmap "),
		prompt.OptionTitle("NMAP Interactive CLI"),
		prompt.OptionSuggestionBGColor(prompt.DarkGray),
		prompt.OptionPrefixTextColor(prompt.Green),
		prompt.OptionDescriptionBGColor(prompt.DarkGray),
		prompt.OptionSelectedDescriptionBGColor(prompt.DarkGray),
		prompt.OptionPreviewSuggestionBGColor(prompt.DarkGray),
		prompt.OptionAddKeyBind(prompt.KeyBind{
			Key: prompt.ControlD,
			Fn: func(b *prompt.Buffer) {
				fmt.Println("\nExiting...")
				os.Exit(0)
			},
		}))
	p.Run()
}
