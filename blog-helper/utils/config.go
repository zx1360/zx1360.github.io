package utils

import (
	"bytes"
	"os"
	"text/template"
	"time"

	"gopkg.in/yaml.v3"
)

// Config 配置文件结构
type Config struct {
	SourceBranch         string `yaml:"source_branch"`
	DefaultCommitMessage string `yaml:"default_commit_message"`
	PostTemplate         string `yaml:"post_template"`
}

// 全局配置实例
var AppConfig Config

// LoadConfig 从文件加载配置
func LoadConfig(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, &AppConfig)
}

// RenderPostTemplate 渲染文章模板
func RenderPostTemplate(title string) (string, error) {
	tpl, err := template.New("post").Parse(AppConfig.PostTemplate)
	if err != nil {
		return "", err
	}

	// 准备模板数据
	data := struct {
		Date  string
		Title string
	}{
		Date:  time.Now().Format("2006-01-02T15:04:05+08:00"),
		Title: title,
	}

	// 执行模板渲染
	var buf bytes.Buffer
	err = tpl.Execute(&buf, data)
	if err != nil {
		return "", err
	}

	// 将缓冲区内容转换为字符串返回
	return buf.String(), nil
}