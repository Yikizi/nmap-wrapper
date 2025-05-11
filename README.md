# Nmap Wrapper

## Problem
Currently nmap is difficult to use. There are 115 different flags and those flags take also arguments.

Also there are over 3000 lines in the manpage, which makes it hard to navigate.

## Idea
We plan to make a nmap CLI tool, that will solve these problems.

**HOW????** - you may ask

By providing:
- [x] Interactive interface
- [x] Tab completion
- [x] Inline documentation
- [x] Smart suggestions
- [x] Command saving
- [x] History

## Interactive Features

This tool enhances the Nmap experience with several interactive features.

### Command Saving (`set` and `list`)

You can save frequently used Nmap command snippets or full commands using the `set` command. This allows you to assign a short name to a longer command or set of flags and arguments, making it easier to reuse complex configurations.

**How to save a command:**

Use the syntax `set <name> <command>`.
For example, to save a common port scan:
```/dev/null/example.txt
set common_ports -p22,80,443,8080
```
Now you can use `common_ports` within your nmap commands.

**How to use a saved command:**

Simply include the saved name in your command line. The tool will automatically substitute the name with the saved command string before execution.
For example, to run a scan using the saved `common_ports` setting on a target:
```/dev/null/example.txt
common_ports scanme.nmap.org
```
This will execute `nmap -p22,80,443,8080 scanme.nmap.org`.

You can even combine multiple saved commands or use them with regular Nmap flags:
```/dev/null/example.txt
set aggressive_scan -A
aggressive_scan common_ports 192.168.1.1
```
This would execute `nmap -A -p22,80,443,8080 192.168.1.1`.

**How to list saved commands:**

Type `list` and press Enter to display all your saved command sets.

```/dev/null/example.txt
list
```

### Autocompletion and Smart Suggestions

The CLI tool provides intelligent autocompletion powered by `go-prompt`. Pressing the `Tab` key will show you available Nmap flags and their descriptions.

**Features:**

-   **Basic Tab Completion:** Provides a list of standard Nmap flags and script categories when you press `Tab`.
-   **Inline Documentation:** Each suggestion includes a brief description to help you understand what the flag or option does.
-   **Smart Suggestions (Last Used Values):** The tool remembers the last value you used for certain flags (like `-p`, `-iL`, `--script-args`, etc.). When you type one of these flags again, the suggestion will automatically include the last used value, often highlighted with `[last used: value]`. This speeds up repeated scans with similar parameters.
-   **Saved Command Suggestions:** Your saved command names will also appear in the suggestions list, allowing you to quickly select them using `Tab`.

To use autocompletion, just start typing a flag or command name and press `Tab`. Navigate the suggestions using the arrow keys and press `Enter` to select one.

## Proof of concept
- Look how far we've come!
![demo](output.gif)

## Running the project
```
go run *.go
```

## Used technologies
- [Nmap](https://nmap.org/) (obviously)
- [go-prompt](https://github.com/c-bata/go-prompt) (for creating CLI tool)

## Contibutors
- [Yikizi](https://github.com/Yikizi)
- [RasmusRaasuke](https://github.com/RasmusRaasuke)
