package main

import (
	"os"
	"testing"
)

const (
	relative = "relative.txt"
	absolute = "/home/akshay/learning-go/unit_testing/absolute.txt"
)

//Single Test
func TestCreateFileSingle(t *testing.T) {

	got, err := createFile(relative)
	defer os.Remove(relative)
	if err != nil {
		//We are expecting nil error so we error out
		t.Fatalf("Expected nil error, Got : %v", err)
	}

	if got != "/home/akshay/learning-go/unit_testing/relative.txt" {
		t.Fatalf("Got : %v , Want : %v", got, "/home/akshay/learning-go/unit_testing/relative.txt")
	}
}

//table driven with subtests
func Test_createFile(t *testing.T) {
	type args struct {
		fileName string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"Empty Name", args{""}, "", true},
		{"Relative Path", args{relative}, "/home/akshay/learning-go/unit_testing/relative.txt", false},
		{"Absolute Path", args{absolute}, absolute, false},
		{"Invalid Path", args{"test/test/test"}, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer os.Remove(tt.args.fileName)
			got, err := createFile(tt.args.fileName)
			if (err != nil) != tt.wantErr {
				t.Errorf("createFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("createFile() = %v, want %v", got, tt.want)
			}
		})
	}
}
