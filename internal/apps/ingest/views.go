package ingest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/santoshkpatro/unbit/internal/models"
	"github.com/santoshkpatro/unbit/internal/utils"
)

func (v *IngestContext) NewEvent(c echo.Context) error {
	var queue = "issue_events"
	var event models.Event
	if err := c.Bind(&event); err != nil {
		utils.RespondFail(c, http.StatusBadRequest, "Failed", nil)
	}

	for _, frame := range event.StackTrace {
		fmt.Printf("     ↳ %s (%s:%d)\n", frame.Function, frame.File, frame.Line)
	}

	token := c.Request().Header.Get("X-Unbit-Token")

	if token == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "missing X-Unbit-Token header",
		})
	}

	payload := models.Payload{
		DSNToken: token,
		Event:    event,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return utils.RespondFail(c, http.StatusInternalServerError, "Failed to marshal event", nil)
	}

	if err := v.Cache.LPush(c.Request().Context(), queue, data).Err(); err != nil {
		return utils.RespondFail(c, http.StatusInternalServerError, "Failed to queue event", nil)
	}

	return utils.RespondOK(c, nil, "Recived success")
}
