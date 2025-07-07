package analyzer

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetLinks(c *gin.Context) {
	// Get pagination parameters from query string
	amountStr := c.DefaultQuery("amount", "10") // Default page size is 10
	pageStr := c.DefaultQuery("page", "1")     // Default page is 1

	// Convert string parameters to integers
	amount, err := strconv.Atoi(amountStr)
	if err != nil || amount < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid amount parameter"})
		return
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page parameter"})
		return
	}

	// Calculate offset
	offset := (page - 1) * amount

	// Get links from database using the new function
	links, err := GetLinksFromDB(amount, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch links"})
		return
	}

	// Add pagination metadata to the response
	c.JSON(http.StatusOK, gin.H{
		"links": links,
		"pagination": gin.H{
			"current_page": page,
			"page_size":    amount,
		},
	})
}