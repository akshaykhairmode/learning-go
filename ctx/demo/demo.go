package demo

import "context"

type myPrivateKey string

type MyExportedKey string

var ctxKey = "ID"
var ctxPrivateKey myPrivateKey = "ID"
var ctxExportedKey MyExportedKey = "ID"
var valueThatShouldNotChange = "SOME_SECRET_STRING"

func Private() context.Context {
	ctx := context.Background()
	dctx := context.WithValue(ctx, "NAME", "John")
	return context.WithValue(dctx, ctxPrivateKey, valueThatShouldNotChange)
}

func Public() context.Context {
	ctx := context.Background()
	dctx := context.WithValue(ctx, "NAME", "John")
	return context.WithValue(dctx, ctxKey, valueThatShouldNotChange)
}

func PublicWithExportedType() context.Context {
	ctx := context.Background()
	dctx := context.WithValue(ctx, "NAME", "John")
	return context.WithValue(dctx, ctxExportedKey, valueThatShouldNotChange)
}
