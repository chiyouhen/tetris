package main

import (
	"fmt"
	"math/rand/v2"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/term"
)

var background = [][]byte{
	[]byte(`XXXXXXXXXXXX............................`),
	[]byte(`X..........X............................`),
	[]byte(`X..........X............................`),
	[]byte(`X..........X............................`),
	[]byte(`X..........X............................`),
	[]byte(`X..........X............................`),
	[]byte(`X..........X............................`),
	[]byte(`X..........X............................`),
	[]byte(`X..........X............................`),
	[]byte(`X..........X............................`),
	[]byte(`X..........X............................`),
	[]byte(`X..........X............................`),
	[]byte(`X..........X............................`),
	[]byte(`X..........X............................`),
	[]byte(`X..........X............................`),
	[]byte(`X..........X............................`),
	[]byte(`X..........X............................`),
	[]byte(`X..........X............................`),
	[]byte(`X..........X............................`),
	[]byte(`X..........X............................`),
	[]byte(`X..........X............................`),
	[]byte(`XXXXXXXXXXXX............................`),
}

var field = []Bricks{
	newBricks(`..........`, nil),
	newBricks(`..........`, nil),
	newBricks(`..........`, nil),
	newBricks(`..........`, nil),
	newBricks(`..........`, nil),
	newBricks(`..........`, nil),
	newBricks(`..........`, nil),
	newBricks(`..........`, nil),
	newBricks(`..........`, nil),
	newBricks(`..........`, nil),
	newBricks(`..........`, nil),
	newBricks(`..........`, nil),
	newBricks(`..........`, nil),
	newBricks(`..........`, nil),
	newBricks(`..........`, nil),
	newBricks(`..........`, nil),
	newBricks(`..........`, nil),
	newBricks(`..........`, nil),
	newBricks(`..........`, nil),
	newBricks(`..........`, nil),
}

var tetrominoes = [][][]byte{
	{
		[]byte(`....`),
		[]byte(`.XX.`),
		[]byte(`.XX.`),
		[]byte(`....`),
	},
	{
		[]byte(`.X..`),
		[]byte(`.X..`),
		[]byte(`.X..`),
		[]byte(`.X..`),
	},
	{
		[]byte(`....`),
		[]byte(`XX..`),
		[]byte(`.XX.`),
		[]byte(`....`),
	},
	{
		[]byte(`....`),
		[]byte(`.XX.`),
		[]byte(`XX..`),
		[]byte(`....`),
	},
	{
		[]byte(`....`),
		[]byte(`.X..`),
		[]byte(`XXX.`),
		[]byte(`....`),
	},
	{
		[]byte(`.X..`),
		[]byte(`.X..`),
		[]byte(`XX..`),
		[]byte(`....`),
	},
	{
		[]byte(`.X..`),
		[]byte(`.X..`),
		[]byte(`.XX.`),
		[]byte(`....`),
	},
}

var tetromino = []Bricks{
	newBricks(`....`, nil),
	newBricks(`....`, nil),
	newBricks(`....`, nil),
	newBricks(`....`, nil),
}
var nextTetromino = []Bricks{
	newBricks(`....`, nil),
	newBricks(`....`, nil),
	newBricks(`....`, nil),
	newBricks(`....`, nil),
}

const (
	width             int = 10
	height            int = 20
	frameBufferWidth      = 22
	frameBufferHeight     = 40
)

type Brick struct {
	Color []byte
	C     byte
}

func emptyBrick() Brick {
	return Brick{
		Color: terminal.Escape.Reset,
		C:     ' ',
	}
}
func transBrick() Brick {
	return Brick{
		Color: terminal.Escape.Reset,
		C:     '.',
	}
}

type Bricks []Brick

func newBricks(s string, color []byte) Bricks {
	bricks := make(Bricks, len(s))
	for i, c := range []byte(s) {
		bricks[i] = Brick{
			Color: color,
			C:     c,
		}
	}
	return bricks
}

var frameBuffer [22][40]Brick

func clearFrameBuffer() {
	for y, line := range frameBuffer {
		for x := range line {
			frameBuffer[y][x] = emptyBrick()
		}
	}
}

func drawBackground() {
	for y, line := range background {
		for x, col := range line {
			if col == '.' {
				continue
			}
			frameBuffer[y][x] = Brick{
				Color: terminal.Escape.Reset,
				C:     col,
			}
		}
	}
}

func drawField() {
	for y, line := range field {
		for x, c := range line {
			if c.C == '.' {
				continue
			}
			frameBuffer[y+1][x+1] = c
		}
	}
}

func placeTetromino(tidx int, x, y int, rotation int) {
	updateTetromino(tetromino, tidx, rotation)
	drawTetromino(tetromino, x+1, y+1)
}

