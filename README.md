# Go-Perf-Handbook

Golang's performance.

## 深入理解 Golang 数据结构

## 测试工具

Golang 自带的测试工具，输入 `go help test` 会输出详细的参数说明。

```shell
usage: go test [build/test flags] [packages] [build/test flags & test binary flags]
```

`go test` 指令会重新编译包中每一个文件名符合 `*_test.go` 模式的文件，这些文件中可以包含功能测试函数、性能测试函数和示例函数，以`_`或者`.`，开头的文件将会被忽略。

以 `_test` 为后缀的包将被编译为一个单独的包，然后与主测试二进制文件链接并运行。go工具将忽略名为 `testdata` 的目录，使其可用来保存测试所需的辅助数据。

作为编译测试二进制文件的一部分，`go test` 指令会在包和源文件上执行 `go vet` 命令（检查时只执行默认集合中的一个高可信的子集，包括`atomic`、`bool`、`buildtags`、`errorsas`、`ifaceassert`、`nilfunc`、`printf`、`stringintconv`），执行`go doc cmd/vet`指令获取更多详细信息，通过 `-vet=off` 来关闭检查。

所有测试的输出和总结都会输出到go命令的标准输出，即使测试代码将输出打印到它自己的标准错误输出。(go命令的标准错误输出用于打印编译测试文件时的错误。)

### 运行模式

1. 本地目录（local directory）模式：没有指定包级别参数时使用这种模式，如 `go test`，`go test -v`。
   1. 此模式下，只会编译本目录下的源文件为测试二进制文件，并运行测试
   2. 此模式下，测试结果是不会被缓存起来的
2. 包列表（package list）模式：显式的使用包级别参数，如 `go test <package_name>`、`go test ./...`、`go test .`。
   1. 此模式下，会编译列表中的每一个包中的源文件为测试二进制文件，并运行测试（测试成功的包，只会输出 `OK`，测试失败的包，将会打印全部的测试输出）
   2. 此模式下，指定 `-bench`、`-v` 参数，那么成功与否都会打印全部的测试输出
   3. 此模式下，只要有一个包测试失败，最终的测试结果就是失败
   4. 此模式下，测试成功的包的结果会被缓存起来，这样可以有效的避免不必要的重复执行，输出结果带上（`cacehd`）

> 关于缓存的匹配规则，如下，有任一条件不符合就不会被缓存，也可以使用 `-count=1` 来显式禁用缓存：
>
> 1. 运行相同的测试二进制文件
> 2. 运行时指定的参数属于可缓存参数，如，`-benchtime`、`-cpu`、`-list`、`-parallel`、 `-run`、 `-short` 和 `-v`
>
> 缓存的测试结果在任何时候都被认为是正在执行，因此在测试成功的包上设置 `-timeout` 是不会起作用的。

### 通用参数

除了可以使用 `build` 的一些参数（执行 `go help build` 获取更详细的参数)，`test` 本身也有一些参数（执行 `go help testflag` 获取更详细的参数）：

