package handlers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"main/internal/model"
	"main/internal/repositories"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

const catAPIURL = "https://api.thecatapi.com/v1/breeds"

type CatBreed struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type CatHandler struct {
	CatRepo *repositories.CatRepository
}

func NewCatHandler(catRepo *repositories.CatRepository) *CatHandler {
	return &CatHandler{CatRepo: catRepo}
}

// @Summary Create a new cat
// @Description Create a new cat with breed validation
// @Tags cats
// @Accept json
// @Produce json
// @Param cat body model.SpyCat true "Cat data"
// @Success 201 {object} model.SpyCat
// @Failure 400 {object} map[string]interface{} "Invalid request body or breed"
// @Failure 500 {object} map[string]interface{} "Failed to create cat"
// @Router /cat [post]
func (h *CatHandler) CreateCat(c *gin.Context) {
	var cat model.SpyCat
	if err := c.ShouldBindJSON(&cat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := validateBreed(cat.Breed); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.CatRepo.Create(&cat)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create cat"})
		return
	}

	c.JSON(http.StatusCreated, cat)
}

// @Summary Get all cats
// @Description Get a list of all cats
// @Tags cats
// @Produce json
// @Success 200 {array} model.SpyCat
// @Failure 500 {object} map[string]interface{}
// @Router /cat [get]
func (h *CatHandler) GetAllCats(c *gin.Context) {
	cats, err := h.CatRepo.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve cats"})
		return
	}
	c.JSON(http.StatusOK, cats)
}

// GetCatByID godoc
// @Summary Get a single spy cat by ID
// @Description Get a single spy cat by its ID
// @Tags cats
// @Produce json
// @Param id path int true "Cat ID"
// @Success 200 {object} model.SpyCat
// @Failure 400 {object} map[string]interface{} "Invalid cat ID"
// @Failure 404 {object} map[string]interface{} "Cat not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /cat/{id} [get]
func (h *CatHandler) GetCatByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cat ID"})
		return
	}

	cat, err := h.CatRepo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cat not found"})
		return
	}

	c.JSON(http.StatusOK, cat)
}

// UpdateCatSalary godoc
// @Summary Update cat's salary
// @Description Update the salary of a spy cat by its ID
// @Tags cats
// @Produce json
// @Param id path int true "Cat ID"
// @Param salary body model.SalaryUpdate true "Salary data"
// @Success 200 {object} map[string]interface{} "Salary updated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid cat ID or request body"
// @Failure 500 {object} map[string]interface{} "Failed to update salary"
// @Router /cat/{id}/salary [put]
func (h *CatHandler) UpdateCatSalary(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cat ID"})
		return
	}

	var updateData model.SalaryUpdate
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	err = h.CatRepo.UpdateSalary(id, updateData.Salary)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update salary"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Salary updated successfully"})
}

// DeleteCat godoc
// @Summary Delete a spy cat
// @Description Delete a spy cat by its ID
// @Tags cats
// @Produce json
// @Param id path int true "Cat ID"
// @Success 200 {object} map[string]interface{} "Cat deleted successfully"
// @Failure 400 {object} map[string]interface{} "Invalid cat ID"
// @Failure 500 {object} map[string]interface{} "Failed to delete cat"
// @Router /cat/{id} [delete]
func (h *CatHandler) DeleteCat(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cat ID"})
		return
	}

	err = h.CatRepo.Delete(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete cat"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cat deleted successfully"})
}

// ValidateBreed checks if the given breed is valid by comparing it to the list of breeds from TheCatAPI
func validateBreed(breed string) error {
	// Make the HTTP request to TheCatAPI
	resp, err := http.Get(catAPIURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("failed to retrieve cat breeds")
	}

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var breeds []CatBreed
	err = json.Unmarshal(body, &breeds)
	if err != nil {
		return err
	}

	// Check if the breed exists in the list of valid breeds
	for _, b := range breeds {
		if b.Name == breed {
			return nil // breed is valid
		}
	}

	return errors.New("invalid breed")
}
