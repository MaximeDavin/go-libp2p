package protocol

// ID is an identifier used to write protocol headers in streams.
type ID string

// ConvertToStrings is a convenience function that takes a slice of protocol.ID and
// converts it to a slice of strings.
func ConvertToStrings(ids []ID) (res []string) {}