- `-args` ：将这个参数后面的部分都传递给测试二进制文件，这个参数一般放在最后。
- `-c` ：将测试二进制文件编译为 `pkg.test`，而不会运行测试文件 (其中 `pkg` 是包导入路径的最后一个元素)，可以使用 `-o` 标志修改文件名。
- `-exec xprog` ：使用 `xprog` 运行测试二进制文件，行为和 `go run `一样。
- `-i` ：（弃用）安装测试依赖的包，而不会运行测试文件，被编译的包都会自动被缓存起来。
- `-json` ：将测试的输出内容转换为 JSON 格式，运行 `go doc test2json 获取关于编码的细节。
- `-o file` ：将测试二进制文件编译到指定的文件中，测试将会继续执行，除非指定了 `-c` 或 `-i` 参数。

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

示例函数没有入参，并且会将结果输出到 `os.Stdout`：

- 如果示例函数体中最后一个注释是以 `Output:` 开头，那么输出将和注释的内容进行比较。
- 如果示例函数体中最后一个注释是以 `Output:` 开头，但是后面没有其他文本内容，同样会被编译和执行，只是没有输出结果。
- 如果示例函数体中最后一个注释是以 `Unordered output:` 开头，那么输出将和注释的内容进行比较，但是每一行的顺序将会被打乱。
- 如果示例函数体中没有以上两种类型的注释，那么只会被编译，不会被执行。

## 标准库

- Golang版本：go1.17.2
- CPU：Intel(R) Xeon(R) Platinum 8255C CPU @ 2.50GHz

### 字符串

生成由1000个字符组成的字符串：

- buildin：使用 `+` 进行字符串相加
- bytes：使用 `bytes.Buffer` 的 `WriteString()` 函数
- strings：使用 `strings.Builder` 的 `WriteString()` 函数

#### 性能

[示例代码](./std/string_concateation/string_concateation_test.go)

```shell
BenchmarkConcatString-2          6898798               176.5 ns/op           530 B/op        0 allocs/op
BenchmarkConcatBuffer-2         157655160              7.764 ns/op           2 B/op          0 allocs/op
BenchmarkConcatBuilder-2        404895285              2.822 ns/op           2 B/op          0 allocs/op
```

### 数字转换

使用标准库 `strconv` 解析 `bool`,`int64`,`float64`。

#### 性能

[示例代码](./std/numeric_conversions/numeric_conversions_test.go)

```shell
BenchmarkParseBool-2            1000000000              0.4962 ns/op           0 B/op          0 allocs/op
BenchmarkParseInt-2             44626659                27.66 ns/op            0 B/op          0 allocs/op
BenchmarkParseFloat-2           25894917                45.81 ns/op            0 B/op          0 allocs/op
```

### 正则表达式

使用标准库 `regexp` 比较是否编译对正则表达式匹配的性能影响。

#### 性能

[示例代码](./std/regular_expressions/regular_expressions_test.go)，是一个电子邮箱匹配的正则表达式。

```shell
BenchmarkMatchString-2                    165367              6859 ns/op            4994 B/op         60 allocs/op
BenchmarkMatchStringCompiled-2           2206146             532.7 ns/op               0 B/op          0 allocs/op
```

### 排序

使用标准库 `sort` 中的算法对1千、1万和10万整数进行排序，时间复杂度是 `o(n*log(n))`。

#### 性能

[示例代码](./std/sorting/sorting_test.go)

```shell
BenchmarkSort1000-2                12427             95886 ns/op              24 B/op          1 allocs/op
BenchmarkSort10000-2                 991           1234771 ns/op              24 B/op          1 allocs/op
BenchmarkSort100000-2                 79          15146117 ns/op              24 B/op          1 allocs/op
```

### 生成随机的数字

使用标准库 `math/rand` 和 `crypto/rand` 生成伪随机数。

#### 性能

[示例代码](./std/random_numbers/random_numbers_test.go)

```shell
BenchmarkMathRand-2             71665480                16.44 ns/op            0 B/op          0 allocs/op
BenchmarkCryptoRand-2             899103                 1318 ns/op           56 B/op          4 allocs/op
```
### 生成随机的字符串

使用标准库 `math/rand` 和 `crypto/rand` 生成长度为16的均匀分部的字符串。

#### 性能

[示例代码](./std/random_strings/random_strings_test.go)

```shell
BenchmarkMatchRandString-2      12583370                93.01 ns/op           32 B/op          2 allocs/op
BenchmarkCryptoRandString-2      1000000                 1009 ns/op           32 B/op          2 allocs/op
```

### 向 `Slice` 中添加元素

向一个 `byte` 切片中添加元素，比较是否触发扩容对性能的影响。

#### 性能

[示例代码](./std/slice_appending/slice_appending_test.go)

```shell
BenchmarkSliceAppend-2                  610151952                1.999 ns/op           5 B/op          0 allocs/op
BenchmarkSliceAppendPreAlloc-2          1000000000               1.073 ns/op           0 B/op          0 allocs/op
```

### 读取 `Map`

比较以 `int` 类型为key和以 `string` 类型为key时，从Map中读取值时的性能差异。
#### 性能

[示例代码](./std/map_access/map_access_test.go)

```shell
BenchmarkMapStringKeys-2        11060889               127.7 ns/op             0 B/op          0 allocs/op
BenchmarkMapIntKeys-2           22202511               50.52 ns/op             0 B/op          0 allocs/op
```

### 对象创建

比较新建对象和使用对象池(`sync.Pool`)的性能差异。

#### 性能

[示例代码](./std/object_creation/object_creation_test.go)

```shell
BenchmarkNoPool-2       24952729                49.64 ns/op           64 B/op          1 allocs/op
BenchmarkPool-2         78339351                15.11 ns/op            0 B/op          0 allocs/op
```

### 哈希函数

使用标准库 `crypto` 和实验库 `"golang.org/x/crypto` 中的多个哈希函数进行对比。

#### 性能

[示例代码](./std/hash_functions/hash_functions_test.go)

