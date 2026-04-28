package conversion

import (
	"fmt"
	"log"
	"slices"
	astarpathfinding "sumac-gen/AStarPathFinding"
	grid2dmap "sumac-gen/Grid2DMap"
	utils "sumac-gen/Utils"
)

type Movement struct {
	Direction grid2dmap.CompassDirectionVector
	Steps     uint32
}

func (movement Movement) Generate() string {
	switch movement.Direction {
	case grid2dmap.NORTH:
		return fmt.Sprintf("move_up(%d);", movement.Steps)
	case grid2dmap.WEST:
		return fmt.Sprintf("move_left(%d);", movement.Steps)
	case grid2dmap.SOUTH:
		return fmt.Sprintf("move_down(%d);", movement.Steps)
	case grid2dmap.EAST:
		return fmt.Sprintf("move_right(%d);", movement.Steps)
	default:
		panic("Looks like I accidently turned on diagonal pathfinding... mb\nif i see turn it off so I get valid moves")
	}

}

func (movement Movement) String() string {
	return movement.Generate()
}

type ObjectiveType byte

const (
	TERMINATION_OBJECTIVE = ObjectiveType(iota)
	NO_OBJECTIVE
	DROP_BALL
	EAT_BALL
	PAINT_TILE
	ROTATE
)

type ObjectiveGroup struct {
	movements   []Movement
	Coordinates []grid2dmap.Coordinate
	Objectives  []Objective
}

func NewObjectiveGroupFromSpec(spec ObjectiveSpec, underlying_map grid2dmap.Grid2DMap) ObjectiveGroup {
	path := astarpathfinding.Pathfind(&underlying_map, spec.Pathfinding_spec, []grid2dmap.Grid2DMapElement{})
	if !path.Successful() {
		panic(path)
	}

	return ObjectiveGroup{
		Coordinates: path.Coordinates,
		Objectives:  spec.Objectives,
	}
}

func NewObjectiveGroupsFromSpecs(specs []ObjectiveSpec, underlying_map grid2dmap.Grid2DMap) []ObjectiveGroup {
	var object_groups = make([]ObjectiveGroup, len(specs))

	for i := range specs {
		object_groups[i] = NewObjectiveGroupFromSpec(specs[i], underlying_map)
	}

	return object_groups
}

func ColourIDToPaintColour(id int) string {
	var colour_names = []string{"red", "blue", "green", "yellow", "cyan", "orange", "white", "black", "gray", "purple"}
	if id >= len(colour_names) {
		log.Println("Attempted to use invalid colour of ", id, " id")
		return "white"
	}
	return colour_names[id]
}

type Objective struct {
	Type     ObjectiveType
	Argument int // if the objective is DROP_BALL for example, this would be the amount of balls to drop
}

func (group *ObjectiveGroup) GenerateCode() string {
	var code string = ""

	for _, movement := range group.FetchMovements() {
		code += movement.String()
	}

	for _, objective := range group.Objectives {
		switch objective.Type {
		case TERMINATION_OBJECTIVE:
			code += "exit();"
		case DROP_BALL:
			code += fmt.Sprintf("spawn_ball(%d);", objective.Argument)
		case EAT_BALL:
			code += fmt.Sprintf("eat_ball(%d);", objective.Argument)
		case ROTATE:
			code += fmt.Sprintf("rotate(%d);", objective.Argument)
		case PAINT_TILE:
			code += fmt.Sprintf("paint('%s');", ColourIDToPaintColour(objective.Argument))
		case NO_OBJECTIVE:
			// nothing C:
		default:
			panic("Woah there, looks like an unknown objective is need?")
		}
	}

	return code
}

