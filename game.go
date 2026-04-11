package main

import (
  "strings"
  "sync"
)

type game struct {
  status         gameStatus
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
  mux            *sync.RWMutex
}

func newGame(selectedMap [][]int) playable {
  g := &game{
    selectedMap:    selectedMap,
    mapIndex:       len(selectedMap) - 1,
    height:         8,
    width:          15,
    playerPosition: 7,
    xRatio:         5,
    yRatio:         3,
    mux:            &sync.RWMutex{},
  }
  // Init radar.
  g.radar = make([][]int, g.height)
  for i := range g.radar {
    g.radar[i] = make([]int, g.width)
  }
  g.v = newViewPort(g.width * g.xRatio, g.height * g.yRatio)
  g.steps = g.yRatio
  return g
}

func (g *game) step() (gameOver bool) {
  g.mux.Lock()
  defer g.mux.Unlock()
  if g.steps > 0 {
    g.steps--
    return false
  }
  g.steps = g.yRatio
  return g.moveForward()
}

func (g *game) render() string {
  g.mux.RLock()
  defer g.mux.RUnlock()
  return g.render2DMode()
}

func (g *game) moveForward() (gameOver bool) {
  mapRow, done := g.readMap()
  if done {
    g.status.gameOver = true
    g.status.won = true
    return true
  }
  // Remove the last row of the radar
  // and prepend a new row.
  shiftDown(g.radar, mapRow)
  // Reposition the player in
  // the new last row of the radar.
  lastRow := g.radar[len(g.radar) - 1]
  // Hit?
  if lastRow[g.playerPosition] == 1 {
    g.status.gameOver = true
    g.status.won = false
  }
  lastRow[g.playerPosition] = 2
  g.radar[len(g.radar) - 1] = lastRow
  return g.status.gameOver
}

func (g *game) render2DMode() string {
  for y, row := range g.radar {
    for x, item := range row {
      from := point{x: g.xRatio * x, y: g.yRatio * y}
      to := point{x: from.x + g.xRatio, y: from.y + g.yRatio}
      // Fighter aircraft?
      if item == 2 {
        g.v.drawFighter(from, to)
        continue
      }
      char := ' '
      switch item {
      case 1:
        char = '╬'
      case 9:
        char = '□'
      }
      to.y -= g.steps // show the blocks row by row
      g.v.drawBlock(char, from, to)
    }
  }
  return g.v.toString()
}

func (g *game) renderRadarMode() string {
  var sb strings.Builder
  for _, row := range g.radar {
    for _, item := range row {
      switch item {
        case 0:
          sb.WriteRune(' ')
        case 1:
          sb.WriteRune('#')
        case 2:
          sb.WriteRune('^')
        case 9:
          sb.WriteRune('*')
      }
    }
    sb.WriteRune('\n')
  }
  return sb.String()
}

func (g *game) readMap() (mapRow []int, done bool) {
  // The whole map have been read.
  if g.mapIndex < 0 {
    // Return blank rows.
    if g.blankRowsCount < g.height {
      g.blankRowsCount++
      return make([]int, g.width), false
    }
    return nil, true
  }
  i := g.mapIndex
  g.mapIndex--
  mapRow = g.selectedMap[i]
  if len(mapRow) == 0 {
    mapRow = make([]int, g.width)
  }
  return mapRow, false
}

func (g *game) moveLeft() {
  if g.playerPosition == 0 {
    return
  }
  g.playerPosition--
}

func (g *game) moveRight() {
  if g.playerPosition == g.width - 1 {
    return
  }
  g.playerPosition++
}

func (g *game) getStatus() gameStatus {
  return g.status
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

func shiftDown[T any](s []T, newElement T) {
  for i := len(s) - 1; i > 0; i-- {
    s[i] = s[i - 1]
  }
  s[0] = newElement
}

