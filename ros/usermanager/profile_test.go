package usermanager

import (
	"reflect"
	"testing"
	"time"

	um "github.com/9d4/netpilot/ros/usermanager/internal/usermanager"
)

func Test_transformProfile(t *testing.T) {
	type args struct {
		src  *um.Profile
		dst  *Profile
		want *Profile
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "one",
			args: args{
				src: &um.Profile{
					Name:          "prof2h",
					Validity:      "2h",
					Comment:       "netpilot-prof2h",
					Price:         131.3,
					StartsWhen:    "first-auth",
					OverrideUsers: "unlimited",
				},
				dst: &Profile{},
				want: &Profile{
					Name:          "prof2h",
					Validity:      time.Duration(time.Hour * 2),
					Comment:       "netpilot-prof2h",
					Price:         131.3,
					StartsWhen:    StartsWhenFirstAuth,
					OverrideUsers: OverrideUsersUnlimited,
				},
			},
		},
		{
			name: "two",
			args: args{
				src: &um.Profile{
					Name:          "prof 10 days",
					Validity:      "1w3d", // this is the example value taken from ROS Rest API
					Comment:       "",
					Price:         100,
					StartsWhen:    "assigned",
					OverrideUsers: "off",
				},
				dst: &Profile{},
				want: &Profile{
					Name:          "prof 10 days",
					Validity:      time.Duration(time.Hour * 24 * 10),
					Comment:       "",
					Price:         100,
					StartsWhen:    StartsWhenAssigned,
					OverrideUsers: OverrideUsersOff,
				},
			},
		},
		{
			name: "three",
			args: args{
				src: &um.Profile{
					Name:          "prof 101d 10:00:00",
					Validity:      "14w3d10h",
					Comment:       "101d 10h",
					Price:         5000.12,
					StartsWhen:    "first-auth",
					OverrideUsers: "10",
				},
				dst: &Profile{},
				want: &Profile{
					Name:          "prof 101d 10:00:00",
					Validity:      time.Duration((time.Hour*24)*101 + (time.Hour * 10)),
					Comment:       "101d 10h",
					Price:         5000.12,
					StartsWhen:    StartsWhenFirstAuth,
					OverrideUsers: 10,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transformProfile(tt.args.src, tt.args.dst)
			if eq := reflect.DeepEqual(*tt.args.dst, *tt.args.want); !eq {
				t.Errorf("\ngot: %+v\nwant: %+v", tt.args.dst, tt.args.want)
			}
		})
	}
}

func Test_transformIntProfile(t *testing.T) {
	type args struct {
		src *Profile
		dst *um.Profile
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transformIntProfile(tt.args.src, tt.args.dst)
		})
	}
}
