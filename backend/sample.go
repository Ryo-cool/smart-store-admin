package main

import (
	"fmt"
	"os"
)

func DoSomething() {
	f, err := os.Open("non_existent_file.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return // ファイルが開けない場合は処理を中断
	}
	defer f.Close()

	data := make([]byte, 100)
	n, err := f.Read(data)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return // ファイルが読み込めない場合は処理を中断
	}
	fmt.Println("Read data:", data[:n])

	// 変数名をより適切に変更
	sampleVariable := 123
	fmt.Println("Variable value:", sampleVariable)

	// 無駄なループを削除し、同等の処理を簡潔に記述
	for i := 0; i < 10; i++ {
		fmt.Println("Matching values:", i)
	}
}
