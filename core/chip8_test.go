package core

import (
	"fmt"
	"testing"
)

func NotImplemented() string {
	msg := fmt.Sprintf(" *** NOT IMPLEMENTED ***\n")
	return msg
}

func TestNewChip8(t *testing.T) {
	chip := NewChip8()
	chip.Memory[0] = 8

	if chip.Memory[0] == 8 {
		t.Log("Test Passed")
	} else {
		t.Error("Memory Write failed")
	}

}

func TestFetchOpCode(t *testing.T) {
	chip := NewChip8()
	//Currently getting set as 0xE000, looks like little endian handling, so
	//created new function SetMem to assign values to memory correctly (BigEndian).
	////chip.Memory[0x200] = 0x00E0
	chip.SetMem(0x200, 0x00E0)
	opcode := chip.fetch()

	if opcode == 0x00E0 {
		t.Log("Test Passed")
	} else {
		t.Error("Retrieved unexpected opcode: ", opcode)
	}
}

func TestDecodeOpCodeType0x1_JP(t *testing.T) {
	chip := NewChip8()
	var op uint16 = Jump | 0x123
	Handle0x1(chip, op)

	if chip.GetPc() == 0x123 {
		t.Log("[OPCODE] JP - Test Passed")
	} else {
		t.Error("[OPCODE] JP failed test")
	}

}

func TestHandlerCalled_JP(t *testing.T) {
	chip := NewChip8()
	var op uint16 = Jump | 0x123
	chip.InstHandlerTable.GetHandler(op)(chip, op)

	if chip.GetPc() == 0x123 {
		t.Log("[OPCODE] JP via InstructionTable - Test Passed")
	} else {
		t.Error("[OPCODE] JP via InstructionTable failed test")
	}

}

func TestHandlerCalled_LD_SET_VX_TO_VX_ADD_VY(t *testing.T) {
	chip := NewChip8()
	chip.SetV(2, 0xF)
	chip.SetV(4, 1)
	var op uint16 = LoadVxAddVy | (0x0024 << 4)
	chip.InstHandlerTable.GetHandler(op)(chip, op)

	if chip.GetV(2) == 16 {
		t.Log("LD, Vx AND Vy via InstructionTable - Test Passed")
	} else {
		msg := fmt.Sprintf("Expected %v, received %v", LoadVxAddVy, chip.GetV(2))
		t.Error(msg)
	}

}

func TestHandlerCalled_LD_SET_VX_TO_VX_SUB_VY(t *testing.T) {
	chip := NewChip8()
	chip.SetV(2, 6)
	chip.SetV(4, 2)
	var op uint16 = LoadVxSubVy | (0x0024 << 4)
	chip.InstHandlerTable.GetHandler(op)(chip, op)

	if chip.GetV(2) == 4 && chip.GetV(0xF) == 1 {
		t.Log("LD, Vx SUB Vy Via InstructionTabe - Test Passed")
	} else {
		msg := fmt.Sprintf("V2 = %v", chip.GetV(2))
		t.Error(msg)
	}

}

func TestDecodeOpCodeRET(t *testing.T) {
	chip := NewChip8()
	chip.SetPc(0x40)
	chip.Push(0x100)
	chip.Sp = 14

	var op uint16 = Return
	Handle0x0(chip, op)

	if chip.GetPc() == 0x100 {
		t.Log("[OPCODE] RET - Test Passed")
	} else {
		t.Error("[OPCODE] RET failed test")
	}

}

func TestDecodeOpCodeCALL(t *testing.T) {
	chip := NewChip8()
	//chip.SetPC(0x40)
	var op uint16 = JumpSub | (0x0040)
	Handle0x2(chip, op)
	if chip.S[chip.Sp+1] == 0x200 && chip.GetPc() == 0x40 {
		t.Log("CALL - Test Passed")
	} else {
		msg := fmt.Sprintf("Expected %v, received %v", 0x40, chip.S[chip.Sp+1])
		t.Error(msg)
	}
}

func TestDecodeOpCode_SNE(t *testing.T) {
	chip := NewChip8()
	chip.SetV(2, 6)

	var op uint16 = SkipVxNeqKk | (0x0207)
	Handle0x4(chip, op)

	if chip.Pc == 514 {
		t.Log("SNE, Skip Next Instruction if Vx != KK - Test Passed")
	} else {
		msg := fmt.Sprintf("Expected PC to equal 514, was actually %v", chip.Pc)
		t.Error(msg)
	}

}

