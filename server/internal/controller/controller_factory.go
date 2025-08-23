package controller

var controllerInstances map[string]IController[any]

func init() {
	controllerInstances = make(map[string]IController[any])
}

// GetController returns the appropriate controller for the task.
func GetController[Schema any](name string) IController[Schema] {
	if controller, ok := controllerInstances[name]; ok {
		// type assertion; equivalent to isinstance(controller, IController[Schema]).
		return controller.(IController[Schema])
	}

	newController := GenericController[Schema]{collectionName: name}
	controllerInstances[name] = newController
	return newController

}
