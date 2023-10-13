package entity

type Module map[string]interface{}

// User represents an EPA user.
type User struct {
	ID           string                   `json:"id,omitempty" firestore:"-"`
	Name         string                   `json:"name,omitempty" firestore:"name"`
	LastName     string                   `json:"lastname,omitempty" firestore:"lastname"`
	Email        string                   `json:"email,omitempty" firestore:"email"`
	Password     string                   `json:"password,omitempty" firestore:"-"`
	PhoneNumber  string                   `json:"phone_number,omitempty" firestore:"phone_number"`
	Birthday     string                   `json:"birthday,omitempty" firestore:"birthday"`
	Subscription string                   `json:"subscription,omitempty" firestore:"subscription"`
	Modules      []map[string]interface{} `json:"modulesWithPermission,omitempty" firestore:"modulesWithPermission"`
	Url          string                   `json:"url_image,omitempty" firestore:"url_image"`
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
