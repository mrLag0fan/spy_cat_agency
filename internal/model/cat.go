package model

type SpyCat struct {
	ID                int     `json:"id"`
	Name              string  `json:"name"`
	ExperienceInYears int     `json:"experience_in_years"`
	Breed             string  `json:"breed"`
	Salary            float64 `json:"salary"`
}

type SalaryUpdate struct {
	Salary float64 `json:"salary"`
}
