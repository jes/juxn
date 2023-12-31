package juxn

type Stack struct {
	data []byte
	ptr  int
}

const STACKSZ = 256

func NewStack() *Stack {
	return &Stack{
		data: make([]byte, STACKSZ),
		ptr:  0,
	}
}

func (s *Stack) Push(val uint16, short bool) bool {
	if short {
		return s.PushShort(val)
	} else {
		return s.PushByte(byte(val))
	}
}

func (s *Stack) PushShort(val uint16) bool {
	if s.ptr > STACKSZ-2 {
		return false
	}
	s.data[s.ptr] = byte(val >> 8)
	s.data[s.ptr+1] = byte(val & 0xff)
	s.ptr += 2
	return true
}

func (s *Stack) PushByte(val byte) bool {
	if s.ptr > STACKSZ-1 {
		return false
	}
	s.data[s.ptr] = val
	s.ptr += 1
	return true
}

func (s *Stack) Pop(short bool) (uint16, bool) {
	if short {
		return s.PopShort()
	} else {
		val, ok := s.PopByte()
		return uint16(val), ok
	}
}

func (s *Stack) PopShort() (uint16, bool) {
	val, ok := s.PeekShort()
	if ok {
		s.ptr -= 2
	}
	return val, ok
}

func (s *Stack) PopByte() (byte, bool) {
	val, ok := s.PeekByte()
	if ok {
		s.ptr -= 1
	}
	return val, ok
}

func (s *Stack) Peek(short bool) (uint16, bool) {
	return s.PeekOffset(short, 0)
}

func (s *Stack) PeekShort() (uint16, bool) {
	return s.PeekShortOffset(0)
}

func (s *Stack) PeekByte() (byte, bool) {
	return s.PeekByteOffset(0)
}

func (s *Stack) PeekOffset(short bool, offset int) (uint16, bool) {
	if short {
		return s.PeekShortOffset(offset)
	} else {
		val, ok := s.PeekByteOffset(offset)
		return uint16(val), ok
	}
}

func (s *Stack) PeekShortOffset(offset int) (uint16, bool) {
	if s.ptr+offset < 2 {
		return 0, false
	}
	return uint16(s.data[s.ptr+offset-1]) + 256*uint16(s.data[s.ptr+offset-2]), true
}

func (s *Stack) PeekByteOffset(offset int) (byte, bool) {
	if s.ptr+offset < 1 {
		return 0, false
	}
	return s.data[s.ptr+offset-1], true
}

func (s *Stack) Size() int {
	return STACKSZ
}

func (s *Stack) Used() int {
	return s.ptr
}
