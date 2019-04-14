package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

const kbdBlDevFilePath string = "/sys/class/leds/asus::kbd_backlight/brightness"

const blDevFilePath string = "/sys/class/backlight/intel_backlight/brightness"

func main() {
	switch argLen := len(os.Args); argLen {
	case 1:
		level := getBrightnessLevel()
		fmt.Println(level)
	case 2:
		switch operationType := os.Args[1]; operationType {
		case "inc":
			increaseBrightnessLevel()
		case "dec":
			decreaseBrightnessLevel()
		case "kbd-on":
			setKbdBackLight(kbdBlOn)
		case "kbd-off":
			setKbdBackLight(kbdBlOff)
		default:
			const bitSize = 32
			newLevel, err := strconv.ParseInt(operationType, 10, bitSize)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Invalid argument. (%s)\n", operationType)
				os.Exit(1)
			}
			setBrightnessLevel(uint32(newLevel))
		}
	default:
		fmt.Fprintf(os.Stderr, "Invalid argument. (%v)\n", os.Args)
		os.Exit(1)
	}
}

type backlightState uint8

const (
	kbdBlOn backlightState = iota
	kbdBlOff
)

func (state backlightState) bytes() []byte {
	switch state {
	case kbdBlOn:
		return []byte("1\n")
	case kbdBlOff:
		return []byte("0\n")
	default:
		log.Fatalf("invalid backlightState")
		return []byte("\n")
	}
}

func setKbdBackLight(state backlightState) {
	sysfile := kbdBlDevFilePath
	info, err := os.Stat(sysfile)
	err = ioutil.WriteFile(sysfile, state.bytes(), info.Mode())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't write to %v\n", sysfile)
		fmt.Fprintln(os.Stderr, err.Error())
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
		fmt.Fprintf(os.Stderr, "Couldn't convert from string to int\n")
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
