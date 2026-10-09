package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/alecthomas/kong"
)

const kbdBlDevFilePath string = "/sys/class/leds/asus::kbd_backlight/brightness"

var blDevFilePath = findBlDevFilePath()

// Device name differs per GPU driver (intel_backlight, amdgpu_bl1, ...).
func findBlDevFilePath() string {
	paths, _ := filepath.Glob("/sys/class/backlight/*/brightness")
	if len(paths) == 0 {
		fmt.Fprintln(os.Stderr, "No backlight device found under /sys/class/backlight")
		os.Exit(1)
	}
	return paths[0]
}

var cli struct {
	Level  levelCmd  `cmd:"" default:"withargs" help:"Print the brightness level, or set it when LEVEL is given."`
	Inc    incCmd    `cmd:"" help:"Increase the brightness level."`
	Dec    decCmd    `cmd:"" help:"Decrease the brightness level."`
	KbdOn  kbdOnCmd  `cmd:"" name:"kbd-on" help:"Turn on the keyboard backlight."`
	KbdOff kbdOffCmd `cmd:"" name:"kbd-off" help:"Turn off the keyboard backlight."`
}

type levelCmd struct {
	Level *uint32 `arg:"" optional:""`
}

func (c *levelCmd) Run() error {
	if c.Level == nil {
		fmt.Println(getBrightnessLevel())
		return nil
	}
	setBrightnessLevel(*c.Level)
	return nil
}

type incCmd struct{}

func (incCmd) Run() error { increaseBrightnessLevel(); return nil }

type decCmd struct{}

func (decCmd) Run() error { decreaseBrightnessLevel(); return nil }

type kbdOnCmd struct{}

func (kbdOnCmd) Run() error { setKbdBackLight(kbdBlOn); return nil }

type kbdOffCmd struct{}

func (kbdOffCmd) Run() error { setKbdBackLight(kbdBlOff); return nil }

func main() {
	ctx := kong.Parse(&cli, kong.Description("Backlight brightness control."))
	ctx.FatalIfErrorf(ctx.Run())
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
	err := os.WriteFile(sysfile, state.bytes(), 0o644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't write to %v\n", sysfile)
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func getBrightnessLevel() uint32 {
	return readSysUint(blDevFilePath)
}

func getMaxBrightnessLevel() uint32 {
	return readSysUint(filepath.Join(filepath.Dir(blDevFilePath), "max_brightness"))
}

func readSysUint(path string) uint32 {
	bytes, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't read %v\n", path)
		os.Exit(1)
	}
	level, err := strconv.ParseUint(strings.TrimSpace(string(bytes)), 10, 32)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't parse %v: %v\n", path, err)
		os.Exit(1)
	}
	return uint32(level)
}

func setBrightnessLevel(newLevel uint32) {
	err := os.WriteFile(blDevFilePath, []byte(strconv.Itoa(int(newLevel))+"\n"), 0o644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't write to %v\n", blDevFilePath)
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

// max_brightness ranges from ~100 (intel) to 65535 (amdgpu), so a fixed step can't fit all devices.
const stepDivisor = 20

func stepLevel(level, maxLevel uint32, up bool) uint32 {
	step := max(maxLevel/stepDivisor, 1)
	if up {
		return min(level+step, maxLevel)
	}
	if level < step {
		return 0
	}
	return level - step
}

func increaseBrightnessLevel() {
	setBrightnessLevel(stepLevel(getBrightnessLevel(), getMaxBrightnessLevel(), true))
}

func decreaseBrightnessLevel() {
	setBrightnessLevel(stepLevel(getBrightnessLevel(), getMaxBrightnessLevel(), false))
}
