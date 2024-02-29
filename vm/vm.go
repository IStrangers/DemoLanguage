package vm

import "github.com/istrangers/demolanguage/parser"

type ValueStack ValueArray

func (self *ValueStack) expand(index int) {
	if index < len(*self) {
		return
	}
	index++
	var newCap int
	if index < 1024 {
		newCap = index * 2
	} else {
		newCap = (index + 1025) &^ 1023
	}
	newValueStack := make(ValueStack, index, newCap)
	copy(newValueStack, *self)
	*self = newValueStack
}

type Stash struct {
	outer       *Stash
	values      ValueArray
	extraArgs   ValueArray
	nameMapping map[string]uint32
}

type Context struct {
	program *Program
	pc      int
	sb      int
	args    int
	result  Value
}

type CallStack []Context

func (self *CallStack) size() int {
	return len(*self)
}

func (self *CallStack) add(ctx Context) {
	*self = append(*self, ctx)
}

func (self *CallStack) pop() Context {
	lastIndex := self.size() - 1
	ctx := (*self)[lastIndex]
	*self = (*self)[:lastIndex]
	return ctx
}

type RefStack []Ref

func (self *RefStack) add(ref Ref) {
	*self = append(*self, ref)
}

func (self *RefStack) size() int {
	return len(*self)
}

func (self *RefStack) peek() Ref {
	return (*self)[self.size()-1]
}

func (self *RefStack) pop() Ref {
	lastIndex := self.size() - 1
	ref := (*self)[lastIndex]
	*self = (*self)[:lastIndex]
	return ref
}

type TryFrame struct {
	exception       *Exception
	callStackLength int
	refLength       int
	sp              int
	stash           *Stash
	catchPos        int
	finallyPos      int
	finallyRet      int
}

type TryStack []TryFrame

func (self *TryStack) size() int {
	return len(*self)
}

func (self *TryStack) add(tryFrame TryFrame) {
	*self = append(*self, tryFrame)
}

type VM struct {
	runtime *Runtime
	program *Program

	pc   int
	sp   int
	sb   int
	args int

	maxCallStackSize int

	stash     *Stash
	callStack CallStack
	stack     ValueStack
	refStack  RefStack
	tryStack  TryStack

	result Value
}

func CreateVM() *VM {
	runtime := CreateRuntime()
	return runtime.vm
}

func (self *VM) RunScript(script string) (Value, error) {
	parser := parser.CreateParser(1, "", script, true, true)
	program, err := parser.Parse()
	if err != nil {
		return nil, err
	}
	compiler := CreateCompiler()
	compiler.compile(program)
	self.program = compiler.program
	self.runTry()
	return self.result, nil
}

func (self *VM) run() {
	for {
		if self.pc >= self.getInstructionSize() {
			break
		}
		self.execInstruction(self.pc)
	}
}

func (self *VM) getInstructionSize() int {
	return self.program.getInstructionSize()
}

func (self *VM) getInstruction(pc int) Instruction {
	return self.program.getInstruction(pc)
}

func (self *VM) execInstruction(pc int) {
	instruction := self.getInstruction(pc)
	instruction.exec(self)
}

func (self *VM) setDefining(name string, value Value) {
	stash := self.stash
	if stash != nil {
		if _, exists := stash.nameMapping[name]; !exists {
			index := uint32(stash.values.size())
			stash.values = append(stash.values, value)
			stash.nameMapping[name] = index
		}
		return
	}
	globalObject := self.runtime.globalObject
	globalObject.self.setProperty(name, value)
}

func (self *VM) getDefining(name string) Value {
	stash := self.stash
	for ; stash != nil; stash = stash.outer {
		if index, exists := stash.nameMapping[name]; exists {
			return stash.values[index]
		}
	}
	globalObject := self.runtime.globalObject
	return globalObject.self.getProperty(name)
}

func (self *VM) getValue(index int) Value {
	return self.program.values[index]
}

func (self *VM) push(value Value) {
	self.stack.expand(self.sp)
	self.stack[self.sp] = value
	self.sp++
}

func (self *VM) pop() Value {
	self.sp--
	return self.stack[self.sp]
}

func (self *VM) peek() Value {
	return self.stack[self.sp-1]
}

func (self *VM) clearStack() {
	sp := self.sp
	stackTail := self.stack[sp:]
	for i := range stackTail {
		stackTail[i] = nil
	}
	self.stack = self.stack[:sp]
}

func (self *VM) newStash() *Stash {
	self.stash = &Stash{
		outer: self.stash,
	}
	return self.stash
}

func (self *VM) saveCtx(ctx *Context) {
	ctx.program, ctx.pc, ctx.sb, ctx.args, ctx.result = self.program, self.pc, self.sb, self.args, self.result
}

func (self *VM) restoreCtx(ctx Context) {
	self.program, self.pc, self.sb, self.args, self.result = ctx.program, ctx.pc, ctx.sb, ctx.args, ctx.result
}

