package template

import (
	"bytes"
	"github.com/docking-tools/register/api"
	"log"
	"os"
)

// env returns the value of the environment variable set
func env(s string) (string, error) {
	return os.Getenv(s), nil
}

func convertGraphTopath(m api.Recmap) (result map[string]string) {
	log.Print(m)
	result = make(map[string]string, 0)
	for k, v := range m {
		var buffer bytes.Buffer
		buffer.WriteString("/")
		buffer.WriteString(k)
		key := buffer.String()
		switch v.(type) {
		case api.Recmap:
			child := convertGraphTopath(v.(api.Recmap))
			for kc, vc := range child {
				var bufferChild bytes.Buffer
				bufferChild.WriteString(key)
				bufferChild.WriteString(kc)
				result[bufferChild.String()] = vc

			}
		case string:
			result[key] = v.(string)

		}
	}
	return
}

func listPathfromGraph(m api.Recmap) (result []string) {
	result, _ = listPathfromGraphInt(m)
	return
}
func listPathfromGraphInt(m api.Recmap) (result []string, hasKey bool) {
	result = make([]string, 0)
	hasKey = false
	for k, v := range m {
		var buffer bytes.Buffer
		buffer.WriteString("/")
		buffer.WriteString(k)
		key := buffer.String()
		switch v.(type) {
		case api.Recmap:
			child, lastKey := listPathfromGraphInt(v.(api.Recmap))
			if lastKey {
				result = append(result, key)
			}
			for _, kc := range child {
				var bufferChild bytes.Buffer
				bufferChild.WriteString(key)
				bufferChild.WriteString(kc)
				result = append(result, bufferChild.String())
			}
		default:
			hasKey = true
		}
		//if hasKey {
		//	result = append(result, key)
		//}
	}
	return
}
