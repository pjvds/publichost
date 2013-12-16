package protocol_test

import (
	"bytes"
	"fmt"
	"github.com/pjvds/publichost/protocol/request"
	"testing"
)

func TestReadRequestSetsHeaderEntry(t *testing.T) {
	key := "MyKey"
	value := "MyValue"

	data := bytes.NewBufferString(fmt.Sprintf("%v:%v\r\n", key, value))

	request, err := request.ReadRequest(data)
	if err != nil {
		t.Errorf("Unexpected error for valid request: %v", err)
	}

	value, ok := request.Header["MyKey"]
	if !ok {
		t.Errorf("Missing MyKey in header")
	} else if value != "MyValue" {
		t.Errorf("MyKey has wrong value %v, expected %v", value, "MyValue")
	}
}