func TestDecodeOpCode_SE(t *testing.T) {
	chip := NewChip8()
	chip.SetV(2, 6)

	var op uint16 = SkipVxEqKk | (0x0206)
	Handle0x3(chip, op)

	if chip.Pc == 0x202 {
		t.Log("SE, Skip Next Instruction if Vx == KK - Test Passed")
	} else {
		msg := fmt.Sprintf("Expected PC to equal 514, was actually %v", chip.Pc)
		t.Error(msg)
	}

}

func TestDecodeOpCode_SE_VX_VY(t *testing.T) {
	chip := NewChip8()
	chip.SetV(2, 6)
	chip.SetV(3, 6)

	var op uint16 = SkipVxEqVy | (0x0230)
	Handle0x5(chip, op)

	if chip.Pc == 0x202 {
		t.Log("SE, Skip Next Instruction if Vx == VY - Test Passed")
	} else {
		msg := fmt.Sprintf("Expected PC to equal 514, was actually %v", chip.Pc)
		t.Error(msg)
	}

}

func TestDecodeOpCode_SE_VX_NEQ_VY(t *testing.T) {
	chip := NewChip8()
	chip.SetV(2, 6)
	chip.SetV(3, 4)

	var op uint16 = SkipVxNeqVy | (0x0230)
	Handle0x9(chip, op)

	if chip.Pc == 0x202 {
		t.Log("SE, Skip Next Instruction if Vx != VY - Test Passed")
	} else {
		msg := fmt.Sprintf("Expected PC to equal 514, was actually %v", chip.Pc)
		t.Error(msg)
	}

}

func TestDecodeOpCode_LD_VX_KK(t *testing.T) {
	chip := NewChip8()
	var op uint16 = LoadVxFromKk | (0x0207)
	Handle0x6(chip, op)

	vx := chip.GetV(2)

	if vx == 7 {
		t.Log("LD, SET V2 = KK - Test Passed")
	} else {
		msg := fmt.Sprintf("Failed to set register correctly, expected val %v, actual %v", 7, chip.GetV(2))
		t.Error(msg)
	}
}

func TestDecodeOpCode_LD_VX_ADD_KK(t *testing.T) {
	chip := NewChip8()
	var op uint16 = LoadVxAddKk | (0x0207)

	chip.SetV(GetRegVx(op), 2)
	Handle0x7(chip, op)
	vx := chip.GetV(GetRegVx(op))

	if vx == 9 {
		t.Log("LD, SET VX = VX ADD KK - Test Passed")
	} else {
		msg := fmt.Sprintf("Failed to set register correctly, expected val %v, actual %v", 9, chip.GetV(GetRegVx(op)))
		t.Error(msg)
	}
}

func TestDecodeOpCodeCALL0x8(t *testing.T) {
	chip := NewChip8()
	var op uint16 = LoadVxFromVy | (0x0024 << 4)
	Handle0x8(chip, op)

	vx := GetRegVx(op)
	vy := GetRegVy(op)

	fmt.Printf("V[%v], v[%v], OP = %v", vx, vy, op)

	chip.SetV(vx, 3)
	chip.SetV(vy, 4)

	chip.SetV(vx, chip.GetV(vy))
	//	chip.V[vx] = chip.V[vy]

	if chip.V[2] == 4 && chip.V[4] == 4 {
		t.Log("LD, V2 set to V4 - Test Passed")
	} else {
		msg := fmt.Sprintf("Expected V2 == 4, V4 == 4")
		t.Error(msg)
	}

}

func TestDecodeOpCodeType0x8_LD_OR(t *testing.T) {
	chip := NewChip8()
	chip.SetV(2, 5)
	chip.SetV(4, 2)
	var op uint16 = LoadVxOrVy | (0x0024 << 4)
	Handle0x8(chip, op)

	r2 := chip.GetV(2)

	if r2 == 7 {
		t.Log("LD, VX = VX OR VY - Test Passed")
	} else {
		t.Logf("Expected value 7 in V2, received %v\n", r2)
	}

}

