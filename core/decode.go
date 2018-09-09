package core

import (
	"fmt"
	"math/rand"
	"time"
)

//Instructions and other constants
const (
	opMask   = 0xF000
	opMask12 = 0xF0FF
	//	opMask8  = 0xF00F
	//	MASK8      = 0x00FF
	mask12     = 0x0FFF
	nibbleMask = 0x000F
	//Instructions
	Clear                = 0x00E0
	Return               = 0x00EE
	SysCall              = 0x0FFF //Not used in modern interpreters apparently
	Jump                 = 0x1000
	JumpSub              = 0x2000
	SkipVxEqKk           = 0x3000 //Skip next instruction if Vx == KK
	SkipVxNeqKk          = 0x4000
	SkipVxEqVy           = 0x5000
	LoadVxFromKk         = 0x6000
	LoadVxAddKk          = 0x7000
	LoadVxFromVy         = 0x8000
	LoadVxOrVy           = 0x8001
	LoadVxAndVy          = 0x8002
	LoadVxXorVy          = 0x8003
	LoadVxAddVy          = 0x8004
	LoadVxSubVy          = 0x8005
	LoadVxShiftR         = 0x8006
	LoadVxVySubVx        = 0x8007
	LoadVxShiftL         = 0x800E
	SkipVxNeqVy          = 0x9000
	LoadIFromNnn         = 0xA000
	JumpPlusA0           = 0xB000
	RndVxAndKk           = 0xC000
	DrawSprite           = 0xD000
	SkipVxEqKey          = 0xE09E
	SkipVxNeqKey         = 0xE0A1
	LoadVxFromDelayTimer = 0xF007
	LoadVxFromK          = 0xF00A
	LoadDelayTimerFromVx = 0xF015
	LoadSoundTimerFromVx = 0xF018
	LoadIAddVx           = 0xF01E
	LoadSpriteCharacter  = 0xF029
	LoadIWithBcdOfVx     = 0xF033
	StoreV0ToVx          = 0xF055
	LoadIToV0ToVx        = 0xF065
)

func shift12(val uint16) uint8 {
	return uint8(val >> 12)
}

func shift8(val uint16) uint8 {
	return uint8(val >> 8)
}

func decimalToBcd(val uint16) uint16 {
	var v uint16
	var h uint16
	var t uint16
	var o uint16

	for i := 0; i < 8; i++ {

		if h >= 5 {
			h += 3
		}

		if t >= 5 {
			t += 3
		}

		if o >= 5 {
			o += 3
		}

		h <<= 1
		h &= 0xF
		h |= ((t & 0x8) >> 3) & 0x000F

		t <<= 1
		t &= 0xF
		t |= ((o & 0x8) >> 3) & 0x000F

		o <<= 1
		o &= 0xF
		o |= ((val & 0x80) >> 7) & 0x000F

		val = val << 1

	}
	v = 0xFFF & ((h << 8) | (t << 4) | (o))
	return v
}

//Handle0x0 handler for Clear and Return instructions
func Handle0x0(chip *Chip8, opcode uint16) {
	switch opcode {
	case Clear:
		chip.ClearScreenMem()
	case Return:
		chip.SetPc(chip.Pop())
	default:
	}
}

//Handle0x1 intruction: Jump to address specified in NNN
//Opcode format 1NNN
func Handle0x1(chip *Chip8, opcode uint16) {
	switch opcode & opMask {
	case Jump:
		val := opcode & mask12
		chip.SetPc(val)
	}
}

//Handle0x2 instruction: Jump Subroutine, pushes the current PC to the stack and
//loads the PC with address specified with NNN
//Opcode format: 2NNN
func Handle0x2(chip *Chip8, opcode uint16) {
	switch (opcode & opMask) >> 12 {
	case JumpSub >> 12:
		chip.Push(chip.Pc)
		addr := opcode & mask12
		chip.Pc = addr
	default:
		fmt.Println("ARSE")
	}
}

//Handle0x3 instruction Skip next instruction if Vx equals NN
//Opcode format: 3XNN
func Handle0x3(chip *Chip8, opcode uint16) {
	switch opcode & opMask >> 12 {
	case SkipVxEqKk >> 12:
		vx := GetRegVx(opcode)
		kk := opcode & 0x00FF
		val := chip.GetV(vx)
		if val == uint8(kk) {
			chip.Pc += 2
		}
	}
}

//Handle0x4 instruction: Skip next instruction if Vx is not equal to value NN
//Opcode format 4XNN
func Handle0x4(chip *Chip8, opcode uint16) {
	switch opcode & opMask >> 12 {
	case SkipVxNeqKk >> 12:
		vx := GetRegVx(opcode)
		kk := opcode & 0x00FF
		val := chip.GetV(vx)
		if val != uint8(kk) {
			chip.Pc += 2
		}
	}
}

