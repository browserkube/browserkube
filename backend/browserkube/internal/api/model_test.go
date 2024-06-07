package api

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimestamp_MarshalJSON(t *testing.T) {
	layout := "2006-01-02"
	ttime, err := time.Parse(layout, "2020-08-13")
	assert.NoError(t, err)

	tests := []struct {
		name    string
		t       Timestamp
		want    []byte
		wantErr bool
	}{
		{
			name:    "positive",
			wantErr: false,
			want:    []byte("1597276800000"),
			t:       Timestamp(ttime),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.t.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MarshalJSON() got = %v, want %v", string(got), string(tt.want))
			}
		})
	}
}
