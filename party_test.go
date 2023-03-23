package party

import "testing"

func TestJson(t *testing.T) {

	method := MethodGet
	url := "http://facebook.com"
	response, err := NewClientParty(method, url).HitClient()
	if err != nil {
		t.Error("ERROR ", *err)
		return
	}

	t.Log("SUCCESS ", response)
}
