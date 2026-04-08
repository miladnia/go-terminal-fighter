package main

func main() {
  clearScreen()
  klgr := newKeyLogger()
  for {
    key := ask(klgr, []menuOption{
      {key: 's', title: "Start"},
      {key: 'e', title: "Exit"},
    })
    if key == 's' {
      start(klgr)
    } else {
      break
    }
  }
  flashMessage("GOOD GAME!")
}

func start(klgr *keyLogger) {
  selectedMap := getMap()
  game := newGame(selectedMap)
  writer := terminalWriter{}
  eng := newEngine(&writer, klgr, map[byte]string{
    'a': "left",
    'd': "right",
    'q': "pause",
  })
  eng.startGame(game, 2)

  for {
    select {
    case <-eng.gameOver:
      showInfo([]string{
        "game over!",
        "",
      })
      return
    case <-eng.gamePaused:
      key := ask(klgr, []menuOption{
        {key: 'r', title: "Resume"},
        {key: 'q', title: "Quit"},
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

