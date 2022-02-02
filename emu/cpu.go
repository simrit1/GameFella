package emu

import (
	"fmt"
	"io/ioutil"
	"os"
)

var (
	MEM_SIZE int32  = 65536
	PC       uint16 = 0x0100
	SP       uint16 = 0xFFFE
	C               = 0
)

type CPU struct {
	mem      *Memory
	reg      *Registers
	flags    *Flags
	sp       uint16
	pc       uint16
	cyc      int
	imeDelay int
	ime      bool
}

func NewCPU() *CPU {
	return &CPU{mem: NewMemory(MEM_SIZE), reg: NewRegisters(), flags: NewFlags(), pc: PC, sp: SP}
}

func (c *CPU) LoadRom(filename string) {
	rom, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	for i := 0; i < len(rom); i++ {
		c.mem.writeByte(uint16(i), rom[i])
	}
}

func (c *CPU) readByte(addr uint16) uint8 {
	return c.mem.readByte(addr)
}

func (c *CPU) readByteHL() uint8 {
	return c.readByte(c.reg.getHL())
}

func (c *CPU) writeByte(addr uint16, val uint8) {
	c.mem.writeByte(addr, val)
}

func (c *CPU) writeByteHL(val uint8) {
	c.writeByte(c.reg.getHL(), val)
}

func (c *CPU) nextByte() uint8 {
	val := c.readByte(c.pc)
	c.pc++
	return val
}

func (c *CPU) nextTwoBytes() uint16 {
	a := c.nextByte()
	b := c.nextByte()
	return (uint16(b) << 8) | uint16(a)
}

func (c *CPU) fetch() uint8 {
	return c.nextByte()
}

func (c *CPU) decode(opcode uint8) (func(*CPU), int) {
	if opcode == 0xCB {
		opcode = c.fetch()
		return CB_INSTRUCTIONS[opcode], CB_CYCLES[opcode]
	}
	return INSTRUCTIONS[opcode], CYCLES[opcode]
}

func (c *CPU) Execute(debug bool) {
	if debug {
		c.print()
	}
	c.checkIME()
	opcode := c.fetch()
	instr, cyc := c.decode(opcode)
	c.cyc += cyc * 4
	instr(c)
}

func (c *CPU) print() {
	fmt.Printf("A: %02X F: %02X B: %02X C: %02X D: %02X E: %02X H: %02X L: %02X SP: %04X PC: 00:%04X (%02X %02X %02X %02X)\n",
		c.reg.A, c.flags.getF(), c.reg.B, c.reg.C, c.reg.D, c.reg.E, c.reg.H, c.reg.L,
		c.sp, c.pc, c.readByte(c.pc), c.readByte(c.pc+1), c.readByte(c.pc+2), c.readByte(c.pc+3))
}

func (c *CPU) checkIME() {
	if c.imeDelay == 2 {
		c.imeDelay--
	} else if c.imeDelay == 1 {
		c.imeDelay--
		c.ime = true
	}
}

func unimplemented(c *CPU) {
	if c.mem.mem[c.pc-2] == 0xCB {
		fmt.Println("CB instruction")
	}
	fmt.Printf("unimplemented: 0x%02X\n", c.mem.mem[c.pc-1])
	os.Exit(0)
}

func flip(val uint8) uint8 {
	if val == 1 {
		return 0
	} else {
		return 1
	}
}

func (c *CPU) add8(a uint8, b uint8, cy uint8) uint8 {
	ans := uint16(a) + uint16(b) + uint16(cy)
	c.flags.setZero8(ans)
	c.flags.N = 0
	c.flags.setHalfCarryAdd8((a & 0xF), (b & 0xF))
	c.flags.setCarryAdd8(ans)
	return uint8(ans)
}

func (c *CPU) add16(a uint16, b uint16, cy uint8) uint16 {
	ans := uint32(a) + uint32(b) + uint32(cy)
	c.flags.setZero16(ans)
	c.flags.N = 0
	c.flags.setHalfCarryAdd16((a & 0xFF), (b & 0xFF))
	c.flags.setCarryAdd16(ans)
	return uint16(ans)
}

