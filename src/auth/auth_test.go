package auth

import (
	"testing"
)

func TestAuth_Check(t *testing.T) {
	{
		a := NewAuth("1247fbc0e373637ff9fb08f0c81722f3", "2e409476eac9c6a7461e684ec4933748", []byte(`{"timestamp":158894385094}`))
		err := a.Check()
		if err != nil {
			t.Fatal(err)
		}
	}
	{
		a := NewAuth("1247fbc0e373637ff9fb08f0c81722f3", "1e409476eac9c6a7461e684ec4933748", []byte(`{"timestamp":158894385094}`))
		err := a.Check()
		if err == nil {
			t.Fatal("should be fail to check")
		}
	}
	t.Log("success")
}