func (group *ObjectiveGroup) ExtractMovements() {
	if len(group.Coordinates) < 2 {
		return
	}

	group.movements = make([]Movement, len(group.Coordinates)-1)

	// Extract movements
	for i := range group.Coordinates[1:] {
		var cord_delta = grid2dmap.Coordinate{
			X: group.Coordinates[i+1].X - group.Coordinates[i].X,
			Y: group.Coordinates[i+1].Y - group.Coordinates[i].Y,
		}

		group.movements[i] = Movement{
			Direction: cord_delta.GetCompassDirectionVector(),
			Steps:     uint32(cord_delta.Hypot()),
		}

	}

	// Collapse movements
	for i := range group.movements[1:] {
		if group.movements[i+1].Direction != group.movements[i].Direction {
			continue
		}
		group.movements[i+1].Steps += group.movements[i].Steps
		group.movements[i].Steps = 0
	}

	group.movements = slices.DeleteFunc(group.movements, func(movement Movement) bool { return movement.Steps == 0 })
}

// Generates movements if not generated and returns them
func (group *ObjectiveGroup) FetchMovements() []Movement {
	if group.movements == nil {
		group.ExtractMovements()
	}

	return group.movements
}

type FoundInteractable struct {
	Coordinate grid2dmap.Coordinate
	Type       InteractableType

	Argument uint32
}

type CommonMapInfo struct {
	Karel_position grid2dmap.Coordinate
	Karel_tile     grid2dmap.Grid2DMapElement
}

// Map info for the start or end
type SnapshotMapInfo struct {
	CommonMapInfo
	Found_interactables []FoundInteractable
}

// Tile
// 0b00000000000000000000000000000000 Empty tile (u32)
// 0b00000000000000001111111111111111 Argument Mask (4 args, 1 nibble each)
// 0b00000000000000000000000000000111 Final Position Rotation Mask
// 0b00000000000000000000000001111000 Ball amount Mask
// 0b00000000000000000000011110000000 Paint color Mask
// 0b10000000000000000000000000000000 Start Position Flag
// 0b01000000000000000000000000000000 Final Position Flag
// 0b00100000000000000000000000000000 Wall Flag
// 0b00010000000000000000000000000000 Ball Flag
// 0b00001000000000000000000000000000 Paint Flag

// Updated tiles flags

type TileFlag grid2dmap.Grid2DMapElement

const (
	POSITION_FLAG = TileFlag(0x80000000 >> iota)
	WALL_FLAG
	BALL_FLAG
	PAINT_FLAG

	START_INTERACTABLE_FLAG = WALL_FLAG
	LAST_INTERACTABLE_FLAG  = PAINT_FLAG

	INTERACTABLE_FLAGS = BALL_FLAG | PAINT_FLAG
)

// Updated tile masks and offsets
const (
	ROTATION_ARG_MASK   = grid2dmap.Grid2DMapElement(0x00000007)
	ROTATION_ARG_OFFSET = grid2dmap.Grid2DMapElement(0)

	BALL_ARG_MASK   = grid2dmap.Grid2DMapElement(0x00000078)
	BALL_ARG_OFFSET = grid2dmap.Grid2DMapElement(3)

	PAINT_ARG_MASK   = grid2dmap.Grid2DMapElement(0x00000780)
	PAINT_ARG_OFFSET = grid2dmap.Grid2DMapElement(7)
)

type InteractableType uint32

const (
	NO_OBJECT = InteractableType(iota)
	BALL
	PAINT
)

func InteractableTypeToFlag(interactable_type InteractableType) TileFlag {
	switch interactable_type {
	case BALL:
		return BALL_FLAG
	case PAINT:
		return PAINT_FLAG
	}

	return 0
}

func GetInteractableTileTypes(tile grid2dmap.Grid2DMapElement) []InteractableType {
	var found_types = make([]InteractableType, 0)
	current_mask := START_INTERACTABLE_FLAG

	for i := 0; ; i++ {
		if (TileFlag(tile) & current_mask) != 0 {
			found_types = append(found_types, InteractableType(i))
		}
		if current_mask == LAST_INTERACTABLE_FLAG {
			break
		}
		current_mask >>= 1
	}
	return found_types
}

