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

	t = "nmap " + t

	parseCommand(t)

	t = checkSudo(t)
	parts := strings.Fields(t)
	var cmd *exec.Cmd
	if len(parts) == 0 {
		fmt.Println("No command entered")
		return
	} else if len(parts) == 1 && parts[0] == "nmap" {
		cmd = exec.Command("nmap", "-h")
	} else {
		cmd = exec.Command(parts[0], parts[1:]...)
	}

	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(output))
}

func main() {
	fmt.Println("NMAP Interactive CLI")
	fmt.Println("Please type `nmap` and press tab for options. 'Ctrl-D' to exit the program.")

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
