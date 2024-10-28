package teardowns

import "errors"

type Func func() error

type Funcs []Func

func New() *Funcs {
	funcs := make(Funcs, 0)
	return &funcs
}

func (t *Funcs) Append(f Func) {
	if t == nil {
		temp := make(Funcs, 0)
		t = &temp
	}
	*t = append(*t, f)
}

func (t *Funcs) AppendMany(fs ...Func) {
	if t == nil {
		temp := make(Funcs, 0, len(fs))
		t = &temp
	}
	*t = append(*t, fs...)
}

func (t *Funcs) Clear() {
	if t == nil {
		return
	}
	*t = nil
}

func (t *Funcs) Execute() error {
	if t == nil {
		return nil
	}

	var err error
	for i := len(*t) - 1; i > 0; i-- {
		fErr := (*t)[i]()
		if fErr != nil {
			err = errors.Join(err, fErr)
		}
	}

	t.Clear()

	return err
}

func (t *Funcs) ToTeardownFunc() Func {
	return func() error {
		return t.Execute()
	}
}
