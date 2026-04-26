package main

import (
	"fmt"
	"strings"
	conversion "sumac-gen/Conversion"
	grid2dmap "sumac-gen/Grid2DMap"
)

func main() {
	start_map := grid2dmap.NewFromSlice(
		[][]grid2dmap.Grid2DMapElement{
			{1, 1, 1, conversion.INTERACTABLE_PLACED_BALLS + 5, 0, 0},
			{0, 0, 1, 0, 1, 0},
			{conversion.START_POSITION, 0, 0, 0, 0, 0},
			{0, 1, 1, 1, 0, 0},
			{1, conversion.INTERACTABLE_PLACED_BALLS + 3, 0, 0, 1, 0},
			{1, 0, 1, 0, 0, 0},
		},
	)

	final_map := grid2dmap.NewFromSlice(
		[][]grid2dmap.Grid2DMapElement{
			{1, 1, 1, conversion.INTERACTABLE_PLACED_BALLS + 10, 0, 0},
			{0, 0, 1, 0, 1, 0},
			{0, 0, 0, 0, 0, 0},
			{0, 1, 1, 1, 0, 0},
			{1, conversion.INTERACTABLE_PLACED_BALLS + 6, 0, 0, 1, 0},
			{1, 0, 1, 0, 0, conversion.END_POSITION},
		},
	)

	core_map, specs := conversion.GetObjectiveSpecs(start_map, final_map)

	object_groups := conversion.NewObjectiveGroupsFromSpecs(specs, core_map)

	code := ""
	for i := range object_groups {
		code += object_groups[i].GenerateCode()
	}

	code = strings.ReplaceAll(code, ";", "\n")

	fmt.Println(code)
}
