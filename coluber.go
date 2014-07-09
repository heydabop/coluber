package main

import (
	"github.com/nsf/termbox-go"
	"log"
	"time"
)

const (
	ColorWall  = termbox.ColorRed
	ColorEmpty = termbox.ColorBlack
	ColorSnake = termbox.ColorGreen
	ColorFood  = termbox.ColorYellow
)

type Cell struct {
	X     int
	Y     int
	Clear bool
	Color termbox.Attribute
}

type Segment struct {
	X   int
	Y   int
	Dir int
}

func gameOver(snake []Segment, board [][]Cell) {
	for i := range snake {
		termbox.SetCell(snake[i].X, snake[i].Y, 0x0000, termbox.ColorBlack, termbox.ColorRed)
	}
	termbox.Flush()
	return
}

func moveSnake(snake []Segment, board [][]Cell, lastDir *int) {
	for {
		for i := range snake {
			if i == len(snake)-1 {
				board[snake[i].Y][snake[i].X].Color = ColorEmpty
				board[snake[i].Y][snake[i].X].Clear = true
				termbox.SetCell(snake[i].X, snake[i].Y, 0x0000, termbox.ColorBlack, ColorEmpty)
			}
			switch snake[i].Dir {
			case 0:
				snake[i].Y = snake[i].Y - 1
				break
			case 1:
				snake[i].X = snake[i].X + 1
				break
			case 2:
				snake[i].Y = snake[i].Y + 1
				break
			case 3:
				snake[i].X = snake[i].X - 1
				break
			}
			if i == 0 {
				if !board[snake[i].Y][snake[i].X].Clear { //collision
					switch snake[i].Dir { //undo move of head segment
					case 0:
						snake[i].Y = snake[i].Y + 1
						break
					case 1:
						snake[i].X = snake[i].X - 1
						break
					case 2:
						snake[i].Y = snake[i].Y - 1
						break
					case 3:
						snake[i].X = snake[i].X + 1
						break
					}
					gameOver(snake, board)
					return
				}
				board[snake[i].Y][snake[i].X].Color = ColorSnake
				board[snake[i].Y][snake[i].X].Clear = false
				termbox.SetCell(snake[i].X, snake[i].Y, 0x0000, termbox.ColorBlack, ColorSnake)
			}
		}
		for j := len(snake) - 1; j > 0; j-- {
			snake[j].Dir = snake[j-1].Dir
		}
		*lastDir = snake[0].Dir
		termbox.Flush()
		time.Sleep(200 * time.Millisecond)
	}
}

func main() {
	err := termbox.Init()
	if err != nil {
		log.Panicln(err)
	}
	defer termbox.Close()
	termbox.SetInputMode(termbox.InputEsc)

	const boardSize = 40
	board := make([][]Cell, boardSize)
	for y := range board {
		board[y] = make([]Cell, boardSize*2)
		for x := range board[y] {
			if x == 0 || y == 0 || x == (boardSize*2)-1 || y == boardSize-1 {
				board[y][x] = Cell{x, y, false, ColorWall}
			} else {
				board[y][x] = Cell{x, y, true, ColorEmpty}
			}
		}
	}
	snake := make([]Segment, 7, 16)
	for i := range snake {
		snake[i] = Segment{40 - i, 20, 1}
		board[20][40-i].Color = ColorSnake
		board[20][40-i].Clear = false
	}
	for _, row := range board {
		for _, cell := range row {
			termbox.SetCell(cell.X, cell.Y, 0x0000, termbox.ColorBlack, cell.Color)
		}
	}
	termbox.Flush()
	lastDir := snake[0].Dir
	go moveSnake(snake, board, &lastDir)
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyArrowUp:
				if lastDir != 2 {
					snake[0].Dir = 0
				}
				break
			case termbox.KeyArrowRight:
				if lastDir != 3 {
					snake[0].Dir = 1
				}
				break
			case termbox.KeyArrowDown:
				if lastDir != 0 {
					snake[0].Dir = 2
				}
				break
			case termbox.KeyArrowLeft:
				if lastDir != 1 {
					snake[0].Dir = 3
				}
				break
			case termbox.KeyCtrlC:
				return
			}
		case termbox.EventError:
			log.Panic(ev.Err)
			break
		}
	}
}
