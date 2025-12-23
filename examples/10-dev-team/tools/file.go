// Package tools 提供开发团队 Agent 使用的工具集
package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/chuanbosi666/agent_go/pkg/tool"
)

// FileTools 文件操作工具集，包含安全路径校验
type FileTools struct {
	// ProjectRoot 项目根目录，所有操作限制在此目录内
	ProjectRoot string
}

// NewFileTools 创建文件工具集
func NewFileTools(projectRoot string) *FileTools {
	return &FileTools{ProjectRoot: projectRoot}
}

// validatePath 校验路径是否在项目目录内
func (ft *FileTools) validatePath(targetPath string) (string, error) {
	// 如果是相对路径，转为绝对路径
	if !filepath.IsAbs(targetPath) {
		targetPath = filepath.Join(ft.ProjectRoot, targetPath)
	}

	absProject, err := filepath.Abs(ft.ProjectRoot)
	if err != nil {
		return "", fmt.Errorf("获取项目根目录绝对路径失败: %w", err)
	}

	absTarget, err := filepath.Abs(targetPath)
	if err != nil {
		return "", fmt.Errorf("获取目标路径绝对路径失败: %w", err)
	}

	// 检查是否在项目目录内
	if !strings.HasPrefix(absTarget, absProject) {
		return "", fmt.Errorf("拒绝访问: 路径 %s 超出项目根目录 %s", targetPath, ft.ProjectRoot)
	}

	return absTarget, nil
}

// CreateReadFileTool 创建读取文件工具
func (ft *FileTools) CreateReadFileTool() tool.FunctionTool {
	return tool.FunctionTool{
		Name:        "read_file",
		Description: "读取指定文件的内容。路径可以是相对于项目根目录的相对路径。",
		ParamsJSONSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"path": map[string]any{
					"type":        "string",
					"description": "要读取的文件路径（相对于项目根目录）",
				},
			},
			"required": []string{"path"},
		},
		OnInvokeTool: func(ctx context.Context, arguments string) (any, error) {
			var params struct {
				Path string `json:"path"`
			}
			if err := json.Unmarshal([]byte(arguments), &params); err != nil {
				return nil, fmt.Errorf("参数解析失败: %w", err)
			}

			// 路径安全校验
			absPath, err := ft.validatePath(params.Path)
			if err != nil {
				return nil, err
			}

			// 读取文件
			content, err := os.ReadFile(absPath)
			if err != nil {
				if os.IsNotExist(err) {
					return fmt.Sprintf("文件不存在: %s", params.Path), nil
				}
				return nil, fmt.Errorf("读取文件失败: %w", err)
			}

			return string(content), nil
		},
	}
}

// CreateWriteFileTool 创建写入文件工具
func (ft *FileTools) CreateWriteFileTool() tool.FunctionTool {
	return tool.FunctionTool{
		Name:        "write_file",
		Description: "写入内容到指定文件。如果文件不存在则创建，如果目录不存在则自动创建目录。",
		ParamsJSONSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"path": map[string]any{
					"type":        "string",
					"description": "要写入的文件路径（相对于项目根目录）",
				},
				"content": map[string]any{
					"type":        "string",
					"description": "要写入的文件内容",
				},
			},
			"required": []string{"path", "content"},
		},
		OnInvokeTool: func(ctx context.Context, arguments string) (any, error) {
			var params struct {
				Path    string `json:"path"`
				Content string `json:"content"`
			}
			if err := json.Unmarshal([]byte(arguments), &params); err != nil {
				return nil, fmt.Errorf("参数解析失败: %w", err)
			}

			// 路径安全校验
			absPath, err := ft.validatePath(params.Path)
			if err != nil {
				return nil, err
			}

			// 创建目录（如果不存在）
			dir := filepath.Dir(absPath)
			if err := os.MkdirAll(dir, 0755); err != nil {
				return nil, fmt.Errorf("创建目录失败: %w", err)
			}

			// 写入文件
			if err := os.WriteFile(absPath, []byte(params.Content), 0644); err != nil {
				return nil, fmt.Errorf("写入文件失败: %w", err)
			}

			return fmt.Sprintf("成功写入文件: %s (%d 字节)", params.Path, len(params.Content)), nil
		},
	}
}