```shell
BenchmarkMD5-2                    657549              1809 ns/op              16 B/op          1 allocs/op
BenchmarkSHA1-2                   799243              1541 ns/op              24 B/op          1 allocs/op
BenchmarkSHA256-2                 262231              5359 ns/op              32 B/op          1 allocs/op
BenchmarkSHA512-2                 446698              2687 ns/op              64 B/op          1 allocs/op
BenchmarkSHA3256-2                265435              4656 ns/op             512 B/op          3 allocs/op
BenchmarkSHA3512-2                132931              9650 ns/op             576 B/op          3 allocs/op
BenchmarkBLAKE2b256-2             703269              1554 ns/op              32 B/op          1 allocs/op
BenchmarkBLAKE2b512-2             762640              1500 ns/op              64 B/op          1 allocs/op
```

### base64编解码

使用标准库 `encoding/base64` 对1千字节的数据进行编解码。

#### 性能

[示例代码](./std/base64/base64_test.go)

```shell
BenchmarkEncode-2         572259              1953 ns/op            2816 B/op          2 allocs/op
BenchmarkDecode-2         607777              2029 ns/op            2560 B/op          2 allocs/op
```

### 文件读写

比较读写1MB文本文件时，是否使用缓冲对性能的影响，使用标准库 `bufio` 作为缓冲I/O。

#### 性能

[示例代码](./std/file_io/file_io_test.go)

```shell
BenchmarkWriteFile-2                   6         180701948 ns/op             121 B/op          3 allocs/op
BenchmarkWriteFileBuffered-2         138           7536411 ns/op            4216 B/op          4 allocs/op
BenchmarkReadFile-2                   13          89048980 ns/op             120 B/op          3 allocs/op
BenchmarkReadFileBuffered-2          228           5225005 ns/op         1604224 B/op     100004 allocs/op
```

### 序列化

比较 `json`、`protobuf` 和 `gob` 在序列化和反序列化时的性能。

#### 性能

[示例代码](./std/serialization/serialization_test.go)

```shell
BenchmarkJSONMarshal-2           2294682               531.1 ns/op           144 B/op          1 allocs/op
BenchmarkJSONUnmarshal-2          541806                2261 ns/op           336 B/op         10 allocs/op
BenchmarkPBMarshal-2             5426214               222.0 ns/op            96 B/op          1 allocs/op
BenchmarkPBUnmarshal-2           2553306               478.9 ns/op           144 B/op          6 allocs/op
BenchmarkGobMarshal-2            2805680               424.7 ns/op            24 B/op          1 allocs/op
BenchmarkGobUnmarshal-2         1000000000         0.0000656 ns/op             0 B/op          0 allocs/op
```

### 压缩

使用标准库 `compress/gzip`，比较压缩和解压缩 100KB 的 JSON 格式的数据时的性能差异。

#### 性能

[示例代码](./std/compression/compression_test.go)

```shell
BenchmarkWrite-2             336           3537089 ns/op            2421 B/op          0 allocs/op
BenchmarkRead-2              970           1333251 ns/op         5862914 B/op         29 allocs/op
```

### URL 解析

查看使用标准库 `net/url` 进行URL解析的性能。

#### 性能

[示例代码](./std/url_parseing/url_parsing_test.go)

```shell
BenchmarkParse-2         2643684               432.8 ns/op           144 B/op          1 allocs/op
```

### 模板渲染

查看使用标准库 `text/template` 进行模板渲染的性能。

#### 性能

[示例代码](./std/templates/templates_test.go)

```shell
BenchmarkExecute-2        496792              2417 ns/op             160 B/op         11 allocs/op
```

### HTTP 服务器

使用标准库 `net/http`，比较 `HTTP` 和 `HTTPS` 协议下服务器的性能差异。

#### 性能

[示例代码](./std/http_server/http_server_test.go)

```shell
BenchmarkHTTP-2                    31004             39801 ns/op            5712 B/op         64 allocs/op
BenchmarkHTTPNoKeepAlive-2          7567            173512 ns/op           17872 B/op        139 allocs/op
BenchmarkHTTPSNoKeepAlive-2          100          11610812 ns/op          189822 B/op       1286 allocs/op
```
## Golang 并发模式

## Golang 高级优化技巧

### 内存对齐

为什么要关心内存对齐：

1. 正在编写的代码在性能（CPU、Memory）方面有一定的要求
2. 正在处理向量方面的指令
3. 某些硬件平台（ARM）体系不支持未对齐的内存访问

为什么要做内存对齐：

