package converter

import (
	"encoding/json"
	"strings"
	"testing"

	"Darvin-Resume/internal/model"
)

func TestConvertBasicInfo(t *testing.T) {
	resume := &model.Resume{
		BasicInfo: model.BasicInfo{
			Name:    "张三",
			Phone:   "13800138000",
			Email:   "zhangsan@example.com",
			GitHub:  "github.com/zhangsan",
			Summary: "全栈工程师",
		},
		Modules: []model.Module{},
	}

	md, err := ConvertResumeToMarkdown(resume)
	if err != nil {
		t.Fatalf("Convert failed: %v", err)
	}

	// Verify name as title
	if !strings.Contains(md, "# 张三") {
		t.Error("Markdown should contain name as title")
	}

	// Verify contact info line
	if !strings.Contains(md, "全栈工程师") {
		t.Error("Markdown should contain summary")
	}
	if !strings.Contains(md, "13800138000") {
		t.Error("Markdown should contain phone")
	}
	if !strings.Contains(md, "zhangsan@example.com") {
		t.Error("Markdown should contain email")
	}
	if !strings.Contains(md, "github.com/zhangsan") {
		t.Error("Markdown should contain GitHub")
	}
}

func TestConvertEducation(t *testing.T) {
	itemsJSON, _ := json.Marshal([]model.EducationItem{
		{
			School:     "清华大学",
			Degree:     "本科",
			Major:      "计算机科学与技术",
			StartDate:  "2018-09",
			EndDate:    "2022-06",
			GPA:        "3.8/4.0",
			Highlights: []string{"GPA排名前10%", "获得国家奖学金"},
		},
	})

	resume := &model.Resume{
		BasicInfo: model.BasicInfo{Name: "测试"},
		Modules: []model.Module{
			{
				Type:    "education",
				Title:   "教育经历",
				Order:   1,
				Visible: true,
				Items:   itemsJSON,
			},
		},
	}

	md, err := ConvertResumeToMarkdown(resume)
	if err != nil {
		t.Fatalf("Convert failed: %v", err)
	}

	// Verify school header
	if !strings.Contains(md, "### 清华大学") {
		t.Error("Should contain school name")
	}
	if !strings.Contains(md, "计算机科学与技术") {
		t.Error("Should contain major")
	}
	if !strings.Contains(md, "本科") {
		t.Error("Should contain degree")
	}

	// Verify date range
	if !strings.Contains(md, "2018-09 - 2022-06") {
		t.Error("Should contain date range")
	}

	// Verify GPA
	if !strings.Contains(md, "GPA: 3.8/4.0") {
		t.Error("Should contain GPA")
	}

	// Verify highlights
	if !strings.Contains(md, "GPA排名前10%") {
		t.Error("Should contain first highlight")
	}
	if !strings.Contains(md, "获得国家奖学金") {
		t.Error("Should contain second highlight")
	}
}

func TestConvertSkills(t *testing.T) {
	itemsJSON, _ := json.Marshal([]model.SkillItem{
		{
			Category: "编程语言",
			Skills:   []string{"Go", "Python", "TypeScript"},
		},
		{
			Category: "框架",
			Skills:   []string{"Vue", "React", "Fiber"},
		},
	})

	resume := &model.Resume{
		BasicInfo: model.BasicInfo{Name: "测试"},
		Modules: []model.Module{
			{
				Type:    "skills",
				Title:   "专业技能",
				Order:   2,
				Visible: true,
				Items:   itemsJSON,
			},
		},
	}

	md, err := ConvertResumeToMarkdown(resume)
	if err != nil {
		t.Fatalf("Convert failed: %v", err)
	}

	// Verify skills section
	if !strings.Contains(md, "## 专业技能") {
		t.Error("Should contain skills section header")
	}
	if !strings.Contains(md, "**编程语言**: Go、Python、TypeScript") {
		t.Error("Should contain first skill category")
	}
	if !strings.Contains(md, "**框架**: Vue、React、Fiber") {
		t.Error("Should contain second skill category")
	}
}

