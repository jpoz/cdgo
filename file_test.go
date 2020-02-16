package cdgo

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

var testBytes = []byte{
	0x001, 0x002, 0x003, 0x004, 0x005, 0x006, 0x007,
	0x008, 0x009, 0x010, 0x011, 0x012, 0x013, 0x014,
	0x015, 0x016, 0x017, 0x018, 0x019, 0x020, 0x021,
	0x022, 0x023, 0x024, 0x025, 0x026, 0x027, 0x028,
	0x029, 0x030, 0x031, 0x032, 0x033, 0x034, 0x035,
	0x036, 0x037, 0x038, 0x039, 0x040, 0x041, 0x042,
}

func TestNewFile(t *testing.T) {
	reader := strings.NewReader("cdg file")
	file := NewFile(reader)
	if file.Reader != reader {
		t.Error("Reader is not the reader passed into the function")
	}

	if len(file.buffer) != PacketSize {
		t.Errorf(
			"Initialized buffer was incorrect length. Expected %d but was %d",
			PacketSize,
			len(file.buffer),
		)
	}
}

func TestNextPacket(t *testing.T) {
	reader := bytes.NewReader(testBytes)
	file := NewFile(reader)

	packet, err := file.NextPacket()
	if err != nil {
		t.Fatalf("Error wasn't nil: %s", err)
	}
	if packet.Command != 0x01 {
		t.Errorf(
			"Packet Command should have been %X but was %X",
			0x01,
			packet.Command,
		)
	}
	if packet.Instruction != 0x02 {
		t.Errorf(
			"Packet Instruction should have been %x but was %x",
			0x02,
			packet.Instruction,
		)
	}
	expected := []byte{
		0x005, 0x006, 0x007, 0x008, 0x009, 0x010,
		0x011, 0x012, 0x013, 0x014, 0x015, 0x016,
		0x017, 0x018, 0x019, 0x020, 0x021, 0x022,
		0x023, 0x024,
	}
	for i := range packet.Data {
		if packet.Data[i] != expected[i] {
			t.Errorf(
				"Packet Data[%d] should have been %x but was %x",
				i,
				expected[i],
				packet.Data[i],
			)
		}
	}

	// If 24 bytes aren't left and unexpected EOF will be returned
	_, err = file.NextPacket()
	if err != io.ErrUnexpectedEOF {
		t.Fatalf("Error  %s", err)
	}
}
