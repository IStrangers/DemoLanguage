package interpreter

type Stash struct {
	outer        *Stash
	valueMapping map[string]Value
}

func (self *Stash) setValue(name string, value Value) {
	self.valueMapping[name] = value
}

func (self *Stash) getValue(name string) Value {
	stash := self
	for stash != nil {
		value, exists := stash.valueMapping[name]
		if exists {
			return value
		}
		stash = self.outer
	}
	return Value{Null, nil}
}

func (self *Stash) contains(name string) bool {
	_, exists := self.valueMapping[name]
	return exists
}
