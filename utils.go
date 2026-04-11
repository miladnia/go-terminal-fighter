package main

import (
  "encoding/json"
  "fmt"
  "os"
  "time"
)

type menuOption struct {
  key      byte
  title    string
  disabled bool
}

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

func ask(klgr *keyLogger, options []menuOption) (selectedKey byte) {
  keys := map[byte]struct{}{}
  fmt.Println(" ┌──────────────────────┐")
  fmt.Println(" │ Please Choose:       │")
  for _, opt := range options {
    if opt.disabled {
      continue
    }
    fmt.Printf(" │ [%c] %-16s │\n", opt.key, opt.title)
    keys[opt.key] = struct{}{}
  }
  fmt.Println(" └──────────────────────┘")
  for {
    selectedKey = <-klgr.C
    if _, ok := keys[selectedKey]; ok {
      break
    }
  }
  return selectedKey
}

func decodeJsonFile(filename string, v any) error {
  file, err := os.Open(filename)
  if err != nil {
    return fmt.Errorf("could not open the file '%v'! %s", filename, err)
  }
  defer file.Close()

  decoder := json.NewDecoder(file)
  if err := decoder.Decode(v); err != nil {
    return fmt.Errorf("could not decode the file '%v'! %s", filename, err)
  }

  return nil
}