//Handle0x5 instruction: Skip the next instruction if Vx equals Vy
//Opcode format: 5XY0
func Handle0x5(chip *Chip8, opcode uint16) {
	switch opcode & opMask >> 12 {
	case SkipVxEqVy >> 12:
		vx := GetRegVx(opcode)
		vy := GetRegVy(opcode)
		r1 := chip.GetV(vx)
		r2 := chip.GetV(vy)
		if r1 == r2 {
			chip.Pc += 2
		}
	}
}

//Handle0x6 instruction: Load Vx from NN
//Opcode format: 6XNN
func Handle0x6(chip *Chip8, opcode uint16) {
	switch opcode & opMask >> 12 {
	case LoadVxFromKk >> 12:
		chip.SetV(GetRegVx(opcode), GetOpVal(opcode))
	}
}

//Handle0x7 instruction: Load Vx with the result of Vx add NN
//Opcode format: 7XNN
func Handle0x7(chip *Chip8, opcode uint16) {
	switch opcode & opMask >> 12 {
	case LoadVxAddKk >> 12:
		vx := GetRegVx(opcode)
		val := chip.GetV(vx)
		chip.SetV(vx, val+GetOpVal(opcode))
	}

}

//Handle0x8 handler for various instructions:
//Instruction: Load Vx from Vy
//Instruction: Load Vx with result of Vx bitwise Or Vy
//Instruction: Load Vx with result of Vx bitwise And Vy
//Instruction: Load Vx with result of Vx bitwise Xor Vy
//Instruction: Load Vx with result of Vx Add Vy
//Instruction: Load Vx with result of Vx Sub Vy
//Instruction: Load Vx with result of Right-Shift Vx
//Instruction: Load Vx with result of Vx Sub Vx
//Instruction: Load Vx with result of Left-Shift Vx
func Handle0x8(chip *Chip8, opcode uint16) {
	switch opcode & 0xF00F {
	case LoadVxFromVy & 0xF00F:
		chip.SetV(GetRegVx(opcode), chip.GetV(GetRegVy(opcode)))
	case LoadVxOrVy & 0xF00F:
		vx := chip.GetV(GetRegVx(opcode))
		vy := chip.GetV(GetRegVy(opcode))
		chip.SetV(GetRegVx(opcode), vx|vy)
	case LoadVxAndVy & 0xF00F:
		vx := chip.GetV(GetRegVx(opcode))
		vy := chip.GetV(GetRegVy(opcode))
		chip.SetV(GetRegVx(opcode), vx&vy)
	case LoadVxXorVy & 0xF00F:
		vx := chip.GetV(GetRegVx(opcode))
		vy := chip.GetV(GetRegVy(opcode))
		chip.SetV(GetRegVx(opcode), vx^vy)
	case LoadVxAddVy & 0xF00F:
		vx := chip.GetV(GetRegVx(opcode))
		vy := chip.GetV(GetRegVy(opcode))
		if uint16(vx)+uint16(vy) > 255 {
			chip.SetV(VF, 1)
		}
		chip.SetV(GetRegVx(opcode), vx+vy)
	case LoadVxSubVy & 0xF00F:
		vx := chip.GetV(GetRegVx(opcode))
		vy := chip.GetV(GetRegVy(opcode))
		if vx > vy {
			chip.SetV(VF, 1)
		} else {
			chip.SetV(VF, 0)
		}
		chip.SetV(GetRegVx(opcode), vx-vy)
	case LoadVxShiftR & 0xF00F:
		vx := chip.GetV(GetRegVx(opcode))
		chip.SetV(VF, vx&0x1)
		chip.SetV(GetRegVx(opcode), vx>>1)
	case LoadVxVySubVx & 0xF00F:
		vx := chip.GetV(GetRegVx(opcode))
		vy := chip.GetV(GetRegVy(opcode))
		if vy > vx {
			chip.SetV(VF, 1)
		} else {
			chip.SetV(VF, 0)
		}
		chip.SetV(GetRegVx(opcode), vy-vx)
	case LoadVxShiftL & 0xF00F:
		vx := chip.GetV(GetRegVx(opcode))
		chip.SetV(VF, (vx&0x80)>>7)
		chip.SetV(GetRegVx(opcode), vx<<1)
	default:
		fmt.Printf("Opcode: %v, Masked %v \n", opcode, LoadVxFromVy&0xF00F)
	}
}

