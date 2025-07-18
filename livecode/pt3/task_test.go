package pt3

import (
	"encoding/json"
	"fmt"
	"os"
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

	cwd, err := os.Getwd()
	require.NoError(t, err)
	fmt.Println(cwd)

	b, err := json.MarshalIndent(tasks, "", "\t")
	require.NoError(t, err)

	// записываем в файл
	err = os.WriteFile("tasks.json", b, 755) // 755 - разрешение
	require.NoError(t, err)                  // прочитали ошибку

	fromFile, err := os.ReadFile("tasks.json")
	require.NoError(t, err)

	var decoded []Task
	json.Unmarshal(fromFile, &decoded)
	require.NoError(t, err)

	require.Equal(t, int64(2), decoded[0].ID)
	fmt.Println(decoded)

	// var decoded []Task
	// json.Unmarshal(b, &decoded)
	// require.NoError(t, err)

	// fmt.Println(decoded)
}
