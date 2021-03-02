package util

import (
	"fmt"
	"os"
	"path"
	"time"
)

func GetDownloadPath() string {
	path, _ := os.UserHomeDir()
	return path
}

func GetDownloadFilename(name string) string {
	if name == "" {
		name = time.Now().Format("2006-01-02")
	}
	filename := fmt.Sprintf("es-%s.log", name)
	return path.Join(GetDownloadPath(), filename)
}