package libocfl

type OCFLObject struct {
	path string
}

func Open(path string) *OCFLObject {
	obj := &OCFLObject{path: path}

	return obj
}

func Validate(obj *OCFLObject) {
	runValidation(obj)
}
