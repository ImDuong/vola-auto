package utils

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"unicode"
)

func GetPythonRunner() (pythonRunner string, pythonVersion string, err error) {
	pythonRunner = "python"
	cmd := exec.Command(pythonRunner, "-V")
	output, err := cmd.CombinedOutput()
	if err == nil {
		pythonVersion := string(output)
		return pythonRunner, pythonVersion, nil
	}

	pythonRunner = "python3"
	cmd = exec.Command(pythonRunner, "-V")
	output, err = cmd.CombinedOutput()
	if err == nil {
		pythonVersion := string(output)
		return pythonRunner, pythonVersion, nil
	}

	pythonRunner = "python2"
	cmd = exec.Command(pythonRunner, "-V")
	output, err = cmd.CombinedOutput()
	if err == nil {
		pythonVersion := string(output)
		return pythonRunner, pythonVersion, nil
	}
	return pythonRunner, pythonVersion, fmt.Errorf("cannot find python")
}

// GetPathInCamelCase converts a Windows path to CamelCase, preserving the case of file names
func GetPathInCamelCase(convertingPath string) string {
	if convertingPath == "" {
		return ""
	}

	// Split the path into its components except the last component
	components := strings.Split(convertingPath, "\\")
	for i := 0; i < len(components)-1; i++ {
		components[i] = toCamelCase(components[i])
	}

	// Preserve the case of the last component if it's a file
	if isDirectory(filepath.Join(components...)) {
		components[len(components)-1] = toCamelCase(components[len(components)-1])
	}

	return strings.Join(components, "\\")
}

// toCamelCase converts a string to CamelCase
func toCamelCase(s string) string {
	var result strings.Builder
	isNewWord := true

	for _, char := range s {
		if unicode.IsLetter(char) {
			if isNewWord {
				result.WriteRune(unicode.ToUpper(char))
				isNewWord = false
			} else {
				result.WriteRune(unicode.ToLower(char))
			}
		} else {
			result.WriteRune(char)
			isNewWord = (char == ' ')
		}
	}

	return result.String()
}

// isDirectory checks if the given path is a directory or not
func isDirectory(path string) bool {
	return !strings.Contains(filepath.Ext(path), ".")
}
