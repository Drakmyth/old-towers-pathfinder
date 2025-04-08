package main

import (
	"bufio"
	"fmt"
	"os"
)

type Tile int

const (
	Tile_Empty Tile = iota
	Tile_Wall
	Tile_Brick
	Tile_OneWay
)

type Level []Tile

func load_level(path string) (*Level, error) {
	file, err := os.Open(path)

	if err != nil {
		return nil, err
	}
	defer file.Close()

	level := Level{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		bytes := scanner.Bytes()
		for _, b := range bytes {
			level = append(level, byte_to_tile(b))
		}
	}

	if err = scanner.Err(); err != nil {
		return nil, err
	}

	return &level, nil
}

func byte_to_tile(b byte) Tile {
	switch b {
	case 'X':
		return Tile_Wall
	case 'B':
		return Tile_Brick
	case '1':
		return Tile_OneWay
	default:
		return Tile_Empty
	}
}

func main() {
	_, err := load_level("level23.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
}