//Handle0x9 Instruction: Skip if Vx != Vy
func Handle0x9(chip *Chip8, opcode uint16) {
	switch opcode & opMask >> 12 {
	case SkipVxNeqVy >> 12:
		vx := GetRegVx(opcode)
		vy := GetRegVy(opcode)
		r1 := chip.GetV(vx)
		r2 := chip.GetV(vy)
		if r1 != r2 {
			chip.Pc += 2
		}
	}
}

//Handle0xA instruction: Load Register I With Embedded Opcode Value NNN
//Opcode format: ANNN
func Handle0xA(chip *Chip8, opcode uint16) {
	switch opcode & opMask >> 12 {
	case LoadIFromNnn >> 12:
		val := GetOpVal12(opcode)
		chip.SetI(val)

	}
}

//Handle0xB instruction: Jump address V0 + NNN
//Opcode format: BNNN
func Handle0xB(chip *Chip8, opcode uint16) {
	switch opcode & opMask >> 12 {
	case JumpPlusA0 >> 12:
		val := GetOpVal12(opcode)
		v0 := chip.GetV(0)
		chip.SetPc(val + uint16(v0))
	}
}

//Handle0xC instruction: Load Register Vx with a bitwise AND operation with NN and a random value between 0 - 255
//Logic: Vx=rand() & NN
//Opcode format: CXNN
func Handle0xC(chip *Chip8, opcode uint16) {
	switch opcode & opMask >> 12 {
	case RndVxAndKk >> 12:
		val := GetOpVal(opcode)
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		rndVal := r.Intn(256)
		chip.SetV(GetRegVx(opcode), uint8(rndVal)&val)
	}
}

//Handle0xD instruction: Draw a Sprite to the screen at location Vx, Vy with height N
//Opcode format: DXYN
func Handle0xD(chip *Chip8, opcode uint16) {
	switch opcode & opMask >> 12 {
	case DrawSprite >> 12:
		vx := GetRegVx(opcode)
		vy := GetRegVy(opcode)
		n := opcode & nibbleMask
		chip.Draw(uint8(n), int32(chip.GetV(vx)), int32(chip.GetV(vy)))
	}
}

//Handle0xE Handler for instructions relate to key presses
func Handle0xE(chip *Chip8, opcode uint16) {
	switch opcode & opMask12 {
	case SkipVxEqKey:
		vx := GetRegVx(opcode)
		v := chip.GetV(vx)
		if chip.Keys[v] == 1 {
			chip.Pc += 2
		}
	case SkipVxNeqKey:
		vx := GetRegVx(opcode)
		v := chip.GetV(vx)
		if chip.Keys[v] == 0 {
			chip.Pc += 2
		}
	case LoadVxFromK:
		vx := GetRegVx(opcode)
		in := chip.DisplayHandler.WaitForInput()
		chip.SetV(vx, uint8(in))
	}
}

//Handle0xF handler for various instructions.
func Handle0xF(chip *Chip8, opcode uint16) {
	switch opcode & opMask12 {
	case LoadVxFromDelayTimer:
		vx := GetRegVx(opcode)
		chip.SetV(vx, chip.DelayTimer)
	case LoadDelayTimerFromVx:
		vx := GetRegVx(opcode)
		chip.SetDT(chip.GetV(vx))
	case LoadSoundTimerFromVx:
		vx := GetRegVx(opcode)
		chip.SetST(chip.GetV(vx))
	case LoadIAddVx:
		vx := GetRegVx(opcode)
		val := chip.GetI() + uint16(chip.GetV(vx))
		chip.SetI(val)
	case LoadSpriteCharacter:
		vx := GetRegVx(opcode)
		n := chip.GetV(vx)
		chip.SetI(GetCharBank(uint16(n)))
	case LoadIWithBcdOfVx:
		//Set I to BCD of Vx")
		vx := GetRegVx(opcode)
		n := chip.GetV(vx)
		bcd := decimalToBcd(uint16(n))
		idx := chip.GetI()
		chip.Memory[idx] = uint8((bcd & 0x0F00) >> 8)
		chip.Memory[idx+1] = uint8((bcd & 0x00F0) >> 4)
		chip.Memory[idx+2] = uint8((bcd & 0x000F))
	case StoreV0ToVx:
		vx := GetRegVx(opcode)
		idx := chip.GetI()
		var i uint8
		for i = 0; i <= vx; i++ {
			chip.Memory[idx] = chip.GetV(i)
			idx++
		}
	case LoadIToV0ToVx:
		vx := GetRegVx(opcode)
		idx := chip.GetI()
		var i uint8
		for i = 0; i <= vx; i++ {
			val := chip.Memory[idx]
			chip.SetV(i, val)
			idx++
		}
	}
}
