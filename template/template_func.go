package template

import (
	"os"
	"github.com/docking-tools/register/api"
	"bytes"
)

// env returns the value of the environment variable set
func env(s string) (string, error) {
	return os.Getenv(s), nil
}

func convertGraphTopath(m api.Recmap) (result map[string]string) {
	result = make(map[string]string,0)
	for k, v := range m {
		var buffer bytes.Buffer
		buffer.WriteString("/")
		buffer.WriteString(k)
		key :=buffer.String()
		switch v.(type){
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