func TestConvertProjects(t *testing.T) {
	itemsJSON, _ := json.Marshal([]model.ProjectItem{
		{
			Name:        "简历制作工具",
			Role:        "独立开发者",
			StartDate:   "2024-01",
			EndDate:     "2024-03",
			Description: "一款基于Markdown的本地简历制作工具",
			Highlights:  []string{"支持实时预览", "支持PDF导出", "支持AI润色"},
			TechStack:   []string{"Go", "Vue3", "SQLite"},
		},
	})

	resume := &model.Resume{
		BasicInfo: model.BasicInfo{Name: "测试"},
		Modules: []model.Module{
			{
				Type:    "projects",
				Title:   "项目经历",
				Order:   3,
				Visible: true,
				Items:   itemsJSON,
			},
		},
	}

	md, err := ConvertResumeToMarkdown(resume)
	if err != nil {
		t.Fatalf("Convert failed: %v", err)
	}

	// Verify project header
	if !strings.Contains(md, "### 简历制作工具 — 独立开发者") {
		t.Error("Should contain project name and role")
	}

	// Verify date range
	if !strings.Contains(md, "2024-01 - 2024-03") {
		t.Error("Should contain date range")
	}

	// Verify description
	if !strings.Contains(md, "一款基于Markdown的本地简历制作工具") {
		t.Error("Should contain description")
	}

	// Verify highlights
	if !strings.Contains(md, "- 支持实时预览") {
		t.Error("Should contain first highlight")
	}

	// Verify tech stack
	if !strings.Contains(md, "**技术栈**: Go、Vue3、SQLite") {
		t.Error("Should contain tech stack")
	}
}

func TestConvertInternship(t *testing.T) {
	itemsJSON, _ := json.Marshal([]model.InternshipItem{
		{
			Company:     "字节跳动",
			Position:    "后端开发实习生",
			StartDate:   "2023-06",
			EndDate:     "2023-09",
			Description: "参与抖音后端服务开发",
			Highlights:  []string{"优化接口响应时间50%", "独立负责活动模块"},
		},
	})

	resume := &model.Resume{
		BasicInfo: model.BasicInfo{Name: "测试"},
		Modules: []model.Module{
			{
				Type:    "internship",
				Title:   "实习经历",
				Order:   4,
				Visible: true,
				Items:   itemsJSON,
			},
		},
	}

	md, err := ConvertResumeToMarkdown(resume)
	if err != nil {
		t.Fatalf("Convert failed: %v", err)
	}

	// Verify company and position
	if !strings.Contains(md, "### 字节跳动 — 后端开发实习生") {
		t.Error("Should contain company and position")
	}

	// Verify date range
	if !strings.Contains(md, "2023-06 - 2023-09") {
		t.Error("Should contain date range")
	}

	// Verify highlights
	if !strings.Contains(md, "- 优化接口响应时间50%") {
		t.Error("Should contain first highlight")
	}
}

func TestConvertCampus(t *testing.T) {
	itemsJSON, _ := json.Marshal([]model.CampusItem{
		{
			Name:        "ACM程序设计竞赛",
			Role:        "队长",
			StartDate:   "2020-09",
			EndDate:     "2021-06",
			Description: "组织团队每周训练，参与省级比赛",
			Highlights:  []string{"省级一等奖", "校内赛冠军"},
		},
	})

	resume := &model.Resume{
		BasicInfo: model.BasicInfo{Name: "测试"},
		Modules: []model.Module{
			{
				Type:    "campus",
				Title:   "校园经历",
				Order:   5,
				Visible: true,
				Items:   itemsJSON,
			},
		},
	}

	md, err := ConvertResumeToMarkdown(resume)
	if err != nil {
		t.Fatalf("Convert failed: %v", err)
	}

	// Verify campus activity
	if !strings.Contains(md, "### ACM程序设计竞赛 — 队长") {
		t.Error("Should contain activity name and role")
	}

	// Verify highlights
	if !strings.Contains(md, "- 省级一等奖") {
		t.Error("Should contain first highlight")
	}
}

