package tenant

import (
	"os"
	"strings"
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

func TestGetMemberGroupIDs(t *testing.T) {
	tn := Tenant{}
	tn.ClientID = os.Getenv("B2C_CLIENT_ID")
	tn.ClientSecret = os.Getenv("B2C_CLIENT_SECRET")
	tn.TenantDomain = os.Getenv("B2C_TENANT_DOMAIN")

	if err := tn.GetAccessToken(); err != nil {
		t.Errorf("Error while obtaining access token: %s", err)
	}

	userObjectID := os.Getenv("B2C_TESTUSER")

	groups, err := tn.GetMemberGroupIDs(userObjectID)
	if err != nil {
		t.Errorf("Error while obtaining member group IDs: %s", err)
	}

	expectedGroups := os.Getenv("B2C_USERGROUPS")
	gotGroups := strings.Join(groups, " ")

	if expectedGroups != gotGroups {
		t.Errorf("Expected groups %q, got %q", expectedGroups, gotGroups)
	}
}

func TestAddGroupMember(t *testing.T) {
	tn := Tenant{}
	tn.ClientID = os.Getenv("B2C_CLIENT_ID")
	tn.ClientSecret = os.Getenv("B2C_CLIENT_SECRET")
	tn.TenantDomain = os.Getenv("B2C_TENANT_DOMAIN")

	if err := tn.GetAccessToken(); err != nil {
		t.Errorf("Error while obtaining access token: %s", err)
	}

	userEmail := os.Getenv("B2C_TESTMAIL")
	aadGroup := os.Getenv("B2C_TESTGROUP")

	if err := tn.AddGroupMember(aadGroup, userEmail); err != nil {
		t.Errorf("Error while adding user %q to group %q: %s", userEmail, aadGroup, err)
	}
}

func TestDeleteGroupMember(t *testing.T) {
	tn := Tenant{}
	tn.ClientID = os.Getenv("B2C_CLIENT_ID")
	tn.ClientSecret = os.Getenv("B2C_CLIENT_SECRET")
	tn.TenantDomain = os.Getenv("B2C_TENANT_DOMAIN")

	if err := tn.GetAccessToken(); err != nil {
		t.Errorf("Error while obtaining access token: %s", err)
	}

	userEmail := os.Getenv("B2C_TESTMAIL")
	aadGroup := os.Getenv("B2C_TESTGROUP")

	if err := tn.DeleteGroupMember(aadGroup, userEmail); err != nil {
		t.Errorf("Error while deleting user %q from group %q: %s", userEmail, aadGroup, err)
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
