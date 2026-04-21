# mc1.8.x-bin

A Golang library for manipulating binary streams specific to the Minecraft 1.8.x protocol.
Designed for high performance and ease of use when creating servers or proxies.

### Features
* **VarInt & VarLong:** Full implementation according to the Minecraft protocol specification.
* **String Handling:** Reads and writes UTF-8 strings with a length prefix.
* **Low Level:** Runs on top of `io.Reader` and `io.Writer` (compatible with TCP Conn, File, or Buffer).

### Installation
```bash
go get [github.com/Kartendsy/mc1.8.x-bin](https://github.com/Kartendsy/mc1.8.x-bin)```

### 4. Usage Example (Quick Start)
This is the most important part. Provide a code example that can be directly copied and run.
```markdown
### 🛠 Usage Example

```go
package main

import (
"bytes"
"fmt"
"[github.com/Kartendsy/mc1.8.x-bin/bin](https://github.com/Kartendsy/mc1.8.x-bin/bin)" // Adjust to your package folder
)

func main() {
buf := new(bytes.Buffer)
writer := bin.NewWriter(buf)

// Write a VarInt
writer.WriteVarInt(128)

fmt.Printf("Encoded: %x\n", buf.Bytes())
}```

[![Go Reference](https://pkg.go.dev/badge/github.com/Kartendsy/mc1.8.x-bin.svg)](https://pkg.go.dev/github.com/Kartendsy/mc1.8.x-bin)
