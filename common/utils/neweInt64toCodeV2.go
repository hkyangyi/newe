package utils

import (
	"crypto/rand"
	"encoding/binary"
	"errors"
	"strings"
)

/*
用法 / 说明

这个文件提供一个“把非负 int64 编成短字符串，再从字符串还原 int64”的实现，适用于生成可读的短 token。

核心 API

- Int64ToDictString(n int64) (string, error)
	- 将 n 按 DictString 作为“进制字典”编码为变长字符串。
- Int64ToDictStringFixedWidth(n int64, width int) (string, error)
	- 将 n 编码为固定长度 width 的字符串；当核心编码长度不足时会进行“补齐”。
- DictStringToInt64(s string) (int64, error)
	- 将字符串还原为 int64。

编码规则

- DictString 是主字典（类似 base-N 的字符表）。编码得到的“核心串”只会包含 DictString 中的字符。
- 仅支持 n >= 0；负数或非法输入会返回 ErrTokenInvalid。

补齐规则（随机字符 + 随机位置）

- FixedWidth 模式下，如果核心串长度 < width：
	- 会从 DictBq 中挑选“不会出现在 DictString 里的字符”作为补齐字符池（避免与核心串混淆）。
	- 随机选取 need = width-len(core) 个补齐字符。
	- 将这些补齐字符逐个插入到输出串的“随机位置”（可能在头部、中间或尾部），最终长度恰好为 width。
- 随机数使用 crypto/rand（加密安全随机）。

解码规则（与随机位置补齐兼容）

- DictStringToInt64 在还原时，会先“过滤掉所有不在 DictString 中的字符”。
	- 由于补齐字符保证不在 DictString 中，因此无论补齐字符插在什么位置，都不会影响还原结果。

示例

- 运行 main() 内的示例：会打印 Encoded / Decoded。
	- Encoded 字符串长度应为 width（例如 10）。
	- Decoded 应等于原始 n。
*/

const (
	DictString = "q7ZbL1vE0yDk3uNwF8aO5sJcT6rGznpQ4hIMXoUjtRBeAlKWdSg"
	DictBq     = "PV2xgH9YimCf"
)

var ErrTokenInvalid = errors.New("token invalid")

// func main() {
// 	// Example usage:
// 	n := int64(50)
// 	s, err := Int64ToDictStringFixedWidth(n, 10)
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println("Encoded:", s)

// 	n2, err := DictStringToInt64(s)
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println("Decoded:", n2)
// }

func Int64ToDictString(n int64) (string, error) {
	return int64ToAlphabet(n, DictString)
}

func Int64ToDictStringFixedWidth(n int64, width int) (string, error) {
	return int64ToAlphabetFixedWidth(n, width, DictString, DictBq)
}

func DictStringToInt64(s string) (int64, error) {
	// If s is produced by FixedWidth, it may contain padding chars from DictBq.
	// Padding chars are chosen to be NOT present in DictString, so we can safely
	// drop any non-DictString chars regardless of where they were inserted.
	buf := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		b := s[i]
		if strings.IndexByte(DictString, b) >= 0 {
			buf = append(buf, b)
		}
	}
	if len(buf) == 0 {
		return 0, ErrTokenInvalid
	}
	return alphabetToInt64(string(buf), DictString)
}

func int64ToAlphabet(n int64, alphabet string) (string, error) {
	if n < 0 {
		return "", ErrTokenInvalid
	}
	if len(alphabet) < 2 {
		return "", ErrTokenInvalid
	}
	if n == 0 {
		return string(alphabet[0]), nil
	}

	base := int64(len(alphabet))
	buf := make([]byte, 0, 16)
	for n > 0 {
		r := n % base
		buf = append(buf, alphabet[int(r)])
		n /= base
	}
	for i, j := 0, len(buf)-1; i < j; i, j = i+1, j-1 {
		buf[i], buf[j] = buf[j], buf[i]
	}
	return string(buf), nil
}

func int64ToAlphabetFixedWidth(n int64, width int, alphabet, padAlphabet string) (string, error) {
	if width <= 0 {
		return "", ErrTokenInvalid
	}
	core, err := int64ToAlphabet(n, alphabet)
	if err != nil {
		return "", err
	}
	if len(core) > width {
		return "", ErrTokenInvalid
	}
	if len(core) == width {
		return core, nil
	}

	// To make padding unambiguous, we only use pad chars that are NOT present in the main alphabet.
	padPool := make([]byte, 0, len(padAlphabet))
	for i := 0; i < len(padAlphabet); i++ {
		ch := padAlphabet[i]
		if !strings.ContainsRune(alphabet, rune(ch)) {
			padPool = append(padPool, ch)
		}
	}
	if len(padPool) == 0 {
		padPool = []byte{alphabet[0]}
	}

	need := width - len(core)
	pad, err := randomFromPool(padPool, need)
	if err != nil {
		return "", err
	}
	out, err := insertRandomly([]byte(core), pad)
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func insertRandomly(core []byte, pad []byte) ([]byte, error) {
	out := make([]byte, len(core))
	copy(out, core)
	for _, p := range pad {
		idx, err := randomIntn(len(out) + 1)
		if err != nil {
			return nil, err
		}
		out = append(out, 0)
		copy(out[idx+1:], out[idx:])
		out[idx] = p
	}
	return out, nil
}

func randomIntn(n int) (int, error) {
	if n <= 0 {
		return 0, ErrTokenInvalid
	}
	var b [4]byte
	un := uint32(n)
	lim := uint32(0xFFFFFFFF - (0xFFFFFFFF % un))
	for {
		if _, err := rand.Read(b[:]); err != nil {
			return 0, err
		}
		v := binary.LittleEndian.Uint32(b[:])
		if v < lim {
			return int(v % un), nil
		}
	}
}

func randomFromPool(pool []byte, length int) ([]byte, error) {
	if length < 0 {
		return nil, ErrTokenInvalid
	}
	if length == 0 {
		return []byte{}, nil
	}
	if len(pool) == 0 {
		return nil, ErrTokenInvalid
	}
	out := make([]byte, length)
	buf := make([]byte, 32)
	i := 0
	max := byte(len(pool))
	lim := byte(256 - (256 % int(max)))
	for i < length {
		if _, err := rand.Read(buf); err != nil {
			return nil, err
		}
		for _, b := range buf {
			if b >= lim {
				continue
			}
			out[i] = pool[int(b%max)]
			i++
			if i == length {
				break
			}
		}
	}
	return out, nil
}

func alphabetToInt64(s string, alphabet string) (int64, error) {
	if s == "" {
		return 0, ErrTokenInvalid
	}
	idx := make(map[rune]int, len(alphabet))
	for i, r := range alphabet {
		idx[r] = i
	}
	base := int64(len(alphabet))
	var n int64
	for _, r := range s {
		v, ok := idx[r]
		if !ok {
			return 0, ErrTokenInvalid
		}
		// overflow check: n*base + v <= MaxInt64
		if n > (int64(^uint64(0)>>1)-int64(v))/base {
			return 0, ErrTokenInvalid
		}
		n = n*base + int64(v)
	}
	return n, nil
}
