package main

import (
	"io"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	stdoutStubName := "stdout_stub.txt"

	t.Run("check no argument command and result code", func(t *testing.T) {
		stdoutStub, _ := os.Create(stdoutStubName)
		defer func() {
			stdoutStub.Close()
			os.Remove(stdoutStubName)
		}()

		// подменим stdout на файл
		stdOut := os.Stdout
		os.Stdout = stdoutStub
		returnCode := RunCmd([]string{"pwd"}, Environment{})
		os.Stdout = stdOut

		path, _ := os.Getwd()
		stdoutStub.Seek(0, 0)
		bytes, _ := io.ReadAll(stdoutStub)

		// Добавим к ожидаемой строке "\n", т.к. в stdout выводится с переносом строки на конце
		require.Equal(t, path+"\n", string(bytes))
		require.Equal(t, 0, returnCode)
	})

	t.Run("check command with argument and env changes", func(t *testing.T) {
		fileFrom, _ := os.Create(stdoutStubName)
		defer func() {
			fileFrom.Close()
			os.Remove(stdoutStubName)
		}()

		os.Setenv("VAR_TO_UPDATE", "321")
		os.Setenv("VAR_TO_REMOVE", "123")
		// проверим что переменные успешно установлены
		require.Equal(t, "321", os.Getenv("VAR_TO_UPDATE"))
		require.Equal(t, "123", os.Getenv("VAR_TO_REMOVE"))

		env := Environment{
			"VAR_TO_CREATE": EnvValue{Value: "NEW_VALUE", NeedRemove: false},
			"VAR_TO_UPDATE": EnvValue{Value: "000", NeedRemove: false},
			"VAR_TO_REMOVE": EnvValue{Value: "", NeedRemove: true},
		}

		// подменим stdout на файл
		stdOut := os.Stdout
		os.Stdout = fileFrom
		RunCmd([]string{"printenv", "VAR_TO_CREATE"}, env)
		RunCmd([]string{"printenv", "VAR_TO_UPDATE"}, env)
		RunCmd([]string{"printenv", "VAR_TO_REMOVE"}, env)
		os.Stdout = stdOut

		fileFrom.Seek(0, 0)
		bytes, _ := io.ReadAll(fileFrom)

		require.Equal(t, "NEW_VALUE\n000\n", string(bytes))
	})

	t.Run("test sleep", func(t *testing.T) {
		start := time.Now()
		returnCode := RunCmd([]string{"sleep", "1"}, Environment{})
		elapsedTime := time.Since(start)

		require.Equal(t, 0, returnCode)
		require.GreaterOrEqual(t, int64(elapsedTime), int64(time.Second))
	})

	t.Run("invalid command argument", func(t *testing.T) {
		returnCode := RunCmd([]string{"sleep", "dfg"}, Environment{})

		require.Equal(t, 1, returnCode)
	})
}
