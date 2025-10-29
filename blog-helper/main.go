package main

import (
	"bufio"
	"fmt"
	"os"

	"blog-helper/commands"
	"blog-helper/utils"
)

func main() {
	// 加载配置文件
	if err := utils.LoadConfig("config.yaml"); err != nil {
		fmt.Printf("加载配置文件失败: %v\n", err)
		return
	}

	fmt.Println("Hugo博客助手 - 输入数字选择操作:")
	fmt.Println("1. 创建新文章")
	fmt.Println("2. 提交并推送更改")
	fmt.Println("按Ctrl+C退出程序")

	scanner := bufio.NewScanner(os.Stdin)

	// 主循环
	for {
		fmt.Print("\n请输入操作编号(1/2): ")
		if !scanner.Scan() {
			fmt.Printf("输入错误: %v\n", scanner.Err())
			continue
		}

		choice := scanner.Text()
		switch choice {
		case "1":
			if err := commands.CreatePost(scanner); err != nil {
				fmt.Printf("操作失败: %v\n", err)
			}
		case "2":
			if err := commands.PushChanges(scanner); err != nil {
				fmt.Printf("操作失败: %v\n", err)
			}
		default:
			fmt.Println("无效的选择，请输入1或2")
		}
	}
}