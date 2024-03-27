package softrender

import (
	"errors"
)

type Surface struct {
	Invisible uint8

	width, height int
	pixels        []uint8
}

func NewSurface(width, height int) (Surface, error) {
	if width < 0 || height < 0 {
		return Surface{},
			errors.New("the width and height must be positive integers")
	}
	return Surface{
			width:  width,
			height: height,
			pixels: make([]uint8, height*width),
		},
		nil
}

func (s *Surface) Width() int {
	return s.width
}

func (s *Surface) Height() int {
	return s.height
}

func (s *Surface) inRange(x, y int) bool {
	return 0 <= x && x < s.width && 0 <= y && y < s.height
}

func (s *Surface) GetPixel(x, y int) (uint8, error) {
	if s.inRange(x, y) {
		return s.pixels[y*s.width+x], nil
	}
	return s.Invisible, errors.New("coordinate not part of Surface")
}

func (s *Surface) SetPixel(color uint8, x, y int) {
	if s.inRange(x, y) {
		s.pixels[y*s.width+x] = color
	}
}

func (s *Surface) Line(color uint8, x0, y0, x1, y1 int) {
	abs := func(i int) int {
		if i < 0 {
			return -1
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

func (s *Surface) RectFill(color uint8, x, y, width, height int) {
	for xOff := range width {
		for yOff := range height {
			s.SetPixel(color, x+xOff, y+yOff)
		}
	}
}

func (s *Surface) RectOutline(color uint8, x, y, width, height int) {
	for xOff := range width {
		s.SetPixel(color, x+xOff, y)
		s.SetPixel(color, x+xOff, y+height-1)
	}
	for yOff := range height {
		s.SetPixel(color, x, y+yOff)
		s.SetPixel(color, x+width-1, y+yOff)
	}
}

func (s *Surface) circleLine(color uint8, centerX, centerY, x, y int) {
	s.Line(color, centerX+x, centerY+y, centerX+x, centerY-y)
	s.Line(color, centerX-x, centerY+y, centerX-x, centerY-y)
	s.Line(color, centerX+y, centerY+x, centerX+y, centerY-x)
	s.Line(color, centerX-y, centerY+x, centerX-y, centerY-x)
}

func (s *Surface) CircleFill(color uint8, x, y, radius int) {
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

func (s *Surface) circlePixel(color uint8, centerX, centerY, x, y int) {
	s.SetPixel(color, centerX+x, centerY+y)
	s.SetPixel(color, centerX-x, centerY+y)
	s.SetPixel(color, centerX+x, centerY-y)
	s.SetPixel(color, centerX-x, centerY-y)
	s.SetPixel(color, centerX+y, centerY+x)
	s.SetPixel(color, centerX-y, centerY+x)
	s.SetPixel(color, centerX+y, centerY-x)
	s.SetPixel(color, centerX-y, centerY-x)
}

func (s *Surface) CircleOutline(color uint8, x, y, radius int) {
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

func (s *Surface) Blit(x, y int, source *Surface) {
	for xOff := range source.width {
		for yOff := range source.height {
			color, _ := source.GetPixel(xOff, yOff)
			if color != source.Invisible {
				s.SetPixel(color, x+xOff, y+yOff)
			}
		}
	}
}

func (s *Surface) BlitPro(
	destX, destY, destWidth, destHeight int,
	source *Surface,
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

func (s *Surface) Fill(color uint8) {
	for i := range s.pixels {
		s.pixels[i] = color
	}
}