1. 平台（移植性）原因：不是所有的硬件平台都能够访问任意地址上的任意数据。例如：特定的硬件平台只允许在特定地址获取特定类型的数据，否则会导致异常情况
2. 性能原因：若访问未对齐的内存，将会导致 CPU 进行两次内存访问，并且要花费额外的时钟周期来处理对齐及运算。而本身就对齐的内存仅需要一次访问就可以完成读取动作

#### 内存布局

```go
type Foo struct {
  A int8 // 1
  B int8 // 1 
  C int8 // 1
}

var f Foo
fmt.Println(unsafe.Sizeof(f)) // 3

type Bar struct {
  x int32 // 4
  y *Foo  // 8 (64位处理器)
  z bool  // 1
}

var b Bar
fmt.Println(unsafe.Sizeof(b)) // 24
```

从上面的示例，乍一看，变量 `b` 的内存大小应该是 13，但其实是 24。这是因为 Go 编译器会按照一定的规则自动进行内存对齐。这样设计是为了减少 CPU 访问内存的次数，从而加大 CPU 的吞吐量。如果不进行对其的话，很可能会增加 CPU 访问内存的次数。

> 因为，CPU 在访问内存时，是按照字长来访问的（64位的处理器，字长是8个字节），所以，CPU 每次访问内存的单位就是8字节，每次加载内存数据可以是若干个字长，也就是8字节的整数倍。如果不进行内存对齐，那么在访问某个结构体时，可能会出现某个字段跨一个字长的情况，此时就需要读取两次内存了。 

### 对齐系数

```go
var b1 Bar

// 结构体变量的对齐系数
fmt.Println(unsafe.Alignof(b1)) // 8

// 结构体变量中每个字段的对齐系数
fmt.Println(unsafe.Alignof(b1.x)) // 4
fmt.Println(unsafe.Alignof(b1.y)) // 8
fmt.Println(unsafe.Alignof(b1.z)) // 1
```

对齐系数，表示这个变量需要按照对齐系数的整数倍进行对齐。

`unsafe.Alignof()`函数的规则：

1. 任意的变量，对齐系数至少为1
2. 结构体变量，对齐系数是所有变量中对齐系数最大的那个变量的对齐系数
3. 数组类型变量，数组中元素类型的对齐系数的倍数

因此，可以通过调整结构体字段的顺序，来降低占用的内存，如交换y和z的位置。

```go
type Bar struct {
  x int32 // 4
  z bool  // 1
  y *Foo  // 8 (64位处理器)
}

var b2 Bar
fmt.Println(unsafe.Sizeof(b2)) // 16
```

### 特殊场景

#### 空结构体字段对齐

如果结构或数组中不包含 `size` 大于零的字段（或元素），则其大小为0。两个不同的0大小变量在内存中可能有相同的地址。

由于空结构体 `struct{}` 的大小为 0:

- 当结构体中包含空结构体类型的字段时，通常不需要进行内存对齐,
- 当空结构体类型作为结构体的最后一个字段时，如果有指向该字段的指针，那么就会返回该结构体之外的地址。为了避免内存泄露会额外进行一次内存对齐。

```go
type Demo1 struct {
  m struct{} // 0
  n int8     // 1
}

var d1 Demo1
fmt.Println(unsafe.Sizeof(d1))  // 1

type Demo2 struct {
  n int8     // 1
  m struct{} // 0
}

var d2 Demo2
fmt.Println(unsafe.Sizeof(d2))  // 2
```

访问 `d2.m` 可能会造成内存泄露，因此会进行一次内存对齐。

在实际编程中通过灵活应用空结构体大小为0的特性能够帮助我们节省很多不必要的内存开销。

- 使用空结构体作为map的值来实现一个类似 Set 的数据结构
- 使用空结构体作为通知类 channel 的元素

#### 原子操作在32位平台要求强制内存对齐

在 x86 平台上原子操作需要强制内存对齐是因为在 32bit 平台下进行 64bit 原子操作要求必须 8 字节对齐，否则程序会 panic。

```go
// src/atomic/doc.go

// BUG(rsc): On 386, the 64-bit functions use instructions unavailable before the Pentium MMX.
//
// On non-Linux ARM, the 64-bit functions use instructions unavailable before the ARMv6k core.
//
// On ARM, 386, and 32-bit MIPS, it is the caller's responsibility
// to arrange for 64-bit alignment of 64-bit words accessed atomically.
// The first word in a variable or in an allocated struct, array, or slice can
// be relied upon to be 64-bit aligned.
```

#### false sharing

结构体内存对齐除了上面的场景外，在一些需要防止 CacheLin e伪共享的时候，也需要进行特殊的字段对齐。例如 `sync.Pool` 中就有这种设计：