func GetTileArguments(tile grid2dmap.Grid2DMapElement, tile_flag TileFlag) uint32 {
	switch tile_flag {
	case BALL_FLAG:
		return uint32((tile & BALL_ARG_MASK) >> BALL_ARG_OFFSET)
	case START_INTERACTABLE_FLAG:
		fallthrough
	case POSITION_FLAG:
		return uint32((tile & ROTATION_ARG_MASK) >> ROTATION_ARG_OFFSET)
	case PAINT_FLAG:
		return uint32((tile & PAINT_ARG_MASK) >> PAINT_ARG_OFFSET)
	}
	return 0
}

func InteractableTileToFound(x, y int, tile grid2dmap.Grid2DMapElement) []FoundInteractable {
	var found_types = GetInteractableTileTypes(tile)
	var found_interactables []FoundInteractable = make([]FoundInteractable, len(found_types))
	for i := range found_interactables {
		found_interactables[i] = FoundInteractable{
			Coordinate: grid2dmap.Coordinate{X: x, Y: y},
			Type:       found_types[i],
			Argument:   GetTileArguments(tile, InteractableTypeToFlag(found_types[i])),
		}
	}

	return found_interactables
}

func GetCommonMapInfo(common_map grid2dmap.Grid2DMap) CommonMapInfo {
	var map_info = CommonMapInfo{}
	for x := range common_map.Width() {
		for y := range common_map.Height() {
			var checked_tile grid2dmap.Grid2DMapElement = common_map.Get(x, y)

			// Looking for jim_bob positioning
			if checked_tile&grid2dmap.Grid2DMapElement(POSITION_FLAG) != 0 {
				if map_info.Karel_tile != grid2dmap.EMPTY {
					panic("Woah there, there are multiple jim bob positions")
				}
				map_info.Karel_position = grid2dmap.Coordinate{X: x, Y: y}
				map_info.Karel_tile = checked_tile
			}
		}
	}
	return map_info
}

func GetSnapshotMapInfo(snapshot_map grid2dmap.Grid2DMap) SnapshotMapInfo {
	var map_info = SnapshotMapInfo{}
	map_info.CommonMapInfo = GetCommonMapInfo(snapshot_map)
	for x := range snapshot_map.Width() {
		for y := range snapshot_map.Height() {
			var checked_tile grid2dmap.Grid2DMapElement = snapshot_map.Get(x, y)

			if (checked_tile & grid2dmap.Grid2DMapElement(INTERACTABLE_FLAGS)) != 0 {
				interactables := InteractableTileToFound(x, y, checked_tile)

				for i := range interactables {
					map_info.Found_interactables = append(
						map_info.Found_interactables,
						interactables[i],
					)
				}
			}
		}
	}
	return map_info
}

type ObjectiveSpec struct {
	Pathfinding_spec astarpathfinding.PathfindingTarget
	Objectives       []Objective
}

