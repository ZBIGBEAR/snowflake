package snowflake

import "errors"

const (
	NB4095 = 4095
)

var CountErr = errors.New("one millisecond can generat max id is 4096")
var SeqStack = errors.New("sequence is bigger than 4095")
