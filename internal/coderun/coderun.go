package coderun

import (
	"BuntServer/internal/models"
	"bytes"
	"fmt"
	"math/rand/v2"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

var languages = map[string]func(*models.CodeRunRequest) *models.CodeRunOutput{
	"javascript": runJavaScript,
	"python":     runPython,
	"java":       runJava,
	"go":         runGo,
}

func runGo(crRequest *models.CodeRunRequest) *models.CodeRunOutput {
	outBuf := new(bytes.Buffer)
	errBuf := new(bytes.Buffer)

	cmd := exec.Command("go", "run", crRequest.MainFilePath)
	cmd.Stdout = outBuf
	cmd.Stderr = errBuf
	cmd.Run()

	return &models.CodeRunOutput{
		Output: outBuf.String(),
		Err:    errBuf.String(),
	}
}

func runJava(crRequest *models.CodeRunRequest) *models.CodeRunOutput {
	outBuf := new(bytes.Buffer)
	errBuf := new(bytes.Buffer)

	cmd := exec.Command("javac", "-d", ".", crRequest.MainFilePath)
	cmd.Stdout = outBuf
	cmd.Stderr = errBuf
	cmd.Run()

	runCmd := exec.Command("java", "Main")
	runCmd.Stdout = outBuf
	runCmd.Stderr = errBuf
	runCmd.Run()

	return &models.CodeRunOutput{
		Output: outBuf.String(),
		Err:    errBuf.String(),
	}
}

func runPython(crRequest *models.CodeRunRequest) *models.CodeRunOutput {
	outBuf := new(bytes.Buffer)
	errBuf := new(bytes.Buffer)

	cmd := exec.Command("py", crRequest.MainFilePath)
	cmd.Stdout = outBuf
	cmd.Stderr = errBuf
	cmd.Run()

	return &models.CodeRunOutput{
		Output: outBuf.String(),
		Err:    errBuf.String(),
	}
}

func runJavaScript(crRequest *models.CodeRunRequest) *models.CodeRunOutput {
	outBuf := new(bytes.Buffer)
	errBuf := new(bytes.Buffer)

	cmd := exec.Command("node", crRequest.MainFilePath)
	cmd.Stdout = outBuf
	cmd.Stderr = errBuf
	cmd.Run()

	return &models.CodeRunOutput{
		Output: outBuf.String(),
		Err:    errBuf.String(),
	}
}

func splitPath(path string) (string, string) {
	i := strings.LastIndex(path, "/")

	if i < 0 {
		return "", path
	}

	dirPath := path[:i]
	filePath := path[i+1:]

	return dirPath, filePath
}

func processCodeRunFiles(parentDir string, crFiles []models.CodeRunFile) {
	for _, file := range crFiles {
		dirPath, _ := splitPath(file.Path)

		os.MkdirAll(parentDir+dirPath, 0700)
		os.WriteFile(parentDir+file.Path, []byte(file.Content), 0644)
	}
}

func ProcessCodeRunRequest(crRequest *models.CodeRunRequest) (*models.CodeRunOutput, error) {
	var output *models.CodeRunOutput
	runner, ok := languages[crRequest.Language]

	if !ok {
		return nil, fmt.Errorf("unsupported language type: %v", crRequest.Language)
	}

	randDir := "tmp"

	startTime := time.Now()
	for i := 0; i < 10; i++ {
		randDir += strconv.Itoa(rand.IntN(10))
	}

	parentDir := "CODERUNNER/" + randDir + "/"
	crRequest.MainFilePath = parentDir + crRequest.MainFilePath
	processCodeRunFiles(parentDir, crRequest.Files)

	// RETURN ERROR IF NO FILES ARE PRESENT
	output = runner(crRequest)

	// Clean up files processed
	os.RemoveAll(parentDir)

	duration := time.Since(startTime)
	output.RunTime = duration

	return output, nil
}
