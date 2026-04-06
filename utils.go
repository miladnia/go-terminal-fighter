package main

import (
  "fmt"
  "os"
  "syscall"
  "time"
  "unsafe"
)

func printSpace(s string, count int) string {
  for i := 0; i < count; i++ {
    s = fmt.Sprint(s, " ")
  }
  return s
}

func clearScreen() {
  // Clear screen
  fmt.Print("\033[2J")
  fmt.Print("\033[H")
}

func flashMessage(msg string) {
  clearScreen()
  fmt.Println(msg)
  time.Sleep(1000 * time.Millisecond)
  clearScreen()
}

func showInfo(lines []string) {
  clearScreen()
  for _, ln := range lines {
    fmt.Println(ln)
    time.Sleep(500 * time.Millisecond)
  }
}

func askToChoose(options map[byte]string) (selectedKey byte) {
  clearScreen()
  fmt.Println("==============")
  fmt.Println("Please Choose:")
  for key, title := range options {
    fmt.Printf("[%c] %s\n", key, title)
  }
  fmt.Println("==============")
  for {
    selectedKey = captureInput()
    if _, ok := options[selectedKey]; ok {
      break
    }
  }
  return selectedKey
}

func captureInput() (key byte) {
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

