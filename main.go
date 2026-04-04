package main

import (
  "fmt"
  "strings"
  "time"
)

type context struct {
  selectedMap [][]int
  mapIndex int
  hRatio int
  vRatio int
  blankLineCount int
  visibleMap [][]int
  height int
  speed int
}

func (ctx *context) appendToVisibleMap(mapLine []int) {
  if len(ctx.visibleMap) >= ctx.height {
    // Remove the oldest map line
    // to free space for the ongoing
    // new line.
    ctx.visibleMap = ctx.visibleMap[1:]
  }
  ctx.visibleMap = append(ctx.visibleMap, mapLine)
}

func (ctx *context) readMapChunk() (chunk [][]int, end bool) {
  // The whole map have been read.
  if ctx.mapIndex <= 0 {
    if ctx.blankLineCount < ctx.height {
      ctx.blankLineCount++
      ctx.appendToVisibleMap([]int{})
      return ctx.visibleMap, false
    }
    return nil, true
  }
  // Add blank lines between map lines.
  if ctx.blankLineCount < ctx.vRatio {
    ctx.blankLineCount++
    ctx.appendToVisibleMap([]int{})
    return ctx.visibleMap, false
  }
  ctx.blankLineCount = 0
  i := ctx.mapIndex
  ctx.mapIndex--
  mapLine := ctx.selectedMap[i]
  ctx.appendToVisibleMap(mapLine)
  return ctx.visibleMap, false
}

func main() {
  startGame()
}

func startGame() {
  ctx := newGame()
  delayMS := time.Duration(500 / ctx.speed)
  ticker := time.NewTicker(delayMS * time.Millisecond)
  for ; ; <-ticker.C {
    complete := render(&ctx)
    if complete {
      break
    }
  }
  printBlankLine(5)
  fmt.Println("GAME OVER!")
  fmt.Print("Do you want to continue? (y/n) ")
  var answer string
  fmt.Scanln(&answer)
  answer = strings.ToLower(answer)
  if answer == "n" || answer == "no" {
    fmt.Println("GOOD GAME!")
    time.Sleep(1000 * time.Millisecond)
    clearScreen()
    return
  }
  startGame()
}

func render(ctx *context) (complete bool) {
  mapChunk, end := ctx.readMapChunk()
  if end {
    return true
  }
  clearScreen()
  for i := len(mapChunk) - 1; i >= 0; i-- {
    mapLine := mapChunk[i]
    for _, item := range mapLine {
      switch item {
        case 0:
          printSpace(ctx.hRatio)
        case 1:
          fmt.Print("#")
      }
    }
    fmt.Println()
  }
  return false
}

func printSpace(count int) {
  for i := 0; i < count; i++ {
    fmt.Print(" ")
  }
}

func printBlankLine(count int) {
  for i := 0; i < count; i++ {
    fmt.Println()
  }
}

func newGame() context {
  selectedMap := getMap()
  return context{
    selectedMap: selectedMap,
    mapIndex: len(selectedMap) - 1,
    hRatio: 5,
    vRatio: 3,
    visibleMap: [][]int{},
    height: 20,
    speed: 10,
  }
}

func getMap() [][]int {
  return [][]int{
    []int{0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0},
    []int{0, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 1, 0},
    []int{0, 0, 0, 0, 0, 1, 1, 0, 1, 1, 0, 0, 0, 0, 0},
    []int{0, 1, 1, 1, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 0},
    []int{0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0},
    []int{0, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 1, 0},
    []int{0, 0, 0, 0, 0, 1, 1, 0, 1, 1, 0, 0, 0, 0, 0},
    []int{0, 1, 1, 1, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 0},
    []int{0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0},
    []int{0, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 1, 0},
    []int{0, 0, 0, 0, 0, 1, 1, 0, 1, 1, 0, 0, 0, 0, 0},
    []int{0, 1, 1, 1, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 0},
    []int{0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0},
    []int{0, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 1, 0},
    []int{0, 0, 0, 0, 0, 1, 1, 0, 1, 1, 0, 0, 0, 0, 0},
    []int{0, 1, 1, 1, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 0},
  }
}

func clearScreen() {
  // Clear screen
  fmt.Print("\033[2J")
  fmt.Print("\033[H")
}
