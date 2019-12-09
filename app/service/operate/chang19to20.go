package operate

import (
	"github.com/gogf/gf/util/gconv"
)

func luhnSum(inVal string) int {
	tmpStr := reverseString(inVal)
	evenSum := 0
	isOdd := false
	oddSum := 0

	for i := 0; i < len(tmpStr); i++ {
		digit := gconv.Int(tmpStr[i : i+1])
		if isOdd {
			oddSum = oddSum + digit
			isOdd = false

		} else {
			//把奇数位数字取出，然后乘以2
			digit = digit * 2
			//把所得数字相加，如果是两位数结果，拆开相加得到A，比如14拆成1+4；
			if digit > 9 {
				digit1 := gconv.Int(gconv.String(digit)[0:1])
				digit2 := gconv.Int(gconv.String(digit)[1:2])
				digit = digit1 + digit2
			}
			evenSum = evenSum + digit
			isOdd = true
		}

	}

	return oddSum + evenSum
}

func LuhnNext(inVal string) int {
	rst := luhnSum(inVal) % 10
	if rst == 0 {
		return 0
	} else {
		return 10 - rst
	}
}

// 反转字符串
func reverseString(s string) string {
	runes := []rune(s)
	for from, to := 0, len(runes)-1; from < to; from, to = from+1, to-1 {
		runes[from], runes[to] = runes[to], runes[from]
	}
	return string(runes)
}
