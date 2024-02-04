package vm

import (
	"fmt"
	"math"
	"strings"
)

type Instruction interface {
	exec(*VM)
}

type InstructionArray []Instruction

func (self *InstructionArray) size() int {
	return len(*self)
}

func (self *InstructionArray) add(instructions ...Instruction) {
	*self = append(*self, instructions...)
}

var (
	Add _Add
	Sub _Sub
	Mul _Mul
	Div _Div
	Mod _Mod

	AND _AND
	OR  _OR
	Not _Not
	Inc _Inc
	Dec _Dec
	Neg _Neg

	EQ _EQ
	NE _NE
	LT _LT
	LE _LE
	GT _GT
	GE _GE

	Pop                 _Pop
	Dup                 _Dup
	SaveResult          _SaveResult
	InitVar             _InitVar
	LoadNull            _LoadNull
	NewObject           _NewObject
	PushArrayValue      _PushArrayValue
	GetPropOrElem       _GetPropOrElem
	GetPropOrElemCallee _GetPropOrElemCallee
	LoadDynamicThis     _LoadDynamicThis
	Ret                 _Ret
)

type _LoadNull struct{}

func (self _LoadNull) exec(vm *VM) {
	vm.push(Const_Null_Value)
	vm.pc++
}

type LoadVal int

func (self LoadVal) exec(vm *VM) {
	vm.push(vm.getValue(int(self)))
	vm.pc++
}

type _Add struct{}

func (self _Add) exec(vm *VM) {
	left := vm.stack[vm.sp-2]
	right := vm.stack[vm.sp-1]

	var value Value
	if left.isString() || right.isString() {
		value = ToStringValue(left.toString() + right.toString())
	} else if left.isFloat() || right.isFloat() {
		value = ToFloatValue(left.toFloat() + right.toFloat())
	} else {
		value = ToIntValue(left.toInt() + right.toInt())
	}

	vm.stack[vm.sp-2] = value
	vm.sp--
	vm.pc++
}

type _Sub struct{}

func (self _Sub) exec(vm *VM) {
	left := vm.stack[vm.sp-2]
	right := vm.stack[vm.sp-1]

	var value Value
	if left.isFloat() || right.isFloat() {
		value = ToFloatValue(left.toFloat() - right.toFloat())
	} else {
		value = ToIntValue(left.toInt() - right.toInt())
	}

	vm.stack[vm.sp-2] = value
	vm.sp--
	vm.pc++
}

type _Mul struct{}

func (self _Mul) exec(vm *VM) {
	left := vm.stack[vm.sp-2]
	right := vm.stack[vm.sp-1]

	var value Value
	if left.isFloat() || right.isFloat() {
		value = ToFloatValue(left.toFloat() * right.toFloat())
	} else {
		value = ToIntValue(left.toInt() * right.toInt())
	}

	vm.stack[vm.sp-2] = value
	vm.sp--
	vm.pc++
}

type _Div struct{}

func (self _Div) exec(vm *VM) {
	left := vm.stack[vm.sp-2]
	right := vm.stack[vm.sp-1]

	vm.stack[vm.sp-2] = ToFloatValue(left.toFloat() / right.toFloat())
	vm.sp--
	vm.pc++
}

type _Mod struct{}

func (self _Mod) exec(vm *VM) {
	left := vm.stack[vm.sp-2]
	right := vm.stack[vm.sp-1]

	var value Value
	if left.isFloat() || right.isFloat() {
		value = ToFloatValue(math.Mod(left.toFloat(), right.toFloat()))
	} else {
		value = ToIntValue(left.toInt() % right.toInt())
	}

	vm.stack[vm.sp-2] = value
	vm.sp--
	vm.pc++
}

type _AND struct{}

func (self _AND) exec(vm *VM) {
	left := vm.stack[vm.sp-2]
	right := vm.stack[vm.sp-1]

	value := ToIntValue(left.toInt() & right.toInt())

	vm.stack[vm.sp-2] = value
	vm.sp--
	vm.pc++
}

