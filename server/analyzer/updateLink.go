package analyzer

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func UpdateLink(c *gin.Context) {
	// Define request structure for the PUT request
	type UpdateRequest struct {
		ID     int    `json:"id" binding:"required"`
		Status string `json:"status" binding:"required"`
	}

	var req UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// Validate status value
	if req.Status != "stop" && req.Status != "pending" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Status must be either 'stop' or 'pending'"})
		return
	}

	// Fetch the existing record to preserve other fields
	var url string
	var title string
	var htmlVersion string
	var h1, h2, h3, h4, h5, h6, internal, external, inaccessible int
	var loginForm bool

	query := "SELECT url, title, html_version, h1, h2, h3, h4, h5, h6, internal_links, external_links, inaccessible_links, login_form_detected FROM results WHERE id = ?"
	err := DB.QueryRow(query, req.ID).Scan(&url, &title, &htmlVersion, &h1, &h2, &h3, &h4, &h5, &h6, &internal, &external, &inaccessible, &loginForm)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		return
	}

	// Create headings map for updateInDB function
	headings := map[string]int{
		"H1": h1,
		"H2": h2,
		"H3": h3,
		"H4": h4,
		"H5": h5,
		"H6": h6,
	}

	// Update the record with the new status
	now := time.Now()
	err = updateInDB(
		req.ID,
		url,
		req.Status,
		now,
		title,
		htmlVersion,
		headings,
		internal,
		external,
		inaccessible,
		loginForm,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update record"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Link status updated successfully",
		"id":      req.ID,
		"status":  req.Status,
	})
}