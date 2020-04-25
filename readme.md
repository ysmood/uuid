# Overview

[![GoDoc](https://godoc.org/github.com/ysmood/uuid?status.svg)](https://pkg.go.dev/github.com/ysmood/uuid)

A distributed unique ID generator inspired by Twitter's Snowflake.

The key difference is it doesn't require a coordinator, it's completely stateless.

Total size is 128 bits, the format looks like this:

```text
32 bits for user defined namespace (4 chars)
56 bits for time in microsecond (from year 2020 to 2262)
16 bits for machine name (65536 machines)
24 bits for cryptographically secure noise (reasonable collision rate)
```

A sample anatomy for id `75756964-000927f291f2d9-3b6a-5253a3` looks like this:

```text
namespace        time          machine    noise
    |             |               |         |
    v             v               v         v
  uuid     2020-04-26T21:26...   3b6a     5253a3   # parsed format

75756964     000927f291f2d9      3b6a     5253a3   # hex format
```

The namespace is usually used to specify the application name. Such as use it to filter a specific app's log in elastic search so that you don't need an extra field to store the filter tag.
