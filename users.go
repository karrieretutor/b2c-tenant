package tenant

import (
	"encoding/json"
	"fmt"
	"log"
)

// The MemberGroupIdsResponse struct contains the list of group objectIds the queried user is part of
type MemberGroupIdsResponse struct {
	GroupIds []string `json:"value"`
}

// UserResponse simply contains the API response (within the 'value' tag) type for our JSON unmarshaler to put the data into
type UserResponse struct {
	Users []User `json:"value"`
}

// The User struct contains the details from the B2C users
type User struct {
	ObjectID       string   `json:"objectId"`
	DisplayName    string   `json:"displayName"`
	EmailAddresses []string `json:"otherMails"`
	ID             string   `json:"id"`
}

// GetMemberGroupIDs returns a list of group objectIds the user is part of
// This is for the /getMemberGroups/ handler that is used by the B2C custom policy
func (t Tenant) GetMemberGroupIDs(UserObjectID string) ([]string, error) {

	parameter := "{\"securityEnabledOnly\": false}"

	ar, err := t.callGraphAPI("/users/"+UserObjectID+"/getMemberGroups", "1.6", "POST", parameter)
	if err != nil {
		return nil, fmt.Errorf("error calling Graph API: %s", err)
	}

	mgr := MemberGroupIdsResponse{}

	err = json.Unmarshal(ar, &mgr)
	if err != nil {
		fmt.Println(err)
	}

	return mgr.GroupIds, nil
}

// GetMemberGroupsDetailed returns a list of group objectIds the user is part of
func (t Tenant) GetMemberGroupsDetailed(UserObjectID string) ([]string, error) {

	parameter := "{\"securityEnabledOnly\": false}"

	ar, err := t.callGraphAPI("/users/"+UserObjectID+"/getMemberGroups", "1.6", "POST", parameter)
	if err != nil {
		msg := "Error in calling API: " + err.Error()
		log.Println(msg)
	}

	mgr := MemberGroupIdsResponse{}

	err = json.Unmarshal(ar, &mgr)
	if err != nil {
		fmt.Println(err)
	}

	// Fetching Details of user group
	groups := make([]string, len(mgr.GroupIds))
	for i, groupID := range mgr.GroupIds {
		group, err := t.GetGroup(groupID)
		if err != nil {
			log.Println(err)
		}
		groups[i] = group.DisplayName
	}

	return groups, nil
}

// GetUsers returns a list of users in the B2C directory
func (t Tenant) GetUsers() ([]User, error) {
	ar, err := t.callGraphAPI("/users/", "1.6", "GET", "")
	if err != nil {
		msg := "Error in calling API: " + err.Error()
		log.Println(msg)
	}

	// fmt.Print(string(ar))

	ur := UserResponse{}

	err = json.Unmarshal(ar, &ur)
	if err != nil {
		fmt.Println(err)
	}

	return ur.Users, nil
}

// GetUser returns a single user's details from the B2C directory
func (t Tenant) GetUser(objectID string) (User, error) {
	ar, err := t.callNewGraphAPI("/users/"+objectID, "GET", "")
	if err != nil {
		msg := "Error in calling API: " + err.Error()
		log.Println(msg)
		return User{}, fmt.Errorf(msg)
	}

	ur := User{}

	err = json.Unmarshal(ar, &ur)
	if err != nil {
		fmt.Println(err)
	}

	return ur, nil
}
