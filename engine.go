package main

import "time"

type playable interface {
  step() (gameOver bool)
  render() string
  moveLeft()
  moveRight()
  getStatus() gameStatus
}

type gameStatus struct {
  gameOver bool
  won      bool
}

type runner interface {
  start(game playable)
  pause()
  resume()
  stop()
}

type engine struct {
  game       playable
  stepping   runner
  rendering  runner
  controller runner
  gamePaused chan struct{}
  gameOver   chan gameStatus
}

func newEngine(w writer, klgr *keyLogger, controllerSetup map[byte]string) *engine {
  e := &engine{
    gamePaused: make(chan struct{}),
    gameOver:   make(chan gameStatus),
  }
  e.rendering = newRenderRunner(42 * time.Millisecond, w)
  e.stepping = newStepRunner(500 * time.Millisecond, func() {
    e.stopGame()
    e.gameOver <- e.game.getStatus()
  })
  e.controller = newControllerRunner(klgr, controllerSetup, func() {
    e.pauseGame()
    e.gamePaused <- struct{}{}
  })
  return e
}

func (e *engine) startGame(game playable, speed int) {
  e.game = game
  go e.rendering.start(game)
  go e.stepping.start(game)
  go e.controller.start(game)
}

func (e *engine) pauseGame() {
  e.controller.pause()
  e.stepping.pause()
  // Let the rendering to be done completely.
  time.Sleep(100 * time.Millisecond)
  e.rendering.pause()
}

func (e *engine) resumeGame() {
  e.rendering.resume()
  e.controller.resume()
  e.stepping.resume()
}

func (e *engine) stopGame() {
  e.stepping.stop()
  e.controller.stop()
  // Let the rendering to be done completely.
  time.Sleep(100 * time.Millisecond)
  e.rendering.stop()
}

type stepRunner struct {
  interval   time.Duration
  ticker     *time.Ticker
  done       chan struct{}
  onGameOver func()
}

func newStepRunner(interval time.Duration, onGameOver func()) *stepRunner {
  return &stepRunner{
    interval:   interval,
    done:       make(chan struct{}),
    onGameOver: onGameOver,
  }
}

func (r *stepRunner) start(game playable) {
  r.ticker = time.NewTicker(r.interval)
  for {
    select {
    case <-r.done:
      r.ticker.Stop()
      return
    case <-r.ticker.C:
      gameOver := game.step()
      if gameOver {
        go r.onGameOver()
      }
    }
  }
}

func (r *stepRunner) pause() {
  r.ticker.Stop()
}

func (r *stepRunner) resume() {
  r.ticker.Reset(r.interval)
}

func (r *stepRunner) stop() {
  r.done <- struct{}{}
}

type renderRunner struct {
  interval time.Duration
  ticker   *time.Ticker
  done     chan struct{}
  w        writer
}

func newRenderRunner(interval time.Duration, w writer) *renderRunner {
  return &renderRunner {
    interval: interval,
    done:     make(chan struct{}),
    w:        w,
  }
}

func (r *renderRunner) start(game playable) {
  r.ticker = time.NewTicker(r.interval)
  for {
    select {
    case <-r.done:
      r.ticker.Stop()
      return
    case <-r.ticker.C:
      frame := game.render()
      r.w.write(frame)
    }
  }
}

func (r *renderRunner) pause() {
  r.ticker.Stop()
}

func (r *renderRunner) resume() {
  r.ticker.Reset(r.interval)
}

func (r *renderRunner) stop() {
  r.done <- struct{}{}
}

type controllerRunner struct {
  klgr     *keyLogger
  setup    map[byte]string
  onPause  func()
  chPause  chan struct{}
  chResume chan struct{}
  chStop   chan struct{}
}

func newControllerRunner(klgr *keyLogger, setup map[byte]string, onPause func()) *controllerRunner {
  return &controllerRunner{
    klgr:     klgr,
    setup:    setup,
    onPause:  onPause,
    chPause:  make(chan struct{}),
    chResume: make(chan struct{}),
    chStop:   make(chan struct{}),
  }
}

func (r *controllerRunner) start(game playable) {
  for {
    select {
    case <-r.chStop:
      return
    case <-r.chPause:
      select {
      case <-r.chResume:
        continue
      case <-r.chStop:
        return
      }
    case key := <-r.klgr.C:
      command := r.setup[key]
      switch command {
      case "left":
        game.moveLeft()
      case "right":
        game.moveRight()
      case "pause":
        go r.onPause()
      }
    }
  }
}

func (r *controllerRunner) pause() {
  r.chPause <- struct{}{}
}

func (r *controllerRunner) resume() {
  r.chResume <- struct{}{}
}

func (r *controllerRunner) stop() {
  r.chStop <- struct{}{}
}

