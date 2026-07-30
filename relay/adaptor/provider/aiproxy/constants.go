package aiproxy

import "github.com/Leon-PanPan/one-api-pro/relay/adaptor/openai"

var ModelList = []string{""}

func init() {
	ModelList = openai.ModelList
}