func TestConvertAwards(t *testing.T) {
	itemsJSON, _ := json.Marshal([]model.AwardItem{
		{
			Name:        "国家奖学金",
			Level:       "国家级",
			Date:        "2021-10",
			Description: "全校仅10人获得",
		},
		{
			Name:  "优秀学生干部",
			Level: "校级",
			Date:  "2022-05",
		},
	})

	resume := &model.Resume{
		BasicInfo: model.BasicInfo{Name: "测试"},
		Modules: []model.Module{
			{
				Type:    "awards",
				Title:   "荣誉奖项",
				Order:   6,
				Visible: true,
				Items:   itemsJSON,
			},
		},
	}

	md, err := ConvertResumeToMarkdown(resume)
	if err != nil {
		t.Fatalf("Convert failed: %v", err)
	}

	// Verify awards
	if !strings.Contains(md, "**国家奖学金** (国家级) - 2021-10") {
		t.Error("Should contain first award with level and date")
	}
	if !strings.Contains(md, "全校仅10人获得") {
		t.Error("Should contain award description")
	}
	if !strings.Contains(md, "**优秀学生干部** (校级)") {
		t.Error("Should contain second award")
	}
}

func TestConvertCertificates(t *testing.T) {
	itemsJSON, _ := json.Marshal([]model.CertificateItem{
		{
			Name:   "英语四级",
			Issuer: "教育部",
			Date:   "2020-09",
			Score:  "580分",
		},
		{
			Name:   "云计算工程师",
			Issuer: "阿里云",
			Date:   "2023-03",
		},
	})

	resume := &model.Resume{
		BasicInfo: model.BasicInfo{Name: "测试"},
		Modules: []model.Module{
			{
				Type:    "certificates",
				Title:   "证书",
				Order:   7,
				Visible: true,
				Items:   itemsJSON,
			},
		},
	}

	md, err := ConvertResumeToMarkdown(resume)
	if err != nil {
		t.Fatalf("Convert failed: %v", err)
	}

	// Verify certificates
	if !strings.Contains(md, "**英语四级** - 教育部 (2020-09) - 580分") {
		t.Error("Should contain first certificate with issuer, date, score")
	}
	if !strings.Contains(md, "**云计算工程师** - 阿里云") {
		t.Error("Should contain second certificate")
	}
}

func TestConvertEvaluation(t *testing.T) {
	itemsJSON, _ := json.Marshal([]model.EvaluationItem{
		{
			Content: "热爱技术，具有良好的编程习惯和代码风格。性格开朗，善于沟通，具有团队合作精神。",
		},
	})

	resume := &model.Resume{
		BasicInfo: model.BasicInfo{Name: "测试"},
		Modules: []model.Module{
			{
				Type:    "evaluation",
				Title:   "自我评价",
				Order:   8,
				Visible: true,
				Items:   itemsJSON,
			},
		},
	}

	md, err := ConvertResumeToMarkdown(resume)
	if err != nil {
		t.Fatalf("Convert failed: %v", err)
	}

	// Verify evaluation content
	if !strings.Contains(md, "## 自我评价") {
		t.Error("Should contain evaluation section header")
	}
	if !strings.Contains(md, "热爱技术，具有良好的编程习惯和代码风格") {
		t.Error("Should contain evaluation content")
	}
}

func TestModuleOrdering(t *testing.T) {
	itemsJSON, _ := json.Marshal([]model.EducationItem{
		{School: "第二学校"},
	})

	// Create modules in non-sorted order
	resume := &model.Resume{
		BasicInfo: model.BasicInfo{Name: "测试"},
		Modules: []model.Module{
			{
				Type:    "education",
				Title:   "教育经历",
				Order:   10, // Later order
				Visible: true,
				Items:   itemsJSON,
			},
			{
				Type:    "skills",
				Title:   "专业技能",
				Order:   5, // Earlier order
				Visible: true,
				Items:   json.RawMessage(`[]`),
			},
		},
	}

	md, err := ConvertResumeToMarkdown(resume)
	if err != nil {
		t.Fatalf("Convert failed: %v", err)
	}

	// Skills (order 5) should appear before Education (order 10)
	skillsIdx := strings.Index(md, "## 专业技能")
	eduIdx := strings.Index(md, "## 教育经历")
	if skillsIdx == -1 || eduIdx == -1 {
		t.Fatal("Both sections should exist in markdown")
	}
	if skillsIdx > eduIdx {
		t.Error("Modules should be ordered by Order field, skills (5) should come before education (10)")
	}
}

