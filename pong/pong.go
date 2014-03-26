package main

import (
        "code.google.com/p/goncurses"
        "fmt"
        "time"
)

/*
 Field describes the current state of things,
 such as players and the location of the ball.
 */
type Field struct {
    h, w int
    p1, p2 *Player
    b *Ball
    moveballchan chan bool
    scorechan chan bool
    done chan bool
}

/*
Function that controls the ball moving in a certain time interval
and collisions.
*/
func (f *Field) MoveBall() {
    for {
        f.b.Move(f.moveballchan)
        // If this is a score, increment the other player's score.
        if f.p1.IsThisScore(f.b) {
            f.p2.score++
                f.b.SetAtCenter(f, 2)
                f.scorechan <- true
                if f.p2.score == 7 {
                    f.done <- true
                }
                return
        } else if f.p2.IsThisScore(f.b) {
            f.p1.score++
                f.b.SetAtCenter(f, 1)
                f.scorechan <- true
                if f.p1.score == 7 {
                    f.done <- true
                }
                return
        }
        // if after this move the ball is at a paddle,
        // its next move should be to deflect off in the opposite direction.
        if f.p1.IsThisDeflection(f.b) || f.p2.IsThisDeflection(f.b) {
            f.b.HitPaddle()
        }
        // if after this move the ball is at the wall,
        // its next move should be to deflect off in the opposite direction.
        if f.b.y == f.h - 2 || f.b.y == 1 {
            f.b.HitWall()
        }
        // Finally, sleep it a certain amount of time.
        time.Sleep(time.Second / time.Duration(f.b.speed))
    }
}


/*
A player's paddle is defined by the top and height.
Bottom is provided as a convenience.
Top is where the paddle is located in the field,
and height is the height of the paddle.
A paddle cannot break the top and bottom walls of the field.
c is the paddle's column
*/
type Player struct {
    t, h, b, c, score int
    movepaddlechan chan int
}

/*
Initialize the player based on the height/width of the field
*/
func InitializePlayer(h, w, number int, movepaddlechan chan int) (p *Player) {
    var c, b, t int
    // Player 1 is at the left
    if number == 1 {
        c = 2
    } else {
        // Player 2 is at the right
        c = w -2
    }
    t = int(h/3)
    h = int(h/4)
    // Set bottom to be the top plus the height minus 1 for 0 indexed
    b = t + h - 1
    p = &Player{t, h, b, c, 0, movepaddlechan}
    return
}

/*
d is delta that the paddle is to move.
Player's paddle cannot move beyond the walls of the field.
*/
func (p *Player) Move(d int, f *Field) bool {
    // If going up
    if d < 0 {
        if d + p.t > 0 {
            p.t += d
            p.b += d
            return true
        }
    } else {
        if p.b + d < f.h - 1 {
            p.t += d
            p.b += d
            return true
        }
    }
    return false
}

/*
determines whether or not the ball is deflected
*/
func (p *Player) IsThisDeflection(b *Ball) bool {
    if b.y >= p.t && b.y <= p.b {
        if b.x == p.c + 1 || b.x == p.c - 1 {
            return true
        }
    }
    return false
}
/*
determines whether or not the ball is a score
*/
func (p *Player) IsThisScore(b *Ball) bool {
    if b.x == p.c {
        return true
    }
    return false
}

/*
func GetPlayerMoves(p *Player, up, down rune) {
}
*/

/*
A ball moves one space by mx and my
depending on where it hits the player's paddle (TODO)
*/
type Ball struct {
    y, x, oldy, oldx, my, mx, speed, hits int
}

func (b *Ball) HitWall() {
    b.my = -1 * b.my
}

/*
Speed increases on paddle hit
*/
func (b *Ball) HitPaddle() {
    b.mx = -1 * b.mx
    b.hits++
    if (b.hits * 2 > b.speed) {
        b.speed += 2;
    }
}

func (b *Ball) Move(c chan bool) {
    b.oldx = b.x
    b.oldy = b.y
    b.x = b.x + b.mx
    b.y = b.y + b.my
    c <- true
}
func (b *Ball) SetAtCenter(f *Field, p int) {
    var d int
    // Set direction of the ball
    if p == 1 {
        d = -1
    } else {
        d = 1
    }
    b.y, b.x, b.my, b.mx = int(f.h/2), int(f.w/2), d, d
}

