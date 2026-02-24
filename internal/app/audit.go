package app

import (
	"context"
	"errors"

	"github.com/envcrypts/envcrypt-cli/internal/config"
	"github.com/google/uuid"
	"github.com/spf13/viper"
)

type AuditOptions struct {
	Limit      int
	Offset     int
	ActorEmail string
	Action     string
	Status     string
	From       string
	To         string
	JSON       bool
}

func (app *App) GetProjectAuditLogs(ctx context.Context, projectName string, opts AuditOptions) (*config.ProjectAuditResponse, error) {
	userId := viper.GetString("user.id")
	if userId == "" {
		return nil, errors.New("missing user id")
	}

	uid, err := uuid.Parse(userId)
	if err != nil {
		return nil, err
	}

	projectRequest := config.GetMemberProjectRequest{
		ProjectName: projectName,
		UserId:      uid,
	}

	var projectResponse config.GetMemberProjectResponse
	err = app.HttpClient.Do(ctx, "POST", "/projects/get", projectRequest, &projectResponse, true)
	if err != nil {
		return nil, err
	}

	req := config.ProjectAuditRequest{
		ProjectID:  projectResponse.ProjectId.String(),
		Limit:      opts.Limit,
		Offset:     opts.Offset,
		ActorEmail: opts.ActorEmail,
		Action:     opts.Action,
		Status:     opts.Status,
		From:       opts.From,
		To:         opts.To,
	}

	var resp config.ProjectAuditResponse
	err = app.HttpClient.Do(ctx, "POST", "/projects/audit", req, &resp, true)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}
