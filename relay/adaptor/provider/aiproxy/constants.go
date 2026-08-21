package aiproxy

import "github.com/modelbus/one-api-pro/relay/adaptor/openai"

var ModelList = []string{""}

func init() {
	ModelList = openai.ModelList
}
