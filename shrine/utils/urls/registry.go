package urls

var registry = &routeRegistry{
	routes: make(map[string]registeredRoute),
}
