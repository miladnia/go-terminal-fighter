package main

import "strings"

type point struct {
  x, y int
}

type canvas struct {
  width  int
  height int
  xRatio int
  yRatio int
  grid   [][]rune
}

func newCanvas(width, height, xRatio, yRatio int) *canvas {
  c := &canvas{
    width:  width * xRatio,
    height: height * yRatio,
    xRatio: xRatio,
    yRatio: yRatio,
  }
  c.grid = make([][]rune, c.height)
  return c
}

func (c *canvas) clean() {
  c.grid = make([][]rune, c.height)
}

func (c *canvas) draw(char rune, p point, rowsToMask int) {
  a := point{
    x: p.x * c.xRatio,
    y: p.y * c.yRatio - rowsToMask,
  }
  b := point{
    x: a.x + c.xRatio - 1,
    y: a.y + c.yRatio - 1,
  }
  c.drawBlock(char, a, b)
}

func (c *canvas) drawBlock(char rune, a, b point) {
  for y := a.y; y <= b.y; y++ {
    if y < 0 {
      continue
    }
    if len(c.grid[y]) == 0 {
      c.grid[y] = make([]rune, c.width)
    }
    for x := a.x; x <= b.x; x++ {
      c.grid[y][x] = char
    }
  }
}

func (c *canvas) drawSprite(f [][]rune, p point) {
  if len(f) != yRatio {
    return
  }
  p.x *= xRatio
  p.y *= yRatio
  for i, row := range f {
    if len(c.grid[p.y + i]) == 0{
      c.grid[p.y + i] = make([]rune, c.width)
    }
    for j, char := range row {
      c.grid[p.y + i][p.x + j] = char
    }
  }
}

func (c *canvas) toString() string {
  var sb strings.Builder
  // By masking the last rows, only one line
  // of characters will be shown in the last row,
  // and when we have 'rowsToMask' (greater than 0)
  // the last row will be disapear line by line smoothly.
  rowsToIgnore := yRatio - 1
  h := len(c.grid) - rowsToIgnore
  g := c.grid[:h]
  for _, row := range g {
    if len(row) == 0 {
      // Create a blank line.
      for i := 0; i < c.width; i++ {
        sb.WriteRune(' ')
      }
      sb.WriteRune('\n')
      continue
    }
    for _, char := range row {
      if char == 0 {
        char = ' '
      }
      sb.WriteRune(char)
    }
    sb.WriteRune('\n')
  }
  return sb.String()
}

