package model

type Mission struct {
	ID        int      `json:"id"`
	CatID     int      `json:"cat_id"`
	Completed bool     `json:"completed"`
	Targets   []Target `json:"targets"`
}
