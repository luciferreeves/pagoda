package env

import (
	"reflect"
	"strconv"
	"strings"
	"time"
)

func setFieldFromEnv(field reflect.Value, envKey, defaultVal string) {
	if field.Type() == reflect.TypeFor[time.Duration]() {
		setDurationField(field, envKey, defaultVal)
		return
	}

	switch field.Kind() {
	case reflect.String:
		field.SetString(getEnv(envKey, defaultVal))
	case reflect.Bool:
		defaultBool, _ := strconv.ParseBool(defaultVal)
		field.SetBool(getEnvBool(envKey, defaultBool))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		defaultInt, _ := strconv.ParseInt(defaultVal, 10, 64)
		field.SetInt(getEnvInt(envKey, defaultInt))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		defaultUint, _ := strconv.ParseUint(defaultVal, 10, 64)
		setUintField(field, envKey, defaultUint)
	case reflect.Float32, reflect.Float64:
		defaultFloat, _ := strconv.ParseFloat(defaultVal, 64)
		field.SetFloat(getEnvFloat(envKey, defaultFloat))
	case reflect.Slice:
		setSliceField(field, envKey, defaultVal)
	}
}

func setUintField(field reflect.Value, envKey string, defaultVal uint64) {
	field.SetUint(getEnvUint(envKey, defaultVal))
}

func setDurationField(field reflect.Value, envKey, defaultVal string) {
	if field.Type() == reflect.TypeFor[time.Duration]() {
		defaultDuration, _ := time.ParseDuration(defaultVal)
		field.Set(reflect.ValueOf(getEnvDuration(envKey, defaultDuration)))
	}
}

func setSliceField(field reflect.Value, envKey, defaultVal string) {
	if field.Type().Elem().Kind() == reflect.String {
		var defaultSlice []string
		if defaultVal != "" {
			parts := strings.SplitSeq(defaultVal, ",")
			for part := range parts {
				trimmed := strings.TrimSpace(part)
				if trimmed != "" {
					defaultSlice = append(defaultSlice, trimmed)
				}
			}
		}
		result := getEnvStringSlice(envKey, defaultSlice)
		field.Set(reflect.ValueOf(result))
	}
}
