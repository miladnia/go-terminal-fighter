package main

import (
  "sync"
)

const (
  width  = 15
  height = 8
  xRatio = 5
  yRatio = 3
)

type game struct {
  status         gameStatus
  gameMap        gameMap
  radar          [][]int
  canvas         *canvas
  fighter        fighter
  mux            sync.RWMutex
  steps          int
}

func newGame(mapGrid [][]int) playable {
  g := &game{
    mux: sync.RWMutex{},
  }
  g.gameMap = gameMap{
    grid: mapGrid,
  }
  g.fighter = fighter{
    x: 7,
    y: len(mapGrid) + height,
  }
  g.canvas = newCanvas(width, height, xRatio, yRatio)
  return g
}

func (g *game) step() (gameOver bool) {
  g.mux.Lock()
  defer g.mux.Unlock()
  g.steps++
  if g.steps < yRatio {
    return false
  }
  g.steps = 0
  return g.moveForward()
}

func (g *game) render() string {
  g.mux.Lock()
  defer g.mux.Unlock()
  return g.render2DMode()
}

func (g *game) moveForward() (gameOver bool) {
  if g.status.gameOver {
    return true
  }

  ok := g.fighter.moveForward()
  if !ok {
    g.status.gameOver = true
    g.status.won = true
    return true
  }
  mapIndex := g.fighter.y
  mapChunk, _ := g.gameMap.getChunk(mapIndex - height + 1, mapIndex)
  g.radar = mapChunk

  fighterRow := g.radar[len(g.radar) - 2]
  // Hit?
  if len(fighterRow) > 0 && fighterRow[g.fighter.x] == 1 {
    g.status.gameOver = true
    g.status.won = false
  }
  return g.status.gameOver
}

func (g *game) render2DMode() string {
  g.canvas.clean()
  for i, row := range g.radar {
    // Empty row?
    if len(row) == 0 {
      continue
    }
    for j, item := range row {
      var char rune
      switch item {
      case 1:
        char = '╬'
      case 9:
        char = '□'
      default:
        continue
      }
      rowsToMask := yRatio - g.steps - 1
      g.canvas.draw(char, point{x: j, y: i}, rowsToMask)
    }
  }

  //   ▲
  // ■■█■■
  //  ■ ■
  // (5 x 3)
  f := [][]rune{
    []rune{' ', ' ', '▲', ' ', ' '},
    []rune{'■', '■', '█', '■', '■'},
    []rune{' ', '■', ' ', '■', ' '},
  }
  p := point{x: g.fighter.x, y: height - 2}
  g.canvas.drawSprite(f, p)

  return g.canvas.toString()
}

func (g *game) moveLeft() {
  if g.fighter.x == 0 {
    return
  }
  g.fighter.moveLeft()
}

func (g *game) moveRight() {
  if g.fighter.x == width - 1 {
    return
  }
  g.fighter.moveRight()
}

func (g *game) getStatus() gameStatus {
  return g.status
}

