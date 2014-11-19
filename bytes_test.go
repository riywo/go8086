package go8086

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var bytesmodRMTests = []struct {
	mod  Mod
	reg  Reg
	rm   RM
	w    Bit
	sreg *SegmentRegister
	out  string
	num  int
}{
	{Mod00, Reg000, RM000, Bit8, nil, "[bx+si],al", 1},
	{Mod00, Reg001, RM001, Bit8, nil, "[bx+di],cl", 1},
	{Mod00, Reg010, RM010, Bit8, nil, "[bp+si],dl", 1},
	{Mod00, Reg011, RM011, Bit8, nil, "[bp+di],bl", 1},
	{Mod00, Reg100, RM100, Bit8, nil, "[si],ah", 1},
	{Mod00, Reg101, RM101, Bit8, nil, "[di],ch", 1},
	{Mod00, Reg110, RM110, Bit8, nil, "[0x12ff],dh", 3},
	{Mod00, Reg111, RM111, Bit8, nil, "[bx],bh", 1},
	{Mod00, Reg000, RM000, Bit16, nil, "[bx+si],ax", 1},
	{Mod00, Reg001, RM001, Bit16, nil, "[bx+di],cx", 1},
	{Mod00, Reg010, RM010, Bit16, nil, "[bp+si],dx", 1},
	{Mod00, Reg011, RM011, Bit16, nil, "[bp+di],bx", 1},
	{Mod00, Reg100, RM100, Bit16, nil, "[si],sp", 1},
	{Mod00, Reg101, RM101, Bit16, nil, "[di],bp", 1},
	{Mod00, Reg110, RM110, Bit16, nil, "[0x12ff],si", 3},
	{Mod00, Reg111, RM111, Bit16, nil, "[bx],di", 1},
	{Mod01, Reg000, RM000, Bit16, nil, "[bx+si-0x1],ax", 2},
	{Mod01, Reg001, RM001, Bit16, nil, "[bx+di-0x1],cx", 2},
	{Mod01, Reg010, RM010, Bit16, nil, "[bp+si-0x1],dx", 2},
	{Mod01, Reg011, RM011, Bit16, nil, "[bp+di-0x1],bx", 2},
	{Mod01, Reg100, RM100, Bit16, nil, "[si-0x1],sp", 2},
	{Mod01, Reg101, RM101, Bit16, nil, "[di-0x1],bp", 2},
	{Mod01, Reg110, RM110, Bit16, nil, "[bp-0x1],si", 2},
	{Mod01, Reg111, RM111, Bit16, nil, "[bx-0x1],di", 2},
	{Mod10, Reg000, RM000, Bit8, nil, "[bx+si+0x12ff],al", 3},
	{Mod10, Reg001, RM001, Bit8, nil, "[bx+di+0x12ff],cl", 3},
	{Mod10, Reg010, RM010, Bit8, nil, "[bp+si+0x12ff],dl", 3},
	{Mod10, Reg011, RM011, Bit8, nil, "[bp+di+0x12ff],bl", 3},
	{Mod10, Reg100, RM100, Bit8, nil, "[si+0x12ff],ah", 3},
	{Mod10, Reg101, RM101, Bit8, nil, "[di+0x12ff],ch", 3},
	{Mod10, Reg110, RM110, Bit8, nil, "[bp+0x12ff],dh", 3},
	{Mod10, Reg111, RM111, Bit8, nil, "[bx+0x12ff],bh", 3},
	{Mod11, Reg000, RM111, Bit8, nil, "bh,al", 1},
	{Mod11, Reg001, RM110, Bit8, nil, "dh,cl", 1},
	{Mod11, Reg010, RM101, Bit8, nil, "ch,dl", 1},
	{Mod11, Reg011, RM100, Bit8, nil, "ah,bl", 1},
	{Mod11, Reg100, RM011, Bit16, nil, "bx,sp", 1},
	{Mod11, Reg101, RM010, Bit16, nil, "dx,bp", 1},
	{Mod11, Reg110, RM001, Bit16, nil, "cx,si", 1},
	{Mod11, Reg111, RM000, Bit16, nil, "ax,di", 1},
	{Mod00, Reg000, RM000, Bit8, ES, "[es:bx+si],al", 1},
	{Mod00, Reg110, RM110, Bit8, CS, "[cs:0x12ff],dh", 3},
	{Mod01, Reg000, RM000, Bit16, SS, "[ss:bx+si-0x1],ax", 2},
	{Mod10, Reg111, RM111, Bit8, DS, "[ds:bx+0x12ff],bh", 3},
}

