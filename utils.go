package main

import (
	"log"
)

type Snippet struct {
	ID       int
	Name     string
	Category string
	Content  []byte
}

type Repository struct {
	Origin string
	Commit string
	Branch string
}

type SnippetLookUp struct {
	Flat    map[int]Snippet
	Grouped map[string][]Snippet
}

func Map[T, V any](ts []T, fn func(T) V) []V {
	result := make([]V, len(ts))

	for i, t := range ts {
		result[i] = fn(t)
	}

	return result
}

func Unwrap[T any](t T, err error) T {
	if err != nil {
		log.Fatal(err)
	}

	return t
}
