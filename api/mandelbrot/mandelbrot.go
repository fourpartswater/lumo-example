package mandelbrot

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
)

const (
	min           = -2.2
	max           = 2.2
	resolution    = 256
	escapeModulos = 2.3
	maxIterations = 64
)

type complex struct {
	r float64
	i float64
}

func (c complex) Add(o complex) complex {
	return complex{
		r: c.r + o.r,
		i: c.i + o.i,
	}
}

func (c complex) Mul(o complex) complex {
	return complex{
		r: (c.r * o.r) - (c.i * o.i),
		i: (c.r * o.i) + (c.i * o.r),
	}
}

func (c complex) Abs() float64 {
	return math.Sqrt((c.r * c.r) + (c.i * c.i))
}

func writeEscapeColor(buffer []uint8, x int, y int, numIterations int) {
	div := float64(numIterations) / float64(maxIterations)
	gray := math.Max(1.0-div, 0.1)
	buffer[(x*resolution+y)*4] = uint8(255.0 * gray)
	buffer[(x*resolution+y)*4+1] = uint8(255.0 * gray)
	buffer[(x*resolution+y)*4+2] = uint8(255.0 * gray)
	buffer[(x*resolution+y)*4+3] = uint8(255.0)
}

func writeColor(buffer []uint8, x int, y int, modulus float64) {
	factor := (modulus / escapeModulos)

	incr := math.Log10(factor * 3.5)
	r := math.Min(math.Abs(8.0*incr)*factor, 1.0)
	g := math.Min(math.Abs(4.0*incr)*factor, 1.0)
	b := math.Min(math.Abs(incr)*factor, 1.0)

	buffer[(x*resolution+y)*4] = uint8(255.0 * r)
	buffer[(x*resolution+y)*4+1] = uint8(255.0 * g)
	buffer[(x*resolution+y)*4+2] = uint8(255.0 * b)
	buffer[(x*resolution+y)*4+3] = uint8(255.0)
}

func mandelbrot(z int, x int, y int) []uint8 {

	dim := math.Pow(2, float64(z))
	scale := (max - min) / dim
	pixelScale := scale / float64(resolution)
	tileXMin := min + (float64(x) * scale)
	tileYMin := min + (float64(y) * scale)

	buffer := make([]uint8, resolution*resolution*4)

	// Iterate through the entire panel, pixel by pixel
	for row := 0; row < resolution; row++ {
		// Calculate the actual y position
		yPos := tileYMin + float64(row)*pixelScale
		for col := 0; col < resolution; col++ {
			// Calculate the actual x position
			xPos := tileXMin + float64(col)*pixelScale
			// Create the complex number for this position
			c := complex{r: xPos, i: yPos}
			z := complex{r: 0, i: 0}
			iterations := 0
			escaped := false
			modulus := 0.0
			// Iterate the fractal equation z = z*z + c until z either
			// escapes or the maximum number of iterations is reached
			for {
				z = z.Mul(z).Add(c)
				modulus = z.Abs()
				escaped = modulus > escapeModulos
				iterations++
				if escaped || iterations > maxIterations {
					break
				}
			}
			// Set the colour according to what stopped the above loop
			px := (resolution - 1) - row
			py := col
			if escaped {
				writeEscapeColor(buffer, px, py, iterations)
			} else {
				writeColor(buffer, px, py, modulus)
			}
		}
	}
	return buffer
}

func handleErr(w http.ResponseWriter, err error) {
	// write error header
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(500)
	// error string
	bytes, err := json.Marshal(map[string]interface{}{
		"success": false,
		"error":   err.Error(),
	})
	if err == nil {
		// write error
		fmt.Fprint(w, string(bytes))
	}
}

// Handler generates a mandelbrot arraybuffer tile.
func Handler(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(r.URL.Path, "/")
	if len(path) < 4 {
		handleErr(w, fmt.Errorf("missing params"))
		return
	}
	z, err := strconv.ParseUint(path[2], 10, 64)
	if err != nil {
		handleErr(w, err)
		return
	}
	x, err := strconv.ParseUint(path[3], 10, 64)
	if err != nil {
		handleErr(w, err)
		return
	}
	y, err := strconv.ParseUint(path[4], 10, 64)
	if err != nil {
		handleErr(w, err)
		return
	}
	if x >= (1 << z) {
		handleErr(w, fmt.Errorf("x parameter out of range"))
		return
	}
	if y >= (1 << z) {
		handleErr(w, fmt.Errorf("y parameter out of range"))
		return
	}
	buffer := mandelbrot(int(z), int(x), int(y))
	// set content type response header
	w.Header().Set("Content-Type", "application/octet-stream")
	w.WriteHeader(200)
	w.Write(buffer)
}
