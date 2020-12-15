package gomongo

import (
	"fmt"

	"gopkg.in/mgo.v2/bson"
)

// Parser type
type Parser struct{}

// Parse Method
// <out> type should have `bson` or `json` flag
func (P Parser) Parse(in interface{}, out *interface{}, toJSON bool) {
	var parsedPayload bson.M
	var isParsed bool

	// parse to bson.M if <string>
	switch in.(type) {
	case string:
		if isNotJSON := bson.UnmarshalJSON([]byte(in.(string)), &parsedPayload); isNotJSON != nil {
			if err := bson.Unmarshal([]byte(in.(string)), &parsedPayload); err != nil {
				fmt.Println(err)
			}
			isParsed = true
		}
		isParsed = true
	}

	var err error
	var jsonBytes []byte
	if toJSON {
		if isParsed {
			if jsonBytes, err = bson.MarshalJSON(parsedPayload); err != nil {
				fmt.Println(err)
			}
			*out = string(jsonBytes)
		} else {
			if jsonBytes, err = bson.MarshalJSON(in); err != nil {
				fmt.Println(err)
			}
			*out = string(jsonBytes)
		}

	} else if !toJSON {
		if isParsed {
			if jsonBytes, err = bson.Marshal(parsedPayload); err != nil {
				fmt.Println(err)
			}
		} else {
			if jsonBytes, err = bson.Marshal(in); err != nil {
				fmt.Println(err)
			}
		}

		if err := bson.Unmarshal(jsonBytes, &out); err != nil {
			fmt.Println(err)
		}
	}
}

// GrabOnKeyword Method
func (P Parser) GrabOnKeyword(in bson.M, out *bson.M, keyword string, toJSON bool) {


	roll := func (key string) {
		for key, value := range in {
			fmt.Println(key, value)
	
			if key == "string"
		}
	}	

}
