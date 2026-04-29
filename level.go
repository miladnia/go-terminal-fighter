package main

import (
  "embed"
  "encoding/json"
  "fmt"
  "path/filepath"
)

const (
  levelsFile = "levels.json"
  mapsDir    = "maps"
)

//go:embed levels.json maps/*.json
var gameFS embed.FS

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
  mapFile := filepath.Join(mapsDir, lvl.MapFile)
  err := decodeJsonFile(mapFile, &m)
  if err != nil {
    panic(err)
  }
  return m
}

func decodeJsonFile(filename string, v any) error {
  data, err := gameFS.ReadFile(filename)
	if err != nil {
    return fmt.Errorf("could not open the file '%v'! (%s)", filename, err)
	}

  if err := json.Unmarshal(data, v); err != nil {
    return fmt.Errorf("could not decode the file '%v'! (%s)", filename, err)
  }

  return nil
}

