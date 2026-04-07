package main

import (
  "os"
//  "time"
)

func main() {
  for {
    key := ask(map[byte]string{
      's': "Start",
      'e': "Exit",
    })
    if key == 's' {
      start()
    } else {
      break
    }
  }
  flashMessage("GOOD GAME!")
}

func start() {
  selectedMap := getMap()
  game := newGame(selectedMap)
  eng := newEngine(os.Stdout, map[byte]string{
    'a': "left",
    'd': "right",
    'q': "pause",
  })
  eng.startGame(game, 10)

  for {
    select {
    case <-eng.gameOver:
      showInfo([]string{
        "game over!",
        "",
      })
      return
    case <-eng.gamePaused:
      key := ask(map[byte]string{
        'r': "Resume",
        'q': "Quit",
      })
      if key == 'r' {
        eng.resumeGame()
      } else {
        return
      }
    }
  }
}

func getMap() [][]int {
  return [][]int{
    []int{0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0},
    []int{0, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 1, 0},
    []int{0, 0, 0, 0, 0, 1, 1, 0, 1, 1, 0, 0, 0, 0, 0},
    []int{0, 1, 1, 1, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 0},
    []int{0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0},
    []int{0, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 1, 0},
    []int{0, 0, 0, 0, 0, 1, 1, 0, 1, 1, 0, 0, 0, 0, 0},
    []int{0, 1, 1, 1, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 0},
    []int{0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0},
    []int{0, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 1, 0},
    []int{0, 0, 0, 0, 0, 1, 1, 0, 1, 1, 0, 0, 0, 0, 0},
    []int{0, 1, 1, 1, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 0},
    []int{0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0},
    []int{0, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 1, 0},
    []int{0, 0, 0, 0, 0, 1, 1, 0, 1, 1, 0, 0, 0, 0, 0},
    []int{0, 1, 1, 1, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 0},
  }
}

