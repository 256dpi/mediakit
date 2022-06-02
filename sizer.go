package mediakit

import "math"

// Size defines a size.
type Size struct {
	Width, Height int
}

// Area returns the area.
func (s Size) Area() int {
	return s.Width * s.Height
}

// Aspect returns the aspect ratio.
func (s Size) Aspect() float64 {
	return float64(s.Width) / float64(s.Height)
}

// Scale returns a scaled size.
func (s Size) Scale(f float64) Size {
	return Size{
		Width:  int(math.Round(float64(s.Width) * f)),
		Height: int(math.Round(float64(s.Height) * f)),
	}
}

// Sizer is a function that applies a function to a size.
type Sizer func(s Size) Size

// KeepSize returns the input size.
func KeepSize() Sizer {
	return func(s Size) Size {
		return s
	}
}

// MaxWidth limits the size by the specified width.
func MaxWidth(max int) Sizer {
	return func(s Size) Size {
		if s.Width <= max {
			return s
		}
		return Size{
			Width:  max,
			Height: int(math.Round(float64(max) / s.Aspect())),
		}
	}
}

// MaxHeight limits the size by the specified height.
func MaxHeight(max int) Sizer {
	return func(s Size) Size {
		if s.Height <= max {
			return s
		}
		return Size{
			Width:  int(math.Round(float64(max) * s.Aspect())),
			Height: max,
		}
	}
}

// MaxArea limits the size by the specified area.
func MaxArea(max int) Sizer {
	return func(s Size) Size {
		if s.Area() <= max {
			return s
		}
		f := math.Sqrt(float64(max) / float64(s.Area()))
		return s.Scale(f)
	}
}

// MaxSize limits the size by the specified size.
func MaxSize(max Size) Sizer {
	return func(s Size) Size {
		if s.Width <= max.Width && s.Height <= max.Height {
			return s
		}
		fw := float64(max.Width) / float64(s.Width)
		fh := float64(max.Height) / float64(s.Height)
		if fw < fh {
			return s.Scale(fw)
		}
		return s.Scale(fh)
	}
}