func (c *CPU) sub(a uint8, b uint8, cy uint8) uint8 {
	cy = flip(cy)
	ans := c.add8(a, ^b, cy)
	c.flags.N = 1
	c.flags.C = flip(c.flags.C)
	return uint8(ans)
}

func (c *CPU) and(a uint8, b uint8) uint8 {
	ans := uint16(a) & uint16(b)
	c.flags.setZero8(ans)
	c.flags.N = 0
	c.flags.H = 1
	c.flags.C = 0
	return uint8(ans)
}

func (c *CPU) or(a uint8, b uint8) uint8 {
	ans := uint16(a) | uint16(b)
	c.flags.setZero8(ans)
	c.flags.N = 0
	c.flags.H = 0
	c.flags.C = 0
	return uint8(ans)
}

func (c *CPU) xor(a uint8, b uint8) uint8 {
	ans := uint16(a) ^ uint16(b)
	c.flags.setZero8(ans)
	c.flags.N = 0
	c.flags.H = 0
	c.flags.C = 0
	return uint8(ans)
}

func (c *CPU) cp(a uint8, b uint8) {
	ans := uint16(a) - uint16(b)
	c.flags.setZero8(ans)
	c.flags.N = 1
	if (^(uint16(a) ^ ans ^ uint16(b)) & 0x10) > 0 {
		c.flags.H = 1
	} else {
		c.flags.H = 0
	}
	if b > a {
		c.flags.C = 1
	} else {
		c.flags.C = 0
	}
}

func (c *CPU) inc8(val uint8) uint8 {
	val++
	c.flags.setZero8(uint16(val))
	c.flags.N = 0
	if (val & 0xf) == 0 {
		c.flags.H = 1
	} else {
		c.flags.H = 0
	}
	return val
}

func (c *CPU) dec8(val uint8) uint8 {
	val--
	c.flags.setZero8(uint16(val))
	c.flags.N = 1
	if (val & 0xf) == 0xf {
		c.flags.H = 1
	} else {
		c.flags.H = 0
	}
	return val
}

func (c *CPU) bit(val uint8, u3 uint8) {
	c.flags.setZero8(uint16((val << u3) & 1))
	c.flags.N = 0
	c.flags.H = 1
}

func (c *CPU) jump(addr uint16, cond bool) {
	if cond {
		c.pc = addr
	}
}

func nop(c *CPU) {
}

func addAB(c *CPU) {
	c.reg.A = c.add8(c.reg.A, c.reg.B, 0)
}

func addAC(c *CPU) {
	c.reg.A = c.add8(c.reg.A, c.reg.C, 0)
}

func addAD(c *CPU) {
	c.reg.A = c.add8(c.reg.A, c.reg.D, 0)
}

func addAE(c *CPU) {
	c.reg.A = c.add8(c.reg.A, c.reg.E, 0)
}

func addAH(c *CPU) {
	c.reg.A = c.add8(c.reg.A, c.reg.H, 0)
}

func addAL(c *CPU) {
	c.reg.A = c.add8(c.reg.A, c.reg.L, 0)
}

func addAHL(c *CPU) {
	c.reg.A = c.add8(c.reg.A, c.readByteHL(), 0)
}

func addAA(c *CPU) {
	c.reg.A = c.add8(c.reg.A, c.reg.A, 0)
}

func adi(c *CPU) {
	c.reg.A = c.add8(c.reg.A, c.nextByte(), 0)
}

func adcAB(c *CPU) {
	c.reg.A = c.add8(c.reg.A, c.reg.B, c.flags.C)
}

func adcAC(c *CPU) {
	c.reg.A = c.add8(c.reg.A, c.reg.C, c.flags.C)
}

func adcAD(c *CPU) {
	c.reg.A = c.add8(c.reg.A, c.reg.D, c.flags.C)
}

func adcAE(c *CPU) {
	c.reg.A = c.add8(c.reg.A, c.reg.E, c.flags.C)
}

func adcAH(c *CPU) {
	c.reg.A = c.add8(c.reg.A, c.reg.H, c.flags.C)
}

func adcAL(c *CPU) {
	c.reg.A = c.add8(c.reg.A, c.reg.L, c.flags.C)
}

func adcAHL(c *CPU) {
	c.reg.A = c.add8(c.reg.A, c.readByteHL(), c.flags.C)
}

