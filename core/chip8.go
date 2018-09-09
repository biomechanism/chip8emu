package core

// typedef unsigned char Uint8;
// void MyCallback(void *userdata, Uint8 *stream, int len);
// void Chip8AudioCallback(void *userdata, Uint8 *stream, int len);
import "C"
import (
	"chip8emu/utils"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"reflect"
	"time"
	"unsafe"

	"gopkg.in/veandco/go-sdl2.v0/sdl"
)

var AudioPlay bool

var chip8 *Chip8

const (
	displayWidth  = 64
	displayHeight = 32
)

//Chip-8 keypad keys.
const (
	ChipKey0 = 0
	ChipKey1 = 1
	ChipKey2 = 2
	ChipKey3 = 3
	ChipKey4 = 4
	ChipKey5 = 5
	ChipKey6 = 6
	ChipKey7 = 7
	ChipKey8 = 8
	ChipKey9 = 9
	ChipKeyA = 10
	ChipKeyB = 11
	ChipKeyC = 12
	ChipKeyD = 13
	ChipKeyE = 14
	ChipKeyF = 15
)

//Chip-8 registers
const (
	V0 = 0
	V1 = 1
	V2 = 2
	V3 = 3
	V4 = 4
	V5 = 5
	V6 = 6
	V7 = 7
	V8 = 8
	V9 = 9
	VA = 10
	VB = 11
	VC = 12
	VD = 13
	VE = 14
	VF = 15
)

//Offsets for the character banks
const (
	CharBank0  = 0x0
	CharBank1  = 0x5
	CharBank2  = 0xA
	CharBank3  = 0xF
	CharBank4  = 0x14
	CharBank5  = 0x19
	CharBank6  = 0x1E
	CharBank7  = 0x23
	CharBank8  = 0x28
	CharBank9  = 0x2D
	CharBank10 = 0x32
	CharBank11 = 0x37
	CharBank12 = 0x3C
	CharBank13 = 0x41
	CharBank14 = 0x46
	CharBank15 = 0x4B
)

//sdl audio generate waveform

//Stores memory addresses associated with each character bank.
var charBank = [16]uint16{
	CharBank0, CharBank1, CharBank2, CharBank3, CharBank4, CharBank5, CharBank6, CharBank7,
	CharBank8, CharBank9, CharBank10, CharBank11, CharBank12, CharBank13, CharBank14, CharBank15,
}

var chars = []uint8{
	0xF0, 0x90, 0x90, 0x90, 0xF0, 0x20, 0x60, 0x20, 0x20, 0x70, 0xF0, 0x10, 0xF0, 0x80, 0xF0, 0xF0, 0x10, 0xF0, 0x10, 0xF0,
	0x90, 0x90, 0xF0, 0x10, 0x10, 0xF0, 0x80, 0xF0, 0x10, 0xF0, 0xF0, 0x80, 0xF0, 0x90, 0xF0, 0xF0, 0x10, 0x20, 0x40, 0x40,
	0xF0, 0x90, 0xF0, 0x90, 0xF0, 0xF0, 0x90, 0xF0, 0x10, 0xF0, 0xF0, 0x90, 0xF0, 0x90, 0x90, 0xE0, 0x90, 0xE0, 0x90, 0xE0,
	0xF0, 0x80, 0x80, 0x80, 0xF0, 0xE0, 0x90, 0x90, 0x90, 0xE0, 0xF0, 0x80, 0xF0, 0x80, 0xF0, 0xF0, 0x80, 0xF0, 0x80, 0x80,
}

type handlerTable struct {
	InstructionTable *map[int]func(*Chip8, uint16)
}

