package schema

import (
	"database/sql"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/lib/pq"
)

func toSnakeCase(str string) string {
	matchFirstCap := regexp.MustCompile("(.)([A-Z][a-z]+)")
	matchAllCap := regexp.MustCompile("([a-z0-9])([A-Z])")

	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func formatValueForInsert(value interface{}) string {
	switch v := value.(type) {
	case string:
		return fmt.Sprintf("$$%s$$", v)
	case time.Time:
		return fmt.Sprintf("'%s'", v.Format("2006-01-02 15:04:05"))
	case *float64:
		return fmt.Sprintf("%f", *v)
	case *uint8:
		return fmt.Sprintf("%d", *v)
	case []string:
		for i, _v := range v {
			v[i] = fmt.Sprintf(`"%s"`, _v)
			v[i] = strings.Replace(v[i], "'", "''", -1)
		}

		stringValue := strings.Join(v, ", ")

		return fmt.Sprintf("'{%s}'", stringValue)
	case pq.StringArray:
		for i, _v := range v {
			v[i] = fmt.Sprintf(`"%s"`, _v)
			v[i] = strings.Replace(v[i], "'", "''", -1)
		}

		stringValue := strings.Join(v, ", ")

		return fmt.Sprintf("'{%s}'", stringValue)
	case sql.NullTime:
		return fmt.Sprintf("'%s'", v.Time.Format("2006-01-02 15:04:05"))
	case sql.NullString:
		return fmt.Sprintf("$$%s$$", v.String)
	default:
		return fmt.Sprintf("%v", v)
	}
}

func parseFieldsToString(schema any, ignore ...string) ([]string, []string) {
	value := reflect.ValueOf(schema)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	var schemaFields []string
	var schemaValues []string
	for i := 0; i < value.NumField(); i++ {

		field := value.Field(i)
		fieldName := toSnakeCase(value.Type().Field(i).Name)
		shouldContinue := true

		if fieldName == "created_at" || fieldName == "updated_at" {
			shouldContinue = false
		}

		for _, _ignore := range ignore {
			if _ignore == fieldName {
				shouldContinue = false
			}
		}

		if !shouldContinue {
			continue
		}

		if field.Kind() == reflect.Ptr && !field.IsNil() {
			schemaFields = append(schemaFields, fieldName)
			schemaValues = append(schemaValues, formatValueForInsert(field.Elem().Interface()))
		} else if field.Kind() != reflect.Ptr && !reflect.DeepEqual(field.Interface(), reflect.Zero(field.Type()).Interface()) {
			schemaFields = append(schemaFields, fieldName)
			schemaValues = append(schemaValues, formatValueForInsert(field.Interface()))
		}

	}

	return schemaFields, schemaValues
}

func ParseFieldsToInsertQuery(schema any, ignore ...string) (string, string) {
	fields, values := parseFieldsToString(schema, ignore...)

	return strings.Join(fields, ", "), strings.Join(values, ", ")
}

func ParseArrayFieldsToInsertQuery(schemas any, ignore ...string) (string, string) {
	var fieldsStr string
	var valuesStr string

	sliceValue := reflect.ValueOf(schemas)
	if sliceValue.Kind() == reflect.Ptr {
		sliceValue = sliceValue.Elem()
	}

	if sliceValue.Kind() == reflect.Array || sliceValue.Kind() == reflect.Slice {
		for i := 0; i < sliceValue.Len(); i++ {
			schema := sliceValue.Index(i).Interface()
			fields, values := parseFieldsToString(schema, ignore...)

			if i == 0 {
				fieldsStr = strings.Join(fields, ", ")
				valuesStr += fmt.Sprintf("(%s)", strings.Join(values, ", "))
			} else {
				valuesStr += fmt.Sprintf(", (%s)", strings.Join(values, ", "))
			}
		}
	}

	return fieldsStr, valuesStr
}

func ParseFieldsToUpdateQuery(schema any, ignore ...string) string {
	fields, values := parseFieldsToString(schema, ignore...)

	var query []string
	for i := range fields {
		query = append(query, fmt.Sprintf("%s = %s", fields[i], values[i]))
	}

	return strings.Join(query, ", ")
}

type QueryParams struct {
	Select    string
	Join      []string
	Where     string
	OrderBy   string
	SortOrder string
	GroupBy   string
	Offset    uint64
	Limit     uint64
}

func PrepareFindQuery(query string, params QueryParams) string {
	query = strings.Replace(query, "?", params.Select, 1)

	if len(params.Join) != 0 {
		join := strings.Join(params.Join, " ")
		query = fmt.Sprintf("%s %s", query, join)
	}

	if len(params.Where) != 0 {
		query = fmt.Sprintf("%s WHERE %s", query, params.Where)
	}

	if params.OrderBy != "" {
		query = fmt.Sprintf("%s ORDER BY %s", query, params.OrderBy)
	}

	if params.SortOrder != "" {
		query = fmt.Sprintf("%s %s", query, params.SortOrder)
	}

	if params.GroupBy != "" {
		query = fmt.Sprintf("%s GROUP BY %s", query, params.GroupBy)
	}

	if params.Offset != 0 {
		query = fmt.Sprintf("%s OFFSET %d", query, params.Offset)
	}

	if params.Limit != 0 {
		query = fmt.Sprintf("%s LIMIT %d", query, params.Limit)
	}

	return query
}
