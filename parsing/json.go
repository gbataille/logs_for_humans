package parsing

import "encoding/json"

func FromJsonGeneric(in []byte) (map[string]interface{}, error) {
	var jsonMap map[string](interface{})
	err := json.Unmarshal(in, &jsonMap)
	if err != nil {
		return nil, err
	}

	return jsonMap, nil
}