func placeNextTetromino(tidx int, x, y int, rotation int) {
	updateTetromino(nextTetromino, tidx, rotation)
	drawTetromino(nextTetromino, x, y)
}

func tetrominoColor(tidx int) []byte {
	switch tidx {
	case 0:
		// O
		return terminal.Escape.Cyan
	case 1:
		// I
		return terminal.Escape.Red
	case 2:
		// Z
		return terminal.Escape.Blue
	case 3:
		// Z
		return terminal.Escape.Blue
	case 4:
		// T
		return terminal.Escape.Yellow
	case 5:
		// L
		return terminal.Escape.Green
	case 6:
		// L
		return terminal.Escape.Green
	}
	return terminal.Escape.Reset
}

func drawTetromino(block []Bricks, x, y int) {
	for yy := 0; yy < 4; yy++ {
		for xx := 0; xx < 4; xx++ {
			c := block[yy][xx]
			if c.C == '.' {
				continue
			}
			frameBuffer[yy+y][xx+x] = c
		}
	}
}

func canPlace(tidx int, x, y int, rotation int) bool {
	updateTetromino(tetromino, tidx, rotation)
	for yy := 0; yy < 4; yy++ {
		for xx := 0; xx < 4; xx++ {
			c := tetromino[yy][xx]
			if c.C == '.' {
				continue
			}
			if yy+y < 0 || yy+y >= height || xx+x < 0 || xx+x >= width || field[yy+y][xx+x].C != '.' {
				return false
			}
		}
	}
	return true
}

func tetrominoBrick(tidx int, xx, yy int, rotation int) byte {
	tetromino := tetrominoes[tidx]
	switch rotation {
	case 0:
		return tetromino[yy][xx]
	case 90:
		return tetromino[3-xx][yy]
	case 180:
		return tetromino[3-yy][3-xx]
	case 270:
		return tetromino[xx][3-yy]
	}
	return 0
}

func updateTetromino(block []Bricks, tidx int, rotation int) {
	for yy := 0; yy < 4; yy++ {
		for xx := 0; xx < 4; xx++ {
			block[yy][xx] = Brick{
				Color: tetrominoColor(tidx),
				C:     tetrominoBrick(tidx, xx, yy, rotation),
			}
		}
	}
}

func fixTetromino(tidx int, x, y int, rotation int) {
	updateTetromino(tetromino, tidx, rotation)
	for yy := 0; yy < 4; yy++ {
		for xx := 0; xx < 4; xx++ {
			c := tetromino[yy][xx]
			if c.C == '.' {
				continue
			}
			field[yy+y][xx+x] = c
		}
	}
}

func score() int {
	for y := height - 1; y >= 0; y-- {
		count := 0
		for x := 0; x < width; x++ {
			if field[y][x].C != '.' {
				count += 1
			}
		}
		if count == width {
			return y
		}
	}
	return -1
}

func drawScore(y int, frame int) {
	switch frame {
	case 0:
		for x := 0; x < width; x++ {
			field[y][x] = Brick{
				Color: terminal.Escape.Yellow,
				C:     '-',
			}
		}
	case 1:
		for x := 0; x < width; x++ {
			field[y][x] = transBrick()
		}
	case 2:
		for {
			if y == 0 {
				for x := 0; x < width; x++ {
					field[y][x] = transBrick()
				}
				return
			}
			for x := 0; x < width; x++ {
				field[y][x] = field[y-1][x]
			}
			y -= 1
		}
	}
}

func render() {
	posCursor(0, 0)
	for _, line := range frameBuffer {
		for _, col := range line {
			terminal.Write(append(col.Color, col.C))
		}
		terminal.Write([]byte("\r\n"))

	}
}

func clear() {
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			field[y][x] = Brick{
				Color: terminal.Escape.Reset,
				C:     '.',
			}
		}
	}
}

func posCursor(x, y int) {
	var bs strings.Builder
	bs.WriteByte(033)
	bs.WriteByte('[')
	bs.WriteString(strconv.Itoa(y))
	bs.WriteByte(';')
	bs.WriteString(strconv.Itoa(x))
	bs.WriteByte('H')
	terminal.Write([]byte(bs.String()))
}

func clearScreen() {
	terminal.Write([]byte{033, 'c'})
}
func hideCursor() {
	terminal.Write(append([]byte{033, '['}, "?25l"...))
}
func showCursor() {
	terminal.Write(append([]byte{033, '['}, "?25h"...))
}

