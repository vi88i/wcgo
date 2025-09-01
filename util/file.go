package util

import (
	"fmt"
	"os"
	"path/filepath"
	"wcgo/constants"
)

func GetFiles(directory string) []string {
  files := []string{}

	filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
      if (err != nil) {
        return err
      }

      if !info.IsDir() {
        files = append(files, path)
      }

      return nil
    },
  )

  return files
}

func IsValidDirectory(dir string) bool {
  if dir == constants.NoDirectory {
    fmt.Println("No directory provided")
    return false
  } else if _, err := os.Stat(dir); err != nil {
    fmt.Printf("Invalid path: %v\n", dir)
    return false
  }
  return true
}