func TestDecodeOpCodeType0x8_LD_AND(t *testing.T) {
	chip := NewChip8()
	chip.SetV(2, 7)
	chip.SetV(4, 2)
	var op uint16 = LoadVxAndVy | (0x0024 << 4)
	Handle0x8(chip, op)

	r2 := chip.GetV(2)

	if r2 == 2 {
		t.Log("LD, VX = VX AND VY - Test Passed")
	} else {
		t.Logf("Expected value 2 in V2, received %v\n", r2)
	}

}

func TestDecodeOpCodeType0x8_LD_ADD_NoOverflow(t *testing.T) {
	chip := NewChip8()
	chip.SetV(2, 0xF)
	chip.SetV(4, 1)
	var op uint16 = LoadVxAddVy | (0x0024 << 4)
	Handle0x8(chip, op)

	if chip.GetV(2) == 16 {
		t.Log("LD, Vx AND Vy - Test Passed")
	} else {
		msg := fmt.Sprintf("Expected %v, received %v", LoadVxAddVy, chip.GetV(2))
		t.Error(msg)
	}

}

func TestDecodeOpCodeType0x8_LD_ADD_Overflow(t *testing.T) {
	chip := NewChip8()
	chip.SetV(2, 0xFD)
	chip.SetV(4, 5)
	var op uint16 = LoadVxAddVy | (0x0024 << 4)
	Handle0x8(chip, op)

	//FIXME: Also check Vx
	if chip.GetV(0xF) == 1 {
		t.Log("LD, Vx AND Vy With Overflow- Test Passed")
	} else {
		msg := fmt.Sprintf("V2 = %v", chip.GetV(2))
		t.Error(msg)
	}

	//if ret == (LD_SET_VX_TO_VX_ADD_VY & 0xF00F) {
	//	t.Log("LD, Vx AND Vy - Test Passed")
	//} else {
	//	msg := fmt.Sprintf("Expected %v, received %v", LD_SET_VX_TO_VX_ADD_VY, ret)
	//	t.Error(msg)
	//}

}

func TestDecodeOpCodeType0x8_LD_SUB(t *testing.T) {
	chip := NewChip8()
	chip.SetV(2, 6)
	chip.SetV(4, 2)
	var op uint16 = LoadVxSubVy | (0x0024 << 4)
	Handle0x8(chip, op)

	if chip.GetV(2) == 4 && chip.GetV(0xF) == 1 {
		t.Log("LD, Vx SUB Vy - Test Passed")
	} else {
		msg := fmt.Sprintf("V2 = %v", chip.GetV(2))
		t.Error(msg)
	}

}

func TestDecodeOpCodeType0x8_LD_SHR(t *testing.T) {
	chip := NewChip8()
	chip.SetV(2, 9)
	var op uint16 = LoadVxShiftR | (0x0200)
	Handle0x8(chip, op)

	if chip.GetV(2) == 4 && chip.GetV(VF) == 1 {
		t.Log("LD, Vx = Vx SHR (LSB == 1, VF = 1) - Test Passed")
	} else {
		msg := fmt.Sprintf("V2 = %v", chip.GetV(2))
		t.Error(msg)
	}

	chip.SetV(3, 8)
	op = LoadVxShiftR | (0x0300)
	Handle0x8(chip, op)

	if chip.GetV(3) == 4 && chip.GetV(VF) == 0 {
		t.Log("LD, Vx = Vx SHR (LSB == 0, VF = 0) - Test Passed")
	} else {
		msg := fmt.Sprintf("V3 = %v", chip.GetV(3))
		t.Error(msg)
	}

}

func TestDecodeOpCodeType0x8_LD_SUBN(t *testing.T) {
	chip := NewChip8()
	chip.SetV(V2, 2)
	chip.SetV(V4, 8)
	var op uint16 = LoadVxVySubVx | (0x0024 << 4)
	Handle0x8(chip, op)

	if chip.GetV(2) == 6 && chip.GetV(0xF) == 1 {
		t.Log("LD, Vy SUB Vx (8-2) - Test Passed")
	} else {
		msg := fmt.Sprintf("V2 = %v\n", chip.GetV(2))
		t.Error(msg)
	}

	chip.SetV(V2, 8)
	chip.SetV(V4, 2)
	op = LoadVxVySubVx | (0x0024 << 4)
	Handle0x8(chip, op)
	fmt.Printf("V2 = %v\n", chip.GetV(2))
	if chip.GetV(2) == 250 && chip.GetV(0xF) == 0 {
		t.Log("LD, Vy SUB Vx (2-8) - Test Passed")
	} else {
		msg := fmt.Sprintf("V2 = %v", chip.GetV(2))
		t.Error(msg)
	}

}