//GetHandler used to retrieve the relevant instruction handler for the given opcode
func (t *handlerTable) GetHandler(opcode uint16) func(*Chip8, uint16) {
	var mask uint16
	switch opcode & 0xF000 {
	case 0x8000:
		mask = 0xF00F
		break
	case 0x0000:
		mask = 0xF0FF
		break
	case 0xE000:
		mask = 0xF0FF
		break
	case 0xF000:
		mask = 0xF0FF
		break
	default:
		mask = 0xF000
	}

	return (*t.InstructionTable)[int(opcode&mask)]
}

func newHandlerTable() *map[int]func(*Chip8, uint16) {
	var InstructionTable = map[int]func(*Chip8, uint16){
		Clear:                Handle0x0,
		Return:               Handle0x0,
		Jump:                 Handle0x1,
		JumpSub:              Handle0x2,
		SkipVxEqKk:           Handle0x3,
		SkipVxNeqKk:          Handle0x4,
		SkipVxEqVy:           Handle0x5,
		LoadVxFromKk:         Handle0x6,
		LoadVxAddKk:          Handle0x7,
		LoadVxFromVy:         Handle0x8,
		LoadVxOrVy:           Handle0x8,
		LoadVxAndVy:          Handle0x8,
		LoadVxXorVy:          Handle0x8,
		LoadVxAddVy:          Handle0x8,
		LoadVxSubVy:          Handle0x8,
		LoadVxShiftR:         Handle0x8,
		LoadVxVySubVx:        Handle0x8,
		LoadVxShiftL:         Handle0x8,
		SkipVxNeqVy:          Handle0x9,
		LoadIFromNnn:         Handle0xA,
		JumpPlusA0:           Handle0xB,
		RndVxAndKk:           Handle0xC,
		DrawSprite:           Handle0xD,
		SkipVxEqKey:          Handle0xE,
		SkipVxNeqKey:         Handle0xE,
		LoadVxFromK:          Handle0xE,
		LoadVxFromDelayTimer: Handle0xF,
		LoadDelayTimerFromVx: Handle0xF,
		LoadSoundTimerFromVx: Handle0xF,
		LoadIAddVx:           Handle0xF,
		LoadSpriteCharacter:  Handle0xF,
		LoadIWithBcdOfVx:     Handle0xF,
		StoreV0ToVx:          Handle0xF,
		LoadIToV0ToVx:        Handle0xF,
	}
	return &InstructionTable
}

type inputHandler struct {
	InputInterface
}

//DisplayInterface the expected methods that a renderer must implement
type DisplayInterface interface {
	Render()
	IsAlive() bool
	HandleInput()
	WaitForInput() uint8
}

//InputInterface the expected method that an input handler must implement
type InputInterface interface {
	Read() uint16
}

//Chip8 the core structure for managing machine state
type Chip8 struct {
	Memory [4096]uint8
	V      [16]uint8
	VMem   [32][64]uint8

	//Address register
	I uint16

	//Program Counter
	Pc uint16
	//Stack
	S [16]uint16
	//Stack Pointer
	Sp uint8

	DelayTimer uint8
	SoundTimer uint8

	DelayTicker *time.Ticker
	SoundTicker *time.Ticker

	Keys [16]uint8

	DTDisabled bool
	STDisabled bool

	DisplayHandler   DisplayInterface
	KBHandler        *inputHandler
	InstHandlerTable *handlerTable
	GfxClipping      bool
	MM               *utils.MachineMonitor
}

//NewChip8 constructor to instantiate a machine
func NewChip8() (c *Chip8) {
	chip := new(Chip8)
	chip8 = chip
	chip.init()
	return chip
}

//NewChip8WithHandler constructor to instantiate a machine with a specified input handler
func NewChip8WithHandler(KBHandler *inputHandler) (c *Chip8) {
	chip := NewChip8()
	chip.KBHandler = KBHandler
	return chip
}

