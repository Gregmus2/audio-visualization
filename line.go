package main

import (
	"github.com/Gregmus2/simple-engine/graphics"
	"github.com/go-gl/gl/v4.6-core/gl"
)

type Line struct {
	x1, y1, x2, y2 float32
	prog  uint32
	shape *graphics.ShapeHelper
}

func (f *ObjectFactory) NewLine(x1, y1, x2, y2 float32) *Line {
	return &Line{
		x1: x1, y1: y1, x2: x2, y2: y2,
		prog:  f.Prog.GetByColor(graphics.White()),
		shape: f.Shape,
	}
}

func (u *Line) Draw(_ float32) error {
	gl.UseProgram(u.prog)
	u.shape.Line(u.x1, u.y1, u.x2, u.y2)
	gl.UseProgram(0)

	return nil
}
