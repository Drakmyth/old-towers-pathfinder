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

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	if err = scanner.Err(); err != nil {
		return nil, err
	}

	return &Level{}, nil
}

func main() {
	_, err := load_level("level23.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
}
