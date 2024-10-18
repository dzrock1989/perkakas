package elasticsearch

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
	"github.com/tigapilarmandiri/perkakas/configs"
)

func TestStoreLog(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	err := godotenv.Load("../../../.env")
	require.NoError(t, err)

	configs.LoadConfigs()

	type Data struct {
		Name string `json:"name"`
	}

	document := Data{
		Name: "go-percobaan",
	}

	err = StoreHistory(context.Background(), "kambin-history-test-1", CreateHistory(time.Now(), "testing_user", "testing_service", "test_table", uuid.NewString(), document))
	require.NoError(t, err)
}
