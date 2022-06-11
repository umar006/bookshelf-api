package model

type Book struct {
	Id        string `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	Year      int16  `json:"year,omitempty"`
	Author    string `json:"author,omitempty"`
	Summary   string `json:"summary,omitempty"`
	Publisher string `json:"publisher,omitempty"`
	PageCount int16  `json:"pageCount,omitempty" db:"page_count"`
	ReadPage  int16  `json:"readPage,omitempty" db:"read_page"`
	Finished  *bool  `json:"finished,omitempty"`
	Reading   *bool  `json:"reading,omitempty"`
	CreatedAt string `json:"insertedAt,omitempty" db:"created_at"`
	UpdatedAt string `json:"updatedAt,omitempty" db:"updated_at"`
}
