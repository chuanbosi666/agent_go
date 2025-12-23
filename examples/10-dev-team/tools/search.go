package tools

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/chuanbosi666/agent_go/pkg/tool"
)

// SearchTools 搜索工具集
type SearchTools struct {
	// ProjectRoot 项目根目录
	ProjectRoot string
	// MaxResults 最大结果数
	MaxResults int
	// IgnorePatterns 忽略的目录模式
	IgnorePatterns []string
}

// NewSearchTools 创建搜索工具集
func NewSearchTools(projectRoot string) *SearchTools {
	return &SearchTools{
		ProjectRoot: projectRoot,
		MaxResults:  100,
		IgnorePatterns: []string{
			".git", "node_modules", "vendor", "__pycache__",
			".idea", ".vscode", "dist", "build", "target",
		},
	}
}

// shouldIgnore 检查路径是否应该被忽略
func (st *SearchTools) shouldIgnore(path string) bool {
	for _, pattern := range st.IgnorePatterns {
		if strings.Contains(path, pattern) {
			return true
		}
	}
	return false
}

// CreateSearchFilesTool 创建文件搜索工具（按 glob 模式）
func (st *SearchTools) CreateSearchFilesTool() tool.FunctionTool {
	return tool.FunctionTool{
		Name:        "search_files",
		Description: "按 glob 模式搜索文件。例如: *.go 匹配所有 Go 文件，**/test_*.py 匹配所有测试文件。",
		ParamsJSONSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"pattern": map[string]any{
					"type":        "string",
					"description": "glob 模式（如 *.go, **/*.ts, src/**/*.java）",
				},
				"dir": map[string]any{
					"type":        "string",
					"description": "搜索目录（相对于项目根目录，可选，默认为项目根目录）",
				},
			},
			"required": []string{"pattern"},
		},
		OnInvokeTool: func(ctx context.Context, arguments string) (any, error) {
			var params struct {
				Pattern string `json:"pattern"`
				Dir     string `json:"dir"`
			}
			if err := json.Unmarshal([]byte(arguments), &params); err != nil {
				return nil, fmt.Errorf("参数解析失败: %w", err)
			}

			// 确定搜索目录
			searchDir := st.ProjectRoot
			if params.Dir != "" {
				searchDir = filepath.Join(st.ProjectRoot, params.Dir)
			}

			var matches []string
			count := 0

			err := filepath.Walk(searchDir, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return nil // 忽略无法访问的文件
				}

				// 检查是否应该忽略
				if st.shouldIgnore(path) {
					if info.IsDir() {
						return filepath.SkipDir
					}
					return nil
				}

				// 跳过目录
				if info.IsDir() {
					return nil
				}

				// 匹配文件名
				name := filepath.Base(path)
				matched, _ := filepath.Match(params.Pattern, name)
				if matched {
					// 返回相对路径
					relPath, _ := filepath.Rel(st.ProjectRoot, path)
					matches = append(matches, relPath)
					count++
					if count >= st.MaxResults {
						return filepath.SkipAll
					}
				}

				return nil
			})

			if err != nil && err != filepath.SkipAll {
				return nil, fmt.Errorf("搜索失败: %w", err)
			}

			// 格式化输出
			var result strings.Builder
			result.WriteString(fmt.Sprintf("搜索模式: %s\n", params.Pattern))
			result.WriteString(fmt.Sprintf("找到 %d 个匹配文件:\n", len(matches)))

			for _, match := range matches {
				result.WriteString(fmt.Sprintf("  %s\n", match))
			}

			if count >= st.MaxResults {
				result.WriteString(fmt.Sprintf("\n(结果已截断，最多显示 %d 个)\n", st.MaxResults))
			}

			return result.String(), nil
		},
	}
}

// CreateSearchContentTool 创建内容搜索工具
func (st *SearchTools) CreateSearchContentTool() tool.FunctionTool {
	return tool.FunctionTool{
		Name:        "search_content",
		Description: "在文件中搜索指定文本内容。返回包含匹配内容的文件和行号。",
		ParamsJSONSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"query": map[string]any{
					"type":        "string",
					"description": "要搜索的文本内容",
				},
				"file_pattern": map[string]any{
					"type":        "string",
					"description": "文件模式过滤（可选，如 *.go 只搜索 Go 文件）",
				},
				"case_sensitive": map[string]any{
					"type":        "boolean",
					"description": "是否区分大小写（默认 false）",
				},
			},
			"required": []string{"query"},
		},
		OnInvokeTool: func(ctx context.Context, arguments string) (any, error) {
			var params struct {
				Query         string `json:"query"`
				FilePattern   string `json:"file_pattern"`
				CaseSensitive bool   `json:"case_sensitive"`
			}
			if err := json.Unmarshal([]byte(arguments), &params); err != nil {
				return nil, fmt.Errorf("参数解析失败: %w", err)
			}

			searchQuery := params.Query
			if !params.CaseSensitive {
				searchQuery = strings.ToLower(searchQuery)
			}

			type Match struct {
				File    string
				Line    int
				Content string
			}

			var matches []Match
			count := 0

			err := filepath.Walk(st.ProjectRoot, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return nil
				}

				// 检查是否应该忽略
				if st.shouldIgnore(path) {
					if info.IsDir() {
						return filepath.SkipDir
					}
					return nil
				}

				// 跳过目录
				if info.IsDir() {
					return nil
				}

				// 检查文件模式
				if params.FilePattern != "" {
					matched, _ := filepath.Match(params.FilePattern, filepath.Base(path))
					if !matched {
						return nil
					}
				}

				// 读取文件
				file, err := os.Open(path)
				if err != nil {
					return nil
				}
				defer file.Close()

				scanner := bufio.NewScanner(file)
				lineNum := 0
				for scanner.Scan() {
					lineNum++
					line := scanner.Text()
					searchLine := line
					if !params.CaseSensitive {
						searchLine = strings.ToLower(line)
					}

					if strings.Contains(searchLine, searchQuery) {
						relPath, _ := filepath.Rel(st.ProjectRoot, path)
						matches = append(matches, Match{
							File:    relPath,
							Line:    lineNum,
							Content: line,
						})
						count++
						if count >= st.MaxResults {
							return filepath.SkipAll
						}
					}
				}

				return nil
			})

			if err != nil && err != filepath.SkipAll {
				return nil, fmt.Errorf("搜索失败: %w", err)
			}

			// 格式化输出
			var result strings.Builder
			result.WriteString(fmt.Sprintf("搜索内容: %q\n", params.Query))
			result.WriteString(fmt.Sprintf("找到 %d 处匹配:\n\n", len(matches)))

			currentFile := ""
			for _, match := range matches {
				if match.File != currentFile {
					result.WriteString(fmt.Sprintf("=== %s ===\n", match.File))
					currentFile = match.File
				}
				// 截断过长的行
				content := match.Content
				if len(content) > 100 {
					content = content[:100] + "..."
				}
				result.WriteString(fmt.Sprintf("  行 %d: %s\n", match.Line, strings.TrimSpace(content)))
			}

			if count >= st.MaxResults {
				result.WriteString(fmt.Sprintf("\n(结果已截断，最多显示 %d 个)\n", st.MaxResults))
			}

			return result.String(), nil
		},
	}
}

