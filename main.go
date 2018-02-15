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
	level := getBrightnessLevel()
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
	//ioutil.WriteFile(blDevFilePath,

	// func WriteFile(filename string, data []byte, perm os.FileMode) error {
	err := ioutil.WriteFile(blDevFilePath, []byte(string(newLevel)), "w")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Couldn't write")
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
