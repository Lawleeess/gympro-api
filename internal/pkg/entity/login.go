package entity

// StandardLoginCredentials is a custom type to represent email/pass credentials
type StandardLoginCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// AuthResponse represents the output when a user signs in.
type AuthResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
	User         User   `json:"user"`
}

// SignWithCustomTokenResp represents the response from firebase
// when a user is logged with a custom token
type SignWithCustomTokenResp struct {
	Kind         string `json:"kind"`
	Token        string `json:"idToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    string `json:"expiresIn"`
	IsNewUser    bool   `json:"isNewUser"`
	LocalID      string `json:"localId"`
}

// FirebaseErrBody represents the body of err returned from firestore.
type FirebaseErrBody struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// FirebaseError represents the error structure returned when something
// Wrong happened in firebase.
type FirebaseError struct {
	FirebaseErrBody `json:"error"`
}
