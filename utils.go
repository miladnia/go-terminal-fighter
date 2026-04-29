package main

import (
  "fmt"
  "time"
)

func clearScreen() {
  fmt.Print("\033[2J")
  goToScreenTop()
}

// Go to (0, 0)
func goToScreenTop() {
  fmt.Print("\033[H")
}

func hideCursor() {
  fmt.Print("\033[?25l")
}

func showCursor() {
  fmt.Print("\033[?25h")
}

func showInfo(lines []string) {
  clearScreen()
  for _, ln := range lines {
    fmt.Println(ln)
    time.Sleep(250 * time.Millisecond)
  }
}

