package channelrouter

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/modelbus/one-api-pro/model"
)

type PriorityRandomSelector struct{}

func (s *PriorityRandomSelector) Select(_ context.Context, candidates []*model.Channel, req *RouteRequest) (*model.Channel, error) {
	if len(candidates) == 0 {
		return nil, fmt.Errorf("no candidates available")
	}

	if req.IgnoreFirstPriority {
		highestPriority := candidates[0].GetPriority()
		startIdx := 0
		for i, ch := range candidates {
			if ch.GetPriority() != highestPriority {
				startIdx = i
				break
			}
		}
		if startIdx > 0 && startIdx < len(candidates) {
			return candidates[startIdx+rand.Intn(len(candidates)-startIdx)], nil
		}
	}

	endIdx := len(candidates)
	highestPriority := candidates[0].GetPriority()
	if highestPriority > 0 {
		for i, ch := range candidates {
			if ch.GetPriority() != highestPriority {
				endIdx = i
				break
			}
		}
	}
	idx := rand.Intn(endIdx)
	return candidates[idx], nil
}