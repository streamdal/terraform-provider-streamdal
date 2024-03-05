package util

import (
	"fmt"
	"strings"

	"github.com/streamdal/streamdal/libs/protos/build/go/protos"
)

func AudienceToStr(audience *protos.Audience) string {
	if audience == nil {
		return ""
	}

	str := strings.ToLower(fmt.Sprintf("%s:%s:%s:%s", audience.ServiceName, audience.OperationType, audience.OperationName, audience.ComponentName))

	return str
}

// AudienceFromStr will parse a string into an Audience.
// If the string is invalid, nil will be returned.
func AudienceFromStr(s string) *protos.Audience {
	if s == "" {
		return nil
	}

	parts := strings.Split(s, ":")
	if len(parts) != 4 {
		return nil
	}

	opType := protos.OperationType_OPERATION_TYPE_UNSET

	if parts[1] == strings.ToLower(protos.OperationType_OPERATION_TYPE_CONSUMER.String()) {
		opType = protos.OperationType_OPERATION_TYPE_CONSUMER
	} else {
		opType = protos.OperationType_OPERATION_TYPE_PRODUCER
	}

	return &protos.Audience{
		ServiceName:   strings.ToLower(parts[0]),
		OperationType: opType,
		OperationName: strings.ToLower(parts[2]),
		ComponentName: strings.ToLower(parts[3]),
	}
}

func AudienceEquals(a, b *protos.Audience) bool {
	if a == nil || b == nil {
		return false
	}

	return AudienceToStr(a) == AudienceToStr(b)
}

func AudienceInList(audience *protos.Audience, list []*protos.Audience) bool {
	for _, a := range list {
		if AudienceEquals(a, audience) {
			return true
		}
	}

	return false
}
