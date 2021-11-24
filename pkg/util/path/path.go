package path

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func FindPath(path string, depth int) (string, error) {
	if strings.HasPrefix(path, "/") {
		_, err := os.Stat(path)
		return path, err
	}

	// 从当前目录开始寻找
	orginPath := path
	for i := 0; i < depth; i++ {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
		path = fmt.Sprintf("../%s", path)
	}

	// 从执行目录开始寻找
	execPath, err := ExecutablePath()
	if err != nil {
		return "", err
	}

	path = orginPath
	for i := 0; i < depth; i++ {
		p := fmt.Sprintf("%s%c%s", execPath, filepath.Separator, path)
		if _, err := os.Stat(p); err == nil {
			return p, nil
		}
		path = fmt.Sprintf("../%s", path)
	}

	return "", fmt.Errorf("%s not found", orginPath)
}

func ExecutablePath() (string, error) {
	dir, err := os.Executable()
	if err != nil {
		return "", err
	}

	return filepath.Dir(dir), nil
}
