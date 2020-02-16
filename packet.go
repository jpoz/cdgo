package cdgo

import (
	"fmt"
	"image/color"
)

// PacketSize is the size of each cdg packet
const PacketSize = 24

const screenWidth = 300
const screenHeight = 216

const visibleWidth = 288
const visibleHeight = 192

const fontWidth = 6
const fontHeight = 12

var fontMasks = []byte{
	0x01,
	0x02,
	0x04,
	0x08,
	0x10,
	0x20,
}

// CDGMask is a mask for each byte instruction in a package
const CDGMask = 0x3F
const cdgCommand = 0x09
const b4 = 0x0F

// Instructions

// MemoryPresetCode is the instruction code for a memeory preset
const MemoryPresetCode = 1

// BorderPresetCode is the instruction code for a border preset
const BorderPresetCode = 2

// TileBlockCode is the instruction code for a tile block
const TileBlockCode = 6

// TileBlockXorCode is the instruction code for a tile block
const TileBlockXorCode = 38

// ScrollPresetCode is the instruction code for a scroll preset
const ScrollPresetCode = 20

// ScrollCopyCode is the instruction code for a scroll copy
const ScrollCopyCode = 24

// TransparentColorCode is the instruction code for a transpaarent color
const TransparentColorCode = 28

// LoadColorTableLowCode is the instruction code to load the first 8 colors
const LoadColorTableLowCode = 30

// LoadColorTableHightCode is the instruction code to load the last 8 colors
const LoadColorTableHightCode = 31

// Packet is 24 bytes of the cdg file
type Packet struct {
	Command     byte
	Instruction byte
	ParityQ     []byte
	Data        []byte
	ParityP     []byte
}

// MemoryPreset represents a memory preset package information
type MemoryPreset struct {
	Color  uint8
	Repeat uint8
}

// ColorTable represents a memory preset packge information
type ColorTable struct {
	Offset int // if it's low 0 high 8
	Colors []color.Color
}

// TileBlock is the data of a given tile
type TileBlock struct {
	Color0 uint8
	Color1 uint8
	Row    uint16
	Column uint16
	XOR    bool
	Pixels []uint8
}

// Scroll represents data for ScrollPreset and ScrollCopy
type Scroll struct {
	Copy    bool
	Color   uint8
	HScroll uint8
	VScroll uint8
	HSCmd   uint8
	HOffset uint8
	VSCmd   uint8
	VOffset uint8
}

// Print will print out a "pretty" versin of the bytes
func (p Packet) Print() {
	fmt.Printf("%08b %08b %x %b %x\n",
		p.Command,
		p.Instruction,
		p.ParityQ,
		p.Data,
		p.ParityP,
	)
}

// IsCgdCommand test if the packet is a cdg command
func (p Packet) IsCgdCommand() bool {
	return p.Command&CDGMask == cdgCommand
}

// InstructionCode returns the cdg instruction code
func (p Packet) InstructionCode() byte {
	return p.Instruction & CDGMask
}

// MemoryPreset returns the data as a memory preset
func (p Packet) MemoryPreset() *MemoryPreset {
	return &MemoryPreset{
		Color:  p.Data[0] & b4,
		Repeat: p.Data[1] & b4,
	}
}

// LoadColorTable will load the colors and return the offset
func (p Packet) LoadColorTable() *ColorTable {
	colorTable := &ColorTable{
		Offset: 0,
		Colors: make([]color.Color, 8),
	}

	// TODO might want to remove this
	if p.InstructionCode() == LoadColorTableHightCode {
		colorTable.Offset = 8
	}

	for i := 0; i < 8; i++ {
		var ColorEntry int
		// fmt.Printf("%08b %08b\n", p.Data[0], p.Data[1])
		ColorEntry = int(p.Data[2*i]&CDGMask) << 8
		// fmt.Printf("%08b R \n", ColorEntry)
		ColorEntry = ColorEntry + int(p.Data[(2*i)+1]&CDGMask)
		// fmt.Printf("%08b GB\n", ColorEntry)
		ColorEntry = ((ColorEntry & 0x3F00) >> 2) | (ColorEntry & 0x003F)

		var c color.RGBA
		r := ColorEntry & 0xF00 >> 8
		// fmt.Printf("R %08b \n", r)
		// fmt.Printf("R %d \n", uint(r)*17)

		g := ColorEntry & 0x0F0 >> 4
		// fmt.Printf("G %08b \n", g)
		// fmt.Printf("G %d \n", uint(g)*17)

		b := ColorEntry & 0x00F
		// fmt.Printf("B %08b \n", g)
		// fmt.Printf("B %d \n", uint(g)*17)

		c.R = uint8(r) * 17
		c.G = uint8(g) * 17
		c.B = uint8(b) * 17
		c.A = 255

		colorTable.Colors[i] = c
	}

	return colorTable
}

// TileBlock will give back a tileblock with indexed colors
func (p Packet) TileBlock(xor bool) *TileBlock {
	tileBlock := &TileBlock{
		Color0: p.Data[0] & CDGMask,
		Color1: p.Data[1] & CDGMask,
		Row:    uint16(p.Data[2]&0x1F) * fontHeight,
		Column: uint16(p.Data[3]&0x3F) * fontWidth,
		Pixels: make([]uint8, fontWidth*fontHeight),
	}

	var y uint16
	for y = 0; y < fontHeight; y++ {
		r := p.Data[y+4] & 0x3F
		var x uint16
		for x = 0; x < fontWidth; x++ {
			c := r & fontMasks[5-x] >> (5 - x)
			// fmt.Printf("%x", c)
			tileBlock.Pixels[x+y*fontWidth] = c
		}
		// fmt.Printf(" - %b\n", r)
	}

	return tileBlock
}

// ScrollPreset will return the scoll data from the packet
func (p Packet) ScrollPreset() *Scroll {
	hScroll := p.Data[1] & 0x3F
	vScroll := p.Data[2] & 0x3F

	return &Scroll{
		Copy:    false,
		Color:   p.Data[0] & 0x0F,
		HScroll: hScroll,
		VScroll: vScroll,
		HSCmd:   (hScroll & 0x30) >> 4,
		HOffset: (hScroll & 0x07),
		VSCmd:   (vScroll & 0x30) >> 4,
		VOffset: (vScroll & 0x0F),
	}
}

// ScrollCopy will return the scoll data from the packet
func (p Packet) ScrollCopy() *Scroll {
	hScroll := p.Data[1] & 0x3F
	vScroll := p.Data[2] & 0x3F

	return &Scroll{
		Copy:    true,
		Color:   p.Data[0] & 0x0F,
		HScroll: hScroll,
		VScroll: vScroll,
		HSCmd:   (hScroll & 0x30) >> 4,
		HOffset: (hScroll & 0x07),
		VSCmd:   (vScroll & 0x30) >> 4,
		VOffset: (vScroll & 0x0F),
	}
}
