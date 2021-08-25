package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"os"
	"runtime"
	"time"

	"./rig"
	"./rig/codec"
)

type ipack struct {
	id int
	s1 string
	s2 string
	s3 string
}
type opack struct {
	id     int
	count  int
	result []float64
	rststr string
}

var ipipe = make(chan ipack, 10000)
var opipe = make(chan opack, 10000)
var npipe = make(chan int) // count of events
var spipe = make(chan int) // status
var msg = make(chan int, 5)

const maxev = 20000000
const maxwk = 16 // max workers

var unlimited bool = false
var limn = 50000

func output(name string) {
	fout, err := os.Create(name)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error creating out file: ", err)
		os.Exit(0)
	}
	defer fout.Close()
	fmt.Fprintln(fout, "id count xt yt kbx kby x1 y1 kx1 ky1 kxt1 kyt1 x2 kx2 rig tot cx2 cy2 cz2 ckx2 cky2")
	total := -1
	storage := make([]opack, maxev)
	occupacy := make([]bool, maxev)
	current := 0
	stored := 0
	spipe <- 0
	dftcnt := 0
loop:
	for {
		select {
		case e := <-opipe:

			storage[e.id] = e
			occupacy[e.id] = true
			stored++
			if stored%100 == 0 {
				//fmt.Fprintln(os.Stderr, stored)
				msg <- 0
			}
			if stored >= total && total > 0 {
				break loop
			}
		case n := <-npipe:

			total = n
			if stored >= n {
				break loop
			}
		default:
			if dftcnt < 15000 {
				dftcnt++
				break
			}
			dftcnt = 0
			enough := true
			for k := 0; k < 20; k++ {
				enough = enough && occupacy[current+k]
				if enough == false {
					break
				}
			}
			if enough {
				for k := 0; k < 20; k++ {
					fmt.Fprint(fout, storage[current].rststr)
					current++
				}
			}
		}
	}
	fmt.Println("\nAll collected.")
	for i := current; i < total; i++ {
		fmt.Fprint(fout, storage[i].rststr)
	}
	spipe <- 1
}
func work(f func(int, string, string, string) (int, []float64)) {
	for {
		x := <-ipipe
		msg <- 1
		//fmt.Println(x.id)
		cnt, rst := f(x.id, x.s1, x.s2, x.s3)
		str := fmt.Sprintf("%d %d", x.id, cnt)
		if rst == nil {
			for i := 0; i < 19; i++ {
				str += " NA"
			}
		} else {
			for k := 0; k < len(rst); k++ {
				str += fmt.Sprintf(" %16.10f", rst[k])
			}
		}
		str += fmt.Sprintln("")
		msg <- 2
		opipe <- opack{id: x.id, count: cnt, result: rst, rststr: str}
		msg <- 3
	}
}

func monitor() {
	wcount := 0
	pcount := 0
	ecount := 0
	totalt := 0.0
	snap := 0
	tk := time.NewTicker(2 * time.Second)
	fmt.Fprintln(os.Stdout, "|     count |  #Working |  #Waiting | Average Speed | Current Speed |")

	for {
		select {
		case num := <-msg:
			if num == 1 {
				wcount += 1
			} else if num == 2 {
				wcount -= 1
				pcount += 1
			} else if num == 3 {
				pcount -= 1
			} else if num == 0 {
				ecount = ecount + 1
			}
		case <-tk.C:
			totalt += 2
			fmt.Fprint(os.Stdout, "\r\033[K")
			fmt.Fprintf(os.Stdout, "| %8dk |  %8d |  %8d | %9.4fke/s | %9.4fke/s |", ecount/10, wcount, pcount, float64(ecount)/totalt/10.0, float64(ecount-snap)/20.0)
			snap = ecount
		}
	}
}