type _OR struct{}

func (self _OR) exec(vm *VM) {
	left := vm.stack[vm.sp-2]
	right := vm.stack[vm.sp-1]

	value := ToIntValue(left.toInt() | right.toInt())

	vm.stack[vm.sp-2] = value
	vm.sp--
	vm.pc++
}

type _Not struct{}

func (self _Not) exec(vm *VM) {
	value := vm.stack[vm.sp-1]

	if value.toBool() {
		value = Const_Bool_False_Value
	} else {
		value = Const_Bool_True_Value
	}

	vm.stack[vm.sp-1] = value
	vm.pc++
}

type _Inc struct{}

func (self _Inc) exec(vm *VM) {
	value := vm.stack[vm.sp-1]

	if value.isFloat() {
		value = ToFloatValue(value.toFloat() + 1.0)
	} else {
		value = ToIntValue(value.toInt() + 1)
	}

	vm.stack[vm.sp-1] = value
	vm.pc++
}

type _Dec struct{}

func (self _Dec) exec(vm *VM) {
	value := vm.stack[vm.sp-1]

	if value.isFloat() {
		value = ToFloatValue(value.toFloat() - 1.0)
	} else {
		value = ToIntValue(value.toInt() - 1)
	}

	vm.stack[vm.sp-1] = value
	vm.pc++
}

type _Neg struct{}

func (self _Neg) exec(vm *VM) {
	value := vm.stack[vm.sp-1]

	if value.isFloat() {
		value = ToFloatValue(-value.toFloat())
	} else {
		value = ToIntValue(-value.toInt())
	}

	vm.stack[vm.sp-1] = value
	vm.pc++
}

type _EQ struct{}

func (self _EQ) exec(vm *VM) {
	left := vm.stack[vm.sp-2]
	right := vm.stack[vm.sp-1]

	value := Const_Bool_False_Value
	if left.equals(right) {
		value = Const_Bool_True_Value
	}

	vm.stack[vm.sp-2] = value
	vm.sp--
	vm.pc++
}

type _NE struct{}

func (self _NE) exec(vm *VM) {
	left := vm.stack[vm.sp-2]
	right := vm.stack[vm.sp-1]

	value := Const_Bool_True_Value
	if left.equals(right) {
		value = Const_Bool_False_Value
	}

	vm.stack[vm.sp-2] = value
	vm.sp--
	vm.pc++
}

type _LT struct{}

func (self _LT) exec(vm *VM) {
	left := vm.stack[vm.sp-2]
	right := vm.stack[vm.sp-1]

	value := lessComp(left, right)

	vm.stack[vm.sp-2] = value
	vm.sp--
	vm.pc++
}

type _LE struct{}

func (self _LE) exec(vm *VM) {
	left := vm.stack[vm.sp-2]
	right := vm.stack[vm.sp-1]

	value := Const_Bool_True_Value
	if lessComp(right, left) == Const_Bool_True_Value {
		value = Const_Bool_False_Value
	}

	vm.stack[vm.sp-2] = value
	vm.sp--
	vm.pc++
}

type _GT struct{}

func (self _GT) exec(vm *VM) {
	left := vm.stack[vm.sp-2]
	right := vm.stack[vm.sp-1]

	value := lessComp(right, left)

	vm.stack[vm.sp-2] = value
	vm.sp--
	vm.pc++
}

type _GE struct{}

func (self _GE) exec(vm *VM) {
	left := vm.stack[vm.sp-2]
	right := vm.stack[vm.sp-1]

	value := Const_Bool_True_Value
	if lessComp(left, right) == Const_Bool_True_Value {
		value = Const_Bool_False_Value
	}

	vm.stack[vm.sp-2] = value
	vm.sp--
	vm.pc++
}

