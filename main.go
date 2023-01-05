package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	cmdd "github.com/Pipexlul/batch-file-manager/command_data"
	td "github.com/Pipexlul/batch-file-manager/translation_data"
)

func greetings() {
	fmt.Println("Welcome to Pipexlul's batch file manager!")
	fmt.Println("This is a simple program that allows you to copy/rename/delete files with ease.")
	fmt.Println("It's meant to be my first take on the Go world, so it's not perfect, but it's a start.")

	fmt.Println("--------------------")

	fmt.Println("Bienvenido al programa de gestión de archivos de Pipexlul!")
	fmt.Println("Este es un programa simple que te permite copiar/renombrar/eliminar archivos con facilidad.")
	fmt.Println("Está hecho para ser mi primer intento en el mundo de Go, así que no es perfecto, pero es un comienzo.")

	fmt.Println("--------------------")

	fmt.Println("Please type the language you want to use (en/es):")
	fmt.Println("Por favor, escribe el idioma que quieres usar (en/es):")
}

func main() {
	greetings()

	var lang string
	var translatedToken string
	var err error

	for {
		fmt.Scanln(&lang)

		lang = strings.ToLower(strings.TrimSpace(lang))

		if lang == "en" || lang == "es" {
			break
		} else {
			fmt.Println("Please type a valid language (en/es):")
			fmt.Println("Por favor, escribe un idioma válido (en/es):")
		}
	}

	// Display chosen language
	if translatedToken, err = td.GetTranslatedString(lang, "langChosen"); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(translatedToken)
	}

	fmt.Println("Enter a command and it's arguments if necessary. To get a list of all commands use the command help.")
	fmt.Println("If you want detailed information about a single command you can do as follows:")
	fmt.Println("help -command copyAll")
	fmt.Println("If you'd like to exit the program, just type exit")

	reader := bufio.NewReader(os.Stdin)

	for {
		userInput, err := reader.ReadString('\n')

		if err != nil {
			fmt.Println(err)
			continue
		}

		userInput = strings.TrimRight(userInput, "\n")
		userInput = strings.TrimRight(userInput, "\r")

		if userInput == "exit" {
			os.Exit(0)
		}

		executeErr := cmdd.TryExecuteCommand(userInput)

		if executeErr != nil {
			fmt.Printf("[Error]: %v\n", executeErr)
		}
	}
}
