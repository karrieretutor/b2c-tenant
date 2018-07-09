package tenant

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

const loginURL string = "https://login.microsoftonline.com/"

// Tenant contains the data of the app registration in Azure AD that has write permissions in the AAD tenant
// also the Access Token that gets returned and will be used for accessing the API
type Tenant struct {
	ClientID     string
	ClientSecret string
	TenantDomain string
	AccessToken  AccessToken
}

// AccessToken contains an OAuth2 access token for use with Azure AD Graph API calls
type AccessToken struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

// AuthCountResponse simply contains the API response (within the 'value' tag) type for our JSON unmarshaler to put the data into
type AuthCountResponse struct {
	Value []AuthenticationCount `json:"value"`
}

// The AuthenticationCount struct contains the details from the B2C authentication count report
type AuthenticationCount struct {
	B2CAuthenticationCount float64 `json:"AuthenticationCount"`
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
}

// The MemberGroupIdsResponse struct contains the list of group objectIds the queried user is part of
type MemberGroupIdsResponse struct {
	GroupIds []string `json:"value"`
}

// GroupResponse simply contains the API response (within the 'value' tag) type for our JSON unmarshaler to put the data into
type GroupResponse struct {
	Groups []Group `json:"value"`
}

// The Group struct contains the details from the B2C users
type Group struct {
	ObjectID    string `json:"objectId"`
	DisplayName string `json:"displayName"`
}

// GetB2CAuthenticationCount returns the count of B2C authentications in the last 30 days
// Unfortunately, the 30 days is a limit of the upstream API
// so we need a way to account for that in a rolling fashion in Prometheus
func (t Tenant) GetB2CAuthenticationCount() (float64, error) {

	ar, err := t.callGraphAPI("/reports/b2cAuthenticationCount/", "beta", "GET", "")
	if err != nil {
		msg := "Error in calling API: " + err.Error()
		log.Println(msg)
	}

	acr := AuthCountResponse{}

	err = json.Unmarshal(ar, &acr)
	if err != nil {
		fmt.Println(err)
	}

	return acr.Value[0].B2CAuthenticationCount, nil
}

// GetMemberGroupIDs returns a list of group objectIds the user is part of
// This is for the /getMemberGroups/ handler that is used by the B2C custom policy
func (t Tenant) GetMemberGroupIDs(UserObjectID string) ([]string, error) {

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

// CallGraphAPI does the API call to the Azure AD Graph API and returns the response as an APIRespone struct
func (t Tenant) callGraphAPI(endpoint string, apiversion string, method string, param string) ([]byte, error) {
	requestString := "https://graph.windows.net/" + t.TenantDomain + endpoint + "?api-version=" + apiversion
	log.Printf("Calling %s \n", requestString)

	// fmt.Println(requestString)

	client := &http.Client{}

	req, err := http.NewRequest(method, requestString, nil)

	if method == "POST" && param != "" {
		postReader := strings.NewReader(param)
		req, err = http.NewRequest(method, requestString, postReader)
	}

	if err != nil {
		fmt.Println(err)
		return []byte{}, err
	}

	req.Header.Add("Authorization", t.AccessToken.TokenType+" "+t.AccessToken.AccessToken)

	if method == "POST" {
		req.Header.Add("Content-Type", "application/json")
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return []byte{}, err
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)

	//fmt.Print(string(bodyBytes))

	if resp.StatusCode != 200 {
		return bodyBytes, err
	}

	return bodyBytes, nil
}

// GetUserDetails

// GetAccessToken returns the access token for API access
func (t *Tenant) GetAccessToken() error {

	authAuthenticatorURL := loginURL + t.TenantDomain + "/oauth2/token?api-version=1.0"

	parameters := url.Values{
		"client_id":     {t.ClientID},
		"client_secret": {t.ClientSecret},
		"grant_type":    {"client_credentials"},
	}

	resp, err := http.PostForm(authAuthenticatorURL, parameters)
	if err != nil {
		fmt.Printf("Error in POSTing the token request: %s\n", err)
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Printf("Error in reading the response body: %s\n", err)
	}

	if resp.StatusCode != 200 {
		fmt.Printf("Error while calling auth endpoint %s\n", authAuthenticatorURL)
		fmt.Print(string(bodyBytes))
	}

	at := AccessToken{}

	err = json.Unmarshal(bodyBytes, &at)
	if err != nil {
		fmt.Printf("Error getting the Access token: %s\n", err)
		return err
	}

	t.AccessToken = at
	return nil
}
