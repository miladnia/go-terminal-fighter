package main

import (
  "fmt"
  "os"
  "syscall"
  "unsafe"
)

func printSpace(s string, count int) string {
  for i := 0; i < count; i++ {
    s = fmt.Sprint(s, " ")
  }
  return s
}

func printBlankLine(count int) {
  for i := 0; i < count; i++ {
    fmt.Println()
  }
}

func clearScreen() {
  // Clear screen
  fmt.Print("\033[2J")
  fmt.Print("\033[H")
}

func captureInput(callback func(key byte) (done bool)) {
  var oldState syscall.Termios
	fd := int(os.Stdin.Fd())

	// Get current terminal state
	if _, _, err := syscall.Syscall6(syscall.SYS_IOCTL, uintptr(fd),
		uintptr(syscall.TCGETS), uintptr(unsafe.Pointer(&oldState)), 0, 0, 0); err != 0 {
		panic(err)
	}

	newState := oldState
	newState.Lflag &^= syscall.ICANON | syscall.ECHO

	// Apply new state (raw-ish mode)
	if _, _, err := syscall.Syscall6(syscall.SYS_IOCTL, uintptr(fd),
		uintptr(syscall.TCSETS), uintptr(unsafe.Pointer(&newState)), 0, 0, 0); err != 0 {
		panic(err)
	}

	defer func() {
		syscall.Syscall6(syscall.SYS_IOCTL, uintptr(fd),
			uintptr(syscall.TCSETS), uintptr(unsafe.Pointer(&oldState)), 0, 0, 0)
	}()

	buf := make([]byte, 1)
  for {
    os.Stdin.Read(buf)
    done := callback(buf[0])
    if done {
      break
    }
  }
}

