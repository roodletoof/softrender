package softrender

import (
	"errors"
)

type Surface[T comparable] struct {
	Invisible T

	width, height int
	pixels        []T
}

func NewSurface[T comparable](width, height int) (Surface[T], error) {
	if width < 0 || height < 0 {
		return Surface[T]{},
			errors.New("the width and height must be positive integers")
	}
	return Surface[T]{
			width:  width,
			height: height,
			pixels: make([]T, height*width),
		},
		nil
}

func (s *Surface[T]) Width() int {
	return s.width
}

func (s *Surface[T]) Height() int {
	return s.height
}

func (s *Surface[T]) inRange(x, y int) bool {
	return 0 <= x && x < s.width && 0 <= y && y < s.height
}

func (s *Surface[T]) GetPixel(x, y int) (T, error) {
	if s.inRange(x, y) {
		return s.pixels[y*s.width+x], nil
	}
	return s.Invisible, errors.New("coordinate not part of Surface")
}

func (s *Surface[T]) SetPixel(color T, x, y int) {
	if s.inRange(x, y) {
		s.pixels[y*s.width+x] = color
	}
}

func (s *Surface[T]) Line(color T, x0, y0, x1, y1 int) {
	abs := func(i int) int {
		if i < 0 {
			return -i
		}
		return i
	}

	dx := abs(x1 - x0)
	sx := -1
	if x0 < x1 {
		sx = 1
	}
	dy := -abs(y1 - y0)
	sy := -1
	if y0 < y1 {
		sy = 1
	}
	errorI := dx + dy

	for {
		s.SetPixel(color, x0, y0)
		if x0 == x1 && y0 == y1 {
			break
		}
		e2 := 2 * errorI
		if e2 >= dy {
			if x0 == x1 {
				break
			}
			errorI += dy
			x0 += sx
		}
		if e2 <= dx {
			if y0 == y1 {
				break
			}
			errorI += dx
			y0 += sy
		}
	}
}

func (s *Surface[T]) RectFill(color T, x, y, width, height int) {
	for xOff := range width {
		for yOff := range height {
			s.SetPixel(color, x+xOff, y+yOff)
		}
	}
}

func (s *Surface[T]) RectOutline(color T, x, y, width, height int) {
	for xOff := range width {
		s.SetPixel(color, x+xOff, y)
		s.SetPixel(color, x+xOff, y+height-1)
	}
	for yOff := range height {
		s.SetPixel(color, x, y+yOff)
		s.SetPixel(color, x+width-1, y+yOff)
	}
}

func (s *Surface[T]) circleLine(color T, centerX, centerY, x, y int) {
	s.Line(color, centerX+x, centerY+y, centerX+x, centerY-y)
	s.Line(color, centerX-x, centerY+y, centerX-x, centerY-y)
	s.Line(color, centerX+y, centerY+x, centerX+y, centerY-x)
	s.Line(color, centerX-y, centerY+x, centerX-y, centerY-x)
}

func (s *Surface[T]) CircleFill(color T, x, y, radius int) {
	outerX := 0
	outerY := radius
	d := 3 - 2*radius
	s.circleLine(color, x, y, outerX, outerY)
	for outerX <= outerY {
		outerX += 1
		if 0 < d {
			outerY -= 1
			d = d + 4*(outerX-outerY) + 10
		} else {
			d = d + 4*outerX + 6
		}
		s.circleLine(color, x, y, outerX, outerY)
	}
}

func (s *Surface[T]) circlePixel(color T, centerX, centerY, x, y int) {
	s.SetPixel(color, centerX+x, centerY+y)
	s.SetPixel(color, centerX-x, centerY+y)
	s.SetPixel(color, centerX+x, centerY-y)
	s.SetPixel(color, centerX-x, centerY-y)
	s.SetPixel(color, centerX+y, centerY+x)
	s.SetPixel(color, centerX-y, centerY+x)
	s.SetPixel(color, centerX+y, centerY-x)
	s.SetPixel(color, centerX-y, centerY-x)
}

func (s *Surface[T]) CircleOutline(color T, x, y, radius int) {
	outerX := 0
	outerY := radius
	d := 3 - 2*radius
	s.circlePixel(color, x, y, outerX, outerY)
	for outerX <= outerY {
		outerX += 1
		if 0 < d {
			outerY -= 1
			d = d + 4*(outerX-outerY) + 10
		} else {
			d = d + 4*outerX + 6
		}
		s.circlePixel(color, x, y, outerX, outerY)
	}
}

func (s *Surface[T]) Blit(x, y int, source *Surface[T]) {
	for xOff := range source.width {
		for yOff := range source.height {
			color, _ := source.GetPixel(xOff, yOff)
			if color != source.Invisible {
				s.SetPixel(color, x+xOff, y+yOff)
			}
		}
	}
}

func (s *Surface[T]) BlitPro(
	destX, destY, destWidth, destHeight int,
	source *Surface[T],
	sourceX, sourceY, sourceWidth, sourceHeight int,
) {
	for destXOff := range destWidth {
		for destYOff := range destHeight {

			currSourceX := sourceX + sourceWidth*destXOff/destWidth
			currSourceY := sourceY + sourceHeight*destYOff/destHeight

			color, _ := source.GetPixel(currSourceX, currSourceY)

			if color != source.Invisible {
				s.SetPixel(color, destX+destXOff, destY+destYOff)
			}
		}
	}
}

func (s *Surface[T]) Fill(color T) {
	for i := range s.pixels {
		s.pixels[i] = color
	}
}
