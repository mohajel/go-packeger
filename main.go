package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run <path/to/your/file.go>")
		return
	}

	filePath := os.Args[1]

	// remove final file name from path
	path := strings.Join(strings.Split(filePath, "/")[:len(strings.Split(filePath, "/"))-1], "/")

	moduleName := getModuleName()
	imports := getImports(filePath)

	fmt.Printf("Module name: %s\n", moduleName)
	fmt.Printf("Imports: %v\n", imports)

	if _, err := os.Stat("./packager-result"); os.IsNotExist(err) {
		err := os.Mkdir("./packager-result", os.ModePerm)
		if err != nil {
			fmt.Println("Error creating directory:", err)
			return
		}
	}

	internalFile, err := os.Create("./packager-result/internal_packages.txt")
	if err != nil {
		fmt.Println("Error creating internal packages file:", err)
		return
	}

	defer internalFile.Close()

	externalFile, err := os.Create("./packager-result/external_packages.txt")
	if err != nil {
		fmt.Println("Error creating external packages file:", err)
		return
	}
	defer externalFile.Close()

	for _, pkg := range imports {
		if isInternal(pkg, moduleName) {
			writeDoc(internalFile, pkg, true)
		} else {
			writeDoc(externalFile, pkg, false)
		}
	}

	// execute go docs -all -src <paht>
	cmd := exec.Command("go", "doc", "-all", "-src", path)
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("Error getting doc for package %s: %v\n", path, err)
		return
	}
	// save the output to a ./packager-result/docs.txt
	docFile, err := os.Create("./packager-result/docs.txt")
	if err != nil {
		fmt.Println("Error creating docs file:", err)
		return
	}
	defer docFile.Close()

	writer := bufio.NewWriter(docFile)
	writer.WriteString(string(output))
	defer writer.Flush()
}

func getModuleName() string {
	cmd := exec.Command("go", "list", "-m")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error getting module name:", err)
		os.Exit(1)
	}
	return strings.TrimSpace(string(output))
}

func getImports(filePath string) []string {
	cmd := exec.Command("go", "list", "-f", "{{.Imports}} {{.TestImports}} {{.XTestImports}}", filePath)
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error getting imports:", err)
		os.Exit(1)
	}
	return ExtractPackageNames(string(output))
}

func isInternal(pkg, moduleName string) bool {
	// Check if the package is an internal package by checking if it contains the module name
	return strings.Contains(pkg, moduleName)
}

func writeDoc(file *os.File, pkg string, isInternal bool) {
	var cmd *exec.Cmd
	if isInternal {
		cmd = exec.Command("go", "doc", "--all", pkg)
	} else {
		cmd = exec.Command("go", "doc", pkg)
	}

	output, err := cmd.Output()
	if err != nil {
		fmt.Fprintf(file, "Error getting doc for package %s: %v\n", pkg, err)
		return
	}

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	writer.WriteString(fmt.Sprintf("Documentation for package %s:\n", pkg))
	writer.WriteString(string(output))
	writer.WriteString("\n\n")
}
