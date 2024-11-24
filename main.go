// Author: Erdet Nasufi <erdet.nasufi@gmail.com> //

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
)

const (
	kTermColorReset   = "\033[0m"
	kTermColorRed     = "\033[31m"
	kTermColorGreen   = "\033[32m"
	kTermColorYellow  = "\033[33m"
	kTermColorBlue    = "\033[34m"
	kTermColorMagenta = "\033[35m"
	kTermColorCyan    = "\033[36m"
	kTermColorGray    = "\033[37m"
	kTermColorWhite   = "\033[97m"
)

const (
	MX_FILE_SIZE = 1048576
)

const (
	MX_PROJECT_FILE = ".mxproject"
	MK_FILE         = "Debug/Core/Src/subdir.mk"
)

func Ok(err error, fn string, ifFailMsg string) {
	if err != nil {
		log.Fatalln(kTermColorRed, fn, ":", ifFailMsg, "\n", err.Error(), kTermColorReset)
	}
}

func findMxProjectFile(root string) *os.File {
	mxProjectFullPath := fmt.Sprintf("%s/%s", root, MX_PROJECT_FILE)
	mxProjectFile, err := os.Open(mxProjectFullPath)

	Ok(err, reflect.Func.String(), "Failed to load .mxproject file.")

	return mxProjectFile
}

func findHeaderPaths(root string, mxFile *os.File) []string {
	buff := make([]byte, MX_FILE_SIZE)
	_, err := mxFile.Read(buff)
	Ok(err, reflect.Func.String(), "Failed to find the HeaderPath.")

	content := string(buff)

	n := strings.Index(content, "HeaderPath=")
	if n < 0 {
		Ok(fmt.Errorf("Index <0"), reflect.Func.String(), "Failed to find the HeaderPath.")
	}
	content = content[n:]

	n = strings.IndexAny(content, "\n")
	if n < 0 {
		Ok(fmt.Errorf("Index <0"), reflect.Func.String(), "Failed to parse the HeaderPath.")
	}
	content = content[11:n]
	content = strings.TrimRight(content, "\n")

	paths := strings.Split(content, ";")
	nHeaderPaths := len(paths)
	if nHeaderPaths == 0 {
		Ok(fmt.Errorf("No path(s) found."), reflect.Func.String(), "Failed to find header paths.")
	}
	return paths
}

func createClangdFile(root string, wall bool) {
	buff, err := os.ReadFile(fmt.Sprintf("%s/%s", root, MK_FILE))
	Ok(err, reflect.Func.String(), "Failed to read subdir.mk")

	content := string(buff)
	n1 := strings.Index(content, "arm-none-eabi-gcc")
	n2 := strings.Index(content[n1:], "\n")
	content = content[n1:(n1 + n2)]

	flags := make([]string, 0)
	var flag string
	for {
		n1 = strings.Index(content, "-D")
		if n1 == -1 {
			break
		}
		content = content[n1+2:]
		n2 = strings.Index(content, " ")
		flag = content[:n2]

		flags = append(flags, flag)
	}

	content = "#.clangd\n"
	content += "CompileFlags:\n"
	content += "\tAdd: [ "
	for _, flag := range flags {
		content += fmt.Sprintf("-D%s, ", flag)
	}

	if wall {
		content += fmt.Sprintf("-Wall, ")
	}

	mxFile := findMxProjectFile(root)
	includePaths := findHeaderPaths(root, mxFile)

	countIncludePaths := len(includePaths)
	for i, includePath := range includePaths {
		content += fmt.Sprintf("-I%s/%s", root, includePath)

		if i < (countIncludePaths - 1) {
			content += ", "
		}
	}
	content += " ]\n\n"

	clangdFile, err := os.Create(fmt.Sprintf("%s/.clangd", root))
	Ok(err, reflect.Func.String(), "Failed to create .clangd file")
	defer clangdFile.Close()

	_, err = clangdFile.Write([]byte(content))
	Ok(err, reflect.Func.String(), "Failed to write .clangd file")

	println(kTermColorGreen, ".clangd file is added successfully.", kTermColorReset)
}

func main() {
	wall := flag.Bool("Wall", false, "Add -Wall flag. By default, not added.")
	flag.Parse()

	pwd, err := os.Getwd()
	Ok(err, reflect.Func.String(), "Failed to find the working directory.")

	createClangdFile(pwd, *wall)
}
