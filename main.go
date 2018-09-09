package main

import (
	"chip8emu/core"
	"chip8emu/opts"
	"chip8emu/view"
	"fmt"
	"sync"

	"github.com/jessevdk/go-flags"
)

var wg = sync.WaitGroup{}

var chip *core.Chip8

func main() {

	opts := opts.Opts{}
	_, err := flags.Parse(&opts)

	if err != nil {
		panic(err)
	}

	chip = core.NewChip8()

	//chip.GfxClipping = core.DrawClippingEnabled

	view.NewSDLDisplayRenderer(chip, &wg, &opts)
	fmt.Printf("FILE: %v\n", opts.File)
	chip.Load(opts.File)
	chip.Start()
	wg.Wait()
}

func GetChip() *core.Chip8 {
	return chip
}
