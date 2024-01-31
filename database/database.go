package database

import (
	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
	"os"
)

type Configuration struct {
	Host      string `json:"host"`
	PrivateIP string `json:"private_ip"`
	PublicIP  string `json:"public_ip"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Database  string `json:"database"`
}

// ---------------------------------------------------------------------------------------------------------
func Connect(ctx context.Context, secretPath string) (*sql.DB, error) {

	dbc, err := getCredentials(ctx, secretPath)

	if err != nil {
		log.Error("Database:Connection:getCredentials:", err)
		return nil, err
	}

	sqlPath := "/cloudsql"

	if len(os.Getenv("DEVELOPMENT")) > 0 {
		if len(os.Getenv("SQLPATH")) == 0 {
			return nil, fmt.Errorf("missing SQLPATH")
		}
		sqlPath = os.Getenv("SQLPATH")
	}

	dbURI := fmt.Sprintf("%s:%s@unix(%s/%s)/%s?autocommit=true&parseTime=true&timeout=5s", dbc.Username, dbc.Password, sqlPath, dbc.Host, dbc.Database)

	link, err := sql.Open("mysql", dbURI)
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %s", err.Error())
	}

	_, _ = link.Exec("SET time_zone = 'Europe/London'")

	return link, err

}

// getCredentials ------------------------------------------------------------------------------------------
func getCredentials(ctx context.Context, secretPath string) (Configuration, error) {

	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		log.Println(fmt.Errorf("failed to create secretmanager client: %v", err))
	}

	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: secretPath,
	}

	// Call the API.
	result, err := client.AccessSecretVersion(ctx, req)
	if err != nil {
		log.Println(fmt.Errorf("failed to access secret version: %v", err))
	}

	dbc := Configuration{}

	err = json.Unmarshal(result.Payload.Data, &dbc)

	return dbc, err

}
