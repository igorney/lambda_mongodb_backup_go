package main

import (
	"context"
	"testing"
)

func TestHandler(t *testing.T) {
	t.Skip("Skipping test because env is not set")
	// Chame sua função handler aqui e verifique os resultados
	result, err := handler(context.Background())
	if err != nil {
		t.Errorf("handler returned an error: %v", err)
	}
	if result != "Backup completed successfully and uploaded to S3" {
		t.Errorf("handler returned unexpected result: %v", result)
	}
}