func TestDecodeOpCodeType0x8_LD_SHL(t *testing.T) {
	chip := NewChip8()
	chip.SetV(V2, 129)
	var op uint16 = LoadVxShiftL | (0x0200)
	Handle0x8(chip, op)

	if chip.GetV(V2) == 2 && chip.GetV(VF) == 1 {
		t.Log("LD, Vx SHL (MSB == 1, VF = 1) - Test Passed")
	} else {
		msg := fmt.Sprintf("V2 = %v, VF = %v", chip.GetV(2), chip.GetV(0xF))
		t.Error(msg)
	}

	chip.SetV(V3, 8)
	op = LoadVxShiftL | (0x0300)
	Handle0x8(chip, op)

	if chip.GetV(3) == 16 && chip.GetV(0xf) == 0 {
		t.Log("LD, Vx SHL (MSB == 0, VF = 0) - Test Passed")
	} else {
		msg := fmt.Sprintf("V3 = %v", chip.GetV(3))
		t.Error(msg)
	}

}

func TestDecodeOpCodeType0xA_LD_I(t *testing.T) {
	chip := NewChip8()
	var op uint16 = LoadIFromNnn | (0x0222)
	Handle0xA(chip, op)

	if chip.GetI() == 0x0222 {
		t.Log("LD, I = KKK - Test Passed")
	} else {
		msg := fmt.Sprintf("I = %v", chip.GetI())
		t.Error(msg)
	}

}

//
//func TestDecodeOpCodeType0x1F_JP(t *testing.T) {
//	chip := NewChip8()
//	var op uint16 = JP | (512)
//
//}

func TestDecodeOpCodeType0xB_JP_NNN_PLUS_A0(t *testing.T) {
	chip := NewChip8()
	var op uint16 = JumpPlusA0 | (0x0222)
	chip.SetV(V0, 1)
	Handle0xB(chip, op)

	if chip.GetPc() == 0x0223 {
		t.Log("JP, NNN PLUS A0 - Test Passed")
	} else {
		msg := fmt.Sprintf("Expected PC = %v, actual PC = %v\n", 0x0223, chip.GetPc())
		t.Error(msg)
	}

}

//func TestDecodeOpCodeType0xC_RND_VX_AND_KK(t *testing.T) {
//	chip := NewChip8()
//	var op uint16 = RND_VX_AND_KK | (0x0222)
//
//}

// func TestDecodeOpCodeType0xD_DRW(t *testing.T) {
// 	chip := NewChip8()
// 	var op uint16 = DrawSprite | (0x0233)
// 	chip.Memory[3] = 0xFF
// 	chip.Memory[4] = 0xAA
// 	chip.Memory[5] = 0xC
// 	chip.SetI(3)
// 	chip.SetV(V2, 1)
// 	chip.SetV(V3, 1)

// 	chip.Memory[GetScreenMemIdx(1, 3)] = 8

// 	fmt.Printf("SCR MEM (1,1): %v\n", chip.Memory[GetScreenMemIdx(1, 1)])
// 	fmt.Printf("SCR MEM (1,2): %v\n", chip.Memory[GetScreenMemIdx(1, 2)])
// 	fmt.Printf("SCR MEM (1,3): %v\n", chip.Memory[GetScreenMemIdx(1, 3)])
// 	fmt.Printf("MEM 5: %v\n", chip.Memory[0xf00+5])

// 	Handle0xD(chip, op)
// 	fmt.Println("--AFTER--")
// 	fmt.Printf("SCR MEM (1,1): %v\n", chip.Memory[GetScreenMemIdx(1, 1)])
// 	fmt.Printf("SCR MEM (1,2): %v\n", chip.Memory[GetScreenMemIdx(1, 2)])
// 	fmt.Printf("SCR MEM (1,3): %v\n", chip.Memory[GetScreenMemIdx(1, 3)])
// 	fmt.Printf("MEM 5: %v\n", chip.Memory[0xf00+5])
// }

func TestDecodeOpCodeType0xE_SKIP_ON_KEYPRESS(t *testing.T) {
	chip := NewChip8()
	var op uint16 = SkipVxEqKey | (0x0200)
	chip.SetV(V2, 2)
	chip.SetKey(V2)
	Handle0xE(chip, op)

	if chip.Pc == 0x202 {
		t.Log("SKP, Skip Next Instruction on KeyPress - Test Passed")
	} else {
		msg := fmt.Sprintf("Expected PC to equal 516, was actually %v", chip.Pc)
		t.Error(msg)
	}

}

