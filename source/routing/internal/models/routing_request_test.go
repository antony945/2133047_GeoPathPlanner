package models_test

import (
	"geopathplanner/routing/internal/models"
	"geopathplanner/routing/internal/utils"
	"testing"
)

func TestNewRoutingRequestFromJson(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		jsonString string
		want       *models.RoutingRequest
		wantErr    bool
	}{
		// TODO: Add test cases.
		{
            name: "RR-NoConstraints",
            wantErr: false,
            jsonString: `{
                "request_id": "1",
                "waypoints": [
                    {
                        "type": "Feature",
                        "properties": {
                            "altitudeUnit": "mt",
                            "altitudeValue": 200
                        },
                        "geometry": {
                            "type": "Point",            
                            "coordinates": [
                                4.433724687935722,
                                50.872778105839274
                            ]
                        }
                    },
                    {
                        "type": "Feature",
                        "properties": {
                            "altitudeUnit": "mt",
                            "altitudeValue": 300
                        },
                        "geometry": {
                            "type": "Point",            
                            "coordinates": [
                                4.46992531620532,
                                50.884400404439646
                            ]
                        }
                    },
                    {
                        "type": "Feature",
                        "properties": {
                            "altitudeUnit": "mt",
                            "altitudeValue": 400
                        },
                        "geometry": {
                            "type": "Point",            
                            "coordinates": [
                                4.45503208121508,
                                50.890383059561145
                            ]
                        }
                    }
                ],
                "search_volume": {
                    "type": "Feature",
                    "properties": {},
                    "geometry": {
                        "type": "Polygon",
                        "coordinates": [
                            [
                            [
                                4.424480223895728,
                                50.89367115387381
                            ],
                            [
                                4.424480223895728,
                                50.867778999101745
                            ],
                            [
                                4.480653371960557,
                                50.867778999101745
                            ],
                            [
                                4.480653371960557,
                                50.89367115387381
                            ],
                            [
                                4.424480223895728,
                                50.89367115387381
                            ]
                            ]
                        ]
                    }
                },
                "parameters": {
                    "algorithm": "rrt",
                    "storage": "list",
                    "sampler": "uniform",
                    "seed": 10,
                    "max_iterations": 10000,
                    "max_step_size_mt": 20,
                    "goal_bias": 0.10
                }
            }`,
        },
        {
            name: "RR-NoSearchVolume",
            wantErr: true,
            jsonString: `{
                "request_id": "1",
                "waypoints": [
                    {
                        "type": "Feature",
                        "properties": {
                            "altitudeUnit": "mt",
                            "altitudeValue": 200
                        },
                        "geometry": {
                            "type": "Point",            
                            "coordinates": [
                                4.433724687935722,
                                50.872778105839274
                            ]
                        }
                    },
                    {
                        "type": "Feature",
                        "properties": {
                            "altitudeUnit": "mt",
                            "altitudeValue": 300
                        },
                        "geometry": {
                            "type": "Point",            
                            "coordinates": [
                                4.46992531620532,
                                50.884400404439646
                            ]
                        }
                    },
                    {
                        "type": "Feature",
                        "properties": {
                            "altitudeUnit": "mt",
                            "altitudeValue": 400
                        },
                        "geometry": {
                            "type": "Point",            
                            "coordinates": [
                                4.45503208121508,
                                50.890383059561145
                            ]
                        }
                    }
                ],
                "constraints": [],
                "search_volume": {},
                "parameters": {
                    "algorithm": "rrtstar",
                    "storage": "rtree",
                    "sampler": "uniform",
                    "seed": 10,
                    "max_iterations": 10000,
                    "max_step_size_mt": 20,
                    "goal_bias": 0.10
                }
            }`,
        },
        {
            name: "RR-NoWaypoint",
            wantErr: true,
            jsonString: `{
                "request_id": "1",
                "waypoints": [
                    {
                        "lat": ciao
                    },
                    {
                        "type": "Feature",
                        "properties": {
                            "altitudeUnit": "mt",
                            "altitudeValue": 300
                        },
                        "geometry": {
                            "type": "Point",            
                            "coordinates": [
                                4.46992531620532,
                                50.884400404439646
                            ]
                        }
                    },
                    {
                        "type": "Feature",
                        "properties": {
                            "altitudeUnit": "mt",
                            "altitudeValue": 400
                        },
                        "geometry": {
                            "type": "Point",            
                            "coordinates": [
                                4.45503208121508,
                                50.890383059561145
                            ]
                        }
                    }
                ],
                "constraints": [
                    {
                        "id": 0,
                        "type": "Feature",
                        "properties": {
                            "altitudeUnit": "mt",
                            "maxAltitudeValue": 400,
                            "minAltitudeValue": 500
                        },
                        "geometry": {
                            "type": "Polygon",
                            "coordinates": [
                            [
                            [
                                4.452166432497052,
                                50.87565312347229
                            ],
                            [
                                4.45389986799907,
                                50.874291381256256
                            ],
                            [
                                4.4498672299481825,
                                50.88096516631941
                            ],
                            [
                                4.451848154983281,
                                50.8813222592305
                            ],
                            [
                                4.45478419822345,
                                50.87636729027514
                            ],
                            [
                                4.4560577382486315,
                                50.87672438629892
                            ],
                            [
                                4.451812846246014,
                                50.883420200613955
                            ],
                            [
                                4.447815509645977,
                                50.88214806230394
                            ],
                            [
                                4.452166432497052,
                                50.87565312347229
                            ]
                            ]
                            ]
                        }
                    },
                    {
                        "id": 1,
                        "type": "Feature",
                        "properties": {
                            "altitudeUnit": "ft",
                            "maxAltitudeValue": 100,
                            "minAltitudeValue": 1000
                        },
                        "geometry": {
                            "type": "Polygon",
                            "coordinates": [
                                [
                                [
                                    4.435823054794525,
                                    50.87917754349243
                                ],
                                [
                                    4.435999901698551,
                                    50.876186530052024
                                ],
                                [
                                    4.443605337154025,
                                    50.878195458959425
                                ],
                                [
                                    4.439678842862065,
                                    50.88446720530686
                                ],
                                [
                                    4.435823054794525,
                                    50.87917754349243
                                ]
                                ]
                            ]
                        }
                    },
                    {
                        "id": 2,
                        "type": "Feature",
                        "properties": {},
                        "geometry": {
                            "type": "Polygon",
                            "coordinates": [
                            [
                                [
                                4.443180802952469,
                                50.87710169486769
                                ],
                                [
                                4.445940351819587,
                                50.874043594523414
                                ],
                                [
                                4.449867003560769,
                                50.875472226326934
                                ],
                                [
                                4.447425915923816,
                                50.88058383172759
                                ],
                                [
                                4.443180802952469,
                                50.87710169486769
                                ]
                            ]
                            ]
                        }
                    }
                ],
                "search_volume": {
                    "type": "Feature",
                    "properties": {},
                    "geometry": {
                        "type": "Polygon",
                        "coordinates": [
                            [
                            [
                                4.424480223895728,
                                50.89367115387381
                            ],
                            [
                                4.424480223895728,
                                50.867778999101745
                            ],
                            [
                                4.480653371960557,
                                50.867778999101745
                            ],
                            [
                                4.480653371960557,
                                50.89367115387381
                            ],
                            [
                                4.424480223895728,
                                50.89367115387381
                            ]
                            ]
                        ]
                    }
                },
                "parameters": {
                    "algorithm": "rrt",
                    "storage": "list",
                    "sampler": "uniform",
                    "seed": 10,
                    "max_iterations": 10000,
                    "max_step_size_mt": 20,
                    "goal_bias": 0.10
                }
            }`,
        },
        {
            name: "RR-WrongWaypoint",
            wantErr: true,
            jsonString: `{
                "request_id": "1",
                "waypoints": [
                    {
                        "lat": ciao
                    },
                    {
                        "type": "Feature",
                        "properties": {
                            "altitudeUnit": "mt",
                            "altitudeValue": 300
                        },
                        "geometry": {
                            "type": "Point",            
                            "coordinates": [
                                4.46992531620532,
                                50.884400404439646
                            ]
                        }
                    },
                    {
                        "type": "Feature",
                        "properties": {
                            "altitudeUnit": "mt",
                            "altitudeValue": 400
                        },
                        "geometry": {
                            "type": "Point",            
                            "coordinates": [
                                4.45503208121508,
                                50.890383059561145
                            ]
                        }
                    }
                ],
                "constraints": [
                    {
                        "id": 0,
                        "type": "Feature",
                        "properties": {
                            "altitudeUnit": "mt",
                            "maxAltitudeValue": 400,
                            "minAltitudeValue": 500
                        },
                        "geometry": {
                            "type": "Polygon",
                            "coordinates": [
                            [
                            [
                                4.452166432497052,
                                50.87565312347229
                            ],
                            [
                                4.45389986799907,
                                50.874291381256256
                            ],
                            [
                                4.4498672299481825,
                                50.88096516631941
                            ],
                            [
                                4.451848154983281,
                                50.8813222592305
                            ],
                            [
                                4.45478419822345,
                                50.87636729027514
                            ],
                            [
                                4.4560577382486315,
                                50.87672438629892
                            ],
                            [
                                4.451812846246014,
                                50.883420200613955
                            ],
                            [
                                4.447815509645977,
                                50.88214806230394
                            ],
                            [
                                4.452166432497052,
                                50.87565312347229
                            ]
                            ]
                            ]
                        }
                    },
                    {
                        "id": 1,
                        "type": "Feature",
                        "properties": {
                            "altitudeUnit": "ft",
                            "maxAltitudeValue": 100,
                            "minAltitudeValue": 1000
                        },
                        "geometry": {
                            "type": "Polygon",
                            "coordinates": [
                                [
                                [
                                    4.435823054794525,
                                    50.87917754349243
                                ],
                                [
                                    4.435999901698551,
                                    50.876186530052024
                                ],
                                [
                                    4.443605337154025,
                                    50.878195458959425
                                ],
                                [
                                    4.439678842862065,
                                    50.88446720530686
                                ],
                                [
                                    4.435823054794525,
                                    50.87917754349243
                                ]
                                ]
                            ]
                        }
                    },
                    {
                        "id": 2,
                        "type": "Feature",
                        "properties": {},
                        "geometry": {
                            "type": "Polygon",
                            "coordinates": [
                            [
                                [
                                4.443180802952469,
                                50.87710169486769
                                ],
                                [
                                4.445940351819587,
                                50.874043594523414
                                ],
                                [
                                4.449867003560769,
                                50.875472226326934
                                ],
                                [
                                4.447425915923816,
                                50.88058383172759
                                ],
                                [
                                4.443180802952469,
                                50.87710169486769
                                ]
                            ]
                            ]
                        }
                    }
                ],
                "search_volume": {
                    "type": "Feature",
                    "properties": {},
                    "geometry": {
                        "type": "Polygon",
                        "coordinates": [
                            [
                            [
                                4.424480223895728,
                                50.89367115387381
                            ],
                            [
                                4.424480223895728,
                                50.867778999101745
                            ],
                            [
                                4.480653371960557,
                                50.867778999101745
                            ],
                            [
                                4.480653371960557,
                                50.89367115387381
                            ],
                            [
                                4.424480223895728,
                                50.89367115387381
                            ]
                            ]
                        ]
                    }
                },
                "parameters": {
                    "algorithm": "rrt",
                    "storage": "list",
                    "sampler": "uniform",
                    "seed": 10,
                    "max_iterations": 10000,
                    "max_step_size_mt": 20,
                    "goal_bias": 0.10
                }
            }`,
        },
        {
            name: "RR-WrongConstraint",
            wantErr: true,
            jsonString: `{
                "request_id": "1",
                "waypoints": [
                    {
                        "type": "Feature",
                        "properties": {
                            "altitudeUnit": "mt",
                            "altitudeValue": 300
                        },
                        "geometry": {
                            "type": "Point",            
                            "coordinates": [
                                4.46992531620532,
                                50.884400404439646
                            ]
                        }
                    },
                    {
                        "type": "Feature",
                        "properties": {
                            "altitudeUnit": "mt",
                            "altitudeValue": 400
                        },
                        "geometry": {
                            "type": "Point",            
                            "coordinates": [
                                4.45503208121508,
                                50.890383059561145
                            ]
                        }
                    }
                ],
                "constraints": [
                    {
                        "type": "lautaro",
                    },
                    {
                        "id": 1,
                        "type": "Feature",
                        "properties": {
                            "altitudeUnit": "ft",
                            "maxAltitudeValue": 100,
                            "minAltitudeValue": 1000
                        },
                        "geometry": {
                            "type": "Polygon",
                            "coordinates": [
                                [
                                [
                                    4.435823054794525,
                                    50.87917754349243
                                ],
                                [
                                    4.435999901698551,
                                    50.876186530052024
                                ],
                                [
                                    4.443605337154025,
                                    50.878195458959425
                                ],
                                [
                                    4.439678842862065,
                                    50.88446720530686
                                ],
                                [
                                    4.435823054794525,
                                    50.87917754349243
                                ]
                                ]
                            ]
                        }
                    },
                    {
                        "id": 2,
                        "type": "Feature",
                        "properties": {},
                        "geometry": {
                            "type": "Polygon",
                            "coordinates": [
                            [
                                [
                                4.443180802952469,
                                50.87710169486769
                                ],
                                [
                                4.445940351819587,
                                50.874043594523414
                                ],
                                [
                                4.449867003560769,
                                50.875472226326934
                                ],
                                [
                                4.447425915923816,
                                50.88058383172759
                                ],
                                [
                                4.443180802952469,
                                50.87710169486769
                                ]
                            ]
                            ]
                        }
                    }
                ],
                "search_volume": {
                    "type": "Feature",
                    "properties": {},
                    "geometry": {
                        "type": "Polygon",
                        "coordinates": [
                            [
                            [
                                4.424480223895728,
                                50.89367115387381
                            ],
                            [
                                4.424480223895728,
                                50.867778999101745
                            ],
                            [
                                4.480653371960557,
                                50.867778999101745
                            ],
                            [
                                4.480653371960557,
                                50.89367115387381
                            ],
                            [
                                4.424480223895728,
                                50.89367115387381
                            ]
                            ]
                        ]
                    }
                },
                "parameters": {
                    "algorithm": "rrt",
                    "storage": "list",
                    "sampler": "uniform",
                    "seed": 10,
                    "max_iterations": 10000,
                    "max_step_size_mt": 20,
                    "goal_bias": 0.10
                }
            }`,
        },
        {
            name: "RR-WrongParameters",
            wantErr: false,
            jsonString: `{
                "request_id": "1",
                "waypoints": [
                    {
                        "type": "Feature",
                        "properties": {
                            "altitudeUnit": "mt",
                            "altitudeValue": 300
                        },
                        "geometry": {
                            "type": "Point",            
                            "coordinates": [
                                4.46992531620532,
                                50.884400404439646
                            ]
                        }
                    },
                    {
                        "type": "Feature",
                        "properties": {
                            "altitudeUnit": "mt",
                            "altitudeValue": 400
                        },
                        "geometry": {
                            "type": "Point",            
                            "coordinates": [
                                4.45503208121508,
                                50.890383059561145
                            ]
                        }
                    }
                ],
                "constraints": [
                    {
                        "id": 1,
                        "type": "Feature",
                        "properties": {
                            "altitudeUnit": "ft",
                            "maxAltitudeValue": 100,
                            "minAltitudeValue": 1000
                        },
                        "geometry": {
                            "type": "Polygon",
                            "coordinates": [
                                [
                                [
                                    4.435823054794525,
                                    50.87917754349243
                                ],
                                [
                                    4.435999901698551,
                                    50.876186530052024
                                ],
                                [
                                    4.443605337154025,
                                    50.878195458959425
                                ],
                                [
                                    4.439678842862065,
                                    50.88446720530686
                                ],
                                [
                                    4.435823054794525,
                                    50.87917754349243
                                ]
                                ]
                            ]
                        }
                    },
                    {
                        "id": 2,
                        "type": "Feature",
                        "properties": {},
                        "geometry": {
                            "type": "Polygon",
                            "coordinates": [
                            [
                                [
                                4.443180802952469,
                                50.87710169486769
                                ],
                                [
                                4.445940351819587,
                                50.874043594523414
                                ],
                                [
                                4.449867003560769,
                                50.875472226326934
                                ],
                                [
                                4.447425915923816,
                                50.88058383172759
                                ],
                                [
                                4.443180802952469,
                                50.87710169486769
                                ]
                            ]
                            ]
                        }
                    }
                ],
                "search_volume": {
                    "type": "Feature",
                    "properties": {},
                    "geometry": {
                        "type": "Polygon",
                        "coordinates": [
                            [
                            [
                                4.424480223895728,
                                50.89367115387381
                            ],
                            [
                                4.424480223895728,
                                50.867778999101745
                            ],
                            [
                                4.480653371960557,
                                50.867778999101745
                            ],
                            [
                                4.480653371960557,
                                50.89367115387381
                            ],
                            [
                                4.424480223895728,
                                50.89367115387381
                            ]
                            ]
                        ]
                    }
                },
                "parameters": {
                    "algorithm": "bbbb",
                    "storage": "aaaa",
                    "sampler": "dddddd",
                    "seed": 10,
                    "max_iterations": 10000,
                    "max_step_size_mt": 20,
                    "goal_bias": 0.10
                }
            }`,
        },
        {
            name: "RR-Complete",
            wantErr: false,
            jsonString: `{
                "request_id": "1",
                "waypoints": [
                    {
                        "type": "Feature",
                        "properties": {
                            "altitudeUnit": "mt",
                            "altitudeValue": 200
                        },
                        "geometry": {
                            "type": "Point",            
                            "coordinates": [
                                4.433724687935722,
                                50.872778105839274
                            ]
                        }
                    },
                    {
                        "type": "Feature",
                        "properties": {
                            "altitudeUnit": "mt",
                            "altitudeValue": 300
                        },
                        "geometry": {
                            "type": "Point",            
                            "coordinates": [
                                4.46992531620532,
                                50.884400404439646
                            ]
                        }
                    },
                    {
                        "type": "Feature",
                        "properties": {
                            "altitudeUnit": "mt",
                            "altitudeValue": 400
                        },
                        "geometry": {
                            "type": "Point",            
                            "coordinates": [
                                4.45503208121508,
                                50.890383059561145
                            ]
                        }
                    }
                ],
                "constraints": [
                {
                    "id": 0,
                    "type": "Feature",
                    "properties": {
                        "altitudeUnit": "mt",
                        "maxAltitudeValue": 400,
                        "minAltitudeValue": 500
                    },
                    "geometry": {
                        "type": "Polygon",
                        "coordinates": [
                        [
                        [
                            4.452166432497052,
                            50.87565312347229
                        ],
                        [
                            4.45389986799907,
                            50.874291381256256
                        ],
                        [
                            4.4498672299481825,
                            50.88096516631941
                        ],
                        [
                            4.451848154983281,
                            50.8813222592305
                        ],
                        [
                            4.45478419822345,
                            50.87636729027514
                        ],
                        [
                            4.4560577382486315,
                            50.87672438629892
                        ],
                        [
                            4.451812846246014,
                            50.883420200613955
                        ],
                        [
                            4.447815509645977,
                            50.88214806230394
                        ],
                        [
                            4.452166432497052,
                            50.87565312347229
                        ]
                        ]
                        ]
                    }
                },
                {
                    "id": 1,
                    "type": "Feature",
                    "properties": {
                        "altitudeUnit": "ft",
                        "maxAltitudeValue": 100,
                        "minAltitudeValue": 1000
                    },
                    "geometry": {
                        "type": "Polygon",
                        "coordinates": [
                            [
                            [
                                4.435823054794525,
                                50.87917754349243
                            ],
                            [
                                4.435999901698551,
                                50.876186530052024
                            ],
                            [
                                4.443605337154025,
                                50.878195458959425
                            ],
                            [
                                4.439678842862065,
                                50.88446720530686
                            ],
                            [
                                4.435823054794525,
                                50.87917754349243
                            ]
                            ]
                        ]
                    }
                },
                {
                    "id": 2,
                    "type": "Feature",
                    "properties": {},
                    "geometry": {
                        "type": "Polygon",
                        "coordinates": [
                        [
                            [
                            4.443180802952469,
                            50.87710169486769
                            ],
                            [
                            4.445940351819587,
                            50.874043594523414
                            ],
                            [
                            4.449867003560769,
                            50.875472226326934
                            ],
                            [
                            4.447425915923816,
                            50.88058383172759
                            ],
                            [
                            4.443180802952469,
                            50.87710169486769
                            ]
                        ]
                        ]
                    }
                }
            ],
                "search_volume": {
                    "type": "Feature",
                    "properties": {},
                    "geometry": {
                        "type": "Polygon",
                        "coordinates": [
                            [
                            [
                                4.424480223895728,
                                50.89367115387381
                            ],
                            [
                                4.424480223895728,
                                50.867778999101745
                            ],
                            [
                                4.480653371960557,
                                50.867778999101745
                            ],
                            [
                                4.480653371960557,
                                50.89367115387381
                            ],
                            [
                                4.424480223895728,
                                50.89367115387381
                            ]
                            ]
                        ]
                    }
                },
                "parameters": {
                    "algorithm": "rrtstar",
                    "storage": "rtree",
                    "sampler": "uniform",
                    "seed": 10,
                    "max_iterations": 10000,
                    "max_step_size_mt": 20,
                    "goal_bias": 0.10
                },
                "received_at": "2025-11-01T10:40:13Z"
            }`,
        },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := models.NewRoutingRequestFromJson(tt.jsonString)
			utils.ExportToJSON(got, "models", tt.name, false)			
			
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("NewRoutingRequestFromJson() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("NewRoutingRequestFromJson() succeeded unexpectedly")
			}
		})
	}
}
