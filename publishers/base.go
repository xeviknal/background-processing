package publishers

type Publishers []Publisher
type Publisher interface {
	Publish()
}

func PublishTasks() {
	publishers := Publishers{
		JobsPublisher{},
	}

	for _, publisher := range publishers {
		go publisher.Publish()
	}
}
