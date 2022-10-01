package main

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	environment := make(Environment, len(files))

	for _, file := range files {
		fileName := file.Name()
		filePath := filepath.Join(dir, fileName)

		envValue, err := readEnvValueFromFile(filePath)
		if err != nil {
			return nil, err
		}
		environment[fileName] = *envValue
	}

	return environment, nil
}

func readEnvValueFromFile(filePath string) (*EnvValue, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bufReader := bufio.NewReader(file)
	lineBytes, _, err := bufReader.ReadLine()

	envValue := ""
	needRemove := true
	switch {
	case err == nil:
		// заменим терминальные нули 0x00 на перенос строки \n
		lineBytes = bytes.ReplaceAll(lineBytes, []byte("\x00"), []byte("\n"))
		envValue = strings.TrimRight(string(lineBytes), " \t\n")
		needRemove = false
	case !errors.Is(err, io.EOF):
		return nil, err
	}

	return &EnvValue{
		Value:      envValue,
		NeedRemove: needRemove,
	}, nil
}
