package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/url"

	"firebase.google.com/go/auth"
	"github.com/CValier/gympro-api/internal/pkg/config"
	"github.com/CValier/gympro-api/internal/pkg/entity"
	"github.com/CValier/gympro-api/internal/pkg/ports"
	"github.com/epa-datos/errors"
)

type authSvc struct {
	client ports.FirebaseCli
}

// NewAuthService returns a new instance for auth service
func NewAuthService(client ports.FirebaseCli) *authSvc {
	return &authSvc{
		client: client,
	}
}

// SignUpWithEmailAndPass creates a new email and password user by issuing an HTTP POST request
// to the Auth signupNewUser endpoint.
// for more details: https://firebase.google.com/docs/reference/rest/auth?hl=en
// It returns the id of the created user as string.
func (a *authSvc) SignUpWithEmailAndPass(email, pass string) (string, error) {
	// Creating a body payload to send it to firebase API.
	bodyReq := entity.StandardLoginCredentials{
		Email:    email,
		Password: pass,
	}

	// Encoding the body payload
	body, err := json.Marshal(bodyReq)
	if err != nil {
		return "", err
	}

	// We need to add the firebase key to the request as url param
	parm := url.Values{}
	parm.Add("key", config.CfgIn.FirebaseKey)

	reqURL := config.CfgIn.FirebaseHost + ":signUp?" + parm.Encode()

	// Creating a new request
	req, err := http.NewRequest(http.MethodPost, reqURL, bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", "application/json")

	// Making the http request
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil || resp.StatusCode != http.StatusOK {
		// Decoding response in the proper err struct.
		dec := json.NewDecoder(resp.Body)
		response := &entity.FirebaseError{}
		if err := dec.Decode(response); err != nil {
			return "", err
		}

		// Check if the error was triggered because the email sent already exists.
		if resp.StatusCode == http.StatusBadRequest && response.Message == "EMAIL_EXISTS" {
			return "", err
		}
		return "", err
	}
	defer resp.Body.Close()

	// Decoding response
	decoder := json.NewDecoder(resp.Body)
	response := &entity.SignWithCustomTokenResp{}

	err = decoder.Decode(response)
	return response.LocalID, err
}

// VerifyToken verifies in Firebase if the given token is valid.
func (a *authSvc) VerifyToken(token string) (*auth.Token, error) {
	return a.client.VerifyToken(token)
}

// GenerateCustomToken creates a token with the claims given encrypted inside the JWT.
func (a *authSvc) GenerateCustomToken(ctx context.Context, userID string, claims map[string]interface{}) (string, error) {
	return a.client.GenerateCustomToken(userID, claims)
}

// RevokeUserTokens removes all the active refresh tokens for a particular user.
func (a *authSvc) RevokeUserTokens(userID string) error {
	return a.client.RevokeUserTokens(userID)
}

// SignInWithPass makes a HTTP request to Firebase API to signIn a user with email and pass
// For more details: https://firebase.google.com/docs/reference/rest/auth?hl=en
func (a *authSvc) SignInWithPass(ctx context.Context, creds *entity.StandardLoginCredentials) (string, error) {
	scope := errors.Operation("auth_service.SignInWithPass")

	// Creating a body payload to send it to firebase REST API
	bodyReq := struct {
		entity.StandardLoginCredentials
		ReturnSecureToken bool `json:"returnSecureToken"`
	}{
		*creds,
		true,
	}

	// Encoding the body payload
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(bodyReq)
	if err != nil {
		log.Fatal(err)
	}

	// We need to add the firebase key to the request as url param
	parm := url.Values{}
	parm.Add("key", config.CfgIn.FirebaseKey)

	reqURL := config.CfgIn.FirebaseHost + ":signInWithPassword?" + parm.Encode()

	// Creating a new request
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, &buf)
	if err != nil {
		return "", errors.Build(
			scope,
			errors.InternalError,
			errors.Message("Failed to create a new request: "+err.Error()),
		)
	}
	req.Header.Add("Content-Type", "application/json")

	// Making the http request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		return "", errors.Build(
			scope,
			errors.NotFound,
			errors.Message("Invalid email/password combination"),
		)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	response := entity.SignWithCustomTokenResp{}

	return response.Token, decoder.Decode(&response)
}

// SignInWithTokenClaims makes a HTTP request to Firebase REST API
// To exchange a custom token for a id token.
// For more details: https://firebase.google.com/docs/reference/rest/auth?hl=en#section-verify-custom-token
func (a *authSvc) SignInWithTokenClaims(ctx context.Context, customToken string) (*entity.SignWithCustomTokenResp, error) {
	scope := errors.Operation("auth_service.SignInWithTokenClaims")

	// Creating a body payload to send it to firebase REST API
	bodyReq := struct {
		Token             string `json:"token"`
		ReturnSecureToken bool   `json:"returnSecureToken"`
	}{
		customToken,
		true,
	}

	// Encoding the body payload
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(bodyReq)
	if err != nil {
		log.Fatal(err)
	}

	// We need to add the firebase key to the request as url param
	parm := url.Values{}
	parm.Add("key", config.CfgIn.FirebaseKey)

	reqURL := config.CfgIn.FirebaseHost + ":signInWithCustomToken?" + parm.Encode()

	// Creating a new request
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, &buf)
	if err != nil {
		return nil, errors.Build(
			scope,
			errors.InternalError,
			errors.Message("Failed to create a new request: "+err.Error()),
		)
	}
	req.Header.Add("Content-Type", "application/json")

	// Making the http request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, errors.Build(
			scope,
			errors.NotFound,
			errors.Message("Invalid custom token"),
		)
	}
	defer resp.Body.Close()

	// Decoding response
	decoder := json.NewDecoder(resp.Body)
	response := &entity.SignWithCustomTokenResp{}

	return response, decoder.Decode(response)
}

// RemoveUser deletes a user with given token ID.
// https://firebase.google.com/docs/reference/rest/auth?hl=en#section-delete-account
func (a *authSvc) RemoveUser(idToken string) {
	// Creating a body payload to send it to firebase API.
	bodyReq := struct {
		IDToken string `json:"idToken"`
	}{
		IDToken: idToken,
	}

	// Encoding the body payload
	body, err := json.Marshal(bodyReq)
	if err != nil {
		return
	}

	// We need to add the firebase key to the request as url param
	parm := url.Values{}
	parm.Add("key", config.CfgIn.FirebaseKey)

	reqURL := config.CfgIn.FirebaseHost + ":delete?" + parm.Encode()

	// Creating a new request
	req, _ := http.NewRequest(http.MethodPost, reqURL, bytes.NewBuffer(body))
	req.Header.Add("Content-Type", "application/json")

	// Making the http request
	client := &http.Client{}
	resp, _ := client.Do(req)
	resp.Body.Close()
}
