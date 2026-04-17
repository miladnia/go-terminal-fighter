package main

import "fmt"

type fighter struct {
  x, y int
}

func (f *fighter) reachedEnd() bool {
  return f.y <= 0
}

func (f *fighter) moveForward() {
  if f.reachedEnd() {
    return
  }
  f.y--
}

func (f *fighter) moveLeft() {
  f.x--
}

func (f *fighter) moveRight() {
  f.x++
}

func (f *fighter) getSprite(w, h int) ([][]rune, error) {
  if w == 5 && h == 3 {
    //   ▲
    // ■■█■■
    //  ■ ■
    // (5 x 3)
    return [][]rune{
      []rune{' ', ' ', '▲', ' ', ' '},
      []rune{'■', '■', '█', '■', '■'},
      []rune{' ', '■', ' ', '■', ' '},
    }, nil
  }

  return nil, fmt.Errorf("fighter: no (%v x %v) sprite.")
}