func TestDecodeOpCodeType0xE_SKIP_ON_NO_KEYPRESS(t *testing.T) {
	chip := NewChip8()
	var op uint16 = SkipVxNeqKey | (0x0200)
	Handle0xE(chip, op)

	if chip.Pc == 0x202 {
		t.Log("SKP, Skip Next Instruction if no KeyPress  - Test Passed")
	} else {
		msg := fmt.Sprintf("Expected PC to equal 516, was actually %v", chip.Pc)
		t.Error(msg)
	}

}

//func TestDecodeOpCodeType0xE_WAIT_FOR_KEYINPUT(t *testing.T) {
//
//	chip := NewChip8WithHandler(new(Handler))
//	var op uint16 = LD_SET_VX_TO_K | (0x0200)
//	v2 := chip.GetV(2)
//	fmt.Printf("V2 (BEFORE) = %v\n", v2)
//	HandleOpCodeType0xE_WAIT_FOR_KEYPRESS(chip, op)
//	v2 = chip.GetV(2)
//	fmt.Printf("V2 (AFTER) = %v\n", v2)
//
//}

func TestDecodeOpCodeType0xF_LD_VX_WITH_DT(t *testing.T) {
	chip := NewChip8()
	var op uint16 = LoadVxFromDelayTimer | (0x0200)
	chip.DisableDelayTimer()
	chip.SetDT(105)
	Handle0xF(chip, op)

	vx := GetRegVx(op)

	if chip.GetV(vx) == 105 {
		t.Log("LD, SET V2 to DT - Test Passed")
	} else {
		msg := fmt.Sprintf("Expected V2 to be 105, was actually %v", chip.GetV(2))
		t.Error(msg)
	}
}

func TestDecodeOpCodeType0xF_LD_DT_WITH_VX(t *testing.T) {
	chip := NewChip8()
	var op uint16 = LoadDelayTimerFromVx | (0x0200)
	chip.DisableDelayTimer()
	chip.SetV(V2, 105)
	Handle0xF(chip, op)

	if chip.GetDT() == 105 {
		t.Log("LD, SET DT to V2 - Test Passed")
	} else {
		msg := fmt.Sprintf("Expected DT to be 105, was actually %v", chip.GetDT())
		t.Error(msg)
	}

}

func TestDecodeOpCodeType0xF_LD_ST_WITH_VX(t *testing.T) {
	chip := NewChip8()
	var op uint16 = LoadSoundTimerFromVx | (0x0200)
	chip.DisableSoundTimer()
	chip.SetV(V2, 105)
	Handle0xF(chip, op)

	if chip.GetST() == 105 {
		t.Log("LD, SET ST to V2 - Test Passed")
	} else {
		msg := fmt.Sprintf("Expected DT to be 105, was actually %v", chip.GetDT())
		t.Error(msg)
	}

}

func TestDecodeOpCodeType0xF_LD_SET_I_TO_I_PLUS_VX(t *testing.T) {
	chip := NewChip8()
	chip.SetI(200)
	chip.SetV(V2, 5)
	var op uint16 = LoadIAddVx | (0x0200)
	Handle0xF(chip, op)

	if chip.GetI() == 205 {
		t.Log("LD, SET I to I + V2 - Test Passed")
	} else {
		msg := fmt.Sprintf("Expected I to be 205, was actually %v", chip.GetI())
		t.Error(msg)
	}

}

func TestDecodeOpCodeType0xF_LD_SET_SPRITE_CHAR_FROM_VX(t *testing.T) {
	chip := NewChip8()
	chip.SetV(2, 3)
	var op uint16 = LoadSpriteCharacter | (0x0200)
	Handle0xF(chip, op)

	bank := chip.GetI()

	if bank == CharBank3 {
		t.Log("LD, SET SPRITE CHAR FROM VX - Test Passed")
	} else {
		msg := fmt.Sprintf("Expected I to be %v, was actually %v", CharBank3, chip.GetI())
		t.Error(msg)
	}

}

