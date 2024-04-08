package auth_test

import (
	"reflect"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/remiges-tech/crux/db/sqlc-gen"
	"github.com/remiges-tech/crux/server/capability"
)

func TestFilterApp(t *testing.T) {
	type args struct {
		appCapDb []sqlc.GetUserCapsAndAppsByRealmRow
		appcaps  []string
		apps     *[]string
	}
	tests := []struct {
		name string
		args args
		want map[string][]string
	}{
		{
			name: "Test case 1 - no entry present in db",
			args: args{
				appCapDb: []sqlc.GetUserCapsAndAppsByRealmRow{},
				appcaps:  []string{"schema", "rules"},
				apps:     &[]string{"nedbank1", "retailbank"},
			},
			want: map[string][]string{
				"nedbank1":   {"schema", "rules"},
				"retailbank": {"schema", "rules"},
			},
		},
		{
			name: "Test case 2 - one entry is already present in db",
			args: args{
				appCapDb: []sqlc.GetUserCapsAndAppsByRealmRow{
					{Cap: "schema", App: pgtype.Text{String: "starmf", Valid: true}},
				},
				appcaps: []string{"schema", "rules"},
				apps:    &[]string{"nedbank1", "retailbank", "starmf"},
			},
			want: map[string][]string{
				"nedbank1":   {"schema", "rules"},
				"retailbank": {"schema", "rules"},
				"starmf":     {"rules"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := capability.FilterApp(tt.args.appCapDb, tt.args.appcaps, tt.args.apps); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FilterApp() = %v, want %v", got, tt.want)
			}
		})
	}
}
