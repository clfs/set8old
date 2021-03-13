package set8

// TaggedMessage represents a message and tag pair, both as byte slices.
type TaggedMessage struct {
	Message, Tag []byte
}
