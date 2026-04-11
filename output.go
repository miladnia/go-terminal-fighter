package main

import (
  "fmt"
)

type writer interface {
  write(s string)
}

type terminalWriter struct {}

func (w *terminalWriter) write(s string) {
  // Go to (0, 0)
  fmt.Print("\033[H")
  fmt.Println(s)
}