func lessComp(left Value, right Value) Value {
	var less bool
	if left.isString() && right.isString() {
		less = strings.Compare(left.toString(), right.toString()) < 0
	} else if left.isInt() || left.isFloat() && right.isInt() || right.isFloat() {
		if left.isFloat() || right.isFloat() {
			less = left.toFloat() < right.toFloat()
		} else {
			less = left.toInt() < right.toInt()
		}
	}
	if less {
		return Const_Bool_True_Value
	}
	return Const_Bool_False_Value
}

type Jeq int

func (self Jeq) exec(vm *VM) {
	vm.sp--
	value := vm.stack[vm.sp]
	if value.toBool() {
		vm.pc += int(self)
	} else {
		vm.pc++
	}
}

type Jeq1 int

func (self Jeq1) exec(vm *VM) {
	value := vm.stack[vm.sp-1]
	if value.toBool() {
		vm.pc += int(self)
	} else {
		vm.sp--
		vm.pc++
	}
}

type JeqNull int

func (self JeqNull) exec(vm *VM) {
	vm.sp--
	value := vm.stack[vm.sp]
	if value.isNull() {
		vm.pc++
	} else {
		vm.pc += int(self)
	}
}

type Jne int

func (self Jne) exec(vm *VM) {
	vm.sp--
	value := vm.stack[vm.sp]
	if !value.toBool() {
		vm.pc += int(self)
	} else {
		vm.pc++
	}
}

type Jne1 int

func (self Jne1) exec(vm *VM) {
	value := vm.stack[vm.sp-1]
	if !value.toBool() {
		vm.pc += int(self)
	} else {
		vm.sp--
		vm.pc++
	}
}

type Jump int

func (self Jump) exec(vm *VM) {
	vm.pc += int(self)
}

type _Pop struct{}

func (self _Pop) exec(vm *VM) {
	vm.sp--
	vm.pc++
}

type _Dup struct{}

func (self _Dup) exec(vm *VM) {
	vm.push(vm.stack[vm.sp-1])
	vm.pc++
}

type _SaveResult struct{}

func (self _SaveResult) exec(vm *VM) {
	vm.result = vm.stack[vm.sp-1]
	vm.sp--
	vm.pc++
}

type ResolveVar string

func (self ResolveVar) exec(vm *VM) {
	vm.refStack.add(&ObjectRef{
		refObject: vm.runtime.globalObject,
		refName:   string(self),
	})
	vm.pc++
}

type _InitVar struct{}

func (self _InitVar) exec(vm *VM) {
	ref := vm.refStack.pop()
	ref.set(vm.stack[vm.sp-1])
	vm.sp--
	vm.pc++
}

type PutVar int

func (self PutVar) exec(vm *VM) {
	ref := vm.refStack.pop()
	ref.set(vm.stack[vm.sp-1])
	vm.sp += int(self)
	vm.pc++
}

type LoadVar string

func (self LoadVar) exec(vm *VM) {
	name := string(self)
	value := vm.getDefining(name)
	if value == nil {
		//wait adjust
		panic(fmt.Sprintf("ReferenceError: '%s' is not defined", name))
	}
	vm.push(value)
	vm.pc++
}

type InitStackVar int

func (self InitStackVar) exec(vm *VM) {
	index := int(self)
	if index <= 0 {
		vm.stack[vm.sb-index] = vm.stack[vm.sp-1]
	} else {
		vm.stack[vm.sb+vm.args+index] = vm.stack[vm.sp-1]
	}
	vm.sp--
	vm.pc++
}

type InitStackVar1 int

func (self InitStackVar1) exec(vm *VM) {
	index := int(self)
	vm.stack[vm.sb+index] = vm.stack[vm.sp-1]

	vm.sp--
	vm.pc++
}

type LoadStackVar int

