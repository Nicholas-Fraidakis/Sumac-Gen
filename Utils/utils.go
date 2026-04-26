package utils

type Tuple[A any, B any] struct {
	A A
	B B
}

func AbsInt(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

/*

func FindRobloxPID() (int, error) {
	ids, err := robotgo.FindIds("Roblox")

	var roblox_pid int = 0
	if err != nil {
		panic(err)
	}

	for _, id := range ids {
		if robotgo.GetTitle(id) == "Roblox" {
			roblox_pid = id
		}
	}

	if roblox_pid == 0 {
		return 0, errors.New("couldn't find roblox's PID")
	}

	return roblox_pid, nil
}

func GetRobloxDimensions() (int, int, int, int) {
	roblox_pid, err := FindRobloxPID()

	if err != nil {
		return 0, 0, 0, 0
	}

	x, y, w, h := robotgo.GetBounds(roblox_pid)

	if runtime.GOOS == "windows" {
		x += 11
		y += 45
		w -= 22
		h -= 56
	}

	return x, y, w, h
}

func getGeneratorDiamensionsGeneralArea() (int, int, int, int) {
	x, y, w, h := GetRobloxDimensions()

	scale := float64(w) / 1940.0
	if float64(w)/float64(h) >= 1940.0/901.0 {
		scale = float64(h) / 901.0
	}
	x, y = x+w/2-int(550.0*scale)/2, y+h/2-int(700.0*scale)/2

	w, h = int(550.0*scale), int(700.0*scale)

	return x, y, w, h
}

func GetGeneratorDiamensions() (int, int, int, int) {
	var new_x, new_y, new_w, new_h int
	x, y, w, h := getGeneratorDiamensionsGeneralArea()
	photo, err := robotgo.Capture(x, y, w, h)

	if err != nil {
		return 0, 0, 0, 0
	}

	new_x, new_y, new_w, new_h = x, y, w, h
	for internal_y := range h {
		if photo.At(w/2, internal_y) != GEN_BORDER_COLOUR {
			continue
		}

		new_y += internal_y
		new_h -= internal_y
		break
	}

	for internal_y := range h {
		if photo.At(w/2, h-internal_y) != GEN_BORDER_COLOUR {
			continue
		}

		new_h -= internal_y
		break
	}

	for internal_x := range w {
		if photo.At(internal_x, h/2) != GEN_BORDER_COLOUR {
			continue
		}

		new_x += internal_x
		new_w -= internal_x
		break
	}

	for internal_x := range w {
		if photo.At(w-internal_x, h/2) != GEN_BORDER_COLOUR {
			continue
		}

		new_w -= internal_x
		break
	}
	return new_x, new_y, new_w, new_h
}

func CaptureRoblox() (*image.RGBA, error) {
	_, err := FindRobloxPID()

	if err != nil {
		return nil, err
	}

	x, y, w, h := GetRobloxDimensions()

	return robotgo.Capture(x, y, w, h)
}

var (
	GEN_BORDER_COLOUR         color.RGBA = color.RGBA{0x14, 0x14, 0x14, 255}
	GEN_BACKGROUND_COLOUR     color.RGBA = color.RGBA{0x0a, 0x0a, 0xa, 255}
	GEN_BACKGROUND_ALT_COLOUR color.RGBA = color.RGBA{0x0e, 0x0e, 0xe, 255}
)

func CaptureGenerator() (*image.RGBA, error) {
	return robotgo.Capture(GetGeneratorDiamensions())
}

func Clamp(n, min, max int) int {
	if n < min {
		n = min
	} else if n > max {
		n = max
	}
	return n
}

func GeneratorPointToPosition(x, y, offset_x, offset_y, image_length int) (int, int) {
	var border_size = image_length / 360
	var tile_size = int(float64(image_length-border_size*8) / 6)
	var px, py = int(
		tile_size*x + tile_size/2 + border_size*Clamp(x-1, 0, 4),
	), int(
		tile_size*y + tile_size/3 + border_size*Clamp(y-1, 0, 4),
	)
	return px, py
}
func accessPointInGenerator(x, y, image_length int, generator_image *image.RGBA) color.RGBA {
	var border_size = image_length / 360
	var tile_size = int(float64(image_length-border_size*8) / 6)
	var px, py = int(
		tile_size*x + tile_size/2 + border_size*Clamp(x-1, 0, 4),
	), int(
		tile_size*y + tile_size/3 + border_size*Clamp(y-1, 0, 4),
	)

	return generator_image.At(
		px, py,
	).(color.RGBA)
}

func ParseGenerator() grid2dmap.Grid2DMap {
	var generator_image *image.RGBA
	var image_length int
	var puzzle = grid2dmap.New(6, 6)
	var generator_colours = []color.RGBA{}

	generator_image, err := CaptureGenerator()
	if err != nil {
		fmt.Println("AAAAAA")
		return puzzle
	}

	image_length = generator_image.Bounds().Size().X

	for y := range 6 {
		for x := range 6 {
			var node_color = accessPointInGenerator(x, y, image_length, generator_image)

			if node_color == GEN_BACKGROUND_COLOUR || node_color == GEN_BACKGROUND_ALT_COLOUR {
				continue
			}

			if slices.Contains(generator_colours, node_color) {
				for n, colour := range generator_colours {
					if colour == node_color {
						puzzle.Set(x, y, grid2dmap.Grid2DMapElement(n+1))
					}
				}
				continue
			}

			generator_colours = append(generator_colours, node_color)

			puzzle.Set(x, y, grid2dmap.Grid2DMapElement(len(generator_colours)))
		}
	}
	return puzzle
}
*/
