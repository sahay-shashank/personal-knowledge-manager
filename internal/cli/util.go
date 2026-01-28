package cli

import (
	"os"
	"os/exec"
)

func tempEditor(content *string) (string, error) {
	tempFile, err := os.CreateTemp("", "pkm-*.data")
	if err != nil {
		return "", err
	}
	defer os.Remove(tempFile.Name())

	if content != nil {
		tempFile.WriteString(*content)
	}

	tempFile.Close()

	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vi"
	}
	cmd := exec.Command(editor, tempFile.Name())
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		return "", err
	}

	newContent, err := os.ReadFile(tempFile.Name())
	if err != nil {
		return "", err
	}
	return string(newContent), nil
}
