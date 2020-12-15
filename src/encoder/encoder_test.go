package encoder

import "testing"

func TestToSnakeCase(t *testing.T) {
  if ToSnakeCase("SomethingSimple") != "_something_simple" {
    t.Fail()
  }

  if ToSnakeCase("SessionToken") != "_session_token" {
    t.Fail()
  }

  if ToSnakeCase("AccessKeyId") != "_access_key_id" {
    t.Fail()
  }

  if ToSnakeCase("SecretAccessKey") != "_secret_access_key" {
    t.Fail()
  }

  if ToSnakeCase("") != "" {
    t.Fail()
  }
}
