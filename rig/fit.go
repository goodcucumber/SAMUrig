package rig

import "./mat"

func fit2(xs, ys []float64) (float64, float64, float64) {
	minlen := len(xs)
	if len(ys) < minlen {
		minlen = len(ys)
	}
	xs2 := make([]float64, minlen)
	for i := 0; i < minlen; i++ {
		xs2[i] = xs[i] * xs[i]
	}
	matx1 := mat.Gen(xs, minlen, 1)
	matx2 := mat.Gen(xs2, minlen, 1)
	f1 := make([]float64, 1)
	f1[0] = 1.0
	mat1 := mat.Gen(f1, minlen, 1)
	maty := mat.Gen(ys, minlen, 1)
	matx := mat.CBind(mat1, mat.CBind(matx1, matx2))
	matxtxi := mat.Inv(mat.Mult(mat.T(matx), matx))
	matxtxixt := mat.Mult(matxtxi, mat.T(matx))
	rst := mat.Mult(matxtxixt, maty)
	return rst.Ele[0][0], rst.Ele[1][0], rst.Ele[2][0]
}

//var Fit = fit2

func fit3(xs, ys []float64) (float64, float64, float64, float64) {
	minlen := len(xs)
	if len(ys) < minlen {
		minlen = len(ys)
	}
	xs2 := make([]float64, minlen)
	xs3 := make([]float64, minlen)
	for i := 0; i < minlen; i++ {
		xs2[i] = xs[i] * xs[i]
		xs3[i] = xs[i] * xs[i] * xs[i]
	}
	matx1 := mat.Gen(xs, minlen, 1)
	matx2 := mat.Gen(xs2, minlen, 1)
	matx3 := mat.Gen(xs3, minlen, 1)
	f1 := make([]float64, 1)
	f1[0] = 1.0
	mat1 := mat.Gen(f1, minlen, 1)
	maty := mat.Gen(ys, minlen, 1)
	matx := mat.CBind(mat1, matx1, matx2, matx3)
	matxtxi := mat.Inv(mat.Mult(mat.T(matx), matx))
	matxtxixt := mat.Mult(matxtxi, mat.T(matx))
	rst := mat.Mult(matxtxixt, maty)
	return rst.Ele[0][0], rst.Ele[1][0], rst.Ele[2][0], rst.Ele[3][0]
}
