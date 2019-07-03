package main

import (
	"image/color"

	"github.com/dianelooney/rave/waves"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

func main() {
	p, err := plot.New()
	if err != nil {
		panic(err)
	}
	p.Title.Text = "Waveform"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"

	// A quadratic function x^2
	f1 := plotter.NewFunction(waves.Experiment)
	f1.Color = color.RGBA{B: 255, A: 255}
	f1.Samples = 1000

	p.Add(f1)
	p.Legend.Add("f1", f1)
	p.Legend.ThumbnailWidth = 0.5 * vg.Inch

	// Set the axis ranges.  Unlike other data sets,
	// functions don't set the axis ranges automatically
	// since functions don't necessarily have a
	// finite range of x and y values.
	p.X.Min = 0
	p.X.Max = 1
	p.Y.Min = -1
	p.Y.Max = 1

	// Save the plot to a PNG file.
	if err := p.Save(6*vg.Inch, 4*vg.Inch, "charts/out.png"); err != nil {
		panic(err)
	}
}
