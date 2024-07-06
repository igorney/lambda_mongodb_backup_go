package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"
)

func handler(ctx context.Context) (string, error) {
	start := time.Now() // Marcar o início da execução

	// Valores hardcoded
	mongodbURI := "[MONGO_URI]"
	backupFolder := "/app/backups/"

	// Certifique-se de que a pasta de backup existe
	if _, err := os.Stat(backupFolder); os.IsNotExist(err) {
		if err := os.MkdirAll(backupFolder, 0755); err != nil {
			return "Failed to create backup folder", err
		}
	}
	folderCreationTime := time.Since(start)
	fmt.Printf("Tempo para criar/verificar pasta de backup: %s\n", folderCreationTime)

	// Construção dos comandos
	timestamp := time.Now().Format("20060102T150405")
	backupName := fmt.Sprintf("%s.dump.gz", timestamp)
	localBackupPath := fmt.Sprintf("%s%s", backupFolder, backupName)
	localLatestPath := fmt.Sprintf("%slatest.dump.gz", backupFolder)

	mongodumpCmd := []string{"mongodump", "--uri", mongodbURI, "--archive=" + localBackupPath, "--gzip"}

	// Execução dos comandos com timeout
	ctx, cancel := context.WithTimeout(ctx, 10*time.Minute) // Timeout de 10 minutos
	defer cancel()

	cmd := exec.CommandContext(ctx, mongodumpCmd[0], mongodumpCmd[1:]...)
	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	fmt.Println("Running mongodump command:", mongodumpCmd)
	mongodumpStart := time.Now()
	err := cmd.Run()
	mongodumpDuration := time.Since(mongodumpStart)
	fmt.Printf("Tempo de execução do mongodump: %s\n", mongodumpDuration)

	if ctx.Err() == context.DeadlineExceeded {
		fmt.Printf("Backup command timed out\n")
		return "Backup command timed out", ctx.Err()
	}
	if err != nil {
		fmt.Printf("Backup failed: %s\n", err)
		fmt.Printf("Stderr: %s\n", stderr.String())
		fmt.Printf("Stdout: %s\n", out.String())
		return fmt.Sprintf("Backup failed: %s\nStderr: %s\nStdout: %s", err, stderr.String(), out.String()), err
	}

	// Copiar o backup para o arquivo mais recente
	copyStart := time.Now()
	input, err := os.ReadFile(localBackupPath)
	if err != nil {
		return fmt.Sprintf("Failed to read backup file: %s", err), err
	}

	if err := os.WriteFile(localLatestPath, input, 0644); err != nil {
		return fmt.Sprintf("Failed to write latest backup file: %s", err), err
	}
	copyDuration := time.Since(copyStart)
	fmt.Printf("Tempo para copiar o arquivo de backup: %s\n", copyDuration)

	totalDuration := time.Since(start)
	fmt.Printf("Tempo total de execução: %s\n", totalDuration)

	return fmt.Sprintf("Backup succeeded: %s", out.String()), nil
}

func main() {
	// Execute localmente
	result, err := handler(context.Background())
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Result:", result)
	}
}
