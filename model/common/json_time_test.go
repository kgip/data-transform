package common

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestJSONTime_MarshalJSON(t *testing.T) {
	now := JSONTime(time.Now())
	bytes, _ := json.Marshal(now)
	fmt.Println(string(bytes))
}

func TestJSONTime_UnmarshalJSON(t *testing.T) {
	str_time := `"2018-10-23"`
	jt := &JSONTime{}
	json.Unmarshal([]byte(str_time), jt)
	fmt.Println(jt)
}

func TestAssign(t *testing.T) {
	tn := time.Now()
	jt := JSONTime(tn.Add(time.Hour))
	rt1 := reflect.TypeOf(tn)
	rt2 := reflect.TypeOf(jt)
	fmt.Println(rt2.AssignableTo(rt1))
}
