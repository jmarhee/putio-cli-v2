package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// readCredentials reads user credentials from environment variables or prompts the user
func readCredentials() (string, string, string, string) {
	var username, password, libraryPath, librarySubpath string

	if os.Getenv("PUTIO_USER") == "" {
		fmt.Print("Enter your Put.io Username: ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		username = scanner.Text()
	} else {
		username = os.Getenv("PUTIO_USER")
	}

	if os.Getenv("PUTIO_PASS") == "" {
		fmt.Print("Enter your Put.io Password: ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		password = scanner.Text()
	} else {
		password = os.Getenv("PUTIO_PASS")
	}

	if os.Getenv("PUTIO_LIBRARY_PATH") == "" {
		fmt.Print("Enter the Plex Library path: ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		libraryPath = scanner.Text()
	} else {
		libraryPath = os.Getenv("PUTIO_LIBRARY_PATH")
	}

	if os.Getenv("PUTIO_LIBRARY_SUBPATH") == "" {
		fmt.Print("Enter the subdirectory to download and unpack to: ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		librarySubpath = scanner.Text()
	} else {
		librarySubpath = os.Getenv("PUTIO_LIBRARY_SUBPATH")
	}

	return username, password, libraryPath, librarySubpath
}

// readConfig reads configuration from a file
func readConfig(configPath string) (string, string, string, string) {
	file, err := os.Open(configPath)
	if err != nil {
		fmt.Println("Error opening config file:", err)
		return "", "", "", ""
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var username, password, libraryPath, librarySubpath string

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "username=") {
			username = strings.TrimPrefix(line, "username=")
		} else if strings.HasPrefix(line, "password=") {
			password = strings.TrimPrefix(line, "password=")
		} else if strings.HasPrefix(line, "library_path=") {
			libraryPath = strings.TrimPrefix(line, "library_path=")
		} else if strings.HasPrefix(line, "library_subpath=") {
			librarySubpath = strings.TrimPrefix(line, "library_subpath=")
		}
	}

	return username, password, libraryPath, librarySubpath
}
