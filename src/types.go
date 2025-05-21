package main

type JsonValue any
type JsonString string
type JsonArray []JsonValue
type JsonObject map[JsonString]JsonValue

//type JsonRoot interface {
//	GetRoot() any
//}
