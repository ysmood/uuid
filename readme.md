# Overview

[![GoDoc](https://godoc.org/github.com/ysmood/uuid?status.svg)](https://pkg.go.dev/github.com/ysmood/uuid)

A distributed unique ID generator inspired by Twitter's Snowflake.

The key difference is it doesn't require a coordinator, it's completely stateless.

Total size is 128 bits, the format looks like this:

```text
24 bits for user defined namespace (3 chars)
56 bits for time in microsecond (from year 2020 to 2262)
16 bits for machine name (65536 machines)
32 bits for cryptographically secure noise (reasonable collision rate)
```

A sample anatomy for id `73736f-0009298b229ba5-6e01-011f8cdb` looks like this:

```text
namespace        time          machine     noise
    |             |               |          |
    v             v               v          v
   sso     2020-04-26T21:26...    n1      011f8cdb   # parsed format

 service       timestamp        Norway01   random    # meaning

  73736f     0009298b229ba5      6e01     011f8cdb   # hex format
```

The namespace is usually used to specify the application name. Such as use it to filter a specific app's log in elastic search so that you don't need an extra field to store the filter tag.

## Code Example

```go
package main

import (
    "fmt"

    "github.com/ysmood/uuid"
)

func main() {
    fmt.Println(uuid.New().Bin())
}
```
