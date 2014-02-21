package main

import (
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
}

/*
A player's paddle is defined by the top and height.
Top is where the paddle is located in the field,
and height is the height of the paddle.
A paddle cannot break the top and bottom walls of the field.
c is the paddle's column
*/
type Player struct {
    t, h, c, score int
}

type PlayerOutOfBounds Player

func (p PlayerOutOfBounds) Error() string {
    return fmt.Sprintln("Player out of bounds")
}

/*
Initialize the player based on the height/width of the field
*/
func InitializePlayer(h, w, number int) (p *Player) {
    var c int
    // Player 1 is at the left
    if number == 1 {
        c = 1
    } else {
        // Player 2 is at the right
        c = w -2
    }
    p = &Player{int(h/3), int(h/4), c, 0}
    return
}

/*
d is delta that the paddle is to move.
Player's paddle cannot move beyond the walls of the field.
*/
func (p *Player) Move(d int, f *Field) error {
    // If going up
    if d < 0 {
        if d + p.t > 0 {
            p.t += d
            return nil
        } else {
            p.t = 0
            return nil
        }
        //return an out of bounds error
        return PlayerOutOfBounds(*p)
    } else {
        if p.t + p.h - 1 + d < f.h {
            p.t += d
            return nil
        } else {
            p.t = f.h - 1
            return nil
        }
        // return an out of bounds error
        return PlayerOutOfBounds(*p)
    }
}

/*
determines whether or not the ball is deflected
*/
func (p *Player) IsThisDeflection(b *Ball) bool {
    if b.y >= p.t && b.y <= p.t + p.h - 1{
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
    y, x, my, mx int
}

func (b *Ball) HitWall() {
    b.my = -1 * b.my
}

func (b *Ball) HitPaddle() {
    b.mx = -1 * b.mx
}

func (b *Ball) Move() {
    b.x = b.x + b.mx
    b.y = b.y + b.my
}
func (b *Ball) SetAtCenter(f *Field, p int) {
    var d int
    if p == 1 {
        d = -1
    } else {
        d = 1
    }
    b.y, b.x, b.my, b.mx = int(f.h/2), int(f.w/2), d, d
}

/*
Function that controls the ball moving in a certain time interval
and collisions.
*/
func MoveBall(b *Ball, f *Field, p1, p2 *Player, score chan bool) {
    for {
        time.Sleep(time.Second / 5)
        b.Move()
        // If this is a score, increment the other player's score.
        if p1.IsThisScore(b) {
            p2.score++
            b.SetAtCenter(f, 2)
            score <- true
            return
        } else if p2.IsThisScore(b) {
            p1.score++
            b.SetAtCenter(f, 1)
            score <- true
            return
        }
        // if after this move the ball is at a paddle,
        // its next move should be to deflect off in the opposite direction.
        if p1.IsThisDeflection(b) || p2.IsThisDeflection(b) {
            b.HitPaddle()
        }
        // if after this move the ball is at the wall,
        // its next move should be to deflect off in the opposite direction.
        if b.y == f.h - 2 || b.y == 1 {
            b.HitWall()
        }
    }
}

/*
This is the int to char representation
of things on the field.
*/
var GameGraphics = map[string]string{"wall": "#", "paddle": "+", "ball": "O"}
/*
Draw the game to the terminal.
Results in annoying blinky experience.
Handles clearing the screen, as well.
*/
func DrawGame(f *Field, p1, p2 *Player, b *Ball, showScore bool) {
    // Clear
    fmt.Println("\033[2J\033[0;0H")
    // Field drawing
    for i := 0; i < f.h; i++ {
        fmt.Print(GameGraphics["wall"]);
        for ii := 0; ii < f.w; ii++ {
            // If this is the top or bottom, draw a wall
            if i == 0 || i == f.h - 1 {
                fmt.Print(GameGraphics["wall"])
            } else if i > 0 && i < f.h - 1 {
                // If this is the center of the screen,
                // display the score string. Else, draw the game.
                if showScore && i == int(f.h/2) {
                    if ii == int(f.w/3) {
                        fmt.Print("     Left Player: ", p1.score, " Right Player: ", p2.score, "     ");
                    }
                } else {
                    // If this is the second column
                    if ii == 1 {
                        // Could be where player 1's paddle is
                        if i >= p1.t && i <= p1.t + p1.h {
                            fmt.Print(GameGraphics["paddle"])
                        } else {
                            fmt.Print(" ");
                        }
                    } else if ii == f.w - 2 {
                        if i >= p2.t && i <= p2.t + p2.h {
                            fmt.Print(GameGraphics["paddle"])
                        } else {
                            fmt.Print(" ");
                        }
                    } else {
                        if b.x == ii && b.y == i {
                            fmt.Print(GameGraphics["ball"])
                        } else {
                            fmt.Print(" ");
                        }
                    }
                }
            }
        }
        fmt.Println(GameGraphics["wall"]);
    }
}

func main() {
    // Initial height and width for this field
    h, w := 20, 50
    // Create a Ball
    b := &Ball{int(h/2), int(w/2), -1, -1}
    // Create 2 players
    p1, p2 := InitializePlayer(h, w, 1), InitializePlayer(h, w, 2)
    // Create a field object
    f := &Field{h, w, p1, p2, b}
    //done := false
    score := make(chan bool)
    go MoveBall(b, f, p1, p2, score)
    //go GetPlayerMoves(p1, 'w', 's')
    //go GetPlayerMoves(p2, 'o', 'l')
    for i := 0; i < 100; i++ {
        select {
            case <-score:
                DrawGame(f, p1, p2, b, true)
                time.Sleep(time.Second * 2)
                go MoveBall(b, f, p1, p2, score)
            default:
                DrawGame(f, p1, p2, b, false)
        }
        time.Sleep(time.Second / 10)
    }
}
