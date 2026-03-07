package collections

type Set[T comparable] map[T]struct{}

func SetOf[T comparable](values ...T) Set[T] {
	set := make(Set[T], len(values))
	for _, value := range values {
		set[value] = struct{}{}
	}
	return set
}

func (set Set[T]) Has(value T) bool {
	_, exists := set[value]
	return exists
}
