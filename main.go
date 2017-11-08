package main

import (
	"fmt"
	"image/color"
	"math/rand"
	"time"

	"github.com/bugra/kmeans"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

const (
	K = 7
)

func main() {
	rand.Seed(int64(time.Now().Unix()))
	n := 1000
	bubbleData := randomTriples(n)

	p, err := plot.New()
	if err != nil {
		panic(err)
	}
	p.Title.Text = "Bubbles"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"

	bs, err := plotter.NewBubbles(bubbleData, vg.Points(1), vg.Points(20))
	if err != nil {
		panic(err)
	}
	bs.Color = color.RGBA{R: 196, B: 128, A: 255}
	bs.MaxRadius = 2
	bs.MinRadius = 1
	p.Add(bs)

	if err := p.Save(4*vg.Inch, 4*vg.Inch, "bubble.png"); err != nil {
		panic(err)
	}

	data := make([][]float64, bs.Len())
	for i := 0; i < bs.Len(); i++ {
		x, y, _ := bs.XYZ(i)
		data[i] = []float64{x, y}
	}

	distanceF := map[string]kmeans.DistanceFunction{
		"SquaredEuclideanDistance": kmeans.SquaredEuclideanDistance,
		"ManhattanDistance":        kmeans.ManhattanDistance,
		"EuclideanDistance":        kmeans.EuclideanDistance,
		"ChebyshevDistance":        kmeans.ChebyshevDistance,
		"HammingDistance":          kmeans.HammingDistance,
		"BrayCurtisDistance":       kmeans.BrayCurtisDistance,
		"CanberraDistance":         kmeans.CanberraDistance,
	}

	for name, fct := range distanceF {

		clusters, err := kmeans.Kmeans(data, K, fct, 10)
		if err != nil {
			fmt.Printf("Error:%v", err)
			return
		}
		drawCluster(clusters, data, name)
	}
}

func drawCluster(clusters []int, data [][]float64, name string) {
	p, err := plot.New()
	if err != nil {
		panic(err)
	}
	p.Title.Text = "Bubbles Cluster"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"

	BS := map[int]*plotter.Bubbles{}

	for index, C := range clusters {

		bs, ok := BS[C]
		if !ok {
			var err2 error
			bs, err2 = plotter.NewBubbles(make(plotter.XYZs, 100), vg.Points(1), vg.Points(20))
			if err2 != nil {
				panic(err2)
			}
			bs.MaxRadius = 5
			bs.MinRadius = 5
			bs.Color = color.RGBA{R: 250 - uint8(250/K*C), B: uint8(250 / K * C), A: 255}
			BS[C] = bs
		}

		bs.XYZs = append(bs.XYZs, struct{ X, Y, Z float64 }{X: data[index][0], Y: data[index][1], Z: 1})
	}

	for _, bs := range BS {
		p.Add(bs)
	}
	if err := p.Save(4*vg.Inch, 4*vg.Inch, "bubble_"+name+".png"); err != nil {
		panic(err)
	}
}

// randomTriples returns some random x, y, z triples
// with some interesting kind of trend.
func randomTriples(n int) plotter.XYZs {
	data := make(plotter.XYZs, n)
	for i := range data {
		//if i == 0 {
		data[i].X = rand.Float64()
		//} else {
		//	data[i].X = data[i-1].X + 2*rand.Float64()
		//}
		data[i].Y = rand.Float64() //data[i].X + 100*rand.Float64()
		data[i].Z = 1              //data[i].X
	}
	return data
}
