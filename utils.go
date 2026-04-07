package main

import (
  "fmt"
  "time"
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
  // Go to (0, 0)
  fmt.Print("\033[H")
}

func flashMessage(msg string) {
  clearScreen()
  time.Sleep(500 * time.Millisecond)
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

func ask(klgr *keyLogger, options map[byte]string) (selectedKey byte) {
  fmt.Println("==============")
  fmt.Println("Please Choose:")
  for key, title := range options {
    fmt.Printf("[%c] %s\n", key, title)
  }
  fmt.Println("==============")
  for {
    selectedKey = <-klgr.C
    if _, ok := options[selectedKey]; ok {
      break
    }
  }
  clearScreen()
  return selectedKey
}

