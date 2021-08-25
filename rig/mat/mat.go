package mat

import (
	"fmt"
	"math"
)

type Mat struct {
	Ele [][]float64
	M   int
	N   int
}

func Gen(o []float64, m, n int) Mat {
	if len(o) < 1 {
		return Mat{nil, 0, 0}
	}
	l := len(o)
	c := make([][]float64, m)
	for i := 0; i < m; i++ {
		c[i] = make([]float64, n)
		for j := 0; j < n; j++ {
			c[i][j] = o[(i*n+j)%l]
		}
	}
	return Mat{c, m, n}
}

func Mult(A, B Mat) Mat {
	if A.Ele == nil || B.Ele == nil {
		return Mat{nil, 0, 0}
	}
	if A.N != B.M {
		return Mat{nil, 0, 0}
	}
	c := make([][]float64, A.M)
	for i := 0; i < A.M; i++ {
		c[i] = make([]float64, B.N)
		for j := 0; j < B.N; j++ {
			c[i][j] = 0
			for k := 0; k < A.N; k++ {
				c[i][j] += A.Ele[i][k] * B.Ele[k][j]
			}
		}
	}
	return Mat{c, A.M, B.N}
}

func Insert0(A Mat, r, c int) Mat {
	m, n := A.M+1, A.N+1
	if r < 0 || r >= m {
		m = A.M
	}
	if c < 0 || c >= n {
		n = A.N
	}
	cc := make([][]float64, m)
	id1 := 0
	id2 := 0
	for i := 0; i < m; i++ {
		cc[i] = make([]float64, n)
		id2 = 0
		for j := 0; j < n; j++ {
			cc[i][j] = A.Ele[id1][id2]
			if i == r || j == c {
				cc[i][j] = 0
			} else {
				id2++
			}
		}
		if i != r {
			id1++
		}
	}
	return Mat{cc, m, n}
}

func Drop(A Mat, r, c int) Mat {
	m, n := A.M-1, A.N-1
	if r < 0 || r > m {
		m = A.M
	}
	if c < 0 || c > n {
		n = A.N
	}
	cc := make([][]float64, m)
	id1 := 0
	id2 := 0
	for i := 0; i < m; i++ {
		cc[i] = make([]float64, n)
		if id1 == r {
			id1++
		}
		id2 = 0
		for j := 0; j < m; j++ {
			if id2 == c {
				id2++
			}
			cc[i][j] = A.Ele[id1][id2]
			id2++
		}
		id1++
	}
	return Mat{cc, m, n}
}

func T(A Mat) Mat {
	c := make([][]float64, A.N)
	for i := 0; i < A.N; i++ {
		c[i] = make([]float64, A.M)
		for j := 0; j < A.M; j++ {
			c[i][j] = A.Ele[j][i]
		}
	}
	return Mat{c, A.N, A.M}
}

func Det(A Mat) float64 {
	if A.M != A.N {
		return math.NaN()
	}
	if A.M == 1 {
		return A.Ele[0][0]
	}
	if A.M == 2 {
		return A.Ele[0][0]*A.Ele[1][1] - A.Ele[0][1]*A.Ele[1][0]
	}
	rst := 0.0
	sign := 1.0
	for i := 0; i < A.N; i++ {
		if i%2 == 0 {
			sign = 1.0
		} else {
			sign = -1.0
		}
		rst += sign * A.Ele[0][i] * Det(Drop(A, 0, i))
	}
	return rst
}

func Rbind(A, B Mat) Mat {
	if A.N != B.N {
		return Mat{nil, 0, 0}
	}
	c := make([][]float64, A.M+B.M)
	for i := 0; i < A.M; i++ {
		c[i] = make([]float64, A.N)
		for j := 0; j < A.N; j++ {
			c[i][j] = A.Ele[i][j]
		}
	}
	for i := 0; i < B.M; i++ {
		c[i] = make([]float64, A.N)
		for j := 0; j < B.N; j++ {
			c[i+A.M][j] = B.Ele[i][j]
		}
	}
	return Mat{c, A.M + B.M, A.N}
}

func CBind2(A, B Mat) Mat {
	if A.M != B.M {
		return Mat{nil, 0, 0}
	}
	c := make([][]float64, A.M)
	for i := 0; i < A.M; i++ {
		c[i] = make([]float64, A.N+B.N)
		for j := 0; j < A.N; j++ {
			c[i][j] = A.Ele[i][j]
		}
		for j := 0; j < B.N; j++ {
			c[i][j+A.N] = B.Ele[i][j]
		}
	}
	return Mat{c, A.M, A.N + B.N}
}
func CBind(A ...Mat) Mat {
	if len(A) == 1 {
		return A[0]
	}
	rst := A[0]
	for _, v := range A[1:len(A)] {
		rst = CBind2(rst, v)
	}
	return rst
}

func Inv(A Mat) Mat {
	if A.M != A.N {
		return Mat{nil, 0, 0}
	}
	c := make([][]float64, A.M)
	sign := 1.0
	det := Det(A)
	if det == 0 && A.M == 3 {
		B := Drop(A, 1, 1)
		C := Inv(B)
		return Insert0(C, 1, 1)
	}
	for i := 0; i < A.M; i++ {
		c[i] = make([]float64, A.N)
		for j := 0; j < A.N; j++ {
			if (i+j)%2 == 0 {
				sign = 1.0
			} else {
				sign = -1.0
			}
			c[i][j] = sign * Det(Drop(A, i, j)) / det
		}
	}
	return T(Mat{c, A.M, A.N})
}

func (this *Mat) Show() {
	for i := 0; i < this.M; i++ {
		for j := 0; j < this.N; j++ {
			fmt.Printf("%+3.5f  ", this.Ele[i][j])
		}
		fmt.Println("")
	}
}
