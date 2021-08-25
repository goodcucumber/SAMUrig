package rig

/*
func fetch6(r [6]float64, bmap []float64) [3]float64 {
	rst := [3]float64{0.0, 0.0, 0.0}
	x := r[0] / 10.0
	y := r[1]/10.0 + 40.0
	z := r[2] / 10.0
	x0 := math.Abs(x)
	y0 := y
	z0 := math.Abs(z)
	if x0 >= 300.0 || z0 >= 300.0 || y0 >= 80.0 || y0 < 0 {
		return rst
	}
	if math.IsNaN(x + y + z) {
		return rst
	}
	xp := int(math.Floor(x0))
	yp := int(math.Floor(y0))
	zp := int(math.Floor(z0))

	xr := x0 - math.Floor(x0)
	yr := y0 - math.Floor(y0)
	zr := z0 - math.Floor(z0)

	for i := 0; i < 3; i++ {
		v1 := bmap[((zp+0)+(xp+0)*302+(yp+0)*302*302)*3+i]
		v2 := bmap[((zp+1)+(xp+0)*302+(yp+0)*302*302)*3+i]
		v3 := bmap[((zp+1)+(xp+1)*302+(yp+0)*302*302)*3+i]
		v4 := bmap[((zp+0)+(xp+1)*302+(yp+0)*302*302)*3+i]
		v5 := bmap[((zp+0)+(xp+0)*302+(yp+1)*302*302)*3+i]
		v6 := bmap[((zp+1)+(xp+0)*302+(yp+1)*302*302)*3+i]
		v7 := bmap[((zp+1)+(xp+1)*302+(yp+1)*302*302)*3+i]
		v8 := bmap[((zp+0)+(xp+1)*302+(yp+1)*302*302)*3+i]
		bxyz := v1*(1-xr)*(1-yr)*(1-zr) + v2*xr*(1-yr)*(1-zr) +
			v3*xr*yr*(1-zr) + v4*(1-xr)*yr*(1-zr) +
			v5*(1-xr)*(1-yr)*zr + v6*xr*(1-yr)*zr +
			v7*xr*yr*zr + v8*(1-xr)*yr*zr
		rst[i] = bxyz
		//rst[i] = (v5*(1-xr)+v6*xr)*(1-yr) + (v7*xr+v8*(1-xr))*yr
	}
	if x < 0 {
		rst[0] = -rst[0]
	}
	if z < 0 {
		rst[2] = -rst[2]
	}
	return rst
}
*/