func TestDecodeOpCodeType0xF_LD_I_WITH_BCD(t *testing.T) {
	chip := NewChip8()
	chip.SetI(100)
	chip.SetV(2, 162)
	var op uint16 = LoadIWithBcdOfVx | (0x0200)
	Handle0xF(chip, op)

	if chip.Memory[100] == 1 && chip.Memory[101] == 6 && chip.Memory[102] == 2 {
		t.Log("LD, LOAD BCD OF VX INTO MEM I - Test Passed")
	} else {
		msg := fmt.Sprintf("Expected I to be loaded with BCD")
		t.Error(msg)
	}
}

func TestDecodeOpCodeType0xF_DECIMAL_TO_BCD(t *testing.T) {
	//chip := NewChip8()
	bcd := decimalToBcd(162)

	if bcd == 354 {
		t.Log("CONVERT DECIMAL TO BCD - Test Passed")
	} else {
		msg := fmt.Sprintf("Expected result %v, was actually %v", 354, bcd)
		t.Error(msg)
	}
}

//FIXME: Run through full register list
func TestDecodeOpCodeType0xF_STORE_V0_TO_VX_AT_I(t *testing.T) {
	chip := NewChip8()
	chip.SetI(200)
	chip.SetV(0, 1)
	chip.SetV(1, 2)
	chip.SetV(2, 3)
	chip.SetV(3, 4)
	chip.SetV(4, 5)
	var op uint16 = StoreV0ToVx | (0x0400)

	Handle0xF(chip, op)

	if Mem(200, chip) == 1 && Mem(201, chip) == 2 && Mem(202, chip) == 3 && Mem(203, chip) == 4 {
		t.Log("STORED V0 TO VX AT I - Test Passed")
	} else {
		msg := fmt.Sprintf("V0 to VX not stored at I correctly")
		t.Error(msg)
	}

}

//FIXME: Run through full register list
func TestDecodeOpCodeType0xF_READ_I_INTO_V0_TO_VX(t *testing.T) {
	chip := NewChip8()
	chip.Memory[200] = 1
	chip.Memory[201] = 2
	chip.Memory[202] = 3
	chip.Memory[203] = 4
	chip.SetI(200)
	var op uint16 = LoadIToV0ToVx | (0x0400)
	Handle0xF(chip, op)

	if chip.GetV(0) == 1 && chip.GetV(1) == 2 && chip.GetV(2) == 3 && chip.GetV(3) == 4 {
		t.Log("READ FROM I INTO V0 TO VX - Test Passed")
	} else {
		msg := fmt.Sprintf("I not read correctly into V0 to VX")
		t.Error(msg)
	}
}

func Mem(idx uint16, c *Chip8) uint8 {
	return c.Memory[idx]
}

//Very rough, needs proper work.
func TestExecuteInstruction(t *testing.T) {
	chip := NewChip8()
	chip.SetV(2, 3)
	chip.SetV(4, 4)
	var op uint16 = LoadVxAddVy | (0x0024 << 4)
	chip.Memory[512] = uint8(op >> 8)
	chip.Memory[513] = uint8(op & 0x00FF)

	var op2 uint16 = Jump | (512)
	chip.Memory[514] = uint8(op2 >> 8)
	chip.Memory[515] = uint8(op2 & 0x00FF)

	inst := chip.fetch()
	chip.execute(inst)

	//Reset PC to 0x200
	inst = chip.fetch()
	chip.execute(inst)

	inst = chip.fetch()
	chip.execute(inst)

	//Reset PC to 0x200
	inst = chip.fetch()
	chip.execute(inst)

	inst = chip.fetch()
	chip.execute(inst)

	if chip.GetV(2) == 15 {
		t.Log("Multiple executions (JP, ADD), ADD VX and VX - Test Passed")
	} else {
		msg := fmt.Sprintf("Incorrect operation, result should have been 15, was %v", chip.GetV(2))
		t.Error(msg)
	}

}

//
//func TestDrawInstruction(t *testing.T) {
//	chip := NewChip8()
//	chip.SetI(0)
//	chip.Memory[0] = 24
//	chip.Memory[1] = 126
//	chip.SetV(2, 10)
//	chip.SetV(3, 10)
//	var inst uint16 = DRW_VX_VY_NIBBLE | (0x0232)
//
//	chip.Execute(inst)
//
//	runtime.LockOSThread()
//	view.NewSDLDisplayRenderer(chip)
//	chip.DoTestStart()
//
//}
