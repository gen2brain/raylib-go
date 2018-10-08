// Package easings - Useful easing functions for values animation
//
// A port of Robert Penner's easing equations (http://robertpenner.com/easing/)
package easings

import (
	"math"
)

// Linear Easing functions

// LinearNone easing
// t: current time, b: begInnIng value, c: change In value, d: duration
func LinearNone(t, b, c, d float32) float32 {
	return c*t/d + b
}

// LinearIn easing
// t: current time, b: begInnIng value, c: change In value, d: duration
func LinearIn(t, b, c, d float32) float32 {
	return c*t/d + b
}

// LinearOut easing
// t: current time, b: begInnIng value, c: change In value, d: duration
func LinearOut(t, b, c, d float32) float32 {
	return c*t/d + b
}

// LinearInOut easing
// t: current time, b: begInnIng value, c: change In value, d: duration
func LinearInOut(t, b, c, d float32) float32 {
	return c*t/d + b
}

// Sine Easing functions

// SineIn easing
// t: current time, b: begInnIng value, c: change In value, d: duration
func SineIn(t, b, c, d float32) float32 {
	return -c*float32(math.Cos(float64(t/d)*(math.Pi/2))) + c + b
}

// SineOut easing
// t: current time, b: begInnIng value, c: change In value, d: duration
func SineOut(t, b, c, d float32) float32 {
	return c*float32(math.Sin(float64(t/d)*(math.Pi/2))) + b
}

// SineInOut easing
// t: current time, b: begInnIng value, c: change In value, d: duration
func SineInOut(t, b, c, d float32) float32 {
	return -c/2*(float32(math.Cos(math.Pi*float64(t/d)))-1) + b
}

// Circular Easing functions

// CircIn easing
// t: current time, b: begInnIng value, c: change In value, d: duration
func CircIn(t, b, c, d float32) float32 {
	t = t / d
	return -c*(float32(math.Sqrt(float64(1-t*t)))-1) + b
}

// CircOut easing
// t: current time, b: begInnIng value, c: change In value, d: duration
func CircOut(t, b, c, d float32) float32 {
	return c*float32(math.Sqrt(1-float64((t/d-1)*t))) + b
}

// CircInOut easing
// t: current time, b: begInnIng value, c: change In value, d: duration
func CircInOut(t, b, c, d float32) float32 {
	t = t / d * 2

	if t < 1 {
		return -c/2*(float32(math.Sqrt(float64(1-t*t)))-1) + b
	}

	t = t - 2
	return c/2*(float32(math.Sqrt(1-float64(t*t)))+1) + b
}

// Cubic Easing functions

// CubicIn easing
// t: current time, b: begInnIng value, c: change In value, d: duration
func CubicIn(t, b, c, d float32) float32 {
	t = t / d
	return c*t*t*t + b
}

// CubicOut easing
// t: current time, b: begInnIng value, c: change In value, d: duration
func CubicOut(t, b, c, d float32) float32 {
	t = t/d - 1
	return c*(t*t*t+1) + b
}

// CubicInOut easing
// t: current time, b: begInnIng value, c: change In value, d: duration
func CubicInOut(t, b, c, d float32) float32 {
	t = t / d * 2
	if t < 1 {
		return (c/2*t*t*t + b)
	}

	t = t - 2
	return c/2*(t*t*t+2) + b
}

// Quadratic Easing functions

// QuadIn easing
// t: current time, b: begInnIng value, c: change In value, d: duration
func QuadIn(t, b, c, d float32) float32 {
	t = t / d
	return c*t*t + b
}

// QuadOut easing
// t: current time, b: begInnIng value, c: change In value, d: duration
func QuadOut(t, b, c, d float32) float32 {
	t = t / d
	return (-c*t*(t-2) + b)
}

// QuadInOut easing
// t: current time, b: begInnIng value, c: change In value, d: duration
func QuadInOut(t, b, c, d float32) float32 {
	t = t / d * 2
	if t < 1 {
		return ((c / 2) * (t * t)) + b
	}

	return -c/2*((t-1)*(t-3)-1) + b
}

// Exponential Easing functions

// ExpoIn easing
// t: current time, b: begInnIng value, c: change In value, d: duration
func ExpoIn(t, b, c, d float32) float32 {
	if t == 0 {
		return b
	}

	return (c*float32(math.Pow(2, 10*float64(t/d-1))) + b)
}

// ExpoOut easing
// t: current time, b: begInnIng value, c: change In value, d: duration
func ExpoOut(t, b, c, d float32) float32 {
	if t == d {
		return (b + c)
	}

	return c*(-float32(math.Pow(2, -10*float64(t/d)))+1) + b
}

