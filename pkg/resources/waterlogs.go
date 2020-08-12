package resources

import (
	"net/http"
	"strings"
	"time"

	"github.com/dgravesa/fountain/pkg/fountain"

	"github.com/dgravesa/fountain/pkg/data"
	"github.com/gin-gonic/gin"
)

// WaterLogsResource contains handlers for the waterlogs endpoint
type WaterLogsResource struct {
	reservoir data.Reservoir
}

// NewWaterLogsResource instantiates a new waterlogs resource
func NewWaterLogsResource(r data.Reservoir) *WaterLogsResource {
	wl := new(WaterLogsResource)
	wl.reservoir = r
	return wl
}

// GetWls gets waterlogs for a user
func (wlr *WaterLogsResource) GetWls(c *gin.Context) {
	userID := c.Param("user")

	// retrieve logs for the user from the reservoir
	userlogs, err := wlr.reservoir.UserWls(userID)
	if err != nil {
		c.Error(err)
		c.Status(http.StatusInternalServerError)
	} else {
		c.JSON(http.StatusOK, userlogs)
	}
}

type newWl struct {
	Amount *float64 `form:"amount" binding:"required,numeric,gt=0"`
}

// PostWl posts a new waterlog for a user
func (wlr *WaterLogsResource) PostWl(c *gin.Context) {
	userID := c.Param("user")

	var nwl newWl
	if err := c.Bind(&nwl); err != nil {
		c.JSON(http.StatusBadRequest, strings.Split(err.Error(), "\n"))
		return
	}

	// create the log
	wl := fountain.WaterLog{
		Time:   time.Now(),
		Amount: *nwl.Amount,
	}

	// push the log into the reservoir
	if err := wlr.reservoir.WriteWl(userID, &wl); err != nil {
		c.Error(err)
		c.Status(http.StatusInternalServerError)
	} else {
		c.JSON(http.StatusCreated, wl)
	}
}
