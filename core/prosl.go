package core

type Prosl interface {
	Convert(yaml string) error
	Validate() error
	Execute() (model.Object, error)
}

type ProslConvertor interface {
	Convert(yaml string)
}
