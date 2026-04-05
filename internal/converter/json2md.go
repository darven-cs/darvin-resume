package converter

import (
	"bytes"
	"encoding/json"
	"sort"
	"strings"

	"open-resume/internal/model"
)

// ConvertResumeToMarkdown converts a resume model to markdown format
func ConvertResumeToMarkdown(resume *model.Resume) (string, error) {
	var buf bytes.Buffer

	// Basic Info Section
	if resume.BasicInfo.Name != "" {
		buf.WriteString("# ")
		buf.WriteString(resume.BasicInfo.Name)
		buf.WriteString("\n")
	}

	// Contact info line
	contactParts := []string{}
	if resume.BasicInfo.Summary != "" {
		contactParts = append(contactParts, resume.BasicInfo.Summary)
	}
	if resume.BasicInfo.Phone != "" {
		contactParts = append(contactParts, resume.BasicInfo.Phone)
	}
	if resume.BasicInfo.Email != "" {
		contactParts = append(contactParts, resume.BasicInfo.Email)
	}
	if resume.BasicInfo.GitHub != "" {
		contactParts = append(contactParts, resume.BasicInfo.GitHub)
	}
	if resume.BasicInfo.Website != "" {
		contactParts = append(contactParts, resume.BasicInfo.Website)
	}
	if len(contactParts) > 0 {
		buf.WriteString(strings.Join(contactParts, " | "))
		buf.WriteString("\n")
	}

	// Check if there are visible modules
	hasVisibleModules := false
	for _, module := range resume.Modules {
		if module.Type != "basicInfo" && module.Visible {
			hasVisibleModules = true
			break
		}
	}

	if hasVisibleModules {
		buf.WriteString("\n---\n\n")
	}

	// Sort modules by order
	sortedModules := make([]model.Module, len(resume.Modules))
	copy(sortedModules, resume.Modules)
	sort.Slice(sortedModules, func(i, j int) bool {
		return sortedModules[i].Order < sortedModules[j].Order
	})

	// Convert each visible module
	for _, module := range sortedModules {
		if !module.Visible || module.Type == "basicInfo" {
			continue
		}

		switch module.Type {
		case "education":
			convertEducation(&buf, module)
		case "skills":
			convertSkills(&buf, module)
		case "projects":
			convertProjects(&buf, module)
		case "internship":
			convertInternship(&buf, module)
		case "campus":
			convertCampus(&buf, module)
		case "awards":
			convertAwards(&buf, module)
		case "certificates":
			convertCertificates(&buf, module)
		case "evaluation":
			convertEvaluation(&buf, module)
		}
		buf.WriteString("\n")
	}

	return buf.String(), nil
}

func convertEducation(buf *bytes.Buffer, module model.Module) {
	buf.WriteString("## ")
	buf.WriteString(module.Title)
	buf.WriteString("\n\n")

	var items []model.EducationItem
	if err := json.Unmarshal(module.Items, &items); err != nil {
		return
	}

	for _, item := range items {
		// School and degree
		header := item.School
		if item.Major != "" {
			header += " — " + item.Major
		}
		if item.Degree != "" {
			header += " (" + item.Degree + ")"
		}
		buf.WriteString("### ")
		buf.WriteString(header)
		buf.WriteString("\n")

		// Date range
		dateRange := formatDateRange(item.StartDate, item.EndDate)
		if dateRange != "" {
			buf.WriteString(dateRange)
			buf.WriteString("\n")
		}

		// GPA
		if item.GPA != "" {
			buf.WriteString("GPA: ")
			buf.WriteString(item.GPA)
			buf.WriteString("\n")
		}

		// Highlights
		for _, highlight := range item.Highlights {
			buf.WriteString("- ")
			buf.WriteString(highlight)
			buf.WriteString("\n")
		}
		buf.WriteString("\n")
	}
}

func convertSkills(buf *bytes.Buffer, module model.Module) {
	buf.WriteString("## ")
	buf.WriteString(module.Title)
	buf.WriteString("\n\n")

	var items []model.SkillItem
	if err := json.Unmarshal(module.Items, &items); err != nil {
		return
	}

	for _, item := range items {
		buf.WriteString("- **")
		buf.WriteString(item.Category)
		buf.WriteString("**: ")
		buf.WriteString(strings.Join(item.Skills, "、"))
		buf.WriteString("\n")
	}
	buf.WriteString("\n")
}

func convertProjects(buf *bytes.Buffer, module model.Module) {
	buf.WriteString("## ")
	buf.WriteString(module.Title)
	buf.WriteString("\n\n")

	var items []model.ProjectItem
	if err := json.Unmarshal(module.Items, &items); err != nil {
		return
	}

	for _, item := range items {
		// Project name and role
		header := item.Name
		if item.Role != "" {
			header += " — " + item.Role
		}
		buf.WriteString("### ")
		buf.WriteString(header)
		buf.WriteString("\n")

		// Date range
		dateRange := formatDateRange(item.StartDate, item.EndDate)
		if dateRange != "" {
			buf.WriteString(dateRange)
			buf.WriteString("\n")
		}

		// Description
		if item.Description != "" {
			buf.WriteString(item.Description)
			buf.WriteString("\n")
		}

		// Highlights
		for _, highlight := range item.Highlights {
			buf.WriteString("- ")
			buf.WriteString(highlight)
			buf.WriteString("\n")
		}

		// Tech stack
		if len(item.TechStack) > 0 {
			buf.WriteString("**技术栈**: ")
			buf.WriteString(strings.Join(item.TechStack, "、"))
			buf.WriteString("\n")
		}
		buf.WriteString("\n")
	}
}

