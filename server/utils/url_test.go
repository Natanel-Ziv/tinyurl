package utils

import "testing"

func TestGenerateRandomShortStringInSIze(t *testing.T) {
	type args struct {
		maxSize int
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{
			name: "Sanity",
			args: args{
				maxSize: 5,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateRandomShortStringInSIze(tt.args.maxSize)
			if err != nil {
				t.Errorf("GenerateRandomShortStringInSIze() = %+v", err)
			}
			t.Logf("GenerateRandomShortStringInSIze() = %v", got)
		})
	}
}
