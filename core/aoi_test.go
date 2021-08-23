package core

import (
	"fmt"
	"reflect"
	"testing"
)

func TestNewAOIManager(t *testing.T) {

	aoiMgr:= NewAOIManager(0,250,5,0,250,5)
	log.Println(aoiMgr)
	//type args struct {
	//	minX  int
	//	maxX  int
	//	cntsX int
	//	minY  int
	//	maxY  int
	//	cntsY int
	//}
	//tests := []struct {
	//	name string
	//	args args
	//	want *AOIManager
	//}{
	//	// TODO: Add test cases.
	//}
	//for _, tt := range tests {
	//	t.Run(tt.name, func(t *testing.T) {
	//		if got := NewAOIManager(tt.args.minX, tt.args.maxX, tt.args.cntsX, tt.args.minY, tt.args.maxY, tt.args.cntsY); !reflect.DeepEqual(got, tt.want) {
	//			t.Errorf("NewAOIManager() = %v, want %v", got, tt.want)
	//		}
	//	})
	//}
}

func TestAOIManager_String(t *testing.T) {
	type fields struct {
		MinX  int
		MaxX  int
		CntsX int
		MinY  int
		MaxY  int
		CntsY int
		grIDs map[int]*GrID
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &AOIManager{
				MinX:  tt.fields.MinX,
				MaxX:  tt.fields.MaxX,
				CntsX: tt.fields.CntsX,
				MinY:  tt.fields.MinY,
				MaxY:  tt.fields.MaxY,
				CntsY: tt.fields.CntsY,
				grIDs: tt.fields.grIDs,
			}
			if got := m.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAOIManager_grIDLength(t *testing.T) {
	type fields struct {
		MinX  int
		MaxX  int
		CntsX int
		MinY  int
		MaxY  int
		CntsY int
		grIDs map[int]*GrID
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &AOIManager{
				MinX:  tt.fields.MinX,
				MaxX:  tt.fields.MaxX,
				CntsX: tt.fields.CntsX,
				MinY:  tt.fields.MinY,
				MaxY:  tt.fields.MaxY,
				CntsY: tt.fields.CntsY,
				grIDs: tt.fields.grIDs,
			}
			if got := m.grIDLength(); got != tt.want {
				t.Errorf("grIDLength() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAOIManager_grIDWIDth(t *testing.T) {
	type fields struct {
		MinX  int
		MaxX  int
		CntsX int
		MinY  int
		MaxY  int
		CntsY int
		grIDs map[int]*GrID
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &AOIManager{
				MinX:  tt.fields.MinX,
				MaxX:  tt.fields.MaxX,
				CntsX: tt.fields.CntsX,
				MinY:  tt.fields.MinY,
				MaxY:  tt.fields.MaxY,
				CntsY: tt.fields.CntsY,
				grIDs: tt.fields.grIDs,
			}
			if got := m.grIDWIDth(); got != tt.want {
				t.Errorf("grIDWIDth() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewAOIManager1(t *testing.T) {
	type args struct {
		minX  int
		maxX  int
		cntsX int
		minY  int
		maxY  int
		cntsY int
	}
	tests := []struct {
		name string
		args args
		want *AOIManager
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewAOIManager(tt.args.minX, tt.args.maxX, tt.args.cntsX, tt.args.minY, tt.args.maxY, tt.args.cntsY); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAOIManager() = %v, want %v", got, tt.want)
			}
		})
	}
}
