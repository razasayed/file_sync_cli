package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func createEmptyFile(name string) {
	d := []byte("")
	check(os.WriteFile(name, d, 0644))
}

func setup() {
	err := os.Mkdir("unit_test_data", 0755)
	check(err)
	err = os.Mkdir("unit_test_data/source_folder_a", 0755)
	check(err)
	createEmptyFile("unit_test_data/source_folder_a/file_a")
	createEmptyFile("unit_test_data/source_folder_a/file_b")
	createEmptyFile("unit_test_data/source_folder_a/file_c")
	err = os.Mkdir("unit_test_data/source_folder_b", 0755)
	check(err)
	createEmptyFile("unit_test_data/source_folder_b/file_a")
	createEmptyFile("unit_test_data/source_folder_b/file_c")
	createEmptyFile("unit_test_data/source_folder_b/file_d")
	err = os.Mkdir("unit_test_data/source_folder_c", 0755)
	check(err)
	createEmptyFile("unit_test_data/source_folder_c/file_a")
	createEmptyFile("unit_test_data/source_folder_c/file_d")
	createEmptyFile("unit_test_data/source_folder_c/file_e")
	err = os.Mkdir("unit_test_data/source_folder_c/dir_a", 0755)
	check(err)
	createEmptyFile("unit_test_data/source_folder_c/dir_a/file_a_a")
	err = os.Mkdir("unit_test_data/source_folder_d", 0755)
	check(err)
	createEmptyFile("unit_test_data/source_folder_d/file_a")
	err = os.Mkdir("unit_test_data/target_folder", 0755)
	check(err)
}

func teardown() {
	os.RemoveAll("unit_test_data")
}

func getFilesInDirectory(directoryPath string) []string {
	result := []string{}
	_ = filepath.Walk(directoryPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		result = append(result, path)
		return nil
	})
	return result[1:]
}

func TestDirectorySync(t *testing.T) {
	sourceDirectoryPaths := []string{
		"unit_test_data/source_folder_a",
		"unit_test_data/source_folder_b",
		"unit_test_data/source_folder_c",
		"unit_test_data/source_folder_d",
	}
	destinationDirectoryPath := "unit_test_data/target_folder"

	for _, sourceDirectoryPath := range sourceDirectoryPaths {
		filesInSourceDirectory := getFilesInDirectory(sourceDirectoryPath)
		sync(sourceDirectoryPath, destinationDirectoryPath)
		filesInDestinationDirectory := getFilesInDirectory(destinationDirectoryPath)
		// Check that the counts match
		if len(filesInSourceDirectory) != len(filesInDestinationDirectory) {
			t.Errorf("Sync failed for source directory %s", sourceDirectoryPath)
			break
		}
		// Check that the file names match
		names_match := true
		for i := 0; i < len(filesInSourceDirectory); i++ {
			source_file_name := strings.TrimPrefix(filesInSourceDirectory[i], sourceDirectoryPath)
			destination_file_name := strings.TrimPrefix(filesInDestinationDirectory[i], destinationDirectoryPath)
			if source_file_name != destination_file_name {
				names_match = false
				break
			}
		}
		if !names_match {
			t.Errorf("Sync failed for source directory %s", sourceDirectoryPath)
			break
		}
	}
}

func TestCopyFile(t *testing.T) {
	sourceDirectoryPath := "unit_test_data/source_folder_a"
	sourceFilePath := sourceDirectoryPath + "/file_a"

	destinationDirectoryPath := "unit_test_data/target_folder"
	destinationFilePath := destinationDirectoryPath + "/file_a"

	copyFile(sourceFilePath, destinationFilePath)

	filesInSourceDirectory := getFilesInDirectory(sourceDirectoryPath)
	sourceFileName := strings.TrimPrefix(filesInSourceDirectory[0], sourceDirectoryPath)

	filesInDestinationDirectory := getFilesInDirectory(destinationDirectoryPath)
	destinationFileName := strings.TrimPrefix(filesInDestinationDirectory[0], destinationDirectoryPath)

	if sourceFileName != destinationFileName {
		t.Errorf("File copy failed for file %s to %s", sourceFilePath, destinationFilePath)
	}
}

func TestCopyDirectory(t *testing.T) {
	sourceDirectoryPath := "unit_test_data/source_folder_c/dir_a"
	destinationDirectoryPath := "unit_test_data/target_folder/dir_a"

	copyDirectory(sourceDirectoryPath, destinationDirectoryPath)

	filesInSourceDirectory := getFilesInDirectory(sourceDirectoryPath)
	filesInDestinationDirectory := getFilesInDirectory(destinationDirectoryPath)

	// Check that the counts match
	if len(filesInSourceDirectory) != len(filesInDestinationDirectory) {
		t.Errorf("Directory copy failed for %s to %s", sourceDirectoryPath, destinationDirectoryPath)
	}

	// Check that the names match
	for i := 0; i < len(filesInSourceDirectory); i++ {
		sourceName := strings.TrimPrefix(filesInSourceDirectory[i], sourceDirectoryPath)
		destinationName := strings.TrimPrefix(filesInDestinationDirectory[i], destinationDirectoryPath)

		if sourceName != destinationName {
			t.Errorf("Directory copy failed for %s to %s", sourceDirectoryPath, destinationDirectoryPath)
			break
		}
	}
}
