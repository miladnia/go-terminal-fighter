package main

import "fmt"

type context struct {
  selectedMap [][]int
  mapIndex int
  blankLineCount int
  radar [][]int
  playerPosition int
  width int
  height int
}

func newGame(selectedMap [][]int) game {
  ctx := context{
    selectedMap: selectedMap,
    mapIndex: len(selectedMap) - 1,
    height: 20,
    width: 15,
    playerPosition: 7,
  }
  ctx.radar = make([][]int, ctx.height)
  return &ctx
}

func (ctx *context) step() (end bool, won bool) {
  return ctx.moveForward()
}

func (ctx *context) render() string {
  var canvas string
  for i := len(ctx.radar) - 1; i >= 0; i-- {
    line := ctx.radar[i]
    for _, item := range line {
      switch item {
        case 0:
          canvas  = fmt.Sprint(canvas, " ")
        case 1:
          canvas = fmt.Sprint(canvas, "*")
        case 2:
          canvas = fmt.Sprint(canvas, "^")
        case 9:
          canvas = fmt.Sprint(canvas, "#")
      }
    }
    canvas = fmt.Sprintln(canvas)
  }
  return canvas
}

func (ctx *context) moveForward() (end bool, won bool) {
  mapLine, done := ctx.readMap()
  if done {
    return true, true
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
  // Hit?
  if firstLine[ctx.playerPosition] == 1 {
    end = true
  }
  firstLine[ctx.playerPosition] = 2
  ctx.radar[0] = firstLine
  return end, false
}

func (ctx *context) readMap() (mapLine []int, done bool) {
  // The whole map have been read.
  if ctx.mapIndex < 0 {
    // Return blank lines.
    if ctx.blankLineCount < ctx.height {
      ctx.blankLineCount++
      return []int{}, false
    }
    return nil, true
  }
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

