package astarpathfinding

import (
	"math"
	"slices"
	grid2dmap "sumac-gen/Grid2DMap"
	utils "sumac-gen/Utils"
)

var (
	Can_move_diagonally       bool   = false
	Max_pathfind_steps        uint32 = 1000
	Cost_fixed_point_accuracy uint64 = 100
)

func CompassDirectionCost(dir grid2dmap.CompassDirectionVector) uint64 {
	if dir <= grid2dmap.EAST {
		return 1 * Cost_fixed_point_accuracy
	} else {
		return uint64(math.Sqrt(2) * float64(Cost_fixed_point_accuracy))
	}
}

type Path struct {
	Coordinates []grid2dmap.Coordinate
	successful  bool
}

func (path *Path) Successful() bool {
	return path.successful
}

type Node struct {
	parent_node *Node

	x, y int

	guaranteed_cost, estimated_cost uint64

	child_nodes    []*Node
	checked_before bool
}

func (node *Node) X() int {
	return node.x
}

func (node *Node) Y() int {
	return node.y
}

func (node *Node) Parent() *Node {
	return node.parent_node
}
func (node *Node) TotalCost() uint64 {
	return node.guaranteed_cost + node.estimated_cost
}

type PathfindingTarget struct {
	Beginning grid2dmap.Coordinate
	Target    grid2dmap.Coordinate
}

func SetPathInGrid(grid *grid2dmap.Grid2DMap, path Path, value grid2dmap.Grid2DMapElement) {
	if !path.Successful() {
		return

	}

	for _, coord := range path.Coordinates {
		grid.Set(coord.X, coord.Y, value)
	}
}

func createNode(parent *Node, step_direction grid2dmap.CompassDirectionVector, target_x, target_y int) *Node {
	var node Node
	var parent_copy Node

	// Resets / Blocks Diagonal steps if disallowed
	if !Can_move_diagonally && step_direction > grid2dmap.EAST {
		step_direction %= grid2dmap.EAST
	}

	if parent != nil {
		parent_copy = *parent
	}

	node.parent_node = parent

	node.guaranteed_cost = parent_copy.guaranteed_cost + CompassDirectionCost(step_direction)

	node.x, node.y = parent_copy.x+step_direction.StepInX(), parent_copy.y+step_direction.StepInY()

	node.estimated_cost = uint64(
		utils.AbsInt(node.x-target_x)+
			utils.AbsInt(node.y-target_y),
	) * Cost_fixed_point_accuracy

	node.parent_node.child_nodes = append(node.parent_node.child_nodes, &node)
	return node.parent_node.child_nodes[len(node.parent_node.child_nodes)-1]
}

func Pathfind(grid *grid2dmap.Grid2DMap, target PathfindingTarget, ignored_elements []grid2dmap.Grid2DMapElement) Path {
	var raw_nodes map[int]*Node = make(map[int]*Node)

	var open_nodes = []int{}

	var current_node = new(Node)
	var current_node_index int
	var path Path

	current_node.x = target.Beginning.X
	current_node.y = target.Beginning.Y

	current_node.estimated_cost = uint64(utils.AbsInt(target.Beginning.X-target.Target.X) + utils.AbsInt(target.Beginning.Y-target.Target.Y))
	current_node_index = grid.PositionTo1DIndex(current_node.x, current_node.y)
	open_nodes = append(open_nodes, grid.PositionTo1DIndex(current_node.x, current_node.y))
	raw_nodes[current_node_index] = current_node

	for current_node.estimated_cost != 0 {
		var successful_directions int = 0

		slices.SortFunc(
			open_nodes,
			func(a, b int) int {
				switch {
				case raw_nodes[a].TotalCost() < raw_nodes[b].TotalCost():
					return -1
				case raw_nodes[a].TotalCost() == raw_nodes[b].TotalCost():
					if raw_nodes[a].guaranteed_cost <= raw_nodes[b].guaranteed_cost {
						return -1
					}
					return 1
				case raw_nodes[a].TotalCost() > raw_nodes[b].TotalCost():
					return 1
				}
				return 0
			},
		)
		if len(open_nodes) == 0 {
			break
		}
		current_node = raw_nodes[open_nodes[0]]
		current_node_index = open_nodes[0]

		open_nodes = slices.DeleteFunc(open_nodes, func(node int) bool { return node == current_node_index })

		for dir := range grid2dmap.CompassDirectionVector(4) {
			attempted_x, attempted_y := current_node.x+dir.StepInX(), current_node.y+dir.StepInY()

			if current_node.parent_node != nil &&
				(current_node.parent_node.x == attempted_x && current_node.parent_node.y == attempted_y) {
				continue
			}

			if attempted_x == target.Target.X && attempted_y == target.Target.Y {
				raw_nodes[grid.PositionTo1DIndex(attempted_x, attempted_y)] = createNode(current_node, dir, target.Target.X, target.Target.Y)
				open_nodes = append(open_nodes, grid.PositionTo1DIndex(attempted_x, attempted_y))
				current_node = raw_nodes[grid.PositionTo1DIndex(attempted_x, attempted_y)]
				break
			}

			if grid.Get(attempted_x, attempted_y) != grid2dmap.EMPTY && !slices.Contains(ignored_elements, grid.Get(attempted_x, attempted_y)) {
				continue
			}

			successful_directions++

			if _, ok := raw_nodes[grid.PositionTo1DIndex(attempted_x, attempted_y)]; ok {
				continue
			}

			raw_nodes[grid.PositionTo1DIndex(attempted_x, attempted_y)] = createNode(current_node, dir, target.Target.X, target.Target.Y)
			open_nodes = append(open_nodes, grid.PositionTo1DIndex(attempted_x, attempted_y))

		}
	}

	for current_node != nil {
		path.Coordinates = slices.Insert(path.Coordinates, 0, grid2dmap.Coordinate{X: current_node.x, Y: current_node.y})
		current_node = current_node.parent_node
	}

	path.Coordinates = slices.Insert(path.Coordinates, 0, grid2dmap.Coordinate{X: target.Beginning.X, Y: target.Beginning.Y})

	path.successful = true

	if coord := path.Coordinates[len(path.Coordinates)-1]; coord.X != target.Target.X || coord.Y != target.Target.Y {
		path.successful = false
		path.Coordinates = []grid2dmap.Coordinate{}
	}
	return path
}
