package application_test

import (
	"github.com/KadimovRus/calc_go/internal/application"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewAgent(t *testing.T) {
	t.Setenv("COMPUTING_POWER", "2")
	t.Setenv("ORCHESTRATOR_URL", "http://example.com")

	agent, err := application.NewAgent()
	assert.NoError(t, err)
	assert.Equal(t, 2, agent.ComputingPower)
	assert.Equal(t, "http://example.com", agent.OrchestratorURL)
}
