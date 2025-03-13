package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	var username, password, libraryPath, librarySubpath, url string

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

	filePath, err := download(url, username, password, libraryPath, librarySubpath)
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
