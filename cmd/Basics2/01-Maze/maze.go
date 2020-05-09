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
	width, height int
	data          [][]byte
}

/** Create an empty maze.
 * @param w The width (must be odd).
 * @param h The height (must be odd).
 */
func newMaze(w int, h int) *Maze {
	m := Maze{w, h, make([][]byte, h)}
	for y := range m.data {
		m.data[y] = make([]byte, w)
		for x := range m.data[y] {
			m.data[y][x] = wall
		}
	}
	for x := 0; x < w; x++ {
		m.data[0][x], m.data[h-1][x] = space, space
	}
	for y := 0; y < h; y++ {
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
	for y := 0; y < m.height; y++ {
		for x := 0; x < m.width; x++ {
			if m.data[y][x] == wall {
				fmt.Printf("[]")
			} else {
				fmt.Printf("  ")
			}
		}
		fmt.Printf("\n")
	}
}

// Maze start
func mainMaze(w, h int) *Maze {
	m := newMaze(w, h)
	// for i := 0; i < 23; i++ {
	// 	fmt.Println(m.data[i])
	// }
	// fmt.Println()
	generateMaze(m)
	// for i := 0; i < 23; i++ {
	// 	fmt.Println(m.data[i])
	// }

	showMaze(m)
	return m
}
