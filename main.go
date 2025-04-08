package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"
)

type Tile int

const (
	Tile_Empty Tile = iota
	Tile_Wall
	Tile_Brick
	Tile_OneWay
)

var solid_tiles = []Tile{Tile_Brick, Tile_Wall}

type Direction int

const (
	Dir_Up Direction = iota
	Dir_Down
	Dir_Left
	Dir_Right
)

// NOTE: For now we'll hardcode to a level width of 18 tiles
const LEVEL_WIDTH = 18

var offsets = map[Direction]int{
	Dir_Up:    -LEVEL_WIDTH,
	Dir_Down:  LEVEL_WIDTH,
	Dir_Left:  -1,
	Dir_Right: 1,
}

// NOTE: Currently only a subset of tile types are supported
var byte_to_tile = map[byte]Tile{
	'B': Tile_Brick,
	'.': Tile_Empty,
	'1': Tile_OneWay,
	'X': Tile_Wall,
	'S': Tile_Empty,
}

// NOTE: Currently only a subset of tile types are supported
var tile_to_byte = map[Tile]byte{
	Tile_Brick:  'B',
	Tile_Empty:  '.',
	Tile_OneWay: '1',
	Tile_Wall:   'X',
}

var dir_to_byte = map[Direction]byte{
	Dir_Up:    '^',
	Dir_Down:  'v',
	Dir_Left:  '<',
	Dir_Right: '>',
}

type Level struct {
	tiles              []Tile
	player_index       int
	previous_direction Direction
}

func (level Level) move(dir Direction) (*Level, error) {
	if dir == level.previous_direction {
		return nil, errors.New("cannot move in the same direction twice")
	}

	tiles := make([]Tile, len(level.tiles))
	copy(tiles, level.tiles)

	new_level := Level{
		tiles:              tiles,
		previous_direction: dir,
	}

	offset := offsets[dir]
	player := level.player_index

	touching_index := player + offsets[level.previous_direction]
	if tiles[touching_index] == Tile_Brick {
		tiles[touching_index] = Tile_Empty
	}

	// NOTE: Because the player cannot reach the edges of the level, we don't
	// have to worry about bounds checking the tiles array
	for next_tile := tiles[player+offset]; !next_tile.is_solid(); next_tile = tiles[player+offset] {
		if tiles[player] == Tile_OneWay {
			tiles[player] = Tile_Wall
		}
		player = player + offset
	}

	new_level.player_index = player
	return &new_level, nil
}

func (level Level) to_string() string {
	level_width := LEVEL_WIDTH
	slice_start := 0
	slice_end := level_width

	builder := strings.Builder{}

	for slice_end < len(level.tiles) {
		next_slice := level.tiles[slice_start:slice_end]
		for i, t := range next_slice {
			if slice_start+i == level.player_index {
				builder.WriteByte(dir_to_byte[level.previous_direction])
			} else {
				builder.WriteByte(tile_to_byte[t])
			}
		}
		builder.WriteByte('\n')

		slice_start += level_width
		slice_end += level_width
	}

	return builder.String()
}

func load_level(path string) (*Level, error) {
	file, err := os.Open(path)

	if err != nil {
		return nil, err
	}
	defer file.Close()

	level := Level{
		// NOTE: Assume the player starts facing down. This is usually true.
		previous_direction: Dir_Down,
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		bytes := scanner.Bytes()
		for _, b := range bytes {
			level.tiles = append(level.tiles, byte_to_tile[b])
			if b == 'S' {
				level.player_index = len(level.tiles) - 1
			}
		}
	}

	if err = scanner.Err(); err != nil {
		return nil, err
	}

	return &level, nil
}

func (tile Tile) is_solid() bool {
	return slices.Contains(solid_tiles, tile)
}

func main() {
	level, err := load_level("level23.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(level.to_string())

	level, _ = level.move(Dir_Up)
	level, _ = level.move(Dir_Left)
	level, err = level.move(Dir_Up)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(level.to_string())
}
