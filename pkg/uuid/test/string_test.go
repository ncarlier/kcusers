package test

import (
	"testing"

	"github.com/ncarlier/kcusers/pkg/uuid"
)

func TestGetUUIDPrefix(t *testing.T) {
	val, ok := uuid.GetUUIDPrefix("2b7efa70-8ecc-47e8-97d3-98a63804320a, foo, bar")
	if !ok {
		t.Error("should extract UUID")
	}
	if val != "2b7efa70-8ecc-47e8-97d3-98a63804320a" {
		t.Error("invalid extracted UUID")
	}
}

func TestIsUUID(t *testing.T) {
	if !uuid.IsUUID("2b7efa70-8ecc-47e8-97d3-98a63804320a") {
		t.Error("should be an UUID")
	}
	if uuid.IsUUID("2b7efa70-8ecc 47e8-97d3-98a63804320a") {
		t.Error("should NOT be an UUID")
	}
}