func TestEmptyModule(t *testing.T) {
	itemsJSON, _ := json.Marshal([]model.EducationItem{
		{School: "清华大学"},
	})

	resume := &model.Resume{
		BasicInfo: model.BasicInfo{Name: "测试"},
		Modules: []model.Module{
			{
				Type:    "education",
				Title:   "教育经历",
				Order:   1,
				Visible: false, // Hidden module
				Items:   itemsJSON,
			},
			{
				Type:    "skills",
				Title:   "专业技能",
				Order:   2,
				Visible: true,
				Items:   json.RawMessage(`[{"category":"编程语言","skills":["Go"]}]`),
			},
		},
	}

	md, err := ConvertResumeToMarkdown(resume)
	if err != nil {
		t.Fatalf("Convert failed: %v", err)
	}

	// Hidden module should not appear
	if strings.Contains(md, "## 教育经历") {
		t.Error("Hidden module should not appear in markdown")
	}
	if strings.Contains(md, "清华大学") {
		t.Error("Hidden module content should not appear in markdown")
	}

	// Visible module should appear
	if !strings.Contains(md, "## 专业技能") {
		t.Error("Visible module should appear in markdown")
	}
}

func TestFullResume(t *testing.T) {
	// Create a complete resume with all module types
	educationJSON, _ := json.Marshal([]model.EducationItem{
		{
			School:     "清华大学",
			Degree:     "本科",
			Major:      "计算机科学",
			StartDate:  "2018-09",
			EndDate:    "2022-06",
			GPA:        "3.9/4.0",
			Highlights: []string{"GPA 3.9", "国家奖学金"},
		},
	})

	skillsJSON, _ := json.Marshal([]model.SkillItem{
		{Category: "编程语言", Skills: []string{"Go", "Python", "TypeScript"}},
	})

	projectJSON, _ := json.Marshal([]model.ProjectItem{
		{
			Name:        "简历工具",
			Role:        "开发者",
			StartDate:   "2024-01",
			EndDate:     "2024-03",
			Description: "本地简历制作工具",
			Highlights:  []string{"支持PDF导出", "支持AI润色"},
			TechStack:   []string{"Go", "Vue3"},
		},
	})

	resume := &model.Resume{
		BasicInfo: model.BasicInfo{
			Name:   "张三",
			Phone:  "13800138000",
			Email:  "zhangsan@example.com",
			GitHub: "github.com/zhangsan",
		},
		Modules: []model.Module{
			{Type: "education", Title: "教育经历", Order: 1, Visible: true, Items: educationJSON},
			{Type: "skills", Title: "专业技能", Order: 2, Visible: true, Items: skillsJSON},
			{Type: "projects", Title: "项目经历", Order: 3, Visible: true, Items: projectJSON},
		},
	}

	md, err := ConvertResumeToMarkdown(resume)
	if err != nil {
		t.Fatalf("Convert failed: %v", err)
	}

	// Verify complete structure
	if !strings.Contains(md, "# 张三") {
		t.Error("Should contain name header")
	}
	if !strings.Contains(md, "13800138000") {
		t.Error("Should contain phone")
	}
	if !strings.Contains(md, "zhangsan@example.com") {
		t.Error("Should contain email")
	}
	if !strings.Contains(md, "github.com/zhangsan") {
		t.Error("Should contain GitHub")
	}
	if !strings.Contains(md, "\n---\n") {
		t.Error("Should contain separator line")
	}
	if !strings.Contains(md, "## 教育经历") {
		t.Error("Should contain education section")
	}
	if !strings.Contains(md, "清华大学") {
		t.Error("Should contain school")
	}
	if !strings.Contains(md, "## 专业技能") {
		t.Error("Should contain skills section")
	}
	if !strings.Contains(md, "## 项目经历") {
		t.Error("Should contain projects section")
	}
	if !strings.Contains(md, "**技术栈**: Go、Vue3") {
		t.Error("Should contain tech stack")
	}
}
