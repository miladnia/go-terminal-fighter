//go:build linux || darwin || freebsd || netbsd || openbsd

package main

import (
  "os"
  "golang.org/x/sys/unix"
)

func readKey() byte {
	fd := int(os.Stdin.Fd())

	// Get current terminal state
	oldState, err := unix.IoctlGetTermios(fd, unix.TCGETS)
	if err != nil {
		panic(err)
	}

	// Create raw state (cbreak mode - minimal changes)
	newState := *oldState
	// Disable canonical mode (line buffering) and echo
	newState.Lflag &^= unix.ICANON | unix.ECHO
	// Set VMIN to 1 (wait for 1 character) and VTIME to 0 (no timeout)
	newState.Cc[unix.VMIN] = 1
	newState.Cc[unix.VTIME] = 0

	// Apply new state
	if err := unix.IoctlSetTermios(fd, unix.TCSETS, &newState); err != nil {
		panic(err)
	}

	// Read single byte
	buf := make([]byte, 1)
	os.Stdin.Read(buf)
	key := buf[0]

	// Restore original state
	if err := unix.IoctlSetTermios(fd, unix.TCSETS, oldState); err != nil {
		panic(err)
	}

	// Convert uppercase to lowercase
	if 'A' <= key && key <= 'Z' {
		key += 'a' - 'A'
	}

	return key
}

