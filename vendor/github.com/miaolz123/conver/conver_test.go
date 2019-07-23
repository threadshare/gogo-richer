package conver

import (
	"log"
	"testing"
)

func TestString(t *testing.T) {
	var i = "\n \t33 4 6999   66666    .677777\n"
	log.Println(i)
	log.Println(Bool(i))
	log.Println(BoolMust(i, true))
	log.Println(Bytes(i))
	log.Println(BytesMust(i))
	log.Println(Float32(i))
	log.Println(Float32Must(i, 6.666))
	log.Println(Float64(i))
	log.Println(Float64Must(i, 6.666))
	log.Println(Int(i))
	log.Println(IntMust(i, 6666))
	log.Println(Int32(i))
	log.Println(Int32Must(i, 6666))
	log.Println(Int64(i))
	log.Println(Int64Must(i, 6666))
	log.Println(String(i))
	log.Println(StringMust(i, "HHHH"))
}
