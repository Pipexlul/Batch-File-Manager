package command_data

import "fmt"

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

func RegisterCommand(c command) {
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

func getBaseCommandOrAliasValue(parsedArgs map[string]string, argToCheck argData) (string, error) {
	if res, found := parsedArgs[argToCheck.name]; found {
		return res, nil
	}

	for _, alias := range argToCheck.aliases {
		if res, found = parsedArgs[alias]; found {
			return res, nil
		}
	}

	return "", fmt.Errorf("Could not find arg of base name %s and neither with alias names: %v", argToCheck.name, argToCheck.aliases)
}

func findCmdDataFromCmdName(cmdName string) (*command, error) {
	for _, cmd := range Commands {
		if cmd.name == cmdName {
			return &cmd, nil
		}
	}

	return nil, fmt.Errorf("Couldn't find any command of base name: %s", cmdName)
}

func findArgDataFromArgName(cmd *command, argName string) (*command, error)

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

	RegisterCommand(copyAllCmd)

	helpCmd := command{
		name:        "help",
		description: "Lists all commands and their descriptions if no arguments are passed otherwise it prints info about a single command",
		arguments: []argData{
			{
				name: "command",
				description: "Command to get information about, if not set, it will print all Commands",
				aliases: []string{
					"cmd",
				},
				required: false,
			}
		},
		helpMsg: `help lists all the possible commands and their data, but if you use the "command" argument and give it the name of another command
		it will print info only about that specific command`,
	}
	helpCmd.execute: func(args map[string]string) error {

	}
}
