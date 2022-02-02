package emu

var (
	CYCLES = []int{
		1, 3, 2, 2, 1, 1, 2, 1, 5, 2, 2, 2, 1, 1, 2, 1,
		0, 3, 2, 2, 1, 1, 2, 1, 3, 2, 2, 2, 1, 1, 2, 1,
		2, 3, 2, 2, 1, 1, 2, 1, 2, 2, 2, 2, 1, 1, 2, 1,
		2, 3, 2, 2, 3, 3, 3, 1, 2, 2, 2, 2, 1, 1, 2, 1,
		1, 1, 1, 1, 1, 1, 2, 1, 1, 1, 1, 1, 1, 1, 2, 1,
		1, 1, 1, 1, 1, 1, 2, 1, 1, 1, 1, 1, 1, 1, 2, 1,
		1, 1, 1, 1, 1, 1, 2, 1, 1, 1, 1, 1, 1, 1, 2, 1,
		2, 2, 2, 2, 2, 2, 0, 2, 1, 1, 1, 1, 1, 1, 2, 1,
		1, 1, 1, 1, 1, 1, 2, 1, 1, 1, 1, 1, 1, 1, 2, 1,
		1, 1, 1, 1, 1, 1, 2, 1, 1, 1, 1, 1, 1, 1, 2, 1,
		1, 1, 1, 1, 1, 1, 2, 1, 1, 1, 1, 1, 1, 1, 2, 1,
		1, 1, 1, 1, 1, 1, 2, 1, 1, 1, 1, 1, 1, 1, 2, 1,
		2, 3, 3, 4, 3, 4, 2, 4, 2, 4, 3, 0, 3, 6, 2, 4,
		2, 3, 3, 0, 3, 4, 2, 4, 2, 4, 3, 0, 3, 0, 2, 4,
		3, 3, 2, 0, 0, 4, 2, 4, 4, 1, 4, 0, 0, 0, 2, 4,
		3, 3, 2, 1, 0, 4, 2, 4, 3, 2, 4, 1, 0, 0, 2, 4,
	}

	CB_CYCLES = []int{
		2, 2, 2, 2, 2, 2, 4, 2, 2, 2, 2, 2, 2, 2, 4, 2,
		2, 2, 2, 2, 2, 2, 4, 2, 2, 2, 2, 2, 2, 2, 4, 2,
		2, 2, 2, 2, 2, 2, 4, 2, 2, 2, 2, 2, 2, 2, 4, 2,
		2, 2, 2, 2, 2, 2, 4, 2, 2, 2, 2, 2, 2, 2, 4, 2,
		2, 2, 2, 2, 2, 2, 3, 2, 2, 2, 2, 2, 2, 2, 3, 2,
		2, 2, 2, 2, 2, 2, 3, 2, 2, 2, 2, 2, 2, 2, 3, 2,
		2, 2, 2, 2, 2, 2, 3, 2, 2, 2, 2, 2, 2, 2, 3, 2,
		2, 2, 2, 2, 2, 2, 3, 2, 2, 2, 2, 2, 2, 2, 3, 2,
		2, 2, 2, 2, 2, 2, 4, 2, 2, 2, 2, 2, 2, 2, 4, 2,
		2, 2, 2, 2, 2, 2, 4, 2, 2, 2, 2, 2, 2, 2, 4, 2,
		2, 2, 2, 2, 2, 2, 4, 2, 2, 2, 2, 2, 2, 2, 4, 2,
		2, 2, 2, 2, 2, 2, 4, 2, 2, 2, 2, 2, 2, 2, 4, 2,
		2, 2, 2, 2, 2, 2, 4, 2, 2, 2, 2, 2, 2, 2, 4, 2,
		2, 2, 2, 2, 2, 2, 4, 2, 2, 2, 2, 2, 2, 2, 4, 2,
		2, 2, 2, 2, 2, 2, 4, 2, 2, 2, 2, 2, 2, 2, 4, 2,
		2, 2, 2, 2, 2, 2, 4, 2, 2, 2, 2, 2, 2, 2, 4, 2,
	}

	INSTRUCTIONS = map[uint8]func(*CPU){
		0x00: nop,
		0x01: ldBC16,
		0x02: ldBCA,
		0x03: incBC,
		0x04: incB,
		0x05: decB,
		0x06: ldiB,
		0x07: unimplemented,
		0x08: ld16SP,
		0x09: addHLBC,
		0x0A: ldABC,
		0x0B: decBC,
		0x0C: incC,
		0x0D: decC,
		0x0E: ldiC,
		0x0F: unimplemented,
		0x10: unimplemented,
		0x11: ldDE16,
		0x12: ldDEA,
		0x13: incDE,
		0x14: incD,
		0x15: decD,
		0x16: ldiD,
		0x17: unimplemented,
		0x18: jr,
		0x19: addHLDE,
		0x1A: ldADE,
		0x1B: decDE,
		0x1C: incE,
		0x1D: decE,
		0x1E: ldiE,
		0x1F: unimplemented,
		0x20: jrnz,
		0x21: ldHL16,
		0x22: ldHLIA,
		0x23: incHL16,
		0x24: incH,
		0x25: decH,
		0x26: ldiH,
		0x27: unimplemented,
		0x28: jrz,
		0x29: addHLHL,
		0x2A: ldAHLI,
		0x2B: decHL16,
		0x2C: incL,
		0x2D: decL,
		0x2E: ldiL,
		0x2F: cpl,
		0x30: jrnc,
		0x31: ldSP16,
		0x32: ldHLDA,
		0x33: incSP,
		0x34: incHL,
		0x35: decHL,
		0x36: ldiHL,
		0x37: unimplemented,
		0x38: jrc,
		0x39: addHLSP,
		0x3A: ldAHLD,
		0x3B: decSP,
		0x3C: incA,
		0x3D: decA,
		0x3E: ldiA,
		0x3F: unimplemented,
		0x40: ldBB,
		0x41: ldBC,
		0x42: ldBD,
		0x43: ldBE,
		0x44: ldBH,
		0x45: ldBL,
		0x46: ldBHL,
		0x47: ldBA,
		0x48: ldCB,
		0x49: ldCC,
		0x4A: ldCD,
		0x4B: ldCE,
		0x4C: ldCH,
		0x4D: ldCL,
		0x4E: ldCHL,
		0x4F: ldCA,
		0x50: ldDB,
		0x51: ldDC,
		0x52: ldDD,
		0x53: ldDE,
		0x54: ldDH,
		0x55: ldDL,
		0x56: ldDHL,
		0x57: ldDA,
		0x58: ldEB,
		0x59: ldEC,
		0x5A: ldED,
		0x5B: ldEE,
		0x5C: ldEH,
		0x5D: ldEL,
		0x5E: ldEHL,
		0x5F: ldEA,
		0x60: ldHB,
		0x61: ldHC,
		0x62: ldHD,
		0x63: ldHE,
		0x64: ldHH,
		0x65: ldHL,
		0x66: ldHHL,
		0x67: ldHA,
		0x68: ldLB,
		0x69: ldLC,
		0x6A: ldLD,
		0x6B: ldLE,
		0x6C: ldLH,
		0x6D: ldLL,
		0x6E: ldLHL,
		0x6F: ldLA,
		0x70: ldHLB,
		0x71: ldHLC,
		0x72: ldHLD,
		0x73: ldHLE,
		0x74: ldHLH,
		0x75: ldHLL,
		0x76: unimplemented,
		0x77: ldHLA,
		0x78: ldAB,
		0x79: ldAC,
		0x7A: ldAD,
		0x7B: ldAE,
		0x7C: ldAH,
		0x7D: ldAL,
		0x7E: ldAHL,
		0x7F: ldAA,
		0x80: addAB,
		0x81: addAC,
		0x82: addAD,
		0x83: addAE,
		0x84: addAH,
		0x85: addAL,
		0x86: addAHL,
		0x87: addAA,
		0x88: adcAB,
		0x89: adcAC,
		0x8A: adcAD,
		0x8B: adcAE,
		0x8C: adcAH,
		0x8D: adcAL,
		0x8E: adcAHL,
		0x8F: adcAA,
		0x90: subAB,
		0x91: subAC,
		0x92: subAD,
		0x93: subAE,
		0x94: subAH,
		0x95: subAL,
		0x96: subAHL,
		0x97: subAA,
		0x98: sbcAB,
		0x99: sbcAC,
		0x9A: sbcAD,
		0x9B: sbcAE,
		0x9C: sbcAH,
		0x9D: sbcAL,
		0x9E: sbcAHL,
		0x9F: sbcAA,
		0xA0: andAB,
		0xA1: andAC,
		0xA2: andAD,
		0xA3: andAE,
		0xA4: andAH,
		0xA5: andAL,
		0xA6: andAHL,
		0xA7: andAA,
		0xA8: xorAB,
		0xA9: xorAC,
		0xAA: xorAD,
		0xAB: xorAE,
		0xAC: xorAH,
		0xAD: xorAL,
		0xAE: xorAHL,
		0xAF: xorAA,
		0xB0: orAB,
		0xB1: orAC,
		0xB2: orAD,
		0xB3: orAE,
		0xB4: orAH,
		0xB5: orAL,
		0xB6: orAHL,
		0xB7: orAA,
		0xB8: cpAB,
		0xB9: cpAC,
		0xBA: cpAD,
		0xBB: cpAE,
		0xBC: cpAH,
		0xBD: cpAL,
		0xBE: cpAHL,
		0xBF: cpAA,
		0xC0: retnz,
		0xC1: unimplemented,
		0xC2: jpnz,
		0xC3: jp,
		0xC4: callnz,
		0xC5: unimplemented,
		0xC6: adi,
		0xC7: unimplemented,
		0xC8: retz,
		0xC9: ret,
		0xCA: jpz,
		0xCB: unimplemented,
		0xCC: callz,
		0xCD: call,
		0xCE: aci,
		0xCF: unimplemented,
		0xD0: retnc,
		0xD1: unimplemented,
		0xD2: jpnc,
		0xD3: unimplemented,
		0xD4: callnc,
		0xD5: unimplemented,
		0xD6: sui,
		0xD7: unimplemented,
		0xD8: retc,
		0xD9: reti,
		0xDA: jpc,
		0xDB: unimplemented,
		0xDC: callc,
		0xDD: unimplemented,
		0xDE: sbi,
		0xDF: unimplemented,
		0xE0: ldh16A,
		0xE1: unimplemented,
		0xE2: ldhCA,
		0xE3: unimplemented,
		0xE4: unimplemented,
		0xE5: unimplemented,
		0xE6: ani,
		0xE7: unimplemented,
		0xE8: addSP,
		0xE9: jpHL,
		0xEA: ld16A,
		0xEB: unimplemented,
		0xEC: unimplemented,
		0xED: unimplemented,
		0xEE: xri,
		0xEF: unimplemented,
		0xF0: ldhA16,
		0xF1: unimplemented,
		0xF2: ldhAC,
		0xF3: di,
		0xF4: unimplemented,
		0xF5: unimplemented,
		0xF6: ori,
		0xF7: unimplemented,
		0xF8: ldHLSP,
		0xF9: ldSPHL,
		0xFA: ldA16,
		0xFB: ei,
		0xFC: unimplemented,
		0xFD: unimplemented,
		0xFE: cpi,
		0xFF: unimplemented,
	}

	CB_INSTRUCTIONS = map[uint8]func(*CPU){
		0x00: unimplemented,
		0x01: unimplemented,
		0x02: unimplemented,
		0x03: unimplemented,
		0x04: unimplemented,
		0x05: unimplemented,
		0x06: unimplemented,
		0x07: unimplemented,
		0x08: unimplemented,
		0x09: unimplemented,
		0x0A: unimplemented,
		0x0B: unimplemented,
		0x0C: unimplemented,
		0x0D: unimplemented,
		0x0E: unimplemented,
		0x0F: unimplemented,
		0x10: unimplemented,
		0x11: unimplemented,
		0x12: unimplemented,
		0x13: unimplemented,
		0x14: unimplemented,
		0x15: unimplemented,
		0x16: unimplemented,
		0x17: unimplemented,
		0x18: unimplemented,
		0x19: unimplemented,
		0x1A: unimplemented,
		0x1B: unimplemented,
		0x1C: unimplemented,
		0x1D: unimplemented,
		0x1E: unimplemented,
		0x1F: unimplemented,
		0x20: unimplemented,
		0x21: unimplemented,
		0x22: unimplemented,
		0x23: unimplemented,
		0x24: unimplemented,
		0x25: unimplemented,
		0x26: unimplemented,
		0x27: unimplemented,
		0x28: unimplemented,
		0x29: unimplemented,
		0x2A: unimplemented,
		0x2B: unimplemented,
		0x2C: unimplemented,
		0x2D: unimplemented,
		0x2E: unimplemented,
		0x2F: unimplemented,
		0x30: unimplemented,
		0x31: unimplemented,
		0x32: unimplemented,
		0x33: unimplemented,
		0x34: unimplemented,
		0x35: unimplemented,
		0x36: unimplemented,
		0x37: unimplemented,
		0x38: unimplemented,
		0x39: unimplemented,
		0x3A: unimplemented,
		0x3B: unimplemented,
		0x3C: unimplemented,
		0x3D: unimplemented,
		0x3E: unimplemented,
		0x3F: unimplemented,
		0x40: bit0B,
		0x41: bit0C,
		0x42: bit0D,
		0x43: bit0E,
		0x44: bit0H,
		0x45: bit0L,
		0x46: bit0HL,
		0x47: bit0A,
		0x48: bit1B,
		0x49: bit1C,
		0x4A: bit1D,
		0x4B: bit1E,
		0x4C: bit1H,
		0x4D: bit1L,
		0x4E: bit1HL,
		0x4F: bit1A,
		0x50: bit2B,
		0x51: bit2C,
		0x52: bit2D,
		0x53: bit2E,
		0x54: bit2H,
		0x55: bit2L,
		0x56: bit2HL,
		0x57: bit2A,
		0x58: bit3B,
		0x59: bit3C,
		0x5A: bit3D,
		0x5B: bit3E,
		0x5C: bit3H,
		0x5D: bit3L,
		0x5E: bit3HL,
		0x5F: bit3A,
		0x60: bit4B,
		0x61: bit4C,
		0x62: bit4D,
		0x63: bit4E,
		0x64: bit4H,
		0x65: bit4L,
		0x66: bit4HL,
		0x67: bit4A,
		0x68: bit5B,
		0x69: bit5C,
		0x6A: bit5D,
		0x6B: bit5E,
		0x6C: bit5H,
		0x6D: bit5L,
		0x6E: bit5HL,
		0x6F: bit5A,
		0x70: bit6B,
		0x71: bit6C,
		0x72: bit6D,
		0x73: bit6E,
		0x74: bit6H,
		0x75: bit6L,
		0x76: bit6HL,
		0x77: bit6A,
		0x78: bit7B,
		0x79: bit7C,
		0x7A: bit7D,
		0x7B: bit7E,
		0x7C: bit7H,
		0x7D: bit7L,
		0x7E: bit7HL,
		0x7F: bit7A,
		0x80: unimplemented,
		0x81: unimplemented,
		0x82: unimplemented,
		0x83: unimplemented,
		0x84: unimplemented,
		0x85: unimplemented,
		0x86: unimplemented,
		0x87: unimplemented,
		0x88: unimplemented,
		0x89: unimplemented,
		0x8A: unimplemented,
		0x8B: unimplemented,
		0x8C: unimplemented,
		0x8D: unimplemented,
		0x8E: unimplemented,
		0x8F: unimplemented,
		0x90: unimplemented,
		0x91: unimplemented,
		0x92: unimplemented,
		0x93: unimplemented,
		0x94: unimplemented,
		0x95: unimplemented,
		0x96: unimplemented,
		0x97: unimplemented,
		0x98: unimplemented,
		0x99: unimplemented,
		0x9A: unimplemented,
		0x9B: unimplemented,
		0x9C: unimplemented,
		0x9D: unimplemented,
		0x9E: unimplemented,
		0x9F: unimplemented,
		0xA0: unimplemented,
		0xA1: unimplemented,
		0xA2: unimplemented,
		0xA3: unimplemented,
		0xA4: unimplemented,
		0xA5: unimplemented,
		0xA6: unimplemented,
		0xA7: unimplemented,
		0xA8: unimplemented,
		0xA9: unimplemented,
		0xAA: unimplemented,
		0xAB: unimplemented,
		0xAC: unimplemented,
		0xAD: unimplemented,
		0xAE: unimplemented,
		0xAF: unimplemented,
		0xB0: unimplemented,
		0xB1: unimplemented,
		0xB2: unimplemented,
		0xB3: unimplemented,
		0xB4: unimplemented,
		0xB5: unimplemented,
		0xB6: unimplemented,
		0xB7: unimplemented,
		0xB8: unimplemented,
		0xB9: unimplemented,
		0xBA: unimplemented,
		0xBB: unimplemented,
		0xBC: unimplemented,
		0xBD: unimplemented,
		0xBE: unimplemented,
		0xBF: unimplemented,
		0xC0: unimplemented,
		0xC1: unimplemented,
		0xC2: unimplemented,
		0xC3: unimplemented,
		0xC4: unimplemented,
		0xC5: unimplemented,
		0xC6: unimplemented,
		0xC7: unimplemented,
		0xC8: unimplemented,
		0xC9: unimplemented,
		0xCA: unimplemented,
		0xCB: unimplemented,
		0xCC: unimplemented,
		0xCD: unimplemented,
		0xCE: unimplemented,
		0xCF: unimplemented,
		0xD0: unimplemented,
		0xD1: unimplemented,
		0xD2: unimplemented,
		0xD3: unimplemented,
		0xD4: unimplemented,
		0xD5: unimplemented,
		0xD6: unimplemented,
		0xD7: unimplemented,
		0xD8: unimplemented,
		0xD9: unimplemented,
		0xDA: unimplemented,
		0xDB: unimplemented,
		0xDC: unimplemented,
		0xDD: unimplemented,
		0xDE: unimplemented,
		0xDF: unimplemented,
		0xE0: unimplemented,
		0xE1: unimplemented,
		0xE2: unimplemented,
		0xE3: unimplemented,
		0xE4: unimplemented,
		0xE5: unimplemented,
		0xE6: unimplemented,
		0xE7: unimplemented,
		0xE8: unimplemented,
		0xE9: unimplemented,
		0xEA: unimplemented,
		0xEB: unimplemented,
		0xEC: unimplemented,
		0xED: unimplemented,
		0xEE: unimplemented,
		0xEF: unimplemented,
		0xF0: unimplemented,
		0xF1: unimplemented,
		0xF2: unimplemented,
		0xF3: unimplemented,
		0xF4: unimplemented,
		0xF5: unimplemented,
		0xF6: unimplemented,
		0xF7: unimplemented,
		0xF8: unimplemented,
		0xF9: unimplemented,
		0xFA: unimplemented,
		0xFB: unimplemented,
		0xFC: unimplemented,
		0xFD: unimplemented,
		0xFE: unimplemented,
		0xFF: unimplemented,
	}
)
