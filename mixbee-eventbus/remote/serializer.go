

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
package remote

var DefaultSerializerID int32 = 0
var serializers []Serializer

func init() {
	RegisterSerializer(newProtoSerializer())
	RegisterSerializer(newJsonSerializer())
}

func RegisterSerializerAsDefault(serializer Serializer) {
	serializers = append(serializers, serializer)
	DefaultSerializerID = int32(len(serializers) - 1)
}

func RegisterSerializer(serializer Serializer) {
	serializers = append(serializers, serializer)
}

type Serializer interface {
	Serialize(msg interface{}) ([]byte, error)
	Deserialize(typeName string, bytes []byte) (interface{}, error)
	GetTypeName(msg interface{}) (string, error)
}

func Serialize(message interface{}, serializerID int32) ([]byte, string, error) {
	res, err := serializers[serializerID].Serialize(message)
	typeName, err := serializers[serializerID].GetTypeName(message)
	return res, typeName, err
}

func Deserialize(message []byte, typeName string, serializerID int32) (interface{}, error) {
	return serializers[serializerID].Deserialize(typeName, message)
}
