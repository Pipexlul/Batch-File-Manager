package command_data

import (
	"fmt"
	"strings"

	inpa "github.com/Pipexlul/batch-file-manager/input_parser"
)

var Commands []command

type argData struct {
	name        string
	description string
	aliases     []string
	required    bool
}

type command struct {
	name        string
	description string
	arguments   []argData
	helpMsg     string

	execute func(args map[string]string) error
}

func registerCommand(c command) {
	Commands = append(Commands, c)
}

func allRequiredArgsPresent(parsedArgs map[string]string, cmdArgs []argData) (bool, []string) {
	var missingArgs []string

	for _, arg := range cmdArgs {
		if arg.required {
			if res, found := parsedArgs[arg.name]; !found || res == "" {
				foundAlias := false

				for _, alias := range arg.aliases {
					if foundAlias {
						break
					}

					res, found = parsedArgs[alias]
					if found && res != "" {
						foundAlias = true
					}
				}

				if !foundAlias {
					missingArgs = append(missingArgs, arg.name)
				}
			}
		}
	}

	return len(missingArgs) == 0, missingArgs
}

func findCmdDataFromCmdName(cmdName string) (*command, error) {
	for _, cmd := range Commands {
		if cmd.name == cmdName {
			return &cmd, nil
		}
	}

	return nil, fmt.Errorf("Couldn't find any command of base name: %s", cmdName)
}

func findArgDataFromArgName(cmd *command, argName string) (*argData, error) {
	for _, comm := range Commands {
		if comm.name == cmd.name {
			for _, argData := range comm.arguments {
				if argData.name == argName {
					return &argData, nil
				}
			}
		}
	}

	return nil, fmt.Errorf("Couldn't retrieve argData of name %s from %s command", argName, cmd.name)
}

func displayCommandInfo(cmd *command, onlyDesc bool) {
	if onlyDesc {
		fmt.Printf("[%s] - [%s]\n", cmd.name, cmd.description)
		return
	}

	var cmdInfo strings.Builder

	cmdInfo.WriteString(fmt.Sprintf("[%s]\n", cmd.name))
	cmdInfo.WriteString(fmt.Sprintf("\t%s\n", cmd.description))
	cmdInfo.WriteString("\tArguments:\n")
	if len(cmd.arguments) == 0 {
		cmdInfo.WriteString("\t\tNo arguments\n")
	} else {
		for _, argData := range cmd.arguments {
			cmdInfo.WriteString(fmt.Sprintf("\t\t%s\n", argData.name))
			cmdInfo.WriteString(fmt.Sprintf("\t\tDescription: %s\n", argData.description))
			cmdInfo.WriteString(fmt.Sprintf("\t\tRequired: %t\n", argData.required))
			if len(argData.aliases) == 0 {
				cmdInfo.WriteString("\t\tNo aliases\n")
			} else {
				cmdInfo.WriteString("\t\tAliases:\n")
				for _, aliasName := range argData.aliases {
					cmdInfo.WriteString(fmt.Sprintf("\t\t\t%s\n", aliasName))
				}
			}
		}
	}

	fmt.Println(cmdInfo.String())
}

func getBaseArgNameOrAliasValue(parsedArgs map[string]string, arg *argData) (string, bool) {
	res, found := parsedArgs[arg.name]

	if found {
		return res, true
	}

	for _, alias := range arg.aliases {
		res, found = parsedArgs[alias]

		if found {
			return res, true
		}
	}

	return "", false
}

func TryExecuteCommand(userInput string) error {
	cmdName := strings.SplitN(userInput, " ", 2)[0]

	userArgs := inpa.ParseInput(userInput)

	var targetCmd *command

	for _, cmd := range Commands {
		if cmd.name == cmdName {
			targetCmd = &cmd
			break
		}
	}

	if targetCmd == nil {
		return fmt.Errorf("Could not find command of name: %s, please check your input.\n", cmdName)
	}

	if success, missingArgs := allRequiredArgsPresent(userArgs, targetCmd.arguments); !success {
		return fmt.Errorf("Missing required arguments for command %s: %v", cmdName, missingArgs)
	}

	executeErr := targetCmd.execute(userArgs)

	return executeErr
}

func init() {
	copyAllCmd := command{
		name:        "copyAll",
		description: "Copy all files from one folder to another",
		arguments: []argData{
			{
				name:        "source",
				description: "Source folder",
				aliases:     []string{"s", "src"},
				required:    true,
			},
			{
				name:        "destination",
				description: "Destination folder",
				aliases:     []string{"d", "dst"},
				required:    true,
			},
		},
		helpMsg: `copyAll copies all the files in the [source] folder path to the [destination] folder path
		Arguments:
			-source
				aliases: -s -src
				purpose: Path of the directory where files will be copied from
				required: Yes
				
			-destination
				aliases: -d -dst
				purpose: Path of the directory where files will be copied to
				required: Yes`,
	}
	copyAllCmd.execute = func(args map[string]string) error {
		if valid, missingArgs := allRequiredArgsPresent(args, copyAllCmd.arguments); !valid {
			return fmt.Errorf("There are missing arguments in your command: %v", missingArgs)
		}

		fmt.Println("Successful copyAll command")
		return nil
	}

	registerCommand(copyAllCmd)

	helpCmd := command{
		name:        "help",
		description: "Lists all commands and their descriptions if no arguments are passed, otherwise it prints info about a single command",
		arguments: []argData{
			{
				name:        "command",
				description: "Command to get information about, if not set, it will print all Commands",
				aliases: []string{
					"cmd",
				},
				required: false,
			},
		},
		helpMsg: `help lists all the possible commands and their data, but if you use the "command" argument and give it the name of another command
		it will print info only about that specific command`,
	}
	helpCmd.execute = func(args map[string]string) error {
		res, found := getBaseArgNameOrAliasValue(args, &helpCmd.arguments[0])

		if found && res != "" {
			var cmd *command
			var err error

			cmd, err = findCmdDataFromCmdName(res)

			if err != nil {
				return fmt.Errorf("No command of name %s found", res)
			}

			displayCommandInfo(cmd, false)
		} else if !found {
			for _, cmd := range Commands {
				displayCommandInfo(&cmd, true)
			}
		}

		return nil
	}

	registerCommand(helpCmd)
}
