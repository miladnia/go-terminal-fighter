package main

import (
  "fmt"
)

type writer interface {
  write(s string)
}

type terminalWriter struct {}

func (w *terminalWriter) write(s string) {
  // Clear screen
  fmt.Print("\033[2J")
  // Go to (0, 0)
  fmt.Print("\033[H")
  fmt.Println(s)
}

