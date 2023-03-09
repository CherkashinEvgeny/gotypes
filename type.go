package gen

import (
	"github.com/dave/jennifer/jen"
	"github.com/pkg/errors"
	"go/types"
	"strconv"
)

func GenerateType(t types.Type) (code jen.Code) {
	switch v := t.(type) {
	case *types.Interface:
		return generateInterfaceType(v)
	case *types.Struct:
		return generateStructType(v)
	case *types.Signature:
		return generateFuncType(v)
	case *types.Tuple:
		return generateTupleType(v)
	case *types.Map:
		return generateMapType(v)
	case *types.Chan:
		return generateChanType(v)
	case *types.Slice:
		return generateSliceType(v)
	case *types.Array:
		return generateArrayType(v)
	case *types.Pointer:
		return generatePointerType(v)
	case *types.Named:
		return generateNamedType(v)
	case *types.Basic:
		return generateBasicType(v)
	default:
		panic(errors.Errorf("unknown type = %v", t))
	}
}

func generateInterfaceType(t *types.Interface) (code jen.Code) {
	n := t.NumMethods()
	methods := make([]jen.Code, n)
	for i := 0; i < n; i++ {
		method := t.Method(i)
		methodType := method.Type()
		if nd, ok := methodType.(*types.Named); ok {
			methods = append(methods, jen.Qual(
				nd.Obj().Pkg().Path(),
				nd.Obj().Name(),
			))
		}
		if s, ok := methodType.(*types.Signature); ok {
			methods = append(methods, jen.Id(method.Name()).Add(generateFuncParams(s)).Add(generateFuncResults(s)))
		}
	}
	return jen.Interface(methods...)
}

func generateStructType(t *types.Struct) (code jen.Code) {
	n := t.NumFields()
	fields := make([]jen.Code, 0, n)
	for i := 0; i < n; i++ {
		field := t.Field(i)
		fields = append(fields, jen.Id(field.Name()).Add(GenerateType(field.Type())))
	}
	return jen.Struct(fields...)
}

func generateFuncType(t *types.Signature) (code jen.Code) {
	return jen.Func().Add(generateFuncParams(t)).Add(generateFuncResults(t))
}

func generateFuncParams(t *types.Signature) (code jen.Code) {
	variables := t.Params()
	n := variables.Len()
	if t.Variadic() {
		n--
	}
	params := make([]jen.Code, 0, n)
	for i := 0; i < n; i++ {
		variable := variables.At(i)
		if variable.Name() != "" {
			params = append(params, jen.Id(variable.Name()).Add(GenerateType(variable.Type())))
		} else {
			params = append(params, GenerateType(variable.Type()))
		}
	}
	if t.Variadic() {
		variable := variables.At(n)
		if variable.Name() != "" {
			params = append(params, jen.Id(variable.Name()).Add(jen.Op("..."), GenerateType(variable.Type())))
		} else {
			params = append(params, jen.Add(jen.Op("..."), GenerateType(variable.Type())))
		}
	}
	return jen.Params(params...)
}

func generateFuncResults(t *types.Signature) (code jen.Code) {
	variables := t.Results()
	n := variables.Len()
	areNamed := false
	params := make([]jen.Code, 0, n)
	for i := 0; i < n; i++ {
		variable := variables.At(i)
		if variable.Name() != "" {
			areNamed = true
			params = append(params, jen.Id(variable.Name()).Add(GenerateType(variable.Type())))
		} else {
			params = append(params, GenerateType(variable.Type()))
		}
	}
	if n > 1 || areNamed {
		code = jen.Params(params...)
	} else {
		code = params[0]
	}
	return code
}

func generateTupleType(t *types.Tuple) (code jen.Code) {
	n := t.Len()
	params := make([]jen.Code, 0, n)
	for i := 0; i < n; i++ {
		params = append(params, GenerateType(t.At(i).Type()))
	}
	return jen.List(params...)
}

func generateMapType(t *types.Map) (code jen.Code) {
	return jen.Map(GenerateType(t)).Add(GenerateType(t))
}

func generateChanType(t *types.Chan) (code jen.Code) {
	switch t.Dir() {
	case types.SendOnly:
		return jen.Chan().Op("<-").Add(GenerateType(t))
	case types.RecvOnly:
		return jen.Op("<-").Chan().Add(GenerateType(t))
	default:
		return jen.Chan().Add(GenerateType(t))
	}
}

func generateSliceType(t *types.Slice) (code jen.Code) {
	return jen.Index().Add(GenerateType(t))
}

func generateArrayType(t *types.Array) (code jen.Code) {
	return jen.Index(jen.Lit(strconv.FormatInt(t.Len(), 10))).Add(GenerateType(t.Elem()))
}

func generatePointerType(t *types.Pointer) (code jen.Code) {
	return jen.Op("*").Add(GenerateType(t))
}

func generateNamedType(t *types.Named) (code jen.Code) {
	obj := t.Obj()
	pkg := obj.Pkg()
	return jen.Qual(pkg.Path(), obj.Name())
}

func generateBasicType(t *types.Basic) (code jen.Code) {
	switch t.String() {
	case "bool":
		return jen.Bool()
	case "int":
		return jen.Int()
	case "int8":
		return jen.Int8()
	case "int16":
		return jen.Int16()
	case "int32":
		return jen.Int32()
	case "int64":
		return jen.Int64()
	case "uint":
		return jen.Uint()
	case "uint8":
		return jen.Uint8()
	case "uint16":
		return jen.Uint16()
	case "uint32":
		return jen.Uint32()
	case "uint64":
		return jen.Uint64()
	case "uintptr":
		return jen.Uintptr()
	case "float32":
		return jen.Float32()
	case "float64":
		return jen.Float64()
	case "complex64":
		return jen.Complex64()
	case "complex128":
		return jen.Complex128()
	case "string":
		return jen.String()
	case "byte":
		return jen.Byte()
	case "rune":
		return jen.Rune()
	default:
		panic("unsupported type")
	}
}
