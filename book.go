package main

type Book struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Year      int16  `json:"year"`
	Author    string `json:"author"`
	Summary   string `json:"summary"`
	Publisher string `json:"publisher"`
	PageCount int16  `json:"pageCount" db:"page_count"`
	ReadPage  int16  `json:"readPage" db:"read_page"`
	Finished  bool   `json:"finished"`
	Reading   bool   `json:"reading"`
	CreatedAt string `json:"insertedAt" db:"created_at"`
	UpdatedAt string `json:"updatedAt" db:"updated_at"`
}
