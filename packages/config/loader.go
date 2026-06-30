package config

import (
	"fmt"
	"os"
	"reflect"
	"strings"
)

type Loader struct {
	providers []Provider
}

type Provider interface {
	Name() string
	Load() (map[string]interface{}, error)
}

func NewLoader() *Loader {
	return &Loader{}
}

func (l *Loader) AddProvider(p Provider) {
	l.providers = append(l.providers, p)
}

func (l *Loader) Load(config interface{}) error {
	for _, p := range l.providers {
		data, err := p.Load()
		if err != nil {
			return fmt.Errorf("provider %s: %w", p.Name(), err)
		}
		if err := applyValues(config, data); err != nil {
			return fmt.Errorf("apply %s: %w", p.Name(), err)
		}
	}
	return Validate(config)
}

func Validate(config interface{}) error {
	val := reflect.ValueOf(config)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return nil
	}
	t := val.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("validate")
		if strings.Contains(tag, "required") {
			fv := val.Field(i)
			if fv.IsZero() {
				return fmt.Errorf("field %s is required", field.Name)
			}
		}
	}
	return nil
}

func applyValues(config interface{}, data map[string]interface{}) error {
	val := reflect.ValueOf(config)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return nil
	}
	t := val.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		envTag := field.Tag.Get("env")
		if envTag == "" {
			envTag = field.Name
		}
		if v, ok := data[envTag]; ok {
			fv := val.Field(i)
			if fv.IsValid() && fv.CanSet() {
				sv := fmt.Sprintf("%v", v)
				switch fv.Kind() {
				case reflect.String:
					fv.SetString(sv)
				}
			}
		}
		_ = os.Getenv
	}
	return nil
}
