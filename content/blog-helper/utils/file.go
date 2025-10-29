package utils

import (
	"os"
	"os/exec"
	"path/filepath"
)

// CreatePostFile 创建文章文件并写入内容
func CreatePostFile(title string, content string) (string, error) {
	// 构建文件路径
	dirPath := filepath.Join("content", "post", title) // 假设Hugo文章放在content/post目录
	filePath := filepath.Join(dirPath, "index.md")

	// 检查目录是否已存在
	if _, err := os.Stat(dirPath); err == nil {
		return "", os.ErrExist
	}

	// 创建多级目录
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return "", err
	}

	// 创建并写入文件
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		// 写入失败时清理已创建的目录
		os.RemoveAll(dirPath)
		return "", err
	}

	return filePath, nil
}

// OpenFileInEditor 用默认编辑器打开文件
func OpenFileInEditor(filePath string) error {
	var cmd *exec.Cmd

	// 根据操作系统选择默认编辑器命令
	switch os := os.Getenv("OS"); os {
	case "Windows_NT":
		cmd = exec.Command("cmd", "/c", "start", filePath)
	default: // Linux, macOS等
		cmd = exec.Command("xdg-open", filePath) // Linux
		// 如果xdg-open不存在，尝试open命令(macOS)
		if err := cmd.Run(); err != nil {
			cmd = exec.Command("open", filePath)
		}
	}

	return cmd.Run()
}