package cdgo

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"os"
)

// Screen is a map of pixels to colors for a cdg
type Screen struct {
	Pixels   []uint8
	ColorMap []color.Color
	Dirty    bool
}

// NewScreen Will create a new screen with blank values
func NewScreen() *Screen {
	screen := &Screen{
		Pixels:   make([]uint8, visibleWidth*visibleHeight),
		ColorMap: make([]color.Color, 16),
	}

	for i := range screen.ColorMap {
		screen.ColorMap[i] = color.Black
	}

	return screen
}

// UpdateWithPacket takes a packet the screen should be updated with
func (s *Screen) UpdateWithPacket(packet *Packet) error {
	switch packet.InstructionCode() {
	case MemoryPresetCode:
		s.MemoryPreset(packet)
	case BorderPresetCode:
	case TileBlockCode:
		s.TileBlock(packet)
	case TileBlockXorCode:
		s.TileBlockXOR(packet)
	case ScrollPresetCode:
	case ScrollCopyCode:
	case TransparentColorCode:
	case LoadColorTableLowCode:
		s.LoadColorTable(packet)
	case LoadColorTableHightCode:
		s.LoadColorTable(packet)
	default:
		return fmt.Errorf("Packet is not a valid CDG command")
	}

	return nil
}

// MemoryPreset updates the screen with the memory preset package
func (s *Screen) MemoryPreset(packet *Packet) {
	mp := packet.MemoryPreset()

	for i := range s.Pixels {
		s.Pixels[i] = mp.Color
	}

	s.Dirty = true
}

// LoadColorTable updates the color map with the given packet
func (s *Screen) LoadColorTable(packet *Packet) {
	colorTable := packet.LoadColorTable()

	// TODO remake color output
	/// fmt.Printf("cdgLoadColorTable0..7\n")
	// fmt.Printf("  Color %d = 0x%X\n", colorIdx, ColorEntry)
	for i, color := range colorTable.Colors {
		s.ColorMap[i+colorTable.Offset] = color
	}

	s.Dirty = true
}

// TileBlock will update the screen with the given tile
func (s *Screen) TileBlock(packet *Packet) {
	tileBlock := packet.TileBlock(false)

	var x, y uint16
	for i, c := range tileBlock.Pixels {
		x = uint16(i) % fontWidth
		y = uint16(i) / fontWidth

		pixelIdx := ((tileBlock.Column + x) + ((tileBlock.Row + y) * visibleWidth))
		if pixelIdx > uint16(len(s.Pixels)-1) {
			// fmt.Println("TileBlock:", pixelIdx, "too bigg")
			continue
		}

		if c == 0x01 {
			s.Pixels[pixelIdx] = tileBlock.Color1
		} else {
			s.Pixels[pixelIdx] = tileBlock.Color0
		}
	}

	s.Dirty = true
}

// TileBlockXOR will update the screen with the given tile
func (s *Screen) TileBlockXOR(packet *Packet) {
	tileBlock := packet.TileBlock(true)

	var x, y uint16
	for i, c := range tileBlock.Pixels {
		x = uint16(i) % fontWidth
		y = uint16(i) / fontWidth

		pixelIdx := ((tileBlock.Column + x) + ((tileBlock.Row + y) * visibleWidth))
		if pixelIdx > uint16(len(s.Pixels)-1) {
			// fmt.Println("TileBlockXOR:", pixelIdx, "too bigg")
			continue
		}

		if c == 0x01 {
			s.Pixels[pixelIdx] ^= tileBlock.Color1
		} else {
			s.Pixels[pixelIdx] ^= tileBlock.Color0
		}
	}

	s.Dirty = true
}

// WriteJPEG writes a jpeg
func (s Screen) WriteJPEG(filename string) error {
	img := image.NewRGBA(image.Rect(0, 0, visibleWidth, visibleHeight))

	for x := 0; x < visibleWidth; x++ {
		for y := 0; y < visibleHeight; y++ {
			pixelIndex := x + (y * visibleWidth)

			if pixelIndex > (len(s.Pixels) - 1) {
				fmt.Println("Image index:", pixelIndex, "too bigg")
				continue
			}

			colorIdx := s.Pixels[pixelIndex]
			if colorIdx > uint8(len(s.ColorMap)-1) {
				fmt.Println(pixelIndex, " color too bigg")
				continue
			}

			img.Set(x, y, s.ColorMap[colorIdx])
		}
	}

	out, err := os.Create(filename)
	defer out.Close()
	if err != nil {
		return err
	}
	return jpeg.Encode(out, img, nil)
}
