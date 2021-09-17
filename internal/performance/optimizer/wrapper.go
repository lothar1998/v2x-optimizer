package optimizer

import (
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/lothar1998/v2x-optimizer/pkg/optimizer"
)

type keyValue struct {
	name  string
	value interface{}
}

func (p keyValue) toEntry() string {
	return fmt.Sprintf("%s:%v", p.name, p.value)
}

type Wrapper struct {
	optimizer.Optimizer
}

func (w *Wrapper) MapKey() string {
	val := reflect.ValueOf(w.Optimizer)

	var params []keyValue
	for i := 0; i < val.NumField(); i++ {
		if val.Field(i).CanInterface() {
			p := keyValue{val.Type().Field(i).Name, val.Field(i).Interface()}
			params = append(params, p)
		}
	}

	sort.Slice(params, func(i, j int) bool {
		return params[i].name < params[j].name
	})

	if len(params) == 0 {
		return val.Type().Name()
	}

	return val.Type().Name() + "," + join(params, ",")
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
