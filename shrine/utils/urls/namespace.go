package urls

func SetNamespace(namespace string) {
	registry.mutex.Lock()
	defer registry.mutex.Unlock()
	registry.currentNamespace = namespace
}
