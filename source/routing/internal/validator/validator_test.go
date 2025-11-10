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

	minAltitude := models.MustNewAltitude(50, models.MT)
	maxAltitude := models.MustNewAltitude(300, models.MT)
	outsideAltitude := models.MustNewAltitude(400, models.MT)
	outsideAltitudeMax := models.MustNewAltitude(500, models.MT)
	insideAltitude := models.MustNewAltitude(100, models.MT)

	sv.SetAltitude(minAltitude, maxAltitude)
	
	w_list_outside_2d := []*models.Waypoint{
		models.MustNewWaypoint(12, 50.895584255438735, 4.4425424745192, insideAltitude),
		models.MustNewWaypoint(13, 50.895263589374, 4.460836941507978, insideAltitude),
	}

	w_list_outside_3d := make([]*models.Waypoint, 0, len(w_list))
	for i, w := range w_list {
		w_list_outside_3d = append(w_list_outside_3d, models.MustNewWaypoint(i, w.Lat, w.Lon, w.Alt))
	}
	w_list_outside_3d[0].SetAltitude(outsideAltitude)
	w_list_outside_3d[1].SetAltitude(outsideAltitude)

	c_list_outside_2d := []*models.Feature3D{
		models.MustNewFeatureFromGeojson(`{
            "id": "9ZKayQc49pYwXFRXKNODi82xGtkStFDh",
            "type": "Feature",
            "properties": {},
            "geometry": {
                "coordinates": [
                    [
                        [
                            4.466753474827641,
                            50.896283325477214
                        ],
                        [
                            4.466682294848738,
                            50.9005694447811
                        ],
                        [
                            4.477525881225631,
                            50.901800827558134
                        ],
                        [
                            4.474624763394445,
                            50.89477663215163
                        ],
                        [
                            4.466753474827641,
                            50.896283325477214
                        ]
                    ]
                ],
                "type": "Polygon"
            }
        }`),
		models.MustNewFeatureFromGeojson(`{
            "id": "Qjtg0lO55NEks55T9SxhtFKcEyV0hCAO",
            "type": "Feature",
            "properties": {},
            "geometry": {
                "coordinates": [
                    [
                        [
                            4.451213420713714,
                            50.89541782155524
                        ],
                        [
                            4.446371675091967,
                            50.89970124216393
                        ],
                        [
                            4.460684775516739,
                            50.903263622243855
                        ],
                        [
                            4.461621997646006,
                            50.89710421848892
                        ],
                        [
                            4.451213420713714,
                            50.89541782155524
                        ]
                    ]
                ],
                "type": "Polygon"
            }
        }`),
	}

	c_on_border := models.MustNewFeatureFromGeojson(`{
            "id": "bP19jhZIftCeogi6xkBykEH056TCsdlN",
            "type": "Feature",
            "properties": {},
            "geometry": {
                "coordinates": [
                    [
                        [
                            4.458586426221046,
                            50.895235380947724
                        ],
                        [
                            4.467764640657151,
                            50.89240797622449
                        ],
                        [
                            4.465019634851103,
                            50.89692221138691
                        ],
                        [
                            4.458586426221046,
                            50.895235380947724
                        ]
                    ]
                ],
                "type": "Polygon"
            }
        }`)

	c_list_outside_3d := make([]*models.Feature3D, 0, len(c_list))
	for _, c := range c_list {
		c_list_outside_3d = append(c_list_outside_3d, models.MustNewFeatureFromGeojsonFeature(c.Feature))
	}
	c_list_outside_3d[0].SetAltitude(outsideAltitude, outsideAltitudeMax)
	c_list_outside_3d[1].SetAltitude(outsideAltitude, outsideAltitudeMax)


	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		searchVolume *models.Feature3D
		waypoints    []*models.Waypoint
		constraints  []*models.Feature3D
		valWaypoints         []*models.Waypoint
		valConstraints        []*models.Feature3D
		wantErr      bool
	}{
		{name: "Val - All wps and obstacles in sv", searchVolume: sv, waypoints: w_list, constraints: c_list, valWaypoints: w_list, valConstraints: c_list, wantErr: false},
		{name: "Val - Some wps not in sv", searchVolume: sv, waypoints: append(w_list, w_list_outside_2d...), constraints: c_list, valWaypoints: w_list, valConstraints: c_list, wantErr: false},
		{name: "Val - Some wps not in sv (altitude)", searchVolume: sv, waypoints: w_list_outside_3d, constraints: c_list, valWaypoints: w_list_outside_3d[2:], valConstraints: c_list, wantErr: false},
		{name: "Val - Some obstacles not in sv", searchVolume: sv, waypoints: w_list, constraints: append(c_list, c_list_outside_2d...), valWaypoints: w_list, valConstraints: c_list, wantErr: false},
		{name: "Val - Some obstacles not in sv (altitude)", searchVolume: sv, waypoints: w_list, constraints: c_list_outside_3d, valWaypoints: w_list, valConstraints: c_list_outside_3d[2:], wantErr: false},
		{name: "Val - One obstacle on border", searchVolume: sv, waypoints: w_list, constraints: append(c_list, c_on_border), valWaypoints: w_list, valConstraints: append(c_list, c_on_border), wantErr: false},
		// not throw an error in this case, the service will handle this
		{name: "Val - Less than 2 wps in sv", searchVolume: sv, waypoints: append(w_list_outside_2d, w_list[0]), constraints: c_list, valWaypoints: []*models.Waypoint{w_list[0]}, valConstraints: c_list, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := validator.NewDefaultValidator()
			valWaypoints, valConstraints, gotErr := v.ValidateInput(tt.searchVolume, tt.waypoints, tt.constraints)
			
			// TODO: Just for testing, put the validated wps in another color
            utils.MarkWaypointsAsInsideSearchVolume(tt.waypoints, valWaypoints...)
            utils.MarkConstraintsAsInsideSearchVolume(tt.constraints, valConstraints...)
			utils.ExportToGeoJSON("validator", tt.waypoints, tt.constraints, tt.name, false, tt.searchVolume.Feature)
			
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
			if !assert.ElementsMatch(t, valWaypoints, tt.valWaypoints) {
				t.Errorf("ValidateInput() = %v, want %v", valWaypoints, tt.valWaypoints)
			}
			if !assert.ElementsMatch(t, valConstraints, tt.valConstraints) {
				t.Errorf("ValidateInput() = %v, want %v", valConstraints, tt.valConstraints)
			}
		})
	}
}
