package gen

import (
	"github.com/pkg/errors"
	"go/types"
)

func AddTypeImports(t types.Type, imports map[*types.Package]struct{}) {
	switch v := t.(type) {
	case *types.Interface:
		AddInterfaceImports(v, imports)
	case *types.Struct:
		AddStructImports(v, imports)
	case *types.Signature:
		AddFuncImports(v, imports)
	case *types.Tuple:
		AddTupleImports(v, imports)
	case *types.Map:
		AddMapImports(v, imports)
	case *types.Chan:
		AddChanImports(v, imports)
	case *types.Slice:
		AddSliceImports(v, imports)
	case *types.Array:
		AddArrayImports(v, imports)
	case *types.Pointer:
		AddPointerImports(v, imports)
	case *types.Named:
		AddNamedImports(v, imports)
	case *types.Basic:
		basicImports(v, imports)
	default:
		panic(errors.Errorf("unknown type = %v", t))
	}
}

func AddInterfaceImports(t *types.Interface, imports map[*types.Package]struct{}) {
	n := t.NumMethods()
	for i := 0; i < n; i++ {
		method := t.Method(i)
		AddTypeImports(method.Type(), imports)
	}
}

func AddStructImports(t *types.Struct, imports map[*types.Package]struct{}) {
	n := t.NumFields()
	for i := 0; i < n; i++ {
		field := t.Field(i)
		AddTypeImports(field.Type(), imports)
	}
}

func AddFuncImports(f *types.Signature, imports map[*types.Package]struct{}) {
	AddTupleImports(f.Params(), imports)
	AddTupleImports(f.Results(), imports)
}

func AddTupleImports(t *types.Tuple, imports map[*types.Package]struct{}) {
	for i := 0; i < t.Len(); i++ {
		param := t.At(i)
		AddTypeImports(param.Type(), imports)
	}
}

func AddMapImports(t *types.Map, imports map[*types.Package]struct{}) {
	AddTypeImports(t.Key(), imports)
	AddTypeImports(t.Elem(), imports)
}

func AddChanImports(t *types.Chan, imports map[*types.Package]struct{}) {
	AddTypeImports(t.Elem(), imports)
}

func AddSliceImports(t *types.Slice, imports map[*types.Package]struct{}) {
	AddTypeImports(t.Elem(), imports)
}

func AddArrayImports(t *types.Array, imports map[*types.Package]struct{}) {
	AddTypeImports(t.Elem(), imports)
}

func AddPointerImports(t *types.Pointer, imports map[*types.Package]struct{}) {
	AddTypeImports(t.Elem(), imports)
}

func AddNamedImports(t *types.Named, imports map[*types.Package]struct{}) {
	obj := t.Obj()
	pkg := obj.Pkg()
	imports[pkg] = struct{}{}
}

func basicImports(_ *types.Basic, _ map[*types.Package]struct{}) {
}
