package meta

func (f facade) MustHave() required {
	return required{request: f.request, context: f.context}
}

func (f facade) Default(defaults string) withDefault {
	return withDefault{request: f.request, context: f.context, defaults: defaults}
}
