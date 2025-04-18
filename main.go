package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	var username, password, libraryPath, librarySubpath, url string
	var background, insecure bool

	flag.StringVar(&username, "username", os.Getenv("PUTIO_USER"), "Put.io Username")
	flag.StringVar(&username, "u", os.Getenv("PUTIO_USER"), "Put.io Username (short)")
	flag.StringVar(&password, "password", os.Getenv("PUTIO_PASS"), "Put.io Password")
	flag.StringVar(&password, "p", os.Getenv("PUTIO_PASS"), "Put.io Password (short)")
	flag.StringVar(&libraryPath, "library_path", os.Getenv("PUTIO_LIBRARY_PATH"), "Target Root Directory (i.e. /mnt/Plex)")
	flag.StringVar(&libraryPath, "l", os.Getenv("PUTIO_LIBRARY_PATH"), "Target Root Directory (short)")
	flag.StringVar(&librarySubpath, "library_subpath", os.Getenv("PUTIO_LIBRARY_SUBPATH"), "Target Subdirectory (i.e. TV or Music)")
	flag.StringVar(&librarySubpath, "s", os.Getenv("PUTIO_LIBRARY_SUBPATH"), "Target Subdirectory (short)")
	flag.StringVar(&url, "url", "", "Put.io Zip URL")
	flag.StringVar(&url, "z", "", "Put.io Zip URL (short)")
	flag.BoolVar(&background, "background", false, "Run download and extraction in the background and log output")
	flag.BoolVar(&insecure, "insecure", false, "Skip TLS certificate verification")

	flag.Parse()

	if url == "" {
		fmt.Println("URL is required")
		os.Exit(1)
	}

	if username == "" || password == "" || libraryPath == "" || librarySubpath == "" {
		if username == "" || password == "" || libraryPath == "" {
			username, password, libraryPath, librarySubpath = readCredentials()
		}
	}

	// Configure HTTP client
	httpClient := &http.Client{}
	if insecure {
		// Skip TLS verification
		httpClient.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}

	// Check if already running in background
	if os.Getenv("PUTIO_BACKGROUND") == "1" {
		filePath, err := downloadWithClient(httpClient, url, username, password, libraryPath, librarySubpath)
		if err != nil {
			fmt.Println("Error downloading file:", err)
			os.Exit(1)
		}

		err = extract(filePath, libraryPath, librarySubpath)
		if err != nil {
			fmt.Println("Error extracting file:", err)
			os.Exit(1)
		}

		fmt.Println("Download and extraction completed successfully.")
		return
	}

	if background {
		logFile, err := os.OpenFile("download.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println("Error opening log file:", err)
			os.Exit(1)
		}
		defer logFile.Close()

		cmd := exec.Command(os.Args[0], os.Args[1:]...)
		cmd.Stdout = logFile
		cmd.Stderr = logFile
		cmd.Env = append(os.Environ(), "PUTIO_BACKGROUND=1") // Set environment variable
		cmd.Dir = libraryPath  // Set the working directory

		// Detach the process
		cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

		err = cmd.Start()
		if err != nil {
			fmt.Println("Error starting background process:", err)
			os.Exit(1)
		}

		fmt.Printf("Process running in background with PID %d\n", cmd.Process.Pid)
		os.Exit(0)
	}

	filePath, err := downloadWithClient(httpClient, url, username, password, libraryPath, librarySubpath)
	if err != nil {
		fmt.Println("Error downloading file:", err)
		os.Exit(1)
	}

	err = extract(filePath, libraryPath, librarySubpath)
	if err != nil {
		fmt.Println("Error extracting file:", err)
		os.Exit(1)
	}

	fmt.Println("Download and extraction completed successfully.")
}
