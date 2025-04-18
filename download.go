package main

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gregdel/pushover"
	"github.com/schollz/progressbar/v3"
)

// downloadWithClient handles the downloading process using a custom HTTP client
func downloadWithClient(client *http.Client, url string, username string, password string, libraryPath string, librarySubpath string) (string, error) {
	dlDir := filepath.Join(libraryPath, librarySubpath)
	if _, err := os.Stat(dlDir); os.IsNotExist(err) {
		os.MkdirAll(dlDir, os.ModePerm)
	}

	fileName := filepath.Base(url)
	filePath := filepath.Join(dlDir, fileName)

	fmt.Printf("Downloading %s to %s\n", url, filePath)

	start := time.Now()
	err := downloadFileWithClient(client, url, filePath)
	if err != nil {
		return "", err
	}
	duration := time.Since(start)
	fmt.Printf("Download completed in %s\n", duration)

	if os.Getenv("PUTIO_NOTIFY") == "1" {
		sendPushoverNotification("Download completed", fmt.Sprintf("File %s downloaded in %s", fileName, duration))
	}

	return filePath, nil
}

// downloadFileWithClient downloads a file from the given URL using a custom HTTP client
func downloadFileWithClient(client *http.Client, url string, filepath string) error {
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	bar := progressbar.DefaultBytes(
		resp.ContentLength,
		"downloading",
	)

	_, err = io.Copy(io.MultiWriter(out, bar), resp.Body)
	return err
}

// extract handles the extraction process
func extract(filePath string, libraryPath string, librarySubpath string) error {
	destDir := filepath.Join(libraryPath, librarySubpath)
	fmt.Printf("Extracting %s to %s\n", filePath, destDir)

	start := time.Now()
	err := extractZip(filePath, destDir)
	if err != nil {
		return err
	}
	duration := time.Since(start)
	fmt.Printf("Extraction completed in %s\n", duration)

	if os.Getenv("PUTIO_NOTIFY") == "1" {
		sendPushoverNotification("Extraction completed", fmt.Sprintf("File %s extracted in %s", filepath.Base(filePath), duration))
	}

	// Remove the archive if PUTIO_CLEAN is set to 1
	if os.Getenv("PUTIO_CLEAN") == "1" {
		err = os.Remove(filePath)
		if err != nil {
			return fmt.Errorf("failed to remove archive: %v", err)
		}
		fmt.Printf("Removed archive %s\n", filePath)
	}

	return nil
}

// extractZip extracts a zip file to the specified directory
func extractZip(src string, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		fpath := filepath.Join(dest, f.Name)

		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			return err
		}

		_, err = io.Copy(outFile, rc)

		outFile.Close()
		rc.Close()

		if err != nil {
			return err
		}
	}
	return nil
}

// sendPushoverNotification sends a notification using the Pushover API
func sendPushoverNotification(title, message string) {
	appToken := os.Getenv("PUSHOVER_TOKEN")
	userKey := os.Getenv("PUSHOVER_USER")

	app := pushover.New(appToken)
	recipient := pushover.NewRecipient(userKey)

	msg := pushover.NewMessageWithTitle(message, title)
	_, err := app.SendMessage(msg, recipient)
	if err != nil {
		fmt.Printf("Failed to send Pushover notification: %v\n", err)
	} else {
		fmt.Println("Pushover notification sent successfully")
	}
}
