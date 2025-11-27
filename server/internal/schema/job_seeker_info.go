package schema

type JobSeekerInfo struct {
	Location string `bson:"location" json:"location" binding:"required,min=2,max=200"`
	Phone    string `bson:"phone" json:"phone" binding:"required,min=8,max=20"`
	LinkedIn string `bson:"linkedIn" json:"linkedIn" binding:"omitempty,url,max=200"`
}
