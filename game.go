package main

import "sync"

const (
  width              = 15
  height             = 8
  xRatio             = 5
  yRatio             = 3
  speedMin           = 1
  speedMax           = 10
  fighterDefaultX    = 7
  fighterBottomSpace = 1
)

type game struct {
  state       gameState
  gameMap     gameMap
  radar       [][]int
  canvas      *canvas
  fighter     fighter
  steps       int
  renderSteps int
  stepsDelay  int
  mux         sync.Mutex
}

func newGame(mapGrid [][]int, speed int) playable {
  if speed < speedMin {
    speed = speedMin
  } else if speed > speedMax {
    speed = speedMax
  }
  g := &game{
    gameMap:    gameMap{grid: mapGrid},
    stepsDelay: speedMax - speed + 1,
    fighter:    fighter{x: fighterDefaultX, y: len(mapGrid) - 1 + height},
    radar:      make([][]int, height),
    canvas:     newCanvas(width, height, xRatio, yRatio),
    mux:        sync.Mutex{},
  }
  return g
}

func (g *game) step() (gameOver bool) {
  if g.state.gameOver {
    return true
  }
  g.steps++
  if g.steps % g.stepsDelay != 0 {
    return
  }
  g.renderSteps++
  if g.renderSteps % yRatio != 0 {
    return
  }
  g.moveForward()
  return g.state.gameOver
}

func (g *game) moveForward() {
  g.mux.Lock()
  defer g.mux.Unlock()
  g.fighter.moveForward()
  mapIndex := g.fighter.y
  mapChunk, _ := g.gameMap.getChunk(mapIndex - height + 1, mapIndex)
  g.radar = mapChunk
  g.updateState()
}

func (g *game) render() string {
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
      rowsToMask := yRatio - (g.renderSteps % yRatio) - 1
      g.canvas.draw(char, point{x: j, y: i}, rowsToMask)
    }
  }

  // Draw fighter.
  g.canvas.drawSprite(&g.fighter, point{
    x: g.fighter.x,
    y: len(g.radar) - 1 - fighterBottomSpace,
  })

  return g.canvas.toString()
}

func (g *game) moveLeft() {
  g.mux.Lock()
  defer g.mux.Unlock()
  if g.fighter.x == 0 {
    return
  }
  g.fighter.moveLeft()
  g.updateState()
}

func (g *game) moveRight() {
  g.mux.Lock()
  defer g.mux.Unlock()
  if g.fighter.x == width - 1 {
    return
  }
  g.fighter.moveRight()
  g.updateState()
}

func (g *game) getState() gameState {
  return g.state
}

func (g *game) updateState() {
  g.state.won = g.fighter.reachedEnd()
  g.state.gameOver = g.state.won || g.hasCollided()
}

func (g *game) hasCollided() bool {
  fighterRow := g.radar[len(g.radar) - 1 - fighterBottomSpace]
  collided := len(fighterRow) > 0 && fighterRow[g.fighter.x] == 1
  return collided
}