func GetObjectiveSpecs(start, final grid2dmap.Grid2DMap) (core_map_geometry grid2dmap.Grid2DMap, objective_specs []ObjectiveSpec) {
	var start_map_info, final_map_info = GetSnapshotMapInfo(start), GetSnapshotMapInfo(final)

	core_map_geometry = start.Clone()

	for x := range core_map_geometry.Width() {
		for y := range core_map_geometry.Height() {
			if core_map_geometry.Get(x, y)&grid2dmap.Grid2DMapElement(WALL_FLAG) == 0 {
				core_map_geometry.Set(x, y, grid2dmap.EMPTY)
			}
		}
	}

	// Gets all the interactable coords to create objectives
	var interactable_coords = make(
		[]grid2dmap.Coordinate,
		0,
	)

	for _, interactable_1 := range start_map_info.Found_interactables {
		if !slices.Contains(interactable_coords, interactable_1.Coordinate) {
			interactable_coords = append(interactable_coords, interactable_1.Coordinate)
		}
	}

	for _, interactable_2 := range final_map_info.Found_interactables {
		if !slices.Contains(interactable_coords, interactable_2.Coordinate) {
			interactable_coords = append(interactable_coords, interactable_2.Coordinate)
		}
	}

	interactable_coords = grid2dmap.SortCoordinateListByClosest(interactable_coords, start_map_info.Karel_position)

	// Creates objective spec
	objective_specs = make(
		[]ObjectiveSpec,
		0,
	)

	// Adds pathfinding targets for specs
	var current_start_coord, current_target_coord grid2dmap.Coordinate
	current_start_coord = start_map_info.Karel_position

	for i := range interactable_coords {
		current_target_coord = interactable_coords[i]
		objective_specs = append(
			objective_specs,
			ObjectiveSpec{
				Pathfinding_spec: astarpathfinding.PathfindingTarget{
					Beginning: current_start_coord,
					Target:    current_target_coord,
				},
			},
		)

		current_start_coord = current_target_coord
	}

	current_target_coord = final_map_info.Karel_position

	objective_specs = append(
		objective_specs,
		ObjectiveSpec{
			Pathfinding_spec: astarpathfinding.PathfindingTarget{
				Beginning: current_start_coord,
				Target:    current_target_coord,
			},
			Objectives: []Objective{{Type: NO_OBJECTIVE}},
		},
	)

	// Adds objectives to objective specs

	for i, interactable_coord := range interactable_coords {
		objective_specs[i].Objectives = make([]Objective, 1)
		objective_specs[i].Objectives[0].Type = NO_OBJECTIVE

		// Gets all interactables part of the tile
		var interactables_start []FoundInteractable = InteractableTileToFound(
			interactable_coord.X, interactable_coord.Y,
			start.Get(interactable_coord.X, interactable_coord.Y),
		)
		var interactables_final []FoundInteractable = InteractableTileToFound(
			interactable_coord.X, interactable_coord.Y,
			final.Get(interactable_coord.X, interactable_coord.Y),
		)

		// Gets all the interactable types to ready for comparison
		var all_interactable_types = make([]InteractableType, 0)

		for _, interactable := range interactables_start {
			if !slices.Contains(all_interactable_types, interactable.Type) {
				all_interactable_types = append(all_interactable_types, interactable.Type)
			}
		}

		for _, interactable := range interactables_final {
			if !slices.Contains(all_interactable_types, interactable.Type) {
				all_interactable_types = append(all_interactable_types, interactable.Type)
			}
		}

		// Compares each interactable
		for _, interactable_type := range all_interactable_types {
			// Gets the two interactables to compare
			var interactable_1 = FoundInteractable{
				Coordinate: interactable_coord,
				Type:       interactable_type,
			}

			for _, potential_interactable := range interactables_start {
				if potential_interactable.Type != interactable_type {
					continue
				}
				interactable_1.Argument = potential_interactable.Argument
			}

			var interactable_2 = FoundInteractable{
				Coordinate: interactable_coord,
				Type:       interactable_type,
			}

			for _, potential_interactable := range interactables_final {
				if potential_interactable.Type != interactable_type {
					continue
				}
				interactable_2.Argument = potential_interactable.Argument
			}

			switch interactable_type {
			case BALL:
				amount_difference := int(interactable_2.Argument - interactable_1.Argument)

				if amount_difference < 0 {
					objective_specs[i].Objectives = append(
						objective_specs[i].Objectives,
						Objective{Type: EAT_BALL, Argument: utils.AbsInt(amount_difference)},
					)

				} else if amount_difference > 0 {
					objective_specs[i].Objectives = append(
						objective_specs[i].Objectives,
						Objective{Type: DROP_BALL, Argument: utils.AbsInt(amount_difference)},
					)
				}
			case PAINT:
				var selected_colour uint32 = interactable_2.Argument
				objective_specs[i].Objectives = append(
					objective_specs[i].Objectives,
					Objective{Type: PAINT_TILE, Argument: int(selected_colour)},
				)
			}

		}

	}

	var final_rotation int = int(
		(final.Get(
			final_map_info.Karel_position.X, final_map_info.Karel_position.Y,
		)&
			ROTATION_ARG_MASK)>>ROTATION_ARG_OFFSET) * 90
	if final_rotation != 0 {
		objective_specs[len(objective_specs)-1].Objectives = append(
			objective_specs[len(objective_specs)-1].Objectives, Objective{
				Type:     ROTATE,
				Argument: final_rotation,
			},
		)
	}

	return core_map_geometry, objective_specs
}
