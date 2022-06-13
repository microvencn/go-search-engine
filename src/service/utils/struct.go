package utils

type KeyValue[K string | int | float64, V string | int | []int] struct {
	key   K
	value V
}

func (kv KeyValue[K, V]) Key() K {
	return kv.key
}

func (kv KeyValue[K, V]) Value() V {
	return kv.value
}

func (kv *KeyValue[K, V]) SetKey(key K) {
	kv.key = key
}

func (kv *KeyValue[K, V]) SetValue(val V) {
	kv.value = val
}
