package tines

import (
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/tines/go-tines/tines"
)

func resourceTinesCredential() *schema.Resource {
	return &schema.Resource{
		Create: resourceTinesCredentialCreate,
		Read:   resourceTinesCredentialRead,
		Update: resourceTinesCredentialUpdate,
		Delete: resourceTinesCredentialDelete,

		Schema: map[string]*schema.Schema{
			"credential_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"mode": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"team_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"folder_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"value": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"jwt_algorithm": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"jwt_payload": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"jwt_auto_generate_time_claims": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"jwt_private_key": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"oauth_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"oauth_token_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"oauth_client_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"oauth_client_secret": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"oauth_scope": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"oauth_grant_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"aws_authentication_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"aws_access_key": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"aws_secret_key": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"aws_assumed_role_arn": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"aws_assumed_role_external_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"http_request_options": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"http_request_location_of_token": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"mtls_client_certificate": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"mtls_client_private_key": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"mtls_root_certificate": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"read_access": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceTinesCredentialCreate(d *schema.ResourceData, meta interface{}) error {

	name := d.Get("name").(string)
	mode := d.Get("mode").(string)
	teamID := d.Get("team_id").(int)
	folderID := d.Get("folder_id").(int)
	value := d.Get("value").(string)
	jwtAlgorithm := d.Get("jwt_algorithm").(string)
	jwtPayload := d.Get("jwt_payload").(string)
	jwtAutoGenerateTimeClaims := d.Get("jwt_auto_generate_time_claims").(bool)
	jwtPrivateKey := d.Get("jwt_private_key").(string)
	oauthURL := d.Get("oauth_url").(string)
	oauthTokenURL := d.Get("oauth_token_url").(string)
	oauthClientID := d.Get("oauth_client_id").(string)
	oauthClientSecret := d.Get("oauth_client_secret").(string)
	oauthScope := d.Get("oauth_scope").(string)
	oauthGrantType := d.Get("oauth_grant_type").(string)
	awsAuthenticationType := d.Get("aws_authentication_type").(string)
	awsAccessKey := d.Get("aws_access_key").(string)
	awsSecretKey := d.Get("aws_secret_key").(string)
	awsAssumedRoleARN := d.Get("aws_assumed_role_arn").(string)
	awsAssumedRoleExternalID := d.Get("aws_assumed_role_external_id").(string)
	httpRequestOptions := d.Get("http_request_options").(string)
	httpLocationOfToken := d.Get("http_request_location_of_token").(string)
	mtlsClientCertificate := d.Get("mtls_client_certificate").(string)
	mtlsClientPrivateKey := d.Get("mtls_client_private_key").(string)
	mtlsRootCertificate := d.Get("mtls_root_certificate").(string)
	readAccess := d.Get("read_access").(string)
	description := d.Get("description").(string)

	tinesClient := meta.(*tines.Client)

	c := tines.Credential{
		Name:                       name,
		Mode:                       mode,
		TeamID:                     teamID,
		FolderID:                   folderID,
		Value:                      value,
		JWTAlgorithm:               jwtAlgorithm,
		JWTPayload:                 jwtPayload,
		JWTAutoGenerateTimeClaims:  jwtAutoGenerateTimeClaims,
		JWTPrivateKey:              jwtPrivateKey,
		OAuthURL:                   oauthURL,
		OAuthTokenURL:              oauthTokenURL,
		OAuthClientID:              oauthClientID,
		OAuthClientSecret:          oauthClientSecret,
		OAuthScope:                 oauthScope,
		OAuthGrantType:             oauthGrantType,
		AWSAuthenticationType:      awsAuthenticationType,
		AWSAccessKey:               awsAccessKey,
		AWSSecretKey:               awsSecretKey,
		AWSAssumedRoleARN:          awsAssumedRoleARN,
		AWSAssumedRoleExternalID:   awsAssumedRoleExternalID,
		HTTPRequestOptions:         httpRequestOptions,
		HTTPRequestLocationOfToken: httpLocationOfToken,
		MTLSClientCertificate:      mtlsClientCertificate,
		MTLSClientPrivateKey:       mtlsClientPrivateKey,
		MTLSRootCertificate:        mtlsRootCertificate,
		ReadAccess:                 readAccess,
		Description:                description,
	}

	credential, _, err := tinesClient.Credential.Create(&c)
	if err != nil {
		return err
	}

	scid := strconv.Itoa(credential.ID)

	d.SetId(scid)

	return resourceTinesCredentialRead(d, meta)
}

func resourceTinesCredentialRead(d *schema.ResourceData, meta interface{}) error {

	tinesClient := meta.(*tines.Client)

	cid, _ := strconv.ParseInt(d.Id(), 10, 32)
	credential, _, err := tinesClient.Credential.Get(int(cid))
	if err != nil {
		log.Printf("[DEBUG] Error: %v", err)
		if strings.Contains(err.Error(), "404") {
			log.Printf("[INFO] Credential %v no longer exists", d.Id())
			d.SetId("")
			return nil
		} else {
			return err
		}
	}

	scid := strconv.Itoa(credential.ID)

	d.SetId(scid)
	d.Set("name", credential.Name)
	d.Set("mode", credential.Mode)
	d.Set("folder_id", credential.FolderID)
	d.Set("team_id", credential.TeamID)
	d.Set("value", credential.Value)
	d.Set("jwt_algorithm", credential.JWTAlgorithm)
	d.Set("jwt_payload", credential.JWTPayload)
	d.Set("jwt_auto_generate_time_claims", credential.JWTAutoGenerateTimeClaims)
	d.Set("jwt_private_key", credential.JWTPrivateKey)
	d.Set("oauth_url", credential.OAuthURL)
	d.Set("oauth_token_url", credential.OAuthTokenURL)
	d.Set("oauth_client_id", credential.OAuthClientID)
	d.Set("oauth_client_secret", credential.OAuthClientSecret)
	d.Set("oauth_scope", credential.OAuthScope)
	d.Set("oauth_grant_type", credential.OAuthGrantType)
	d.Set("aws_authentication_type", credential.AWSAuthenticationType)
	d.Set("aws_access_key", credential.AWSAccessKey)
	d.Set("aws_secret_key", credential.AWSSecretKey)
	d.Set("aws_assumed_role_arn", credential.AWSAssumedRoleARN)
	d.Set("aws_assumed_role_external_id", credential.AWSAssumedRoleExternalID)
	d.Set("http_request_options", credential.HTTPRequestOptions)
	d.Set("http_request_location_of_token", credential.HTTPRequestLocationOfToken)
	d.Set("mtls_client_certificate", credential.MTLSClientCertificate)
	d.Set("mtls_client_private_key", credential.MTLSClientPrivateKey)
	d.Set("mtls_root_certificate", credential.MTLSRootCertificate)
	d.Set("read_access", credential.ReadAccess)
	d.Set("description", credential.Description)

	return nil
}

func resourceTinesCredentialDelete(d *schema.ResourceData, meta interface{}) error {

	tinesClient := meta.(*tines.Client)
	cid, _ := strconv.ParseInt(d.Id(), 10, 32)
	_, err := tinesClient.Credential.Delete(int(cid))
	if err != nil {
		return err
	}
	d.SetId("")

	return nil
}

func resourceTinesCredentialUpdate(d *schema.ResourceData, meta interface{}) error {

	tinesClient := meta.(*tines.Client)

	cid, _ := strconv.ParseInt(d.Id(), 10, 32)
	name := d.Get("name").(string)
	mode := d.Get("mode").(string)
	teamID := d.Get("team_id").(int)
	folderID := d.Get("folder_id").(int)
	value := d.Get("value").(string)
	jwtAlgorithm := d.Get("jwt_algorithm").(string)
	jwtPayload := d.Get("jwt_payload").(string)
	jwtAutoGenerateTimeClaims := d.Get("jwt_auto_generate_time_claims").(bool)
	jwtPrivateKey := d.Get("jwt_private_key").(string)
	oauthURL := d.Get("oauth_url").(string)
	oauthTokenURL := d.Get("oauth_token_url").(string)
	oauthClientID := d.Get("oauth_client_id").(string)
	oauthClientSecret := d.Get("oauth_client_secret").(string)
	oauthScope := d.Get("oauth_scope").(string)
	oauthGrantType := d.Get("oauth_grant_type").(string)
	awsAuthenticationType := d.Get("aws_authentication_type").(string)
	awsAccessKey := d.Get("aws_access_key").(string)
	awsSecretKey := d.Get("aws_secret_key").(string)
	awsAssumedRoleARN := d.Get("aws_assumed_role_arn").(string)
	awsAssumedRoleExternalID := d.Get("aws_assumed_role_external_id").(string)
	httpRequestOptions := d.Get("http_request_options").(string)
	httpLocationOfToken := d.Get("http_request_location_of_token").(string)
	mtlsClientCertificate := d.Get("mtls_client_certificate").(string)
	mtlsClientPrivateKey := d.Get("mtls_client_private_key").(string)
	mtlsRootCertificate := d.Get("mtls_root_certificate").(string)
	readAccess := d.Get("read_access").(string)
	description := d.Get("description").(string)

	gr := tines.Credential{
		Name:                       name,
		Mode:                       mode,
		TeamID:                     teamID,
		FolderID:                   folderID,
		Value:                      value,
		JWTAlgorithm:               jwtAlgorithm,
		JWTPayload:                 jwtPayload,
		JWTAutoGenerateTimeClaims:  jwtAutoGenerateTimeClaims,
		JWTPrivateKey:              jwtPrivateKey,
		OAuthURL:                   oauthURL,
		OAuthTokenURL:              oauthTokenURL,
		OAuthClientID:              oauthClientID,
		OAuthClientSecret:          oauthClientSecret,
		OAuthScope:                 oauthScope,
		OAuthGrantType:             oauthGrantType,
		AWSAuthenticationType:      awsAuthenticationType,
		AWSAccessKey:               awsAccessKey,
		AWSSecretKey:               awsSecretKey,
		AWSAssumedRoleARN:          awsAssumedRoleARN,
		AWSAssumedRoleExternalID:   awsAssumedRoleExternalID,
		HTTPRequestOptions:         httpRequestOptions,
		HTTPRequestLocationOfToken: httpLocationOfToken,
		MTLSClientCertificate:      mtlsClientCertificate,
		MTLSClientPrivateKey:       mtlsClientPrivateKey,
		MTLSRootCertificate:        mtlsRootCertificate,
		ReadAccess:                 readAccess,
		Description:                description,
	}

	credential, _, err := tinesClient.Credential.Update(int(cid), &gr)
	if err != nil {
		return err
	}

	scid := strconv.Itoa(credential.ID)

	d.SetId(scid)

	return resourceTinesCredentialRead(d, meta)
}
