package utils

import "testing"


var user *User

type User struct {
	UserName string `json:"username"`
	NickName string `json:"nickname"`
	Age      int
	Birthday string
	Sex      string
	Email    string
	Phone    string
}

func TestMain(m *testing.M) {

	user = &User{
		UserName: "user1",
		NickName: "测试",
		Age:      18,
		Birthday: "2008/8/8",
		Sex:      "男",
		Email:    "test@example.com",
		Phone:    "110",
	}


	m.Run()
}

func TestJson2Object(t *testing.T) {
	str := Object2Json(user)
	t.Log("str", str)

	obj := Json2Object(str, user)
	t.Log("obj", obj)
}


