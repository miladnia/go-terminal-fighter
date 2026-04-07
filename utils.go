package main

import (
  "fmt"
  "time"
)

type menuOption struct {
  key byte
  title string
}

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

func ask(klgr *keyLogger, options []menuOption) (selectedKey byte) {
  keys := map[byte]struct{}{}
  fmt.Println("==============")
  fmt.Println("Please Choose:")
  for _, opt := range options {
    fmt.Printf("[%c] %s\n", opt.key, opt.title)
    keys[opt.key] = struct{}{}
  }
  fmt.Println("==============")
  for {
    selectedKey = <-klgr.C
    if _, ok := keys[selectedKey]; ok {
      break
    }
  }
  clearScreen()
  return selectedKey
}

