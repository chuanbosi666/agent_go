package tools

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/chuanbosi666/agent_go/pkg/tool"
)

// ExecTools 命令执行工具集
type ExecTools struct {
	// ProjectRoot 项目根目录，命令执行的工作目录
	ProjectRoot string
	// Timeout 命令执行超时时间
	Timeout time.Duration
	// AllowedCommands 允许执行的命令白名单（空表示允许所有）
	AllowedCommands []string
}

// NewExecTools 创建命令执行工具集
func NewExecTools(projectRoot string) *ExecTools {
	return &ExecTools{
		ProjectRoot: projectRoot,
		Timeout:     60 * time.Second, // 默认60秒超时
		AllowedCommands: []string{
			"go", "git", "npm", "node", "python", "pip",
			"make", "cargo", "rustc", "javac", "java",
			"mvn", "gradle", "dotnet", "cmake",
		},
	}
}

// isCommandAllowed 检查命令是否在白名单中
func (et *ExecTools) isCommandAllowed(command string) bool {
	if len(et.AllowedCommands) == 0 {
		return true
	}
	for _, allowed := range et.AllowedCommands {
		if command == allowed {
			return true
		}
	}
	return false
}

// CreateExecCommandTool 创建执行命令工具
func (et *ExecTools) CreateExecCommandTool() tool.FunctionTool {
	return tool.FunctionTool{
		Name:        "exec_command",
		Description: fmt.Sprintf("在项目目录中执行系统命令。允许的命令: %s", strings.Join(et.AllowedCommands, ", ")),
		ParamsJSONSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"command": map[string]any{
					"type":        "string",
					"description": "要执行的命令（如 go, git, npm 等）",
				},
				"args": map[string]any{
					"type": "array",
					"items": map[string]any{
						"type": "string",
					},
					"description": "命令参数列表",
				},
				"workdir": map[string]any{
					"type":        "string",
					"description": "工作目录（相对于项目根目录，可选，默认为项目根目录）",
				},
			},
			"required": []string{"command", "args"},
		},
		OnInvokeTool: func(ctx context.Context, arguments string) (any, error) {
			var params struct {
				Command string   `json:"command"`
				Args    []string `json:"args"`
				Workdir string   `json:"workdir"`
			}
			if err := json.Unmarshal([]byte(arguments), &params); err != nil {
				return nil, fmt.Errorf("参数解析失败: %w", err)
			}

			// 检查命令是否允许
			if !et.isCommandAllowed(params.Command) {
				return nil, fmt.Errorf("命令 %s 不在允许列表中。允许的命令: %s",
					params.Command, strings.Join(et.AllowedCommands, ", "))
			}

			// 确定工作目录
			workdir := et.ProjectRoot
			if params.Workdir != "" {
				workdir = filepath.Join(et.ProjectRoot, params.Workdir)
			}

			// 创建带超时的上下文
			execCtx, cancel := context.WithTimeout(ctx, et.Timeout)
			defer cancel()

			// 执行命令
			cmd := exec.CommandContext(execCtx, params.Command, params.Args...)
			cmd.Dir = workdir

			var stdout, stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr

			err := cmd.Run()

			// 格式化输出
			var result strings.Builder
			result.WriteString(fmt.Sprintf("命令: %s %s\n", params.Command, strings.Join(params.Args, " ")))
			result.WriteString(fmt.Sprintf("工作目录: %s\n", workdir))
			result.WriteString("---\n")

			if stdout.Len() > 0 {
				result.WriteString("标准输出:\n")
				result.WriteString(stdout.String())
				result.WriteString("\n")
			}

			if stderr.Len() > 0 {
				result.WriteString("错误输出:\n")
				result.WriteString(stderr.String())
				result.WriteString("\n")
			}

			if err != nil {
				if execCtx.Err() == context.DeadlineExceeded {
					result.WriteString(fmt.Sprintf("\n命令执行超时（超过 %v）\n", et.Timeout))
				} else {
					result.WriteString(fmt.Sprintf("\n执行失败: %v\n", err))
				}
			} else {
				result.WriteString("\n执行成功\n")
			}

			return result.String(), nil
		},
	}
}

