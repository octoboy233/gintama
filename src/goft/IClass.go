package goft

type IClass interface {
	Build(goft *Goft)
	Name() string
}
