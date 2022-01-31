package emu

type Registers struct {
	A uint8
	B uint8
	C uint8
	D uint8
	E uint8
	H uint8
	L uint8
}

func (r *Registers) getAF(f uint8) uint16 {
	return (uint16(r.A) << 8) | uint16(f)
}

func (r *Registers) getBC() uint16 {
	return (uint16(r.B) << 8) | uint16(r.C)
}

func (r *Registers) getDE() uint16 {
	return (uint16(r.D) << 8) | uint16(r.E)
}

func (r *Registers) getHL() uint16 {
	return (uint16(r.H) << 8) | uint16(r.L)
}

func (r *Registers) setBC(val uint16) {
	r.B = uint8(val >> 8)
	r.C = uint8(val & 0xff)
}

func (r *Registers) setDE(val uint16) {
	r.D = uint8(val >> 8)
	r.E = uint8(val & 0xff)
}

func (r *Registers) setHL(val uint16) {
	r.H = uint8(val >> 8)
	r.L = uint8(val & 0xff)
}
