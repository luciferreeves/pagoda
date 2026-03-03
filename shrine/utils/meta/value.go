package meta

func (f facade) MustHave() required {
	return required{request: f.Request, context: f.context}
}

func (f facade) Default(defaults string) withDefault {
	return withDefault{request: f.Request, context: f.context, defaults: defaults}
}
