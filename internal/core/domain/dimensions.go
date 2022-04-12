package domain

type Dimensions struct {
	Width  int
	Height int
	Depth  int
}

func NewDimensions(width, height, depth int) Dimensions {
	return Dimensions{
		Width:  width,
		Height: height,
		Depth:  depth,
	}
}
