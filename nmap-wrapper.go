package main

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/c-bata/go-prompt"
)

func completer(d prompt.Document) []prompt.Suggest {
	return CustomFilter(Suggestions, d.GetWordBeforeCursor())
}

func checkSudo(cmd string) string {
	for _, v := range SudoRequiredFlags {
		if strings.Contains(cmd, v) {
			fmt.Printf("The flag %s requires sudo, you may need to enter your password\n", v)
			return "sudo " + cmd
		}
	}
	return cmd
}

func main() {
	fmt.Println("Please type `nmap` and press tab for options. 'exit' or 'quit' to exit the program.")
	for {
		t := prompt.Input("> ", completer)

		if strings.ToLower(t) == "exit" || strings.ToLower(t) == "quit" || strings.ToLower(t) == "q" {
			fmt.Println("Exiting program...")
			break
		}

		t = checkSudo(t)
		parts := strings.Fields(t)
		var cmd *exec.Cmd
		if len(parts) == 0 {
			fmt.Println("No command entered")
			continue
		} else if len(parts) == 1 && parts[0] == "nmap" {
			cmd = exec.Command("nmap", "-h")
		} else {
			cmd = exec.Command(parts[0], parts[1:]...)
		}

		output, err := cmd.Output()
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}
		fmt.Println(string(output))
	}
}
