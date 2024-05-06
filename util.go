package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
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
	fmt.Println("Repomigrate - tool for editing Github Repository Secrets/Variables")
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
	repo_name := flag.String("repo", "", "GitHub Repo name")
	owner_username := flag.String("username", "", "GitHub Repo owner")
	flag.Parse()

	if *help {
		printUsage()
		os.Exit(0)
	}

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
