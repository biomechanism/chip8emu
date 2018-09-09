package utils

type MachineState struct {
}

type MachineMonitor struct {
	active      bool
	cmdRunStep  bool
	breakPoints map[uint16]bool
}

func NewMachineMonitor() *MachineMonitor {
	return &MachineMonitor{breakPoints: make(map[uint16]bool)}
}

func (mm *MachineMonitor) Activate() {
	mm.active = true
}

func (mm *MachineMonitor) Deactivate() {
	mm.active = false
}

func (mm *MachineMonitor) IsActive() bool {
	return mm.active
}

func (mm *MachineMonitor) IsRunStep() bool {
	return mm.cmdRunStep
}

func (mm *MachineMonitor) SetRunStep() {
	mm.cmdRunStep = true
}

func (mm *MachineMonitor) ClearRunStep() {
	mm.cmdRunStep = false
}

func (mm *MachineMonitor) Reset() {
	mm.cmdRunStep = false
}

func (mm *MachineMonitor) SetBP(address uint16) {
	mm.breakPoints[address] = true
}

func (mm *MachineMonitor) ClrBP(address uint16) {
	mm.breakPoints[address] = false
}

func (mm *MachineMonitor) IsBP(address uint16) bool {
	//	println("-----IsBP----")
	//	println(address)
	//	println(mm.breakPoints[address])
	return mm.breakPoints[address]
}
