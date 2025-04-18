package main

import (
	"archive/zip"
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

// createDummyZip creates a dummy zip file for testing
func createDummyZip() ([]byte, error) {
	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)

	files := []struct {
		Name, Body string
	}{
		{"file1.txt", "This is the content of file1.txt"},
		{"file2.txt", "This is the content of file2.txt"},
	}

	for _, file := range files {
		f, err := zipWriter.Create(file.Name)
		if err != nil {
			return nil, err
		}
		_, err = f.Write([]byte(file.Body))
		if err != nil {
			return nil, err
		}
	}

	err := zipWriter.Close()
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func TestDownloadAndExtract(t *testing.T) {
	// Create a dummy zip file
	zipData, err := createDummyZip()
	if err != nil {
		t.Fatalf("Failed to create dummy zip: %v", err)
	}

	// Set up a local HTTP server to serve the dummy zip file
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/zip")
		w.Write(zipData)
	}))
	defer ts.Close()

	// Set up test directories
	libraryPath := "./test_library"
	librarySubpath := "test_subpath"
	os.MkdirAll(filepath.Join(libraryPath, librarySubpath), os.ModePerm)
	defer os.RemoveAll(libraryPath)

	// Set up HTTP client
	httpClient := &http.Client{}

	// Run the download and extract functions
	filePath, err := downloadWithClient(httpClient, ts.URL, "", "", libraryPath, librarySubpath)
	if err != nil {
		t.Fatalf("Download failed: %v", err)
	}

	err = extract(filePath, libraryPath, librarySubpath)
	if err != nil {
		t.Fatalf("Extraction failed: %v", err)
	}

	// Verify the extracted files
	for _, fileName := range []string{"file1.txt", "file2.txt"} {
		filePath := filepath.Join(libraryPath, librarySubpath, fileName)
		content, err := ioutil.ReadFile(filePath)
		if err != nil {
			t.Fatalf("Failed to read extracted file %s: %v", fileName, err)
		}

		expectedContent := "This is the content of " + fileName
		if string(content) != expectedContent {
			t.Errorf("Content mismatch for %s: got %s, want %s", fileName, content, expectedContent)
		}
	}
}

func TestBackgroundFlag(t *testing.T) {
	// Create a dummy zip file
	zipData, err := createDummyZip()
	if err != nil {
		t.Fatalf("Failed to create dummy zip: %v", err)
	}

	// Set up a local HTTP server to serve the dummy zip file
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/zip")
		w.Write(zipData)
	}))
	defer ts.Close()

	// Set up test directories
	libraryPath := "./test_library"
	librarySubpath := "test_subpath"
	os.MkdirAll(filepath.Join(libraryPath, librarySubpath), os.ModePerm)
	defer os.RemoveAll(libraryPath)

	// Build the main application
	buildCmd := exec.Command("go", "build", "-o", "putio-cli-v2")
	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr
	err = buildCmd.Run()
	if err != nil {
		t.Fatalf("Failed to build application: %v", err)
	}
	defer os.Remove("putio-cli-v2") // Remove the binary after the test

	// Run the main application with the background flag
	cmd := exec.Command("./putio-cli-v2", "--background", "--url", ts.URL, "--library_path", libraryPath, "-s", librarySubpath)
	cmd.Env = append(os.Environ(), "PUTIO_BACKGROUND=1") // Ensure environment is set
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Start()
	if err != nil {
		t.Fatalf("Failed to start background process: %v", err)
	}

	// Wait for a longer period to allow the background process to complete
	time.Sleep(30 * time.Second)

	// Verify the extracted files
	for _, fileName := range []string{"file1.txt", "file2.txt"} {
		filePath := filepath.Join(libraryPath, librarySubpath, fileName)
		content, err := ioutil.ReadFile(filePath)
		if err != nil {
			t.Fatalf("Failed to read extracted file %s: %v", fileName, err)
		}

		expectedContent := "This is the content of " + fileName
		if string(content) != expectedContent {
			t.Errorf("Content mismatch for %s: got %s, want %s", fileName, content, expectedContent)
		}
	}
}

func TestInsecureFlag(t *testing.T) {
	// Create a dummy zip file
	zipData, err := createDummyZip()
	if err != nil {
		t.Fatalf("Failed to create dummy zip: %v", err)
	}

	// Set up a local HTTP server to serve the dummy zip file
	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/zip")
		w.Write(zipData)
	}))
	defer ts.Close()

	// Set up test directories
	libraryPath := "./test_library"
	librarySubpath := "test_subpath"
	os.MkdirAll(filepath.Join(libraryPath, librarySubpath), os.ModePerm)
	defer os.RemoveAll(libraryPath)

	// Build the main application
	buildCmd := exec.Command("go", "build", "-o", "putio-cli-v2")
	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr
	err = buildCmd.Run()
	if err != nil {
		t.Fatalf("Failed to build application: %v", err)
	}
	defer os.Remove("putio-cli-v2") // Remove the binary after the test

	// Run the main application with the insecure flag
	cmd := exec.Command("./putio-cli-v2", "--insecure", "--url", ts.URL, "--library_path", libraryPath, "-s", librarySubpath)
	cmd.Env = append(os.Environ(), "PUTIO_BACKGROUND=1") // Ensure environment is set
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		t.Fatalf("Failed to run process with insecure flag: %v", err)
	}

	// Verify the extracted files
	for _, fileName := range []string{"file1.txt", "file2.txt"} {
		filePath := filepath.Join(libraryPath, librarySubpath, fileName)
		content, err := ioutil.ReadFile(filePath)
		if err != nil {
			t.Fatalf("Failed to read extracted file %s: %v", fileName, err)
		}

		expectedContent := "This is the content of " + fileName
		if string(content) != expectedContent {
			t.Errorf("Content mismatch for %s: got %s, want %s", fileName, content, expectedContent)
		}
	}
}
