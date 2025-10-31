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
	fingerprint := ComputeFingerprint(event)

	var projectId string
	err := db.Get(&projectId, "SELECT id FROM projects WHERE dsn_token = $1", dsnToken)
	if err != nil {
		log.Println("❌ Error fetching project for DSN token:", err)
		return
	}

	var issue models.Issue
	err = db.Get(&issue, "SELECT * FROM issues WHERE project_id = $1 AND fingerprint = $2", projectId, fingerprint)
	if err != nil {
		// Issue does not exist, create a new one
		summary := event.Message
		issueID := utils.GenerateID("iss")
		lastSeenAt := event.Timestamp

		_, err = db.Exec(`
			INSERT INTO issues (id, project_id, summary, fingerprint, last_seen_at, event_count, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, 1, NOW(), NOW())
		`, issueID, projectId, summary, fingerprint, lastSeenAt)

		if err != nil {
			log.Println("❌ Error creating new issue:", err)
			return
		}

		// Now let's insert event
		_, err = db.Exec(`
			INSERT INTO events (id, issue_id, message, level, timestamp, stacktrace, created_at)
			VALUES ($1, $2, $3, $4, $5, $6, NOW())
		`, utils.GenerateID("evt"), issueID, event.Message, event.Level, event.Timestamp, StackTraceToJSON(event.StackTrace))
		if err != nil {
			log.Println("❌ Error inserting event for new issue:", err)
			return
		}

		log.Println("✅ Created new issue with ID:", issueID)
	} else {
		issueID := issue.ID
		// Issue exists, update it
		_, err = db.Exec(`
			UPDATE issues
			SET last_seen_at = $1, event_count = event_count + 1, updated_at = NOW()
			WHERE id = $2
		`, event.Timestamp, issueID)

		if err != nil {
			log.Println("❌ Error updating issue:", err)
			return
		}

		// Insert event
		_, err = db.Exec(`
			INSERT INTO events (id, issue_id, message, level, timestamp, stacktrace, created_at)
			VALUES ($1, $2, $3, $4, $5, $6, NOW())
		`, utils.GenerateID("evt"), issueID, event.Message, event.Level, event.Timestamp, StackTraceToJSON(event.StackTrace))

		if err != nil {
			log.Println("❌ Error inserting event for existing issue:", err)
			return
		}

		log.Println("✅ Updated existing issue with ID:", issueID)
	}

	fmt.Printf("Processing event for project %s, issue %s\n", projectId, issue.ID)
}
