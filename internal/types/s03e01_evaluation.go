package types

import (
	"fmt"
	"strings"
	"time"
)

type AnswerRechecksS03E01 struct {
	Recheck []string `json:"recheck"`
}

type AnswerS03E01 struct {
	Task   string               `json:"task"`
	ApiKey string               `json:"apikey"`
	Answer AnswerRechecksS03E01 `json:"answer"`
}

type EvaluationSensorReading struct {
	SensorType      string  `json:"sensor_type,omitempty"`
	Timestamp       int64   `json:"timestamp,omitempty"`
	TemperatureK    float64 `json:"temperature_K,omitempty"`
	PressureBar     float64 `json:"pressure_bar,omitempty"`
	WaterLevelM     float64 `json:"water_level_meters,omitempty"`
	VoltageSupplyV  float64 `json:"voltage_supply_v,omitempty"`
	HumidityPercent float64 `json:"humidity_percent,omitempty"`
	OperatorNotes   string  `json:"operator_notes,omitempty"`
}

// Value returns the numeric field matching the JSON name.
// ok is false for unknown or non-numeric fields.
func (r EvaluationSensorReading) Value(field string) (float64, bool) {
	switch field {
	case "temperature_K":
		return r.TemperatureK, true
	case "pressure_bar":
		return r.PressureBar, true
	case "water_level_meters":
		return r.WaterLevelM, true
	case "voltage_supply_v":
		return r.VoltageSupplyV, true
	case "humidity_percent":
		return r.HumidityPercent, true
	}
	return 0, false
}

func (r EvaluationSensorReading) String() string {
	var b strings.Builder

	// b.WriteString("── Sensor Reading ──────────────────────\n")

	// if r.SensorType != "" {
	fmt.Fprintf(&b, "  %-16s %s\n", "Type:", r.SensorType)
	// }
	// if r.Timestamp != 0 {
	t := time.Unix(r.Timestamp, 0).UTC()
	fmt.Fprintf(&b, "  %-16s %s (%d)\n", "Timestamp:", t.Format("2006-01-02 15:04:05 UTC"), r.Timestamp)
	// }
	// if r.TemperatureK != 0 {
	fmt.Fprintf(&b, "  %-16s %.2f K\n", "Temperature:", r.TemperatureK)
	// }
	// if r.PressureBar != 0 {
	fmt.Fprintf(&b, "  %-16s %.2f bar\n", "Pressure:", r.PressureBar)
	// }
	// if r.WaterLevelM != 0 {
	fmt.Fprintf(&b, "  %-16s %.2f m\n", "Water level:", r.WaterLevelM)
	// }
	// if r.VoltageSupplyV != 0 {
	fmt.Fprintf(&b, "  %-16s %.2f V\n", "Voltage:", r.VoltageSupplyV)
	// }
	// if r.HumidityPercent != 0 {
	fmt.Fprintf(&b, "  %-16s %.1f %%\n", "Humidity:", r.HumidityPercent)
	// }
	// if r.OperatorNotes != "" {
	fmt.Fprintf(&b, "  %-16s %s\n", "Notes:", r.OperatorNotes)
	// }

	// b.WriteString("────────────────────────────────────────")
	b.WriteString("")
	return b.String()
}

type Capability string

const (
	CapTemperature Capability = "temperature"
	CapPressure    Capability = "pressure"
	CapWater       Capability = "water"
	CapVoltage     Capability = "voltage"
	CapHumidity    Capability = "humidity"
)

var FieldsFor = map[Capability][]string{
	CapTemperature: {"temperature_K"},
	CapPressure:    {"pressure_bar"},
	CapWater:       {"water_level_meters"},
	CapVoltage:     {"voltage_supply_v"},
	CapHumidity:    {"humidity_percent"},
}

func (r EvaluationSensorReading) Capabilities() map[Capability]bool {
	caps := make(map[Capability]bool)
	for _, part := range strings.Split(r.SensorType, "/") {
		part = strings.TrimSpace(strings.ToLower(part))
		if part != "" {
			caps[Capability(part)] = true
		}
	}
	return caps
}

func (r EvaluationSensorReading) Has(c Capability) bool {
	return r.Capabilities()[c]
}

func Validate(r EvaluationSensorReading) []string {
	var problems []string
	caps := r.Capabilities() // parse once, reuse for both directions

	for cap := range caps {

		fields, known := FieldsFor[cap]

		if !known {
			problems = append(problems, fmt.Sprintf("unknown sensor capability %q", cap))
			continue
		}
		for _, field := range fields {
			v, ok := r.Value(field)

			if !ok {
				continue
			}

			// Consistency: a declared capability must have a real measurement.
			if v == 0 {
				problems = append(problems,
					fmt.Sprintf("%s declared but %s is zero/missing", cap, field))
				continue // don't also range-check a value we know is absent
			}

			// Range check only for present values.
			if lim, hasLim := Limits[field]; hasLim && !lim.Contains(v) {
				problems = append(problems,
					fmt.Sprintf("%s out of range: %v (want %v–%v)", field, v, lim.Min, lim.Max))
			}
		}
	}

	for cap, fields := range FieldsFor {
		if caps[cap] {
			continue // declared — already handled above
		}
		for _, field := range fields {
			if v, ok := r.Value(field); ok && v != 0 {
				problems = append(problems,
					fmt.Sprintf("%s not declared but %s = %v", cap, field, v))
			}
		}
	}
	return problems
}

type Range struct {
	Min, Max float64
}

func (r Range) Contains(v float64) bool {
	return v >= r.Min && v <= r.Max
}

var Limits = map[string]Range{
	"temperature_K":      {553, 873},
	"pressure_bar":       {60, 160},
	"water_level_meters": {5.0, 15.0},
	"voltage_supply_v":   {229.0, 231.0},
	"humidity_percent":   {40.0, 80.0},
}
