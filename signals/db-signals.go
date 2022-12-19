package signals

type DBSignals struct {
	Type      string `json:"type" bson:"type"`
	Timestamp int64  `json:"timestamp" bson:"timestamp"`
}
