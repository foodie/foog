package json

import (
	"encoding/json"
)

//定义空类型
type JsonSerializer struct{}

//新建一个类型
func New() *JsonSerializer {
	return &JsonSerializer{}
}

//json encode
func (this *JsonSerializer) Encode(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

//json decode
func (this *JsonSerializer) Decode(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}
