package configuration

import (
	"context"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
	"path/filepath"
)

func SetupFirebase() (*firebase.App, *auth.Client) {
	serviceAccountKeyFilePath, err := filepath.Abs("./serviceAccountKey.json")
	if err != nil {
		panic("Unable to load serviceAccountKeys.json file")
	}
	opt := option.WithCredentialsFile(serviceAccountKeyFilePath)
	// Firebase admin SDK initialization
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		panic("Firebase load error")
	}
	// Firebase Auth
	authenticator, err := app.Auth(context.Background())
	if err != nil {
		panic("Firebase load error")
	}
	return app, authenticator
}
