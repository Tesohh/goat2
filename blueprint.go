package goat

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strconv"
	"strings"
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

// // TODO: Also check if the server has this type already defined in it's schemas, and if that's the case then use that as a $ref
// func (b fieldBlueprint) AsOpenAPIParameter(description *string) openapi31.Parameter {
// 	return openapi31.Parameter{
// 		Name:        b.ParamName,
// 		In:          openapi31.ParameterIn(b.GetFrom), // WARN: In doesn't include body!
// 		Description: description,
// 		Required:    new(bool), // TODO: add optional and required prameters
// 		Deprecated:  new(bool), // TODO: add deprecation
// 		Schema:      schema,
// 		Content:     map[string]openapi31.MediaType{},
// 		Style:       &"",
// 		Example:     &nil,
// 	}
// }

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

func compileBlueprints(v any) []fieldBlueprint {
	blueprints := []fieldBlueprint{}

	rv := reflect.ValueOf(v)
	t := rv.Type()

	n := t.NumField()
	for i := 0; i < n; i++ {
		bp := fieldBlueprint{}
		field := t.Field(i)
		bp.FieldName = field.Name

		tag := field.Tag.Get("goat")
		tags := strings.Split(tag, ",")

		if len(tags) > 0 && tags[0] != "" {
			bp.ParamName = tags[0]
		} else {
			bp.ParamName = strings.ToLower(bp.FieldName)
		}

		if len(tags) > 1 {
			bp.GetFrom = tags[1]
		} else {
			if field.Type.Kind() == reflect.Struct {
				bp.GetFrom = "body"
			} else {
				bp.GetFrom = "query"
			}
		}

		bp.Type = field.Type

		blueprints = append(blueprints, bp)
	}

	return blueprints
}
