// Maze generator in Go
// Joe Wingbermuehle
// 2012-08-07

package main

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	wall  = 1
	space = 0
)

// Maze type
type Maze struct {
	width, height uint8
	data          [][]byte
}

/** Create an empty maze.
 * @param w The width (must be odd).
 * @param h The height (must be odd).
 */
func newMaze(w uint8, h uint8) *Maze {
	m := Maze{w, h, make([][]byte, h)}
	for y := range m.data {
		m.data[y] = make([]byte, w)
		for x := range m.data[y] {
			m.data[y][x] = wall
		}
	}
	for x := uint8(0); x < w; x++ {
		m.data[0][x], m.data[h-1][x] = space, space
	}
	for y := uint8(0); y < h; y++ {
		m.data[y][0], m.data[y][w-1] = space, space
	}
	return &m
}

/** Start carving a maze at the specified coordinates. */
func carveMaze(m *Maze, r *rand.Rand, x int, y int) {
	directions := [][]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
	d := r.Intn(4)
	for i := 0; i < 4; i++ {
		dx, dy := directions[d][0], directions[d][1]
		ax, ay := x+dx, y+dy
		bx, by := ax+dx, ay+dy
		if m.data[ay][ax] == wall && m.data[by][bx] == wall {
			m.data[ay][ax], m.data[by][bx] = space, space
			carveMaze(m, r, bx, by)
		}
		d = (d + 1) % 4
	}
}

/** Generate a maze. */
func generateMaze(m *Maze) {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	m.data[2][2] = space
	carveMaze(m, r, 2, 2)
	m.data[1][2] = space
	m.data[m.height-2][m.width-3] = space
}

/** Show a generated maze. */
func showMaze(m *Maze) {
	wallCount := 0
	for w := m.width - 1; w > 0; w-- {
		for h := uint8(0); h < m.height; h++ {
			//fmt.Println("w,h", w, h)
			if m.data[h][w] == wall {
				fmt.Printf("[]")
				wallCount++
			} else {
				fmt.Printf("  ")
			}
		}
		fmt.Printf("\n")
	}
	fmt.Println("wallCount", wallCount)
}

// Maze start
func designMaze(w, h uint8) *Maze {
	m := newMaze(w, h)
	generateMaze(m)
	showMaze(m)
	return m
}
