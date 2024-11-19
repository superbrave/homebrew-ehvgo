package util

import (
	"bufio"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

func GetRepositoryName() string {
  gitDir, err := traverseParentDirectories()
  if err != nil {
    return ""
  }

  cFile, err := os.Open(filepath.Join(gitDir, string(os.PathSeparator), ".git", string(os.PathSeparator), "config"))
  if err != nil {
    return ""
  }
  defer cFile.Close()

  scanner := bufio.NewScanner(cFile)
  for scanner.Scan() {
    line := strings.TrimSpace(scanner.Text())

    if strings.HasPrefix(line, "url =") {
      url := strings.TrimSpace(strings.TrimPrefix(line, "url ="))
      parts := strings.Split(strings.TrimSuffix(url, ".git"), "/")
      if len(parts) > 0 {
        return parts[len(parts)-1]
      }
    }
  }

  return ""
}

func traverseParentDirectories() (string, error) {
  dir,_ := os.Getwd()

  for {
    gitDir := filepath.Join(dir, ".git")

    if stat, err := os.Stat(gitDir); err == nil && stat.IsDir() {
      return dir, nil
    }

    parent := filepath.Dir(dir)
    if parent == dir {
      break
    }

    dir = parent
  }

  return "", errors.New("No .git directory found in any of the parent directories")
}