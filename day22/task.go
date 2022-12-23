package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

var input = "day22/input.txt"

type Task struct {
	grid       map[Point]string
	start      Point
	directions []string
	max        Point
	facing     string
	cube       bool
	face       string
}

type Point struct {
	x, y int
}

func main() {
	start := time.Now()

	part1, part2 := run()

	log.Println(fmt.Sprintf("Part 1: %d == 144244", part1))

	log.Println(fmt.Sprintf("Part 2: %d 107388 < x < 151106, not 113298, 116142, 117169", part2))

	elapsed := time.Since(start)
	log.Printf("Binomial took %s", elapsed)
}

func run() (int, int) {
	task := readInput()

	//fmt.Println(task)

	//part1 := task.mapPath()

	task.cube = true
	part2 := task.mapPath()

	return 0, part2
	//return part1, part2
}

func (t *Task) mapPath() int {
	position := t.start
	t.facing = "right"
	t.face = "A"
	for _, direction := range t.directions {
		//fmt.Println(position)
		switch direction {
		case "R":
			t.rotateClockwise()
		case "L":
			t.rotateCounterClockwise()
		default:
			position = t.move(position, direction)
		}
	}
	//fmt.Println(position, facing)
	return t.password(position)
}

func (t *Task) password(position Point) int {
	var facingValue int
	switch t.facing {
	case "right":
		facingValue = 0
	case "left":
		facingValue = 2
	case "up":
		facingValue = 3
	case "down":
		facingValue = 1
	}
	return 1000*(position.y+1) + 4*(position.x+1) + facingValue
}

func (t *Task) move(position Point, direction string) Point {
	moves, _ := strconv.Atoi(direction)

	for i := 0; i < moves; i++ {
		switch t.facing {
		case "right":
			position = t.moveStep(position, Point{1, 0})
		case "up":
			position = t.moveStep(position, Point{0, -1})
		case "left":
			position = t.moveStep(position, Point{-1, 0})
		case "down":
			position = t.moveStep(position, Point{0, 1})
		}
	}

	return position
}

func (t *Task) moveStep(position Point, next Point) Point {
	if t.cube == true {
		return t.moveStepCube(position, next)
	} else {
		return t.moveStepPlane(position, next)
	}
}

func (t *Task) moveStepCube(position Point, next Point) Point {
	//fmt.Println(fmt.Sprintf("%+v", position), t.face)
	nextPosition, blocked := t.getNextCubePosition(position, next)
	if blocked {
		return position
	}
	//fmt.Println(fmt.Sprintf("%+v", nextPosition), t.face)
	if t.grid[nextPosition] == "" {
		panic(fmt.Sprintf("%+v", nextPosition) + " " + t.face)
	}
	return nextPosition
}

func (t *Task) getNextCubePosition(position Point, next Point) (Point, bool) {
	position = Point{position.x + next.x, position.y + next.y}

	nextPosition := position
	switch t.face {
	case "A":
		// 50 < x < 100
		// 0 < y < 50
		if position.x < 50 {
			nextPosition = Point{0, 149 - position.y%50}
			if t.grid[nextPosition] == "#" {
				return position, true
			}
			t.checkIsAnEdge(nextPosition)
			t.face = "E"
			t.facing = "right"
		}
		if position.x > 99 {
			t.face = "B"
		}
		if position.y < 0 {
			nextPosition = Point{0, 150 + position.x%50}
			if t.grid[nextPosition] == "#" {
				return position, true
			}
			t.face = "F"
			t.facing = "right"
			t.checkIsAnEdge(nextPosition)

		}
		if position.y > 49 {
			t.face = "C"
		}
	case "B":
		// 100 < x < 150
		// 0 < y < 50
		if position.x < 100 {
			t.face = "A"
		}
		if position.x > 149 {
			nextPosition = Point{99, 149 - position.y%50}
			if t.grid[nextPosition] == "#" {
				return position, true
			}
			t.face = "D"
			t.facing = "left"
			t.checkIsAnEdge(nextPosition)
		}
		if position.y < 0 {
			nextPosition = Point{position.x % 50, 199}
			if t.grid[nextPosition] == "#" {
				return position, true
			}
			t.face = "F"
			t.facing = "up"
			t.checkIsAnEdge(nextPosition)
		}
		if position.y > 49 {
			nextPosition = Point{99, 50 + position.x%50}
			if t.grid[nextPosition] == "#" {
				return position, true
			}
			t.face = "C"
			t.facing = "left"
			t.checkIsAnEdge(nextPosition)
		}
	case "C":
		// 50 < x < 100
		// 50 < y < 100
		if position.x < 50 {
			nextPosition = Point{position.y % 50, 100}
			if t.grid[nextPosition] == "#" {
				return position, true
			}
			t.face = "E"
			t.facing = "down"
			t.checkIsAnEdge(nextPosition)
		}
		if position.x > 99 {
			nextPosition = Point{100 + position.y%50, 49}
			if t.grid[nextPosition] == "#" {
				return position, true
			}
			t.face = "B"
			t.facing = "up"
			t.checkIsAnEdge(nextPosition)
		}
		if position.y < 50 {
			t.face = "A"
		}
		if position.y > 99 {
			t.face = "D"
		}
	case "D":
		// 50 < x < 100
		// 100 < y < 150
		if position.x < 50 {
			t.face = "E"
		}
		if position.x > 99 {
			nextPosition = Point{149, 49 - position.y%50}
			if t.grid[nextPosition] == "#" {
				return position, true
			}
			t.face = "B"
			t.facing = "left"
			t.checkIsAnEdge(nextPosition)
		}
		if position.y < 100 {
			t.face = "C"
		}
		if position.y > 149 {
			nextPosition = Point{49, 150 + position.x%50}
			if t.grid[nextPosition] == "#" {
				return position, true
			}
			t.face = "F"
			t.facing = "left"
			t.checkIsAnEdge(nextPosition)
		}
	case "E":
		// 0 < x < 50
		// 100 < y < 150
		if position.x < 0 {
			nextPosition = Point{50, 49 - position.y%50}
			if t.grid[nextPosition] == "#" {
				return position, true
			}
			t.face = "A"
			t.facing = "right"
			t.checkIsAnEdge(nextPosition)
		}
		if position.x > 49 {
			t.face = "D"
		}
		if position.y < 100 {
			nextPosition = Point{50, 50 + position.x%50}
			if t.grid[nextPosition] == "#" {
				return position, true
			}
			t.face = "C"
			t.facing = "right"
			t.checkIsAnEdge(nextPosition)
		}
		if position.y > 149 {
			t.face = "F"
		}
	case "F":
		// 0 < x < 50
		// 150 < y < 200
		if position.x < 0 {
			nextPosition = Point{50 + position.y%50, 0}
			if t.grid[nextPosition] == "#" {
				return position, true
			}
			t.face = "A"
			t.facing = "down"
			t.checkIsAnEdge(nextPosition)
		}
		if position.x > 49 {
			nextPosition = Point{50 + position.y%50, 149}
			if t.grid[nextPosition] == "#" {
				return position, true
			}
			t.face = "D"
			t.facing = "up"
			t.checkIsAnEdge(nextPosition)
		}
		if position.y < 150 {
			t.face = "E"
		}
		if position.y > 199 {
			nextPosition = Point{100 + position.x%50, 0}
			if t.grid[nextPosition] == "#" {
				return position, true
			}
			t.face = "B"
			t.facing = "down"
			t.checkIsAnEdge(nextPosition)
		}
	}

	return nextPosition, false
}

