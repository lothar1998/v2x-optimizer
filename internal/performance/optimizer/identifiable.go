package optimizer

import (
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/lothar1998/v2x-optimizer/pkg/optimizer"
)

type keyValue struct {
	key   string
	value interface{}
}

func (p keyValue) toEntry() string {
	return fmt.Sprintf("%s:%v", p.key, p.value)
}

type Identifiable interface {
	optimizer.Optimizer
	Identifier() string
}

type IdentifiableOptimizer struct {
	optimizer.Optimizer
}

func (w *IdentifiableOptimizer) Identifier() string {
	rValue := reflect.ValueOf(w.Optimizer)
	if rValue.Kind() == reflect.Ptr {
		rValue = rValue.Elem()
	}

	rType := rValue.Type()

	var params []keyValue

	for i := 0; i < rValue.NumField(); i++ {
		field := rValue.Field(i)

		if field.CanInterface() {
			p := keyValue{rType.Field(i).Name, field.Interface()}
			params = append(params, p)
		}
	}

	sort.Slice(params, func(i, j int) bool {
		return params[i].key < params[j].key
	})

	if len(params) == 0 {
		return rType.Name()
	}

	return rType.Name() + "," + join(params, ",")
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
