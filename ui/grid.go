package ui

import (
	bmath "github.com/javanhut/BifrostEngine/m/v2/math"
)

type Grid struct {
	Size      float32
	Divisions int
	Color     [4]float32
	Visible   bool
	Position  bmath.Vector3
}

func NewGrid() *Grid {
	return &Grid{
		Size:      20.0,
		Divisions: 20,
		Color:     [4]float32{0.3, 0.3, 0.3, 0.5},
		Visible:   true,
		Position:  bmath.NewVector3(0, 0, 0),
	}
}

func (g *Grid) GetLines() []bmath.Vector3 {
	lines := make([]bmath.Vector3, 0)
	halfSize := g.Size / 2
	step := g.Size / float32(g.Divisions)
	
	// Generate lines along X axis
	for i := 0; i <= g.Divisions; i++ {
		z := -halfSize + float32(i)*step
		lines = append(lines, 
			bmath.NewVector3(-halfSize, g.Position.Y, z),
			bmath.NewVector3(halfSize, g.Position.Y, z),
		)
	}
	
	// Generate lines along Z axis
	for i := 0; i <= g.Divisions; i++ {
		x := -halfSize + float32(i)*step
		lines = append(lines,
			bmath.NewVector3(x, g.Position.Y, -halfSize),
			bmath.NewVector3(x, g.Position.Y, halfSize),
		)
	}
	
	return lines
}

func (g *Grid) SetSize(size float32) {
	if size > 0 {
		g.Size = size
	}
}

func (g *Grid) SetDivisions(divisions int) {
	if divisions > 0 && divisions <= 100 {
		g.Divisions = divisions
	}
}