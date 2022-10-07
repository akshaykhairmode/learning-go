package demo

import "context"

type myPrivateKey string

type MyExportedKey string

//NOTE :: All our keys will have the same string value that is ID.

//This is a normal string key, can be used anyone by outside the package
var ctxKey = "ID"

//This key is of type myPrivateKey and the type is only accessible by this package.
var ctxPrivateKey myPrivateKey = "ID"

//This key is of type MyExportedKey and the type is accessible to whoever imports this package
var ctxExportedKey MyExportedKey = "ID"

//Some value which I may not want others to change.
var valueThatShouldNotChange = "SOME_SECRET_STRING"

//returns a context which has ID key but the datatype is not exported.
func Private() context.Context {
	ctx := context.Background()
	dctx := context.WithValue(ctx, "NAME", "John") //Default value
	return context.WithValue(dctx, ctxPrivateKey, valueThatShouldNotChange)
}

//Returns a context which has ID key but the datatype is string that anyone can create.
func Public() context.Context {
	ctx := context.Background()
	dctx := context.WithValue(ctx, "NAME", "John") //Default value
	return context.WithValue(dctx, ctxKey, valueThatShouldNotChange)
}

//Returns a context which has ID key but the datatype is exported.
func PublicWithExportedType() context.Context {
	ctx := context.Background()
	dctx := context.WithValue(ctx, "NAME", "John") //Default value
	return context.WithValue(dctx, ctxExportedKey, valueThatShouldNotChange)
}
