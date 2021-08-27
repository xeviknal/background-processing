package publishers

type JobsPublisher struct{}

// Checking that JobsPublisher implements
//  Publisher interface at compile time
var _ Publisher = (*JobsPublisher)(nil)

func (t JobsPublisher) Publish() {

}
