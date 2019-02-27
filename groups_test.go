package tenant

import (
	"fmt"
	"os"
	"testing"
)

func TestGetGroupMembers(t *testing.T) {
	tn := Tenant{}
	tn.ClientID = os.Getenv("B2C_CLIENT_ID")
	tn.ClientSecret = os.Getenv("B2C_CLIENT_SECRET")
	tn.TenantDomain = os.Getenv("B2C_TENANT_DOMAIN")

	if err := tn.GetGraphAccessToken(); err != nil {
		t.Errorf("Error while obtaining access token: %s", err)
	}

	groupID := os.Getenv("B2C_TESTGROUP")

	//_, err := tn.callGraphAPI("/users/"+userObjectID, "1.6", "GET", "")
	members, err := tn.GetGroupMembers(groupID)
	if err != nil {
		t.Errorf("Error while reading member list: %s", err)
	}

	if len(members) == 0 {
		t.Errorf("Error, group is not containing any members")
	}

	fmt.Println(members)
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
