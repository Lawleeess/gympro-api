package entity

type Routine struct {
	ID          string `json:"id,omitempty" firestore:"id"`
	MuscleGroup string `json:"muscle_group,omitempty" firestore:"muscle_group"`
	Name        string `json:"name,omitempty" firestore:"name"`
	Description string `json:"description,omitempty" firestore:"description"`
	VideoUrl    string `json:"video_url,omitempty" firestore:"video_url"`
	ImageUrl    string `json:"url_image,omitempty" firestore:"url_image"`
}
