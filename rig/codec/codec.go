package codec

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"os"
	"runtime"
)

type BMap struct {
	//Backend [301 * 81 * 301 * 3]float64
	Backend []float64
	BX      []float64
}

//Encode field map to bin
//name1 lower
//name2 upper
func Encode(name1, name2 string, r float64, outname string) BMap {
	fmt.Println("Start")
	var bmap BMap
	bmap.Backend = make([]float64, 302*83*302*3)
	bmap.BX = nil
	f1, e1 := os.Open(name1)
	if e1 != nil {
		fmt.Fprintln(os.Stderr, "Error[encoder]: f1 error: ", e1)
		return bmap
	}
	defer f1.Close()
	f2, e2 := os.Open(name1)
	if e2 != nil {
		fmt.Fprintln(os.Stderr, "Error[encoder]: f2 error: ", e2)
		return bmap
	}
	defer f2.Close()
	fout, e3 := os.Create(outname)
	if e3 != nil {
		fmt.Fprintln(os.Stderr, "Error[encoder]: fout error: ", e3)
		return bmap
	}
	defer fout.Close()
	bmap.BX = bmap.Backend
	for i := 0; i < 3*302*302*83; i++ {
		bmap.Backend[i] = 0.0
	}
	dummy := ""
	for i := 0; i < 8; i++ {
		fmt.Fscanln(f1, &dummy, &dummy, &dummy, &dummy, &dummy, &dummy, &dummy, &dummy)
		fmt.Fscanln(f2, &dummy, &dummy, &dummy, &dummy, &dummy, &dummy, &dummy, &dummy)
	}
	br1 := bufio.NewReader(f1)
	br2 := bufio.NewReader(f2)
	ratio := 0.0

	oldratio := 0.0
	/*readf := func(bbx, bby, bbz []float64, s string, r0 float64) {
		var px, py, pz, bx, by, bz float64
		n, e := fmt.Sscan(s, &px, &py, &pz, &bx, &by, &bz)
		if n != 6 || e != nil {
			panic(e)
		}
		x := int(px) / 10
		y := int((py + 410.0)) / 10
		z := int(pz) / 10
		bbx[z+x*302+y*302*302] += bx * r0
		bby[z+x*302+y*302*302] += by * r0
		bbz[z+x*302+y*302*302] += bz * r0
		//cpipe <- 1
	}*/
	runtime.GC()
	fmt.Println("DO!")
	for cnt := 0; cnt < 301*81*301; cnt++ {
		//if (cnt+1)%300*80*30 == 0 {
		//	runtime.GC()
		//}
		var px, py, pz, bx, by, bz float64
		l, _, _ := br1.ReadLine()
		n, e := fmt.Sscan(string(l), &px, &py, &pz, &bx, &by, &bz)
		if n != 6 || e != nil {
			panic(e)
		}
		x := int(px) / 10
		y := int((py + 410.0)) / 10
		z := int(pz) / 10
		bmap.Backend[(z+x*302+y*302*302)*3] += bx * (1.0 - r)
		bmap.Backend[(z+x*302+y*302*302)*3+1] += by * (1.0 - r)
		bmap.Backend[(z+x*302+y*302*302)*3+2] += bz * (1.0 - r)
		ratio = float64(cnt+1) / (301.0 * 81.0 * 301.0) * 100.0
		if ratio-oldratio > 0.2 || ratio == 100.0 {
			oldratio = ratio
			fmt.Fprint(os.Stdout, "\r\033[K")
			fmt.Printf("Encoding1: %4.1f %%", ratio)
		}
	}
	fmt.Println("")
	oldratio = 0.0
	runtime.GC()
	for cnt := 0; cnt < 301*81*301; cnt++ {
		var px, py, pz, bx, by, bz float64
		l, _, _ := br2.ReadLine()
		n, e := fmt.Sscan(string(l), &px, &py, &pz, &bx, &by, &bz)
		if n != 6 || e != nil {
			panic(e)
		}
		x := int(px) / 10
		y := int((py + 410.0)) / 10
		z := int(pz) / 10
		bmap.Backend[(z+x*302+y*302*302)*3] += bx * r
		bmap.Backend[(z+x*302+y*302*302)*3+1] += by * r
		bmap.Backend[(z+x*302+y*302*302)*3+2] += bz * r

		//readf(bmap.BX, bmap.BY, bmap.BZ, string(l), r)
		ratio = float64(cnt+1) / (301.0 * 81.0 * 301.0) * 100.0
		if ratio-oldratio > 0.2 || ratio == 100.0 {
			oldratio = ratio
			fmt.Fprint(os.Stdout, "\r\033[K")
			fmt.Printf("Encoding2: %4.1f %%", ratio)
		}
	}

	binary.Write(fout, binary.LittleEndian, bmap.Backend)
	fmt.Print("\nEncode done.\n")
	return bmap
}

func Load(name string) BMap {
	var bmap BMap
	bmap.Backend = make([]float64, 302*83*302*3)
	bmap.BX = nil
	f, e := os.Open(name)
	if e != nil {
		fmt.Fprintln(os.Stderr, "Error[load]: open file :", e)
		return bmap
	}
	defer f.Close()
	e = binary.Read(f, binary.LittleEndian, bmap.Backend)
	if e != nil {
		fmt.Fprintln(os.Stderr, "Error[load]: binary.Read :", e)
		return bmap
	}
	bmap.BX = bmap.Backend
	return bmap
}

func Conv2(bmap *BMap) []float64 {
	/*
		rst := make([]float64, 301*81*301*3)
		for i := 0; i < 301; i++ {
			for j := 0; j < 81; j++ {
				for k := 0; k < 301; k++ {
					rst[(k+i*301+j*301*301)*3+0] = bmap.BX[i+j*301+k*301*81]
					rst[(k+i*301+j*301*301)*3+1] = bmap.BY[i+j*301+k*301*81]
					rst[(k+i*301+j*301*301)*3+2] = bmap.BZ[i+j*301+k*301*81]
				}
			}
		}
		return rst
	*/
	return bmap.Backend
}
