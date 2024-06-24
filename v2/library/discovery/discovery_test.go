package discovery

import (
	"testing"

	"github.com/eduardooliveira/stLib/v2/library/entities"
)

func Test_discoverer_Run(t *testing.T) {
	type fields struct {
		root         string
		rootAsset    *entities.Asset
		currentAsset *entities.Asset
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Test 1",
			fields: fields{
				root: "/Users/eduardooliveira/mmp_dev/agent/testdata2",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &discoverer{
				root:         tt.fields.root,
				rootAsset:    tt.fields.rootAsset,
				currentAsset: tt.fields.currentAsset,
			}
			if err := d.Run(); (err != nil) != tt.wantErr {
				t.Errorf("discoverer.Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
