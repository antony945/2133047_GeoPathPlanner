package validator_test

import (
	"geopathplanner/routing/internal/models"
	"geopathplanner/routing/internal/utils"
	"geopathplanner/routing/internal/validator"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultValidator_ValidateInput(t *testing.T) {
	sv, w_list, c_list, _ := utils.SetupTestScenario()

	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		searchVolume *models.Feature3D
		waypoints    []*models.Waypoint
		constraints  []*models.Feature3D
		want         []*models.Waypoint
		want2        []*models.Feature3D
		wantErr      bool
	}{
		{searchVolume: sv, waypoints: w_list, constraints: c_list, want: w_list, want2: c_list, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := validator.NewDefaultValidator()
			got, got2, gotErr := v.ValidateInput(tt.searchVolume, tt.waypoints, tt.constraints)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("ValidateInput() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("ValidateInput() succeeded unexpectedly")
			}
			// TODO: update the condition below to compare got with tt.want.
			if !assert.ElementsMatch(t, got, tt.want) {
				t.Errorf("ValidateInput() = %v, want %v", got, tt.want)
			}
			if !assert.ElementsMatch(t, got2, tt.want2) {
				t.Errorf("ValidateInput() = %v, want %v", got2, tt.want2)
			}
		})
	}
}
