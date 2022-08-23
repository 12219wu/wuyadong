package util

// 0的ASCII码值
const zero = byte(48)

// 9的ASCII码值
const nine = byte(57)

// ToBCD 数字字符（ASCII）转BCD码
func ToBCD(src []byte) []byte {
	length := len(src)
	if length <= 0 {
		return nil
	}

	bts := make([]byte, (length+1)>>1)
	for i := 0; i < length; i += 2 {
		check(src[i])
		bts[i>>1] = (src[i] - zero) << 4
		if i+1 < length {
			check(src[i+1])
			bts[i>>1] |= src[i+1] - zero
		}
	}

	return bts
}

// ParseBCD 解析BCD码为数字字符（ASCII）
func ParseBCD(src []byte) []byte {
	length := len(src)
	if length <= 0 {
		return nil
	}

	bts := make([]byte, length<<1)
	for idx, bt := range src {
		bts[idx<<1] = (bt >> 4) + zero
		bts[(idx<<1)+1] = (bt & 0x0F) + zero
	}

	return bts
}

// 检查字符是否是数字字符
func check(bt byte) {
	if bt < zero || bt > nine {
		panic("BCD must digital character.")
	}
}
