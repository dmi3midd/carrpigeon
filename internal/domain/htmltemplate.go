package domain

import "time"

type HTMLTemplate struct {
	ID        string    `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Content   string    `json:"content" db:"content"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

func (t *HTMLTemplate) Metadata() *HTMLTemplateMetadata {
	return &HTMLTemplateMetadata{
		ID:        t.ID,
		Name:      t.Name,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}
}

type HTMLTemplateMetadata struct {
	ID        string    `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
