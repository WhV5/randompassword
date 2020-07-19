/**
* @Author : henry
* @Data: 2020-07-19 15:03
* @Note:
**/

package main

import (
	"fmt"
	"math/rand"
	"time"
)

// ASCII   48~57 数字  65~90 大写字母 97~122 小写字母
var StrBytes []byte
var Num int

func main() {
	Num = 10
	for i := 0; i < 10; i++ {
		pwd := setPwd()
		fmt.Println(pwd)
	}
}

func setPwd() string {
	intChan := func(done <-chan interface{}) <-chan int {
		intStream := make(chan int)
		go func() {
			defer close(intStream)
			for {
				select {
				case <-done:
					return
				case intStream <- rand.Intn(26):
				}
			}
		}()
		return intStream
	}

	randomPwd := func(done <-chan interface{}, intStream <-chan int) <-chan int {
		numStream := make(chan int, 1)
		capStream := make(chan int, 1)
		lowStream := make(chan int, 1)
		pwdStream := make(chan int)
		go func() {
			defer close(numStream)
			defer close(capStream)
			defer close(lowStream)
			defer close(pwdStream)
			for i := range intStream {
				if i > 10 {
					select {
					case <-done:
						return
					case capStream <- i:
						pwdStream <- <-capStream + 65
					case lowStream <- i:
						pwdStream <- <-lowStream + 97
					}
				} else {
					select {
					case <-done:
						return
					case numStream <- i:
						pwdStream <- <-numStream + 48
					case capStream <- i:
						pwdStream <- <-capStream + 65
					case lowStream <- i:
						pwdStream <- <-lowStream + 97
					}
				}
			}
		}()
		return pwdStream
	}

	rand.Seed(time.Now().UnixNano())
	done := make(chan interface{})
	n := 10
	bytes := make([]byte, n)
	for i := 0; i < n; i++ {
		if i == n {
			close(done)
		}
		num := intChan(done)
		pwd := randomPwd(done, num)
		bytes[i] = byte(<-pwd)
	}

	return string(bytes)
}
