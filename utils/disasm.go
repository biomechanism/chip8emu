package utils

import (
	"encoding/base64"
	"fmt"
	"log"
)

type Inst struct {
	PrgAddress uint32
	Inst       uint16
	OpCode     uint16
	OpAddr     uint16
	Vx         uint16
	Vy         uint16
	Mnemonic   string
	EmitStr    string
}

func StrToByte(rom string) []byte {
	bytes, error := base64.StdEncoding.DecodeString(rom)

	if error != nil {
		log.Fatal(error.Error())
	}

	return bytes

}

func Disassemble(prog []byte) []*Inst {

	opStore := make([]*Inst, 0)

	//var inst uint16
	for i := 0; i < len(prog)-1; i += 2 {
		inst := uint16(prog[i]) << 8
		inst |= uint16(prog[i+1])
		Op := Decode(inst)
		opStore = append(opStore, Op)
	}

	return opStore
}

func Decode(inst uint16) *Inst {
	switch {
	case 0x00E0 == (inst & 0x0F0):
		return &Inst{
			Inst:     inst,
			Mnemonic: "CLS",
			OpCode:   inst & 0x00FF,
			EmitStr:  "\t" + "CLS" + "\t\t" + toHexStr(inst&0x00F0),
		}
	case 0x00EE == (inst & 0x00FF):
		return &Inst{
			Inst:     inst,
			Mnemonic: "RTS",
			OpCode:   inst & 0x00FF,
			EmitStr:  "\t" + "RTS" + "\t\t" + toHexStr(inst&0x00FF),
		}
	case 0x1000 == (inst & 0xF000):
		return &Inst{
			Inst:     inst,
			Mnemonic: "JMP",
			OpCode:   inst & 0xF000,
			OpAddr:   inst & 0x0FFF,
			EmitStr:  "\t" + "JMP" + "\t\t" + toHexStr(inst&0x0FFF),
		}
	case 0x2000 == (inst & 0xF000):
		return &Inst{
			Inst:     inst,
			Mnemonic: "JSP",
			OpCode:   inst & 0xF000,
			OpAddr:   inst & 0x0FFF,
			EmitStr:  "\t" + "JSR" + "\t\t" + toHexStr(inst&0x0FFF),
		}
	case 0x3000 == (inst & 0xF000):
		return &Inst{
			Inst:     inst,
			Mnemonic: "SKEQ",
			Vx:       inst & 0x0F00,
			OpCode:   inst & 0xF000,
			OpAddr:   inst & 0x00FF,
			EmitStr:  "\t" + "SKEQ" + "\t" + vxToStr(inst) + ", " + toHexStr(inst&0x00FF),
		}
	case 0x4000 == (inst & 0xF000):
		return &Inst{
			Inst:     inst,
			Mnemonic: "SKNE",
			Vx:       inst & 0x0F00,
			OpCode:   inst & 0xF000,
			OpAddr:   inst & 0x00FF,
			EmitStr:  "\t" + "SKNE" + "\t" + vxToStr(inst) + ", " + toHexStr(inst&0x00FF),
		}
	case 0x5000 == (inst & 0xF000):
		return &Inst{
			Inst:     inst,
			Mnemonic: "SKEQ",
			Vx:       inst & 0x0F00,
			Vy:       inst & 0x00F0,
			OpCode:   inst & 0xF000,
			EmitStr:  "\t" + "SKEQ" + "\t" + vxToStr(inst) + ", " + vyToStr(inst),
		}
	case 0x6000 == (inst & 0xF000):
		return &Inst{
			Inst:     inst,
			Mnemonic: "MOVE",
			Vx:       inst & 0x0F00,
			OpCode:   inst & 0xF000,
			OpAddr:   inst & 0x00FF,
			EmitStr:  "\t" + "MOVE" + "\t" + vxToStr(inst) + ", " + toHexStr(inst&0x00FF),
		}
	case 0x7000 == (inst & 0xF000):
		return &Inst{
			Inst:     inst,
			Mnemonic: "ADD",
			Vx:       inst & 0x0F00,
			OpCode:   inst & 0xF000,
			OpAddr:   inst & 0x00FF,
			EmitStr:  "\t" + "ADD" + "\t\t" + vxToStr(inst) + ", " + toHexStr(inst&0x00FF),
		}
	case 0x8000 == (inst & 0xF00F):
		return &Inst{
			Inst:     inst,
			Mnemonic: "MOVE",
			Vx:       inst & 0x0F00,
			Vy:       inst & 0x00F0,
			OpCode:   inst & 0xF000,
			EmitStr:  "\t" + "MOVE" + "\t" + vxToStr(inst) + ", " + vyToStr(inst),
		}
	case 0x8001 == (inst & 0xF00F):
		return &Inst{
			Inst:     inst,
			Mnemonic: "OR",
			Vx:       inst & 0x0F00,
			Vy:       inst & 0x00F0,
			OpCode:   inst & 0xF00F,
			EmitStr:  "\t" + "OR" + "\t\t" + vxToStr(inst) + ", " + vyToStr(inst),
		}
	case 0x8002 == (inst & 0xF00F):
		return &Inst{
			Inst:     inst,
			Mnemonic: "AND",
			Vx:       inst & 0x0F00,
			Vy:       inst & 0x00F0,
			OpCode:   inst & 0xF00F,
			EmitStr:  "\t" + "AND" + "\t\t" + vxToStr(inst) + ", " + vyToStr(inst),
		}
	case 0x8003 == (inst & 0xF00F):
		return &Inst{
			Inst:     inst,
			Mnemonic: "XOR",
			Vx:       inst & 0x0F00,
			Vy:       inst & 0x00F0,
			OpCode:   inst & 0xF00F,
			EmitStr:  "\t" + "XOR" + "\t\t" + vxToStr(inst) + ", " + vyToStr(inst),
		}
	case 0x8004 == (inst & 0xF00F):
		return &Inst{
			Inst:     inst,
			Mnemonic: "ADD",
			Vx:       inst & 0x0F00,
			Vy:       inst & 0x00F0,
			OpCode:   inst & 0xF00F,
			EmitStr:  "\t" + "ADD" + "\t\t" + vxToStr(inst) + ", " + vyToStr(inst),
		}
	case 0x8005 == (inst & 0xF00F):
		return &Inst{
			Inst:     inst,
			Mnemonic: "SUB",
			Vx:       inst & 0x0F00,
			Vy:       inst & 0x00F0,
			OpCode:   inst & 0xF00F,
			EmitStr:  "\t" + "SUB" + "\t\t" + vxToStr(inst) + ", " + vyToStr(inst),
		}
	case 0x8006 == (inst & 0xF00F):
		return &Inst{
			Inst:     inst,
			Mnemonic: "SHR",
			Vx:       inst & 0x0F00,
			OpCode:   inst & 0xF00F,
			EmitStr:  "\t" + "SHR" + "\t\t" + vxToStr(inst),
		}
	case 0x8007 == (inst & 0xF00F):
		return &Inst{
			Inst:     inst,
			Mnemonic: "RSUB",
			Vx:       inst & 0x0F00,
			Vy:       inst & 0x00F0,
			OpCode:   inst & 0xF00F,
			EmitStr:  "\t" + "RSUB" + "\t" + vxToStr(inst) + ", " + vyToStr(inst),
		}
	case 0x800E == (inst & 0xF00F):
		return &Inst{
			Inst:     inst,
			Mnemonic: "SHL",
			Vx:       inst & 0x0F00,
			OpCode:   inst & 0xF00F,
			EmitStr:  "\t" + "SHL" + "\t\t" + vxToStr(inst),
		}

	case 0x9000 == (inst & 0xF000):
		return &Inst{
			Inst:     inst,
			Mnemonic: "SKNE",
			Vx:       inst & 0x0F00,
			Vy:       inst & 0x00F0,
			OpCode:   inst & 0xF000,
			EmitStr:  "\t" + "SKNE" + "\t" + vxToStr(inst) + ", " + vyToStr(inst),
		}
	case 0xA000 == (inst & 0xF000):
		return &Inst{
			Inst:     inst,
			Mnemonic: "MOVE",
			OpCode:   inst & 0xF000,
			EmitStr:  "\t" + "MOVE" + "\t" + "I" + ", " + toHexStr(inst&0x0FFF),
		}
	case 0xB000 == (inst & 0xF000):
		return &Inst{
			Inst:     inst,
			Mnemonic: "JMP",
			OpCode:   inst & 0xF000,
			EmitStr:  "\t" + "JMP" + "\t\t" + "V0" + " + " + toHexStr(inst&0x0FFF),
		}
	case 0xC000 == (inst & 0xF000):
		return &Inst{
			Inst:     inst,
			Mnemonic: "RAND",
			Vx:       inst & 0x0F00,
			OpCode:   inst & 0xF000,
			EmitStr:  "\t" + "RAND" + "\t" + vxToStr(inst) + ", " + toHexStr(inst&0x00FF),
		}
	case 0xD000 == (inst & 0xF000):
		return &Inst{
			Inst:     inst,
			Mnemonic: "DRAW",
			Vx:       inst & 0x0F00,
			Vy:       inst & 0x00F0,
			OpCode:   inst & 0xF000,
			EmitStr:  "\t" + "DRAW" + "\t" + vxToStr(inst) + ", " + vyToStr(inst) + ", " + toHexStr(inst&0x000F),
		}
	case 0xE09E == (inst & 0xF0FF):
		return &Inst{
			Inst:     inst,
			Mnemonic: "SKPR",
			OpCode:   inst & 0xF0FF,
			EmitStr:  "\t" + "SKPR" + "\t" + "KEY " + toHexStr(inst&0x0F00),
		}
	case 0xE0A1 == (inst & 0xF0FF):
		return &Inst{
			Inst:     inst,
			Mnemonic: "SKUP",
			OpCode:   inst & 0xF0FF,
			EmitStr:  "\t" + "SKUP" + "\t" + "KEY " + toHexStr(inst&0x0F00),
		}
	case 0xF007 == (inst & 0xF00F):
		return &Inst{
			Inst:     inst,
			Vx:       inst & 0xF00F,
			Mnemonic: "MOVDT",
			OpCode:   inst & 0xF00F,
			EmitStr:  "\t" + "MOVDT" + "\t" + vxToStr(inst),
		}
	case 0xF00A == (inst & 0xF00F):
		return &Inst{
			Inst:     inst,
			Vx:       inst & 0xF0FF,
			Mnemonic: "KEYPR",
			OpCode:   inst & 0xF00F,
			EmitStr:  "\t" + "KEYPR" + "\t" + vxToStr(inst),
		}
	case 0xF015 == (inst & 0xF0FF):
		return &Inst{
			Inst:     inst,
			Vx:       inst & 0xF0FF,
			Mnemonic: "SETDT",
			OpCode:   inst & 0xF0FF,
			EmitStr:  "\t" + "SETDT" + "\t" + vxToStr(inst),
		}
	case 0xF018 == (inst & 0xF0FF):
		return &Inst{
			Inst:     inst,
			Vx:       inst & 0xF0FF,
			Mnemonic: "SETST",
			OpCode:   inst & 0xF0FF,
			EmitStr:  "\t" + "SETST" + "\t" + vxToStr(inst),
		}
	case 0xF01E == (inst & 0xF0FF):
		return &Inst{
			Inst:     inst,
			Vx:       inst & 0xF0FF,
			Mnemonic: "ADDI",
			OpCode:   inst & 0xF0FF,
			EmitStr:  "\t" + "ADDI" + "\t" + vxToStr(inst),
		}
	case 0xF029 == (inst & 0xF0FF):
		return &Inst{
			Inst:     inst,
			Vx:       inst & 0xF0FF,
			Mnemonic: "FONT",
			OpCode:   inst & 0xF0FF,
			EmitStr:  "\t" + "FONT" + "\t" + vxToStr(inst),
		}
	case 0xF033 == (inst & 0xF0FF):
		return &Inst{
			Inst:     inst,
			Vx:       inst & 0xF0FF,
			Mnemonic: "BCD",
			OpCode:   inst & 0xF0FF,
			EmitStr:  "\t" + "BCD" + "\t\t" + vxToStr(inst),
		}
	case 0xF055 == (inst & 0xF0FF):
		return &Inst{
			Inst:     inst,
			Vx:       inst & 0xF0FF,
			Mnemonic: "STOTE",
			OpCode:   inst & 0xF0FF,
			EmitStr:  "\t" + "STORE" + "\t" + vxToStr(inst),
		}
	case 0xF065 == (inst & 0xF0FF):
		return &Inst{
			Inst:     inst,
			Vx:       inst & 0xF0FF,
			Mnemonic: "LOAD",
			OpCode:   inst & 0xF0FF,
			EmitStr:  "\t" + "LOAD" + "\t" + vxToStr(inst),
		}
	default:
		return &Inst{
			Inst:    inst,
			EmitStr: "INVALID OPCODE",
		}
	}
}

func vxToStr(inst uint16) string {
	return regToStr(uint8((inst & 0x0F00) >> 8))
}

func vyToStr(inst uint16) string {
	return regToStr(uint8((inst & 0x00F0) >> 4))
}

func regToStr(reg uint8) string {
	return fmt.Sprintf("V%x", reg)
}

func toHexStr(inst uint16) string {
	return fmt.Sprintf("%x", inst)
}
