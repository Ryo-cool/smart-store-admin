package main

import (
	"fmt"
	"os"
)

func DoSomething() {
	f, _ := os.Open("non_existent_file.txt") // エラーを無視
	defer f.Close()                          // ファイルが開けなかった場合でも実行される

	data := make([]byte, 100)
	n, _ := f.Read(data) // エラーを無視
	fmt.Println("Read data:", data[:n])

	// 名前が適切でない
	myVar := 123
	fmt.Println("Variable value:", myVar)

	// 無駄なループ
	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			if i == j {
				fmt.Println("Matching values:", i)
			}
		}
	}
}

func main() {
	DoSomething() // メイン関数に直接書くのはベストプラクティスではない
}
