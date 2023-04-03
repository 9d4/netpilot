package prefix

import "testing"

func Test_keyPrefix_Get(t *testing.T) {
	type args struct {
		uuid       string
		attributes []string
	}
	tests := []struct {
		name string
		k    keyPrefix
		args args
		want string
	}{
		{
			name: "Get thing:uuid-00123:a:b",
			k:    "thing:",
			args: args{
				uuid:       "uuid-00123",
				attributes: []string{"a", "b"},
			},
			want: "thing:uuid-00123:a:b",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.k.Get(tt.args.uuid, tt.args.attributes...); got != tt.want {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_keyPrefix_GetStatus(t *testing.T) {
	type args struct {
		uuid string
	}
	tests := []struct {
		name string
		k    keyPrefix
		args args
		want string
	}{
		{
			name: "Status (uuid-100-000)",
			k:    "stuff:",
			args: args{uuid: "uuid-100-000"},
			want: "stuff:uuid-100-000:status",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.k.Status(tt.args.uuid); got != tt.want {
				t.Errorf("Status() = %v, want %v", got, tt.want)
			}
		})
	}
}
