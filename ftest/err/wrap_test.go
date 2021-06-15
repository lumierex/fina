package err

import (
	"log"
	"testing"
)

func Test_handler(t *testing.T) {
	bytes, err := handler()
	if err != nil {
		t.Logf("err: %+v", err)
		//t.Log(err)
		return
	}
	log.Println(string(bytes))
}
