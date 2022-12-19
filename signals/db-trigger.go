package signals

type DBTrigger struct {
	SignalType string `json:"signal_type" bson:"signal_type"`
	Action     string `json:"action" bson:"action"`
}
