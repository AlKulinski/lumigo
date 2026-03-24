package domain

type Point struct {
	X int
	Y int
}

func NewPoint(x, y int) Point {
	return Point{X: x, Y: y}
}

func (p Point) DistanceTo(other Point) int {
	dx := p.X - other.X
	dy := p.Y - other.Y
	return dx*dx + dy*dy
}
