package vm

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
	exception      *Exception
	callStackDepth int
	sp             int
	catchPos       int
	finallyPos     int
	finallyRet     int
}

type TryStack []TryFrame

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
	vm := &VM{
		runtime: runtime,
		stash: &Stash{
			values:      ValueArray{},
			extraArgs:   ValueArray{},
			nameMapping: make(map[string]uint32),
		},
		maxCallStackSize: 999,
	}
	return vm
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

func (self *VM) handlingThrow(arg any) *Exception {
	return nil
}

func (self *VM) pushTryFrame(catchPos, finallyPos int) {
	self.tryStack.add(TryFrame{})
}

func (self *VM) popTryFrame() {
	self.tryStack = self.tryStack[:len(self.tryStack)-1]
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
