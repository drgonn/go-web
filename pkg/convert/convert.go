package convert

import "strconv"

type StrTo string

func (s StrTo) String() string {
	return string(s)
}
// 数字字符串转int，非数字会产生0和error
func (s StrTo) Int() (int, error) {
	v, err := strconv.Atoi(s.String())
	return v, err
}
// 数字字符串转int，非数字会产生0
func (s StrTo) MustInt() int {
	v, _ := s.Int()
	return v
}

// 数字字符串转uint32，非数字会产生0和error
func (s StrTo) UInt32() (uint32, error) {
	v, err := strconv.Atoi(s.String())
	return uint32(v), err
}

func (s StrTo) MustUInt32() uint32 {
	v, _ := s.UInt32()
	return v
}