func (c *Chip8) init() {
	c.Sp = 15
	c.Pc = 0x200
	var address uint16
	var i uint16
	for i = 0; i < 16; i++ {
		charBank[i] = address
		address += 5
	}

	for i, e := range chars {
		c.Memory[i] = e
	}

	//c.VMem = c.Memory[3840:]
	c.InstHandlerTable = &handlerTable{}
	c.InstHandlerTable.InstructionTable = newHandlerTable()
	//c.GfxClipping = DrawClippingDisabled
	c.MM = utils.NewMachineMonitor()
	initBeep()
}

//LoadROM load a ROM provided as sa base64 encoded string.
func (c *Chip8) LoadROM(rom string) {
	data, err := base64.StdEncoding.DecodeString(rom)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	for i := range data {
		c.Memory[0x200+i] = data[i]
	}
}

//startDelayTimer start delay time with number.
//Thi number will tick down at a rate of 60Hz
func (c *Chip8) startDelayTimer(num uint8) {

	c.DelayTimer = num

	if c.DTDisabled {
		return
	}

	if c.DelayTicker != nil {
		c.DelayTicker.Stop()
	}

	c.DelayTicker = time.NewTicker(time.Second / 60)
	go func() {
		for range c.DelayTicker.C {
			c.DelayTimer--

			if c.DelayTimer == 0 {
				break
			}
		}
		c.DelayTicker.Stop()
		c.DelayTicker = nil
	}()

}

//startSoundTimer start sound timer with specified number.
//This number will tick down at 60Hz
func (c *Chip8) startSoundTimer(num uint8) {

	c.SoundTimer = num

	if c.STDisabled {
		return
	}

	if c.SoundTicker != nil {
		c.SoundTicker.Stop()
	}

	c.SoundTicker = time.NewTicker(time.Second / 60)
	go func() {

		for range c.SoundTicker.C {

			if c.SoundTimer == 0 {
				continue
			}

			if c.SoundTimer > 0 {
				playBeep()
				c.SoundTimer--
			}

		}

	}()

}

//Start start the machine running
func (c *Chip8) Start() {
	alive := c.DisplayHandler.IsAlive()
	for alive {
		if c.doFDECycle() {
			ins := c.fetch()
			//FIXME: Allow error code to be returned from execute, so we can stop the processor
			//FDE cycle if an invalid instruction is encountered.
			c.execute(ins)
			c.MM.Reset()
		}
		c.DisplayHandler.HandleInput()
		alive = c.DisplayHandler.IsAlive()
		sdl.Delay(2)
	}
}

func (c *Chip8) doFDECycle() bool {

	if !c.MM.IsActive() {
		if c.MM.IsBP(c.GetPc()) {
			c.MM.Activate()
			return false
		}
		return true
	}

	if c.MM.IsActive() {
		if c.MM.IsRunStep() {
			return true
		}

		return false
	}

	return true
}

//Load loads a ROM into memory.
//TODO: change this so it doesn't make use of the base64 coversion method, as
//it is doing needess back and forth coversion to load to memory.
func (c *Chip8) Load(filePath string) {

	var mem []byte
	var err error
	if mem, err = ioutil.ReadFile(filePath); err != nil {
		log.Fatal("Unable to load ROM file")
	}

	rom := base64.StdEncoding.EncodeToString(mem)

	c.LoadROM(rom)
}

func (c *Chip8) fetch() (opcode uint16) {
	var op = uint16(c.Memory[c.Pc]) << 8
	op |= uint16(c.Memory[c.Pc+1])
	c.SetPc(c.Pc + 2)
	return op
}

func (c *Chip8) execute(inst uint16) {
	h := c.InstHandlerTable.GetHandler(inst)
	h(c, inst)
}

//SetMem sets the specified memory location to the passed value
func (c *Chip8) SetMem(mem, val uint16) {
	c.Memory[mem] = uint8(val >> 8)
	c.Memory[mem+1] = uint8(val & 0xFF)
}

//SetV set the specified register to the passed value
func (c *Chip8) SetV(reg, val uint8) {
	c.V[reg] = val
}

