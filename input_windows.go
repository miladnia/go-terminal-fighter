//go:build windows

package main

import (
	"os"
	"golang.org/x/sys/windows"
)

func readKey() byte {
	fd := windows.Handle(os.Stdin.Fd())

	// Get current console mode
	var oldMode uint32
	if err := windows.GetConsoleMode(fd, &oldMode); err != nil {
		panic(err)
	}

	// Disable echo and line input (raw mode equivalents)
	// Clear: ENABLE_ECHO_INPUT (don't show typed chars)
	// Clear: ENABLE_LINE_INPUT (read keys immediately, don't wait for Enter)
	// Clear: ENABLE_PROCESSED_INPUT (don't process Ctrl+C, etc.)
	newMode := oldMode &^ (windows.ENABLE_ECHO_INPUT | windows.ENABLE_LINE_INPUT | windows.ENABLE_PROCESSED_INPUT)

	// Apply new mode
	if err := windows.SetConsoleMode(fd, newMode); err != nil {
		panic(err)
	}

	// Read single byte
	buf := make([]byte, 1)
	os.Stdin.Read(buf)
	key := buf[0]

	// Restore original mode
	if err := windows.SetConsoleMode(fd, oldMode); err != nil {
		panic(err)
	}

	// Convert uppercase to lowercase
	if 'A' <= key && key <= 'Z' {
		key += 'a' - 'A'
	}

	return key
}