func (t *Task) moveStepPlane(position Point, next Point) Point {
	var blocked bool
	var nextPosition Point
	nextPosition = Point{position.x + next.x, position.y + next.y}
	switch t.grid[nextPosition] {
	case "":
		nextPosition, blocked = t.getNextPosition(nextPosition, next)
		if blocked {
			return position
		}
		position = nextPosition
	case "#":
		return position
	case ".":
		position = nextPosition
	}
	return position
}

func (t *Task) getNextPosition(position, next Point) (Point, bool) {
	var nextPosition Point
	if next.x != 0 {
		nextPosition = Point{getNext(position.x+next.x, t.max.x), position.y}
	} else {
		nextPosition = Point{position.x, getNext(position.y+next.y, t.max.y)}
	}

	switch t.grid[nextPosition] {
	case "":
		return t.getNextPosition(nextPosition, next)
	case "#":
		return position, true
	case ".":
		return nextPosition, false
	}
	panic("invalid cell")
}

func getNext(next, max int) int {
	if next < 0 {
		next = max
	}
	if next > max {
		next = 0
	}
	return next
}

func (t *Task) rotateCounterClockwise() {
	switch t.facing {
	case "right":
		t.facing = "up"
	case "up":
		t.facing = "left"
	case "left":
		t.facing = "down"
	case "down":
		t.facing = "right"
	}
}

func (t *Task) rotateClockwise() {
	switch t.facing {
	case "right":
		t.facing = "down"
	case "down":
		t.facing = "left"
	case "left":
		t.facing = "up"
	case "up":
		t.facing = "right"
	}
}

func (t *Task) checkIsAnEdge(p Point) {
	if t.grid[Point{p.x + 1, p.y}] != "" && t.grid[Point{p.x - 1, p.y}] != "" &&
		t.grid[Point{p.x, p.y + 1}] != "" && t.grid[Point{p.x, p.y - 1}] != "" {
		panic(fmt.Sprintf("not an edge: %+v, %s", p, t.face))
	}
}

func readInput() Task {
	file, err := os.Open(input)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	task := Task{
		grid: map[Point]string{},
		cube: false,
	}

	y := 0
	gridSetting := true
	startSet := false
	maxX := 0

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		in := scanner.Text()
		if in != "" {
			if gridSetting {
				for i, c := range in {
					if string(c) != " " {
						if !startSet {
							task.start = Point{i, y}
							startSet = true
						}
						task.grid[Point{i, y}] = string(c)
						if maxX < i {
							maxX = i
						}
					}
				}
				y++
			} else {
				moves := ""
				for _, c := range in {
					_, err = strconv.Atoi(string(c))
					if err != nil {
						task.directions = append(task.directions, moves, string(c))
						moves = ""
					} else {
						moves += string(c)
					}
				}
				if moves != "" {
					task.directions = append(task.directions, moves)
				}
			}
		} else {
			gridSetting = false
		}
	}

	task.max = Point{maxX, y - 1}

	return task
}