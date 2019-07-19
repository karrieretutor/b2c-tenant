package tenant

import (
	"fmt"
	"os"
	"strings"
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

func TestSearchUser(t *testing.T) {
	tn := Tenant{}
	tn.ClientID = os.Getenv("B2C_CLIENT_ID")
	tn.ClientSecret = os.Getenv("B2C_CLIENT_SECRET")
	tn.TenantDomain = os.Getenv("B2C_TENANT_DOMAIN")

	if err := tn.GetGraphAccessToken(); err != nil {
		t.Errorf("Error while obtaining access token: %s", err)
	}

	userEmail := os.Getenv("B2C_TESTMAIL")

	users, err := tn.SearchUser(userEmail)
	if err != nil {
		t.Errorf("Error while searching for user %q: %s", userEmail, err)
	}

	// Test if each returned user actually contains the searched for address
	for _, user := range users {
		if !strings.Contains(user.EmailAddresses[0], userEmail) {
			t.Errorf("Got user %s, expected %s", user.ID, os.Getenv("B2C_TESTUSER"))
		}
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

func TestAddDeleteUser(t *testing.T) {
	tn := Tenant{}
	tn.ClientID = os.Getenv("B2C_CLIENT_ID")
	tn.ClientSecret = os.Getenv("B2C_CLIENT_SECRET")
	tn.TenantDomain = os.Getenv("B2C_TENANT_DOMAIN")

	if err := tn.GetAccessToken(); err != nil {
		t.Errorf("Error while obtaining access token: %s", err)
	}

	testmail := os.Getenv("B2C_TESTMAIL")
	testmail = "b2c-cli.test." + testmail

	user, err := tn.AddUser(testmail, "b2c-cli-testuser", "laskdAföla!!ghhgaöaöa")
	if err != nil {
		t.Fatalf("Error while adding User: %s", err)
	}

	expectedMail := testmail
	if len(user.EmailAddresses) <= 0 {
		t.Logf("Error while adding user returned user has no mail addresses set, user looks like:\n%v",
			*user)
		t.Fail()
	} else if user.EmailAddresses[0] != expectedMail {
		t.Logf("Error: expected returned user to have mail %s but has %s instead, user looks like\n%v",
			expectedMail, user.EmailAddresses[0], *user)
		t.Fail()
	}
	if t.Failed() {
		//cleanup
		t.Log("Deleting user of failed test....")
		if err = tn.DeleteUser(user.ObjectID); err != nil {
			t.Errorf("Ran into error while creating user, tried to delete it but this also failed:\n%s",
				err)
		} else {
			t.Log(" deleted")
		}
		t.FailNow()
	}

	fmt.Println("Deleting test user")
	if err = tn.DeleteUser(user.ObjectID); err != nil {
		t.Errorf("failed to delete test user")
	}
}
