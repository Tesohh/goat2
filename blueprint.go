package goat

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strconv"
)

type fieldBlueprint struct {
	FieldName string
	ParamName string
	Type      reflect.Type
	GetFrom   string // query, path or body
}

func (b fieldBlueprint) String() string {
	return fmt.Sprintf("{ %s:%s of type %s getFrom:%s }", b.FieldName, b.ParamName, b.Type.String(), b.GetFrom)
}

func (b fieldBlueprint) Cast(s string) (any, error) {
	switch b.Type.Kind() {
	case reflect.Int:
		return strconv.Atoi(s)
	case reflect.Float64:
		return strconv.ParseFloat(s, 64)
	case reflect.Float32:
		value, err := strconv.ParseFloat(s, 32)
		return float32(value), err
	case reflect.String:
		return s, nil
	}

	// TODO: go has a reflect.Convert function

	return reflect.Zero(b.Type).Interface(), fmt.Errorf("cannot cast from string to %s", b.Type.Kind())
}

func (b fieldBlueprint) SetField(params reflect.Value, s *Server, r *http.Request) error {
	field := params.FieldByName(b.FieldName)

	if !field.IsValid() {
		return fmt.Errorf("field %s is invalid", b.FieldName)
	}
	if !field.CanSet() {
		return fmt.Errorf("cannot set params field %s", b.FieldName)
	}

	if b.GetFrom == "query" {
		raw := r.URL.Query().Get(b.ParamName)

		v, err := b.Cast(raw)
		if err != nil {
			return err
		}
		field.Set(reflect.ValueOf(v))
	} else if b.GetFrom == "path" {
		raw := r.PathValue(b.ParamName)

		v, err := b.Cast(raw)
		if err != nil {
			return err
		}
		field.Set(reflect.ValueOf(v))
	} else if b.GetFrom == "body" {

		if r.Body == nil {
			return fmt.Errorf("body is empty")
		}

		structField := params.FieldByName(b.FieldName)
		v := reflect.New(structField.Type()).Interface()
		err := json.NewDecoder(r.Body).Decode(&v)
		if err == io.EOF {
			return fmt.Errorf("JSON body is empty (EOF)")
		} else if err != nil {
			return err
		}

		field.Set(reflect.ValueOf(v).Elem())
	} else {
		fmt.Println("Unknown GetFrom option", b.GetFrom)
	}

	return nil
}

func compileBlueprints(v any) ([]fieldBlueprint, error) {
	blueprints := []fieldBlueprint{}

	rv := reflect.ValueOf(v)
	t := rv.Type()

	n := t.NumField()
	for i := 0; i < n; i++ {
		bp := fieldBlueprint{}
		field := t.Field(i)
		bp.FieldName = field.Name
		bp.Type = field.Type

		if field.Type.Kind() == reflect.Struct {
			bp.GetFrom = "body"
		} else {
			query, queryOk := field.Tag.Lookup("query")
			path, pathOk := field.Tag.Lookup("path")

			if queryOk && pathOk {
				return nil, fmt.Errorf("cannot have set both `query` and `path` tag on field %s of %s", field.Name, t.Name())
			} else if !queryOk && !pathOk {
				return nil, fmt.Errorf("need to set either `query` or `path` tag on field %s of %s", field.Name, t.Name())
			}

			if queryOk {
				if query == "" {
					return nil, fmt.Errorf("tag `query` cannot be empty on field %s of %s", field.Name, t.Name())
				}
				bp.ParamName = query
				bp.GetFrom = "query"
			}
			if pathOk {
				if path == "" {
					return nil, fmt.Errorf("tag `query` cannot be empty on field %s of %s", field.Name, t.Name())
				}
				bp.ParamName = path
				bp.GetFrom = "path"
			}
		}

		blueprints = append(blueprints, bp)
	}

	return blueprints, nil
}
