package main

type gameMap struct {
  grid [][]int
}

func (m *gameMap) getChunk(a, b int) (chunk [][]int, ok bool) {
  if a > b {
    return nil, false
  }
  chunkLength := b - a + 1
  chunk = make([][]int, chunkLength)
  lastIndex := len(m.grid) - 1
  if a > lastIndex || b < 0 {
    return chunk, true
  }
  ci := 0
  for i := a; i <= b; i++ {
    if 0 <= i && i <= lastIndex {
      chunk[ci] = m.grid[i]
    }
    ci++
  }
  return chunk, true
}

