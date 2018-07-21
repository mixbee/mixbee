

/***************************************************
Copyright 2016 https://github.com/AsynkronIT/protoactor-go

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*****************************************************/
package zmqremote

import (
	"bytes"
	"fmt"
	"reflect"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/proto"
)

type jsonSerializer struct {
	jsonpb.Marshaler
	jsonpb.Unmarshaler
}

func newJsonSerializer() Serializer {
	return &jsonSerializer{
		Marshaler: jsonpb.Marshaler{},
		Unmarshaler: jsonpb.Unmarshaler{
			AllowUnknownFields: true,
		},
	}
}

func (j *jsonSerializer) Serialize(msg interface{}) ([]byte, error) {
	if message, ok := msg.(*JsonMessage); ok {
		return []byte(message.Json), nil
	} else if message, ok := msg.(proto.Message); ok {

		str, err := j.Marshaler.MarshalToString(message)
		if err != nil {
			return nil, err
		}

		return []byte(str), nil
	}
	return nil, fmt.Errorf("msg must be proto.Message")
}

func (j *jsonSerializer) Deserialize(typeName string, b []byte) (interface{}, error) {
	protoType := proto.MessageType(typeName)
	if protoType == nil {
		m := &JsonMessage{
			TypeName: typeName,
			Json:     string(b),
		}
		return m, nil
	}
	t := protoType.Elem()

	intPtr := reflect.New(t)
	instance, ok := intPtr.Interface().(proto.Message)
	if ok {
		r := bytes.NewReader(b)
		j.Unmarshaler.Unmarshal(r, instance)

		return instance, nil
	} else {
		return nil, fmt.Errorf("msg must be proto.Message")
	}
}

func (j *jsonSerializer) GetTypeName(msg interface{}) (string, error) {
	if message, ok := msg.(*JsonMessage); ok {
		return message.TypeName, nil
	} else if message, ok := msg.(proto.Message); ok {
		typeName := proto.MessageName(message)

		return typeName, nil
	}

	return "", fmt.Errorf("msg must be proto.Message")
}
