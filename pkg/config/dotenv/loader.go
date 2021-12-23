package dotenv

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type nodeName = string
type nodeFullName = string

type node struct {
	Name     nodeName
	Value    []byte
	Children map[nodeName]*node
}

func (n *node) addOrGetChildNode(name nodeName, value []byte) *node {
	if n.Children == nil {
		n.Children = make(map[nodeName]*node)
	}

	if _, found := n.Children[name]; !found {
		n.Children[name] = &node{
			Name:  name,
			Value: value,
		}
	}

	return n.Children[name]
}

type tree struct {
	Root *node
}

func (t *tree) build(fn nodeFullName, value []byte) {
	curNode := t.Root
	keySegs := strings.Split(string(fn), ".")
	keyCnt := len(keySegs)
	for i, key := range keySegs {
		if i == keyCnt-1 {
			curNode = curNode.addOrGetChildNode(key, value)
		} else {
			curNode = curNode.addOrGetChildNode(key, nil)
		}
	}
}

func (n *node) resolve() map[string]interface{} {
	ret := make(map[string]interface{})
	if n.Children == nil || len(n.Children) == 0 {
		ret[n.Name] = (json.RawMessage)(n.Value)
		return ret
	}
	for _, cn := range n.Children {
		if cn.Children == nil || len(cn.Children) == 0 {
			ret[cn.Name] = (json.RawMessage)(cn.Value)
		} else {
			ret[cn.Name] = cn.resolve()
		}
	}
	return ret
}

func Load(file string, input interface{}) error {
	fh, err := os.Open(file)
	if err != nil {
		return err
	}
	defer fh.Close()

	tree := &tree{Root: &node{}}

	scanner := bufio.NewScanner(fh)
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

		fullNodeName, value := nodeFullName(res[1]), res[2]
		tree.build(fullNodeName, value)
	}

	data := tree.Root.resolve()
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	err = json.Unmarshal(jsonData, &input)
	if err != nil {
		return err
	}

	return nil
}
