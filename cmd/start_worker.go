package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/santoshkpatro/unbit/internal/config"
	"github.com/spf13/cobra"
)

var startWorkerCmd = &cobra.Command{
	Use:   "start_worker",
	Short: "Start Worker",
	RunE: func(cmd *cobra.Command, args []string) error {
		return startWorker()
	},
}

func startWorker() error {
	const (
		queueName      = "unbit_jobs"
		maxConcurrency = 5
	)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cacheClient, _ := config.NewRedisConnection(ctx)
	if err := cacheClient.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("redis ping failed: %w", err)
	}
	fmt.Println("Connected to Redis, waiting for jobs...")
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		fmt.Println("Shutting down worker...")
		cancel()
	}()

	sem := make(chan struct{}, maxConcurrency)
	var wg sync.WaitGroup

	for {
		select {
		case <-ctx.Done():
			wg.Wait()
			fmt.Println("All jobs complete. Exiting.")
			return nil
		default:
			// Blocking pop (waits until a job appears)
			res, err := cacheClient.BLPop(ctx, 0*time.Second, queueName).Result()
			if err != nil {
				if ctx.Err() != nil {
					return nil // Context cancelled, exit
				}
				fmt.Println("Error reading from Redis:", err)
				continue
			}

			if len(res) < 2 {
				continue
			}
			payload := res[1]

			// Limit concurrency
			sem <- struct{}{}
			wg.Add(1)
			go func(p string) {
				defer wg.Done()
				defer func() { <-sem }()
				if err := processJob(p); err != nil {
					fmt.Println("Job failed:", err)
				} else {
					fmt.Println("Job done!")
				}
			}(payload)
		}
	}
}

func processJob(payload string) error {
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(payload), &data); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}

	fmt.Printf("Processing job: %+v\n", data)

	// Simulate work
	time.Sleep(2 * time.Second)
	return nil
}
