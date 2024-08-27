package main

import (
	"fmt"
	"os"
)

func main() {
	// 从环境变量中读取 GitHub Action 的输入
	input := os.Getenv("INPUT_MY_INPUT")
	if input == "" {
		fmt.Println("No input provided")
		return
	}

	// 打印输出以便在 GitHub Action 日志中查看
	fmt.Printf("Hello, %s!\n", input)
}
