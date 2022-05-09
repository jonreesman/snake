package main

import (
	"fmt"
	"image/color"
	"math/rand"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func main() {
	/* Main Game App setup */
	a := app.New()
	w := a.NewWindow("Hello World")
	w.SetMaster()

	/* Set up for main menu/window 1 */
	clock := widget.NewLabel("")

	gameSize := 32
	dd := widget.NewSelect(
		[]string{"16x16", "32x32", "64x64"},
		func(s string) {
			switch s {
			case "16x16":
				gameSize = 16
			case "32x32":
				gameSize = 32
			case "64x64":
				gameSize = 64
			default:
				gameSize = 32
			}
			fmt.Printf("gameSize: %d ", gameSize)
		},
	)

	startButton := widget.NewButton("Start", func() {
		fmt.Println("Game open")
		gameWindow := a.NewWindow("Game")
		startGame(gameWindow, gameSize)
		gameWindow.SetMaster()
		w.Hide()
	})

	updateTime(clock)

	mainWindowContent := container.New(layout.NewVBoxLayout(), clock, dd, startButton, layout.NewSpacer())

	w.SetContent(mainWindowContent)
	w.Resize(fyne.NewSize(100, 100))

	go func() {
		for range time.Tick(time.Second) {
			updateTime(clock)
		}
	}()

	w.Show()
	a.Run()
}

func startGame(gameWindow fyne.Window, gameSize int) {
	gridLayout := make([]fyne.CanvasObject, gameSize*gameSize)
	for i := range gridLayout {
		gridLayout[i] = canvas.NewRectangle(color.White)
	}
	grid := container.New(layout.NewGridLayoutWithColumns(gameSize), gridLayout...)
	gameWindow.SetContent(grid)
	gameWindow.Resize(fyne.NewSize(600, 600))

	game := gameInstance{
		gridLayout:      gridLayout,
		gameSize:        gameSize,
		dotPosition:     rand.Intn(gameSize * gameSize),
		playerPosition:  gameSize * gameSize / 3,
		playerQueue:     make([]int, 1),
		playerLength:    1,
		playerDirection: Left,
		score:           0,
	}
	game.playerQueue[0] = gameSize * gameSize / 3
	game.gridLayout[game.dotPosition] = canvas.NewRectangle(color.RGBA{R: 255})

	gameWindow.Canvas().SetOnTypedKey(func(k *fyne.KeyEvent) {
		switch k.Name {
		case "Left":
			if game.playerDirection != Right {
				game.playerDirection = Left
			}
		case "Right":
			if game.playerDirection != Left {
				game.playerDirection = Right
			}
		case "Up":
			if game.playerDirection != Down {
				game.playerDirection = Up
			}
		case "Down":
			if game.playerDirection != Up {
				game.playerDirection = Down
			}
		}
	})
	go runGame(&gameWindow, &game, grid)
	gameWindow.Show()
}

type Direction int64

const (
	Up Direction = iota
	Down
	Left
	Right
)

type gameInstance struct {
	gridLayout      []fyne.CanvasObject
	gameSize        int
	dotPosition     int
	playerPosition  int //Grid is 32 blocks in length
	playerLength    int
	playerQueue     []int
	playerDirection Direction
	score           int
}

func runGame(gameWindow *fyne.Window, game *gameInstance, grid *fyne.Container) {
	for t := range time.Tick(time.Millisecond) {
		if t.UnixMilli()%100 != 0 {
			continue
		}
		switch game.playerDirection {
		case Up:
			game.playerPosition -= game.gameSize
			if game.playerPosition < 0 {
				game.playerPosition += game.gameSize * game.gameSize
				fmt.Println("Kill1")
				deathHandler(gameWindow, game.score)

			}
		case Down:
			game.playerPosition += game.gameSize
			if game.playerPosition > game.gameSize*game.gameSize-1 {
				game.playerPosition -= game.gameSize * game.gameSize
				fmt.Println("Kill2")
				deathHandler(gameWindow, game.score)

			}
		case Left:
			if game.playerPosition%game.gameSize == 0 {
				fmt.Println("Kill3")
				deathHandler(gameWindow, game.score)

			}
			game.playerPosition--
			if game.playerPosition < 0 {
				game.playerPosition = game.gameSize*game.gameSize - 1
				fmt.Println("Kill4")
				deathHandler(gameWindow, game.score)

			}
		case Right:
			if (game.playerPosition+1)%game.gameSize == 0 {
				fmt.Println("Kill5")
				deathHandler(gameWindow, game.score)

			}
			game.playerPosition++
			if game.playerPosition > game.gameSize*game.gameSize-1 {
				game.playerPosition = 0
				fmt.Println("Kill6")
				deathHandler(gameWindow, game.score)

			}
		}

		if game.playerPosition == game.dotPosition {
			game.playerLength++
			game.score++
			fmt.Printf("Score %d\n", game.score)

			game.dotPosition = rand.Intn(game.gameSize * game.gameSize)
			game.gridLayout[game.dotPosition] = canvas.NewRectangle(color.RGBA{R: 255})
		}

		if collisionCheck(game.playerPosition, game.playerQueue) {
			deathHandler(gameWindow, game.score)
		}

		game.playerQueue = append(game.playerQueue, game.playerPosition)

		game.gridLayout[game.playerPosition] = canvas.NewRectangle(color.Black)
		if game.playerLength < len(game.playerQueue) {
			game.gridLayout[game.playerQueue[0]] = canvas.NewRectangle(color.White)
			game.playerQueue = game.playerQueue[1:]

		}
		grid.Refresh()
	}
}

func deathHandler(gameWindow *fyne.Window, score int) {
	(*gameWindow).SetContent(widget.NewLabel("Score: " + strconv.Itoa(score)))
}

func collisionCheck(position int, playerQueue []int) bool {
	for _, p := range playerQueue {
		if position == p {
			return true
		}
	}
	return false
}

func updateTime(clock *widget.Label) {
	formatted := time.Now().Format("Time: 03:04:05")
	clock.SetText(formatted)
}
