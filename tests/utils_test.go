package tests

import (
	"os"
	"testing"

	"github.com/google/uuid"

	gcs "github.com/aos-dev/go-service-gcs"
	ps "github.com/aos-dev/go-storage/v3/pairs"
	"github.com/aos-dev/go-storage/v3/types"
)

func setupTest(t *testing.T) types.Storager {
	t.Log("Setup test for gcs")

	store, err := gcs.NewStorager(
		ps.WithCredential(os.Getenv("STORAGE_GCS_CREDENTIAL")),
		ps.WithName(os.Getenv("STORAGE_GCS_NAME")),
		ps.WithWorkDir("/"+uuid.New().String()+"/"),
		gcs.WithProjectID(os.Getenv("STORAGE_GCS_PROJECT_ID")),
	)
	if err != nil {
		t.Errorf("new storager: %v", err)
	}
	return store
}
