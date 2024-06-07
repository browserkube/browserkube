package browserimage

import "testing"

func TestParseImageType(t *testing.T) {
	type args struct {
		img string
	}
	tests := []struct {
		name    string
		args    args
		want    ImageType
		wantErr bool
	}{
		{
			name: "selenium",
			args: args{img: "selenium/standalone-chrome:117.0"},
			want: ImageTypeSelenium,
		},
		{
			name: "aerokube-selenoid",
			args: args{img: "selenoid/vnc_chrome:117.0"},
			want: ImageTypeSelenoid,
		},
		{
			name: "aerokube-playwright",
			args: args{img: "playwright/chrome:playwright-1.39.0"},
			want: ImageTypeAerokube,
		},
		{
			name: "aerokube-playwright",
			args: args{img: "quay.io/browser/playwright-firefox:playwright-1.39.0"},
			want: ImageTypeAerokube,
		},
		{
			name: "aerokube-chrome",
			args: args{img: "quay.io/browser/chrome:121"},
			want: ImageTypeAerokube,
		},
		{
			name: "aerokube-cdtp",
			args: args{img: "cdtp/chrome:93.0"},
			want: ImageTypeAerokube,
		},
		{
			name: "aerokube-selenoid",
			args: args{img: "selenoid/vnc_chrome:117.0"},
			want: ImageTypeSelenoid,
		},
		{
			name: "microsoft",
			args: args{img: "mcr.microsoft.com/playwright:v1.39.0-jammy"},
			want: ImageTypeMicrosoft,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseImageType(tt.args.img)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseImageType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseImageType() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestImageType_VncPass(t *testing.T) {
	tests := []struct {
		name string
		it   ImageType
		want string
	}{
		{
			name: "aerokube",
			want: "browserkube",
			it:   ImageTypeAerokube,
		},
		{
			name: "selenium",
			want: "browserkube",
			it:   ImageTypeSelenium,
		},
		{
			name: "aerokube",
			want: "selenoid",
			it:   ImageTypeSelenoid,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.it.VncPass(); got != tt.want {
				t.Errorf("VncPass() = %v, want %v", got, tt.want)
			}
		})
	}
}
