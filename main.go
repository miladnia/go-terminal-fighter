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
  radar [][]int
  height int
  width int
  speed int
  gameRunning bool
  playerPosition int
}

func (ctx *context) updateRadar(newLine []int) {
    // Remove the first line of the radar
    // and append a new line.
    ctx.radar = ctx.radar[1:]
    ctx.radar = append(ctx.radar, newLine)
    // Reposition the player in
    // the new first line of the radar.
    fisrtLine := ctx.radar[0]
    if len(fisrtLine) == 0 {
      fisrtLine = make([]int, ctx.width)
    }
    fisrtLine[ctx.playerPosition] = 2
    ctx.radar[0] = fisrtLine
}

func (ctx *context) moveForward() (end bool) {
  // The whole map have been read.
  if ctx.mapIndex <= 0 {
    if ctx.blankLineCount < ctx.height {
      ctx.blankLineCount++
      ctx.updateRadar([]int{})
      return false
    }
    return true
  }
  // Add blank lines between map lines.
  if ctx.blankLineCount < ctx.vRatio {
    ctx.blankLineCount++
    ctx.updateRadar([]int{})
    return false
  }
  ctx.blankLineCount = 0
  i := ctx.mapIndex
  ctx.mapIndex--
  mapLine := ctx.selectedMap[i]
  ctx.updateRadar(mapLine)
  return false
}

func (ctx *context) submitControllerAction(key byte) {
  if 'A' <= key && key <= 'Z' {
    key += 'a' - 'A'
  }
  switch key {
  case 'a':
    ctx.moveLeft()
  case 'd':
    ctx.moveRight()
  }
}

func (ctx *context) moveLeft() {
  if ctx.playerPosition == 0 {
    return
  }
  ctx.playerPosition--
}

func (ctx *context) moveRight() {
  if ctx.playerPosition == ctx.width - 1 {
    return
  }
  ctx.playerPosition++
}

func startGame() {
  ctx := newGame()
  ctx.gameRunning = true
  go captureInput(func(key byte) bool {
    ctx.submitControllerAction(key)
    return !ctx.gameRunning
  })
  go func() {
    ticker := time.NewTicker(50 * time.Millisecond)
    for ; ctx.gameRunning; <-ticker.C {
      render(ctx)
    }
  }()
  delayMS := time.Duration(500 / ctx.speed)
  ticker := time.NewTicker(delayMS * time.Millisecond)
  for ; ; <-ticker.C {
    end := ctx.moveForward()
    if end {
      ctx.gameRunning = false
      break
    }
  }
}

func render(ctx *context) {
  var canvas string
  for i := len(ctx.radar) - 1; i >= 0; i-- {
    line := ctx.radar[i]
    for _, item := range line {
      switch item {
        case 0:
          canvas  = printSpace(canvas, ctx.hRatio)
        case 1:
          canvas = fmt.Sprint(canvas, "*")
        case 2:
          canvas = fmt.Sprint(canvas, "^")
      }
    }
    canvas = fmt.Sprintln(canvas)
  }
  clearScreen()
  fmt.Print(canvas)
}

func newGame() *context {
  selectedMap := getMap()
  ctx := context{
    selectedMap: selectedMap,
    mapIndex: len(selectedMap) - 1,
    hRatio: 1,
    vRatio: 3,
    radar: [][]int{},
    height: 20,
    width: 15,
    speed: 2,
    playerPosition: 7,
  }
  ctx.radar = make([][]int, ctx.height)
  return &ctx
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

func main() {
  for {
    startGame()
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
      break
    }
  }
}