func adcAA(c *CPU) {
	c.reg.A = c.add8(c.reg.A, c.reg.A, c.flags.C)
}

func aci(c *CPU) {
	c.reg.A = c.add8(c.reg.A, c.nextByte(), c.flags.C)
}

func addHLBC(c *CPU) {
	c.reg.setHL(c.add16(c.reg.getHL(), c.reg.getBC(), 0))
}

func addHLDE(c *CPU) {
	c.reg.setHL(c.add16(c.reg.getHL(), c.reg.getDE(), 0))
}

func addHLHL(c *CPU) {
	c.reg.setHL(c.add16(c.reg.getHL(), c.reg.getHL(), 0))
}

func addHLSP(c *CPU) {
	c.reg.setHL(c.add16(c.reg.getHL(), c.sp, 0))
}

func addSP(c *CPU) {
	e8 := int8(c.nextByte())
	if e8 < 0 {
		e8 *= -1
		c.add16(c.sp, ^uint16(e8), 0)
	} else {
		c.add16(c.sp, uint16(e8), 0)
	}
	c.flags.Z = 0
	c.flags.N = 0
}

func subAB(c *CPU) {
	c.reg.A = c.sub(c.reg.A, c.reg.B, 0)
}

func subAC(c *CPU) {
	c.reg.A = c.sub(c.reg.A, c.reg.C, 0)
}

func subAD(c *CPU) {
	c.reg.A = c.sub(c.reg.A, c.reg.D, 0)
}

func subAE(c *CPU) {
	c.reg.A = c.sub(c.reg.A, c.reg.E, 0)
}

func subAH(c *CPU) {
	c.reg.A = c.sub(c.reg.A, c.reg.H, 0)
}

func subAL(c *CPU) {
	c.reg.A = c.sub(c.reg.A, c.reg.L, 0)
}

func subAHL(c *CPU) {
	c.reg.A = c.sub(c.reg.A, c.readByteHL(), 0)
}

func subAA(c *CPU) {
	c.reg.A = c.sub(c.reg.A, c.reg.A, 0)
}

func sui(c *CPU) {
	c.reg.A = c.sub(c.reg.A, c.nextByte(), 0)
}

func sbcAB(c *CPU) {
	c.reg.A = c.sub(c.reg.A, c.reg.B, c.flags.C)
}

func sbcAC(c *CPU) {
	c.reg.A = c.sub(c.reg.A, c.reg.C, c.flags.C)
}

func sbcAD(c *CPU) {
	c.reg.A = c.sub(c.reg.A, c.reg.D, c.flags.C)
}

func sbcAE(c *CPU) {
	c.reg.A = c.sub(c.reg.A, c.reg.E, c.flags.C)
}

func sbcAH(c *CPU) {
	c.reg.A = c.sub(c.reg.A, c.reg.H, c.flags.C)
}

func sbcAL(c *CPU) {
	c.reg.A = c.sub(c.reg.A, c.reg.L, c.flags.C)
}

func sbcAHL(c *CPU) {
	c.reg.A = c.sub(c.reg.A, c.readByteHL(), c.flags.C)
}

func sbcAA(c *CPU) {
	c.reg.A = c.sub(c.reg.A, c.reg.A, c.flags.C)
}

func sbi(c *CPU) {
	c.reg.A = c.sub(c.reg.A, c.nextByte(), c.flags.C)
}

func andAB(c *CPU) {
	c.reg.A = c.and(c.reg.A, c.reg.B)
}

func andAC(c *CPU) {
	c.reg.A = c.and(c.reg.A, c.reg.C)
}

func andAD(c *CPU) {
	c.reg.A = c.and(c.reg.A, c.reg.D)
}

func andAE(c *CPU) {
	c.reg.A = c.and(c.reg.A, c.reg.E)
}

func andAH(c *CPU) {
	c.reg.A = c.and(c.reg.A, c.reg.H)
}

func andAL(c *CPU) {
	c.reg.A = c.and(c.reg.A, c.reg.L)
}

func andAHL(c *CPU) {
	c.reg.A = c.and(c.reg.A, c.readByteHL())
}

