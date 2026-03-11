package models

type TemplateStatus string

const (
	TemplateStatusDraft     TemplateStatus = "draft"
	TemplateStatusPublished TemplateStatus = "published"
)

func (s TemplateStatus) IsValid() bool {
	switch s {
	case TemplateStatusDraft, TemplateStatusPublished:
		return true
	default:
		return false
	}
}
