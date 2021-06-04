package tines

import (
	"context"
	"encoding/json"
	"reflect"
	"strings"

	"fmt"

	"github.com/fatih/structs"

	"github.com/trivago/tgo/tcontainer"
)

// CredentialService handles fields for the Tines instance / API.
type CredentialService struct {
	client *Client
}

// Credential structure
type Credential struct {
	ID                         int    `json:"id" structs:"id,omitempty"`
	Name                       string `json:"name" structs:"name,omitempty"`
	Mode                       string `json:"mode" structs:"mode,omitempty"`
	TeamID                     int    `json:"team_id" structs:"team_id,omitempty"`
	FolderID                   int    `json:"folder_id" structs:"folder_id,omitempty"`
	ReadAccess                 string `json:"read_access" structs:"read_access,omitempty"`
	Value                      string `json:"value" structs:"value,omitempty"`
	JWTAlgorithm               string `json:"jwt_algorithm" structs:"jwt_algorithm,omitempty"`
	JWTPayload                 string `json:"jwt_payload" structs:"jwt_payload,omitempty"`
	JWTAutoGenerateTimeClaims  bool   `json:"jwt_auto_generate_time_claims" structs:"jwt_auto_generate_time_claims,omitempty"`
	JWTPrivateKey              string `json:"jwt_payload" structs:"jwt_payload,omitempty"`
	OAuthURL                   string `json:"oauth_url" structs:"oauth_url,omitempty"`
	OAuthTokenURL              string `json:"oauth_token_url" structs:"oauth_token_url,omitempty"`
	OAuthClientID              string `json:"oauth_client_id" structs:"oauth_client_id,omitempty"`
	OAuthClientSecret          string `json:"oauth_client_secret" structs:"oauth_client_secret,omitempty"`
	OAuthScope                 string `json:"oauth_scope" structs:"oauth_scope,omitempty"`
	OAuthGrantType             string `json:"oauth_grant_type" structs:"oauth_grant_type,omitempty"`
	AWSAuthenticationType      string `json:"aws_authentication_type" structs:"aws_authentication_type,omitempty"`
	AWSAccessKey               string `json:"aws_access_key" structs:"aws_access_key,omitempty"`
	AWSSecretKey               string `json:"aws_secret_key" structs:"aws_secret_key,omitempty"`
	AWSAssumedRoleARN          string `json:"aws_assumed_role_arn" structs:"aws_assumed_role_arn,omitempty"`
	AWSAssumedRoleExternalID   string `json:"aws_assumed_role_external_id" structs:"aws_assumed_role_external_id,omitempty"`
	HTTPRequestOptions         string `json:"http_request_options" structs:"http_request_options,omitempty"`
	HTTPRequestLocationOfToken string `json:"http_request_location_of_token" structs:"http_request_location_of_token,omitempty"`
	MTLSClientCertificate      string `json:"mtls_client_certificate" structs:"mtls_client_certificate,omitempty"`
	MTLSClientPrivateKey       string `json:"mtls_client_private_key" structs:"mtls_client_private_key,omitempty"`
	MTLSRootCertificate        string `json:"mtls_root_certificate" structs:"mtls_root_certificate,omitempty"`
	Unknowns                   tcontainer.MarshalMap
}

// MarshalJSON is a custom JSON marshal function for the Credential* structs.
// It handles Tines options and maps those from / to "Unknowns" key.
func (i *Credential) MarshalJSON() ([]byte, error) {
	m := structs.Map(i)
	unknowns, okay := m["Unknowns"]
	if okay {
		// if unknowns present, shift all key value from unknown to a level up
		for key, value := range unknowns.(tcontainer.MarshalMap) {
			m[key] = value
		}
		delete(m, "Unknowns")
	}
	return json.Marshal(m)
}

