# conver [![GoDoc](https://godoc.org/github.com/miaolz123/conver?status.svg)](https://godoc.org/github.com/miaolz123/conver) [![Build Status](https://travis-ci.org/miaolz123/conver.svg?branch=master)](https://travis-ci.org/miaolz123/conver)

### Golang Type Reflect And Conver

```go
package main

import (
	"fmt"

	"github.com/miaolz123/conver"
)

func main() {
	val := "1 234.567\t\n"
	fmt.Println(conver.Float64Must(val)) // output: 1234.567
	fmt.Println(conver.IntMust(val))     // output: 1235
	val = "ok"
	fmt.Println(conver.BoolMust(val))           // output: true
	fmt.Println(conver.Float64Must(val, 666.6)) // output: 666.6
	fmt.Println(conver.IntMust(val, 6666))      // output: 666
}
```