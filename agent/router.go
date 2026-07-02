package agent

import (
	"strings"
)

func RouteTool(query string, toolService *ToolService) string {
	if strings.Contains(query, "posts") {
		return toolService.GetPosts()
	}

	if strings.Contains(query, "post") {
		return toolService.GetPostByID("26")
	}
	return "没找到工具"
}
