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
			{0, 0, 0, 0, 0, 0, 0},
			{0, grid2dmap.Grid2DMapElement(conversion.PAINT_FLAG | (1 << conversion.PAINT_ARG_OFFSET)), 0, 0, 0, 0, 0},
			{0, 0, grid2dmap.Grid2DMapElement(conversion.BALL_FLAG | (4 << conversion.BALL_ARG_OFFSET)), 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0},
			{grid2dmap.Grid2DMapElement(conversion.POSITION_FLAG), 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0},
		},
	)

	final_map := grid2dmap.NewFromSlice(
		[][]grid2dmap.Grid2DMapElement{
			{0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0},
			{0, 0, grid2dmap.Grid2DMapElement(conversion.BALL_FLAG | (8 << conversion.BALL_ARG_OFFSET)), 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0},
			{grid2dmap.Grid2DMapElement(conversion.POSITION_FLAG) | grid2dmap.Grid2DMapElement(conversion.BALL_FLAG|(8<<conversion.BALL_ARG_OFFSET)) | 2, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0},
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
