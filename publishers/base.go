package publishers

type Publishers []Publisher
type Publisher interface {
	Publish()
}

func PublishTasks() {
	publishers := Publishers{}

	for _, publisher := range publishers {
		publisher.Publish()
	}
}
