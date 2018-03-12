package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

const blDevFilePath string = "/sys/class/backlight/intel_backlight/brightness"

func main() {
	switch argLen := len(os.Args); argLen {
	case 1:
		level := getBrightnessLevel()
		fmt.Println(level)
	case 2:
		switch incOrDecOrNum := os.Args[1]; incOrDecOrNum {
		case "inc":
			increaseBrightnessLevel()
		case "dec":
			decreaseBrightnessLevel()
		default:
			const bitSize = 32
			newLevel, err := strconv.ParseInt(incOrDecOrNum, 10, bitSize)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Invalid argument. (%s)\n", incOrDecOrNum)
				os.Exit(1)
			}
			setBrightnessLevel(uint32(newLevel))
		}
	default:
		fmt.Fprintln(os.Stderr, "Invalid argument. (%v)", os.Args)
		os.Exit(1)
	}
}

func getBrightnessLevel() uint32 {
	bytes, err := ioutil.ReadFile(blDevFilePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Could'nt read brightness device file")
		os.Exit(1)
	}
	// fmt.Printf("%v\n", bytes)
	// fmt.Printf("%s\n", bytes)
	// fmt.Printf("%d\n", bytes)

	// level := binary.BigEndian.Uint32(bytes)
	// fmt.Print(level)
	// 825241648
	valueAsText := strings.TrimRight(string(bytes), "\n")
	level, err := strconv.Atoi(valueAsText)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Couldn't convert from string to int\n")
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}
	return uint32(level)
}

func setBrightnessLevel(newLevel uint32) {
	info, err := os.Stat(blDevFilePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't get fileinfo (%v)\n", blDevFilePath)
		os.Exit(1)
	}

	err = ioutil.WriteFile(
		blDevFilePath,
		[]byte(strconv.Itoa(int(newLevel))+"\n"),
		info.Mode(),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't write to %v\n", blDevFilePath)
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func increaseBrightnessLevel() {
	level := getBrightnessLevel()
	setBrightnessLevel(level + 100)
}

func decreaseBrightnessLevel() {
	level := getBrightnessLevel()
	setBrightnessLevel(level - 100)
}
