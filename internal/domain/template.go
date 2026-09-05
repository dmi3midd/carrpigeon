package domain

import "time"

type Template struct {
	ID        string    `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Content   string    `json:"content" db:"content"`
	IsHTML    bool      `json:"is_html" db:"is_html"`
	Fields    []string  `json:"fields" db:"fields"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

func (t *Template) Metadata() *TemplateMetadata {
	return &TemplateMetadata{
		ID:        t.ID,
		Name:      t.Name,
		IsHTML:    t.IsHTML,
		Fields:    t.Fields,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}
}

type TemplateMetadata struct {
	ID        string    `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	IsHTML    bool      `json:"is_html" db:"is_html"`
	Fields    []string  `json:"fields" db:"fields"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
