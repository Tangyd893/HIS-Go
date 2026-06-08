package assistant

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Department 本院科室
type Department struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// deptAPIResponse 用户服务 /api/user/departments 响应
type deptAPIResponse struct {
	Code    int          `json:"code"`
	Message string       `json:"message"`
	Data    []Department `json:"data"`
}

// DepartmentResolver 科室解析器（跨服务调用 user 服务）
type DepartmentResolver struct {
	userServiceURL string
	httpCli        *http.Client
}

// NewDepartmentResolver 创建科室解析器
func NewDepartmentResolver(userServiceURL string) *DepartmentResolver {
	return &DepartmentResolver{
		userServiceURL: strings.TrimRight(userServiceURL, "/"),
		httpCli: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// ListAll 获取本院全部科室
func (r *DepartmentResolver) ListAll() ([]Department, error) {
	url := r.userServiceURL + "/api/user/departments"
	resp, err := r.httpCli.Get(url)
	if err != nil {
		return nil, fmt.Errorf("调用用户服务失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	var apiResp deptAPIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("解析科室数据失败: %w", err)
	}

	if apiResp.Code != 0 {
		return nil, fmt.Errorf("用户服务返回错误: %s", apiResp.Message)
	}

	return apiResp.Data, nil
}

// MatchDeptTypes 将知识库 dept_types 与本院科室做交集匹配
// deptTypes: 知识库推荐的科室类型（如 ["内科", "呼吸内科"]）
// 返回：本院中匹配的科室名称列表
func (r *DepartmentResolver) MatchDeptTypes(deptTypes []string) ([]Department, error) {
	allDepts, err := r.ListAll()
	if err != nil {
		return nil, err
	}

	var matched []Department
	for _, dt := range deptTypes {
		dtLower := strings.ToLower(strings.TrimSpace(dt))
		for _, dept := range allDepts {
			deptNameLower := strings.ToLower(dept.Name)
			// 模糊匹配：科室名包含 dept_type 或 dept_type 包含科室名
			if strings.Contains(deptNameLower, dtLower) || strings.Contains(dtLower, deptNameLower) {
				// 去重
				found := false
				for _, m := range matched {
					if m.ID == dept.ID {
						found = true
						break
					}
				}
				if !found {
					matched = append(matched, dept)
				}
			}
		}
	}

	return matched, nil
}
