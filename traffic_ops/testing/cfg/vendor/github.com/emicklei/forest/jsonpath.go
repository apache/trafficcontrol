package forest

import (
	"net/http"
	"strconv"
	"strings"
)

// JSONPath returns the value found by following the dotted path in a JSON document hash.
// E.g .chapters.0.title in  { "chapters" : [{"title":"Go a long way"}] }
func JSONPath(t T, r *http.Response, dottedPath string) interface{} {
	var value interface{}
	ExpectJSONHash(t, r, func(doc map[string]interface{}) {
		value = pathFindIn(0, strings.Split(dottedPath, ".")[1:], doc)
	})
	return value
}

// JSONArrayPath returns the value found by following the dotted path in a JSON array.
// E.g .1.title in  [ {"title":"Go a long way"}, {"title":"scary scala"} ]
func JSONArrayPath(t T, r *http.Response, dottedPath string) interface{} {
	var value interface{}
	ExpectJSONArray(t, r, func(list []interface{}) {
		value = pathFindIn(0, strings.Split(dottedPath, ".")[1:], list)
	})
	return value
}

func pathFindIn(index int, tokens []string, here interface{}) interface{} {
	//.Printf("%d %q %d, %v\n", index, tokens, len(tokens), here)
	if here == nil {
		return here
	}
	if index == len(tokens) {
		return here
	}
	token := tokens[index]
	if len(token) == 0 {
		return here
	}
	i, err := strconv.Atoi(token)
	if err == nil {
		// try index into array
		array, ok := here.([]interface{})
		if ok {
			if i >= len(array) {
				return nil
			}
			return pathFindIn(index+1, tokens, array[i])
		}
		return nil
	}
	// try key into hash
	hash, ok := here.(map[string]interface{})
	if ok {
		return pathFindIn(index+1, tokens, hash[token])
	}
	return nil
}
