package main

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
)

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	tree, err := getTreeItems(path, printFiles, "")
	if err != nil {
		return err
	}

	_, err = fmt.Fprintln(out, strings.Join(tree, "\n"))
	if err != nil {
		return err
	}

	return nil
}

func getTreeItems(path string, printFiles bool, prefix string) ([]string, error) {
	var treeItems []string

	f, err := os.Open(path)
	if err != nil {
		return treeItems, err
	}

	fileInfo, err := f.Readdir(-1)

	err = f.Close()
	if err != nil {
		return treeItems, err
	}

	if !printFiles {
		fileInfo = excludeFiles(fileInfo)
	}

	sort.Slice(fileInfo, func(i, j int) bool {
		return fileInfo[i].Name() < fileInfo[j].Name()
	})

	lastIndex := len(fileInfo) - 1

	for i, file := range fileInfo {
		isLast := i == lastIndex
		treeItems = append(treeItems, getItemString(file, prefix, isLast))

		if file.IsDir() {
			items, err := getTreeItems(path+"/"+file.Name(), printFiles, getNextLevelPrefix(prefix, isLast))

			if err != nil {
				return treeItems, err
			}

			treeItems = append(treeItems, items...)
		}
	}

	return treeItems, nil
}

func getNextLevelPrefix(prefix string, isLastItem bool) string {
	if !isLastItem {
		return prefix + "│\t"
	}
	return prefix + "\t"
}

func getItemString(item os.FileInfo, prefix string, isLast bool) (string) {
	if isLast {
		prefix += "└───"
	} else {
		prefix += "├───"
	}

	prefix += item.Name()

	if !item.IsDir() {
		var fileSize string
		if item.Size() == 0 {
			fileSize = "empty"
		} else {
			fileSize = strconv.FormatInt(item.Size(), 10) + "b"
		}
		prefix += " (" + fileSize + ")"
	}

	return prefix
}

func excludeFiles(fileInfo []os.FileInfo) []os.FileInfo {
	var dirs []os.FileInfo

	for _, file := range fileInfo {
		if file.IsDir() {
			dirs = append(dirs, file)
		}
	}

	return dirs
}
