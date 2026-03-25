package coordinates

import (
	"encoding/json"
	"fmt"
	"math"
)

type Coordinate struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}


func (c *Coordinate) Print() {
	fmt.Printf("Latitude:  %f\n", c.Latitude)
	fmt.Printf("Longitude: %f\n", c.Longitude)
	fmt.Println()
}

func UnmarshalCoordinates(jsonString string) ([]Coordinate, error) {
	var coordinates []Coordinate
	err := json.Unmarshal([]byte(jsonString), &coordinates)
	if err != nil {
		fmt.Println("Error unmarshalling coordinates json.")
		return nil, err
	}
	return coordinates, nil
}

func Haversine(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadius = 6371.0
	lat1Rad := lat1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	deltaLat := (lat2 - lat1) * math.Pi / 180
	deltaLon := (lon2 - lon1) * math.Pi / 180

	a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(deltaLon/2)*math.Sin(deltaLon/2)

	return earthRadius * 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
}

type HaversineInput struct {
	Lat1 float64 `json:"lat1"`
	Lon1 float64 `json:"lon1"`
	Lat2 float64 `json:"lat2"`
	Lon2 float64 `json:"lon2"`
}
