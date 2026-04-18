package main

import "fmt"

type option struct {
  key      byte
  title    string
  disabled bool
}

type dialogue struct {
  klgr *keyLogger
}

func (d *dialogue) ask(options []option) (selectedKey byte) {
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
    selectedKey = <-d.klgr.C
    if _, ok := keys[selectedKey]; ok {
      break
    }
  }
  return selectedKey
}

func (d *dialogue) cleanAsk(options []option) (selectedKey byte) {
  fmt.Print("\033[2J") // Clear screen
  fmt.Print("\033[H") // Go to (0, 0)
  return d.askAndClear(options)
}

func (d *dialogue) askAndClear(options []option) (selectedKey byte) {
  selectedKey = d.ask(options)
  fmt.Print("\033[2J") // Clear screen
  fmt.Print("\033[H") // Go to (0, 0)
  return selectedKey
}

func (d *dialogue) askTransparent(options []option) (selectedKey byte) {
  fmt.Print("\033[H") // Go to (0, 0)
  return d.ask(options)
}

func (d *dialogue) askAnyKey() {
  fmt.Println(" ┌──────────────────────────────────┐")
  fmt.Println(" │ Press any key to continue.       │")
  fmt.Println(" └──────────────────────────────────┘")
  <-d.klgr.C
}

