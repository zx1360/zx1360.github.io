package commands

import (
	"bufio"
	"fmt"
	"os"

	"blog-helper/utils"
)

// CreatePost 处理创建文章命令
func CreatePost(scanner *bufio.Scanner) error {
	fmt.Print("请输入文章标题: ")
	if !scanner.Scan() {
		return scanner.Err()
	}
	title := scanner.Text()
	if title == "" {
		return fmt.Errorf("文章标题不能为空")
	}

	// 渲染模板内容
	content, err := utils.RenderPostTemplate(title)
	if err != nil {
		return fmt.Errorf("渲染模板失败: %v", err)
	}

	// 创建文章文件
	filePath, err := utils.CreatePostFile(title, content)
	if err != nil {
		if os.IsExist(err) {
			return fmt.Errorf("错误: 同名文章目录已存在")
		}
		return fmt.Errorf("创建文章文件失败: %v", err)
	}

	fmt.Printf("文章已创建: %s\n", filePath)

	// 打开文件
	if err := utils.OpenFileInEditor(filePath); err != nil {
		fmt.Printf("提示: 无法自动打开编辑器，可手动打开文件: %s\n", filePath)
	}

	return nil
}