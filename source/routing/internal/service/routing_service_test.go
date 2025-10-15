package service_test

import (
	"geopathplanner/routing/internal/models"
	"geopathplanner/routing/internal/service"
	"testing"
)

func TestRoutingService_HandleRoutingRequest(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		input *models.RoutingRequest
		want  *models.RoutingResponse
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rs := service.NewRoutingService()
			got := rs.HandleRoutingRequest(tt.input)
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("HandleRoutingRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}
