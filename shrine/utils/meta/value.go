package meta

func (self facade) MustHave() required {
	return required{request: self.Request, context: self.context}
}

func (self facade) Default(defaults string) withDefault {
	return withDefault{request: self.Request, context: self.context, defaults: defaults}
}