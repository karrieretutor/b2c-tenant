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

// CallGraphAPI does the API call to the Azure AD Graph API and returns the response as an APIRespone struct
func (t Tenant) callGraphAPI(endpoint string, apiversion string, method string, param string) ([]byte, error) {
	requestString := "https://graph.windows.net/" + t.TenantDomain + endpoint + "?api-version=" + apiversion

	// fmt.Println(requestString)

	client := &http.Client{}

	if method == "GET" && param != "" {
		requestString = requestString + "&" + param
	}

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

	log.Printf("Calling %s \n", req.URL)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return []byte{}, err
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)

	if resp.StatusCode > 204 {
		err = fmt.Errorf("Failed API call; status code: %s", resp.Status)
		return bodyBytes, err
	}

	return bodyBytes, nil
}

// CallNewGraphAPI does the API call to the Azure AD Graph API and returns the response as an APIRespone struct
func (t Tenant) callNewGraphAPI(endpoint string, method string, param string) ([]byte, error) {
	requestString := "https://graph.microsoft.com/beta" + endpoint

	// fmt.Println(requestString)

	client := &http.Client{}

	if method == "odatanext" {
		requestString = endpoint
		method = "GET"
	}

	if method == "GET" && param != "" {
		requestString = requestString + "?" + param
	}

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

	log.Printf("Calling %s \n", req.URL)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return []byte{}, err
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)

	if resp.StatusCode > 204 {
		err = fmt.Errorf("Failed API call; status code: %s", resp.Status)
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
		return fmt.Errorf("error in reading the response body: %s", err)
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("error while calling auth endpoint %s: %s", authAuthenticatorURL, string(bodyBytes))
	}

	at := AccessToken{}

	err = json.Unmarshal(bodyBytes, &at)
	if err != nil {
		return fmt.Errorf("Error getting the Access token: %s", err)
	}

	t.AccessToken = at
	return nil
}

// GetGraphAccessToken returns the access token for API access
func (t *Tenant) GetGraphAccessToken() error {

	authAuthenticatorURL := loginURL + t.TenantDomain + "/oauth2/v2.0/token"

	parameters := url.Values{
		"client_id":     {t.ClientID},
		"client_secret": {t.ClientSecret},
		"grant_type":    {"client_credentials"},
		"scope":         {"https://graph.microsoft.com/.default"},
	}

	resp, err := http.PostForm(authAuthenticatorURL, parameters)
	if err != nil {
		return fmt.Errorf("Error in POSTing the token request: %s\n", err)
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return fmt.Errorf("error in reading the response body: %s", err)
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("error while calling auth endpoint %s: %s", authAuthenticatorURL, string(bodyBytes))
	}

	at := AccessToken{}

	err = json.Unmarshal(bodyBytes, &at)
	if err != nil {
		return fmt.Errorf("Error getting the Access token: %s", err)
	}

	t.AccessToken = at
	return nil
}