//SetI sets the memory address regiter to the provided address
func (c *Chip8) SetI(val uint16) {
	c.I = val
}

//GetI get the current address pointed to by the address register
func (c *Chip8) GetI() uint16 {
	return c.I
}

//GetV gets the value of the specified register
func (c *Chip8) GetV(reg uint8) uint8 {
	return c.V[reg]
}

//SetPc sets the Program Counter to the specified address
func (c *Chip8) SetPc(val uint16) {
	c.Pc = val
}

//GetPc gets the current address in the Program Counter
func (c *Chip8) GetPc() uint16 {
	return c.Pc
}

//Push pushes the specified value onto the stack
func (c *Chip8) Push(val uint16) {
	c.S[c.Sp] = val
	c.Sp--
}

//Pop pops the current value off the stack
func (c *Chip8) Pop() uint16 {
	c.Sp++
	return c.S[c.Sp]
}

//SetKey set the specified key as pressed
func (c *Chip8) SetKey(key uint8) {
	c.Keys[key] = 1
}

//ClrKey sets the specified key as released
func (c *Chip8) ClrKey(key uint8) {
	c.Keys[key] = 0
}

//GetKey gets the current state of the key, 1 == pressed, 0 == released
func (c *Chip8) GetKey(key uint8) uint8 {
	return c.Keys[key]
}

//GetST gets the Sound Timer
func (c *Chip8) GetST() uint8 {
	return c.SoundTimer
}

//SetST initialises the sound timer with the specified value
func (c *Chip8) SetST(val uint8) {
	if c.SoundTicker == nil {
		c.startSoundTimer(val)
		return
	}

	c.SoundTimer = val
}

//SetDT initialises the delay timer with the specified value
func (c *Chip8) SetDT(val uint8) {
	c.startDelayTimer(val)
}

//GetDT retrieve the delay timer value
func (c *Chip8) GetDT() uint8 {
	return c.DelayTimer
}

//DisableDelayTimer disables the delay timer from ticking down
func (c *Chip8) DisableDelayTimer() {
	c.DTDisabled = true
}

//EnableDelayTiemr enables the delay timer for ticking down
func (c *Chip8) EnableDelayTiemr() {
	c.DTDisabled = false
}

//DisableSoundTimer disables the sound timer from ticking down
func (c *Chip8) DisableSoundTimer() {
	c.STDisabled = true
}

//EnableSoundTiemr enables the sound timer for ticking down
func (c *Chip8) EnableSoundTiemr() {
	c.STDisabled = false
}

//GetRegVx get the Vx register number form the opcode
func GetRegVx(opcode uint16) uint8 {
	return uint8((opcode & 0x0F00) >> 8)
}

//GetRegVy get the Vy regiter number from the opcode
func GetRegVy(opcode uint16) uint8 {
	return uint8((opcode & 0x00F0) >> 4)
}

//GetOpVal get the one byte value from the opcode
func GetOpVal(opcode uint16) uint8 {
	return uint8(opcode & 0x00FF)
}

//GetOpVal12 get the 12 bit value from the opcode
func GetOpVal12(opcode uint16) uint16 {
	return uint16(0x0FFF & opcode)
}

//GetCharBank get the character bank address for the specified number
func GetCharBank(n uint16) uint16 {
	return charBank[n]
}

//Draw sprite to screen at coordinates x,y of hight n
func (c *Chip8) Draw(n uint8, x, y int32) {
	c.WriteVMem(int(x), int(y), int(n))
}

//WriteVMem write sprite to video memory
func (c *Chip8) WriteVMem(x, y, n int) {
	c.SetV(VF, 0)
	iReg := c.GetI()

	for i := 0; i < n; i++ {

		if !(y < displayHeight && y >= 0) {
			continue
		}

		data := c.Memory[iReg]
		dx := x

		for j := 0; j < 8; j++ {

			if !(dx < displayWidth && dx >= 0) {
				break
			}

			on := (0x80 & data) > 1
			data <<= 1
			if on {
				if c.VMem[y][dx] > 0 {
					c.SetV(VF, 1)
				}
				c.VMem[y][dx] ^= 0xFF
			}
			dx++
		}
		iReg++
		y++
	}
}

