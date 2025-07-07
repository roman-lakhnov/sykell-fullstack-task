package analyzer

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AddRequest struct {
	URLs []string `json:"urls"`
}

func AddLinks(c *gin.Context) {
	var req AddRequest
	if err := c.BindJSON(&req); err != nil || len(req.URLs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON or missing URLs"})
		return
	}

	for _, link := range req.URLs {
		// Save to DB
		err := SaveToDB(link,"pending", nil, "", "", map[string]int{}, 0,  0, 0, false)
		if err != nil {
			fmt.Printf("Failed to save to DB for %s: %v\n", link, err)
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "added for analysis"})
}
