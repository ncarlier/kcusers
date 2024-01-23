package uuid

import "regexp"

const uuidPattern = "[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}"

var uuidRe = regexp.MustCompile("^" + uuidPattern + "$")
var uuidPrefixRe = regexp.MustCompile("^" + uuidPattern)

// GetUUIDPrefix return UUID prfix of a string
func GetUUIDPrefix(value string) (uuid string, ok bool) {
	matches := uuidPrefixRe.FindStringSubmatch(value)
	if len(matches) == 0 {
		return "", false
	}
	return matches[0], true
}

// IsUUID return true if value is a UUID
func IsUUID(value string) bool {
	return uuidRe.MatchString(value)
}
