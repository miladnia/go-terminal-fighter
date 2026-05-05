package main

import (
  "embed"
  "encoding/json"
  "fmt"
  "path"
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
  err := decodeJsonFile(&levels, levelsFile)
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
  err := decodeJsonFile(&m, mapsDir, lvl.MapFile)
  if err != nil {
    panic(err)
  }
  return m
}

// decodeJsonFile unmarshals JSON data from the embedded filesystem into dest.
// The pathSegments are joined as path components (directories followed by filename),
// not as multiple separate file paths. Example: decodeJsonFile(&data, "maps", "blue.json").
func decodeJsonFile(dest any, pathSegments ...string) error {
  fullPath := path.Join(pathSegments...)
  rawJSON, err := gameFS.ReadFile(fullPath)
	if err != nil {
    return fmt.Errorf("could not open the file '%v'! (%s)", fullPath, err)
	}

  if err := json.Unmarshal(rawJSON, dest); err != nil {
    return fmt.Errorf("could not decode the file '%v'! (%s)", fullPath, err)
  }

  return nil
}

