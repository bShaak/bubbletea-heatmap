package bubbleteaheatmap

import (
	"encoding/json"
	"log"
	"os"
	"testing"
	"time"
)

func readFromFile(filename string) []CalDataPoint {
	var fileData []CalDataPoint

	// Get Data from File
	content, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}

	err = json.Unmarshal(content, &fileData)
	if err != nil {
		log.Fatal("Error during Unmarshall(): ", err)
	}

	return fileData
}

func TestGetIndexDate(t *testing.T) {
	now := time.Now()
	expectedX, expectedY := getDateIndex(now, now, 52)

	actualDate := getIndexDate(expectedX, expectedY, now, 52)
	actualX, actualY := getDateIndex(actualDate, now, 52)

	if actualDate.Month() != now.Month() &&
		actualDate.Year() != now.Year() &&
		actualDate.Day() != now.Day() {
		t.Fatalf("Date doesn't match: %v %v", actualDate, now)
	}

	if actualX != expectedX || actualY != expectedY {
		t.Fatalf("Index Date doesn't match")
	}
}

func TestManyDates(t *testing.T) {
	now := time.Now()
	refDate := now

	for i := 0; i < 350; i++ {

		expectedX, expectedY := getDateIndex(now, refDate, 52)

		actualDate := getIndexDate(expectedX, expectedY, refDate, 52)
		actualX, actualY := getDateIndex(actualDate, refDate, 52)

		if actualDate.Month() != now.Month() ||
			actualDate.Year() != now.Year() ||
			actualDate.Day() != now.Day() {
			t.Fatalf("Date doesn't match: %v %v", actualDate, now)
		}

		if actualX != expectedX || actualY != expectedY {
			t.Fatalf("Index Date doesn't match. now: %v \n actualDate: %v \n actualXY: %v, %v \n expectedXY: %v, %v",
				now,
				actualDate,
				actualX, actualY,
				expectedX, expectedY)
		}
		now = now.AddDate(0, 0, -1)
	}
}

func TestFileDates(t *testing.T) {
	// Get list of dates from file
	now := time.Now()

	calData := readFromFile("./tests/test.json")

	for _, v := range calData {

		expectedX, expectedY := getDateIndex(v.Date, now, 52)

		actualDate := getIndexDate(expectedX, expectedY, now, 52)
		actualX, actualY := getDateIndex(actualDate, now, 52)

		if actualDate.Month() != v.Date.Local().Month() ||
			actualDate.Year() != v.Date.Local().Year() ||
			actualDate.Day() != v.Date.Local().Day() {
			t.Fatalf("Date doesn't match: \nresult: %v %v (%v %v) \nexpected:%v %v (%v %v)", actualDate, actualDate.Weekday(), actualX, actualY, v.Date.Local(), v.Date.Local().Weekday(), expectedX, expectedY)
		}

		if actualX != expectedX || actualY != expectedY {
			t.Fatalf("Index Date doesn't match. testdate: %v \n actualDate: %v \n actualXY: %v, %v \n expectedXY: %v, %v",
				v.Date,
				actualDate,
				actualX, actualY,
				expectedX, expectedY)
		}
	}
}

func TestDateOutsideRange(t *testing.T) {
	now := time.Now()

	dateInRange := now.AddDate(0, -1, -1)
	var testData []CalDataPoint
	parseCalToView(testData, now, 52)

	testData = append(testData, CalDataPoint{Date: dateInRange, Value: 1.0})
	parseCalToView(testData, now, 52)

	testData = append(testData, CalDataPoint{Date: now, Value: 0.0})
	parseCalToView(testData, now, 52)

	dateOutsideRange := now.AddDate(-1, -1, 0)
	dataOutsideRange := CalDataPoint{Date: dateOutsideRange, Value: 0.5}

	testData = append(testData, dataOutsideRange)
	parseCalToView(testData, now, 52)

	dateOutsideRangeFuture := now.AddDate(0, 1, 0)
	dataOutsideRangeFuture := CalDataPoint{Date: dateOutsideRangeFuture, Value: 0.5}

	testData = append(testData, dataOutsideRangeFuture)
	parseCalToView(testData, now, 52)
}
