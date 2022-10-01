package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	t.Run("usual working case", func(t *testing.T) {
		dir := "./testdata/env"

		env, _ := ReadDir(dir)

		require.Equal(t, env["BAR"].Value, "bar")
		require.Equal(t, env["BAR"].NeedRemove, false)

		require.Equal(t, env["FOO"].Value, "   foo\nwith new line")
		require.Equal(t, env["FOO"].NeedRemove, false)

		require.Equal(t, env["UNSET"].Value, "")
		require.Equal(t, env["UNSET"].NeedRemove, true)

		require.Equal(t, env["HELLO"].Value, "\"hello\"")
		require.Equal(t, env["ADDED"].NeedRemove, false)

		require.Equal(t, env["EMPTY"].Value, "")
		require.Equal(t, env["EMPTY"].NeedRemove, false)
	})

	t.Run("invalid dir case", func(t *testing.T) {
		dir := "./testdata/not_existing_dir"

		env, err := ReadDir(dir)

		require.Nil(t, env)
		require.Error(t, err)
	})
}