func TestBytesmodRM(t *testing.T) {
	bs := Bytes{0x00, 0xff, 0x12}
	for _, test := range bytesmodRMTests {
		opr1, opr2, readBytes := bs.modRM(test.mod, test.reg, test.rm, test.w, test.sreg)
		assert.Equal(t, test.out, opr1.Disasm()+","+opr2.Disasm())
		assert.Equal(t, test.num, len(readBytes))
	}
}

var bytesGetOperandByModRMTests = []struct {
	bytes Bytes
	w     Bit
	sreg  *SegmentRegister
	out   string
	reg   Reg
}{
	{Bytes{0x00}, Bit8, nil, "[bx+si],al", Reg000},
	{Bytes{0x06, 0x34, 0x12}, Bit8, nil, "[0x1234],al", Reg000},
	{Bytes{0x94, 0x00, 0xff}, Bit8, nil, "[si-0x100],dl", Reg010},
	{Bytes{0x00}, Bit8, ES, "[es:bx+si],al", Reg000},
	{Bytes{0x06, 0x34, 0x12}, Bit8, CS, "[cs:0x1234],al", Reg000},
	{Bytes{0x94, 0x00, 0xff}, Bit8, SS, "[ss:si-0x100],dl", Reg010},
}

func TestBytesGetOperandByModRM(t *testing.T) {
	for _, test := range bytesGetOperandByModRMTests {
		reg, opr1, opr2, readBytes := test.bytes.GetOperandByModRM(test.w, test.sreg)
		assert.Equal(t, test.out, opr1.Disasm()+","+opr2.Disasm())
		assert.Equal(t, test.reg, reg)
		assert.Equal(t, test.bytes, readBytes)
	}
}

var bytesGetOperandByModRMDataTests = []struct {
	bytes  Bytes
	w      Bit
	signed Signed
	sreg   *SegmentRegister
	out    string
	reg    Reg
}{
	{Bytes{0x00, 0x00}, Bit8, Unsign, nil, "[bx+si],0x0", Reg000},
	{Bytes{0x00, 0xff, 0xff}, Bit16, Unsign, nil, "[bx+si],0xffff", Reg000},
	{Bytes{0x00, 0x00}, Bit16, Sign, nil, "[bx+si],+0x0", Reg000},
	{Bytes{0x00, 0x00}, Bit8, Unsign, ES, "[es:bx+si],0x0", Reg000},
	{Bytes{0x00, 0xff, 0xff}, Bit16, Unsign, CS, "[cs:bx+si],0xffff", Reg000},
	{Bytes{0x00, 0x00}, Bit16, Sign, SS, "[ss:bx+si],+0x0", Reg000},
}

func TestBytesGetOperandByModRMData(t *testing.T) {
	for _, test := range bytesGetOperandByModRMDataTests {
		reg, opr1, opr2, readBytes := test.bytes.GetOperandByModRMData(test.w, test.signed, test.sreg)
		assert.Equal(t, test.out, opr1.Disasm()+","+opr2.Disasm())
		assert.Equal(t, test.reg, reg)
		assert.Equal(t, test.bytes, readBytes)
	}
}

var bytesGetOperandOfAccImmTests = []struct {
	bytes Bytes
	w     Bit
	immw  Bit
	out   string
}{
	{Bytes{0x00}, Bit8, Bit8, "al,0x0"},
	{Bytes{0x34, 0xff}, Bit16, Bit16, "ax,0xff34"},
}

func TestBytesGetOperandOfAccImm(t *testing.T) {
	for _, test := range bytesGetOperandOfAccImmTests {
		opr1, opr2, readBytes := test.bytes.GetOperandOfAccImm(test.w, test.immw)
		assert.Equal(t, test.out, opr1.Disasm()+","+opr2.Disasm())
		assert.Equal(t, test.bytes, readBytes)
	}
}
