package models

import (
	"time"

	"github.com/lib/pq"
)

type Setting struct {
	OrgName          string         `db:"org_name" json:"orgName"`
	OrgURL           string         `db:"org_url" json:"orgUrl"`
	SupportEmail     string         `db:"support_email" json:"supportEmail"`
	SchemaMigrations pq.StringArray `db:"schema_migrations" json:"-"`
	AllowInvite      bool           `db:"allow_invite" json:"allowInvite"`
	AllowLogin       bool           `db:"allow_login" json:"allowLogin"`
	IsMaintenanceOn  bool           `db:"is_maintenance_on" json:"isMaintenanceOn"`
	CreatedAt        time.Time      `db:"created_at" json:"createdAt"`
	UpdatedAt        time.Time      `db:"updated_at" json:"-"`
}
