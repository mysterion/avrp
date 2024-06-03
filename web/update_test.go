package web

import (
	"reflect"
	"testing"
)

func TestLatestRelease(t *testing.T) {
	tests := []struct {
		name    string
		want    Release
		wantErr bool
	}{
		// TODO: Add test cases.
		// r, err := web.LatestRelease()
		// if err != nil {
		// 	panic(err)
		// }
		// fmt.Println("Latest Release: ", r.Tag)

		// fmt.Println("All Releases: ")

		// rall, err := web.AllReleases()
		// if err != nil {
		// 	panic(err)
		// }

		// for _, release := range rall {
		// 	fmt.Println(release.Tag)
		// }
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LatestRelease()
			if (err != nil) != tt.wantErr {
				t.Errorf("LatestRelease() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LatestRelease() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAllReleases(t *testing.T) {
	tests := []struct {
		name    string
		want    []Release
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := AllReleases()
			if (err != nil) != tt.wantErr {
				t.Errorf("AllReleases() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AllReleases() = %v, want %v", got, tt.want)
			}
		})
	}
}
