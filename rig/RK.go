package rig

import "math"

/*
func eqf(xn [6]float64, bmap []float64) [6]float64 {
	//rst := [6]float64{xn[1], 0.0, xn[3], 0.0, xn[5], 0.0}
	rst := [6]float64{xn[3], xn[4], xn[5], 0.0, 0.0, 0.0}
	//bn := FetchB(xn[:], bmap)
	bn := fetch6(xn, bmap)

	rst[3] = -xn[4]*bn[2] + xn[5]*bn[1]
	rst[4] = -xn[5]*bn[0] + xn[3]*bn[2]
	rst[5] = -xn[3]*bn[1] + xn[4]*bn[0]
	return rst
}

func CompEqf(x [6]float64, bmp *codec.BMap) {
	bmp1 := codec.Conv(bmp)
	bmp2 := codec.Conv2(bmp)
	bmp3 := codec.Conv2(bmp)
	fmt.Println("ori ", x)
	fmt.Println("bmp  ", fetch6(x, bmp1))
	fmt.Println("eqf0 ", eqf(x, bmp1))
	fmt.Println("eqf0c", eqfc(x, bmp3))
	fmt.Println("eqfa5", eqfa5(x, bmp1))
	fmt.Println("eqfar", eqfar(x, bmp2))
	fmt.Println("")
}

func eqfc(xn [6]float64, bmp []float64) [6]float64 {
	bf := [3]float64{0.0, 0.0, 0.0}
	rst := [6]float64{xn[3], xn[4], xn[5], 0.0, 0.0, 0.0}
	x := xn[0] * 0.1
	y := xn[1]*0.1 + 40.0
	z := xn[2] * 0.1
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
		v1 := bmp[((xp>>1)+(yp>>1)*151+(zp>>1)*151*41)*24+(xp&1)+2*(yp&1)+4*(zp&1)+i*8]
		v2 := bmp[(((xp+1)>>1)+(yp>>1)*151+(zp>>1)*151*41)*24+((xp+1)&1)+2*(yp&1)+4*(zp&1)+i*8]
		v3 := bmp[(((xp+1)>>1)+((yp+1)>>1)*151+(zp>>1)*151*41)*24+((xp+1)&1)+2*((yp+1)&1)+4*(zp&1)+i*8]
		v4 := bmp[((xp>>1)+((yp+1)>>1)*151+(zp>>1)*151*41)*24+(xp&1)+2*((yp+1)&1)+4*(zp&1)+i*8]
		v5 := bmp[((xp>>1)+(yp>>1)*151+((zp+1)>>1)*151*41)*24+(xp&1)+2*(yp&1)+4*((zp+1)&1)+i*8]
		v6 := bmp[(((xp+1)>>1)+(yp>>1)*151+((zp+1)>>1)*151*41)*24+((xp+1)&1)+2*(yp&1)+4*((zp+1)&1)+i*8]
		v7 := bmp[(((xp+1)>>1)+((yp+1)>>1)*151+((zp+1)>>1)*151*41)*24+((xp+1)&1)+2*((yp+1)&1)+4*((zp+1)&1)+i*8]
		v8 := bmp[((xp>>1)+((yp+1)>>1)*151+((zp+1)>>1)*151*41)*24+(xp&1)+2*((yp+1)&1)+4*((zp+1)&1)+i*8]
		bxyz := v1*(1-xr)*(1-yr)*(1-zr) + v2*xr*(1-yr)*(1-zr) +
			v3*xr*yr*(1-zr) + v4*(1-xr)*yr*(1-zr) +
			v5*(1-xr)*(1-yr)*zr + v6*xr*(1-yr)*zr +
			v7*xr*yr*zr + v8*(1-xr)*yr*zr
		bf[i] = bxyz
		//rst[i] = (v5*(1-xr)+v6*xr)*(1-yr) + (v7*xr+v8*(1-xr))*yr
	}
	if x < 0 {
		bf[0] = -bf[0]
	}
	if z < 0 {
		bf[2] = -bf[2]
	}
	rst[3] = -xn[4]*bf[2] + xn[5]*bf[1]
	rst[4] = -xn[5]*bf[0] + xn[3]*bf[2]
	rst[5] = -xn[3]*bf[1] + xn[4]*bf[0]

	return rst
}
*/
var cntcnt = 0

func next(xn [6]float64, bmap []float64, h float64) [6]float64 {
	//fmt.Println("   k1i",xn)
	k1 := eqfar(xn, bmap)
	//fmt.Println("    k1",k1,b)
	var k2 [6]float64
	for i := 0; i < 6; i++ {
		//k2[i] = xn[i] + k1[i]*0.5*h
		k2[i] = math.FMA(k1[i], 0.5*h, xn[i])
	}
	//fmt.Println("   k2i",k2)
	k2 = eqfar(k2, bmap)
	//fmt.Println("    k2",k2,b)
	var k3 [6]float64
	for i := 0; i < 6; i++ {
		//k3[i] = xn[i] + k2[i]*0.5*h
		k3[i] = math.FMA(k2[i], 0.5*h, xn[i])
	}
	//fmt.Println("   k3i",k3)
	k3 = eqfar(k3, bmap)
	//fmt.Println("    k3",k3,b)

	var k4 [6]float64
	for i := 0; i < 6; i++ {
		//k4[i] = xn[i] + k3[i]*h
		k4[i] = math.FMA(k3[i], 0.5*h, xn[i])
	}
	//fmt.Println("   k4i",k4)
	k4 = eqfar(k4, bmap)
	//fmt.Println("    k4",k4,b)

	var rst [6]float64
	for i := 0; i < 6; i++ {
		//rst[i] = xn[i] + (h/6.0)*(k1[i]+2*k2[i]+2*k3[i]+k4[i])
		tmp1 := math.FMA(k2[i], 2.0, k1[i])
		tmp2 := math.FMA(k3[i], 2.0, k4[i])
		rst[i] = math.FMA(tmp1+tmp2, h/6.0, xn[i])
	}
	//fmt.Println(cntcnt, rst)
	//cntcnt++
	return rst
}
var Next = next