func andAA(c *CPU) {
	c.reg.A = c.and(c.reg.A, c.reg.A)
}

func ani(c *CPU) {
	c.reg.A = c.and(c.reg.A, c.nextByte())
}

func orAB(c *CPU) {
	c.reg.A = c.or(c.reg.A, c.reg.B)
}

func orAC(c *CPU) {
	c.reg.A = c.or(c.reg.A, c.reg.C)
}

func orAD(c *CPU) {
	c.reg.A = c.or(c.reg.A, c.reg.D)
}

func orAE(c *CPU) {
	c.reg.A = c.or(c.reg.A, c.reg.E)
}

func orAH(c *CPU) {
	c.reg.A = c.or(c.reg.A, c.reg.H)
}

func orAL(c *CPU) {
	c.reg.A = c.or(c.reg.A, c.reg.L)
}

func orAHL(c *CPU) {
	c.reg.A = c.or(c.reg.A, c.readByteHL())
}

func orAA(c *CPU) {
	c.reg.A = c.or(c.reg.A, c.reg.A)
}

func ori(c *CPU) {
	c.reg.A = c.or(c.reg.A, c.nextByte())
}

func xorAB(c *CPU) {
	c.reg.A = c.xor(c.reg.A, c.reg.B)
}

func xorAC(c *CPU) {
	c.reg.A = c.xor(c.reg.A, c.reg.C)
}

func xorAD(c *CPU) {
	c.reg.A = c.xor(c.reg.A, c.reg.D)
}

func xorAE(c *CPU) {
	c.reg.A = c.xor(c.reg.A, c.reg.E)
}

func xorAH(c *CPU) {
	c.reg.A = c.xor(c.reg.A, c.reg.H)
}

func xorAL(c *CPU) {
	c.reg.A = c.xor(c.reg.A, c.reg.L)
}

func xorAHL(c *CPU) {
	c.reg.A = c.xor(c.reg.A, c.readByteHL())
}

func xorAA(c *CPU) {
	c.reg.A = c.xor(c.reg.A, c.reg.A)
}

func xri(c *CPU) {
	c.reg.A = c.xor(c.reg.A, c.nextByte())
}

func cpAB(c *CPU) {
	c.cp(c.reg.A, c.reg.B)
}

func cpAC(c *CPU) {
	c.cp(c.reg.A, c.reg.C)
}

func cpAD(c *CPU) {
	c.cp(c.reg.A, c.reg.D)
}

func cpAE(c *CPU) {
	c.cp(c.reg.A, c.reg.E)
}

func cpAH(c *CPU) {
	c.cp(c.reg.A, c.reg.H)
}

func cpAL(c *CPU) {
	c.cp(c.reg.A, c.reg.L)
}

func cpAHL(c *CPU) {
	c.cp(c.reg.A, c.readByteHL())
}

func cpAA(c *CPU) {
	c.cp(c.reg.A, c.reg.A)
}

func cpi(c *CPU) {
	c.cp(c.reg.A, c.nextByte())
}

func cpl(c *CPU) {
	c.reg.A ^= 255
	c.flags.N = 1
	c.flags.H = 1
}

func decB(c *CPU) {
	c.reg.B = c.dec8(c.reg.B)
}

func decC(c *CPU) {
	c.reg.C = c.dec8(c.reg.C)
}

func decD(c *CPU) {
	c.reg.D = c.dec8(c.reg.D)
}

func decE(c *CPU) {
	c.reg.E = c.dec8(c.reg.E)
}

func decH(c *CPU) {
	c.reg.H = c.dec8(c.reg.H)
}

func decL(c *CPU) {
	c.reg.L = c.dec8(c.reg.L)
}

func decHL(c *CPU) {
	c.writeByteHL(c.dec8(c.readByteHL()))
}

func decA(c *CPU) {
	c.reg.A = c.dec8(c.reg.A)
}

func decBC(c *CPU) {
	c.reg.setBC(c.reg.getBC() - 1)
}

func decDE(c *CPU) {
	c.reg.setDE(c.reg.getDE() - 1)
}

func decHL16(c *CPU) {
	c.reg.setHL(c.reg.getHL() - 1)
}

func decSP(c *CPU) {
	c.sp--
}

