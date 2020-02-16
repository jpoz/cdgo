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

	screen := cdgo.NewScreen()

	packetCount := 0

	for {
		packet, err := cdgFile.NextPacket()
		if err != nil {
			// TODO if EOF just break
			log.Fatal(err)
			break
		}

		if packet.IsCgdCommand() {
			err = screen.UpdateWithPacket(packet)
			if err != nil {
				log.Fatal(err)
				break
			}

			if screen.Dirty {
				// TODO make dir before writing JPEGs
				err := screen.WriteJPEG(fmt.Sprintf("images/image%07d.jpg", packetCount))
				if err != nil {
					log.Fatal(err)
				}

				screen.Dirty = false
			}

			// if packetCount > 10000 {
			// 	log.Fatal("DONE")

			// 	break
			// }

			packetCount++
		}
	}

}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
