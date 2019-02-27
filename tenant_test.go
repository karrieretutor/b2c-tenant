package tenant

import (
	"os"
	"testing"
)

func TestGetAccessToken(t *testing.T) {
	tn := Tenant{}
	tn.ClientID = os.Getenv("B2C_CLIENT_ID")
	tn.ClientSecret = os.Getenv("B2C_CLIENT_SECRET")
	tn.TenantDomain = os.Getenv("B2C_TENANT_DOMAIN")

	if err := tn.GetAccessToken(); err != nil {
		t.Errorf("Error while obtaining access token: %s", err)
	}
}

func TestGetGraphAccessToken(t *testing.T) {
	tn := Tenant{}
	tn.ClientID = os.Getenv("B2C_CLIENT_ID")
	tn.ClientSecret = os.Getenv("B2C_CLIENT_SECRET")
	tn.TenantDomain = os.Getenv("B2C_TENANT_DOMAIN")

	if err := tn.GetGraphAccessToken(); err != nil {
		t.Errorf("Error while obtaining access token: %s", err)
	}
}

func TestCallNewGraphAPI(t *testing.T) {
	tn := Tenant{}
	tn.ClientID = os.Getenv("B2C_CLIENT_ID")
	tn.ClientSecret = os.Getenv("B2C_CLIENT_SECRET")
	tn.TenantDomain = os.Getenv("B2C_TENANT_DOMAIN")

	if err := tn.GetGraphAccessToken(); err != nil {
		t.Errorf("Error while obtaining access token: %s", err)
	}

	userObjectID := os.Getenv("B2C_TESTUSER")

	_, err := tn.callNewGraphAPI("/users/"+userObjectID, "GET", "")
	if err != nil {
		t.Errorf("Error while reading user: %s", err)
	}
}

func TestCallGraphAPI(t *testing.T) {
	tn := Tenant{}
	tn.ClientID = os.Getenv("B2C_CLIENT_ID")
	tn.ClientSecret = os.Getenv("B2C_CLIENT_SECRET")
	tn.TenantDomain = os.Getenv("B2C_TENANT_DOMAIN")

	if err := tn.GetAccessToken(); err != nil {
		t.Errorf("Error while obtaining access token: %s", err)
	}

	userObjectID := os.Getenv("B2C_TESTUSER")

	_, err := tn.callGraphAPI("/users/"+userObjectID, "1.6", "GET", "")
	if err != nil {
		t.Errorf("Error while reading user: %s", err)
	}
}