func tprintf(x, y int, msg string, args ...any) {
	s := fmt.Sprintf(msg, args...)
	a := make([]Brick, 0, len([]byte(s)))

	for xx, c := range []byte(s) {
		b := Brick{
			Color: nil,
			C:     c,
		}
		if xx == 0 {
			b.Color = terminal.Escape.Reset
		}
		a = append(a, b)
	}
	copy(frameBuffer[y][x:], a)
}

var terminal *term.Terminal

const maxAutoDown = 20

var tidx, nextTidx int
var x, y, rotation, nextRotation int

func spotTetromino() {
	tidx = nextTidx
	nextTidx = rand.IntN(7)
	rotation = nextRotation
	nextRotation = rand.IntN(4) * 90
	x = 3
	y = 0
}

func main() {
	debug := false
	if len(os.Args) == 2 {
		if os.Args[1] == "--debug" {
			debug = true
		}

	}
	oldState, err := term.MakeRaw(int(os.Stderr.Fd()))
	if err != nil {
		panic(err)
	}
	defer term.Restore(int(os.Stdout.Fd()), oldState)
	terminal = term.NewTerminal(os.Stdout, "")
	terminal.SetSize(frameBufferWidth, frameBufferHeight)
	clearScreen()
	defer clearScreen()
	hideCursor()
	defer showCursor()
	tkr := time.NewTicker(time.Millisecond * 50)
	input := make(chan byte, 1)
	go func() {
		var buf [1]byte
		for {
			_, err := os.Stderr.Read(buf[:])
			if err == nil {
				input <- buf[0]
			}
		}
	}()
	spotTetromino()
	tidx = rand.IntN(7)
	rotation = rand.IntN(4) * 90
	updateTetromino(tetromino, tidx, rotation)
	updateTetromino(nextTetromino, nextTidx, nextRotation)

	fmt.Printf("tidx: %d, rotation: %d\n", tidx, rotation)
	autoDown := maxAutoDown
	scoreline := -1
	scoreframe := -1
	totalScore := 0
	combo := 0
	var gameover bool

	var c byte
	for ; true; <-tkr.C {
		select {
		case c = <-input:
			if !gameover && tidx > -1 {
				switch c {
				case 'h':
					if canPlace(tidx, x-1, y, rotation) {
						x -= 1
					}
				case 'l':
					if canPlace(tidx, x+1, y, rotation) {
						x += 1
					}
				case 'k':
					if debug {
						if canPlace(tidx, x, y-1, rotation) {
							y -= 1
						}
					}
				case 'j':
					if canPlace(tidx, x, y+1, rotation) {
						y += 1
						autoDown = maxAutoDown
					}
				case 'f':
					if canPlace(tidx, x, y, (rotation+90)%360) {
						rotation = (rotation + 90) % 360
					}
				case 'p':
					if debug {
						fixTetromino(tidx, x, y, rotation)
						spotTetromino()
					}
				}
			}
			if gameover {
				if c == 'n' {
					clear()
					totalScore = 0
					spotTetromino()
					gameover = false

				}
			}
			if c == 'q' {
				return
			}
		default:
		}

		if !gameover {
			autoDown -= 1
			if autoDown == 0 {
				if canPlace(tidx, x, y+1, rotation) {
					y += 1
				} else {
					fixTetromino(tidx, x, y, rotation)
					tidx = -1
				}
				autoDown = maxAutoDown
			}
		}

		clearFrameBuffer()
		drawBackground()
		drawField()

		if tidx > -1 {
			placeTetromino(tidx, x, y, rotation)
		} else {
			if scoreframe > -1 {
				drawScore(scoreline, scoreframe)
				scoreframe += 1
				if scoreframe == 3 {
					scoreframe = -1
				}
			} else {
				scoreline = score()
				if scoreline == -1 {
					combo = 0
					if !canPlace(nextTidx, 3, 0, nextRotation) {
						gameover = true
					} else {
						spotTetromino()
					}
				} else {
					combo += 1
					scoreframe = 0
					totalScore += combo
				}
			}
		}
		tprintf(13, 0, "Next:")
		placeNextTetromino(nextTidx, 13, 1, nextRotation)
		if gameover {
			tprintf(13, 6, "GAME OVER!")
		}

		tprintf(13, 5, "Score: %d", totalScore)

		if debug {
			tprintf(13, 6, "tidx: %d, rotation: %d", tidx, rotation)
			tprintf(13, 7, "auto: %d", autoDown)
			tprintf(13, 8, "scoreline: %d", scoreline)
			tprintf(13, 9, "scoreframe: %d", scoreframe)
			tprintf(13, 10, "%c", '\U00002014')
		}
		tprintf(13, 19, "[H]left [L]right [J]down")
		tprintf(13, 20, "[F]rotate [Q]quit")

		render()
		posCursor(x+2, y+2)

	}
}
