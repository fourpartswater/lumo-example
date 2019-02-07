package lod

import (
	"math"
	"sort"
)

const (
	bytesPerComponent = 4
	pointStride       = 2
)

// LOD takes the input point array and sorts it by morton code. It then
// generates an offset array which match the byte offsets into the point buffer
// for each LOD. This is used at runtime to only render quadrants of the
// generated tile.
func LOD(data []float32, lod int) ([]float32, []int) {
	// get the points array sorted by morton code
	points := sortPoints(data)

	// generate codes for the sorted points
	codes := make([]int, len(points)/pointStride)
	for i := 0; i < len(points); i += pointStride {
		codes[i/2] = Morton(points[i], points[i+1])
	}

	// calc number of partitions and partition stride
	partitions := math.Pow(4, float64(lod))
	paritionStride := maxMorton / int(partitions)

	// set offsets
	offsets := make([]int, int(partitions))
	// init offsets as -1
	for i := range offsets {
		offsets[i] = -1
	}
	// set the offsets to the least byte in the array
	for i := len(codes) - 1; i >= 0; i-- {
		code := codes[i]
		j := code / paritionStride
		offsets[j] = i * (bytesPerComponent * pointStride)
	}
	// fill empty offsets up with next entries to ensure easy LOD
	for i := len(offsets) - 1; i >= 0; i-- {
		if offsets[i] == -1 {
			if i == len(offsets)-1 {
				offsets[i] = len(points) * bytesPerComponent
			} else {
				offsets[i] = offsets[i+1]
			}
		}
	}
	return points, offsets
}

// EncodeLOD generates the point LOD offsets and encodes them as a byte array.
func EncodeLOD(data []float32, lod int) []byte {
	// get sorted points and offsets
	points, offsets := LOD(data, lod)
	// encode the results
	return encodeLOD(points, offsets)
}

// encodeLOD generates the point LOD offsets and encodes them as a byte array.
func encodeLOD(points []float32, offsets []int) []byte {
	// encode data
	dataBytes := EncodeFloat32(points)
	// encode offsets
	offsetBytes := EncodeInt(offsets)
	// data length
	dataLength := make([]byte, 4)
	binary.LittleEndian.PutUint32(
		dataLength,
		uint32(len(dataBytes)))
	// offset length
	offsetLength := make([]byte, 4)
	binary.LittleEndian.PutUint32(
		offsetLength,
		uint32(len(offsetBytes)))
	// combine the buffers
	a := len(dataLength)
	b := len(offsetLength)
	c := len(dataBytes)
	d := len(offsetBytes)
	// create buffer
	bytes := make([]byte, a+b+c+d)
	// copy into buffer
	copy(bytes[0:a], dataLength)
	copy(bytes[a:a+b], offsetLength)
	copy(bytes[a+b:a+b+c], dataBytes)
	copy(bytes[a+b+c:a+b+c+d], offsetBytes)
	// return buffer
	return bytes
}

func sortPoints(data []float32) []float32 {
	points := make(pointArray, len(data)/pointStride)
	for i := 0; i < len(data); i += pointStride {
		x := data[i]
		y := data[i+1]
		points[i/pointStride] = [pointStride]float32{x, y}
	}
	// sort the points
	sort.Sort(points)
	// convert to flat array
	res := make([]float32, len(points)*pointStride)
	for i, point := range points {
		res[i*pointStride] = point[0]
		res[i*pointStride+1] = point[1]
	}
	return res
}

type pointArray [][pointStride]float32 // x, y

func (p pointArray) Len() int {
	return len(p)
}
func (p pointArray) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}
func (p pointArray) Less(i, j int) bool {
	return Morton(p[i][0], p[i][1]) < Morton(p[j][0], p[j][1])
}