/*
This is the int to char representation
of things on the field.
*/
var GameGraphics = map[string]byte{"wall": '#', "paddle": '+', "ball": 'O'}
/*
Initial Drawing of the game field + paddle + ball
*/
func InitialDrawGame(f *Field, stdscr *goncurses.Window) {
    // Field drawing
    for i := 0; i < f.h; i++ {
        stdscr.MovePrintf(i, 0, "%c", GameGraphics["wall"])
        for ii := 0; ii < f.w; ii++ {
            // If this is the top or bottom, draw a wall
            if i == 0 || i == f.h - 1 {
                stdscr.MovePrintf(i, ii, "%c", GameGraphics["wall"])
            } else if i > 0 && i < f.h - 1 {
                // If this is the center of the screen,
                // display the score string. Else, draw the game.
                // If this is the second column
                if ii == f.p1.c {
                    // Could be where player 1's paddle is
                    if i >= f.p1.t && i <= f.p1.b {
                        stdscr.MovePrintf(i, ii, "%c", GameGraphics["paddle"])
                    }
                } else if ii == f.p2.c {
                    if i >= f.p2.t && i <= f.p2.b {
                        stdscr.MovePrintf(i, ii, "%c", GameGraphics["paddle"])
                    }
                } else {
                    if f.b.x == ii && f.b.y == i {
                        stdscr.MovePrintf(i, ii, "%c", GameGraphics["ball"])
                    }
                }
            }
        }
        stdscr.MovePrintf(i, f.w, "%c", GameGraphics["wall"])
    }
    stdscr.Print("\n")
    stdscr.Println("Player 1: 0")
    stdscr.Println("Player 2: 0")
    stdscr.Refresh()
}

/*
Gets user input the goncurses way
*/
func TakeUserInput(f *Field, stdscr *goncurses.Window) {
    for {
        ch := stdscr.GetChar()
        switch byte(ch) {
             case 'w':
                move := f.p1.Move(-1, f)
                if move {
                    f.p1.movepaddlechan <- -1
                }
             case 's':
                move := f.p1.Move(1, f)
                if move {
                    f.p1.movepaddlechan <- 1
                }
             case 'o':
                move := f.p2.Move(-1, f)
                if move {
                    f.p2.movepaddlechan <- -1
                }
             case 'l':
                move := f.p2.Move(1, f)
                if move {
                    f.p2.movepaddlechan <- 1
                }
             case 'q':
                 f.done <- true
         }
    }
}

/*
Runs as a go routine waiting for signals from other functions.
Passes off the the appropriate draw function when data received on a chan.
*/
func DrawAction(f *Field, stdscr *goncurses.Window) {
    var delta int
    for {
        select {
            case delta = <-f.p1.movepaddlechan:
                DrawPaddleMove(stdscr, f.p1, delta)
            case delta = <-f.p2.movepaddlechan:
                DrawPaddleMove(stdscr, f.p2, delta)
            case <-f.moveballchan:
                DrawBallMove(stdscr, f.b)
            case <- f.scorechan:
                DrawScores(stdscr, f)
                go f.MoveBall()
        }
        stdscr.Refresh()
    }
}

func DrawBallMove(stdscr *goncurses.Window, b *Ball) {
    // Unset the old one
    stdscr.MovePrint(b.oldy, b.oldx, " ")
    // Set the new one
    stdscr.MovePrintf(b.y, b.x, "%c", GameGraphics["ball"])
    return
}

/*
d is the delta y
*/
func DrawPaddleMove(stdscr *goncurses.Window, p *Player, d int) {
    // if it is going up, remove the bottom most symbol
    // and add one to the top
    if d < 0 {
        stdscr.MovePrint(p.b + 1, p.c, " ")
        stdscr.MovePrintf(p.t, p.c, "%c", GameGraphics["paddle"])
    } else {
        stdscr.MovePrint(p.t-1, p.c, " ")
        stdscr.MovePrintf(p.b, p.c, "%c", GameGraphics["paddle"])
    }
    return
}

func DrawScores(stdscr *goncurses.Window, f *Field) {
    stdscr.MovePrint(f.h, 10, f.p1.score)
    stdscr.MovePrint(f.h+1, 10, f.p2.score)
}

func main() {
    stdscr, err := goncurses.Init()
    if err != nil {
        fmt.Println("goncurses failed to intialize: ", err)
        return
    }
    defer goncurses.End()
    goncurses.Echo(false)
    goncurses.Cursor(0)
    // Initial height and width for this field
    h, w := 20, 50
    // Create a Ball
    b := &Ball{int(h/2), int(w/2), int(h/2), int(w/2), -1, -1, 10, 0}
    // Create 2 players
    movepaddle1chan := make(chan int)
    movepaddle2chan := make(chan int)
    p1, p2 := InitializePlayer(h, w, 1, movepaddle1chan), InitializePlayer(h, w, 2, movepaddle2chan)
    // Create a field object
    moveballchan   := make(chan bool)
    scorechan      := make(chan bool)
    done           := make(chan bool)
    f := &Field{h, w, p1, p2, b, moveballchan, scorechan, done}
    InitialDrawGame(f, stdscr)
    go f.MoveBall()
    go TakeUserInput(f, stdscr)
    go DrawAction(f, stdscr)
    <-done
    stdscr.Clear()
    stdscr.Println("Game over!")
    time.Sleep(time.Second * 2)
    return
}
