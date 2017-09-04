package mutator

import (
	"reflect"
)

func Create(functionPointer interface{}) interface{} {
	//use reflection to find out what we need about the function providewd
	funcReflect := reflect.ValueOf(functionPointer)
	function := funcReflect.Elem()

	outType := function.Type().Out(0)

	//Start mutator function ----------------------------------------
	//this is the function that we will use as the template for the specific
	//mutator that we are going to create
	var mutate = func(in []reflect.Value) (out []reflect.Value) {
		inValue := in[0]
		outValue := reflect.New(outType).Elem()
		//scan for all tags in input that point to the ∆ output type
		for i := 0; i < inValue.NumField(); i++ {
			tag := inValue.Type().Field(i).Tag.Get("∆." + outType.String())
			//if we found one and it points to a setable field then try to set it
			if outValue.FieldByName(tag).CanSet() {
				outValue.FieldByName(tag).Set(inValue.Field(i))
			}
		}
		//the return type is an array of reflect values of length one of type
		return []reflect.Value{outValue}
	}
	//End mutator function ----------------------------------------

	//create the mutator function using reflection
	mutator := reflect.MakeFunc(function.Type(), mutate)

	//set the pointer in the original function pointer and also return it
	function.Set(mutator)
	return functionPointer
}
