package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"main/internal/model"
	"main/internal/store"
)

type MissionRepository struct {
	db *sql.DB
}

func NewMissionRepository(store store.Store) *MissionRepository {
	return &MissionRepository{db: store.DB}
}

// AssignCat - Призначає кота до місії
func (r *MissionRepository) AssignCat(missionID int, catID int) error {
	query := `UPDATE missions SET cat_id = $1 WHERE id = $2 AND cat_id IS NULL`
	result, err := r.db.Exec(query, catID, missionID)
	if err != nil {
		return err
	}

	rowsAffcted, err := result.RowsAffected()
	if rowsAffcted == 0 || err != nil {
		return fmt.Errorf("mission is already assigned or does not exist")
	}

	return nil
}

func (r *MissionRepository) UpdateNotes(targetID int, notes string) error {
	query := `
        UPDATE targets 
        SET notes = $1 
        WHERE id = $2 
        AND complete = FALSE 
        AND mission_id IN (SELECT id FROM missions WHERE complete = FALSE)
    `
	result, err := r.db.Exec(query, notes, targetID)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("cannot update notes: mission or target is completed")
	}

	return nil
}

// Create creates a new mission with targets in the database
func (r *MissionRepository) Create(mission *model.Mission) error {
	tx, err := r.db.BeginTx(context.Background(), nil)
	if err != nil {
		return fmt.Errorf("unable to start transaction: %v", err)
	}
	defer tx.Rollback()

	query := `INSERT INTO missions (cat_id, complete) VALUES ($1, $2) RETURNING id`
	err = tx.QueryRow(query, mission.CatID, mission.Completed).Scan(&mission.ID)
	if err != nil {
		return fmt.Errorf("unable to create mission: %v", err)
	}

	for _, target := range mission.Targets {
		targetQuery := `INSERT INTO targets (mission_id, name, country, notes, complete) VALUES ($1, $2, $3, $4, $5)`
		_, err := r.db.Exec(targetQuery, mission.ID, target.Name, target.Country, target.Notes, target.Complete)
		if err != nil {
			return fmt.Errorf("unable to create target: %v", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("unable to commit transaction: %v", err)
	}

	return nil
}

// Delete deletes a mission from the database within a transaction
func (r *MissionRepository) Delete(missionID int) error {
	tx, err := r.db.BeginTx(context.Background(), nil)
	if err != nil {
		return fmt.Errorf("unable to start transaction: %v", err)
	}
	defer tx.Rollback()

	query := `SELECT cat_id FROM missions WHERE id = $1`
	var catID int
	row := tx.QueryRow(query, missionID)
	if err := row.Scan(&catID); err != nil {
		return fmt.Errorf("unable to find mission: %v", err)
	}
	if catID != 0 {
		return fmt.Errorf("mission is already assigned to a cat and cannot be deleted")
	}

	targetQuery := `DELETE FROM targets WHERE mission_id = $1`
	_, err = tx.Exec(targetQuery, missionID)
	if err != nil {
		return fmt.Errorf("unable to delete targets of the mission: %v", err)
	}

	missionQuery := `DELETE FROM missions WHERE id = $1`
	_, err = tx.Exec(missionQuery, missionID)
	if err != nil {
		return fmt.Errorf("unable to delete mission: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("unable to commit transaction: %v", err)
	}

	return nil
}

// Update marks a mission as completed
func (r *MissionRepository) Update(missionID int) error {
	query := `UPDATE missions SET complete = $1 WHERE id = $2`
	commandTag, err := r.db.Exec(query, true, missionID)
	if err != nil {
		return fmt.Errorf("unable to update mission: %v", err)
	}
	rowsAffcted, err := commandTag.RowsAffected()
	if rowsAffcted == 0 || err != nil {
		return fmt.Errorf("no mission found with id %d", missionID)
	}

	return nil
}

func (r *MissionRepository) MarkTargetAsComplete(targetID int) error {
	query := `UPDATE targets SET complete = TRUE WHERE id = $1 AND complete = FALSE`
	result, err := r.db.Exec(query, targetID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("target not found or already completed")
	}

	return nil
}

// DeleteTarget deletes a target from a mission within a transaction
func (r *MissionRepository) DeleteTarget(targetID int) error {
	tx, err := r.db.BeginTx(context.Background(), nil)
	if err != nil {
		return fmt.Errorf("unable to start transaction: %v", err)
	}
	defer tx.Rollback()

	targetQuery := `SELECT complete FROM targets WHERE id = $1`
	var isCompleted bool
	row := tx.QueryRow(targetQuery, targetID)
	if err := row.Scan(&isCompleted); err != nil {
		return fmt.Errorf("unable to find target: %v", err)
	}
	if isCompleted {
		return fmt.Errorf("target is completed and cannot be deleted")
	}

	query := `DELETE FROM targets WHERE id = $1`
	_, err = tx.Exec(query, targetID)
	if err != nil {
		return fmt.Errorf("unable to delete target: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("unable to commit transaction: %v", err)
	}

	return nil
}

// AddTarget adds a new target to an existing mission within a transaction
func (r *MissionRepository) AddTarget(missionID int, target *model.Target) error {
	tx, err := r.db.BeginTx(context.Background(), nil)
	if err != nil {
		return fmt.Errorf("unable to start transaction: %v", err)
	}
	defer tx.Rollback()

	missionQuery := `SELECT complete FROM missions WHERE id = $1`
	var isMissionCompleted bool
	row := tx.QueryRow(missionQuery, missionID)
	if err := row.Scan(&isMissionCompleted); err != nil {
		return fmt.Errorf("unable to find mission: %v", err)
	}
	if isMissionCompleted {
		return fmt.Errorf("mission is completed and no new targets can be added")
	}

	query := `INSERT INTO targets (mission_id, name, country, notes, complete) VALUES ($1, $2, $3, $4, $5)`
	_, err = tx.Exec(query, missionID, target.Name, target.Country, target.Notes, target.Complete)
	if err != nil {
		return fmt.Errorf("unable to add target: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("unable to commit transaction: %v", err)
	}

	return nil
}

// GetAll retrieves all missions from the database
func (r *MissionRepository) GetAll() ([]model.Mission, error) {
	query := `
        SELECT 
            m.id AS mission_id, 
            m.cat_id, 
            m.complete, 
            t.id AS target_id, 
            t.name, 
            t.country, 
            t.notes, 
            t.complete AS target_complete
        FROM missions m
        LEFT JOIN targets t ON m.id = t.mission_id
        ORDER BY m.id, t.id
    `

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve missions: %v", err)
	}
	defer rows.Close()

	missionsMap := make(map[int]*model.Mission)
	for rows.Next() {
		var missionID, catID, targetID sql.NullInt32
		var complete, targetComplete sql.NullBool
		var name, country, notes sql.NullString

		if err := rows.Scan(&missionID, &catID, &complete, &targetID, &name, &country, &notes, &targetComplete); err != nil {
			return nil, fmt.Errorf("unable to scan row: %v", err)
		}

		if _, exists := missionsMap[int(missionID.Int32)]; !exists {
			missionsMap[int(missionID.Int32)] = &model.Mission{
				ID:        int(missionID.Int32),
				CatID:     int(catID.Int32),
				Completed: complete.Bool,
				Targets:   []model.Target{},
			}
		}

		if targetID.Valid {
			target := model.Target{
				ID:       int(targetID.Int32),
				Name:     name.String,
				Country:  country.String,
				Notes:    notes.String,
				Complete: targetComplete.Bool,
			}
			missionsMap[int(missionID.Int32)].Targets = append(missionsMap[int(missionID.Int32)].Targets, target)
		}
	}

	var missions []model.Mission
	for _, mission := range missionsMap {
		missions = append(missions, *mission)
	}

	return missions, nil
}

func (r *MissionRepository) GetByID(id int) (model.Mission, error) {
	var mission model.Mission

	query := `
        SELECT m.id, m.cat_id, m.complete, 
               t.id, t.name, t.country, t.notes, t.complete
        FROM missions m
        LEFT JOIN targets t ON m.id = t.mission_id
        WHERE m.id = $1
    `

	rows, err := r.db.Query(query, id)
	if err != nil {
		return model.Mission{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var target model.Target
		err := rows.Scan(&mission.ID, &mission.CatID, &mission.Completed,
			&target.ID, &target.Name, &target.Country, &target.Notes, &target.Complete)
		if err != nil {
			return model.Mission{}, err
		}

		if target.ID != 0 {
			mission.Targets = append(mission.Targets, target)
		}
	}

	if mission.ID == 0 {
		return model.Mission{}, fmt.Errorf("mission not found")
	}

	return mission, nil
}