// CreateGoTestTool 创建 Go 测试专用工具
func (et *ExecTools) CreateGoTestTool() tool.FunctionTool {
	return tool.FunctionTool{
		Name:        "go_test",
		Description: "运行 Go 测试。可以指定测试包、运行特定测试用例、设置覆盖率等。",
		ParamsJSONSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"package": map[string]any{
					"type":        "string",
					"description": "要测试的包路径（如 ./... 表示所有包，./pkg/... 表示 pkg 下所有包）",
				},
				"run": map[string]any{
					"type":        "string",
					"description": "要运行的测试名称正则表达式（可选）",
				},
				"verbose": map[string]any{
					"type":        "boolean",
					"description": "是否显示详细输出（默认 true）",
				},
				"coverage": map[string]any{
					"type":        "boolean",
					"description": "是否生成覆盖率报告（默认 false）",
				},
			},
			"required": []string{"package"},
		},
		OnInvokeTool: func(ctx context.Context, arguments string) (any, error) {
			var params struct {
				Package  string `json:"package"`
				Run      string `json:"run"`
				Verbose  *bool  `json:"verbose"`
				Coverage bool   `json:"coverage"`
			}
			if err := json.Unmarshal([]byte(arguments), &params); err != nil {
				return nil, fmt.Errorf("参数解析失败: %w", err)
			}

			// 构建参数
			args := []string{"test"}

			// 默认启用 verbose
			if params.Verbose == nil || *params.Verbose {
				args = append(args, "-v")
			}

			// 覆盖率
			if params.Coverage {
				args = append(args, "-cover")
			}

			// 指定测试用例
			if params.Run != "" {
				args = append(args, "-run", params.Run)
			}

			// 包路径
			args = append(args, params.Package)

			// 创建带超时的上下文
			execCtx, cancel := context.WithTimeout(ctx, et.Timeout)
			defer cancel()

			// 执行测试
			cmd := exec.CommandContext(execCtx, "go", args...)
			cmd.Dir = et.ProjectRoot

			var stdout, stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr

			err := cmd.Run()

			// 格式化输出
			var result strings.Builder
			result.WriteString(fmt.Sprintf("运行命令: go %s\n", strings.Join(args, " ")))
			result.WriteString("---\n")

			if stdout.Len() > 0 {
				result.WriteString(stdout.String())
			}

			if stderr.Len() > 0 {
				result.WriteString(stderr.String())
			}

			if err != nil {
				if execCtx.Err() == context.DeadlineExceeded {
					result.WriteString(fmt.Sprintf("\n测试超时（超过 %v）\n", et.Timeout))
				} else {
					result.WriteString(fmt.Sprintf("\n测试失败: %v\n", err))
				}
			} else {
				result.WriteString("\n所有测试通过\n")
			}

			return result.String(), nil
		},
	}
}

// CreateGoBuildTool 创建 Go 构建工具
func (et *ExecTools) CreateGoBuildTool() tool.FunctionTool {
	return tool.FunctionTool{
		Name:        "go_build",
		Description: "构建 Go 项目。检查代码是否可以正常编译。",
		ParamsJSONSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"package": map[string]any{
					"type":        "string",
					"description": "要构建的包路径（如 ./... 表示所有包）",
				},
				"output": map[string]any{
					"type":        "string",
					"description": "输出文件名（可选）",
				},
			},
			"required": []string{"package"},
		},
		OnInvokeTool: func(ctx context.Context, arguments string) (any, error) {
			var params struct {
				Package string `json:"package"`
				Output  string `json:"output"`
			}
			if err := json.Unmarshal([]byte(arguments), &params); err != nil {
				return nil, fmt.Errorf("参数解析失败: %w", err)
			}

			// 构建参数
			args := []string{"build"}

			if params.Output != "" {
				args = append(args, "-o", params.Output)
			}

			args = append(args, params.Package)

			// 创建带超时的上下文
			execCtx, cancel := context.WithTimeout(ctx, et.Timeout)
			defer cancel()

			// 执行构建
			cmd := exec.CommandContext(execCtx, "go", args...)
			cmd.Dir = et.ProjectRoot

			var stdout, stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr

			err := cmd.Run()

			// 格式化输出
			var result strings.Builder
			result.WriteString(fmt.Sprintf("运行命令: go %s\n", strings.Join(args, " ")))
			result.WriteString("---\n")

			if stdout.Len() > 0 {
				result.WriteString(stdout.String())
			}

			if stderr.Len() > 0 {
				result.WriteString(stderr.String())
			}

			if err != nil {
				if execCtx.Err() == context.DeadlineExceeded {
					result.WriteString(fmt.Sprintf("\n构建超时（超过 %v）\n", et.Timeout))
				} else {
					result.WriteString(fmt.Sprintf("\n构建失败: %v\n", err))
				}
			} else {
				result.WriteString("\n构建成功\n")
			}

			return result.String(), nil
		},
	}
}

// GetAllTools 返回所有执行工具
func (et *ExecTools) GetAllTools() []tool.FunctionTool {
	return []tool.FunctionTool{
		et.CreateExecCommandTool(),
		et.CreateGoTestTool(),
		et.CreateGoBuildTool(),
	}
}