func (self LoadStackVar) exec(vm *VM) {
	index := int(self)
	var value Value
	if index <= 0 {
		arg := -index
		if arg > vm.args {
			value = Const_Null_Value
		} else {
			value = vm.stack[vm.sb+arg]
		}
	} else {
		value = vm.stack[vm.sb+vm.args+index]
	}
	if value == nil {
		//wait adjust
	}
	vm.push(value)
	vm.pc++
}

type LoadStackVar1 int

func (self LoadStackVar1) exec(vm *VM) {
	index := int(self)
	var value Value
	if index > 0 {
		value = vm.stack[vm.sb+index]
	} else {
		value = vm.stack[vm.sb]
	}
	if value == nil {
		value = Const_Null_Value
	}
	vm.push(value)
	vm.pc++
}

type PutStackVar int

func (self PutStackVar) exec(vm *VM) {
	index := int(self)
	if index > 0 {
		vm.stack[vm.sb+vm.args+index] = vm.stack[vm.sp-1]
	} else {
		//wait adjust
		panic("Illegal stack var index")
	}
	vm.sp--
	vm.pc++
}

type PutStackVar1 int

func (self PutStackVar1) exec(vm *VM) {
	index := int(self)
	if index > 0 {
		vm.stack[vm.sb+index] = vm.stack[vm.sp-1]
	} else {
		//wait adjust
		panic("Illegal stack var index")
	}
	vm.sp--
	vm.pc++
}

type InitStashVar int

func (self InitStashVar) exec(vm *VM) {
	value := vm.stack[vm.sp-1]
	deepLevel := int(self) >> 24
	index := int(self) & 0x00FFFFFF
	stash := vm.stash
	for i := 0; i < deepLevel; i++ {
		stash = stash.outer
	}
	stash.values[index] = value
	vm.sp--
	vm.pc++
}

type LoadStashVar int

func (self LoadStashVar) exec(vm *VM) {
	deepLevel := int(self) >> 24
	index := int(self) & 0x00FFFFFF
	stash := vm.stash
	for i := 0; i < deepLevel; i++ {
		stash = stash.outer
	}
	value := stash.values[index]
	vm.push(value)
	vm.pc++
}

type PutStashVar int

func (self PutStashVar) exec(vm *VM) {
	InitStashVar(self).exec(vm)
}

type LoadDynamicCallee string

func (self LoadDynamicCallee) exec(vm *VM) {
	name := string(self)
	value := vm.getDefining(name)
	if value == nil {
		//wait adjust
	}
	vm.push(Const_Null_Value)
	vm.push(value)
	vm.pc++
}

type BindDefining struct {
	funs []string
	vars []string
}

func (self BindDefining) exec(vm *VM) {
	start := vm.sp - len(self.funs)
	for i, fun := range self.funs {
		value := vm.stack[start+i]
		vm.setDefining(fun, value)
	}
	for _, v := range self.vars {
		vm.setDefining(v, nil)
	}
	vm.sp = start
	vm.pc++
}

type _NewObject struct{}

func (self _NewObject) exec(vm *VM) {
	obj := vm.runtime.newObject()
	vm.push(obj)
	vm.pc++
}

type AddProp string

func (self AddProp) exec(vm *VM) {
	obj := vm.stack[vm.sp-2]
	value := vm.stack[vm.sp-1]
	obj.toObject().self.setProperty(string(self), value)
	vm.sp--
	vm.pc++
}

type GetProp string

func (self GetProp) exec(vm *VM) {
	obj := vm.stack[vm.sp-1]
	if obj == nil {
		//wait adjust
		panic(fmt.Sprintf("Cannot read property '%s' of undefined", self))
	}
	value := obj.toObject().self.getPropertyOrDefault(string(self), Const_Null_Value)
	vm.stack[vm.sp-1] = value
	vm.pc++
}

type GetPropCallee string

func (self GetPropCallee) exec(vm *VM) {
	obj := vm.stack[vm.sp-1]
	if obj == nil {
		//wait adjust
		panic(fmt.Sprintf("Cannot read property '%s' of undefined", self))
	}
	value := obj.toObject().self.getPropertyOrDefault(string(self), Const_Null_Value)
	vm.push(value)
	vm.pc++
}