func convertInternship(buf *bytes.Buffer, module model.Module) {
	buf.WriteString("## ")
	buf.WriteString(module.Title)
	buf.WriteString("\n\n")

	var items []model.InternshipItem
	if err := json.Unmarshal(module.Items, &items); err != nil {
		return
	}

	for _, item := range items {
		// Company and position
		header := item.Company
		if item.Position != "" {
			header += " — " + item.Position
		}
		buf.WriteString("### ")
		buf.WriteString(header)
		buf.WriteString("\n")

		// Date range
		dateRange := formatDateRange(item.StartDate, item.EndDate)
		if dateRange != "" {
			buf.WriteString(dateRange)
			buf.WriteString("\n")
		}

		// Description
		if item.Description != "" {
			buf.WriteString(item.Description)
			buf.WriteString("\n")
		}

		// Highlights
		for _, highlight := range item.Highlights {
			buf.WriteString("- ")
			buf.WriteString(highlight)
			buf.WriteString("\n")
		}
		buf.WriteString("\n")
	}
}

func convertCampus(buf *bytes.Buffer, module model.Module) {
	buf.WriteString("## ")
	buf.WriteString(module.Title)
	buf.WriteString("\n\n")

	var items []model.CampusItem
	if err := json.Unmarshal(module.Items, &items); err != nil {
		return
	}

	for _, item := range items {
		// Name and role
		header := item.Name
		if item.Role != "" {
			header += " — " + item.Role
		}
		buf.WriteString("### ")
		buf.WriteString(header)
		buf.WriteString("\n")

		// Date range
		dateRange := formatDateRange(item.StartDate, item.EndDate)
		if dateRange != "" {
			buf.WriteString(dateRange)
			buf.WriteString("\n")
		}

		// Description
		if item.Description != "" {
			buf.WriteString(item.Description)
			buf.WriteString("\n")
		}

		// Highlights
		for _, highlight := range item.Highlights {
			buf.WriteString("- ")
			buf.WriteString(highlight)
			buf.WriteString("\n")
		}
		buf.WriteString("\n")
	}
}

func convertAwards(buf *bytes.Buffer, module model.Module) {
	buf.WriteString("## ")
	buf.WriteString(module.Title)
	buf.WriteString("\n\n")

	var items []model.AwardItem
	if err := json.Unmarshal(module.Items, &items); err != nil {
		return
	}

	for _, item := range items {
		buf.WriteString("- **")
		buf.WriteString(item.Name)
		buf.WriteString("**")
		if item.Level != "" {
			buf.WriteString(" (")
			buf.WriteString(item.Level)
			buf.WriteString(")")
		}
		if item.Date != "" {
			buf.WriteString(" - ")
			buf.WriteString(item.Date)
		}
		buf.WriteString("\n")
		if item.Description != "" {
			buf.WriteString("  ")
			buf.WriteString(item.Description)
			buf.WriteString("\n")
		}
	}
	buf.WriteString("\n")
}

func convertCertificates(buf *bytes.Buffer, module model.Module) {
	buf.WriteString("## ")
	buf.WriteString(module.Title)
	buf.WriteString("\n\n")

	var items []model.CertificateItem
	if err := json.Unmarshal(module.Items, &items); err != nil {
		return
	}

	for _, item := range items {
		buf.WriteString("- **")
		buf.WriteString(item.Name)
		buf.WriteString("**")
		if item.Issuer != "" {
			buf.WriteString(" - ")
			buf.WriteString(item.Issuer)
		}
		if item.Date != "" {
			buf.WriteString(" (")
			buf.WriteString(item.Date)
			buf.WriteString(")")
		}
		if item.Score != "" {
			buf.WriteString(" - ")
			buf.WriteString(item.Score)
		}
		buf.WriteString("\n")
	}
	buf.WriteString("\n")
}

func convertEvaluation(buf *bytes.Buffer, module model.Module) {
	buf.WriteString("## ")
	buf.WriteString(module.Title)
	buf.WriteString("\n\n")

	var items []model.EvaluationItem
	if err := json.Unmarshal(module.Items, &items); err != nil {
		return
	}

	for _, item := range items {
		buf.WriteString(item.Content)
		buf.WriteString("\n\n")
	}
}

func formatDateRange(startDate, endDate string) string {
	if startDate == "" && endDate == "" {
		return ""
	}
	if startDate != "" && endDate != "" {
		return startDate + " - " + endDate
	}
	if startDate != "" {
		return startDate + " - 至今"
	}
	return " - " + endDate
}