//ClearScreenMem used to clear the video memory
func (c *Chip8) ClearScreenMem() {
	vmem := c.Memory[3840:]
	for i := 0; i < 256; i++ {
		vmem[i] = 0
	}

	for y := 0; y < 32; y++ {
		for x := 0; x < 64; x++ {
			c.VMem[y][x] = 0
		}
	}
}

func initBeep() {
	sdl.PauseAudio(true)
	desired := &sdl.AudioSpec{}
	desired.Freq = 44100
	desired.Format = sdl.AUDIO_S8
	desired.Channels = 1
	desired.Samples = 2048
	//desired.Size = 176400
	desired.Callback = sdl.AudioCallback(C.MyCallback)
	sdl.OpenAudio(desired, nil)

}

func playBeep() {
	sdl.PauseAudio(false)
}

func stopBeep() {
	sdl.PauseAudio(true)
}

// func callback() {
// 	ass := C.int
// 	D := C.Uint8
// }

const (
	toneHz   = 300
	sampleHz = 44100

//	dPhase   = 2 * math.Pi * toneHz / sampleHz
)

//MyCallback foo
//export MyCallback
func MyCallback(userdata unsafe.Pointer, stream *C.Uint8, length C.int) {
	//	toneHz := 440
	//sampleHz := 48000
	dPhase := 2 * math.Pi * toneHz / sampleHz
	n := int(length)

	fmt.Printf("BUFFER LEN: %d\n", n)

	hdr := reflect.SliceHeader{Data: uintptr(unsafe.Pointer(stream)), Len: n, Cap: n}
	buf := *(*[]C.Uint8)(unsafe.Pointer(&hdr))

	var phase float64
	for i := 0; i < n; i += 2 {
		phase += dPhase
		sample := C.Uint8((math.Sin(phase) + 0.999999) * 128)
		buf[i] = sample
		buf[i+1] = sample
		// buf[i+2] = sample
		// buf[i+3] = sample
	}
	sdl.PauseAudio(true)
}

//export Chip8AudioCallback
func Chip8AudioCallback(userdata unsafe.Pointer, stream *C.Uint8, length C.int) {
	//y(t) = A * sin(2*PI*frequency*time + phase)
	sampleFrequency := 44100.0
	amplitude := 100.0
	//phase := 0.0
	time := 1.0 / sampleFrequency
	//bufferSize := int(frequency) * 4
	//ticker := 1.0 / 60

	//dPhase := 2 * math.Pi * toneHz / sampleHz

	//chip := main.GetChip()
	//n := int((chip8.GetST() * 16))
	//n += int(n * n)
	n := int(length)
	//	maxWrite := int(ticker / time)

	fmt.Printf("LEN: %d\n", n)

	hdr := reflect.SliceHeader{Data: uintptr(unsafe.Pointer(stream)), Len: n, Cap: n}
	buf := *(*[]C.Uint8)(unsafe.Pointer(&hdr))

	//var phase float64
	//st := chip8.GetST()
	//var writeCnt = uint8(frequency/60.0) * st
	var t float64
	for i := 0; i < n; i += 2 {
		//phase += dPhase
		//sample := C.Uint8((math.Sin(phase) + 0.2) * 128)
		y := math.Sin(2*math.Pi*440/sampleFrequency*time + 0.0)
		sample := C.Uint8(amplitude * y)
		t += time
		//if writeCnt < maxWrite {
		buf[i] = sample
		buf[i+1] = sample
		//}

		//buf[i+2] = sample
		//buf[i+3] = sample
	}

	//sdl.Delay(uint32(chip8.GetST() * 16))
	sdl.PauseAudio(true)
}
