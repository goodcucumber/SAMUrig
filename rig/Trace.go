package rig

import (
	"math"
)

const MAXSTEPS = 50000

var C0 float64
var C1 float64
var C2 float64
var C3 float64

//Trace
//start {x, y, z, vx, vy, vz} in mm and mm/s
//dist distance from SAMURAI center to FDC2(K) plane in mm
//angle angle fo FDC2(K) in degree
//h step length in s (if m/q = 1)
//bmap magnetic filed map
func Trace(start [6]float64, dist float64, angle float64, h float64, bmap []float64) ([4]float64, [6]float64) {
	var x1 [6]float64
	x1 = start
	var x2 [6]float64
	cos := math.Cos(angle * math.Pi / 180)
	sin := math.Sin(angle * math.Pi / 180)
	tot := 0.0
	for i := 0; i < MAXSTEPS; i++ {
		x2 = next(x1, bmap, h)
		if x2[0]*sin+x2[2]*cos >= dist {
			d1 := x1[0]*sin + x1[2]*cos
			d2 := x2[0]*sin + x2[2]*cos
			k := (dist - d1) / (d2 - d1)
			for i := 0; i < 6; i++ {
				x2[i] = x1[i] + k*(x2[i]-x1[i])
			}
			tot += math.Sqrt((x2[0]-x1[0])*(x2[0]-x1[0]) + (x2[1]-x1[1])*(x2[1]-x1[1]) + (x2[2]-x1[2])*(x2[2]-x1[2]))
			break
		}
		tot += math.Sqrt((x2[0]-x1[0])*(x2[0]-x1[0]) + (x2[1]-x1[1])*(x2[1]-x1[1]) + (x2[2]-x1[2])*(x2[2]-x1[2]))
		x1 = x2
	}
	var rst [4]float64
	rst[0] = x2[0]*cos - x2[2]*sin
	rst[1] = x2[1]
	rst[2] = x2[0]*sin + x2[2]*cos
	rst[3] = tot
	return rst, x2
}

//TraceWithPath
//start {x, y, z, vx, vy, vz} in mm and mm/s
//dist distance from SAMURAI center to FDC2(K) plane in mm
//angle angle fo FDC2(K) in degree
//h step length in s (if m/q = 1)
//bmap magnetic filed map
func TraceWithPath(start [6]float64, dist float64, angle float64, h float64, bmap []float64) ([4]float64, [][6]float64) {
	path := make([][6]float64, MAXSTEPS+1)
	pathidx := 1
	path[0] = start
	var x1 [6]float64
	x1 = start
	var x2 [6]float64
	cos := math.Cos(angle * math.Pi / 180)
	sin := math.Sin(angle * math.Pi / 180)
	tot := 0.0
	for i := 0; i < MAXSTEPS; i++ {
		x2 = next(x1, bmap, h)
		if x2[0]*sin+x2[2]*cos >= dist {
			d1 := x1[0]*sin + x1[2]*cos
			d2 := x2[0]*sin + x2[2]*cos
			k := (dist - d1) / (d2 - d1)
			for i := 0; i < 6; i++ {
				x2[i] = x1[i] + k*(x2[i]-x1[i])
			}
			path[pathidx] = x2
			pathidx++
			tot += math.Sqrt((x2[0]-x1[0])*(x2[0]-x1[0]) + (x2[1]-x1[1])*(x2[1]-x1[1]) + (x2[2]-x1[2])*(x2[2]-x1[2]))
			break
		}
		tot += math.Sqrt((x2[0]-x1[0])*(x2[0]-x1[0]) + (x2[1]-x1[1])*(x2[1]-x1[1]) + (x2[2]-x1[2])*(x2[2]-x1[2]))
		x1 = x2
		path[pathidx] = x2
		pathidx++
	}
	var rst [4]float64
	rst[0] = x2[0]*cos - x2[2]*sin
	rst[1] = x2[1]
	rst[2] = x2[0]*sin + x2[2]*cos
	rst[3] = tot
	return rst, path[0:pathidx]
}

func initRig(endPoint float64) float64 {
	return C0 + C1*endPoint + C2*endPoint*endPoint
}

func Init(bmap []float64, iAng, oAng, tDist, fdc0Dist, fdc2Dist float64) ([]float64, []float64) {

	rigs := make([]float64, 10000)
	poss := make([]float64, 10000)
	
	cos := math.Cos(oAng/180.0*math.Pi)
	sin := math.Sin(oAng/180.0*math.Pi)
	sta := [6]float64{fdc0Dist * sin, 0.0, fdc0Dist * cos, 0.0, 0.0, 0.0}
	for i := 0; i < 10000; i++ {
		rigs[i] = 3000.0 + 0.2*float64(i)
		sta[3] = -1.0 * rigs[i] * sin
		sta[4] = 0.0
		sta[5] = rigs[i] * cos
		tmp, _ := Trace(sta, fdc2Dist, oAng, 0.001, bmap)
		poss[i] = tmp[0]
		rigs[i] = rigs[i] * 2.99792458 / 10.0
	}
	C0, C1, C2, C3 = fit3(poss, rigs)
	return rigs, poss
}

func GuessRig(position float64) float64 {
	//return C0 + C1*position + C2*position*position + C3*position*position*position
	//fmt.Printf("guess: %x %f\n",math.Float64bits(position),position)
	rst := math.FMA(C1,position, C0)
	//fmt.Printf("guess: %x %f\n",math.Float64bits(rst),rst)
	rst = math.FMA(C2,position*position,rst)
	//fmt.Printf("guess: %x %f\n",math.Float64bits(rst),rst)
	rst = math.FMA(C3,position*position*position,rst)
	//fmt.Printf("guess: %x %f\n",math.Float64bits(rst),rst)
	return rst
}
