package main

import "fmt"

const levelsFile = "levels.json"

type level struct {
  No      int
  MapFile string `json:"map-file"`
  Speed   int    `json:"speed"`
  IsFinal bool
}

func getLevel(levelNo int) (level, error) {
  var levels []level
  err := decodeJsonFile(levelsFile, &levels)
  if err != nil {
    return level{}, err
  }
  i := levelNo - 1
  // Out of range?
  if i < 0 || i >= len(levels) {
    return level{}, fmt.Errorf("level '%d' is out of range.", levelNo)
  }
  lvl := levels[i]
  lvl.No = levelNo
  if i == len(levels) - 1 {
    lvl.IsFinal = true
  }
  return lvl, nil
}

func (lvl *level) getMap() (m [][]int) {
  err := decodeJsonFile(lvl.MapFile, &m)
  if err != nil {
    panic(err)
  }
  return m
}

