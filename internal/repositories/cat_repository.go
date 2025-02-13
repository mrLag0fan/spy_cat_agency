package repositories

import (
	"database/sql"
	"fmt"
	"main/internal/model"
	"main/internal/store"
)

type CatRepository struct {
	db *sql.DB
}

func NewCatRepository(store store.Store) *CatRepository {
	return &CatRepository{db: store.DB}
}

// Create creates a new cat in the database
func (r *CatRepository) Create(cat *model.SpyCat) error {
	query := `INSERT INTO cats (name, years_of_experience, breed, salary) VALUES ($1, $2, $3, $4) RETURNING id`
	err := r.db.QueryRow(query, cat.Name, cat.ExperienceInYears, cat.Breed, cat.Salary).Scan(&cat.ID)
	if err != nil {
		return fmt.Errorf("unable to create cat: %v", err)
	}
	return nil
}

// GetAll retrieves all cats from the database
func (r *CatRepository) GetAll() ([]model.SpyCat, error) {
	rows, err := r.db.Query("SELECT id, name, years_of_experience, breed, salary FROM cats")
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve cats: %v", err)
	}
	defer rows.Close()

	var cats []model.SpyCat
	for rows.Next() {
		var cat model.SpyCat
		if err := rows.Scan(&cat.ID, &cat.Name, &cat.ExperienceInYears, &cat.Breed, &cat.Salary); err != nil {
			return nil, fmt.Errorf("unable to scan cat: %v", err)
		}
		cats = append(cats, cat)
	}
	return cats, nil
}

// Delete removes a spy cat from the database
func (r *CatRepository) Delete(catID int) error {
	query := `DELETE FROM cats WHERE id = $1`
	commandTag, err := r.db.Exec(query, catID)
	if err != nil {
		return fmt.Errorf("unable to delete cat: %v", err)
	}
	rowsAffcted, err := commandTag.RowsAffected()
	if rowsAffcted == 0 || err != nil {
		return fmt.Errorf("no cat found with id %d", catID)
	}
	return nil
}

// UpdateSalary updates the salary of a spy cat in the database
func (r *CatRepository) UpdateSalary(catID int, newSalary float64) error {
	query := `UPDATE cats SET salary = $1 WHERE id = $2`
	commandTag, err := r.db.Exec(query, newSalary, catID)
	if err != nil {
		return fmt.Errorf("unable to update salary for cat with id %d: %v", catID, err)
	}
	rowsAffcted, err := commandTag.RowsAffected()
	if rowsAffcted == 0 || err != nil {
		return fmt.Errorf("no cat found with id %d", catID)
	}
	return nil
}

// GetByID retrieves a single spy cat from the database by its ID
func (r *CatRepository) GetByID(catID int) (*model.SpyCat, error) {
	query := `SELECT id, name, years_of_experience, breed, salary FROM cats WHERE id = $1`
	row := r.db.QueryRow(query, catID)

	var cat model.SpyCat
	if err := row.Scan(&cat.ID, &cat.Name, &cat.ExperienceInYears, &cat.Breed, &cat.Salary); err != nil {
		if err.Error() == "no rows in result set" {
			return nil, fmt.Errorf("no cat found with id %d", catID)
		}
		return nil, fmt.Errorf("unable to retrieve cat with id %d: %v", catID, err)
	}
	return &cat, nil
}
