package gamut

import (
	"image/color"
	"math"

	colorful "github.com/lucasb-eyer/go-colorful"
	"github.com/muesli/clusters"
	"github.com/muesli/kmeans"
)

// A ColorGenerator checks whether a point in the three dimensional CIELab space
// is suitable for color generation.
type ColorGenerator interface {
	Valid(col colorful.Color) bool
	Granularity() (l, c float64)
}

// BroadGranularity is used for wider color spaces, e.g. by the PastelGenerator
type BroadGranularity struct {
}

// FineGranularity is used for tighter color spaces, e.g. by the SimilarHueGenerator
type FineGranularity struct {
}

// SimilarHueGenerator produces colors with a similar hue as the given color
type SimilarHueGenerator struct {
	FineGranularity
	Color color.Color
}

// WarmGenerator produces "warm" colors
type WarmGenerator struct {
	BroadGranularity
}

// HappyGenerator produces "happy" colors
type HappyGenerator struct {
	BroadGranularity
}

// PastelGenerator produces "pastel" colors
type PastelGenerator struct {
	BroadGranularity
}

// Granularity returns BroadGranularity's default values
func (g BroadGranularity) Granularity() (l, c float64) {
	return 0.05, 0.1
}

// Granularity returns FineGranularity's default values
func (g FineGranularity) Granularity() (l, c float64) {
	return 0.01, 0.01
}

// distanceDegrees returns the distance between two angles on a circle
// e.g. the distance between 5 degrees and 355 degress is 10, not 350
func distanceDegrees(a1, a2 float64) float64 {
	mod := math.Mod(math.Abs(a1-a2), 360.0)
	if mod > 180.0 {
		return 360.0 - mod
	}

	return mod
}

// Valid returns true if the given color has a similar hue as the original color
func (gen SimilarHueGenerator) Valid(col colorful.Color) bool {
	cf, _ := colorful.MakeColor(gen.Color)
	h, c, l := cf.Hcl()
	hc, cc, lc := col.Hcl()

	if cc < c-0.35 || cc > c+0.35 {
		return false
	}
	if lc < l-0.6 || lc > l+0.6 {
		return false
	}
	if distanceDegrees(h, hc) > 7 {
		return false
	}
	if cf.DistanceCIE94(col) > 0.2 {
		return false
	}

	return true
}

// Valid returns true if the color is considered a "warm" color
func (cc WarmGenerator) Valid(col colorful.Color) bool {
	_, c, l := col.Hcl()
	return 0.1 <= c && c <= 0.4 && 0.2 <= l && l <= 0.5
}

// Valid returns true if the color is considered a "happy" color
func (cc HappyGenerator) Valid(col colorful.Color) bool {
	_, c, l := col.Hcl()
	return 0.3 <= c && 0.4 <= l && l <= 0.8
}

// Valid returns true if the color is considered a "pastel" color
func (cc PastelGenerator) Valid(col colorful.Color) bool {
	_, s, v := col.Hsv()
	return 0.2 <= s && s <= 0.4 && 0.7 <= v && v <= 1.0
}

// ColorObservation is a wrapper around colorful.Color, implementing the
// clusters.Observation interface
type ColorObservation struct {
	colorful.Color
}

// Coordinates returns the data points of a Lab color value
func (c ColorObservation) Coordinates() clusters.Coordinates {
	l, a, b := c.Lab()
	return clusters.Coordinates{l, a, b}
}

// Distance calculates the distance between two ColorObservations in the Lab
// color space
func (c ColorObservation) Distance(pos clusters.Coordinates) float64 {
	c2 := colorful.Lab(pos[0], pos[1], pos[2])
	return c.DistanceLab(c2)
}

// Generate returns a slice with the requested amount of colors, generated by
// the provided ColorGenerator.
func Generate(count int, generator ColorGenerator) ([]color.Color, error) {
	// Create data points in the CIE L*a*b color space
	// l for lightness channel
	// a, b for color channels
	var cc []color.Color
	dl, dab := generator.Granularity()

	var d clusters.Observations
	for l := 0.0; l <= 1.0; l += dl {
		for a := -1.0; a < 1.0; a += dab {
			for b := -1.0; b < 1.0; b += dab {
				col := colorful.Lab(l, a, b)
				// col = colorful.Hcl(a*360.0, b, c)

				if !col.IsValid() || !generator.Valid(col) {
					continue
				}

				d = append(d, ColorObservation{col})
			}
		}
	}

	// Enable graph generation (.png files) for each iteration
	// km, _ := kmeans.NewWithOptions(0.02, Plotter{})
	km, err := kmeans.NewWithOptions(0.02, nil)
	if err != nil {
		return cc, err
	}

	// Partition the color space into multiple clusters
	clusters, err := km.Partition(d, count)
	if err != nil {
		return cc, err
	}

	for _, c := range clusters {
		col := colorful.Lab(c.Center[0], c.Center[1], c.Center[2]).Clamped()
		cc = append(cc, col)
	}

	return cc, nil
}