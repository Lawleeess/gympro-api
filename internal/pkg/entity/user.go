package entity

type Module map[string]interface{}

// User represents an EPA user.
type User struct {
	ID           string `json:"id" firestore:"-"`
	Name         string `json:"name" binding:"required" firestore:"name"`
	LastName     string `json:"lastname" binding:"required" firestore:"lastname"`
	Email        string `json:"email" binding:"required" firestore:"email"`
	Password     string `json:"password,omitempty" binding:"required" firestore:"-"`
	PhoneNumber  string `json:"phone_number,omitempty" binding:"required" firestore:"phone_number"`
	Birthday     string `json:"birthday,omitempty" binding:"required" firestore:"birthday"`
	Subscription string `json:"subscription,omitempty" firestore:"subscription"`
	Role         string `json:"role,omitempty" firestore:"role"`
}

// UpdateClientPermissionsReq represents a request to change clientsWithPermissions for a user with an array of ClientsIDsToUpdate entity
type UpdateClientPermissionsReq struct {
	Clients []ClientsIDsToUpdate `json:"clients"`
}

// ClientsIDsToUpdate contains the ID for clientsWithPermissionsArray when its requested a change, used to avoid send extra info to firestore field
type ClientsIDsToUpdate struct {
	ID string `json:"id" firestore:"id"`
}

// UsersResponse means the response struct to be returned on users management endpoint
type UsersResponse struct {
	TotalItems int     `json:"totalItems"`
	Items      []*User `json:"items"`
}
