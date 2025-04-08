package logging

type KeyVal struct {
	Key string
	Val interface{}
}

func NewKeyVal(key string, val interface{}) KeyVal {
	return KeyVal{Key: key, Val: val}
}

func NewMessage(msg string) KeyVal {
	return NewKeyVal("msg", msg)
}

func NewError(err error) []KeyVal {
	return []KeyVal{
		NewKeyVal("error", true),
		NewKeyVal("msg", func() string {
			if err != nil {
				return err.Error()
			}
			return ""
		}()),
	}
}
