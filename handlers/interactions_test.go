package handlers

import (
	"testing"
)

func TestCalculateScore(t *testing.T) {
	tests := []struct {
		name            string
		interactionType string
		value           float64
		want            float64
		wantErr         bool
	}{
		{
			name:            "view",
			interactionType: "view",
			value:           1,
			want:            1,
			wantErr:         false,
		},
		{
			name:            "like",
			interactionType: "like",
			value:           1,
			want:            2,
			wantErr:         false,
		},
		{
			name:            "invalid",
			interactionType: "invalid",
			value:           1,
			want:            0,
			wantErr:         true,
		},
		{
			name:            "watch_time",
			interactionType: "watch_time",
			value:           100,
			want:            50,
			wantErr:         false,
		},

		{
			name:            "comment",
			interactionType: "comment",
			value:           1,
			want:            3,
			wantErr:         false,
		},
	}
	for _, tt := range tests {
		score, err := calculateScore(tt.interactionType, tt.value)
		if (err != nil) != tt.wantErr {
			t.Errorf("%s: calculateScore() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if score != tt.want {
			t.Errorf("%s: calculateScore() = %v, want %v", tt.name, score, tt.want)
		}

	}
}
