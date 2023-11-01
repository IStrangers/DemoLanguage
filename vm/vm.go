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

type TryFrame struct {
	exception      *Exception
	callStackDepth uint
	sp             int
	catchPos       int
	finallyPos     int
	finallyRet     int
}

type TryFrameArray []TryFrame

func (self *TryFrameArray) add(tryFrame TryFrame) {
	*self = append(*self, tryFrame)
}

type VM struct {
	program *Program
	stack   ValueStack
	pc      int
	sp      int

	tryStack TryFrameArray
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

func (self *VM) getValue(index uint) Value {
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
	defer func() {
		if err := recover(); err != nil {
			ex = self.handlingThrow(err)
		}
	}()
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
