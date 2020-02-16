package cdgo

import (
	"fmt"
	"io"
)

// File is a wrapper around byte array to help parsing a CDG file
type File struct {
	Reader io.Reader

	buffer []byte
}

// NewFile given a byte array will return a new file
func NewFile(reader io.Reader) *File {
	return &File{
		Reader: reader,
		buffer: make([]byte, PacketSize),
	}
}

// NextPacket will return the next cdg packet or return an error
func (f File) NextPacket() (*Packet, error) {
	packet := &Packet{}
	err := f.ReadInto(packet)
	if err != nil {
		return nil, err
	}

	return packet, nil
}

// ReadInto take a packet and replace its contents with the contents of the next
// packet or returns error
func (f File) ReadInto(packet *Packet) error {
	n, err := io.ReadFull(f.Reader, f.buffer)
	if err != nil {
		return err
	}

	if n != PacketSize {
		return fmt.Errorf("Incorrect bytes read %d should be %d", n, PacketSize)
	}

	packet.Command = f.buffer[0]
	packet.Instruction = f.buffer[1]
	packet.ParityQ = f.buffer[2:4]
	packet.Data = f.buffer[4:20] // lol 4:20
	packet.ParityP = f.buffer[20:24]

	return nil
}
