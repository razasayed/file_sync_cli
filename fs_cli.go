package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func main() {
	// Read the flags from the sync subcommand
	syncCmd := flag.NewFlagSet("sync", flag.ExitOnError)
	sourceDirPtr := syncCmd.String("s", "", "source directory")
	destinationDirPtr := syncCmd.String("d", "", "destination directory")
	if len(os.Args) < 2 {
		fmt.Println("expected 'sync' subcommand")
		os.Exit(1)
	}
	syncCmd.Parse(os.Args[2:])

	// Sync source directory to destination directory
	if err := sync(*sourceDirPtr, *destinationDirPtr); err != nil {
		fmt.Println("Err:", err)
	}
}

func sync(sourceDirectoryPath, destinationDirectoryPath string) error {
	// Make paths absolute
	sourceDir, err := filepath.Abs(sourceDirectoryPath)
	if err != nil {
		return err
	}

	destinationDir, err := filepath.Abs(destinationDirectoryPath)
	if err != nil {
		return err
	}

	// Check that source path exists and is a directory
	info, err := os.Stat(sourceDir)
	if err != nil || !info.IsDir() {
		return err
	}

	// Create the destination directory if needed
	if err = os.MkdirAll(destinationDir, info.Mode()); err != nil {
		return err
	}

	// Delete the items from destination directory that do not exist in source directory
	err = filepath.Walk(destinationDir, func(destinationItemPath string, fileinfo os.FileInfo, err error) error {

		destinationItemName := strings.TrimPrefix(destinationItemPath, destinationDir)
		sourceItemPath := filepath.Join(sourceDir, destinationItemName)

		if _, err := os.Stat(sourceItemPath); err != nil {
			if os.IsNotExist(err) {
				if err := os.RemoveAll(destinationItemPath); err != nil {
					return err
				}
			} else {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	// Copy new items into destination directory
	err = filepath.Walk(sourceDir, func(sourceItemPath string, fileinfo os.FileInfo, err error) error {

		sourceItemName := strings.TrimPrefix(sourceItemPath, sourceDir)
		destinationItemPath := filepath.Join(destinationDir, sourceItemName)

		if _, err := os.Stat(destinationItemPath); err != nil {
			if os.IsNotExist(err) {
				if fileinfo.IsDir() {
					if err = copyDirectory(sourceItemPath, destinationItemPath); err != nil {
						return err
					}
				} else {
					if err = copyFile(sourceItemPath, destinationItemPath); err != nil {
						return err
					}
				}
			} else {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func copyFile(sourcePath, destinationPath string) error {
	// Open the original file
	originalFile, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer originalFile.Close()

	// Create the new file
	newFile, err := os.Create(destinationPath)
	if err != nil {
		return err
	}
	defer newFile.Close()

	// Copy
	_, err = io.Copy(newFile, originalFile)
	if err != nil {
		return err
	}

	return nil
}

func copyDirectory(sourcePath, destinationPath string) error {
	var err error
	var sourceInfo os.FileInfo
	var fileInfos []os.FileInfo

	if sourceInfo, err = os.Stat(sourcePath); err != nil {
		return err
	}

	if err = os.MkdirAll(destinationPath, sourceInfo.Mode()); err != nil {
		return err
	}

	if fileInfos, err = ioutil.ReadDir(sourcePath); err != nil {
		return err
	}

	for _, info := range fileInfos {
		sourcePath := path.Join(sourcePath, info.Name())
		destinationPath := path.Join(destinationPath, info.Name())

		if info.IsDir() {
			if err = copyDirectory(sourcePath, destinationPath); err != nil {
				return err
			}
		} else {
			if err = copyFile(sourcePath, destinationPath); err != nil {
				return err
			}
		}
	}

	return nil
}
