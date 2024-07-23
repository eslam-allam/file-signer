package slice

import "fmt"

func Filter[T any](items []T, predicate func(T) bool) (filtered []T) {
	filtered = make([]T, 0)
	for _, item := range items {
		if predicate(item) {
			filtered = append(filtered, item)
		}
	}
	return filtered
}

func Map[T any, S any](items []T, trasform func(T) S) (mapped []S) {
	mapped = make([]S, len(items))
	for i, item := range items {
		mapped[i] = trasform(item)
	}
	return mapped
}

func MapWithErr[T any, S any](items []T, trasform func(T) (S, error)) (mapped []S, err error) {
	mapped = make([]S, len(items))
	for i, item := range items {
		mapped[i], err = trasform(item)
		if err != nil {
			return nil, fmt.Errorf("failed to map collection: %w", err)
		}
	}
	return mapped, nil
}

func AnyMatch[T any](items []T, predicate func(T) bool) bool {
	for _, item := range items {
		if predicate(item) {
			return true
		}
	}
	return false
}

func AllMatch[T any](items []T, predicate func(T) bool) bool {
	for _, item := range items {
		if !predicate(item) {
			return false
		}
	}
	return true
}

func ShiftSliceRight[T any](arr []T) []T {
	if len(arr) <= 1 {
		return arr
	}

	// Get the last element
	last := arr[len(arr)-1]

	// Create a new slice with the last element at the beginning
	shifted := append([]T{last}, arr[:len(arr)-1]...)

	return shifted
}

func ShiftSliceLeft[T any](arr []T) []T {
	if len(arr) <= 1 {
		return arr
	}

	// Get the last element
	first := arr[0]

	// Create a new slice with the last element at the beginning
	shifted := append(arr[1:], first)

	return shifted
}
