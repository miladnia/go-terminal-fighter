package main

import (
  "fmt"
  "strings"
)

type context struct {
  selectedMap    [][]int
  mapIndex       int
  blankRowsCount int
  radar          [][]int
  playerPosition int
  width          int
  height         int
  steps          int
  v              *viewport
  xRatio         int
  yRatio         int
}

func newGame(selectedMap [][]int) game {
  ctx := &context{
    selectedMap:    selectedMap,
    mapIndex:       len(selectedMap) - 1,
    height:         8,
    width:          15,
    playerPosition: 7,
    xRatio:         5,
    yRatio:         3,
  }
  // Init radar.
  ctx.radar = make([][]int, ctx.height)
  for i := range ctx.radar {
    ctx.radar[i] = make([]int, ctx.width)
  }
  ctx.v = newViewPort(ctx.width * ctx.xRatio, ctx.height * ctx.yRatio)
  ctx.steps = ctx.yRatio
  return ctx
}

func (ctx *context) step() (end bool, won bool) {
  if ctx.steps > 0 {
    ctx.steps--
    return false, false
  }
  ctx.steps = ctx.yRatio
  return ctx.moveForward()
}

func (ctx *context) moveForward() (end bool, won bool) {
  mapRow, done := ctx.readMap()
  if done {
    return true, true
  }
  // Remove the last row of the radar
  // and prepend a new row.
  shiftDown(ctx.radar, mapRow)
  // Reposition the player in
  // the new last row of the radar.
  lastRow := ctx.radar[len(ctx.radar) - 1]
  // Hit?
  if lastRow[ctx.playerPosition] == 1 {
    end = true
  }
  lastRow[ctx.playerPosition] = 2
  ctx.radar[len(ctx.radar) - 1] = lastRow
  return end, false
}

func (ctx *context) render() string {
  return ctx.render2DMode()
}

func (ctx *context) render2DMode() string {
  for y, row := range ctx.radar {
    for x, item := range row {
      // Fighter aircraft?
      if item == 2 {
        from := point{x: ctx.xRatio * x, y: ctx.yRatio * y}
        to := point{x: from.x + ctx.xRatio, y: from.y + ctx.yRatio}
        ctx.v.drawFighter(from, to)
        continue
      }
      char := ' '
      switch item {
      case 1:
        char = '╬'
      case 9:
        char = '□'
      }
      from := point{x: ctx.xRatio * x, y: ctx.yRatio * y}
      to := point{x: from.x + ctx.xRatio, y: from.y + ctx.yRatio - ctx.steps}
      ctx.v.drawBlock(char, from, to)
    }
  }
  return ctx.v.toString()
}

func (ctx *context) renderRadarMode() string {
  var canvas string
  for _, row := range ctx.radar {
    for _, item := range row {
      switch item {
        case 0:
          canvas = fmt.Sprint(canvas, " ")
        case 1:
          canvas = fmt.Sprint(canvas, "#")
        case 2:
          canvas = fmt.Sprint(canvas, "^")
        case 9:
          canvas = fmt.Sprint(canvas, "*")
      }
    }
    canvas = fmt.Sprintln(canvas)
  }
  return canvas
}

func (ctx *context) readMap() (mapRow []int, done bool) {
  // The whole map have been read.
  if ctx.mapIndex < 0 {
    // Return blank rows.
    if ctx.blankRowsCount < ctx.height {
      ctx.blankRowsCount++
      return make([]int, ctx.width), false
    }
    return nil, true
  }
  i := ctx.mapIndex
  ctx.mapIndex--
  mapRow = ctx.selectedMap[i]
  if len(mapRow) == 0 {
    mapRow = make([]int, ctx.width)
  }
  return mapRow, false
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

func shiftDown[T any](s []T, newElement T) {
  for i := len(s) - 1; i > 0; i-- {
    s[i] = s[i - 1]
  }
  s[0] = newElement
}

type point struct {
  x, y int
}

type viewport struct {
  width  int
  height int
  canvas [][]rune
}

func newViewPort(width, height int) *viewport {
  return &viewport{
    width:  width,
    height: height,
    canvas: make([][]rune, height),
  }
}

func (v *viewport) drawBlock(char rune, a, b point) {
  for y := a.y; y < b.y; y++ {
    if len(v.canvas[y]) == 0 {
      v.canvas[y] = make([]rune, v.width)
    }
    for x := a.x; x < b.x; x++ {
      v.canvas[y][x] = char
    }
  }
}

func (v *viewport) drawFighter(a, b point) {
  for y := a.y; y < b.y; y++ {
    if len(v.canvas[y]) == 0 {
      v.canvas[y] = make([]rune, v.width)
    }
    // Clean the block for the fighter.
    for x := a.x; x < b.x; x++ {
      v.canvas[y][x] = ' '
    }
  }
  //   ▲
  // ■■█■■
  //  ■ ■
  v.canvas[a.y][a.x + 2] = '▲'
  v.canvas[a.y + 1][a.x] = '■'
  v.canvas[a.y + 1][a.x + 1] = '■'
  v.canvas[a.y + 1][a.x + 2] = '█'
  v.canvas[a.y + 1][a.x + 3] = '■'
  v.canvas[a.y + 1][a.x + 4] = '■'
  v.canvas[a.y + 2][a.x + 1] = '■'
  v.canvas[a.y + 2][a.x + 3] = '■'
}

func (v *viewport) toString() string {
  var sb strings.Builder
  for _, row := range v.canvas {
    for _, char := range row {
      sb.WriteRune(char)
    }
    sb.WriteRune('\n')
  }
  return sb.String()
}

