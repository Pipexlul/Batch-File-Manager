package command_data

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"

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
		for i, argData := range cmd.arguments {
			cmdInfo.WriteString(fmt.Sprintf("\t\t(%d) %s\n", i+1, argData.name))
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
			cmdInfo.WriteString("\n")
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

	userArgs := inpa.ParseInputNew(userInput)

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
		return fmt.Errorf("Missing required arguments for command %s: [%s]", cmdName, strings.Join(missingArgs, ", "))
	}

	executeErr := targetCmd.execute(userArgs)

	return executeErr
}

func CopyFile(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("Not a typical file: %s", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	dest, err := os.Create(dst)
	if err != nil {
		return 0, err
	}

	defer func() {
		if err := dest.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	return io.Copy(dest, source)
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
	}
	copyAllCmd.execute = func(args map[string]string) error {
		srcFolder, _ := getBaseArgNameOrAliasValue(args, &copyAllCmd.arguments[0])
		dstFolder, _ := getBaseArgNameOrAliasValue(args, &copyAllCmd.arguments[1])

		srcFolderStr := strings.TrimSuffix(srcFolder, "/")
		dstFolderStr := strings.TrimSuffix(dstFolder, "/")

		srcFS := os.DirFS(srcFolderStr)
		dstFS := os.DirFS(dstFolderStr)

		filesToCopy := make([]string, 0)

		srcFiles, err := fs.ReadDir(srcFS, ".")
		if err != nil {
			return err
		}

		for _, srcFile := range srcFiles {
			if !srcFile.IsDir() {
				filesToCopy = append(filesToCopy, fmt.Sprintf("%s/%s", srcFolderStr, srcFile.Name()))
			}
		}

		if len(filesToCopy) == 0 {
			return fmt.Errorf("No files to copy")
		}

		filesToBeOverwritten := make([]string, 0, len(filesToCopy))

		dstFiles, err := fs.ReadDir(dstFS, ".")
		if err != nil {
			return err
		}

		for _, dstFile := range dstFiles {
			if !dstFile.IsDir() {
				for _, fileToCopy := range filesToCopy {
					if strings.HasSuffix(fileToCopy, dstFile.Name()) {
						filesToBeOverwritten = append(filesToBeOverwritten, dstFile.Name())
					}
				}
			}
		}

		if len(filesToBeOverwritten) > 0 {
			fmt.Println("The following files were found in destination folder and will be overwritten:")
			for _, fileToBeOverwritten := range filesToBeOverwritten {
				fmt.Println(fileToBeOverwritten)
			}

			var userResponse string
			fmt.Println("Do you want to continue? (y/n)")
			fmt.Scanln(&userResponse)

			if strings.ToLower(strings.TrimSpace(userResponse)) != "y" {
				return fmt.Errorf("User aborted copyAll command")
			}
		}

		fmt.Println("The following files will be copied:")
		for _, fileToCopy := range filesToCopy {
			fmt.Println(fileToCopy)
		}
		fmt.Printf("Into folder: %s\n", dstFolderStr)

		var userResponse string
		fmt.Println("Do you want to continue? (y/n)")
		fmt.Scanln(&userResponse)

		if strings.ToLower(strings.TrimSpace(userResponse)) != "y" {
			return fmt.Errorf("User aborted copyAll command")
		}

		var wg sync.WaitGroup
		completedWithErrors := false

		for _, fileToCopy := range filesToCopy {
			wg.Add(1)
			go func(fileToCopy string) {
				defer wg.Done()

				_, fileName := filepath.Split(fileToCopy)

				_, err := CopyFile(fileToCopy, fmt.Sprintf("%s/%s", dstFolderStr, fileName))
				if err != nil {
					fmt.Println(err)
					completedWithErrors = true
				}
			}(fileToCopy)
		}

		wg.Wait()

		if completedWithErrors {
			return fmt.Errorf("copyAll command completed with errors")
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
