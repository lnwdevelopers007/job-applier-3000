package schema

// CollectionEntity is the schema that get saved to MongoDB database.
type CollectionEntity interface {
	GetCollectionName() string
}
