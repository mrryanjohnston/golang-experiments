# Pong

Classic game implemented in go.

## Installation

You'll need `pkg-config` and `ncurses` first. Read
[this](https://github.com/rthornton128/goncurses/wiki)
to learn how to do it.

Then, do:

```
go get github.com/rthornton128/goncurses
```

Then:

```
go build pong.go
```

Finally:

```
./pong
```

## Controls

- `w` and `s` control the left paddle.
- `o` and `l` controll the right paddle.
- `q` quits the game.

## Win condition
Game is over when one player gets 7 points.

## Libraries
[goncurses](https://code.google.com/p/goncurses/)
