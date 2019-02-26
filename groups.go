package tenant

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
)

// GroupResponse simply contains the API response (within the 'value' tag) type for our JSON unmarshaler to put the data into
type GroupResponse struct {
	Groups    []Group `json:"value"`
	ODataNext string  `json:"@odata.nextLink"`
}

// GroupMemberResponse simply contains the API response (within the 'value' tag) type for our JSON unmarshaler to put the data into
type GroupMemberResponse struct {
	GroupMembers []User `json:"value"`
	ODataNext    string `json:"@odata.nextLink"`
}

// The Group struct contains the details from the B2C users
type Group struct {
	ObjectID    string `json:"objectId"`
	DisplayName string `json:"displayName"`
}

// GetGroup returns object of group in the B2C directory
func (t Tenant) GetGroup(GroupObjectID string) (Group, error) {

	parameter := "{\"securityEnabledOnly\": false}"

	ar, err := t.callGraphAPI("/groups/"+GroupObjectID, "1.6", "GET", parameter)
	if err != nil {
		msg := "Error in calling API: " + err.Error()
		log.Println(msg)
	}

	fmt.Print(string(ar))

	group := Group{}

	err = json.Unmarshal(ar, &group)
	if err != nil {
		fmt.Println(err)
	}

	return group, nil
}

// GetGroupMembers returns a groups's direct members from the B2C directory
func (t Tenant) GetGroupMembers(objectID string) ([]User, error) {
	ar, err := t.callNewGraphAPI("/groups/"+objectID+"/members", "GET", "")
	if err != nil {
		msg := "Error in calling API: " + err.Error()
		log.Println(msg)
		return []User{}, fmt.Errorf(msg)
	}

	gmr := GroupMemberResponse{}

	err = json.Unmarshal(ar, &gmr)
	if err != nil {
		fmt.Println(err)
	}

	return gmr.GroupMembers, nil
}

// GetGroups returns a list of groups in the B2C directory
func (t Tenant) GetGroups() ([]Group, error) {
	ar, err := t.callGraphAPI("/groups/", "1.6", "GET", "")
	if err != nil {
		msg := "Error in calling API: " + err.Error()
		log.Println(msg)
	}

	gr := GroupResponse{}

	err = json.Unmarshal(ar, &gr)
	if err != nil {
		fmt.Println(err)
	}

	return gr.Groups, nil
}

// AddGroupMember adds all users with the supplied email address to the specified AAD group
func (t Tenant) AddGroupMember(aadGroup, userEmail string) error {
	if aadGroup == "" {
		return fmt.Errorf("no AAD group specified")
	}

	if userEmail == "" {
		return fmt.Errorf("no user email specified")
	}

	encodedFilter := url.QueryEscape("otherMails/any(x:x eq '" + userEmail + "')")

	response, err := t.callGraphAPI("/users", "1.6", "GET", "$filter="+encodedFilter)
	if err != nil {
		return fmt.Errorf("error while reading user: %s", err)
	}

	ur := UserResponse{}

	err = json.Unmarshal(response, &ur)
	if err != nil {
		return fmt.Errorf("error unmarshaling JSON response: %s", err)
	}

	if len(ur.Users) == 0 {
		return fmt.Errorf("no user with email %s exists", userEmail)
	}

	for _, user := range ur.Users {
		parameter := "{\"url\": \"https://graph.windows.net/" + t.TenantDomain + "/directoryObjects/" + user.ObjectID + "\"}"

		response, err = t.callGraphAPI("/groups/"+aadGroup+"/$links/members", "1.6", "POST", parameter)
		if err != nil {
			return fmt.Errorf("error while adding user: %s\n%s", err, string(response))
		}
		log.Printf("added user %s to group %s", user.ObjectID, aadGroup)
	}

	return nil
}

// DeleteGroupMember adds all users with the supplied email address to the specified AAD group
func (t Tenant) DeleteGroupMember(aadGroup, userEmail string) error {
	if aadGroup == "" {
		return fmt.Errorf("no AAD group specified")
	}

	if userEmail == "" {
		return fmt.Errorf("no user email specified")
	}

	encodedFilter := url.QueryEscape("otherMails/any(x:x eq '" + userEmail + "')")

	response, err := t.callGraphAPI("/users", "1.6", "GET", "$filter="+encodedFilter)
	if err != nil {
		return fmt.Errorf("error while reading user: %s", err)
	}

	ur := UserResponse{}

	err = json.Unmarshal(response, &ur)
	if err != nil {
		return fmt.Errorf("error unmarshaling JSON response: %s", err)
	}

	if len(ur.Users) == 0 {
		return fmt.Errorf("no user with email %s exists", userEmail)
	}

	for _, user := range ur.Users {
		response, err = t.callGraphAPI("/groups/"+aadGroup+"/$links/members/"+user.ObjectID, "1.6", "DELETE", "")
		if err != nil {
			return fmt.Errorf("error while deleting user: %s\n%s", err, string(response))
		}
		log.Printf("deleted user %s from group %s", user.ObjectID, aadGroup)
	}

	return nil
}
