package browserkubeutil

import (
	"context"
	"testing"
)

func TestSeleniumUP(t *testing.T) {
	type args struct {
		ctx context.Context
		u   string
	}
	type testCase struct {
		name    string
		args    args
		wantErr bool
	}
	tests := []testCase{
		{
			name: "empty",
			args: args{
				ctx: context.Background(),
				u:   "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SeleniumUP(tt.args.ctx, tt.args.u); (err != nil) != tt.wantErr {
				t.Errorf("SeleniumUP() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