// CreateFindSymbolTool 创建符号搜索工具
func (st *SearchTools) CreateFindSymbolTool() tool.FunctionTool {
	return tool.FunctionTool{
		Name:        "find_symbol",
		Description: "搜索代码中的符号定义（函数、类型、变量等）。支持 Go 语言。",
		ParamsJSONSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"name": map[string]any{
					"type":        "string",
					"description": "要搜索的符号名称",
				},
				"type": map[string]any{
					"type":        "string",
					"enum":        []string{"func", "type", "var", "const", "any"},
					"description": "符号类型（默认 any 表示所有类型）",
				},
			},
			"required": []string{"name"},
		},
		OnInvokeTool: func(ctx context.Context, arguments string) (any, error) {
			var params struct {
				Name       string `json:"name"`
				SymbolType string `json:"type"`
			}
			if err := json.Unmarshal([]byte(arguments), &params); err != nil {
				return nil, fmt.Errorf("参数解析失败: %w", err)
			}

			if params.SymbolType == "" {
				params.SymbolType = "any"
			}

			// 构建搜索模式
			var patterns []string
			switch params.SymbolType {
			case "func":
				patterns = []string{fmt.Sprintf("func %s", params.Name), fmt.Sprintf("func (%s", params.Name)}
			case "type":
				patterns = []string{fmt.Sprintf("type %s ", params.Name)}
			case "var":
				patterns = []string{fmt.Sprintf("var %s ", params.Name), fmt.Sprintf("%s :=", params.Name)}
			case "const":
				patterns = []string{fmt.Sprintf("const %s ", params.Name)}
			default:
				patterns = []string{
					fmt.Sprintf("func %s", params.Name),
					fmt.Sprintf("type %s ", params.Name),
					fmt.Sprintf("var %s ", params.Name),
					fmt.Sprintf("const %s ", params.Name),
				}
			}

			type Match struct {
				File    string
				Line    int
				Content string
			}

			var matches []Match

			err := filepath.Walk(st.ProjectRoot, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return nil
				}

				// 只搜索 Go 文件
				if info.IsDir() {
					if st.shouldIgnore(path) {
						return filepath.SkipDir
					}
					return nil
				}

				if !strings.HasSuffix(path, ".go") {
					return nil
				}

				// 读取文件
				file, err := os.Open(path)
				if err != nil {
					return nil
				}
				defer file.Close()

				scanner := bufio.NewScanner(file)
				lineNum := 0
				for scanner.Scan() {
					lineNum++
					line := scanner.Text()

					for _, pattern := range patterns {
						if strings.Contains(line, pattern) {
							relPath, _ := filepath.Rel(st.ProjectRoot, path)
							matches = append(matches, Match{
								File:    relPath,
								Line:    lineNum,
								Content: line,
							})
							break
						}
					}
				}

				return nil
			})

			if err != nil {
				return nil, fmt.Errorf("搜索失败: %w", err)
			}

			// 格式化输出
			var result strings.Builder
			result.WriteString(fmt.Sprintf("搜索符号: %s (类型: %s)\n", params.Name, params.SymbolType))
			result.WriteString(fmt.Sprintf("找到 %d 处定义:\n\n", len(matches)))

			for _, match := range matches {
				result.WriteString(fmt.Sprintf("%s:%d\n", match.File, match.Line))
				result.WriteString(fmt.Sprintf("  %s\n\n", strings.TrimSpace(match.Content)))
			}

			return result.String(), nil
		},
	}
}

// GetAllTools 返回所有搜索工具
func (st *SearchTools) GetAllTools() []tool.FunctionTool {
	return []tool.FunctionTool{
		st.CreateSearchFilesTool(),
		st.CreateSearchContentTool(),
		st.CreateFindSymbolTool(),
	}
}