func (self *VM) pushCtx() {
	if self.callStack.size() > self.maxCallStackSize {
		//wait adjust
		panic("StackOverflowError")
	}
	ctx := Context{}
	self.saveCtx(&ctx)
	self.callStack.add(ctx)
}

func (self *VM) popCtx() {
	ctx := self.callStack.pop()
	self.restoreCtx(ctx)
}

func (self *VM) halted() bool {
	pc := self.pc
	return pc < 0 || pc >= self.getInstructionSize()
}

func (self *VM) throw(arg any) {
	if ex := self.handlingThrow(arg); ex != nil {
		panic(ex)
	}
}

func (self *VM) formatToException(arg any) (ex *Exception) {
	switch arg := arg.(type) {
	case *Exception:
		ex = arg
	case *Object:
		ex = &Exception{value: arg}
	case Value:
		ex = &Exception{value: arg}
	}
	return
}

func getFunctionName(stack ValueStack, sb int) string {
	if sb > 0 {
		function := stack[sb-1]
		if function.isObject() {
			return function.toObject().self.getPropertyOrDefault("name", Const_Empty_String_Value).toString()
		}
	}
	return ""
}

func (self *VM) captureStack(stackFrameArray StackFrameArray, ctxOffset int) StackFrameArray {
	if self.program != nil || self.sb > 0 {
		var functionName string
		if self.program != nil {
			functionName = self.program.functionName
		} else {
			functionName = getFunctionName(self.stack, self.sb)
		}
		stackFrameArray = append(stackFrameArray, StackFrame{program: self.program, pc: self.pc, functionName: functionName})
	}
	for i := self.callStack.size() - 1; i > ctxOffset-1; i-- {
		stackFrame := self.callStack[i]
		if stackFrame.program != nil || stackFrame.sb > 0 {
			var functionName string
			if stackFrame.program != nil {
				functionName = stackFrame.program.functionName
			} else {
				functionName = getFunctionName(self.stack, stackFrame.sb)
			}
			stackFrameArray = append(stackFrameArray, StackFrame{program: stackFrame.program, pc: stackFrame.pc, functionName: functionName})
		}
	}
	return stackFrameArray
}

func (self *VM) restoreStacks(refLength int) (ex *Exception) {
	//wait adjust
	refTail := self.refStack[refLength:]
	for i := range refTail {
		refTail[i] = nil
	}
	self.refStack = self.refStack[:refLength]
	return
}

func (self *VM) handlingThrow(arg any) *Exception {
	ex := self.formatToException(arg)
	if ex.stack == nil {
		ex.stack = self.captureStack(make(StackFrameArray, 0, self.callStack.size()+1), 0)
	}
	for self.tryStack.size() > 0 {
		tryFrame := &self.tryStack[self.tryStack.size()-1]
		if tryFrame.catchPos == -1 && tryFrame.finallyPos == -1 || ex == nil && tryFrame.catchPos != -2 {
			tryFrame.exception = nil
			self.popTryFrame()
			continue
		}
		if tryFrame.callStackLength < self.callStack.size() {
			context := self.callStack[tryFrame.callStackLength]
			self.program, self.result, self.pc, self.sb, self.args = context.program, context.result, context.pc, context.sb, context.args
			self.callStack = self.callStack[:tryFrame.callStackLength]
		}
		self.sp = tryFrame.sp
		self.stash = tryFrame.stash
		_ = self.restoreStacks(tryFrame.refLength)

		if tryFrame.catchPos == -2 {
			break
		}
		if tryFrame.catchPos >= 0 {
			self.push(ex.value)
			self.pc = tryFrame.catchPos
			tryFrame.catchPos = -1
			return nil
		}
		if tryFrame.finallyPos >= 0 {
			tryFrame.exception = ex
			self.pc = tryFrame.finallyPos
			tryFrame.finallyPos = -1
			tryFrame.finallyRet = -2
			return nil
		}
	}
	if ex == nil {
		panic(arg)
	}
	return ex
}

func (self *VM) pushTryFrame(catchPos, finallyPos int) {
	self.tryStack.add(TryFrame{
		callStackLength: self.callStack.size(),
		refLength:       self.refStack.size(),
		sp:              self.sp,
		stash:           self.stash,
		catchPos:        catchPos,
		finallyPos:      finallyPos,
		finallyRet:      -1,
	})
}

func (self *VM) popTryFrame() {
	self.tryStack = self.tryStack[:self.tryStack.size()-1]
}

func (self *VM) runTryInner() (ex *Exception) {
	//defer func() {
	//	if err := recover(); err != nil {
	//		ex = self.handlingThrow(err)
	//	}
	//}()
	self.run()
	return
}

func (self *VM) runTry() *Exception {
	self.pushTryFrame(-2, -1)
	defer self.popTryFrame()

	for {
		ex := self.runTryInner()
		if ex != nil || self.halted() {
			return ex
		}
	}
}
