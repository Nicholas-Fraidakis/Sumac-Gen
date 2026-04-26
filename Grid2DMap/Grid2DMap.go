package grid2dmap

import (
	"fmt"
	"math"
	"slices"
)

type CompassDirectionVector byte

func (dir CompassDirectionVector) StepInX() int {
	switch dir {
	case WEST, NORTHWEST, SOUTHWEST:
		return -1
	case EAST, NORTHEAST, SOUTHEAST:
		return 1
	default:
		return 0
	}
}

func (dir CompassDirectionVector) StepInY() int {
	switch dir {
	case NORTH, NORTHWEST, NORTHEAST:
		return -1
	case SOUTH, SOUTHWEST, SOUTHEAST:
		return 1
	default:
		return 0
	}
}

const (
	NORTH CompassDirectionVector = iota
	SOUTH
	WEST
	EAST
	NORTHWEST
	NORTHEAST
	SOUTHWEST
	SOUTHEAST
)

type Coordinate struct {
	X, Y int
}

func (cord Coordinate) GetCompassDirectionVector() CompassDirectionVector {
	var direction CompassDirectionVector = 0xff

	if cord.X < 0 {
		direction = WEST
	} else if cord.X > 0 {
		direction = EAST
	}

	if cord.Y < 0 {
		if direction == 0xff {
			direction = NORTH
		} else {
			direction += 2
		}

	} else if cord.Y > 0 {
		if direction == 0xff {
			direction = SOUTH
		} else {
			direction += 4
		}
	}

	return direction
}

func (cord Coordinate) Hypot() float64 {
	return math.Sqrt(float64(cord.X*cord.X) + float64(cord.Y*cord.Y))
}

func SortCoordinateListByClosest(coordinates []Coordinate, start_coord Coordinate) []Coordinate {
	sorted_coords := make([]Coordinate, 0, len(coordinates))
	remaining_coords := slices.Clone(coordinates)

	current_coord := start_coord

	for i := range len(coordinates) {
		var closest_coordinate_index = -1
		var closest_coordinate_distance float64 = 0xffffffff
		for j := range len(remaining_coords) {
			if closest_coordinate_index == -1 {
				closest_coordinate_index = j
				continue
			}

			if i == j {
				continue
			}

			var checked_distance = Coordinate{
				X: current_coord.X - coordinates[j].X,
				Y: current_coord.Y - coordinates[j].Y,
			}.Hypot()

			if checked_distance < closest_coordinate_distance {
				closest_coordinate_index = j
				closest_coordinate_distance = checked_distance
			}

		}

		current_coord = remaining_coords[closest_coordinate_index]
		sorted_coords = append(sorted_coords, current_coord)
		remaining_coords = slices.Delete(remaining_coords, closest_coordinate_index, closest_coordinate_index+1)
	}

	return sorted_coords
}

type Grid2DMapElement uint32

const (
	EMPTY         Grid2DMapElement = 0
	OUT_OF_BOUNDS Grid2DMapElement = 0xffffffff
)

type Grid2DMap struct {
	underlying_map [][]Grid2DMapElement
}

func New(x, y int) Grid2DMap {
	underlying_map := make([][]Grid2DMapElement, y)

	for i := range y {
		underlying_map[i] = make([]Grid2DMapElement, x)
	}

	return Grid2DMap{
		underlying_map: underlying_map,
	}
}

func (grid *Grid2DMap) Clone() Grid2DMap {
	var cloned_grid Grid2DMap

	cloned_grid.underlying_map = make([][]Grid2DMapElement, grid.Height())

	for y := range grid.Height() {
		cloned_grid.underlying_map[y] = make([]Grid2DMapElement, grid.Width())
		for x := range grid.Width() {
			cloned_grid.underlying_map[y][x] = grid.underlying_map[y][x]
		}
	}

	return cloned_grid
}

