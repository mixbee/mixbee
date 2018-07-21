

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
	"fmt"
	"reflect"

	"github.com/gogo/protobuf/proto"
)

type protoSerializer struct{}

func newProtoSerializer() Serializer {
	return &protoSerializer{}
}

func (protoSerializer) Serialize(msg interface{}) ([]byte, error) {
	if message, ok := msg.(proto.Message); ok {
		bytes, err := proto.Marshal(message)
		if err != nil {
			return nil, err
		}

		return bytes, nil
	}
	return nil, fmt.Errorf("msg must be proto.Message")
}

func (protoSerializer) Deserialize(typeName string, bytes []byte) (interface{}, error) {
	protoType := proto.MessageType(typeName)
	if protoType == nil {
		return nil, fmt.Errorf("Unknown message type %v", typeName)
	}
	t := protoType.Elem()

	intPtr := reflect.New(t)
	instance := intPtr.Interface().(proto.Message)
	proto.Unmarshal(bytes, instance)

	return instance, nil
}

func (protoSerializer) GetTypeName(msg interface{}) (string, error) {
	if message, ok := msg.(proto.Message); ok {
		typeName := proto.MessageName(message)

		return typeName, nil
	}
	return "", fmt.Errorf("msg must be proto.Message")
}