// UnmarshalJSON is a custom JSON marshal function for the Credential structs.
// It handles Tines custom options and maps those from / to "Unknowns" key.
func (i *Credential) UnmarshalJSON(data []byte) error {

	// Do the normal unmarshalling first
	// Details for this way: http://choly.ca/post/go-json-marshalling/
	type Alias Credential
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(i),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	totalMap := tcontainer.NewMarshalMap()
	err := json.Unmarshal(data, &totalMap)
	if err != nil {
		return err
	}

	t := reflect.TypeOf(*i)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tagDetail := field.Tag.Get("json")
		if tagDetail == "" {
			// ignore if there are no tags
			continue
		}
		options := strings.Split(tagDetail, ",")

		if len(options) == 0 {
			return fmt.Errorf("no tags options found for %s", field.Name)
		}
		// the first one is the json tag
		key := options[0]
		if _, okay := totalMap.Value(key); okay {
			delete(totalMap, key)
		}

	}
	i = (*Credential)(aux.Alias)
	// all the tags found in the struct were removed. Whatever is left are unknowns to struct
	i.Unknowns = totalMap
	return nil

}

// GetWithContext returns a credential for the given credential id.
func (s *CredentialService) GetWithContext(ctx context.Context, credentialID int) (*Credential, *Response, error) {
	apiEndpoint := fmt.Sprintf("user_credentials/%v", credentialID)
	req, err := s.client.NewRequestWithContext(ctx, "GET", apiEndpoint, nil)
	if err != nil {
		return nil, nil, err
	}

	credential := new(Credential)
	resp, err := s.client.Do(req, credential)
	if err != nil {
		return nil, resp, err
	}

	return credential, resp, nil
}

// Get wraps GetWithContext using the background context.
func (s *CredentialService) Get(credentialID int) (*Credential, *Response, error) {
	return s.GetWithContext(context.Background(), credentialID)
}

// DeleteWithContext deletes a credential.
func (s *CredentialService) DeleteWithContext(ctx context.Context, credentialID int) (*Response, error) {
	apiEndpoint := fmt.Sprintf("user_credentials/%v", credentialID)
	req, err := s.client.NewRequestWithContext(ctx, "DELETE", apiEndpoint, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// Delete wraps DeleteWithContext using the background context.
func (s *CredentialService) Delete(credentialID int) (*Response, error) {
	return s.DeleteWithContext(context.Background(), credentialID)
}

// CreateWithContext creates a credential.
func (s *CredentialService) CreateWithContext(ctx context.Context, credential *Credential) (*Credential, *Response, error) {
	apiEndpoint := fmt.Sprintf("user_credentials")
	req, err := s.client.NewRequestWithContext(ctx, "POST", apiEndpoint, credential)
	if err != nil {
		return nil, nil, err
	}

	credentialresp := new(Credential)
	resp, err := s.client.Do(req, credentialresp)
	if err != nil {
		return nil, resp, err
	}

	return credentialresp, resp, err
}

// Create wraps CreateWithContext using the background context.
func (s *CredentialService) Create(credential *Credential) (*Credential, *Response, error) {
	return s.CreateWithContext(context.Background(), credential)
}

// UpdateWithContext updates a credential.
func (s *CredentialService) UpdateWithContext(ctx context.Context, credentialID int, credential *Credential) (*Credential, *Response, error) {
	apiEndpoint := fmt.Sprintf("user_credentials/%v", credentialID)
	req, err := s.client.NewRequestWithContext(ctx, "PUT", apiEndpoint, credential)
	if err != nil {
		return nil, nil, err
	}

	credentialresp := new(Credential)
	resp, err := s.client.Do(req, credentialresp)
	if err != nil {
		return nil, resp, err
	}

	return credentialresp, resp, err
}

// Update wraps UpdateWithContext using the background context.
func (s *CredentialService) Update(globalResouceID int, credential *Credential) (*Credential, *Response, error) {
	return s.UpdateWithContext(context.Background(), globalResouceID, credential)
}