type _GetPropOrElem struct{}

func (self _GetPropOrElem) exec(vm *VM) {
	obj := vm.stack[vm.sp-2]
	prop := vm.stack[vm.sp-1]
	if obj == nil {
		//wait adjust
		panic(fmt.Sprintf("Cannot read property '%s' of undefined", self))
	}
	value := obj.toObject().getOrDefault(prop, Const_Null_Value)
	vm.stack[vm.sp-2] = value
	vm.sp--
	vm.pc++
}

type _GetPropOrElemCallee struct{}

func (self _GetPropOrElemCallee) exec(vm *VM) {
	obj := vm.stack[vm.sp-2]
	prop := vm.stack[vm.sp-1]
	if obj == nil {
		//wait adjust
		panic(fmt.Sprintf("Cannot read property '%s' of undefined", self))
	}
	value := obj.toObject().getOrDefault(prop, Const_Null_Value)
	vm.stack[vm.sp-1] = value
	vm.pc++
}

type NewArray uint32

func (self NewArray) exec(vm *VM) {
	arr := vm.runtime.newArray(make(ValueArray, 0, self))
	vm.push(arr)
	vm.pc++
}

type _PushArrayValue struct{}

func (self _PushArrayValue) exec(vm *VM) {
	obj := vm.stack[vm.sp-2]
	value := vm.stack[vm.sp-1]
	arrayObj := obj.toObject().self.(*ArrayObject)
	arrayObj.values = append(arrayObj.values, value)
	arrayObj.length++
	vm.sp--
	vm.pc++
}

type NewFun struct {
	funDefinition string
	name          string
	program       *Program
}

func (self NewFun) exec(vm *VM) {
	fun := vm.runtime.newFun(self.name)
	fun.funDefinition = TrimWhitespace(self.funDefinition)
	fun.program = self.program
	fun.stash = vm.stash
	vm.push(Object{fun})
	vm.pc++
}

type EnterFun struct {
	stackSize int
	args      int
}

func (self EnterFun) exec(vm *VM) {
	vm.sb = vm.sp - 1 - vm.args
	d := self.args - vm.args
	if d > 0 {
		ss := vm.sp + d + self.stackSize
		vm.stack.expand(ss)
		vs := vm.stack[vm.sp : ss-self.stackSize]
		for index := range vs {
			vs[index] = Const_Null_Value
		}
		vm.args = self.args
		vm.sp = ss
	} else if self.stackSize > 0 {
		ss := vm.sp + self.stackSize
		vm.stack.expand(ss)
		vs := vm.stack[vm.sp:ss]
		for index := range vs {
			vs[index] = nil
		}
		vm.sp = ss
	}
	vm.pc++
}

type EnterFunStash struct {
	argsToStash bool
	stackSize   int
	stashSize   int
	args        int
}

func (self EnterFunStash) exec(vm *VM) {
	stash := vm.newStash()
	stash.values = make(ValueArray, self.stashSize)

	sp := vm.sp
	vm.sb = sp - vm.args - 1
	ss := self.stackSize
	ea := 0
	if self.argsToStash {
		offset := vm.args - self.args
		copy(stash.values, vm.stack[sp-vm.args:sp])
		if offset > 0 {
			stash.extraArgs = make(ValueArray, offset)
			copy(stash.extraArgs, vm.stack[sp-offset:])
		} else {
			vs := stash.values[vm.args:self.args]
			for i := range vs {
				vs[i] = Const_Null_Value
			}
		}
		sp -= vm.args
	} else {
		d := self.args - vm.args
		if d > 0 {
			ss += d
			ea = d
			vm.args = self.args
		}
	}
	vm.stack.expand(sp + ss - 1)
	if ea > 0 {
		vs := vm.stack[sp : vm.sp+ea]
		for i := range vs {
			vs[i] = Const_Null_Value
		}
	}
	vs := vm.stack[sp+ea : sp+ss]
	for i := range vs {
		vs[i] = nil
	}
	vm.sp = sp + ss
	vm.pc++
}

