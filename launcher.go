package main

import "fmt"

type state struct {
  lvl level
}

type launcher struct {
  k *keyLogger
  e *engine
  p *profile
  s state
}

func newLauncher() *launcher {
  writer := terminalWriter{}
  keyLogger := newKeyLogger()
  engine := newEngine(&writer, keyLogger, map[byte]string{
    'a': "left",
    'd': "right",
    'q': "pause",
  })
  return &launcher{
    k: keyLogger,
    e: engine,
  }
}

func (l *launcher) init() {
  clearScreen()
  hideCursor()
  defer showCursor()
  for {
    key := ask(l.k, []menuOption{
      {key: 'n', title: "New Game"},
      {key: 'c', title: "Continue", disabled: l.p == nil},
      {key: 'e', title: "Exit"},
    })
    if key == 'e' {
      return
    }
    if key == 'n' {
      l.p = &profile{
        levelNo: 1,
        playingLevelNo: 1,
      }
    }
    err := l.launch()
    if err != nil {
      fmt.Println(err)
      return
    }
  }
}

func (l *launcher) launch() (err error) {
  err = l.startGame()
  if err != nil {
    return
  }

  for {
    select {
    case status := <-l.e.gameOver:
      if status.won {
        if l.isFinalLevel() {
          showFinalMessage()
          return
        }
        l.levelUp()
        showLevelUpMessage()
      } else {
        showGameOverMessage()
      }
      key := ask(l.k, []menuOption{
        {key: 'c', title: "Continue"},
        {key: 'q', title: "Quit"},
      })
      if key == 'q' {
        return
      }
      err = l.startGame()
      if err != nil {
        return
      }
    case <-l.e.gamePaused:
      goToScreenTop()
      key := ask(l.k, []menuOption{
        {key: 'r', title: "Resume"},
        {key: 'q', title: "Quit"},
      })
      if key == 'q' {
        return
      }
      if key == 'r' {
        l.e.resumeGame()
      }
    }
  }
}

func (l *launcher) startGame() error {
  levelNo := l.p.playingLevelNo
  lvl, err := getLevel(levelNo)
  if err != nil {
    return err
  }
  l.s.lvl = lvl
  game := newGame(lvl.getMap())
  l.e.startGame(game, lvl.Speed)
  return nil
}

func (l *launcher) isFinalLevel() bool {
  return l.s.lvl.IsFinal
}

func (l *launcher) levelUp() {
  l.p.playingLevelNo++
  if l.p.playingLevelNo > l.p.levelNo {
    l.p.levelNo = l.p.playingLevelNo
  }
}

func showGameOverMessage() {
  showInfo([]string{
    "GAME OVER!",
    "YOU HAD A HIT, BUT DON'T GIVE UP!",
    "LET'S PLAY ONE MORE TIME...",
  })
}

func showLevelUpMessage() {
  showInfo([]string{
    "GOOD GAME :)",
    "YOU PASSED THE LEVEL!",
    "LET'S JUMP TO THE NEXT CHALLANGE!",
  })
}

func showFinalMessage() {
  showInfo([]string{
    "OMG!",
    "YOU DID IT!",
    "YOU FINISHED THE GAME!",
  })
}

