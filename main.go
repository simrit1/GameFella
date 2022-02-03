package main

import (
	"github.com/is386/GoBoy/emu"
)

var (
	DEBUG = false
	ROM   = "roms/5.gb"
)

func main() {
	gb := emu.NewGameBoy(DEBUG)
	gb.Run(ROM)
}
