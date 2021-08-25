package rig

import (
	"math"
)

const maxerr = 0.002

func SmartFindTrace(start [6]float64, dist float64, angle float64, h float64, target float64, bmap []float64) (float64, [4]float64, int, [6]float64) {
	rig1 := GuessRig(target)
	sta1 := start
	for i := 3; i < 6; i += 1 {
		sta1[i] = start[i] * rig1 / 2.99792458 * 10.0
	}
	pos1, fin1 := Trace(sta1, dist, angle, h, bmap)
	count := 1
	//fmt.Println("DBG0: ", count, math.Abs(pos1[0]-target))
	if math.Abs(pos1[0]-target) < maxerr {
		return rig1, pos1, count, fin1
	}
	rig2 := rig1 + (target-pos1[0])*C1
	sta2 := start
	for i := 3; i < 6; i += 1 {
		sta2[i] = start[i] * rig2 / 2.99792458 * 10.0
	}
	pos2, fin2 := Trace(sta2, dist, angle, h, bmap)
	count = 2
	//fmt.Println("DBG1: ", count, math.Abs(pos2[0]-target))
	if math.Abs(pos2[0]-target) < maxerr {
		return rig2, pos2, count, fin2
	}
	if math.Abs(pos2[0]-target) > math.Abs(pos1[0]-target) {
		rig1, rig2 = rig2, rig1
		pos1, pos2 = pos2, pos1
		sta1, sta2 = sta2, sta1
		fin1, fin2 = fin2, fin1
	}
	for count < 10 {
		if (pos1[0]-target)*(pos2[0]-target) < 0.0 {
			rig3 := (rig1 + rig2) / 2.0
			sta3 := start
			for i := 3; i < 6; i += 1 {
				sta3[i] = start[i] * rig3 / 2.99792458 * 10.0
			}
			pos3, fin3 := Trace(sta3, dist, angle, h, bmap)
			count += 1
			//fmt.Println("DBGa: ", count, math.Abs(pos3[0]-target), rig1, rig2, rig3)
			if (pos3[0]-target)*(pos1[0]-target) < 0 {
				rig2 = rig3
				pos2 = pos3
				sta2 = sta3
				fin2 = fin3
			} else {
				rig1 = rig3
				pos1 = pos3
				sta1 = sta3
				fin1 = fin3
			}
		} else {
			rig3 := rig1 + (rig2-rig1)/(pos2[0]-pos1[0])*(target-pos1[0])
			sta3 := start
			for i := 3; i < 6; i += 1 {
				sta3[i] = start[i] * rig3 / 2.99792458 * 10.0
			}
			pos3, fin3 := Trace(sta3, dist, angle, h, bmap)
			count += 1
			//fmt.Println("DBGb: ", count, math.Abs(pos3[0]-target), rig1, rig2, rig3)
			rig1 = rig3
			pos1 = pos3
			sta1 = sta3
			fin1 = fin3
		}
		if math.Abs(pos2[0]-target) > math.Abs(pos1[0]-target) {
			rig1, rig2 = rig2, rig1
			pos1, pos2 = pos2, pos1
			sta1, sta2 = sta2, sta1
			fin1, fin2 = fin2, fin1
		}
		if math.Abs(pos2[0]-target) < maxerr {
			return rig2, pos2, count, fin2
		}
	}
	return rig2, pos2, count, fin2
}

func SmartFindTrace2(start [6]float64, dist float64, angle float64, h float64, target float64, bmap []float64) (float64, [4]float64, int, [6]float64) {
	rig1 := GuessRig(target)
	sta1 := start
	for i := 3; i < 6; i++ {
		sta1[i] = start[i] * rig1 / 2.99792458 * 10.0
	}
	pos1, fin1 := Trace(sta1, dist, angle, h, bmap)
	count := 1
	//fmt.Println("DBG0: ", count, math.Abs(pos1[0]-target))
	if math.Abs(pos1[0]-target) < maxerr {
		return rig1, pos1, count, fin1
	}
	//return rig1, pos1, count, fin1
	rig2 := math.FMA(target-pos1[0], C1, rig1)
	sta2 := start
	for i := 3; i < 6; i++ {
		sta2[i] = start[i] * rig2 / 2.99792458 * 10.0
	}
	pos2, fin2 := Trace(sta2, dist, angle, h, bmap)
	count = 2
	//fmt.Println("DBG1: ", count, math.Abs(pos2[0]-target))
	if math.Abs(pos2[0]-target) < maxerr {
		return rig2, pos2, count, fin2
	}
	//return rig2, pos2, count, fin2
	if math.Abs(pos2[0]-target) > math.Abs(pos1[0]-target) {
		rig1, rig2 = rig2, rig1
		pos1, pos2 = pos2, pos1
		sta1, sta2 = sta2, sta1
		fin1, fin2 = fin2, fin1
	}
	for count < 10 {
		//rig3 := rig1 + (rig2-rig1)/(pos2[0]-pos1[0])*(target-pos1[0])
		rig3 := math.FMA(target-pos1[0], (rig2-rig1)/(pos2[0]-pos1[0]), rig1)
		sta3 := start
		for i := 3; i < 6; i++ {
			sta3[i] = start[i] * rig3 / 2.99792458 * 10.0
		}
		pos3, fin3 := Trace(sta3, dist, angle, h, bmap)
		count += 1
		//fmt.Println("DBGb: ", count, math.Abs(pos3[0]-target), rig1, rig2, rig3)
		for math.Abs(pos3[0]-target) >= math.Abs(pos1[0]-target) {
			rig3 = (rig3 + rig2) * 0.5
			for i := 3; i < 6; i++ {
				sta3[i] = start[i] * rig3 / 2.99792458 * 10.0
			}
			pos3, fin3 = Trace(sta3, dist, angle, h, bmap)
			count += 1
			if count >= 10 {
				break
			}
			//fmt.Println("DBGc: ", count, math.Abs(pos3[0]-target), rig1, rig2, rig3)
		}
		rig1 = rig3
		pos1 = pos3
		sta1 = sta3
		fin1 = fin3
		if math.Abs(pos2[0]-target) > math.Abs(pos1[0]-target) {
			rig1, rig2 = rig2, rig1
			pos1, pos2 = pos2, pos1
			sta1, sta2 = sta2, sta1
			fin1, fin2 = fin2, fin1
		}
		if math.Abs(pos2[0]-target) < maxerr {
			return rig2, pos2, count, fin2
		}
	}
	return rig2, pos2, count, fin2
}
