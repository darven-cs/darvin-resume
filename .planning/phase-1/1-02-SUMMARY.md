# Phase 1 Plan 01-02: SQLite 存储层与简历 JSON Schema

## 执行摘要

SQLite 数据库正确初始化，简历结构化 JSON Schema 定义完整，CRUD 操作可用。

## 目标

SQLite 数据库正确初始化，简历结构化 JSON Schema 定义完整，CRUD 操作可用

### 依赖
- Plan 01-01 完成（项目骨架存在）

## 已完成任务

| Task | Name | Commit | Files |
| ---- | ---- | ------ | ----- |
| T1 | 数据库连接管理 | 25b2674 | internal/database/db.go |
| T2 | 数据库迁移 | 25b2674 | internal/database/migrations/*.sql |
| T3 | 简历 JSON Schema | 25b2674 | internal/model/resume.go |
| T4 | CRUD 操作 | 25b2674 | internal/service/resume.go |
| T5 | 应用启动初始化 | 25b2674 | app.go |
| T6 | 单元测试-CRUD | 25b2674 | internal/service/resume_test.go |
| T7 | 单元测试-转换 | 25b2674 | internal/converter/json2md_test.go |

## 技术实现

### 数据库连接管理 (internal/database/db.go)
- 使用 modernc.org/sqlite 纯 Go SQLite 驱动
- 数据库文件路径：用户数据目录/open-resume/data.db
  - Linux: ~/.local/share/open-resume/
  - macOS: ~/Library/Application Support/open-resume/
  - Windows: %APPDATA%/open-resume/
- 启用 WAL 模式、外键约束
- 连接池配置（MaxOpenConns=1, MaxIdleConns=1）
- 优雅关闭
- 自动运行 goose 迁移

### 数据库迁移 (internal/database/migrations/)
使用 goose 管理迁移：
- 001_create_resumes_table.sql: 创建 resumes 表
- 002_create_settings_table.sql: 创建 settings 表

### 简历 JSON Schema (internal/model/resume.go)
完整定义所有数据结构：
- Resume, BasicInfo, Module
- EducationItem, SkillItem, ProjectItem
- InternshipItem, CampusItem, AwardItem
- CertificateItem, EvaluationItem
- ResumeListItem

### CRUD 操作 (internal/service/resume.go)
```go
type ResumeService interface {
    Create(ctx context.Context, resume *Resume) error
    GetByID(ctx context.Context, id string) (*Resume, error)
    List(ctx context.Context) ([]*ResumeListItem, error)
    Update(ctx context.Context, resume *Resume) error
    Delete(ctx context.Context, id string) error
    UpdateJSON(ctx context.Context, id string, jsonData string) error
}
```

### JSON→Markdown 转换 (internal/converter/json2md.go)
支持所有模块类型的转换：
- BasicInfo: 姓名标题、联系方式行
- Education: 学校/专业/学位/时间/亮点
- Skills: 分类加粗、技能顿号分隔
- Projects: 项目名/角色/时间/描述/技术栈
- Internship, Campus, Awards, Certificates, Evaluation

## 测试结果

```
go test ./internal/service/ -v -count=1
=== RUN   TestCreateResume      --- PASS
=== RUN   TestGetResume         --- PASS
=== RUN   TestListResumes       --- PASS
=== RUN   TestUpdateResume      --- PASS
=== RUN   TestDeleteResume      --- PASS
=== RUN   TestUpdateJSON        --- PASS
PASS (6 tests)

go test ./internal/converter/ -v -count=1
=== RUN   TestConvertBasicInfo      --- PASS
=== RUN   TestConvertEducation      --- PASS
=== RUN   TestConvertSkills         --- PASS
=== RUN   TestConvertProjects       --- PASS
=== RUN   TestConvertInternship     --- PASS
=== RUN   TestConvertCampus        --- PASS
=== RUN   TestConvertAwards        --- PASS
=== RUN   TestConvertCertificates   --- PASS
=== RUN   TestConvertEvaluation    --- PASS
=== RUN   TestModuleOrdering       --- PASS
=== RUN   TestEmptyModule          --- PASS
=== RUN   TestFullResume           --- PASS
PASS (12 tests)
```

## 偏差记录

### Auto-fixed Issues

**1. [Rule 2 - Missing Functionality] BasicInfo 序列化缺失**
- **发现于:** Task 4 (CRUD 实现)
- **问题:** Create 和 Update 函数没有将 BasicInfo 正确序列化到 json_data 字段
- **修复:** 添加 BasicInfo 作为 basicInfo 模块存储，GetByID 时提取还原
- **Commit:** 25b2674

## 关键文件

| File | Purpose |
| ---- | ------- |
| internal/database/db.go | 数据库连接、迁移初始化 |
| internal/database/migrations/*.sql | goose 迁移文件 |
| internal/model/resume.go | 简历数据结构 |
| internal/service/resume.go | CRUD 业务逻辑 |
| internal/converter/json2md.go | JSON→Markdown 转换 |
| app.go | 应用启动时初始化数据库 |

## 下一步

Plan 01-03: Bridge 层绑定、JSON-Markdown 正向同步、前端框架
