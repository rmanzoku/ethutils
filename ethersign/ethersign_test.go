package ethersign

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/cheekybits/is"
)

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

func initTesting(t *testing.T) is.I {
	is := is.New(t)
	return is
}

func print(in interface{}) {
	if reflect.TypeOf(in).Kind() == reflect.Struct {
		in, _ = json.Marshal(in)
		in = string(in.([]byte))
	}
	fmt.Println(in)
}
