package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"github.com/santoshkpatro/unbit/internal/models"
	"github.com/santoshkpatro/unbit/internal/utils"
)

func StartWorker(cache *redis.Client, db *sqlx.DB, queue string) {
	ctx := context.Background()

	log.Println("🚀 Worker started, listening on queue:", queue)

	for {
		result, err := cache.BLPop(ctx, 0, queue).Result()
		if err != nil {
			log.Println("❌ Error fetching job from queue:", err)
			continue
		}

		payloadJSON := result[1]
		var payload models.Payload
		if err := json.Unmarshal([]byte(payloadJSON), &payload); err != nil {
			log.Println("❌ Error unmarshaling event:", err)
			continue
		}

		go handleEvent(db, payload.DSNToken, payload.Event)
	}
}

func handleEvent(db *sqlx.DB, dsnToken string, event models.Event) {
	tx, err := db.Beginx()
	if err != nil {
		log.Println("begin tx:", err)
		return
	}
	defer tx.Rollback() // no-op if Commit succeeds

	fingerprint := ComputeFingerprint(event.Properties)

	var projectId string
	if err = tx.Get(&projectId, `SELECT id FROM projects WHERE dsn_token = $1`, dsnToken); err != nil {
		log.Println("❌ project lookup:", err)
		return
	}

	// unique index on (project_id, fingerprint) recommended
	var issueId string
	newIssueId := utils.GenerateID("isu")
	if err = tx.Get(&issueId, `
		INSERT INTO issues (id, project_id, fingerprint, status, event_count)
		VALUES ($1, $2, $3, 'unresolved', 0)
		ON CONFLICT (project_id, fingerprint)
		DO UPDATE SET updated_at = NOW()
		RETURNING id
	`, newIssueId, projectId, fingerprint); err != nil {
		log.Println("❌ upsert issue:", err)
		return
	}

	eventId := utils.GenerateID("evt")
	if _, err = tx.Exec(`
		INSERT INTO events (id, issue_id, timestamp, properties, project_id, event_type)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, eventId, issueId, event.Timestamp, PropertiesToJSON(event.Properties), projectId, `issues`); err != nil {
		log.Println("❌ insert event:", err)
		return
	}

	if _, err = tx.Exec(`
		UPDATE issues
		SET event_count = event_count + 1, updated_at = NOW()
		WHERE id = $1
	`, issueId); err != nil {
		log.Println("❌ bump count:", err)
		return
	}

	if _, err = tx.Exec(`
		UPDATE projects
		SET total_events = total_events + 1
		WHERE id = $1
	`, projectId); err != nil {
		log.Println("❌ bump project event count:", err)
		return
	}

	if err = tx.Commit(); err != nil {
		log.Println("❌ commit:", err)
		return
	}

	fmt.Printf("Processing event for project %s, issue %s\n", projectId, eventId)
}

func PropertiesToJSON(props models.Properties) string {
	jsonData, err := json.Marshal(props)
	if err != nil {
		return "{}"
	}
	return string(jsonData)
}
