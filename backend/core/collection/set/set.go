package set

import (
	"gitlab.com/alienspaces/go-mud/backend/core/collection/mmap"
	"gitlab.com/alienspaces/go-mud/backend/core/collection/slice"
)

type Set[T comparable] map[T]struct{}

func New[T comparable](elements ...T) Set[T] {
	s := Set[T]{}

	for _, e := range elements {
		s[e] = struct{}{}
	}

	return s
}

func (s Set[T]) Add(elements ...T) Set[T] {
	for _, e := range elements {
		s[e] = struct{}{}
	}

	return s
}

func (s Set[T]) Contains(t T) bool {
	_, ok := s[t]
	return ok
}

func FromSlice[T comparable](slices ...[]T) Set[T] {
	set := Set[T]{}

	for _, items := range slices {
		for _, e := range items {
			set[e] = struct{}{}
		}
	}

	return set
}

func FromSliceWithKey[T any, K comparable](keyFn func(T) K, slices ...[]T) Set[K] {
	set := Set[K]{}

	for _, items := range slices {
		for _, e := range items {
			set[keyFn(e)] = struct{}{}
		}
	}

	return set
}

func FromSlicePtr[T comparable](items *[]T) Set[T] {
	if items == nil {
		return Set[T]{}
	}

	return FromSlice(*items)
}

// IsCompleteBipartiteGraph checks whether every node in setA and setB form a complete bipartite graph.
// In setA, a Set[T] map value is a node, with the map key as the node label. In setB, a key-value pair is a node, and the key is the node label.
// For two nodes to be connected, the element in setB must be in the Set[T] value of setA.
func IsCompleteBipartiteGraph[T comparable](setA map[T]Set[T], setB Set[T]) (disconnectedNodeALabels []T, disconnectedNodeBLabels []T) {
	var connectedB []T

	for labelA, nodeA := range setA {
		intersection := Intersection(nodeA, setB)
		if len(intersection) == 0 {
			disconnectedNodeALabels = append(disconnectedNodeALabels, labelA)
			continue
		}
		connectedB = append(connectedB, intersection...)
	}

	return disconnectedNodeALabels, Difference(FromSlice(connectedB), setB)
}

func Intersection[T comparable](setA Set[T], setB Set[T]) []T {
	var smallerSet Set[T]
	var biggerSet Set[T]

	if len(setA) < len(setB) {
		smallerSet = setA
		biggerSet = setB
	} else {
		smallerSet = setB
		biggerSet = setA
	}

	var intersection []T

	for k := range smallerSet {
		if _, ok := biggerSet[k]; ok {
			intersection = append(intersection, k)
		}
	}

	return intersection
}

func Union[T comparable](sets ...Set[T]) Set[T] {
	union := Set[T]{}

	for _, s := range sets {
		for k := range s {
			union[k] = struct{}{}
		}
	}

	return union
}

func Difference[T comparable](setA Set[T], setB Set[T]) []T {
	var diff []T

	for k := range setA {
		if _, ok := setB[k]; !ok {
			diff = append(diff, k)
		}
	}

	for k := range setB {
		if _, ok := setA[k]; !ok {
			diff = append(diff, k)
		}
	}

	return diff
}

// FindDuplicates returns a slice of duplicates if any.
func FindDuplicates[T comparable](items []T) []T {
	duplicatesSet := Set[T]{}

	set := Set[T]{}
	for _, x := range items {
		if _, ok := set[x]; ok {
			duplicatesSet[x] = struct{}{}
		}

		set[x] = struct{}{}
	}

	duplicates := make([]T, len(duplicatesSet))
	i := 0
	for k := range duplicatesSet {
		duplicates[i] = k
		i++
	}

	return duplicates
}

// FindUnique returns the slice of unique values.
func FindUnique[T comparable](values []T) []T {
	found := map[T]bool{}
	var unique []T

	for _, value := range values {
		if _, ok := found[value]; ok {
			continue
		}

		unique = append(unique, value)
		found[value] = true
	}

	return unique
}

// IsSubset return true if sourceSlice is a subset of targetSlice
func IsSubset[T comparable](sourceSlice []T, targetSlice []T) bool {
	targetSliceMap := map[T]bool{}
	for _, target := range targetSlice {
		targetSliceMap[target] = true
	}

	for _, source := range sourceSlice {
		if !targetSliceMap[source] {
			return false
		}
	}
	return true
}

// Separate splits setA and setB into three separate sets of elements: unique elements in setA, common elements in setA and setB, and unique elements in setB.
// The keyFn is used to compare the elements of setA and setB.
//
// This implementation uses setA to populate the first and second return values, and setB to populate the third return value.
// Consequently, 'uniqueA' and 'common' return values will contain whatever data is contained in setA.
// (i.e., if setA is a slice of structs, some fields of setA may be populated, whereas some fields in setB may not be populated, and vice versa).
// If setA has more information than setB, the 'uniqueA' and 'common' return values will contain more field-level data than the 'uniqueB' return value, and vice versa.
func Separate[T any, K comparable](setA []T, setB []T, keyFn func(T) K) (uniqueA []T, common []T, uniqueB []T) {
	aMap := mmap.FromSlice(keyFn, setA)

	for _, b := range setB {
		bKey := keyFn(b)

		if _, ok := aMap[bKey]; ok {
			common = append(common, aMap[bKey])
			delete(aMap, bKey) // remaining aMap values are unique to setA
		} else {
			uniqueB = append(uniqueB, b)
		}
	}

	return slice.FromMap(aMap), common, uniqueB
}