func main() {
	//cpuf, _ := os.Create("cpuprofd1.txt")
	runtime.GOMAXPROCS(16)
	var fbdc *os.File
	var ffdc *os.File
	var fdck *os.File
	var err error
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "add parameter(s)! <1> = bdc file; <2> = fdc0 file; <3> = fdc2k file; <4> = output; [<5>] = dump track only for a specific event")
		return
	}

	fbdc, err = os.Open(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	defer fbdc.Close()

	ffdc, err = os.Open(os.Args[2])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	defer ffdc.Close()

	fdck, err = os.Open(os.Args[3])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	defer fdck.Close()

	byb, _ := ioutil.ReadAll(fbdc)
	brb := bytes.NewReader(byb)
	bb := bufio.NewReader(brb)
	byf, _ := ioutil.ReadAll(ffdc)
	brf := bytes.NewReader(byf)
	bf := bufio.NewReader(brf)
	byk, _ := ioutil.ReadAll(fdck)
	brk := bytes.NewReader(byk)
	bk := bufio.NewReader(brk)


	bb.ReadLine()
	bf.ReadLine()
	bk.ReadLine()
	
	bmp0 := codec.Load("bmap.bin")
	if bmp0.BX == nil {
		fmt.Fprintln(os.Stderr, "Make field binary:")
		bmp0 = codec.Encode("180703-1,40T-3000.table", "180703-1,45T-3000.table", 400.0/500.0, "bmap.bin")
	}
	
	bmp := bmp0.Backend
	
	//beamline angle: 30 deg.
	//direction to FDC2k: 15 dig.
	//distance from target to SAMURAI center: 4500
	//distance from FDC0 to SAMURAI center: 3000
	//distance from SAMURAI center to FDC2: 4000
	rig.Init(bmp, 30, 15, 4500, 3000, 4000)
	
	fwk := func(id int, s1, s2, s3 string) (int, []float64) {
		//beamline angle: 30 deg.
		//direction to FDC2k: 15 dig.
		//distance from target to SAMURAI center: 4500
		//distance from FDC0 to SAMURAI center: 3000
		//distance from SAMURAI center to FDC2: 4000
		cnt, rst := rig.RigWork(s1, s2, s3, bmp, 30, 15, 4500, 3000, 4000)
		return cnt, rst
	}
	if len(os.Args) == 5 {
		
		go output(os.Args[4])
		go monitor()
		<-spipe
		for i := 0; i < maxwk; i++ {
			go work(fwk)
		}
		evtn := 0

		for unlimited || evtn < limn {
			s1, p1, e1 := bb.ReadLine()
			s2, p2, e2 := bf.ReadLine()
			s3, p3, e3 := bk.ReadLine()

			if p1 && p2 && p3 == false {
				fmt.Fprintln(os.Stderr, "Error while reading...")
				fmt.Fprintln(os.Stderr, "    ", p1, e1)
				fmt.Fprintln(os.Stderr, "    ", p2, e2)
				fmt.Fprintln(os.Stderr, "    ", p3, e3)
				return
			}
			if e1 != nil && e1 != io.EOF {
				fmt.Fprintln(os.Stderr, "error on bdc file: ", e1)
			}
			if e2 != nil && e2 != io.EOF {
				fmt.Fprintln(os.Stderr, "error on bdc file: ", e2)
			}
			if e3 != nil && e3 != io.EOF {
				fmt.Fprintln(os.Stderr, "error on bdc file: ", e3)
			}
			if e1 == io.EOF || e2 == io.EOF || e3 == io.EOF {
				fmt.Fprintln(os.Stderr, e1, e2, e3)
				break
			}
			ipipe <- ipack{id: evtn, s1: string(s1), s2: string(s2), s3: string(s3)}
			evtn++
		}
		npipe <- evtn
		<-spipe
	} else {
		target := 0
		fmt.Sscan(os.Args[5], &target)
		evtn := -1
		for {
			s1, p1, e1 := bb.ReadLine()
			s2, p2, e2 := bf.ReadLine()
			s3, p3, e3 := bk.ReadLine()
			evtn++
			if evtn == target {
				fout, _ := os.Create("dump.txt")
				//beamline angle: 30 deg.
				//direction to FDC2k: 15 dig.
				//distance from target to SAMURAI center: 4500
				//distance from FDC0 to SAMURAI center: 3000
				//distance from SAMURAI center to FDC2: 4000
				cnt, rst := rig.RigWork(string(s1), string(s2), string(s3), bmp, 30, 15, 4500, 3000, 4000)
				x1 := rst[4]
				y1 := rst[5]
				kx := rst[8]
				ky := rst[9]
				rig0 := rst[12]
				cos := math.Cos(math.Pi / 6.0)
				sin := math.Sin(math.Pi / 6.0)
				kn := math.Sqrt(1.0 + kx*kx + ky*ky)
				vx := kx / kn
				vy := ky / kn
				vz := 1.0 / kn
				vxn := vx*cos - vz*sin
				vzn := vx*sin + vz*cos
				start := [6]float64{3000*sin + x1*cos, y1, x1*sin - 3000*cos, vxn*rig0/ 2.99792458 * 10.0, vy*rig0/ 2.99792458 * 10.0, vzn*rig0/ 2.99792458 * 10.0}
	
				_, rst2 := rig.TraceWithPath(start, 4000, 15, 0.001, bmp)
				start[0] = rig0
				fmt.Fprint(fout, evtn, cnt)
				for i := 0; i < len(rst); i++ {
					fmt.Fprintf(fout, "%16.10f ", rst[i])
				}
				fmt.Fprintln(fout, "")
				for i := 0; i < len(rst2); i++ {
					fmt.Fprintf(fout, "%16.10f %16.10f ", rst2[i][0], rst2[i][1])
					fmt.Fprintf(fout, "%16.10f %16.10f ", rst2[i][2], rst2[i][3])
					fmt.Fprintf(fout, "%16.10f %16.10f\n", rst2[i][4], rst2[i][5])
				return
			}
			if p1 && p2 && p3 == false {
				fmt.Fprintln(os.Stderr, "Error while reading...")
				fmt.Fprintln(os.Stderr, "    ", p1, e1)
				fmt.Fprintln(os.Stderr, "    ", p2, e2)
				fmt.Fprintln(os.Stderr, "    ", p3, e3)
				return
			}
			if e1 != nil && e1 != io.EOF {
				fmt.Fprintln(os.Stderr, "error on bdc file: ", e1)
			}
			if e2 != nil && e2 != io.EOF {
				fmt.Fprintln(os.Stderr, "error on bdc file: ", e2)
			}
			if e3 != nil && e3 != io.EOF {
				fmt.Fprintln(os.Stderr, "error on bdc file: ", e3)
			}
			if e1 == io.EOF || e2 == io.EOF || e3 == io.EOF {
				fmt.Fprintln(os.Stderr, e1, e2, e3)
				break
			}
		}
	}
	fmt.Println("")
}
