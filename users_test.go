package tenant

import (
	"os"
	"testing"
)

func TestGetUser(t *testing.T) {
	tn := Tenant{}
	tn.ClientID = os.Getenv("B2C_CLIENT_ID")
	tn.ClientSecret = os.Getenv("B2C_CLIENT_SECRET")
	tn.TenantDomain = os.Getenv("B2C_TENANT_DOMAIN")

	if err := tn.GetGraphAccessToken(); err != nil {
		t.Errorf("Error while obtaining access token: %s", err)
	}

	userObjectID := os.Getenv("B2C_TESTUSER")

	//_, err := tn.callGraphAPI("/users/"+userObjectID, "1.6", "GET", "")
	user, err := tn.GetUser(userObjectID)
	if err != nil {
		t.Errorf("Error while reading user: %s", err)
	}

	if user.ID != userObjectID {
		t.Errorf("Object ID of returned user is %s, should be: %s", user.ID, userObjectID)
	}
}
