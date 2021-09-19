package optimizer

import (
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/lothar1998/v2x-optimizer/internal/identifiable"

	"github.com/lothar1998/v2x-optimizer/pkg/optimizer"
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

type IdentifiableWrapper struct {
	optimizer.Optimizer
}

func (w *IdentifiableWrapper) Identifier() string {
	val := reflect.ValueOf(w.Optimizer)
	vType := val.Type()

	var params []keyValue

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)

		if field.CanInterface() {
			p := keyValue{vType.Field(i).Name, field.Interface()}
			params = append(params, p)
		}
	}

	sort.Slice(params, func(i, j int) bool {
		return params[i].key < params[j].key
	})

	if len(params) == 0 {
		return vType.Name()
	}

	return vType.Name() + "," + join(params, ",")
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
