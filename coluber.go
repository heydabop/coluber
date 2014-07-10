package main

import (
	"github.com/nsf/termbox-go"
	"log"
	"math/rand"
	"strconv"
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
		termbox.SetCell(snake[i].X, snake[i].Y, 0x0020, termbox.ColorBlack, termbox.ColorRed)
	}
	termbox.Flush()
	return
}

func renderScore(score uint64) {
	scoreStr := strconv.FormatUint(score, 10)
	for i := range scoreStr {
		termbox.SetCell(i, 0, rune(scoreStr[i]), termbox.ColorBlack, ColorWall)
	}
}

func initGame() ([]Segment, [][]Cell) {
	boardSizeY := 40
	boardSizeX := boardSizeY * 2
	termSizeX, termSizeY := termbox.Size()
	if termSizeY < boardSizeY {
		boardSizeY = termSizeY
	}
	if termSizeX < boardSizeX {
		boardSizeX = termSizeX
	}
	board := make([][]Cell, boardSizeY)
	for y := range board {
		board[y] = make([]Cell, boardSizeX)
		for x := range board[y] {
			if x == 0 || y == 0 || x == (boardSizeX)-1 || y == boardSizeY-1 {
				board[y][x] = Cell{x, y, false, ColorWall}
			} else {
				board[y][x] = Cell{x, y, true, ColorEmpty}
			}
		}
	}
	snake := make([]Segment, 7, 1024)
	for i := range snake {
		snake[i] = Segment{40 - i, 20, 1}
		board[20][40-i].Color = ColorSnake
		board[20][40-i].Clear = false
	}
	makeWalls(board)
	for _, row := range board {
		for _, cell := range row {
			termbox.SetCell(cell.X, cell.Y, 0x0020, termbox.ColorBlack, cell.Color)
		}
	}
	termbox.Flush()
	return snake, board
}

func makeWalls(board [][]Cell) {
	for i := 0; i < 6; i++ {
		length := rand.Intn(25) + 7
		y := rand.Intn(len(board))
		x := rand.Intn(len(board[0]))
		dir := rand.Intn(4)
		for j := 0; j < length; j++ {
			if !(y >= 0 && y < len(board) && x >= 0 && x < len(board[0])) ||
				!board[y][x].Clear {
				break
			}
			board[y][x].Clear = false
			board[y][x].Color = ColorWall
			change := rand.Intn(10)
			if change >= 8 {
				if change == 9 {
					dir = (dir + 1) % 4
				} else {
					dir = (dir - 1) % 4
				}
			}
			switch dir {
			case 0:
				y -= 1
				break
			case 1:
				x += 1
				break
			case 2:
				y += 1
				break
			case 3:
				x -= 1
				break
			}
		}
	}
}

func moveSnake(snake []Segment, board [][]Cell, lastDir *int, gameOverC chan bool) {
	score := uint64(0)
	renderScore(score)
	ticker := time.NewTicker(100 * time.Millisecond)
	food := false
	for {
		if !food {
			randX := rand.Intn(len(board[0]))
			randY := rand.Intn(len(board))
			for !board[randY][randX].Clear {
				randX = rand.Intn(len(board[0]))
				randY = rand.Intn(len(board))
			}
			board[randY][randX].Clear = false
			board[randY][randX].Color = ColorFood
			termbox.SetCell(randX, randY, 0x0020, termbox.ColorBlack, ColorFood)
			food = true
		}
		for i := 0; i < len(snake); i++ {
			if i == len(snake)-1 {
				board[snake[i].Y][snake[i].X].Color = ColorEmpty
				board[snake[i].Y][snake[i].X].Clear = true
				termbox.SetCell(snake[i].X, snake[i].Y, 0x0020, termbox.ColorBlack, ColorEmpty)
			}
			switch snake[i].Dir {
			case 0:
				snake[i].Y -= 1
				break
			case 1:
				snake[i].X += 1
				break
			case 2:
				snake[i].Y += 1
				break
			case 3:
				snake[i].X -= 1
				break
			}
			if i == 0 {
				if !board[snake[0].Y][snake[0].X].Clear { //collision
					if board[snake[0].Y][snake[0].X].Color == ColorFood { //with food
						food = false
						newX := snake[len(snake)-1].X
						newY := snake[len(snake)-1].Y
						switch snake[len(snake)-1].Dir {
						case 0:
							newY += 1
							break
						case 1:
							newX -= 1
							break
						case 2:
							newY -= 1
							break
						case 3:
							newX += 1
							break
						}
						snake = append(snake, Segment{newX, newY, snake[len(snake)-1].Dir})
						board[newY][newX].Color = ColorSnake
						board[newY][newX].Clear = false
						score++
						renderScore(score)
					} else { //with wall
						switch snake[0].Dir { //undo move of head segment
						case 0:
							snake[0].Y += 1
							break
						case 1:
							snake[0].X -= 1
							break
						case 2:
							snake[0].Y -= 1
							break
						case 3:
							snake[0].X += 1
							break
						}
						gameOver(snake, board)
						gameOverC <- true
						return
					}
				}
				board[snake[0].Y][snake[0].X].Color = ColorSnake
				board[snake[0].Y][snake[0].X].Clear = false
				termbox.SetCell(snake[0].X, snake[0].Y, 0x0020, termbox.ColorBlack, ColorSnake)
			}
		}
		for j := len(snake) - 1; j > 0; j-- {
			snake[j].Dir = snake[j-1].Dir
		}
		*lastDir = snake[0].Dir
		termbox.Flush()
		<-ticker.C
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	err := termbox.Init()
	if err != nil {
		log.Panicln(err)
	}
	defer termbox.Close()
	termbox.SetInputMode(termbox.InputEsc)

OutsideGameLoop:
	for {
		gameLoop()
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			if ev.Ch != 0 { //letter key
				switch ev.Ch {
				case 'q':
					return
				case 'n':
					continue OutsideGameLoop
				}
			} else {
				switch ev.Key {
				case termbox.KeyCtrlC:
					return
				}
			}
		case termbox.EventError:
			log.Panic(ev.Err)
			break
		}
	}
}

func gameLoop() {
	snake, board := initGame()
	lastDir := snake[0].Dir
	gameOverC := make(chan bool, 1)
	go moveSnake(snake, board, &lastDir, gameOverC)
	for {
		select {
		case <-gameOverC:
			return
		default:
			switch ev := termbox.PollEvent(); ev.Type {
			case termbox.EventKey:
				if ev.Ch == 0 { //not letter key
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
				} else {
					switch ev.Ch {
					case 'w':
						if lastDir != 2 {
							snake[0].Dir = 0
						}
						break
					case 'd':
						if lastDir != 3 {
							snake[0].Dir = 1
						}
						break
					case 's':
						if lastDir != 0 {
							snake[0].Dir = 2
						}
						break
					case 'a':
						if lastDir != 1 {
							snake[0].Dir = 3
						}
						break
					}
				}
				break
			case termbox.EventError:
				log.Panic(ev.Err)
				break
			}
		}
	}
}
