package utils

import (
	"bufio"
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math"
	"math/rand"
	"strings"
	"time"
)

var (
	Newekey   = "2GrE5hsXx0URawilALpj6zNnQ47odFqm"
	Newescret = "1389bcefgktuvyBCDHIJKMOPSTVWYZ"
)

//生成密钥
func RandKey(length int) (string, string) {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		index := r.Intn(len(bytes))
		result = append(result, bytes[index])
		bytes = append(bytes[:index], bytes[index+1:]...)
	}
	return string(result), string(bytes)
}

//加密
func NeweEncode(num int64) string {
	str := encode(Newekey, int(num))
	fmt.Println("str--------------", str)
	strs := setin(str, Newescret)
	return strs
}

//解密
func NeweDecode(str string) int64 {
	num := decode(Newekey, str)
	return int64(num)
}

//加密
func encode(key string, num int) string {
	keybytes := []byte(key)
	new_num_str := ""
	var remainder int
	var remainder_string string
	for num != 0 {
		remainder = num % 32
		remainder_string = string(keybytes[remainder])
		new_num_str = remainder_string + new_num_str
		num = num / 32
	}
	return new_num_str
}

//解密
func decode(key, numstr string) int {
	keybytes := []byte(key)
	var keydict = make(map[byte]float64)
	for i := 0; i < len(keybytes); i++ {
		keydict[keybytes[i]] = float64(i)
	}
	var new_num float64
	nNum := len(strings.Split(numstr, "")) - 1
	numbytes := []byte(numstr)
	for _, v := range numbytes {
		if _, ok := keydict[v]; !ok {
			nNum = nNum - 1
		}
	}
	for _, v := range numbytes {
		if _, ok := keydict[v]; ok {
			new_num = new_num + keydict[v]*math.Pow(32, float64(nNum))
			nNum = nNum - 1
		}
	}
	return int(new_num)
}

//插入字符串
func setin(val, dict string) string {
	fmt.Println(val)
	length := 10 - len(val)
	// val 如果小于0或等于0
	if length <= 0 {
		return val
	}
	str := randstr(dict, length)
	if len(val) == 0 {
		return str
	}
	valbytes := []byte(val)
	strbytes := []byte(str)
	newbytes := make([]byte, 10)
	copy(newbytes[0:], valbytes)
	for i, v := range strbytes {
		r := rand.New(rand.NewSource(time.Now().UnixNano() + int64(i)))
		index := r.Intn(len(val) + i)
		fmt.Println(index)
		copy(newbytes[index+1:], newbytes[index:10])
		newbytes[index] = v
	}
	return string(newbytes)
}

//随机字符串
func randstr(dict string, length int) string {
	bytes := []byte(dict)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		index := r.Intn(len(bytes))
		result = append(result, bytes[index])
		bytes = append(bytes[:index], bytes[index+1:]...)
	}
	return string(result)
}

func EncodeMD5(value string) string {
	m := md5.New()
	m.Write([]byte(value))
	return hex.EncodeToString(m.Sum(nil))
}

func MD5Sum(txt string) (sum string) {
	h := md5.New()
	buf := bufio.NewWriterSize(h, 128)
	buf.WriteString(txt)
	buf.Flush()
	sign := make([]byte, hex.EncodedLen(h.Size()))
	hex.Encode(sign, h.Sum(nil))
	sum = string(bytes.ToUpper(sign))
	return
}

func getDerivedKey(password string, salt []byte, count int) ([]byte, []byte) {
	key := md5.Sum([]byte(password + string(salt)))
	for i := 0; i < count-1; i++ {
		key = md5.Sum(key[:])
	}
	return key[:8], key[8:]
}

//PBE加密
func Encrypt(password string, obtenationIterations int, plainText string, salt []byte) (string, error) {
	padNum := byte(8 - len(plainText)%8)
	for i := byte(0); i < padNum; i++ {
		plainText += string(padNum)
	}

	dk, iv := getDerivedKey(password, salt, obtenationIterations)

	block, err := des.NewCipher(dk)

	if err != nil {
		return "", err
	}

	encrypter := cipher.NewCBCEncrypter(block, iv)
	encrypted := make([]byte, len(plainText))

	encrypter.CryptBlocks(encrypted, []byte(plainText))
	//fmt.Println(string(encrypted))
	encodedStr := hex.EncodeToString(encrypted)
	//fmt.Println(encodedStr)
	//return base64.StdEncoding.EncodeToString(encrypted), nil
	return encodedStr, nil
}

//PBE解密
func Decrypt(password string, obtenationIterations int, cipherText string, salt []byte) (string, error) {
	//msgBytes, err := base64.StdEncoding.DecodeString(cipherText)
	msgBytes, err := hex.DecodeString(cipherText)
	if err != nil {
		return "", err
	}

	dk, iv := getDerivedKey(password, salt, obtenationIterations)

	block, err := des.NewCipher(dk)

	if err != nil {
		return "", err
	}

	decrypter := cipher.NewCBCDecrypter(block, iv)
	decrypted := make([]byte, len(msgBytes))
	decrypter.CryptBlocks(decrypted, msgBytes)

	decryptedString := strings.TrimRight(string(decrypted), "\x01\x02\x03\x04\x05\x06\x07\x08")

	return decryptedString, nil
}
