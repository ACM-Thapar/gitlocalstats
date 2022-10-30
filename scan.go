package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"strings"
)

func scan(folder string) {
	fmt.Printf("Found folders:\n\n")
	repositories := recursiveScanFolder(folder)
	filepath := getDotFilePath()
	addNewSliceElementsToFile(filepath, repositories)
	fmt.Printf("\n\nSuccessfully added\n\n")
}

func recursiveScanFolder(folder string) []string {
	return scanGitFolders(make([]string, 0), folder)
}

func getDotFilePath() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	dotFile := usr.HomeDir + "/.gitlocalstats"

	return dotFile
}

func addNewSliceElementsToFile(filepath string, newRepos []string) {
	existingRepos := parseFileLinesToSlice(filepath)
	repos := joinSlices(newRepos, existingRepos)
	dumpStringsSliceToFile(repos, filepath)
}

func parseFileLinesToSlice(filepath string) []string {
	f := openFile(filepath)
	defer f.Close()

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		if err != io.EOF {
			panic(err)
		}
	}
	return lines
}

func openFile(filepath string) *os.File {
	f, err := os.OpenFile(filepath, os.O_APPEND|os.O_WRONLY, 0755)
	if err != nil {
		if os.IsNotExist(err) {
			_, err := os.Create(filepath)
			if err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
	}
	return f
}

func joinSlices(new []string, existing []string) []string {
	for _, i := range new {
		if !sliceContains(existing, i) {
			existing = append(existing, i)
		}
	}
	return existing
}

func sliceContains(slice []string, value string) bool {

	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

func dumpStringsSliceToFile(repos []string, filepath string) {
	content := strings.Join(repos, "\n")
	ioutil.WriteFile(filepath, []byte(content), 0755)
}

func scanGitFolders(folders []string, folder string) []string {
	folder = strings.TrimSuffix(folder, "/")

	f, err := os.Open(folder)
	if err != nil {
		log.Fatal(err)
	}
	files, err := f.ReadDir(-1)
	f.Close()
	if err != nil {
		log.Fatal(err)
	}

	var path string

	for _, file := range files {
		if file.IsDir() {
			path = folder + "/" + file.Name()
			if file.Name() == ".git" {
				path = strings.TrimSuffix(path, "/.git")
				fmt.Println(path)
				folders = append(folders, path)
				continue
			}

			if file.Name() == "vendor" || file.Name() == "node_modules" {
				continue
			}
			folders = scanGitFolders(folders, path)
		}
	}
	return folders
}