```go
type poolLocal struct {
  poolLocalInternal

  // Prevents false sharing on widespread platforms with
  // 128 mod (cache line size) = 0 .
  pad [128 - unsafe.Sizeof(poolLocalInternal{})%128]byte
}
```
结构体中的pad字段就是为了防止 false sharing 而设计的。

> 当不同的线程同时读写同一个 cache line 上不同数据时就可能发生 false sharing。false sharing 会导致多核处理器上严重的系统性能下降。

#### hot path

hot path 是指执行非常频繁的指令序列。

在访问结构体的第一个字段时，我们可以直接使用结构体的指针来访问第一个字段（结构体变量的内存地址就是其第一个字段的内存地址）。

如果要访问结构体的其他字段，除了结构体指针外，还需要计算与第一个值的偏移(calculate offset)。在机器码中，偏移量是随指令传递的附加值，CPU 需要做一次偏移值与指针的加法运算，才能获取要访问的值的地址。因为，访问第一个字段的机器代码更紧凑，速度更快。

下面的代码是标准库 sync.Once 中的使用示例，通过将常用字段放置在结构体的第一个位置上减少 CPU 

```go
// src/sync/once.go 

// Once is an object that will perform exactly one action.
//
// A Once must not be copied after first use.
type Once struct {
  // done indicates whether the action has been performed.
  // It is first in the struct because it is used in the hot path.
  // The hot path is inlined at every call site.
  // Placing done first allows more compact instructions on some architectures (amd64/386),
  // and fewer instructions (to calculate offset) on other architectures.
  done uint32
  m    Mutex
}
```

### Goroutine

####  泄漏检测器

> 具有监控存活的 goroutine 数量功能的 APM (Application Performance Monitoring) 应用程序性能监控可以轻松查出 goroutine 泄漏。goroutine 泄漏会导致内存中存活的 goroutine 数量不断上升，直到服务宕机为止。因此，可以在代码部署之前，通过一些方法来检查程序中是否存在泄漏。

[goleak](https://github.com/uber-go/goleak) 可以监控当前测试代码中泄漏的 goroutine，由Uber开源可与单元测试结合使用。启用泄漏检测的唯一要求就是在测试代码结束之前，调用 goleak 库来检测泄漏的 goroutine。事实上，goleak 检测了所有的 goroutine 而不是只检测泄漏的 goroutine。

goleak是基于标准库`runtime.Stack`获取到 goroutine 的堆栈信息的。goleak 解析所有的 goroutine 出并通过以下规则过滤 go 标准库中产生的 goroutine：

1. 由 go test 创建来运行测试逻辑的 goroutine。例如上图中的第二个 goroutine
2. 由 runtime 创建的 goroutine，例如监听信号接收的 goroutine，参阅 [Go: gsignal, Master of goroutine](https://medium.com/a-journey-with-go/go-gsignal-master-of-signals-329f7ff39391)
3. 当前运行的 goroutine，例如上图的第一个 goroutine

[示例代码](goroutine/goleak/example.go),输出结果如下所示：

```bash
Running tool: /usr/local/go/bin/go test -timeout 30s -run ^Test_leak$ example/goroutine/goleak

=== RUN   Test_leak
=== RUN   Test_leak/leak_goroutine
=== CONT  Test_leak
    /home/ubuntu/Go-Perf-Handbook/goroutine/goleak/leaks.go:78: found unexpected goroutines:
        [Goroutine 20 in state sleep, with time.Sleep on top of the stack:
        goroutine 20 [sleep]:
        time.Sleep(0xdf8475800)
                /usr/local/go/src/runtime/time.go:193 +0x12e
        example/goroutine/goleak.leak.func1()
                /home/ubuntu/Go-Perf-Handbook/goroutine/goleak/example.go:7 +0x25
        created by example/goroutine/goleak.leak
                /home/ubuntu/Go-Perf-Handbook/goroutine/goleak/example.go:6 +0x25
        ]
--- FAIL: Test_leak (0.44s)
    --- PASS: Test_leak/leak_goroutine (0.00s)
FAIL
FAIL    example/goroutine/goleak        0.446s
```

从报错信息中可以看到详细的 Goroutine 堆栈信息，可以快速调试并了解发生泄露的 Goroutine。

goleak存在的缺陷：

1. 三方库或者运行在后台中，遗漏的 goroutine 将会造成虚假的结果(无 goroutine 泄漏)
2. 如果在其他未使用 goleak 的测试代码中使用了 goroutine，那么泄漏结果也是错误的。如果这个 goroutine 一直运行到下次使用 goleak 的代码， 则结果也会被这个 goroutine 影响，发生错误。