package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"
)

func handler(ctx context.Context) (string, error) {
	// Valores hardcoded
	mongodbHost := "mongodb"
	mongodbPort := "27017"
	mongodbUser := "root"
	mongodbPass := "example"
	mongodbAuthDB := "admin"
	mongodbDB := "testdb"
	backupFolder := "/app/backups/"

	// Certifique-se de que a pasta de backup existe
	if _, err := os.Stat(backupFolder); os.IsNotExist(err) {
		if err := os.MkdirAll(backupFolder, 0755); err != nil {
			return "Failed to create backup folder", err
		}
	}

	// Construção dos comandos
	timestamp := time.Now().Format("20060102T150405")
	backupName := fmt.Sprintf("%s.dump.gz", timestamp)
	localBackupPath := fmt.Sprintf("%s%s", backupFolder, backupName)
	localLatestPath := fmt.Sprintf("%slatest.dump.gz", backupFolder)

	mongodumpCmd := []string{"mongodump", "--host", mongodbHost, "--port", mongodbPort, "--archive=" + localBackupPath, "--gzip", "--authenticationDatabase=" + mongodbAuthDB}
	if mongodbUser != "" {
		mongodumpCmd = append(mongodumpCmd, "--username", mongodbUser)
	}
	if mongodbPass != "" {
		mongodumpCmd = append(mongodumpCmd, "--password", mongodbPass)
	}
	if mongodbDB != "" {
		mongodumpCmd = append(mongodumpCmd, "--db", mongodbDB)
	}

	// Execução dos comandos
	if err := exec.Command(mongodumpCmd[0], mongodumpCmd[1:]...).Run(); err != nil {
		return "Backup failed", err
	}

	// Copiar o backup para o arquivo mais recente
	input, err := os.ReadFile(localBackupPath)
	if err != nil {
		return "Failed to read backup file", err
	}

	if err := os.WriteFile(localLatestPath, input, 0644); err != nil {
		return "Failed to write latest backup file", err
	}

	return "Backup succeeded", nil
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
