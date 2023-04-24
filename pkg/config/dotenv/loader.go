package dotenv

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"regexp"
)

func Load(file string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	reg := regexp.MustCompile(`^\s*([\w.]+)\s*=\s*(.*)\s*$`)

	lineNo := 0
	for scanner.Scan() {
		lineNo++
		line := []byte(scanner.Text())
		line = bytes.TrimSpace(line)

		// 空行
		if len(line) == 0 {
			continue
		}

		// 注释行
		if bytes.TrimSpace(line)[0] == '#' {
			continue
		}

		res := reg.FindSubmatch(line)
		if len(res) != 3 {
			return fmt.Errorf("config format error @line %d", lineNo)
		}
		os.Setenv(string(res[1]), string(res[2]))
	}

	return nil
}
