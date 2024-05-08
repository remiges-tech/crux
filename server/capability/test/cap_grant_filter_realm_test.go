package auth_test

import (
	"reflect"
	"testing"

	"github.com/remiges-tech/crux/server/capability"
)

func Test_filterRealm(t *testing.T) {
	type args struct {
		realmCapDb []string
		realmcaps  []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "Test case 1",
			args: args{
				realmCapDb: []string{"root", "config", "auth"},
				realmcaps:  []string{"auth", "report"},
			},
			want: []string{"report"},
		},
		{
			name: "Test case 2",
			args: args{
				realmCapDb: []string{"report", "root"},
				realmcaps:  []string{"report", "root"},
			},
			want: []string{},
		},
		{
			name: "Test case 3",
			args: args{
				realmCapDb: []string{"auth", "config"},
				realmcaps:  []string{"auth", "config", "root"},
			},
			want: []string{"root"},
		},
		{
			name: "Test case 4",
			args: args{
				realmCapDb: []string{"report", "auth"},
				realmcaps:  []string{"root"},
			},
			want: []string{"root"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := capability.FilterRealm(tt.args.realmCapDb, tt.args.realmcaps); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("filterRealm() = %v, want %v", got, tt.want)
			}
		})
	}
}
