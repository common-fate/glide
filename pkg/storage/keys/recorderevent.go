package keys

const RecorderEventKey = "RECORDER_EVENT#"

type recorderEventKeys struct {
	PK1    string
	SK1    func(eventID string) string
	GSI1PK func(requestID string) string
	GSI1SK func(eventID string) string
}

var RecorderEvent = recorderEventKeys{
	PK1:    RecorderEventKey,
	SK1:    func(eventID string) string { return eventID },
	GSI1PK: func(requestID string) string { return RecorderEventKey + requestID },
	GSI1SK: func(eventID string) string { return eventID },
}
