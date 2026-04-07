package main

import (
  "os"
  "syscall"
  "unsafe"
)

type keyLogger struct {
  C <-chan byte
}

func newKeyLogger() *keyLogger {
  c := make(chan byte)
  klgr := &keyLogger{
    C: c,
  }
  go func() {
    for {
      key := readKey()
      c <- key
    }
  }()
  return klgr
}

func readKey() (key byte) {
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
  os.Stdin.Read(buf)

  key = buf[0]
  if 'A' <= key && key <= 'Z' {
    key += 'a' - 'A'
  }

  return key
}