func (grid *Grid2DMap) CloneTo(dest *Grid2DMap) {
	w1, h1 := grid.Size()
	w2, h2 := dest.Size()
	if w1 != w2 || h1 != h2 {
		return
	}

	for y := range grid.Height() {
		for x := range grid.Width() {
			dest.underlying_map[y][x] = grid.underlying_map[y][x]
		}
	}
}

func NewFromSlice(grid_slice [][]Grid2DMapElement) Grid2DMap {
	underlying_map := grid_slice

	return Grid2DMap{
		underlying_map: underlying_map,
	}
}

func (grid *Grid2DMap) Get(x, y int) Grid2DMapElement {
	if y >= len(grid.underlying_map) || y < 0 {
		return OUT_OF_BOUNDS
	}
	if x >= len(grid.underlying_map[0]) || x < 0 {
		return OUT_OF_BOUNDS
	}

	return grid.underlying_map[y][x]
}

func (grid *Grid2DMap) Set(x, y int, value Grid2DMapElement) {
	if y >= len(grid.underlying_map) || y < 0 {
		return
	}
	if x >= len(grid.underlying_map[0]) || x < 0 {
		return
	}

	grid.underlying_map[y][x] = value
}

func (grid *Grid2DMap) PositionTo1DIndex(x, y int) int {
	if x > grid.Width() || x < 0 ||
		y > grid.Height() || y < 0 {
		return grid.Width()*grid.Height() + 1
	}

	return x + y*grid.Width()
}
func (grid *Grid2DMap) Size() (x, y int) {
	return grid.Width(), grid.Height()
}

func (grid *Grid2DMap) Width() int {
	if len(grid.underlying_map) <= 0 {
		return 0
	}

	return len(grid.underlying_map[0])
}

func (grid *Grid2DMap) Height() int {
	return len(grid.underlying_map)
}

func (grid *Grid2DMap) Print() {
	for i := range grid.underlying_map {
		fmt.Println(grid.underlying_map[i])
	}
}

func (grid *Grid2DMap) DeleteAll(target_elements Grid2DMapElement) {
	for y := range grid.Height() {
		for x := range grid.Width() {
			if grid.Get(x, y) == target_elements {
				grid.Set(x, y, EMPTY)
			}
		}
	}
}

func (grid *Grid2DMap) ReplaceAll(find, replace Grid2DMapElement) {
	for _, row := range grid.underlying_map {
		for i, element := range row {
			if element == find {
				row[i] = replace
			}
		}
	}
}

type FindAndRunEntry struct {
	Element  Grid2DMapElement
	Function func(x, y int, tile Grid2DMapElement)
}

func FindAndRun(grid_map Grid2DMap, find_and_run_entries []FindAndRunEntry) {
	for x := range grid_map.Width() {
		for y := range grid_map.Height() {
			var checked_tile Grid2DMapElement = grid_map.Get(x, y)

			for i := range find_and_run_entries {
				if find_and_run_entries[i].Element == checked_tile {
					find_and_run_entries[i].Function(x, y, checked_tile)
				}
			}
			/*// Looking for jim_bob positioning
			if checked_tile >= START_POSITION && checked_tile < END_POSITION {
				if grid_map.Start_tile != grid2dmap.EMPTY {
					panic("Woah there, there are multiple start positions")
				}
				grid_map.Karel_start_position = grid2dmap.Coordinate{X: x, Y: y}
				grid_map.Start_tile = checked_tile
			}

			if checked_tile >= END_POSITION && checked_tile < END_OF_COMMON_TILES {
				if grid_map.End_tile != grid2dmap.EMPTY {
					panic("Woah there, there are multiple end positions")
				}
				grid_map.Karel_end_position = grid2dmap.Coordinate{X: x, Y: y}
				grid_map.End_tile = checked_tile
			}*/
		}
	}
}

func FindAndRunSingle(grid_map Grid2DMap, find_and_run_entry FindAndRunEntry) {
	FindAndRun(grid_map, []FindAndRunEntry{find_and_run_entry})
}
