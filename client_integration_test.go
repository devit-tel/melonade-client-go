// +build integration

package melonade_client_go

import (
	"context"
	"testing"

	"github.com/devit-tel/goxid"
	"github.com/stretchr/testify/require"
)

func TestClient_StartWorkflow(t *testing.T) {
	const processManagerEndpoint = "http://localhost:8083/api/process-manager"

	t.Run("Start workflow", func(t *testing.T) {
		client := New(processManagerEndpoint)
		uuid := goxid.New()

		t.Run("success", func(t *testing.T) {
			resp, err := client.StartWorkflow(context.Background(), "simple", "1", uuid.Gen(), nil)

			require.NoError(t, err)
			require.True(t, resp.Success)
			require.NotEmpty(t, resp.Data.CreateTime)
			require.Equal(t, resp.Data.Status, WORKFLOW_STATE_RUNNING)
			require.Nil(t, resp.Error)
		})

		t.Run("fail workflow not found", func(t *testing.T) {
			resp, err := client.StartWorkflow(context.Background(), "not_found_workflow", "1", "unknown_id", nil)

			require.NoError(t, err)
			require.False(t, resp.Success)
			require.Empty(t, resp.Data)
			require.Equal(t, resp.Error.Message, "Workflow not found")
			require.NotEmpty(t, resp.Error.Stack)
		})
	})
}
