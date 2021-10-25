package optimizer

import (
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/lothar1998/v2x-optimizer/internal/identifiable"
	"github.com/lothar1998/v2x-optimizer/pkg/optimizer"
)

const (
	TagIncludeKey   = "id_include"
	TagIncludeValue = "true"

	TagRenameKey = "id_rename"
	TagNameKey   = "id_name"
)

type keyValue struct {
	key   string
	value interface{}
}

func (p keyValue) toEntry() string {
	return fmt.Sprintf("%s:%v", p.key, p.value)
}

type IdentifiableOptimizer interface {
	identifiable.Identifiable
	optimizer.Optimizer
}

type IdentifiableAdapter struct {
	optimizer.Optimizer
}

func (w *IdentifiableAdapter) Identifier() string {
	rValue := reflect.ValueOf(w.Optimizer)
	if rValue.Kind() == reflect.Ptr {
		rValue = rValue.Elem()
	}

	var name string
	var params []keyValue

	for i := 0; i < rValue.NumField(); i++ {
		field := rValue.Field(i)
		fieldType := rValue.Type().Field(i)

		if includeTagValue, ok := fieldType.Tag.Lookup(TagIncludeKey); !ok || includeTagValue != TagIncludeValue {
			if _, ok = fieldType.Tag.Lookup(TagNameKey); ok {
				name = field.Interface().(string)
			}
			continue
		}

		var p keyValue

		if renameTagValue, ok := fieldType.Tag.Lookup(TagRenameKey); ok {
			p = keyValue{renameTagValue, field.Interface()}
		} else {
			p = keyValue{fieldType.Name, field.Interface()}
		}

		params = append(params, p)
	}

	sort.Slice(params, func(i, j int) bool {
		return params[i].key < params[j].key
	})

	if name == "" {
		name = rValue.Type().Name()
	}

	if len(params) == 0 {
		return name
	}

	return name + "," + join(params, ",")
}

func join(elems []keyValue, sep string) string {
	switch len(elems) {
	case 0:
		return ""
	case 1:
		return elems[0].toEntry()
	}

	var b strings.Builder
	b.WriteString(elems[0].toEntry())
	for _, s := range elems[1:] {
		b.WriteString(sep)
		b.WriteString(s.toEntry())
	}
	return b.String()
}
