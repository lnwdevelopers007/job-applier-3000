package schema

type JobSeekerInfo struct {
	Location string `bson:"location" json:"location" binding:"required"`
	Phone    string `bson:"phone" json:"phone" binding:"required"`
	LinkedIn string `bson:"linkedIn" json:"linkedIn"`
}
