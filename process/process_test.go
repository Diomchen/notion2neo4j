package process

import (
	"fmt"
	"testing"
)

func TestReadCSV(t *testing.T) {
	csvFilePath := "C:\\Users\\thgy\\Downloads\\aee088f0-0835-432b-b3d1-01cff038238d_Export-a5f5a71a-f4af-4607-befc-47391251e0fe\\网吧 1b4ee31ce44080279e25e167f564e0bb.csv"
	data, err := ProcessCSV(csvFilePath)

	if err != nil {
		t.Error("Error reading CSV file")
	}

	fmt.Printf("%+v\n", data)

}
