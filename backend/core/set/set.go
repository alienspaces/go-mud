package set

// IsCompleteBipartiteGraph checks whether every node in setA and setB form a complete bipartite graph.
// In setA, a map[string]struct{} map value is a node, with the map key as the node label. In setB, a key-value pair is a node, and the key is the node label.
// For two nodes to be connected, the element in setB must be in the map[string]struct{} value of setA.
func IsCompleteBipartiteGraph(setA map[string]map[string]struct{}, setB map[string]struct{}) (disconnectedNodeALabels []string, disconnectedNodeBLabels []string) {
	var connectedB []string

	for labelA, nodeA := range setA {
		intersection := Intersection(nodeA, setB)
		if len(intersection) == 0 {
			disconnectedNodeALabels = append(disconnectedNodeALabels, labelA)
			continue
		}

		for _, labelB := range intersection {
			connectedB = append(connectedB, labelB)
		}
	}

	return disconnectedNodeALabels, Difference(FromStrSlice(connectedB), setB)
}

func Intersection(setA map[string]struct{}, setB map[string]struct{}) []string {
	var smallerSet map[string]struct{}
	var biggerSet map[string]struct{}

	if len(setA) < len(setB) {
		smallerSet = setA
		biggerSet = setB
	} else {
		smallerSet = setB
		biggerSet = setA
	}

	var intersection []string

	for k := range smallerSet {
		if _, ok := biggerSet[k]; ok {
			intersection = append(intersection, k)
		}
	}

	return intersection
}

func Difference(setA map[string]struct{}, setB map[string]struct{}) []string {
	var diff []string

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

func FromStrSlice(items []string) map[string]struct{} {
	set := map[string]struct{}{}

	for _, c := range items {
		set[c] = struct{}{}
	}

	return set
}

func FromStrSlicePtr(items *[]string) map[string]struct{} {
	if items == nil {
		return map[string]struct{}{}
	}

	return FromStrSlice(*items)
}

// FindDuplicates returns the duplicates if any.
func FindDuplicates(items []string) []string {
	duplicatesSet := map[string]struct{}{}

	set := map[string]struct{}{}
	for _, x := range items {
		if _, ok := set[x]; ok {
			duplicatesSet[x] = struct{}{}
		}

		set[x] = struct{}{}
	}

	duplicates := make([]string, len(duplicatesSet))
	i := 0
	for k := range duplicatesSet {
		duplicates[i] = k
		i++
	}

	return duplicates
}
