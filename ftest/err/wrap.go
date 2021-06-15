package err

import (
	"encoding/json"
	"github.com/pkg/errors"
	"math"
)
import e "errors"

//import "errors"

var errorNotFound = e.New("数据为取到")

func handler() ([]byte, error) {
	return service()
}

func service() ([]byte, error) {
	return dao()
}

func dao() ([]byte, error) {

	json, err := json.Marshal(math.Inf(1))
	if err != nil {
		//return nil, err
		return nil, errors.Wrap(err, "解析错误")
	}
	return json, nil
}
