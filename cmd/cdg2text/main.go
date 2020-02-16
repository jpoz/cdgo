package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jpoz/cdgo"
)

func main() {
	filepath := os.Args[1]

	file, err := os.Open(filepath)
	check(err)

	cdgFile := cdgo.NewFile(file)

	for {
		packet, err := cdgFile.NextPacket()
		if err != nil {
			log.Fatal(err)
			break
		}

		if packet.IsCgdCommand() {
			switch packet.InstructionCode() {
			case cdgo.MemoryPresetCode:
				mp := packet.MemoryPreset()
				fmt.Printf("cdgMemoryPresetCode [Color=%d, Repeat=%d]\n", mp.Color, mp.Repeat)
			case cdgo.BorderPresetCode:
				fmt.Printf("cdgBorderPreset TBD\n")
			case cdgo.TileBlockCode:
				tileBlock := packet.TileBlock(false)

				fmt.Printf(
					"cdgTileBlockNormal [Color0=%d, Color1=%d, ColIndex=%d, RowIndex=%d]\n",
					tileBlock.Color0,
					tileBlock.Color1,
					tileBlock.Column,
					tileBlock.Row,
				)
			case cdgo.TileBlockXorCode:
				tileBlock := packet.TileBlock(true)

				fmt.Printf(
					"cdgTileBlockXOR [Color0=%d, Color1=%d, ColIndex=%d, RowIndex=%d]\n",
					tileBlock.Color0,
					tileBlock.Color1,
					tileBlock.Column,
					tileBlock.Row,
				)
			case cdgo.ScrollPresetCode:
				scroll := packet.ScrollPreset()
				fmt.Printf(
					"cdgScrollPreset [color=%d, hSCmd=%d vSCmd=%d hOffset=%d, vOffset=%d]\n",
					scroll.Color,
					scroll.HSCmd,
					scroll.VSCmd,
					scroll.HOffset,
					scroll.VOffset,
				)
			case cdgo.ScrollCopyCode:
				scroll := packet.ScrollPreset()
				fmt.Printf(
					"cdgScrollCopy [color=%d, hSCmd=%d vSCmd=%d hOffset=%d, vOffset=%d]\n",
					scroll.Color,
					scroll.HSCmd,
					scroll.VSCmd,
					scroll.HOffset,
					scroll.VOffset,
				)
			case cdgo.TransparentColorCode:
				fmt.Printf("cdgDefineTransparentColor TBD\n")
			case cdgo.LoadColorTableLowCode:
			case cdgo.LoadColorTableHightCode:
			default:
				fmt.Println("Packet is not a valid CDG instruction")
			}
		} else {
			fmt.Println("Packet is not a valid CDG command")
		}
	}
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
