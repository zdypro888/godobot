package draw

import (
	"os"

	"google.golang.org/protobuf/proto"
)

func LoadTrajectories(path string) (*Signature, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var trajectories Signature
	if err := proto.Unmarshal(data, &trajectories); err != nil {
		return nil, err
	}
	return &trajectories, nil
}
