// runner.go - Бесконечный бегун на Go
package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/eiannone/keyboard"
)

const (
	WIDTH      = 40
	HEIGHT     = 10
	GROUND     = HEIGHT - 1
	PLAYER     = '@'
	OBSTACLE   = '#'
	RECORD_FILE = "runner_record.json"
)

type RunnerGame struct {
	playerX   int
	playerY   float64
	playerVy  float64
	onGround  bool
	obstacles []int
	score     int
	speed     int
	gameOver  bool
	paused    bool
	running   bool
	record    int
	rand      *rand.Rand
}

func NewRunnerGame() *RunnerGame {
	return &RunnerGame{
		playerX:   5,
		playerY:   GROUND,
		playerVy:  0,
		onGround:  true,
		obstacles: make([]int, 0),
		score:     0,
		speed:     1,
		gameOver:  false,
		paused:    false,
		running:   true,
		record:    loadRecord(),
		rand:      rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func loadRecord() int {
	file, err := os.Open(RECORD_FILE)
	if err != nil {
		return 0
	}
	defer file.Close()
	var data map[string]int
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&data); err != nil {
		return 0
	}
	if val, ok := data["record"]; ok {
		return val
	}
	return 0
}

func saveRecord(record int) {
	data := map[string]int{"record": record}
	file, _ := os.Create(RECORD_FILE)
	defer file.Close()
	encoder := json.NewEncoder(file)
	encoder.Encode(data)
}

func (g *RunnerGame) jump() {
	if g.onGround && !g.gameOver && !g.paused {
		g.playerVy = -3.0
		g.onGround = false
	}
}

func (g *RunnerGame) spawnObstacle() {
	g.obstacles = append(g.obstacles, WIDTH-1)
}

func (g *RunnerGame) updatePhysics() {
	g.playerVy += 0.2
	g.playerY += g.playerVy
	if g.playerY >= GROUND {
		g.playerY = GROUND
		g.playerVy = 0
		g.onGround = true
	}
	if g.playerY < 0 {
		g.playerY = 0
		g.playerVy = 0
	}
}

func (g *RunnerGame) update() {
	if g.gameOver || g.paused {
		return
	}

	// Движение препятствий
	for i := len(g.obstacles) - 1; i >= 0; i-- {
		g.obstacles[i]--
		if g.obstacles[i] < 0 {
			g.obstacles = append(g.obstacles[:i], g.obstacles[i+1:]...)
		}
	}

	// Столкновение
	playerFloor := int(g.playerY)
	for _, x := range g.obstacles {
		if x == g.playerX && (playerFloor >= GROUND-1 || playerFloor == GROUND-1) {
			g.gameOver = true
			if g.score > g.record {
				g.record = g.score
				saveRecord(g.record)
			}
			return
		}
	}

	g.score++
	g.speed = 1 + g.score/10

	if g.rand.Float64() < 0.2+min(0.4, float64(g.score)/200.0) {
		g.spawnObstacle()
	}

	g.updatePhysics()
}

func (g *RunnerGame) draw() {
	clearScreen()
	fmt.Println(stringRepeat("═", WIDTH))
	fmt.Printf("  Счёт: %d   Рекорд: %d   Скорость: %d\n", g.score, g.record, g.speed)
	fmt.Println(stringRepeat("═", WIDTH))

	grid := make([][]rune, HEIGHT)
	for i := range grid {
		grid[i] = make([]rune, WIDTH)
		for j := range grid[i] {
			grid[i][j] = ' '
		}
	}
	for _, x := range g.obstacles {
		if x >= 0 && x < WIDTH {
			grid[GROUND][x] = OBSTACLE
		}
	}
	y := int(g.playerY)
	if y >= 0 && y < HEIGHT {
		grid[y][g.playerX] = PLAYER
	}

	for y := 0; y < HEIGHT; y++ {
		fmt.Print('│')
		for x := 0; x < WIDTH; x++ {
			fmt.Print(string(grid[y][x]))
		}
		fmt.Println('│')
	}
	fmt.Println(stringRepeat("═", WIDTH))
	status := "ИГРА"
	if g.paused {
		status = "ПАУЗА"
	}
	fmt.Printf("  %s  |  Пробел/↑ - прыжок  |  P - пауза  |  R - рестарт  |  Q - выход\n", status)
}

func (g *RunnerGame) reset() {
	g.playerY = GROUND
	g.playerVy = 0
	g.onGround = true
	g.obstacles = nil
	g.score = 0
	g.speed = 1
	g.gameOver = false
	g.paused = false
}

func clearScreen() {
	cmd := exec.Command("clear")
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func stringRepeat(s string, n int) string {
	res := ""
	for i := 0; i < n; i++ {
		res += s
	}
	return res
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func (g *RunnerGame) handleInput() {
	for g.running {
		char, key, err := keyboard.GetKey()
		if err != nil {
			continue
		}
		switch key {
		case keyboard.KeyArrowUp:
			g.jump()
		case keyboard.KeySpace:
			g.jump()
		default:
			switch char {
			case 'p', 'P':
				g.paused = !g.paused
			case 'r', 'R':
				g.reset()
			case 'q', 'Q':
				g.running = false
				return
			}
		}
	}
}

func (g *RunnerGame) run() {
	go g.handleInput()

	lastUpdate := time.Now()
	ticker := time.NewTicker(16 * time.Millisecond)
	defer ticker.Stop()

	for g.running {
		select {
		case <-ticker.C:
			now := time.Now()
			if now.Sub(lastUpdate) >= time.Second/60 {
				g.update()
				g.draw()
				lastUpdate = now
			}
		}
	}
}

func main() {
	game := NewRunnerGame()
	game.run()
	fmt.Println("Игра завершена.")
}
