# cdgo

CDG library for Go (Golang)

### Example

The example below will load a file and print the first instruction packet

```go
file, _ := os.Open("path/to/cdg_file.cdg")
cdgFile := cdgo.NewFile(file)

packet, _ := cdgFile.NextPacket()

packet.Print()
```

### `cmd`

#### cdg2images: Write each frame of a CDG file

```shell
go run cmd/cdg2images/main.go test.cdg
```

#### cdg2text: Debugging tool. Will print human readable instructions of a CDG


```shell
go run cmd/cdg2text/main.go test.cdg
```

### References:

https://jbum.com/cdg_revealed.html

https://goughlui.com/2019/03/31/tech-flashback-the-cdgraphics-format-cdg/

