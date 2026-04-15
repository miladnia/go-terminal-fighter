package main

const (
  width    = 15
  height   = 8
  xRatio   = 5
  yRatio   = 3
  speedMin = 1
  speedMax = 10
  fighterX = 7
)

type game struct {
  status      gameStatus
  gameMap     gameMap
  radar       [][]int
  canvas      *canvas
  fighter     fighter
  steps       int
  renderSteps int
  stepsDelay  int
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
    fighter:    fighter{x: fighterX, y: len(mapGrid) - 1 + height},
    canvas:     newCanvas(width, height, xRatio, yRatio),
  }
  return g
}

func (g *game) step() (gameOver bool) {
  g.steps++
  if g.steps % g.stepsDelay != 0 {
    return
  }
  g.renderSteps++
  if g.renderSteps % yRatio != 0 {
    return
  }
  return g.moveForward()
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

