package main

import (
  "io"
  "fmt"
  "time"
)

type game interface {
  render() string
  step() (end bool)
  moveLeft()
  moveRight()
}

type engine struct {
  w io.Writer
  klgr *keyLogger
  controllerSetup map[byte]string
  g game
  movementInterval time.Duration
  renderInterval time.Duration
  renderTicker *time.Ticker
  movementTicker *time.Ticker
  gamePaused chan struct{}
  gameOver chan struct{}
  renderingDone chan struct{}
  controllerDone chan struct{}
}

func newEngine(w io.Writer, klgr *keyLogger, controllerSetup map[byte]string) *engine {
  return &engine{
    w: w,
    klgr: klgr,
    controllerSetup: controllerSetup,
    gamePaused: make(chan struct{}),
    gameOver: make(chan struct{}),
    renderingDone: make(chan struct{}),
    controllerDone: make(chan struct{}),
  }
}

func (eng *engine) startGame(g game, speed int) {
  eng.g = g
  eng.movementInterval = time.Duration(500 / speed) * time.Millisecond
  eng.renderInterval = 42 * time.Millisecond
  go eng.startRendering()
  go eng.startMovement()
  go eng.listenToController()
}

func (eng *engine) pauseGame() {
  eng.pauseMovement()
  eng.pauseRendering()
  eng.gamePaused <- struct{}{}
}

func (eng *engine) resumeGame() {
  go eng.listenToController()
  eng.resumeRendering()
  time.Sleep(1500 * time.Millisecond)
  eng.resumeMovement()
}

func (eng *engine) stopGame() {
  eng.renderingDone <- struct{}{}
  eng.controllerDone <- struct{}{}
  eng.gameOver <- struct{}{}
}

func (eng *engine) startMovement() {
  eng.movementTicker = time.NewTicker(eng.movementInterval)
  for ; ; <-eng.movementTicker.C {
    end := eng.g.step()
    if end {
      eng.stopGame()
      break
    }
  }
}

func (eng *engine) pauseMovement() {
  eng.movementTicker.Stop()
}

func (eng *engine) resumeMovement() {
  eng.movementTicker.Reset(eng.movementInterval)
}

func (eng *engine) startRendering() {
  eng.renderTicker = time.NewTicker(eng.renderInterval)
  for {
    select {
    case <-eng.renderingDone:
      eng.renderTicker.Stop()
      return
    case <-eng.renderTicker.C:
      frame := eng.g.render()
      clearScreen()
      fmt.Fprintln(eng.w, frame)
    }
  }
  for ; ; <-eng.renderTicker.C {
    frame := eng.g.render()
    clearScreen()
	  fmt.Fprintln(eng.w, frame)
  }
}

func (eng *engine) resumeRendering() {
  eng.renderTicker.Reset(eng.renderInterval)
}

func (eng *engine) pauseRendering() {
  eng.renderTicker.Stop()
  time.Sleep(100 * time.Millisecond)
}

func (eng *engine) listenToController() {
  for {
    select {
    case <-eng.controllerDone:
      return
    case key := <-eng.klgr.C:
      command, ok := eng.controllerSetup[key]
      if !ok {
        continue
      }
      switch command {
      case "left":
        eng.g.moveLeft()
      case "right":
        eng.g.moveRight()
      case "pause":
        eng.pauseGame()
        return
      }
    }
  }
}

