package cdgo

import "testing"

type CdgCommandTest struct {
	Command      byte
	IsCgdCommand bool
}

var CdgCommandTests = []CdgCommandTest{
	CdgCommandTest{
		0b1001, true,
	},
	CdgCommandTest{
		0b1000, false,
	},
	CdgCommandTest{
		0b01001001, true,
	},
}

func TestIsCgdCommand(t *testing.T) {
	for i, test := range CdgCommandTests {
		p := Packet{
			Command:     test.Command,
			Instruction: 0x02,
			ParityQ:     []byte{},
			Data:        []byte{},
			ParityP:     []byte{},
		}

		if p.IsCgdCommand() != test.IsCgdCommand {
			t.Errorf(
				"TestIsCgdCommand: Expected test %d, %X to return %t form IsCdgCommand but it returned %t",
				i,
				p.Command,
				test.IsCgdCommand,
				p.IsCgdCommand(),
			)
		}
	}
}

type InstructionTest struct {
	InstructionByte byte
	InstructionCode byte
}

var InstructionTests = []InstructionTest{
	InstructionTest{
		0b1001, 0x9,
	},
	InstructionTest{
		0b1000, 0x8,
	},
	InstructionTest{
		0b01001001, 0x9,
	},
}

func TestInstructionCode(t *testing.T) {
	for i, test := range InstructionTests {
		p := Packet{
			Command:     0x09,
			Instruction: test.InstructionByte,
			ParityQ:     []byte{},
			Data:        []byte{},
			ParityP:     []byte{},
		}

		if p.InstructionCode() != test.InstructionCode {
			t.Errorf(
				"TestInstructionCode: Expected test %d, %X to return %X form Instruction but it returned %X",
				i,
				p.Instruction,
				test.InstructionCode,
				p.InstructionCode(),
			)
		}
	}
}
