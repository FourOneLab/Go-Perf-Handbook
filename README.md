# Go-Perf-Handbook

Golang's performance.

## 深入理解 Golang 数据结构



## 测试工具

Golang 自带的测试工具，输入 `go help test` 会输出详细的参数说明。

```shell
usage: go test [build/test flags] [packages] [build/test flags & test binary flags]
```

`go test` 指令会重新编译包中每一个文件名符合`*_test.go` 模式的文件，这些文件中可以包含功能测试函数、性能测试函数和示例函数，以`_`或者`.`，开头的文件将会被忽略。

以`_test`为后缀的包将被编译为一个单独的包，然后与主测试二进制文件链接并运行。go工具将忽略名为`testdata`的目录，使其可用来保存测试所需的辅助数据。

作为编译测试二进制文件的一部分，`go test`指令会在包和源文件上执行`go vet`命令（检查时只执行默认集合中的一个高可信的子集，包括`atomic`、`bool`、`buildtags`、`errorsas`、`ifaceassert`、`nilfunc`、`printf`、`stringintconv`），执行`go doc cmd/vet`指令获取更多详细信息，通过`-vet=off`来关闭检查。

所有测试的输出和总结都会输出到go命令的标准输出，即使测试代码将输出打印到它自己的标准错误输出。(go命令的标准错误输出用于打印编译测试文件时的错误。)

### 运行模式

1. **本地目录（local directory）**模式：没有指定包级别参数时使用这种模式，如`go test`，`go test -v`。
   1. 此模式下，只会编译本目录下的源文件为测试二进制文件，并运行测试
   2. 此模式下，测试结果是不会被缓存起来的
2. **包列表（package list）**模式：显式的使用包级别参数，如 `go test <package_name>`、`go test ./...`、`go test .`。
   1. 此模式下，会编译列表中的每一个包中的源文件为测试二进制文件，并运行测试（测试成功的包，只会输出 `OK`，测试失败的包，将会打印全部的测试输出）
   2. 此模式下，指定`-bench`、`-v`参数，那么成功与否都会打印全部的测试输出
   3. 此模式下，只要有一个包测试失败，最终的测试结果就是失败
   4. 此模式下，测试成功的包的结果会被缓存起来，这样可以有效的避免不必要的重复执行，输出结果带上（`cacehd`）

> 关于缓存的匹配规则，如下，有任一条件不符合就不会被缓存，也可以使用 `-count=1`来显式禁用缓存：
>
> 1. 运行相同的测试二进制文件
> 2. 运行时指定的参数属于可缓存参数，如，`-benchtime`、` -cpu`、`-list`、`-parallel`、 `-run`、 `-short` 和 `-v`
>
> 缓存的测试结果在任何时候都被认为是正在执行，因此在测试成功的包上设置`-timeout`是不会起作用的。

### 通用参数

除了可以使用 `build` 的一些参数（执行`go help build`获取更详细的参数)，`test` 本身也有一些参数（执行`go help testflag`获取更详细的参数）：

- `-args` ：将这个参数后面的部分都传递给测试二进制文件，这个参数一般放在最后。
- `-c` ：将测试二进制文件编译为`pkg.test`，而不会运行测试文件 (其中`pkg`是包导入路径的最后一个元素)，可以使用`-o`标志修改文件名。
- `-exec xprog` ：使用 `xprog` 运行测试二进制文件，行为和`go run`一样。
- `-i` ：（弃用）安装测试依赖的包，而不会运行测试文件，被编译的包都会自动被缓存起来。
- `-json` ：将测试的输出内容转换为 JSON 格式，运行`go doc test2json`获取关于编码的细节。
- `-o file` ：将测试二进制文件编译到指定的文件中，测试将会继续执行，除非指定了`-c`或`-i`参数。

### 专用参数



### 测试函数类别

注意：不同类型的测试函数，函数的入参是不同的。

```go
// 功能测试函数，函数命名是 TestXxx 
func TestXxx(t *testing.T) { ... }

// 性能测试函数，函数命名是 BenchmarkXxx
func BenchmarkXxx(b *testing.B) { ... }

// 示例函数，函数命名是 ExampleXxx
func ExamplePrintln() {
		Println("The output of\nthis example.")
		// Output: The output of
		// this example.
}

func ExamplePerm() {
		for _, value := range Perm(4) {
			fmt.Println(value)
		}

		// Unordered output: 4
		// 2
		// 1
		// 3
		// 0
}
```

示例函数没有入参，并且会将结果输出到`os.Stdout`：

- 如果示例函数体中最后一个注释是以 `Output:` 开头，那么输出将和注释的内容进行比较。
- 如果示例函数体中最后一个注释是以 `Output:` 开头，但是后面没有其他文本内容，同样会被编译和执行，只是没有输出结果。
- 如果示例函数体中最后一个注释是以`Unordered output:`开头，那么输出将和注释的内容进行比较，但是每一行的顺序将会被打乱。
- 如果示例函数体中没有以上两种类型的注释，那么只会被编译，不会被执行。

## 标准库

测试机的配置：Intel(R) Core(TM) i5-8257U CPU @ 1.40GHz

|类别|标准库|测试函数|性能|
|---|---|---|---|
|字符串|buildin|+||
|字符串|bytes|Buffer||
|字符串|strings|Builder||
|||||
|||||
|||||
|||||



## Golang 并发模式

