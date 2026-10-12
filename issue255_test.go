package mergo_test

import (
	"testing"

	"dario.cat/mergo"
)

type issue255Person struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name,omitempty"`
	Password  string `json:"-"`
	Age       int
}

func TestIssue255MapStructToMapWithTaggedKeys(t *testing.T) {
	person := issue255Person{
		FirstName: "Ada",
		LastName:  "Lovelace",
		Password:  "secret",
		Age:       36,
	}

	actual := map[string]interface{}{}
	if err := mergo.Map(&actual, person, mergo.WithMapKeyTag("json")); err != nil {
		t.Fatal(err)
	}

	if actual["first_name"] != "Ada" {
		t.Fatalf("expected first_name key to be mapped from tag, got %#v", actual)
	}
	if actual["last_name"] != "Lovelace" {
		t.Fatalf("expected tag name to be used, got %#v", actual)
	}
	if actual["password"] != "secret" {
		t.Fatalf("expected '-' tag to fall back to default field name, got %#v", actual)
	}
	if actual["age"] != 36 {
		t.Fatalf("expected untagged field to fall back to default field name, got %#v", actual)
	}
}