// CreateListDirTool 创建列出目录内容工具
func (ft *FileTools) CreateListDirTool() tool.FunctionTool {
	return tool.FunctionTool{
		Name:        "list_dir",
		Description: "列出指定目录的内容，包括文件和子目录。",
		ParamsJSONSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"path": map[string]any{
					"type":        "string",
					"description": "要列出的目录路径（相对于项目根目录），空字符串表示项目根目录",
				},
			},
			"required": []string{"path"},
		},
		OnInvokeTool: func(ctx context.Context, arguments string) (any, error) {
			var params struct {
				Path string `json:"path"`
			}
			if err := json.Unmarshal([]byte(arguments), &params); err != nil {
				return nil, fmt.Errorf("参数解析失败: %w", err)
			}

			// 如果路径为空，使用项目根目录
			targetPath := params.Path
			if targetPath == "" {
				targetPath = ft.ProjectRoot
			}

			// 路径安全校验
			absPath, err := ft.validatePath(targetPath)
			if err != nil {
				return nil, err
			}

			// 读取目录
			entries, err := os.ReadDir(absPath)
			if err != nil {
				if os.IsNotExist(err) {
					return fmt.Sprintf("目录不存在: %s", params.Path), nil
				}
				return nil, fmt.Errorf("读取目录失败: %w", err)
			}

			// 格式化输出
			var result strings.Builder
			result.WriteString(fmt.Sprintf("目录 %s 内容:\n", params.Path))
			for _, entry := range entries {
				if entry.IsDir() {
					result.WriteString(fmt.Sprintf("  [目录] %s/\n", entry.Name()))
				} else {
					info, _ := entry.Info()
					size := int64(0)
					if info != nil {
						size = info.Size()
					}
					result.WriteString(fmt.Sprintf("  [文件] %s (%d 字节)\n", entry.Name(), size))
				}
			}

			if len(entries) == 0 {
				result.WriteString("  (空目录)\n")
			}

			return result.String(), nil
		},
	}
}

// CreateAppendFileTool 创建追加写入文件工具
func (ft *FileTools) CreateAppendFileTool() tool.FunctionTool {
	return tool.FunctionTool{
		Name:        "append_file",
		Description: "在指定文件末尾追加内容。如果文件不存在则创建。",
		ParamsJSONSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"path": map[string]any{
					"type":        "string",
					"description": "要追加的文件路径（相对于项目根目录）",
				},
				"content": map[string]any{
					"type":        "string",
					"description": "要追加的内容",
				},
			},
			"required": []string{"path", "content"},
		},
		OnInvokeTool: func(ctx context.Context, arguments string) (any, error) {
			var params struct {
				Path    string `json:"path"`
				Content string `json:"content"`
			}
			if err := json.Unmarshal([]byte(arguments), &params); err != nil {
				return nil, fmt.Errorf("参数解析失败: %w", err)
			}

			// 路径安全校验
			absPath, err := ft.validatePath(params.Path)
			if err != nil {
				return nil, err
			}

			// 创建目录（如果不存在）
			dir := filepath.Dir(absPath)
			if err := os.MkdirAll(dir, 0755); err != nil {
				return nil, fmt.Errorf("创建目录失败: %w", err)
			}

			// 打开文件（追加模式）
			file, err := os.OpenFile(absPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				return nil, fmt.Errorf("打开文件失败: %w", err)
			}
			defer file.Close()

			// 写入内容
			if _, err := file.WriteString(params.Content); err != nil {
				return nil, fmt.Errorf("写入内容失败: %w", err)
			}

			return fmt.Sprintf("成功追加内容到文件: %s (%d 字节)", params.Path, len(params.Content)), nil
		},
	}
}

// GetAllTools 返回所有文件操作工具
func (ft *FileTools) GetAllTools() []tool.FunctionTool {
	return []tool.FunctionTool{
		ft.CreateReadFileTool(),
		ft.CreateWriteFileTool(),
		ft.CreateListDirTool(),
		ft.CreateAppendFileTool(),
	}
}
