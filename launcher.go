package main

import (
  "fmt"
  "time"
)

type state struct {
  lvl level
}

type launcher struct {
  d dialogue
  e *engine
  p *profile
  s state
}

func newLauncher() *launcher {
  writer := &terminalWriter{}
  keyLogger := newKeyLogger()
  engine := newEngine(writer, keyLogger, map[byte]string{
    'a': "left",
    'd': "right",
    'q': "pause",
  })
  return &launcher{
    d: dialogue{klgr: keyLogger},
    e: engine,
  }
}

func (l *launcher) init() {
  hideCursor()
  defer showCursor()
  for {
    key := l.d.cleanAsk([]option{
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
    case gameState := <-l.e.gameOver:
      if gameState.won {
        if l.isFinalLevel() {
          showFinalMessage()
          l.d.askAnyKey()
          return
        }
        l.levelUp()
        showLevelUpMessage()
      } else {
        if !gameState.won {
          time.Sleep(1000 * time.Millisecond)
        }
        showGameOverMessage()
      }
      key := l.d.askAndClear([]option{
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
      key := l.d.askTransparent([]option{
        {key: 'r', title: "Resume"},
        {key: 'q', title: "Quit"},
      })
      if key == 'q' {
        l.e.stopGame()
        return
      }
      if key == 'r' {
        time.Sleep(1500 * time.Millisecond)
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
  game := newGame(lvl.getMap(), lvl.Speed)
  l.e.startGame(game)
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

