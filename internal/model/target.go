package model

type Target struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Country  string `json:"country"`
	Notes    string `json:"notes"`
	Complete bool   `json:"complete"`
}

type NoteUpdate struct {
	Notes string `json:"notes"`
}
