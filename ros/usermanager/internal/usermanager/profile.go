package usermanager

import (
	"encoding/json"

	"github.com/9d4/netpilot/ros/board"
)

type Profile struct {
	Name          string  `json:"name"`
	Validity      string  `json:"validity"`
	Comment       string  `json:"comment"`
	Price         float64 `json:"price,string"`
	StartsWhen    string  `json:"starts-when"`
	OverrideUsers string  `json:"override-shared-users"`
}

func GetProfiles(b *board.Board) ([]*Profile, error) {
	res, err := b.Cli().R().Get(b.Url("/user-manager/profile"))
	if err != nil {
		return nil, err
	}

	var p []*Profile
	err = json.Unmarshal(res.Body(), &p)
	if err != nil {
		return nil, err
	}
	return p, nil
}
