package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
)

var (
	VERBOSE     = false
	PRINT       = false
	PRINTMODE   = "normal"
	PRETTY      = false
	CONCISE     = false
	VARSONLY    = false
	SECRETSONLY = false
)

func EditSecretMessages() {
	fmt.Print("Enter new value for secret: ")
}
func EditVariableMessages() {
	fmt.Print("Enter new value for variable: ")
}
func ChoiceMessages() {
	fmt.Println("Which variable do you want to edit?")
	fmt.Println("Enter in a format vN or sN (variable/secret number N)")
	fmt.Print("Your input: ")
}
func printUsage() {
	fmt.Println("Varmigrate - tool for editing Github Repository Secrets/Variables")
	fmt.Println("Options:")
	flag.PrintDefaults()
}

func GetInput(messages func()) string {
	scanner := bufio.NewScanner(os.Stdin)

	messages()

	// Wait for input
	scanned := scanner.Scan()

	// Check if input was received successfully
	if !scanned {
		fmt.Println("Failed to read input:", scanner.Err())
		return "fail"
	}

	// Get the input text
	input := scanner.Text()
	fmt.Println(input)
	return input
}

func parseArgs() (*RepoData, error) {

	help := flag.Bool("help", false, "Print usage information")
	verbose := flag.Bool("verbose", false, "Verbose mode")
	print := flag.Bool("print", false, "Print-only mode")
	printmode := flag.String("mode", "normal", "Print mode (normal/json)")
	concise := flag.Bool("concise", false, "Only print data")
	pretty := flag.Bool("pretty", false, "Pretty-print json")
	varsOnly := flag.Bool("vars-only", false, "Print only repository variables")
	secretsOnly := flag.Bool("secrets-only", false, "Print only repository secrets")

	repo_name := flag.String("repo", "", "GitHub Repo name")
	owner_username := flag.String("username", "", "GitHub Repo owner")
	flag.Parse()

	if *help {
		printUsage()
		os.Exit(0)
	}
	if *verbose {
		VERBOSE = true
	}
	if *print {
		PRINT = true
	}
	if *pretty {
		PRETTY = true
	}
	if *concise {
		CONCISE = true
	}
	if *varsOnly {
		VARSONLY = true
	}
	if *secretsOnly {
		SECRETSONLY = true
	}

	PRINTMODE = *printmode

	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		log.Println()
		return nil, fmt.Errorf("no token provided")
	}
	if *owner_username == "" {
		// log.Println("No repo username provided")
		return nil, fmt.Errorf("no repo username provided")
	}
	if *repo_name == "" {
		return nil, fmt.Errorf("no repo name provided")
	}

	return &RepoData{
		Username:    *owner_username,
		Repo:        *repo_name,
		GithubToken: token,
	}, nil
}
