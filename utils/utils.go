package utils

import (
	"fmt"
	"os/exec"
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
