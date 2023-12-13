package vm

var arrayProps = map[string]Value{
	"get": Object{&NativeFunObject{fun: func(call NativeFunCall) Value {
		this := call.this.toObject().self.(*ArrayObject)
		args := call.args
		if len(args) <= 0 {
			return nil
		}
		return this.getValueByIndex(args[0].(IntValue), nil)
	}}},
	"add": Object{&NativeFunObject{fun: func(call NativeFunCall) Value {
		this := call.this.toObject().self.(*ArrayObject)
		args := call.args
		if len(args) <= 0 {
			return nil
		}
		for _, arg := range args {
			this.values = append(this.values, arg)
		}
		this.length = uint32(this.values.size())
		return nil
	}}},
	"remove": Object{&NativeFunObject{fun: func(call NativeFunCall) Value {
		this := call.this.toObject().self.(*ArrayObject)
		args := call.args
		if len(args) <= 0 {
			return nil
		}
		return this.values.remove(int(args[0].toInt()))
	}}},
	"size": Object{&NativeFunObject{fun: func(call NativeFunCall) Value {
		this := call.this.toObject().self.(*ArrayObject)
		return ToIntValue(int64(this.values.size()))
	}}},
}
