package main

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"github.com/c-bata/go-prompt"
)

func completer(d prompt.Document) []prompt.Suggest {
	word := strings.Trim(d.GetWordBeforeCursor(), " ")
	switch word {
	case "--script":
		return ScriptCategories
	default:
		return CustomFilter(Suggestions, word)
	}
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
	fmt.Println("Please type `nmap` and press tab for options. 'Ctrl-D' to exit the program.")
	fmt.Println("Press Ctrl-T to pick a target")
	var target string
	p := prompt.New(
		func(input string) {
			execute(input)
		},
		completer,
		prompt.OptionPrefix("> "),
		prompt.OptionTitle("My CLI"),
		prompt.OptionSuggestionBGColor(prompt.DarkGray),
		prompt.OptionPrefixTextColor(prompt.Green),
		prompt.OptionDescriptionBGColor(prompt.DarkGray),
		prompt.OptionSelectedDescriptionBGColor(prompt.DarkGray),
		prompt.OptionPreviewSuggestionBGColor(prompt.DarkGray),
		prompt.OptionInitialBufferText("nmap "),
		prompt.OptionAddKeyBind(prompt.KeyBind{
			Key: prompt.ControlT,
			Fn: func(b *prompt.Buffer) {
				b.InsertText(target, false, true)
			},
		}),
		prompt.OptionAddKeyBind(prompt.KeyBind{
			Key: prompt.ControlS,
			Fn: func(b *prompt.Buffer) {
				_, err := fmt.Scanln(&target)
				if err != nil {
					fmt.Println("Error reading input:", err)
					return
				}
				fmt.Println("You saved:", target)
			},
		}),
	)
	p.Run()
}
