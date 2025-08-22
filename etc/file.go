package etc

import (
	"fmt"
	"os"
	"path/filepath"
)

func ReadFile(filePath string) (string, error) {
	extension := filepath.Ext(filePath)
	if extension != ".cm" {
		fmt.Println("not a Cutter(.cm) file: " + filePath)
		os.Exit(0)
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func WriteFile(filePath string, content string) error {
	err := os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		return err
	}
	return nil
}
