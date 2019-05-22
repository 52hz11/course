package models

import (
	"math/rand"
	"time"

	"github.com/bitly/go-simplejson"
)

//本来按道理一个函数就够的，但是abort和body的接收一个是[]byte一个是string，因为懒得在外面转了所以就分成两个算了
func ErrJson(msg string) string {
	json := simplejson.New()
	json.Set("status", "failed")
	json.Set("msg", msg)
	str, _ := json.Encode()
	return string(str)
}

func SuccessJson() []byte {
	json := simplejson.New()
	json.Set("status", "success")
	str, _ := json.Encode()
	return str
}

//生成key相关
func init() {
	rand.Seed(time.Now().UnixNano())
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func GenerateKey() string {
	b := make([]byte, 20)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
