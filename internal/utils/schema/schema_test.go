package schema

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

type schemaTest struct {
	ID        uint64
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

var validSchema = schemaTest{
	ID:        100,
	Name:      "John",
	CreatedAt: time.Now(),
	UpdatedAt: time.Now(),
}

func TestToSnakeCase(t *testing.T) {
	strValue := "HelloHowAreYou"
	expectedValue := "hello_how_are_you"
	value := toSnakeCase(strValue)
	if value != expectedValue {
		t.Errorf("got: value = %s | expected: value = %s", value, expectedValue)
	}
}

func TestFormatValueForInsert(t *testing.T) {
	strValue := "Hello, how are you?"
	timeValue := time.Date(2000, 1, 1, 0, 0, 0, 0, time.Local)
	floatValue := 100.1
	uint8Value := uint8(100)
	strSliceValue := []string{"Hello", "how are you?"}

	expectedValue := "$$Hello, how are you?$$"
	value := formatValueForInsert(strValue)
	if value != expectedValue {
		t.Errorf("got: value = %s | expected: value = %s", value, expectedValue)
	}

	expectedValue = "'2000-01-01 00:00:00'"
	value = formatValueForInsert(timeValue)
	if value != expectedValue {
		t.Errorf("got: value = %s | expected: value = %s", value, expectedValue)
	}

	expectedValue = fmt.Sprintf("%f", floatValue)
	value = formatValueForInsert(&floatValue)
	if value != expectedValue {
		t.Errorf("got: value = %s | expected: value = %s", value, expectedValue)
	}

	expectedValue = fmt.Sprintf("%d", uint8Value)
	value = formatValueForInsert(&uint8Value)
	if value != expectedValue {
		t.Errorf("got: value = %s | expected: value = %s", value, expectedValue)
	}

	expectedValue = `'{"Hello", "how are you?"}'`
	value = formatValueForInsert(strSliceValue)
	if value != expectedValue {
		t.Errorf("got: value = %s | expected: value = %s", value, expectedValue)
	}
}

func TestParseFieldsToString(t *testing.T) {
	expectedFields, expectedValues := []string{"id", "name"}, []string{"100", "$$John$$"}
	fields, values := parseFieldsToString(validSchema)
	if !reflect.DeepEqual(fields, expectedFields) || !reflect.DeepEqual(values, expectedValues) {
		t.Errorf("got: fields = %v and values = %v | expected: fields = %v and values = %v", fields, values, expectedFields, expectedValues)
	}

	expectedFields, expectedValues = []string{"name"}, []string{"$$John$$"}
	fields, values = parseFieldsToString(validSchema, "id")
	if !reflect.DeepEqual(fields, expectedFields) || !reflect.DeepEqual(values, expectedValues) {
		t.Errorf("got: fields = %v and values = %v | expected: fields = %v and values = %v", fields, values, expectedFields, expectedValues)
	}
}

func TestParseFieldsToInsertQuery(t *testing.T) {
	expectedFields, expectedValues := "id, name", "100, $$John$$"
	fields, values := ParseFieldsToInsertQuery(validSchema)
	if !reflect.DeepEqual(fields, expectedFields) || !reflect.DeepEqual(values, expectedValues) {
		t.Errorf("got: fields = %v and values = %v | expected: fields = %v and values = %v", fields, values, expectedFields, expectedValues)
	}

	expectedFields, expectedValues = "name", "$$John$$"
	fields, values = ParseFieldsToInsertQuery(validSchema, "id")
	if !reflect.DeepEqual(fields, expectedFields) || !reflect.DeepEqual(values, expectedValues) {
		t.Errorf("got: fields = %v and values = %v | expected: fields = %v and values = %v", fields, values, expectedFields, expectedValues)
	}
}

func TestParseFieldsToIUpdateQuery(t *testing.T) {
	expectedValue := "id = 100, name = $$John$$"
	value := ParseFieldsToUpdateQuery(validSchema)
	if !reflect.DeepEqual(value, expectedValue) {
		t.Errorf("got: value = %v | expected: value = %v", value, expectedValue)
	}

	expectedValue = "name = $$John$$"
	value = ParseFieldsToUpdateQuery(validSchema, "id")
	if !reflect.DeepEqual(value, expectedValue) {
		t.Errorf("got: value = %v | expected: value = %v", value, expectedValue)
	}
}
