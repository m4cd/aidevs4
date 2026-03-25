package dates

import (
	"fmt"
	"time"
)

func CalculateAgeColonYYYYMMDD(birthDateString string) (a int, err error) {
	birthDate, err := time.Parse("2006-01-02", birthDateString)
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return 1, err
	}

	today := time.Now()

	age := today.Year() - birthDate.Year()

	// If the birthday hasn't occurred this year yet, subtract one from age
	if today.YearDay() < birthDate.YearDay() {
		age--
	}

	return age, nil
}

func ExtractYearYYYYMMDD(birthDateString string) (a int32, err error) {
	birthDate, err := time.Parse("2006-01-02", birthDateString)
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return 1, err
	}

	return int32(birthDate.Year()), nil
}
