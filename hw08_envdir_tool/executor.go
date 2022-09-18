package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	commandName := strings.Trim(cmd[0], " \t\n") // линтер ругается если напрямую передаем входящий аргумент в exec.Command()
	commandArgs := cmd[1:]

	command := exec.Command(commandName, commandArgs...)

	currentEnvVariablesMap := getCurrentEnvVariables()
	newEnvVariables := updateEnvVariables(currentEnvVariablesMap, env)
	envSlice := createEnvSlice(newEnvVariables)

	command.Env = envSlice
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	if err := command.Run(); err != nil {
		var exitErr *exec.ExitError
		errors.As(err, &exitErr)

		return exitErr.ExitCode()
	}

	return 0
}

func getCurrentEnvVariables() map[string]string {
	result := make(map[string]string)
	for _, envString := range os.Environ() {
		s := strings.Split(envString, "=")
		result[s[0]] = s[1]
	}

	return result
}

func updateEnvVariables(envVariables map[string]string, environmentToApply Environment) map[string]string {
	newEnvVariables := copyStringMap(envVariables)
	for envName, envValue := range environmentToApply {
		_, ok := newEnvVariables[envName]
		if !envValue.NeedRemove {
			newEnvVariables[envName] = envValue.Value
		} else if ok {
			delete(newEnvVariables, envName)
		}
	}

	return newEnvVariables
}

// Создает слайс с переменными окружения формата os.Environ(), т.е. VAR1=VAL1 VAR2=VAL2 ...
func createEnvSlice(envMap map[string]string) []string {
	envs := make([]string, 0, len(envMap))
	for envName, envValue := range envMap {
		envs = append(envs, fmt.Sprintf("%v=%v", envName, envValue))
	}

	return envs
}

func copyStringMap(originalMap map[string]string) map[string]string {
	targetMap := make(map[string]string)
	for key, value := range originalMap {
		targetMap[key] = value
	}

	return targetMap
}
