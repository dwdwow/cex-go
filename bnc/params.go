package bnc

import (
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func structToParams(s any) (params url.Values) {
	params = url.Values{}
	typ := reflect.TypeOf(s)
	val := reflect.ValueOf(s)
	kind := val.Kind()
	// s may be nil, invalid value
	if !val.IsValid() {
		return
	}
	if kind != reflect.Struct {
		panic(fmt.Sprintf("bnc: input kind %v is not a struct", kind))
	}
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if !field.IsExported() {
			continue
		}
		name := field.Name
		v := val.FieldByName(name)
		if !v.IsValid() {
			panic(fmt.Sprintf("bnc: field %v is invalid", name))
		}
		// must look up json tag first, because json field name is  in json tag
		tag, ok := field.Tag.Lookup("json")
		if ok {
			// field tag name may be omitempty, so can not use strings.Contains(tag, "omitempty").
			tags := strings.Split(tag, ",")
			if len(tags) > 1 &&
				strings.Contains(tags[1], "omitempty") && v.IsZero() {
				continue
			}
			name = strings.Trim(tags[0], " ")
		}
		tag, ok = field.Tag.Lookup("default")
		if ok && v.IsZero() {
			params.Add(name, tag)
			continue
		}
		params.Add(name, formatAtom(v, false))
	}
	return
}

func formatAtom(v reflect.Value, inSlice bool) string {
	switch v.Kind() {
	case reflect.String:
		if inSlice {
			return `"` + v.String() + `"`
		}
		return v.String()
	case reflect.Uint64:
		return strconv.FormatUint(v.Uint(), 10)
	case reflect.Float64:
		return strconv.FormatFloat(v.Float(), 'f', -1, 64)
	case reflect.Int64:
		return strconv.FormatInt(v.Int(), 10)
	case reflect.Bool:
		return strconv.FormatBool(v.Bool())
	case reflect.Slice, reflect.Array:
		elems := make([]string, v.Len())
		for i := range v.Len() {
			elems[i] = formatAtom(v.Index(i), true)
		}
		return "[" + strings.Join(elems, ",") + "]"
	case reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32:
		return strconv.FormatInt(v.Int(), 10)
	case reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uintptr:
		return strconv.FormatUint(v.Uint(), 10)
	case reflect.Float32:
		return strconv.FormatFloat(v.Float(), 'f', -1, 64)
	default:
		panic(fmt.Sprintf("bnc: field %v is %v", v.Type().String(), v.Kind()))
	}
}

func structToSignedQuery(s any, key string) (query string) {
	params := structToParams(s)
	params.Set("timestamp", strconv.FormatInt(time.Now().UnixMilli(), 10))
	query = params.Encode()
	if key != "" {
		query += "&signature=" + SignByHmacSHA256ToHex(query, key)
	}
	return
}