type EnterFunBody struct {
	EnterBlock
}

func (self EnterFunBody) exec(vm *VM) {
	if self.stashSize > 0 {
		stash := vm.newStash()
		stash.values = make(ValueArray, self.stashSize)
	}

	nsp := vm.sp + self.stackSize
	if self.stackSize > 0 {
		vm.stack.expand(nsp - 1)
		vs := vm.stack[vm.sp:nsp]
		for index := range vs {
			vs[index] = nil
		}
	}
	vm.sp = nsp
	vm.pc++
}

type EnterBlock struct {
	stackSize int
	stashSize int
}

func (self EnterBlock) exec(vm *VM) {
	if self.stashSize > 0 {
		stash := vm.newStash()
		stash.values = make(ValueArray, self.stashSize)
	}
	vm.stack.expand(vm.sp + self.stackSize - 1)
	vs := vm.stack[vm.sp : vm.sp+self.stackSize]
	for i := range vs {
		vs[i] = nil
	}
	vm.sp += self.stackSize
	vm.pc++
}

type LeaveBlock struct {
	stackSize int
	popStash  bool
}

func (self LeaveBlock) exec(vm *VM) {
	if self.popStash {
		vm.stash = vm.stash.outer
	}
	if self.stackSize > 0 {
		vm.sp -= self.stackSize
	}
	vm.pc++
}

type EnterTry struct {
	catchOffset   int
	finallyOffset int
}

func (self EnterTry) exec(vm *VM) {
	var catchPos, finallyPos int
	if self.catchOffset > 0 {
		catchPos = vm.pc + self.catchOffset
	} else {
		catchPos = -1
	}
	if self.finallyOffset > 0 {
		finallyPos = vm.pc + self.finallyOffset
	} else {
		finallyPos = -1
	}
	vm.pushTryFrame(catchPos, finallyPos)
	vm.pc++
}

type LeaveTry struct {
}

func (self LeaveTry) exec(vm *VM) {
	tryStack := vm.tryStack[vm.tryStack.size()-1]
	if tryStack.finallyPos >= 0 {
		tryStack.finallyRet = vm.pc + 1
		vm.pc = tryStack.finallyPos
		tryStack.catchPos = -1
		tryStack.finallyPos = -1
	} else {
		vm.popTryFrame()
		vm.pc++
	}
}

type EnterFinally struct {
}

func (self EnterFinally) exec(vm *VM) {
	tryStack := vm.tryStack[vm.tryStack.size()-1]
	tryStack.finallyPos = -1
	vm.pc++
}

type LeaveFinally struct {
}

func (self LeaveFinally) exec(vm *VM) {
	tryStack := vm.tryStack[vm.tryStack.size()-1]
	ex, ret := tryStack.exception, tryStack.finallyRet
	tryStack.exception = nil
	vm.popTryFrame()
	if ex != nil {
		vm.throw(ex)
		return
	} else {
		if ret != -1 {
			vm.pc = ret
		} else {
			vm.pc++
		}
	}
}

type _LoadDynamicThis struct{}

func (self _LoadDynamicThis) exec(vm *VM) {
	vm.push(vm.runtime.globalObject)
	vm.pc++
}

type Call uint32

func (self Call) exec(vm *VM) {
	n := int(self)
	value := vm.stack[vm.sp-1-n]
	if !value.isObject() {
		//wait adjust
		panic("Value is not a function: " + value.toString())
	}
	object := value.toObject()
	object.self.vmCall(vm, n)
}

type _Ret struct{}

func (self _Ret) exec(vm *VM) {
	vm.stack[vm.sb-1] = vm.stack[vm.sp-1]
	vm.sp = vm.sb
	vm.popCtx()
	vm.pc++
}
