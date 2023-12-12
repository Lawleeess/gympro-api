package entity

type Module map[string]interface{}

// User represents an EPA user.
type User struct {
	ID           string                   `json:"id,omitempty" firestore:"-"`
	OobCode      string                   `json:"oobCode,omitempty" firestore:"oobCode"`
	IsVerified   bool                     `json:"isVerified,omitempty" firestore:"isVerified"`
	Name         string                   `json:"name,omitempty" firestore:"name"`
	LastName     string                   `json:"lastname,omitempty" firestore:"lastname"`
	Email        string                   `json:"email,omitempty" firestore:"email"`
	Password     string                   `json:"password,omitempty" firestore:"password"`
	PhoneNumber  string                   `json:"phone_number,omitempty" firestore:"phone_number"`
	Birthday     string                   `json:"birthday,omitempty" firestore:"birthday"`
	Subscription string                   `json:"subscription" firestore:"subscription"`
	Modules      []map[string]interface{} `json:"modulesWithPermission" firestore:"modulesWithPermission"`
	Url          string                   `json:"url_image,omitempty" firestore:"url_image"`
	UserRole     string                   `json:"user_role,omitempty" firestore:"user_role"`
	UserProgress UserProgress             `json:"userProgress,omitempty" firestore:"userProgress"`
	UserGoals    UserGoals                `json:"userGoals,omitempty" firestore:"userGoals"`
	UserRoutine  UserRoutine              `json:"userRoutine" firestore:"userRoutine"`
}

type UserProgress struct {
	Age      int     `json:"age,omitempty" firestore:"age"`
	Gender   string  `json:"gender,omitempty" firestore:"gender"`
	Height   int     `json:"height,omitempty" firestore:"height"`
	Weight   float64 `json:"weight,omitempty" firestore:"weight"`
	Activity string  `json:"activity,omitempty" firestore:"activity"`
	Goal     string  `json:"goal,omitempty" firestore:"goal"`
}

type UserGoals struct {
	IMC     string `json:"imc,omitempty" firestore:"imc"`
	BMR     string `json:"bmr,omitempty" firestore:"bmr"`
	TDEE    string `json:"tdee,omitempty" firestore:"tdee"`
	Goal    string `json:"goal,omitempty" firestore:"goal"`
	Protein string `json:"protein,omitempty" firestore:"protein"`
	Carbs   string `json:"carbs,omitempty" firestore:"carbs"`
	Fat     string `json:"fat,omitempty" firestore:"fat"`
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

type UserRoutine struct {
	Monday    []Routine `json:"monday" firestore:"monday"`
	Tuesday   []Routine `json:"tuesday" firestore:"tuesday"`
	Wednesday []Routine `json:"wednesday" firestore:"wednesday"`
	Thursday  []Routine `json:"thursday" firestore:"thursday"`
	Friday    []Routine `json:"friday" firestore:"friday"`
	Saturday  []Routine `json:"saturday" firestore:"saturday"`
}