func incB(c *CPU) {
	c.reg.B = c.inc8(c.reg.B)
}

func incC(c *CPU) {
	c.reg.C = c.inc8(c.reg.C)
}

func incD(c *CPU) {
	c.reg.D = c.inc8(c.reg.D)
}

func incE(c *CPU) {
	c.reg.E = c.inc8(c.reg.E)
}

func incH(c *CPU) {
	c.reg.H = c.inc8(c.reg.H)
}

func incL(c *CPU) {
	c.reg.L = c.inc8(c.reg.L)
}

func incHL(c *CPU) {
	c.writeByteHL(c.inc8(c.readByteHL()))
}

func incA(c *CPU) {
	c.reg.A = c.inc8(c.reg.A)
}

func incBC(c *CPU) {
	c.reg.setBC(c.reg.getBC() + 1)
}

func incDE(c *CPU) {
	c.reg.setDE(c.reg.getDE() + 1)
}

func incHL16(c *CPU) {
	c.reg.setHL(c.reg.getHL() + 1)
}

func incSP(c *CPU) {
	c.sp++
}

func bit0B(c *CPU) {
	c.bit(c.reg.B, 0)
}

func bit0C(c *CPU) {
	c.bit(c.reg.C, 0)
}

func bit0D(c *CPU) {
	c.bit(c.reg.D, 0)
}

func bit0E(c *CPU) {
	c.bit(c.reg.E, 0)
}

func bit0H(c *CPU) {
	c.bit(c.reg.H, 0)
}

func bit0L(c *CPU) {
	c.bit(c.reg.L, 0)
}

func bit0A(c *CPU) {
	c.bit(c.reg.A, 0)
}

func bit1B(c *CPU) {
	c.bit(c.reg.B, 1)
}

func bit1C(c *CPU) {
	c.bit(c.reg.C, 1)
}

func bit1D(c *CPU) {
	c.bit(c.reg.D, 1)
}

func bit1E(c *CPU) {
	c.bit(c.reg.E, 1)
}

func bit1H(c *CPU) {
	c.bit(c.reg.H, 1)
}

func bit1L(c *CPU) {
	c.bit(c.reg.L, 1)
}

func bit1A(c *CPU) {
	c.bit(c.reg.A, 1)
}

func bit2B(c *CPU) {
	c.bit(c.reg.B, 2)
}

func bit2C(c *CPU) {
	c.bit(c.reg.C, 2)
}

func bit2D(c *CPU) {
	c.bit(c.reg.D, 2)
}

func bit2E(c *CPU) {
	c.bit(c.reg.E, 2)
}

func bit2H(c *CPU) {
	c.bit(c.reg.H, 2)
}

func bit2L(c *CPU) {
	c.bit(c.reg.L, 2)
}

func bit2A(c *CPU) {
	c.bit(c.reg.A, 2)
}

func bit3B(c *CPU) {
	c.bit(c.reg.B, 3)
}

func bit3C(c *CPU) {
	c.bit(c.reg.C, 3)
}

func bit3D(c *CPU) {
	c.bit(c.reg.D, 3)
}

func bit3E(c *CPU) {
	c.bit(c.reg.E, 3)
}

func bit3H(c *CPU) {
	c.bit(c.reg.H, 3)
}

func bit3L(c *CPU) {
	c.bit(c.reg.L, 3)
}

func bit3A(c *CPU) {
	c.bit(c.reg.A, 3)
}

func bit4B(c *CPU) {
	c.bit(c.reg.B, 4)
}

func bit4C(c *CPU) {
	c.bit(c.reg.C, 4)
}

func bit4D(c *CPU) {
	c.bit(c.reg.D, 4)
}

func bit4E(c *CPU) {
	c.bit(c.reg.E, 4)
}

func bit4H(c *CPU) {
	c.bit(c.reg.H, 4)
}

func bit4L(c *CPU) {
	c.bit(c.reg.L, 4)
}

func bit4A(c *CPU) {
	c.bit(c.reg.A, 4)
}

func bit5B(c *CPU) {
	c.bit(c.reg.B, 5)
}

func bit5C(c *CPU) {
	c.bit(c.reg.C, 5)
}

