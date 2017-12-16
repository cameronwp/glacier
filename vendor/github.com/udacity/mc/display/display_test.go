package display

import (
	"testing"
)

func TestGenerateInfo(t *testing.T) {
	type testI struct {
		Name    string
		Country string
		ID      int
	}

	person1 := testI{
		"person",
		"place",
		1,
	}

	person2 := testI{
		"someone",
		"where",
		2,
	}

	onePerson := []testI{person1, person2}

	header, rows, err := generateTableInfo(onePerson)

	if err != nil {
		t.Error(err)
	}

	if len(header) != 3 {
		t.Errorf("Expected 3 items in the header, found %d", len(header))
	}

	if len(rows) != 2 {
		t.Errorf("Expected 2 rows, found %d", len(rows))
	}

	if header[0] != "Name" {
		t.Errorf("Header is in wrong order: %+v", header)
	}

	if rows[0][0] != "person" {
		t.Errorf("Rows are in wrong order: %+v", rows[0][0])
	}
}
