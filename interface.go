package gen

import (
	"github.com/dave/jennifer/jen"
	"go/types"
)

type Interface struct {
	Name  string
	Value *types.Interface
}

func FindAllInterfaces(pkg *types.Package) (items []Interface) {
	pkgScope := pkg.Scope()
	names := pkgScope.Names()
	for _, name := range names {
		obj := pkgScope.Lookup(name)
		t := obj.Type()
		named, ok := t.(*types.Named)
		if !ok {
			continue
		}
		iface, ok := named.Underlying().(*types.Interface)
		if !ok {
			continue
		}
		items = append(items, Interface{Name: name, Value: iface})
	}
	return items
}

func FilterInterfaces(items []Interface, names []string) (filteredItems []Interface) {
	namesMap := make(map[string]struct{}, len(names))
	for _, name := range names {
		namesMap[name] = struct{}{}
	}
	filteredItems = make([]Interface, 0, len(names))
	for _, item := range items {
		_, found := namesMap[item.Name]
		if found {
			filteredItems = append(filteredItems, item)
		}
	}
	return filteredItems
}

func ImplementInterfaceMethods(
	typeName string,
	receiver string,
	iface *types.Interface,
	f func(name string, signature *types.Signature) jen.Code,
) (funcs []jen.Code) {
	n := iface.NumMethods()
	funcs = make([]jen.Code, 0, n)
	for i := 0; i < n; i++ {
		method := iface.Method(i)
		methodType := method.Type()
		signature, ok := methodType.(*types.Signature)
		if !ok {
			continue
		}
		funcs = append(funcs, GenerateMethod(
			typeName,
			receiver,
			method.Name(),
			signature,
			func() jen.Code {
				return f(method.Name(), signature)
			},
		))
	}
	return funcs
}

func GenerateMethod(
	typeName string,
	receiver string,
	method string,
	signature *types.Signature,
	bodyFunc func() jen.Code,
) (code jen.Code) {
	return jen.Func().Params(
		jen.Id(receiver).Id(typeName),
	).Id(method).Add(
		generateFuncParams(signature),
	).Add(
		generateFuncResults(signature),
	).BlockFunc(func(body *jen.Group) {
		body.Add(bodyFunc())
	})
}
