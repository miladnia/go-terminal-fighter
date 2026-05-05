# Terminal Fighter

A nostalgic remake of the classic Brick Game racing experience.

![A screenshot of the game Terminal Fighter](./docs/screenshot.png)

## Installation

### Binary Release (Linux/Windows)

You can manually download a binary release from [the release page](https://github.com/miladnia/go-terminal-fighter/releases).

### Go

> [!NOTE]
> You'll need to [install Go](https://golang.org/doc/install)

```sh
go install github.com/miladnia/go-terminal-fighter@latest
```

## Build and Run

Required [Go](https://golang.org/doc/install) version >= 1.22

```sh
git clone https://github.com/miladnia/go-terminal-fighter.git
cd go-terminal-fighter
go build
./go-terminal-fighter
```

## Customization

Maps and levels are formatted in JSON files:

```
- maps/
- levels.json
```

## About The Brick Game

The Brick Game, originated in China and Russia in the early 1990s, includes games using a 10 × 20 block grid as a crude, low resolution dot matrix screen. [Wikipedia](https://en.wikipedia.org/wiki/Handheld_electronic_game)

![A photo of the handheld game console Brick Game](./docs/brick-game-01.png)
