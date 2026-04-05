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
  radar [][]int
  height int
  width int
  speed int
  gameRunning bool
  playerPosition int
}

func (ctx *context) moveForward() (end bool) {
    mapLine, end := ctx.readMap()
    if end {
      return true
    }
    // Remove the first line of the radar
    // and append a new line.
    ctx.radar = ctx.radar[1:]
    ctx.radar = append(ctx.radar, mapLine)
    // Reposition the player in
    // the new first line of the radar.
    firstLine := ctx.radar[0]
    if len(firstLine) == 0 {
      firstLine = make([]int, ctx.width)
    }
    firstLine[ctx.playerPosition] = 2
    ctx.radar[0] = firstLine
    return false
}

func (ctx *context) readMap() (mapLine []int, end bool) {
  // The whole map have been read.
  if ctx.mapIndex <= 0 {
    // Return blank lines.
    if ctx.blankLineCount < ctx.height {
      ctx.blankLineCount++
      return []int{}, false
    }
    return nil, true
  }
  // Add blank lines between map lines.
  if ctx.blankLineCount < ctx.vRatio {
    ctx.blankLineCount++
    return []int{}, false
  }
  ctx.blankLineCount = 0
  i := ctx.mapIndex
  ctx.mapIndex--
  mapLine = ctx.selectedMap[i]
  return mapLine, false
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

func (ctx *context) runCommand(cmd string) {
  switch cmd {
  case "left":
    ctx.moveLeft()
  case "right":
    ctx.moveRight()
  }
}

func startGame(ctx *context) {
  ctx.gameRunning = true
  go func() {
    ticker := time.NewTicker(50 * time.Millisecond)
    for ; ctx.gameRunning; <-ticker.C {
      render(ctx)
    }
  }()
  go func() {
    delayMS := time.Duration(500 / ctx.speed)
    ticker := time.NewTicker(delayMS * time.Millisecond)
    for ; ctx.gameRunning; <-ticker.C {
      end := ctx.moveForward()
      if end {
        ctx.gameRunning = false
        break
      }
    }
  }()
  for ctx.gameRunning {
    key := captureInput()
    switch key {
    case 'q':
      ctx.gameRunning = false
    default:
      cmd := translateKeyToGameCommand(key)
      ctx.runCommand(cmd)
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

func translateKeyToGameCommand(key byte) (cmd string) {
  switch key {
  case 'a':
    cmd = "left"
  case 'd':
    cmd = "right"
  }
  return cmd
}

func newGame() context {
  selectedMap := getMap()
  ctx := context{
    selectedMap: selectedMap,
    mapIndex: len(selectedMap) - 1,
    hRatio: 1,
    vRatio: 3,
    height: 20,
    width: 15,
    speed: 2,
    playerPosition: 7,
  }
  ctx.radar = make([][]int, ctx.height)
  return ctx
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
  ctx := newGame()
  for {
    startGame(&ctx)
    printBlankLine(5)
    fmt.Println("GAME OVER!")
    fmt.Println()
    key := askToChoose(map[byte]string{
      'r': "Resume",
      'n': "New Game",
      'q': "Quit",
    })
    if key == 'n' {
      ctx = newGame()
    } else if key == 'r' {
      continue
    } else if key == 'q' {
      break
    }
  }
  flashMessage("GOOD GAME!")
}

