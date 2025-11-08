package setting

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/santoshkpatro/unbit/internal/utils"
)

func setNested(out map[string]any, dotted string, val any) {
	parts := strings.Split(dotted, ".")
	m := out
	for i, p := range parts {
		if i == len(parts)-1 {
			// leaf
			m[p] = val
			return
		}
		// ensure intermediate map exists
		if next, ok := m[p]; ok {
			if nm, ok := next.(map[string]any); ok {
				m = nm
				continue
			}
			// collision with non-map; overwrite with a new map
		}
		nm := map[string]any{}
		m[p] = nm
		m = nm
	}
}

func (v *SettingContext) SettingMeta(c echo.Context) error {
	keys := []string{
		"org.siteName",
		"org.rootUrl",
		"org.supportEmail",
		"system.maintenanceMode",
		"system.maintenanceMessage",
		"ui.theme",
		"ui.language",
	}

	query, args, err := sqlx.In(`
		SELECT key, value
		FROM settings
		WHERE key IN (?)
	`, keys)
	if err != nil {
		return utils.RespondFail(c, http.StatusBadRequest, "No setting", nil)
	}

	query = v.DB.Rebind(query)

	rows, err := v.DB.Queryx(query, args...)
	if err != nil {
		return utils.RespondFail(c, http.StatusBadRequest, "No setting", nil)
	}
	defer rows.Close()

	out := make(map[string]any)

	for rows.Next() {
		var key string
		var raw []byte

		if err := rows.Scan(&key, &raw); err != nil {
			return err
		}

		var val any
		if err := json.Unmarshal(raw, &val); err != nil {
			return err
		}

		setNested(out, key, val)
	}

	return utils.RespondOK(c, out, "")

}
