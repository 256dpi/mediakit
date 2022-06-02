package mediakit

import "math"

// Size defines a size.
type Size struct {
	Width, Height float64
}

// Area returns the area.
func (s Size) Area() float64 {
	return s.Width * s.Height
}

// Aspect returns the aspect ratio.
func (s Size) Aspect() float64 {
	return s.Width / s.Height
}

// Scale returns a scaled size.
func (s Size) Scale(f float64) Size {
	return Size{
		Width:  s.Width * f,
		Height: s.Height * f,
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
func MaxWidth(max float64) Sizer {
	return func(s Size) Size {
		if s.Width <= max {
			return s
		}
		return Size{
			Width:  max,
			Height: max / s.Aspect(),
		}
	}
}

// MaxHeight limits the size by the specified height.
func MaxHeight(max float64) Sizer {
	return func(s Size) Size {
		if s.Height <= max {
			return s
		}
		return Size{
			Width:  max * s.Aspect(),
			Height: max,
		}
	}
}

// MaxArea limits the size by the specified area.
func MaxArea(max float64) Sizer {
	return func(s Size) Size {
		if s.Area() <= max {
			return s
		}
		f := math.Sqrt(max / s.Area())
		return s.Scale(f)
	}
}

// MaxSize limits the size by the specified size.
func MaxSize(max Size) Sizer {
	return func(s Size) Size {
		if s.Width <= max.Width && s.Height <= max.Height {
			return s
		}
		fw := max.Width / s.Width
		fh := max.Height / s.Height
		if fw < fh {
			return s.Scale(fw)
		}
		return s.Scale(fh)
	}
}
