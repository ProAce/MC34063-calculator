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
		//Start a reader from the terminal
		reader := bufio.NewReader(os.Stdin)

		//Enter and parse all the variables
		fmt.Print("Enter Vin (V): ")
		input, _ := reader.ReadString('\n')
		//Check wether it is a Unix or a Windows terminal
		if strings.HasSuffix(input, "\r\n") {
			check = "\r\n"
		} else {
			check = "\n"
		}
		input = strings.TrimSuffix(input, check)
		vin, _ = strconv.ParseFloat(input, 64)

		vinperc = termInput("Enter Vin delta (%): ", reader, check)

		vout = termInput("Enter Vout (V): ", reader, check)

		ripple = termInput("Enter ripple (mV): ", reader, check)
		ripple /= 1000

		i = termInput("Enter Iout (mA): ", reader, check)
		i /= 1000

		frequency = termInput("Enter frequency (KHz): ", reader, check)
		frequency *= 1000

		vf = termInput("Enter forward voltage diode (V): ", reader, check)

		vinmin = (vin - (vin * (vinperc / 100)))

		fmt.Println()

		//Check if it should be in boost or buck mode and start the corresponding program
		if vin > vout {
			vsat = 1
			buck()
		} else {
			vsat = .45
			boost()
		}

		//Loop till the users tells the program to start again or to quit
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
	//Boost calculations based on the MC34063 datasheet
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
	//Buck calculations based on the MC34063 datasheet
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

//Round all values to two decimals
func round(x float64) float64 {
	return math.Round(x/.01) * .01
}

//Parse the input from the terminal from a string to a float
func termInput(input string, reader *bufio.Reader, check string) float64 {
	fmt.Print(input)
	in, _ := reader.ReadString('\n')
	in = strings.TrimSuffix(in, check)
	out, _ := strconv.ParseFloat(in, 64)
	return out
}
