package handlers

import (
	"main/internal/model"
	"main/internal/repositories"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MissionHandler struct {
	MissionRepo *repositories.MissionRepository
}

func NewMissionHandler(missionRepo *repositories.MissionRepository) *MissionHandler {
	return &MissionHandler{MissionRepo: missionRepo}
}

// CreateMission godoc
// @Summary Create a mission with targets
// @Description Create a new mission and its associated targets
// @Tags missions
// @Accept json
// @Produce json
// @Param mission body model.Mission true "Mission details with targets"
// @Success 201 {object} model.Mission "Mission created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 500 {object} map[string]interface{} "Failed to create mission"
// @Router /mission [post]
func (h *MissionHandler) CreateMission(c *gin.Context) {
	var mission model.Mission
	if err := c.ShouldBindJSON(&mission); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	err := h.MissionRepo.Create(&mission)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create mission"})
		return
	}

	c.JSON(http.StatusCreated, mission)
}

// DeleteMission godoc
// @Summary Delete a mission (only if it’s not assigned to a cat)
// @Description Delete a mission, but only if it's not assigned to a cat. Returns an error if the mission is assigned to a cat.
// @Tags missions
// @Produce json
// @Param id path int true "Mission ID"
// @Success 200 {object} map[string]interface{} "Mission deleted successfully"
// @Failure 400 {object} map[string]interface{} "Invalid mission ID"
// @Failure 500 {object} map[string]interface{} "Failed to delete mission"
// @Router /mission/{id} [delete]
func (h *MissionHandler) DeleteMission(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid mission ID"})
		return
	}

	err = h.MissionRepo.Delete(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete mission"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Mission deleted successfully"})
}

// CompleteMission godoc
// @Summary Mark a mission as complete
// @Description Mark a mission as completed in the system.
// @Tags missions
// @Produce json
// @Param id path int true "Mission ID"
// @Success 200 {object} map[string]interface{} "Mission marked as complete"
// @Failure 400 {object} map[string]interface{} "Invalid mission ID"
// @Failure 500 {object} map[string]interface{} "Failed to complete mission"
// @Router /mission/{id}/complete [put]
func (h *MissionHandler) CompleteMission(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid mission ID"})
		return
	}

	err = h.MissionRepo.Update(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to complete mission"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Mission marked as complete"})
}

// UpdateTargetNotes godoc
// @Summary Update notes for a target (only if not completed)
// @Description Update the notes for a mission target if it has not been marked as complete.
// @Tags missions
// @Produce json
// @Param target_id path int true "Target ID"
// @Param notes body model.NoteUpdate true "Updated notes"
// @Success 200 {object} map[string]interface{} "Notes updated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid target ID or request body"
// @Failure 500 {object} map[string]interface{} "Failed to update notes"
// @Router /mission/targets/{target_id}/notes [put]
func (h *MissionHandler) UpdateTargetNotes(c *gin.Context) {
	targetID, err := strconv.Atoi(c.Param("target_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid target ID"})
		return
	}

	var noteUpdate model.NoteUpdate
	if err := c.ShouldBindJSON(&noteUpdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	err = h.MissionRepo.UpdateNotes(targetID, noteUpdate.Notes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update notes"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notes updated successfully"})
}

// MarkTargetAsComplete godoc
// @Summary Mark a mission target as complete
// @Description Marks a specified mission target as complete if found.
// @Tags missions
// @Produce json
// @Param target_id path int true "Target ID"
// @Success 200 {object} map[string]interface{} "Target marked as complete"
// @Failure 400 {object} map[string]interface{} "Invalid target ID"
// @Failure 404 {object} map[string]interface{} "Target not found"
// @Router /mission/targets/{target_id}/complete [put]
func (h *MissionHandler) MarkTargetAsComplete(c *gin.Context) {
	targetID, err := strconv.Atoi(c.Param("target_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid target ID"})
		return
	}

	err = h.MissionRepo.MarkTargetAsComplete(targetID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "target marked as complete"})
}

// DeleteTarget godoc
// @Summary Delete a target from a mission
// @Description Deletes a specified target from a mission by its ID.
// @Tags missions
// @Produce json
// @Param target_id path int true "Target ID"
// @Success 200 {object} map[string]interface{} "Target deleted successfully"
// @Failure 400 {object} map[string]interface{} "Invalid target ID"
// @Failure 500 {object} map[string]interface{} "Failed to delete target"
// @Router /mission/targets/{target_id} [delete]
func (h *MissionHandler) DeleteTarget(c *gin.Context) {
	targetID, err := strconv.Atoi(c.Param("target_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid target ID"})
		return
	}

	err = h.MissionRepo.DeleteTarget(targetID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete target"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Target deleted successfully"})
}

// AddTarget godoc
// @Summary Add a target to an existing mission
// @Description Adds a new target to a specified mission by its ID.
// @Tags missions
// @Produce json
// @Param id path int true "Mission ID"
// @Param target body model.Target true "Target to add"
// @Success 200 {object} map[string]interface{} "Target added successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 400 {object} map[string]interface{} "Invalid mission ID"
// @Failure 500 {object} map[string]interface{} "Failed to add target"
// @Router /mission/{id}/targets [post]
func (h *MissionHandler) AddTarget(c *gin.Context) {
	missionID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid mission ID"})
		return
	}

	var target model.Target
	if err := c.ShouldBindJSON(&target); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	err = h.MissionRepo.AddTarget(missionID, &target)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add target"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Target added successfully"})
}

// AssignCatToMission godoc
// @Summary Assign a cat to a mission
// @Description Assigns a specified cat to an existing mission by its ID.
// @Tags missions
// @Produce json
// @Param id path int true "Mission ID"
// @Param cat_id body int true "Cat ID"
// @Success 200 {object} map[string]interface{} "Cat assigned to mission"
// @Failure 400 {object} map[string]interface{} "Invalid request body or mission ID"
// @Failure 500 {object} map[string]interface{} "Failed to assign cat"
// @Router /mission/{id}/assign-cat [post]
func (h *MissionHandler) AssignCatToMission(c *gin.Context) {
	missionID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid mission ID"})
		return
	}

	var assignData struct {
		CatID int `json:"cat_id"`
	}
	if err := c.ShouldBindJSON(&assignData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	err = h.MissionRepo.AssignCat(missionID, assignData.CatID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign cat"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cat assigned to mission"})
}

// GetAllMissions godoc
// @Summary Get all missions
// @Description Retrieves a list of all missions
// @Tags missions
// @Produce json
// @Success 200 {array} model.Mission "List of missions"
// @Failure 500 {object} map[string]interface{} "Failed to retrieve missions"
// @Router /mission [get]
func (h *MissionHandler) GetAllMissions(c *gin.Context) {
	missions, err := h.MissionRepo.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve missions"})
		return
	}
	c.JSON(http.StatusOK, missions)
}

// GetMissionByID godoc
// @Summary Get a single mission by ID
// @Description Retrieves a mission by its ID
// @Tags missions
// @Produce json
// @Param id path int true "Mission ID"
// @Success 200 {object} model.Mission "Mission details"
// @Failure 400 {object} map[string]interface{} "Invalid mission ID"
// @Failure 404 {object} map[string]interface{} "Mission not found"
// @Router /mission/{id} [get]
func (h *MissionHandler) GetMissionByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid mission ID"})
		return
	}

	mission, err := h.MissionRepo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Mission not found"})
		return
	}

	c.JSON(http.StatusOK, mission)
}
