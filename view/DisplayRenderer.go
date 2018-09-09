package view

import (
	"chip8emu/core"
	"chip8emu/opts"
	"runtime"
	"sync"
	"time"

	"gopkg.in/veandco/go-sdl2.v0/sdl"
)

const (
	PixelWidth   = 10
	PixelHeight  = 10
	ScreenWidth  = 64 * PixelWidth
	ScreenHeight = 32 * PixelHeight
)

var keyboard2Chip8 = map[sdl.Keycode]uint8{
	sdl.K_1: core.ChipKey1,
	sdl.K_2: core.ChipKey2,
	sdl.K_3: core.ChipKey3,
	sdl.K_4: core.ChipKeyC,
	sdl.K_q: core.ChipKey4,
	sdl.K_w: core.ChipKey5,
	sdl.K_e: core.ChipKey6,
	sdl.K_r: core.ChipKeyD,
	sdl.K_a: core.ChipKey7,
	sdl.K_s: core.ChipKey8,
	sdl.K_d: core.ChipKey9,
	sdl.K_f: core.ChipKeyE,
	sdl.K_z: core.ChipKeyA,
	sdl.K_x: core.ChipKey0,
	sdl.K_c: core.ChipKeyB,
	sdl.K_v: core.ChipKeyF,
}

type SDLDisplayRenderer struct {
	Alive     bool
	SdlWindow *sdl.Window
	cpu       *core.Chip8
	vmem      []uint8
	ExitChan  chan<- int
	WaitGroup *sync.WaitGroup
	bgColour  uint32
	xSize     int
	ySize     int
	xDisplay  int
	yDisplay  int
	pixel     *sdl.Rect
}

func NewSDLDisplayRenderer(cpu *core.Chip8, wg *sync.WaitGroup, opts *opts.Opts) *SDLDisplayRenderer {
	runtime.LockOSThread()
	renderer := &SDLDisplayRenderer{}
	renderer.WaitGroup = wg

	renderer.Alive = true
	renderer.WaitGroup.Add(1)
	renderer.bgColour = opts.BgColour

	if renderer.xSize = opts.XSize; renderer.xSize == 0 {
		renderer.xSize = PixelWidth
	}

	if renderer.ySize = opts.YSize; renderer.ySize == 0 {
		renderer.ySize = PixelHeight
	}

	if renderer.xDisplay = 64 * renderer.xSize; renderer.xDisplay == 0 {
		renderer.xDisplay = ScreenWidth
	}

	if renderer.yDisplay = 32 * renderer.ySize; renderer.yDisplay == 0 {
		renderer.yDisplay = ScreenHeight
	}

	renderer.pixel = &sdl.Rect{X: 0, Y: 0, W: int32(renderer.xSize), H: int32(renderer.ySize)}

	renderer.Init(cpu)

	go renderer.Render()
	return renderer
}

func (r *SDLDisplayRenderer) Init(cpu *core.Chip8) {

	r.cpu = cpu
	r.cpu.DisplayHandler = r
	r.InitDisplay()
	r.vmem = cpu.Memory[3840:]

}

func (r *SDLDisplayRenderer) Render() {
	refresh := time.NewTicker(time.Second / 50)

	for r.Alive {

		select {
		case <-refresh.C:
			surface, _ := r.SdlWindow.GetSurface()

			rect := sdl.Rect{X: 0, Y: 0, W: int32(r.xDisplay), H: int32(r.yDisplay)}
			surface.FillRect(&rect, r.bgColour) //0xff00aa00)

			//TODO: Do some timing regarding creating a new rectangle for every pixel vs
			//one or a pool of rectangles reused.

			for y, v1 := range r.cpu.VMem {
				for x, v2 := range v1 {
					if v2 > 0 {
						r.pixel.X = int32(x * int(r.xSize))
						r.pixel.Y = int32(y * int(r.ySize))
						surface.FillRect(r.pixel, uint32(0xffffffff))
					}
				}
			}

			r.SdlWindow.UpdateSurface()

		}
	}

	r.WaitGroup.Done()
}

func (r *SDLDisplayRenderer) InitDisplay() {

	sdl.Init(sdl.INIT_EVERYTHING)
	window, err := sdl.CreateWindow("Chip8", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		r.xDisplay, r.yDisplay, sdl.WINDOW_SHOWN)

	r.SdlWindow = window

	if err != nil {
		panic(err)
	}

	surface, _ := r.SdlWindow.GetSurface()
	rect := sdl.Rect{0, 0, int32(r.xDisplay), int32(r.yDisplay)}
	surface.FillRect(&rect, 0xffff0000)

	r.SdlWindow.UpdateSurface()
	//sdl.Delay(4000)

}

func (r *SDLDisplayRenderer) Shutdown() {
	r.Alive = false
	r.WaitGroup.Wait()
	r.SdlWindow.Destroy()
	sdl.Quit()

}

func (r *SDLDisplayRenderer) IsAlive() bool {
	return r.Alive
}

func (r *SDLDisplayRenderer) WaitForInput() uint8 {
	event := sdl.WaitEvent()

	switch e := event.(type) {
	case *sdl.KeyDownEvent:
		return r.handleAwaitKeyPress(*e)
	}

	return 0
}

func (r *SDLDisplayRenderer) HandleInput() {

	event := sdl.PollEvent()

	if event != nil {
		r.handleKeyPress(event)
	}

}

func (r *SDLDisplayRenderer) handleAwaitKeyPress(event sdl.Event) uint8 {
	chipKey := keyboard2Chip8[event.(sdl.KeyDownEvent).Keysym.Sym]
	return chipKey
}

func (r *SDLDisplayRenderer) handleKeyPress(event sdl.Event) {
	r.doKeyPress(event)
}

func (r *SDLDisplayRenderer) doKeyPress(event sdl.Event) {

	switch e := event.(type) {
	case *sdl.KeyDownEvent:
		chipKey, exists := keyboard2Chip8[e.Keysym.Sym]
		if exists {
			r.cpu.SetKey(chipKey)
		}

		if e.Keysym.Sym == sdl.K_F1 {
			print("Activate")
			r.cpu.MM.Activate()
		}

		if e.Keysym.Sym == sdl.K_F2 {
			print("Deactivate")
			r.cpu.MM.Deactivate()
		}

		if e.Keysym.Sym == sdl.K_F3 {
			print("Run Step")
			r.cpu.MM.SetRunStep()
		}

		break
	case *sdl.KeyUpEvent:
		chipKey, exists := keyboard2Chip8[e.Keysym.Sym]
		if exists {
			r.cpu.ClrKey(chipKey)
		}
		break
	case *sdl.QuitEvent:
		r.Shutdown()
	}

}
