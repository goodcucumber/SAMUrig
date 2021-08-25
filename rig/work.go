package rig

import (
	"fmt"
	"math"
)

func RigWork(BDC, FDC0, FDC2K string, bmap []float64, iAng, oAng, tDist, fdc0Dist, fdc2Dist float64) (int, []float64) {
	//BDC:    id kx ky tx ty rx ry x1 y1 x2 y2; tx,ty => target position
	//FDC0:   idx L0 L1 L2 L3 L4 L5 L6 L7 bx by kx ky rx ry
	//FDC2K:  idx w1 t1 p1 w2 t2 p2 w3 t3 p3 k b r; b => intercept at FDC2 fdc2Dist

	dummy := ""
	id0 := -1
	id1 := -1
	id2 := -1
	xt := math.NaN()
	yt := math.NaN()
	kbx := math.NaN()
	kby := math.NaN()
	x1 := math.NaN()
	y1 := math.NaN()
	k1x := math.NaN()
	k1y := math.NaN()
	x2 := math.NaN()
	k2x := math.NaN()
	fmt.Sscan(BDC, &id0, &kbx, &kby, &xt, &yt, &dummy, &dummy, &dummy, &dummy, &dummy, &dummy)
	//bdc:        right hand coordinate. beam=z+ (front), up=y+, left=x+
	//others:     left hand.             beam=z+ (front), up=y+, right=x+
	//all output: left
	fmt.Sscan(FDC0, &id1, &k1x, &x1, &dummy, &k1y, &y1, &dummy, &dummy, &dummy, &dummy, &dummy, &dummy, &dummy, &dummy, &dummy)
	fmt.Sscan(FDC2K, &id2, &dummy, &dummy, &dummy, &dummy, &dummy, &dummy, &dummy, &dummy, &dummy, &k2x, &x2, &dummy)

	if math.IsNaN(xt) || math.IsNaN(x1) || math.IsNaN(x2) {
		return 0, nil
	}
	if id0 != id1 || id0 != id2 || id1 != id2 || id0 < 0 {
		return 0, nil
	}
	cos := math.Cos(iAng/180.0*math.Pi)
	sin := math.Sin(iAng/180.0*math.Pi)

	kx := (x1 + xt) / (tDist-fdc0Dist)
	ky := (y1 - yt) / (tDist-fdc0Dist)
	kn := math.Sqrt(1.0 + kx*kx + ky*ky)
	vx := kx / kn
	vy := ky / kn
	vz := 1.0 / kn
	vxn := vx*cos - vz*sin
	vzn := vx*sin + vz*cos
	start := [6]float64{fdc0Dist*sin + x1*cos, y1, x1*sin - fdc0Dist*cos, vxn, vy, vzn}
	rig, pos, cnt, fin := SmartFindTrace2(start, fdc2Dist, oAng, 0.001, x2, bmap)
	rst := make([]float64, 19)
	rst[0] = -xt
	rst[1] = yt
	rst[2] = -kbx
	rst[3] = kby
	rst[4] = x1
	rst[5] = y1
	rst[6] = k1x
	rst[7] = k1y
	rst[8] = kx
	rst[9] = ky
	rst[10] = x2
	rst[11] = k2x
	rst[12] = rig
	rst[13] = pos[3] // path length
	rst[14] = pos[0]
	rst[15] = pos[1]
	rst[16] = pos[2]
	v2z := fin[3]*math.Sin(oAng/180.0*math.Pi) + fin[5]*math.Cos(oAng/180.0*math.Pi)
	v2y := fin[4]
	v2x := fin[3]*math.Cos(oAng/180.0*math.Pi) - fin[5]*math.Sin(oAng/180.0*math.Pi)
	rst[17] = v2x / v2z
	rst[18] = v2y / v2z
	return cnt, rst
}
