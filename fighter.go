package main

type fighter struct {
  x, y int
}

func (f *fighter) moveForward() (ok bool) {
  if f.y <= 0 {
    return false
  }
  f.y--
  return true
}

func (f *fighter) moveLeft() {
  f.x--
}

func (f *fighter) moveRight() {
  f.x++
}

