package pt3

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestJson(t *testing.T) {
	tasks := []Task{
		{
			ID:     2,
			Status: "started",
			Cfg: Config{
				StartAt:  "now",
				Timeline: "before 28.10",
			},
		},
		{
			ID:     3,
			Status: "not started",
		},
	}

	b, err := json.Marshal(tasks)
	require.NoError(t, err)

	fmt.Println(b)
}
