package pb

import (
	"errors"
	"github.com/golang/protobuf/proto"
)

var errWrongValueType = errors.New("protobuf: convert on wrong type value")

type ProtobufSerializer struct{}

//定义protobuf
func New() *ProtobufSerializer {
	return &ProtobufSerializer{}
}

//proto Encode
func (s *ProtobufSerializer) Encode(v interface{}) ([]byte, error) {
	pb, ok := v.(proto.Message)
	if !ok {
		return nil, errWrongValueType
	}
	return proto.Marshal(pb)
}

//proto decode
func (s *ProtobufSerializer) Decode(data []byte, v interface{}) error {
	pb, ok := v.(proto.Message)
	if !ok {
		return errWrongValueType
	}
	return proto.Unmarshal(data, pb)
}
