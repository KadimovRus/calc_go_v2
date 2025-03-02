package application

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/KadimovRus/calc_go/internal/pkg/calculation"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

const (
	retryDelay = 2 * time.Second
	pollDelay  = 1 * time.Second
)

type Agent struct {
	ComputingPower  int
	OrchestratorURL string
}

type AgentConfig struct {
	ComputingPower  int
	OrchestratorURL string
}

func NewAgent(opts ...AgentConfig) (*Agent, error) {

	config := AgentConfig{
		ComputingPower:  1,
		OrchestratorURL: "http://localhost:8080",
	}

	if len(opts) > 0 {
		config = opts[0]
	}

	if cpStr := os.Getenv("COMPUTING_POWER"); cpStr != "" {
		cp, err := strconv.Atoi(cpStr)
		if err != nil {
			return nil, fmt.Errorf("неверный формат COMPUTING_POWER: %v", err)
		}
		if cp < 1 {
			return nil, fmt.Errorf("COMPUTING_POWER должно быть больше 0")
		}
		config.ComputingPower = cp
	}

	if url := os.Getenv("ORCHESTRATOR_URL"); url != "" {
		config.OrchestratorURL = url
	}

	return &Agent{
		ComputingPower:  config.ComputingPower,
		OrchestratorURL: config.OrchestratorURL,
	}, nil
}

func (a *Agent) Run(ctx context.Context) error {
	var wg sync.WaitGroup

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	for i := 0; i < a.ComputingPower; i++ {
		wg.Add(1)

		go func(id int) {
			defer wg.Done()

			a.worker(id)
		}(i)
	}

	go func() {
		wg.Wait()
		cancel()
	}()

	select {
	case <-ctx.Done():
		log.Printf("Agent shutting down...")
		return ctx.Err()
	}
}

func (a *Agent) worker(id int) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	for {
		resp, err := client.Get(a.OrchestratorURL + "/internal/task")
		if err != nil {
			log.Printf("Worker %d: error getting task: %v", id, err)
			time.Sleep(retryDelay)
			continue
		}

		defer resp.Body.Close()

		if resp.StatusCode == http.StatusNotFound {
			time.Sleep(pollDelay)
			continue
		}

		var taskResp struct {
			Task struct {
				ID            string  `json:"id"`
				Arg1          float64 `json:"arg1"`
				Arg2          float64 `json:"arg2"`
				Operation     string  `json:"operation"`
				OperationTime int     `json:"operation_time"`
			} `json:"task"`
		}

		err = json.NewDecoder(resp.Body).Decode(&taskResp)
		if err != nil {
			log.Printf("Worker %d: error decoding task: %v", id, err)
			time.Sleep(pollDelay)
			continue
		}

		task := taskResp.Task
		log.Printf("Worker %d: received task %s: %f %s %f, simulating %d ms",
			id, task.ID, task.Arg1, task.Operation, task.Arg2, task.OperationTime)

		time.Sleep(time.Duration(task.OperationTime) * time.Millisecond)

		result, err := calculation.Compute(task.Operation, task.Arg1, task.Arg2)
		if err != nil {
			log.Printf("Worker %d: error computing task %s: %v", id, task.ID, err)
			continue
		}

		err = a.postTaskResult(client, task.ID, result)
		if err != nil {
			log.Printf("Worker %d: error posting result for task %s: %v", id, task.ID, err)
			continue
		}

		log.Printf("Worker %d: successfully completed task %s with result %f",
			id, task.ID, result)
	}
}

func (a *Agent) postTaskResult(client *http.Client, taskID string, result float64) error {
	payload := map[string]interface{}{
		"id":     taskID,
		"result": result,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal result: %w", err)
	}

	resp, err := client.Post(a.OrchestratorURL+"/internal/task",
		"application/json", bytes.NewReader(payloadBytes))
	if err != nil {
		return fmt.Errorf("failed to post result: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("error response: %s", string(body))
	}

	return nil
}
