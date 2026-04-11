package main

import "time"

type game interface {
  render() string
  step() (end bool, won bool)
  moveLeft()
  moveRight()
}

type gameStatus struct {
  won bool
}

type engine struct {
  w                writer
  klgr             *keyLogger
  controllerSetup  map[byte]string
  g                game
  movementInterval time.Duration
  renderInterval   time.Duration
  renderTicker     *time.Ticker
  movementTicker   *time.Ticker
  gamePaused       chan struct{}
  gameOver         chan gameStatus
  renderingDone    chan struct{}
  controllerDone   chan struct{}
}

func newEngine(w writer, klgr *keyLogger, controllerSetup map[byte]string) *engine {
  return &engine{
    w:               w,
    klgr:            klgr,
    controllerSetup: controllerSetup,
    gamePaused:      make(chan struct{}),
    gameOver:        make(chan gameStatus),
    renderingDone:   make(chan struct{}),
    controllerDone:  make(chan struct{}),
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
  eng.resumeRendering()
  time.Sleep(2000 * time.Millisecond)
  go eng.listenToController()
  eng.resumeMovement()
}

func (eng *engine) stopGame(s gameStatus) {
  eng.controllerDone <- struct{}{}
  if !s.won {
    time.Sleep(2000 * time.Millisecond)
  }
  eng.renderingDone <- struct{}{}
  eng.gameOver <- s
}

func (eng *engine) startMovement() {
  eng.movementTicker = time.NewTicker(eng.movementInterval)
  for ; ; <-eng.movementTicker.C {
    end, won := eng.g.step()
    if end {
      eng.stopGame(gameStatus{
        won: won,
      })
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
      eng.w.write(frame)
    }
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