func bit5D(c *CPU) {
	c.bit(c.reg.D, 5)
}

func bit5E(c *CPU) {
	c.bit(c.reg.E, 5)
}

func bit5H(c *CPU) {
	c.bit(c.reg.H, 5)
}

func bit5L(c *CPU) {
	c.bit(c.reg.L, 5)
}

func bit5A(c *CPU) {
	c.bit(c.reg.A, 5)
}

func bit6B(c *CPU) {
	c.bit(c.reg.B, 6)
}

func bit6C(c *CPU) {
	c.bit(c.reg.C, 6)
}

func bit6D(c *CPU) {
	c.bit(c.reg.D, 6)
}

func bit6E(c *CPU) {
	c.bit(c.reg.E, 6)
}

func bit6H(c *CPU) {
	c.bit(c.reg.H, 6)
}

func bit6L(c *CPU) {
	c.bit(c.reg.L, 6)
}

func bit6A(c *CPU) {
	c.bit(c.reg.A, 6)
}

func bit7B(c *CPU) {
	c.bit(c.reg.B, 7)
}

func bit7C(c *CPU) {
	c.bit(c.reg.C, 7)
}

func bit7D(c *CPU) {
	c.bit(c.reg.D, 7)
}

func bit7E(c *CPU) {
	c.bit(c.reg.E, 7)
}

func bit7H(c *CPU) {
	c.bit(c.reg.H, 7)
}

func bit7L(c *CPU) {
	c.bit(c.reg.L, 7)
}

func bit7A(c *CPU) {
	c.bit(c.reg.A, 7)
}

func bit0HL(c *CPU) {
	c.bit(c.readByteHL(), 0)
}

func bit1HL(c *CPU) {
	c.bit(c.readByteHL(), 1)
}

func bit2HL(c *CPU) {
	c.bit(c.readByteHL(), 2)
}

func bit3HL(c *CPU) {
	c.bit(c.readByteHL(), 3)
}

func bit4HL(c *CPU) {
	c.bit(c.readByteHL(), 4)
}

func bit5HL(c *CPU) {
	c.bit(c.readByteHL(), 5)
}

func bit6HL(c *CPU) {
	c.bit(c.readByteHL(), 6)
}

func bit7HL(c *CPU) {
	c.bit(c.readByteHL(), 7)
}

func ldBB(c *CPU) {
}

func ldBC(c *CPU) {
	c.reg.B = c.reg.C
}

func ldBD(c *CPU) {
	c.reg.B = c.reg.D
}

func ldBE(c *CPU) {
	c.reg.B = c.reg.E
}

func ldBH(c *CPU) {
	c.reg.B = c.reg.H
}

func ldBL(c *CPU) {
	c.reg.B = c.reg.L
}

func ldBA(c *CPU) {
	c.reg.B = c.reg.A
}

func ldCB(c *CPU) {
	c.reg.C = c.reg.B
}

func ldCC(c *CPU) {}

func ldCD(c *CPU) {
	c.reg.C = c.reg.D
}

func ldCE(c *CPU) {
	c.reg.C = c.reg.E
}

func ldCH(c *CPU) {
	c.reg.C = c.reg.H
}

func ldCL(c *CPU) {
	c.reg.C = c.reg.L
}

func ldCA(c *CPU) {
	c.reg.C = c.reg.A
}

func ldDB(c *CPU) {
	c.reg.D = c.reg.B
}

func ldDC(c *CPU) {
	c.reg.D = c.reg.C
}

func ldDD(c *CPU) {}

func ldDE(c *CPU) {
	c.reg.D = c.reg.E
}

func ldDH(c *CPU) {
	c.reg.D = c.reg.H
}

func ldDL(c *CPU) {
	c.reg.D = c.reg.L
}

func ldDA(c *CPU) {
	c.reg.D = c.reg.A
}

func ldEB(c *CPU) {
	c.reg.E = c.reg.B
}

func ldEC(c *CPU) {
	c.reg.E = c.reg.C
}

func ldED(c *CPU) {
	c.reg.E = c.reg.D
}

func ldEE(c *CPU) {}

func ldEH(c *CPU) {
	c.reg.E = c.reg.H
}

func ldEL(c *CPU) {
	c.reg.E = c.reg.L
}