// ExpoInOut easing
// t: current time, b: begInnIng value, c: change In value, d: duration
func ExpoInOut(t, b, c, d float32) float32 {
	if t == 0 {
		return b
	}
	if t == d {
		return (b + c)
	}

	t = t / d * 2

	if t < 1 {
		return (c/2*float32(math.Pow(2, 10*float64(t-1))) + b)
	}

	t = t - 1
	return (c/2*(-float32(math.Pow(2, -10*float64(t)))+2) + b)
}

// Back Easing functions

// BackIn easing
// t: current time, b: begInnIng value, c: change In value, d: duration
func BackIn(t, b, c, d float32) float32 {
	s := float32(1.70158)
	t = t / d
	return c*t*t*((s+1)*t-s) + b
}

// BackOut easing
// t: current time, b: begInnIng value, c: change In value, d: duration
func BackOut(t, b, c, d float32) float32 {
	s := float32(1.70158)
	t = t/d - 1
	return c*(t*t*((s+1)*t+s)+1) + b
}

// BackInOut easing
// t: current time, b: begInnIng value, c: change In value, d: duration
func BackInOut(t, b, c, d float32) float32 {
	s := float32(1.70158)
	s = s * 1.525
	t = t / d * 2

	if t < 1 {
		return c/2*(t*t*((s+1)*t-s)) + b
	}

	t = t - 2
	return c/2*(t*t*((s+1)*t+s)+2) + b
}

// Bounce Easing functions

// BounceIn easing
// t: current time, b: begInnIng value, c: change In value, d: duration
func BounceIn(t, b, c, d float32) float32 {
	return (c - BounceOut(d-t, 0, c, d) + b)
}

// BounceOut easing
// t: current time, b: begInnIng value, c: change In value, d: duration
func BounceOut(t, b, c, d float32) float32 {
	t = t / d
	if t < (1 / 2.75) {
		return (c*(7.5625*t*t) + b)
	} else if t < (2 / 2.75) {
		t = t - (1.5 / 2.75)
		return c*(7.5625*t*t+0.75) + b
	} else if t < (2.5 / 2.75) {
		t = t - (2.25 / 2.75)
		return c*(7.5625*t*t+0.9375) + b
	}

	t = t - (2.625 / 2.75)
	return c*(7.5625*t*t+0.984375) + b
}

// BounceInOut easing
// t: current time, b: begInnIng value, c: change In value, d: duration
func BounceInOut(t, b, c, d float32) float32 {
	if t < d/2 {
		return BounceIn(t*2, 0, c, d)*0.5 + b
	}

	return BounceOut(t*2-d, 0, c, d)*0.5 + c*0.5 + b
}

// Elastic Easing functions

// ElasticIn easing
// t: current time, b: begInnIng value, c: change In value, d: duration
func ElasticIn(t, b, c, d float32) float32 {
	if t == 0 {
		return b
	}

	t = t / d

	if t == 1 {
		return b + c
	}

	p := d * 0.3
	a := c
	s := p / 4
	postFix := a * float32(math.Pow(2, 10*float64(t-1)))

	return -(postFix * float32(math.Sin(float64(t*d-s)*(2*math.Pi)/float64(p)))) + b
}

// ElasticOut easing
// t: current time, b: begInnIng value, c: change In value, d: duration
func ElasticOut(t, b, c, d float32) float32 {
	if t == 0 {
		return b
	}

	t = t / d

	if t == 1 {
		return b + c
	}

	p := d * 0.3
	a := c
	s := p / 4

	return a*float32(math.Pow(2, -10*float64(t)))*float32(math.Sin(float64(t*d-s)*(2*math.Pi)/float64(p))) + c + b
}

// ElasticInOut easing
// t: current time, b: begInnIng value, c: change In value, d: duration
func ElasticInOut(t, b, c, d float32) float32 {
	if t == 0 {
		return b
	}

	t = t / d * 2

	if t == 2 {
		return b + c
	}

	p := d * (0.3 * 1.5)
	a := c
	s := p / 4

	if t < 1 {
		t = t - 1
		postFix := a * float32(math.Pow(2, 10*float64(t)))
		return -0.5*(postFix*float32(math.Sin(float64(t*d-s)*(2*math.Pi)/float64(p)))) + b
	}

	t = t - 1
	postFix := a * float32(math.Pow(2, -10*(float64(t))))
	return postFix*float32(math.Sin(float64(t*d-s)*(2*math.Pi)/float64(p)))*0.5 + c + b
}
