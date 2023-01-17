package main

import (
	"archive/zip"
	_ "embed"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
)

//go:embed liblava.zip
var liblava []byte

//go:embed PresentMon.zip
var PresentMon []byte

var (
	Temp           string
	LavaPath       string
	PresentMonPath string
)

func Cleanup() {
	os.RemoveAll(LavaPath)
	os.Remove(PresentMonPath)

	if err := os.MkdirAll(filepath.Join(os.Getenv("AppData"), "liblava", "lava triangle"), os.ModePerm); err != nil {
		log.Println(err)
	}

	if err := os.RemoveAll(tempFolder); err != nil {
		log.Println(err)
	}
}

func init() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		Cleanup()
		os.Exit(0)
	}()
}

func UnzipFiles() {
	Temp = os.TempDir()

	err := os.WriteFile(filepath.Join(Temp, "liblava.zip"), liblava, 0644)
	if err != nil {
		log.Fatal(err)
	}
	err = unzipSource(filepath.Join(Temp, "liblava.zip"), Temp)
	if err != nil {
		log.Fatal(err)
	}
	os.Remove(filepath.Join(Temp, "liblava.zip"))
	LavaPath = filepath.Join(Temp, "liblava", "lava-triangle.exe")
	if !FileExists(LavaPath) {
		panic(LavaPath)
	}

	err = os.WriteFile(Temp+"\\PresentMon.zip", PresentMon, 0644)
	if err != nil {
		log.Fatal(err)
	}
	err = unzipSource(Temp+"\\PresentMon.zip", Temp)
	if err != nil {
		log.Fatal(err)
	}
	os.Remove(Temp + "\\PresentMon.zip")
	PresentMonPath = filepath.Join(Temp, "PresentMon-1.7.0-x64.exe")
	if !FileExists(PresentMonPath) {
		panic(PresentMonPath)
	}
}

func unzipSource(source, destination string) error {
	// 1. Open the zip file
	reader, err := zip.OpenReader(source)
	if err != nil {
		return err
	}
	defer reader.Close()

	// 2. Get the absolute destination path
	destination, err = filepath.Abs(destination)
	if err != nil {
		return err
	}

	// 3. Iterate over zip files inside the archive and unzip each of them
	for _, f := range reader.File {
		err := unzipFile(f, destination)
		if err != nil {
			return err
		}
	}

	return nil
}

func unzipFile(f *zip.File, destination string) error {
	// 4. Check if file paths are not vulnerable to Zip Slip
	filePath := filepath.Join(destination, f.Name)
	if !strings.HasPrefix(filePath, filepath.Clean(destination)+string(os.PathSeparator)) {
		return fmt.Errorf("invalid file path: %s", filePath)
	}

	// 5. Create directory tree
	if f.FileInfo().IsDir() {
		if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
			return err
		}
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
		return err
	}

	// 6. Create a destination file for unzipped content
	destinationFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	// 7. Unzip the content of a file and copy it to the destination file
	zippedFile, err := f.Open()
	if err != nil {
		return err
	}
	defer zippedFile.Close()

	if _, err := io.Copy(destinationFile, zippedFile); err != nil {
		return err
	}
	return nil
}