func ldEA(c *CPU) {
	c.reg.E = c.reg.A
}

func ldHB(c *CPU) {
	c.reg.H = c.reg.B
}

func ldHC(c *CPU) {
	c.reg.H = c.reg.C
}

func ldHD(c *CPU) {
	c.reg.H = c.reg.D
}

func ldHE(c *CPU) {
	c.reg.H = c.reg.E
}

func ldHH(c *CPU) {}

func ldHL(c *CPU) {
	c.reg.H = c.reg.L
}

func ldHA(c *CPU) {
	c.reg.H = c.reg.A
}

func ldLB(c *CPU) {
	c.reg.L = c.reg.B
}

func ldLC(c *CPU) {
	c.reg.L = c.reg.C
}

func ldLD(c *CPU) {
	c.reg.L = c.reg.D
}

func ldLE(c *CPU) {
	c.reg.L = c.reg.E
}

func ldLH(c *CPU) {
	c.reg.L = c.reg.H
}

func ldLL(c *CPU) {}

func ldLA(c *CPU) {
	c.reg.L = c.reg.A
}

func ldAB(c *CPU) {
	c.reg.A = c.reg.B
}

func ldAC(c *CPU) {
	c.reg.A = c.reg.C
}

func ldAD(c *CPU) {
	c.reg.A = c.reg.D
}

func ldAE(c *CPU) {
	c.reg.A = c.reg.E
}

func ldAH(c *CPU) {
	c.reg.A = c.reg.H
}

func ldAL(c *CPU) {
	c.reg.A = c.reg.L
}

func ldAA(c *CPU) {}

func ldiB(c *CPU) {
	c.reg.B = c.nextByte()
}

func ldiC(c *CPU) {
	c.reg.C = c.nextByte()
}

func ldiD(c *CPU) {
	c.reg.D = c.nextByte()
}

func ldiE(c *CPU) {
	c.reg.E = c.nextByte()
}

func ldiH(c *CPU) {
	c.reg.H = c.nextByte()
}

func ldiL(c *CPU) {
	c.reg.L = c.nextByte()
}

func ldiA(c *CPU) {
	c.reg.A = c.nextByte()
}

func ldBC16(c *CPU) {
	c.reg.setBC(c.nextTwoBytes())
}

func ldDE16(c *CPU) {
	c.reg.setDE(c.nextTwoBytes())
}

func ldHL16(c *CPU) {
	c.reg.setHL(c.nextTwoBytes())
}

func ldHLB(c *CPU) {
	c.writeByte(c.reg.getHL(), c.reg.B)
}

func ldHLC(c *CPU) {
	c.writeByte(c.reg.getHL(), c.reg.C)
}

func ldHLD(c *CPU) {
	c.writeByte(c.reg.getHL(), c.reg.D)
}

func ldHLE(c *CPU) {
	c.writeByte(c.reg.getHL(), c.reg.E)
}

func ldHLH(c *CPU) {
	c.writeByte(c.reg.getHL(), c.reg.H)
}

func ldHLL(c *CPU) {
	c.writeByte(c.reg.getHL(), c.reg.L)
}

func ldiHL(c *CPU) {
	c.writeByte(c.reg.getHL(), c.nextByte())
}

func ldBHL(c *CPU) {
	c.reg.B = c.readByte(c.reg.getHL())
}

func ldCHL(c *CPU) {
	c.reg.C = c.readByte(c.reg.getHL())
}

func ldDHL(c *CPU) {
	c.reg.D = c.readByte(c.reg.getHL())
}

func ldEHL(c *CPU) {
	c.reg.E = c.readByte(c.reg.getHL())
}

func ldHHL(c *CPU) {
	c.reg.H = c.readByte(c.reg.getHL())
}

func ldLHL(c *CPU) {
	c.reg.L = c.readByte(c.reg.getHL())
}

func ldBCA(c *CPU) {
	c.writeByte(c.reg.getBC(), c.reg.A)
}

func ldDEA(c *CPU) {
	c.writeByte(c.reg.getDE(), c.reg.A)
}

func ldHLA(c *CPU) {
	c.writeByte(c.reg.getHL(), c.reg.A)
}

