package main

import (
  "fmt"
  "time"
)

type context struct {
  selectedMap [][]int
  mapIndex int
  hRatio int
  vRatio int
  blankLineCount int
}

func (ctx *context) readMapLine() []int {
  if ctx.blankLineCount < ctx.vRatio {
    ctx.blankLineCount++
    return []int{}
  }
  ctx.blankLineCount = 0
  mapLine := ctx.selectedMap[ctx.mapIndex]
  ctx.mapIndex++
  if ctx.mapIndex == len(ctx.selectedMap) {
    ctx.mapIndex = 0
  }
  return mapLine
}

func main() {
  startGame()
}

func startGame() {
  ctx := newGame()
  ticker := time.NewTicker(500 * time.Millisecond)
  for ; ; <-ticker.C {
    render(&ctx)
  }
}

func render(ctx *context) {
  mapLine := ctx.readMapLine()
  if len(mapLine) == 0 {
    fmt.Println()
    return
  }
  for _, item := range mapLine {
    switch item {
      case 0:
        printSpace(ctx.hRatio)
      case 1:
        fmt.Print(item)
    }
  }
  fmt.Println()
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
  return context{
    selectedMap: getMap(),
    hRatio: 5,
    vRatio: 3,
  }
}

func getMap() [][]int {
  return [][]int{
    []int{0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0},
    []int{0, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 1, 0},
    []int{0, 0, 0, 0, 0, 1, 1, 0, 1, 1, 0, 0, 0, 0, 0},
    []int{0, 1, 1, 1, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 0},
  }
}

