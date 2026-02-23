package app

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/envcrypts/envcrypt-cli/internal/config"
	"github.com/google/uuid"
	"github.com/spf13/viper"
)

func (a *App) ExportSnapshot(ctx context.Context, projectName, filename string) (string, error) {
	userIDStr := viper.GetString("user.id")
	if userIDStr == "" {
		return "", errors.New("user not authenticated")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return "", fmt.Errorf("invalid user ID: %w", err)
	}

	req := config.SnapshotExportRequest{
		ProjectName: projectName,
		UserID:      userID,
	}

	var resp config.SnapshotExportResponse
	err = a.HttpClient.Do(ctx, "POST", "/projects/snapshot/export", req, &resp, true)
	if err != nil {
		return "", err
	}

	data, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal snapshot: %w", err)
	}

	err = os.WriteFile(filename, data, 0600)
	if err != nil {
		return "", fmt.Errorf("failed to write snapshot file: %w", err)
	}

	absPath, err := filepath.Abs(filename)
	if err != nil {
		absPath = filename
	}

	return absPath, nil
}

func (a *App) ImportSnapshot(ctx context.Context, newProjectName, filename string) (string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("failed to read snapshot file: %w", err)
	}

	var file config.SnapshotExportResponse
	err = json.Unmarshal(data, &file)
	if err != nil {
		return "", fmt.Errorf("failed to parse snapshot file: %w", err)
	}

	userIDStr := viper.GetString("user.id")
	if userIDStr == "" {
		return "", errors.New("user not authenticated")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return "", fmt.Errorf("invalid user ID: %w", err)
	}

	req := config.SnapshotImportRequest{
		NewProjectName: newProjectName,
		UserID:         userID,
		Snapshot:       file.Snapshot,
		Checksum:       file.Checksum,
	}

	var resp config.SnapshotImportResponse
	err = a.HttpClient.Do(ctx, "POST", "/projects/snapshot/import", req, &resp, true)
	if err != nil {
		return "", err
	}

	return resp.NewProjectID.String(), nil
}