func ld16A(c *CPU) {
	c.writeByte(c.nextTwoBytes(), c.reg.A)
}

func ldh16A(c *CPU) {
	addr := uint16(c.nextByte()) + 0xFF00
	if addr <= 0xFFFF {
		c.writeByte(addr, c.reg.A)
	}
}

func ldhCA(c *CPU) {
	addr := uint16(c.reg.C) + 0xFF00
	if addr <= 0xFFFF {
		c.writeByte(addr, c.reg.A)
	}
}

func ldABC(c *CPU) {
	c.reg.A = c.readByte(c.reg.getBC())
}

func ldADE(c *CPU) {
	c.reg.A = c.readByte(c.reg.getDE())
}

func ldAHL(c *CPU) {
	c.reg.A = c.readByte(c.reg.getHL())
}

func ldA16(c *CPU) {
	c.reg.A = c.readByte(c.nextTwoBytes())
}

func ldhA16(c *CPU) {
	addr := uint16(c.nextByte()) + 0xFF00
	if addr <= 0xFFFF {
		c.reg.A = c.readByte(addr)
	}
}

func ldhAC(c *CPU) {
	addr := uint16(c.reg.C) + 0xFF00
	if addr <= 0xFFFF {
		c.reg.A = c.readByte(addr)
	}
}

func ldHLIA(c *CPU) {
	ldHLA(c)
	c.reg.setHL(c.reg.getHL() + 1)
}

func ldHLDA(c *CPU) {
	ldHLA(c)
	c.reg.setHL(c.reg.getHL() - 1)
}

func ldAHLI(c *CPU) {
	ldAHL(c)
	c.reg.setHL(c.reg.getHL() + 1)
}

func ldAHLD(c *CPU) {
	ldAHL(c)
	c.reg.setHL(c.reg.getHL() - 1)
}

func ldSP16(c *CPU) {
	c.sp = c.nextTwoBytes()
}

func ld16SP(c *CPU) {
	addr := c.nextTwoBytes()
	c.writeByte(addr, uint8(c.sp&0xFF))
	c.writeByte(addr+1, uint8(c.sp>>8))
}

func ldHLSP(c *CPU) {
	og := c.sp
	e8 := int8(c.nextByte())
	if e8 < 0 {
		e8 *= -1
		c.add16(c.sp, ^uint16(e8), 0)
	} else {
		c.add16(c.sp, uint16(e8), 0)
	}
	c.reg.setHL(c.sp)
	c.sp = og
	c.flags.Z = 0
	c.flags.N = 0
}

func ldSPHL(c *CPU) {
	c.sp = c.reg.getHL()
}

func jp(c *CPU) {
	c.jump(c.nextTwoBytes(), true)
}

func jpz(c *CPU) {
	c.jump(c.nextTwoBytes(), c.flags.Z == 1)
}

func jpnz(c *CPU) {
	c.jump(c.nextTwoBytes(), c.flags.Z == 0)
}

func jpc(c *CPU) {
	c.jump(c.nextTwoBytes(), c.flags.C == 1)
}

func jpnc(c *CPU) {
	c.jump(c.nextTwoBytes(), c.flags.C == 0)
}

func jpHL(c *CPU) {
	c.jump(c.reg.getHL(), true)
}

func jr(c *CPU) {
	e8 := int8(c.nextByte())
	c.jump(uint16(int(c.pc)+int(e8)), true)
}

func jrz(c *CPU) {
	e8 := int8(c.nextByte())
	c.jump(uint16(int(c.pc)+int(e8)), c.flags.Z == 1)
}

func jrnz(c *CPU) {
	e8 := int8(c.nextByte())
	c.jump(uint16(int(c.pc)+int(e8)), c.flags.Z == 0)
}

func jrc(c *CPU) {
	e8 := int8(c.nextByte())
	c.jump(uint16(int(c.pc)+int(e8)), c.flags.C == 1)
}

func jrnc(c *CPU) {
	e8 := int8(c.nextByte())
	c.jump(uint16(int(c.pc)+int(e8)), c.flags.C == 0)
}

func di(c *CPU) {
	c.ime = false
	c.imeDelay = 0
}

func ei(c *CPU) {
	c.imeDelay = 2
}
