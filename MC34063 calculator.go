package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

var vin, vinperc, vinmin, vout, ripple, i, vf, vsat, frequency float64

func main() {
	retry := true
	var check string

	for retry == true {
		reader := bufio.NewReader(os.Stdin)

		fmt.Print("Enter Vin (V): ")
		input, _ := reader.ReadString('\n')
		if strings.HasSuffix(input, "\r\n") {
			check = "\r\n"
		} else {
			check = "\n"
		}
		input = strings.TrimSuffix(input, check)
		vin, _ = strconv.ParseFloat(input, 64)

		fmt.Print("Enter Vin delta (%): ")
		input, _ = reader.ReadString('\n')
		input = strings.TrimSuffix(input, check)
		vinperc, _ = strconv.ParseFloat(input, 64)

		fmt.Print("Enter Vout (V): ")
		input, _ = reader.ReadString('\n')
		input = strings.TrimSuffix(input, check)
		vout, _ = strconv.ParseFloat(input, 64)

		fmt.Print("Enter ripple (mV): ")
		input, _ = reader.ReadString('\n')
		input = strings.TrimSuffix(input, check)
		ripple, _ = strconv.ParseFloat(input, 64)
		ripple /= 1000

		fmt.Print("Enter Iout (mA): ")
		input, _ = reader.ReadString('\n')
		input = strings.TrimSuffix(input, check)
		i, _ = strconv.ParseFloat(input, 64)
		i /= 1000

		fmt.Print("Enter frequency (KHz): ")
		input, _ = reader.ReadString('\n')
		input = strings.TrimSuffix(input, check)
		frequency, _ = strconv.ParseFloat(input, 64)
		frequency *= 1000

		fmt.Print("Enter forward voltage diode (V): ")
		input, _ = reader.ReadString('\n')
		input = strings.TrimSuffix(input, check)
		vf, _ = strconv.ParseFloat(input, 64)

		vinmin = (vin - (vin * (vinperc / 100)))

		fmt.Println()

		if vin > vout {
			vsat = 1
			buck()
		} else {
			vsat = .45
			boost()
		}

		br := false
		for br == false {
			br = true
			retry = false

			fmt.Println()
			fmt.Print("Want to do another calculation? (Y/N)")
			input, _ = reader.ReadString('\n')
			input = strings.TrimSuffix(input, check)
			fmt.Println()
			if (input == "y") || (input == "Y") {
				retry = true
			} else if (input == "n") || (input == "N") {
			} else {
				fmt.Print("False entry, enter Y or N")
				fmt.Println()
				br = false
			}
		}
	}
}

func boost() {
	TonToff := (vout + vf - vinmin) / (vinmin - vsat)
	time := 1 / frequency
	toff := time / (TonToff + 1)
	ton := time - toff
	ct := (4.0 * math.Pow(10, -5)) * ton
	ipk := 2 * i * (TonToff + 1)
	rsc := .3 / ipk
	lmin := ((vinmin - vsat) / ipk) * ton
	co := 9 * ((i * ton) / ripple)

	co *= math.Pow(10, 6)
	ct *= math.Pow(10, 12)
	ton *= math.Pow(10, 6)
	lmin *= math.Pow(10, 6)
	ipk *= math.Pow(10, 3)

	fmt.Println("Vout:", round(vout), "V")
	fmt.Println("Vin(min):", round(vinmin), "V")
	fmt.Println("Ton:", round(ton), "us")
	fmt.Println("Ct:", round(ct), "pF")
	fmt.Println("Ipeak:", round(ipk), "mA")
	fmt.Println("Rsc:", round(rsc), "Ohm")
	fmt.Println("L(min):", round(lmin), "uH")
	fmt.Println("Co:", round(co), "uF")
}

func buck() {
	TonToff := (vout + vf) / (vinmin - vsat - vout)
	time := (1 / frequency)
	toff := time / (TonToff + 1)
	ton := time - toff
	ct := (4.0 * math.Pow(10, -5)) * ton
	ipk := 2 * i
	rsc := .3 / ipk
	lmin := ((vinmin - vsat - vout) / ipk) * ton
	co := (ipk * time) / (8 * ripple)

	co *= math.Pow(10, 6)
	ct *= math.Pow(10, 12)
	ton *= math.Pow(10, 6)
	lmin *= math.Pow(10, 6)
	ipk *= math.Pow(10, 3)

	fmt.Println("Vout:", round(vout), "V")
	fmt.Println("Vin(min):", round(vinmin), "V")
	fmt.Println("Ton:", round(ton), "us")
	fmt.Println("Ct:", round(ct), "pF")
	fmt.Println("Ipeak:", round(ipk), "mA")
	fmt.Println("Rsc:", round(rsc), "Ohm")
	fmt.Println("L(min):", round(lmin), "uH")
	fmt.Println("Co:", round(co), "uF")
}

func round(x float64) float64 {
	return math.Round(x/.01) * .01
}
