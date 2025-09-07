package gromb

// LRC lrc sum.
type LRC struct {
	sum uint8
}

// Reset rest lrc sum.
func (sf *LRC) Reset() *LRC {
	sf.sum = 0
	return sf
}

// Push data in sum.
func (sf *LRC) Push(data ...uint8) *LRC {
	for _, b := range data {
		sf.sum += b
	}
	return sf
}

// Value got lrc value.
func (sf *LRC) Value() uint8 {
	return uint8(-int8(sf.sum))
}

func LRCCalcul(data []uint8) uint8 {
	sum := uint8(0)
	for _, b := range data {
		sum += b
	}
	return uint8(^int8(sum) + 1)
}
