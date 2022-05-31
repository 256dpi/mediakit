package mediakit

import "math"

// Size defines a size.
type Size struct {
	W, H float64
}

// Area returns the area.
func (s Size) Area() float64 {
	return s.W * s.H
}

// Aspect returns the aspect ratio.
func (s Size) Aspect() float64 {
	return s.W / s.H
}

// Scale returns a scaled size.
func (s Size) Scale(f float64) Size {
	return Size{
		W: s.W * f,
		H: s.H * f,
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
		if s.W <= max {
			return s
		}
		return Size{
			W: max,
			H: max / s.Aspect(),
		}
	}
}

// MaxHeight limits the size by the specified height.
func MaxHeight(max float64) Sizer {
	return func(s Size) Size {
		if s.H <= max {
			return s
		}
		return Size{
			W: max * s.Aspect(),
			H: max,
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
		if s.W <= max.W && s.H <= max.H {
			return s
		}
		fw := max.W / s.W
		fh := max.H / s.H
		if fw < fh {
			return s.Scale(fw)
		}
		return s.Scale(fh)
	}
}
