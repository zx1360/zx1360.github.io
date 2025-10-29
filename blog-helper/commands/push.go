package commands

import (
	"bufio"
	"fmt"
	"os/exec"

	"blog-helper/utils"
)

// PushChanges 处理推送代码命令
func PushChanges(scanner *bufio.Scanner) error {
	fmt.Print("请输入提交信息(直接回车使用默认): ")
	if !scanner.Scan() {
		return scanner.Err()
	}
	message := scanner.Text()

	// 使用默认提交信息
	if message == "" {
		message = utils.AppConfig.DefaultCommitMessage
	}

	// 执行git命令
	commands := [][]string{
		{"git", "add", "."},
		{"git", "commit", "-m", message},
		{"git", "push", "origin", utils.AppConfig.SourceBranch},
	}

	for _, cmdArgs := range commands {
		cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
		cmd.Stdout = fmt.Println // 输出命令执行结果
		cmd.Stderr = fmt.Println // 输出错误信息

		if err := cmd.Run(); err != nil {
			return fmt.Errorf("执行命令 %s 失败: %v", cmdArgs[0], err)
		}
	}

	fmt.Println("代码已成功推送到远程仓库")
	return nil
}