package usermanager

import (
	"strconv"
	"time"

	um "github.com/9d4/netpilot/ros/usermanager/internal/usermanager"
	str2duration "github.com/xhit/go-str2duration/v2"
)

type StartsWhen string
type OverrideUsers int

const (
	StartsWhenFirstAuth StartsWhen = "first-auth"
	StartsWhenAssigned  StartsWhen = "assigned"
)
const (
	OverrideUsersInvalid   = -2
	OverrideUsersOff       = -1
	OverrideUsersUnlimited = 0
)

type Profile struct {
	Name          string
	Validity      time.Duration
	Comment       string
	Price         float64
	StartsWhen    StartsWhen
	OverrideUsers int
}

// func GetProfiles(b *board.Board) ([]*Profile, error) {
// 	profiles, err := um.GetProfiles(b)
// 	if err != nil {
// 		return nil, err
// 	}
// }

func transformProfile(src *um.Profile, dst *Profile) {
	dst.Name = src.Name
	dst.Comment = src.Comment
	dst.Price = src.Price

	switch StartsWhen(src.StartsWhen) {
	case StartsWhenAssigned:
		dst.StartsWhen = StartsWhenAssigned
	case StartsWhenFirstAuth:
		dst.StartsWhen = StartsWhenFirstAuth
	default:
		dst.StartsWhen = ""
	}

	switch src.OverrideUsers {
	case "off":
		dst.OverrideUsers = OverrideUsersOff
	case "unlimited":
		dst.OverrideUsers = OverrideUsersUnlimited
	default:
		num, err := strconv.Atoi(src.OverrideUsers)
		if err != nil {
			dst.OverrideUsers = OverrideUsersInvalid
			break
		}
		dst.OverrideUsers = num
	}

	dur, err := str2duration.ParseDuration(src.Validity)
	if err == nil {
		dst.Validity = dur
	}
}

func transformIntProfile(src *Profile, dst *um.Profile) {

}
