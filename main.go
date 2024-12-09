package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth      = 320
	screenHeight     = 240
	initialLiveCells = screenWidth * screenHeight / 10
	cellSize         = 4 // RGBA bytes per cell
)

type World struct {
	area   []bool
	width  int
	height int
	buffer []bool // Reused for state updates
}

// NewWorld creates a new world and initializes it.
func NewWorld(width, height int, maxInitLiveCells int) *World {
	w := &World{
		area:   make([]bool, width*height),
		width:  width,
		height: height,
		buffer: make([]bool, width*height),
	}
	w.init(maxInitLiveCells)
	return w
}

// init populates the world with random live cells.
func (w *World) init(maxLiveCells int) {
	rand.Seed(time.Now().UnixNano()) // Ensure randomness per run
	for i := 0; i < maxLiveCells; i++ {
		x := rand.Intn(w.width)
		y := rand.Intn(w.height)
		w.area[y*w.width+x] = true
	}
}

// neighbourCount calculates the number of live neighbors for cell (x, y).
func (w *World) neighbourCount(x, y int) int {
	count := 0
	for j := -1; j <= 1; j++ {
		for i := -1; i <= 1; i++ {
			if i == 0 && j == 0 {
				continue
			}
			nx, ny := x+i, y+j
			if nx >= 0 && ny >= 0 && nx < w.width && ny < w.height {
				if w.area[ny*w.width+nx] {
					count++
				}
			}
		}
	}
	return count
}

// Update progresses the world's state by one generation.
func (w *World) Update() {
	for y := 0; y < w.height; y++ {
		for x := 0; x < w.width; x++ {
			idx := y*w.width + x
			neighbors := w.neighbourCount(x, y)
			w.buffer[idx] = false
			if neighbors == 3 || (neighbors == 2 && w.area[idx]) {
				w.buffer[idx] = true
			}
		}
	}
	// Swap buffers
	w.area, w.buffer = w.buffer, w.area
}

// Draw converts the world's state to RGBA pixels.
func (w *World) Draw(pix []byte) {
	for i, live := range w.area {
		offset := i * cellSize
		if live {
			pix[offset], pix[offset+1], pix[offset+2], pix[offset+3] = 0xff, 0xff, 0xff, 0xff
		} else {
			pix[offset], pix[offset+1], pix[offset+2], pix[offset+3] = 0, 0, 0, 0
		}
	}
}

type Game struct {
	world  *World
	pixels []byte
}

func (g *Game) Update() error {
	g.world.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.pixels == nil {
		g.pixels = make([]byte, screenWidth*screenHeight*cellSize)
	}
	g.world.Draw(g.pixels)
	screen.WritePixels(g.pixels)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	game := &Game{
		world: NewWorld(screenWidth, screenHeight, initialLiveCells),
	}

	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("~Game of Life~")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
