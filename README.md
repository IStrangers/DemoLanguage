

# DemoLanguage

## 介绍

DemoLanguage 是一个用 Go 语言实现的动态编程语言。该项目包含完整的语言实现，包括词法分析器、语法解析器、抽象语法树（AST）定义、解释器以及基于字节码的虚拟机编译器。

## 软件架构

项目采用经典的编程语言实现架构，主要包含以下组件：

### 核心组件

- **AST（抽象语法树）**：定义了语言的所有语法节点类型，包括语句（Statement）、表达式（Expression）和声明（Declaration）
- **Parser（解析器）**：包含词法分析器（Lexer）和语法解析器（Parser），负责将源代码转换为 AST
- **Interpreter（解释器）**：直接遍历执行 AST 的解释器实现
- **VM（虚拟机）**：基于字节码的编译器，将 AST 编译为字节码后由虚拟机执行

### 运行时支持

- **内置类型**：数组、对象、全局函数等内置类型的实现
- **值类型系统**：支持整数、浮点数、字符串、布尔值、对象、函数等多种值类型
- **作用域管理**：支持词法作用域和闭包
- **异常处理**：支持 try/catch/finally 异常处理机制

### 指令系统

虚拟机支持丰富的指令集，包括：
- 算术运算：加、减、乘、除、取余
- 逻辑运算：与、或、非、比较运算
- 流程控制：条件跳转、循环控制
- 函数调用：函数定义、调用和返回
- 对象操作：属性访问、数组操作

## 语言特性

### 基本语法

```dl
// 变量声明
var x = 10;
var name = "DemoLanguage";

// 函数定义
fun greet(name) {
    println("Hello, " + name);
}

// 箭头函数
var add = (a, b) => a + b;

// 类和对象
class Person {
    init(name) {
        this.name = name;
    }
    
    sayHello() {
        println("Hello, I'm " + this.name);
    }
}
```

### 控制流

```dl
// 条件判断
if (age >= 18) {
    println("Adult");
} else {
    println("Minor");
}

// 循环
for (var i = 0; i < 10; i++) {
    println(i);
}

// Switch 语句
switch (day) {
    case 1:
        println("Monday");
        break;
    default:
        println("Other day");
}
```

### 异常处理

```dl
try {
    riskyOperation();
} catch (e) {
    println("Error: " + e);
} finally {
    cleanup();
}
```

## 安装教程

### 环境要求

- Go 1.20 或更高版本

### 安装步骤

1. 克隆项目到本地：

```bash
git clone https://gitee.com/QQXQQ/DemoLanguage.git
cd DemoLanguage
```

2. 构建项目：

```bash
go build -o demo_language
```

3. 验证安装：

```bash
./demo_language example/example.dl
```

## 使用说明

### 运行脚本

直接运行 `.dl` 文件：

```bash
./demo_language your_script.dl
```

### REPL 模式

项目包含一个 REPL 实现，可以交互式地执行代码：

```bash
go run main.go
```

### 示例脚本

项目提供了多个示例脚本供参考：

- `example/example.dl`：基础语法示例
- `example/vm_example.dl`：虚拟机特性示例

### 调试功能

虚拟机支持指令级别的调试，可以输出字节码执行过程：

```go
// 在代码中启用调试
program.DumpInstructions(func(format string, args ...interface{}) {
    fmt.Printf(format, args...)
})
```

## 项目结构

```
DemoLanguage/
├── ast/                    # 抽象语法树定义
│   └── node.go            # 所有语法节点类型
├── parser/                # 词法分析和语法解析
│   ├── lexer.go           # 词法分析器
│   ├── parser.go          # 语法解析器
│   └── expression.go      # 表达式解析
│   └── statement.go        # 语句解析
├── interpreter/           # 解释器实现
│   ├── interpreter.go     # 解释器主逻辑
│   ├── evaluate_*.go     # 表达式和语句求值
│   └── builtin_*.go       # 内置类型
├── vm/                    # 虚拟机和编译器
│   ├── compiler.go        # 字节码编译器
│   ├── vm.go             # 虚拟机核心
│   ├── instruction.go     # 指令定义
│   └── runtime.go         # 运行时支持
├── file/                  # 文件处理
├── token/                 # 词法标记定义
├── example/               # 示例代码
└── main.go               # 程序入口
```

## 测试

运行项目测试：

```bash
go test ./...
```

## 许可证

本项目遵循 MIT 许可证开源。