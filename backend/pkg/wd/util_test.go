package wd

import "testing"

func Test_parseSessionID(t *testing.T) {
	type args struct {
		uPath string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantCmd string
		wantErr bool
	}{
		{
			args:    args{uPath: "/wd/hub/session/session-id/element/element-id/selected"},
			name:    "positive",
			want:    "session-id",
			wantCmd: "/element/element-id/selected",
			wantErr: false,
		},
		{
			args:    args{uPath: "/wd/hub/session/session-id"},
			name:    "delete session",
			want:    "session-id",
			wantCmd: "",
			wantErr: false,
		},
		{
			args:    args{uPath: "/wd/hub/session/session-id/"},
			name:    "delete session 2",
			want:    "session-id",
			wantCmd: "/",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := ParseSessionPath(tt.args.uPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseSessionID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parseSessionID() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.wantCmd {
				t.Errorf("parseSessionID() got1 = %v, want %v", got1, tt.wantCmd)
			}
		})
	}
}

func TestReplaceSession(t *testing.T) {
	type args struct {
		path    string
		session string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "create session",
			args:    args{path: "/session", session: "xxx"},
			wantErr: true,
		},
		{
			name:    "delete session",
			args:    args{path: "/session/:sessionId", session: "xxx"},
			want:    "/session/xxx",
			wantErr: false,
		},
		{
			name:    "timeouts",
			args:    args{path: "/session/:sessionId/timeouts", session: "xxx"},
			want:    "/session/xxx/timeouts",
			wantErr: false,
		},
		{
			name:    "window size",
			args:    args{path: "/session/:sessionId/window/:windowHandle/size", session: "xxx"},
			want:    "/session/xxx/window/:windowHandle/size",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReplaceSession(tt.args.path, tt.args.session)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReplaceSession() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ReplaceSession() got = %v, want %v", got, tt.want)
			}
		})
	}
}
