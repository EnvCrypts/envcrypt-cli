package app

import (
	"context"
	"crypto/rand"
	"errors"
	"log"

	"github.com/envcrypts/envcrypt-cli/internal/config"
	cryptoutils "github.com/envcrypts/envcrypt-cli/internal/crypto"
	"github.com/google/uuid"
	"github.com/spf13/viper"
)

func (app *App) CreateProject(ctx context.Context, projectName string) error {
	email := viper.GetString("user.email")
	if email == "" {
		return errors.New("no user email found")
	}

	userReq := config.UserKeyRequestBody{
		Email: email,
	}

	var userResp config.UserKeyResponseBody
	if err := app.HttpClient.Do(ctx, "POST", "/users/search", userReq, &userResp, false); err != nil {
		return err
	}

	prk := make([]byte, 32)
	if _, err := rand.Read(prk); err != nil {
		return err
	}

	wrappedKey, err := cryptoutils.WrapPRKForUser(prk, userResp.PublicKey)
	if err != nil {
		return err
	}

	projectReq := config.ProjectCreateRequest{
		Name:               projectName,
		UserId:             userResp.UserId,
		WrappedPRK:         wrappedKey.WrappedPRK,
		WrapNonce:          wrappedKey.WrapNonce,
		EphemeralPublicKey: wrappedKey.WrapEphemeralPub,
	}

	var projectResp config.ProjectCreateResponse
	if err := app.HttpClient.Do(ctx, "POST", "/projects/create", projectReq, &projectResp, true); err != nil {
		return err
	}

	return nil
}

func (app *App) ListProjects(ctx context.Context) (*config.ListProjectResponse, error) {
	userId := viper.GetString("user.id")
	if userId == "" {
		return nil, errors.New("no user id found")
	}

	uid, err := uuid.Parse(userId)
	if err != nil {
		return nil, err
	}

	projectsReq := config.ListProjectRequest{
		UserId: uid,
	}

	var projectsRes config.ListProjectResponse
	if err := app.HttpClient.Do(ctx, "POST", "/projects/list", projectsReq, &projectsRes, true); err != nil {
		log.Print(projectsRes)
		log.Print(err)
		return nil, err
	}

	return &projectsRes, nil
}

func (app *App) DeleteProject(ctx context.Context, projectName string) error {
	email, userId := viper.GetString("user.email"), viper.GetString("user.id")

	uid, err := uuid.Parse(userId)
	if err != nil {
		return err
	}

	if email == "" || uid == uuid.Nil {
		return errors.New("user not authenticated")
	}

	deleteReq := config.ProjectDeleteRequest{
		ProjectName: projectName,
		UserId:      uid,
	}
	var deleteResp config.ProjectDeleteResponse

	if err := app.HttpClient.Do(ctx, "POST", "/projects/delete", deleteReq, &deleteResp, true); err != nil {
		return err
	}

	return nil
}


func (app *App) RotatePRK(ctx context.Context, projectID string) (int32, error) {
	userEmail, userId := viper.GetString("user.email"), viper.GetString("user.id")
	if userEmail == "" || userId == "" {
		return 0, errors.New("user not authenticated")
	}

	uid, err := uuid.Parse(userId)
	if err != nil {
		return 0, err
	}

	projID, err := uuid.Parse(projectID)
	if err != nil {
		return 0, err
	}

	userPriv, err := cryptoutils.LoadPrivateKey(userEmail)
	if err != nil {
		return 0, errors.New("could not load private key")
	}

	// Step 1: Init
	initReq := config.RotateInitRequest{
		ProjectID: projID,
		UserID:    uid,
	}

	var initResp config.RotateInitResponse
	if err := app.HttpClient.Do(ctx, "POST", "/projects/rotate/init", initReq, &initResp, true); err != nil {
		return 0, err
	}

	// Step 2: Crypto (Local Only)
	var myWrappedPRK *config.WrappedKey
	for _, wp := range initResp.WrappedPRKs {
		if wp.UserID == uid {
			myWrappedPRK = &wp
			break
		}
	}

	if myWrappedPRK == nil {
		return 0, errors.New("user's wrapped PRK not found in init response")
	}

	wrappedKey := &cryptoutils.WrappedKey{
		WrappedPRK:       myWrappedPRK.WrappedPRK,
		WrapNonce:        myWrappedPRK.WrapNonce,
		WrapEphemeralPub: myWrappedPRK.EphemeralPublicKey,
	}

	oldPRK, err := cryptoutils.UnwrapPRK(wrappedKey, userPriv)
	if err != nil {
		return 0, errors.New("could not unwrap old PRK")
	}

	// Generate new PRK
	newPRK := make([]byte, 32)
	if _, err := rand.Read(newPRK); err != nil {
		return 0, errors.New("could not generate new PRK")
	}

	// Re-wrap PRK for each member
	newWrappedPRKs := make([]config.WrappedKey, len(initResp.MemberPublicKeys))
	for i, member := range initResp.MemberPublicKeys {
		wrapped, err := cryptoutils.WrapPRKForUser(newPRK, member.PublicKey)
		if err != nil {
			return 0, errors.New("could not wrap new PRK for member")
		}

		newWrappedPRKs[i] = config.WrappedKey{
			UserID:             member.UserID,
			WrappedPRK:         wrapped.WrappedPRK,
			WrapNonce:          wrapped.WrapNonce,
			EphemeralPublicKey: wrapped.WrapEphemeralPub,
		}
	}

	// Re-wrap DEKs
	newWrappedDEKs, err := cryptoutils.RewrapDEKs(oldPRK, newPRK, initResp.WrappedDEKs)
	if err != nil {
		return 0, errors.New("could not re-wrap DEKs")
	}

	// Step 3: Commit
	commitReq := config.RotateCommitRequest{
		ProjectID:          projID,
		UserID:             uid,
		ExpectedPRKVersion: initResp.PRKVersion,
		NewWrappedPRKs:     newWrappedPRKs,
		NewWrappedDEKs:     newWrappedDEKs,
	}

	var commitResp config.RotateCommitResponse
	if err := app.HttpClient.Do(ctx, "POST", "/projects/rotate/commit", commitReq, &commitResp, true); err != nil {
		return 0, err
	}

	return commitResp.NewPRKVersion, nil
}
