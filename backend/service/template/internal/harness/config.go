package harness

import "gitlab.com/alienspaces/go-mud/backend/service/template/internal/record"

// DataConfig -
type DataConfig struct {
	TemplateConfig []TemplateConfig
}

// TemplateConfig -
type TemplateConfig struct {
	Record *record.Template
}
