# go build 常用技巧

## 1. 如何动态注入go程序的版本等信息？

有时候我们想在go程序中注入编译时间，编译的go版本(多人协同时可能go版本不同)，编译的处理器架构等信息，在进行发布。那一般怎么操作呢？在开源项目中我们可以看到很多这种样例，

1. [k8s 版本信息动态编译配置，version.sh](https://github.com/kubernetes/hack/lib/version.sh)
2. [k8s-release-版本信息](https://github.com/kubernetes/vendor/k8s.io/component-base/version/version.go)
3. [etcd-mackfile](https://github.com/kubernetes/vendor/go.etcd.io/bbolt/Makefile)

以下是k8s中很多组件引用的版本基础包，在打包编译时会动态修改里面的部分变量。

```golang
/*
Copyright 2019 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package version

import (
	"fmt"
	"runtime"

	apimachineryversion "k8s.io/apimachinery/pkg/version"
)

// Get returns the overall codebase version. It's for detecting
// what code a binary was built from.
func Get() apimachineryversion.Info {
	// These variables typically come from -ldflags settings and in
	// their absence fallback to the settings in ./base.go
	return apimachineryversion.Info{
		Major:        gitMajor,
		Minor:        gitMinor,
		GitVersion:   gitVersion,
		GitCommit:    gitCommit,
		GitTreeState: gitTreeState,
		BuildDate:    buildDate,
		GoVersion:    runtime.Version(),
		Compiler:     runtime.Compiler,
		Platform:     fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
}

```

参考这些开源项目来尝试一个gobuild的案例。

```golang
package main

import (
	"fmt"
	"os"
	"runtime"
)

var buildtime = ""

func main() {
	args := os.Args
	if len(args) == 2 && (args[1] == "--version" || args[1] == "-v") {
		fmt.Printf("Build Time: %s\n", buildtime)
		fmt.Printf("go version: %s\n", runtime.Version())
		fmt.Printf("Platform: %s:%s\n", runtime.GOOS, runtime.GOARCH)
	}
}

```

编译命令

```bash
# build/sh

#!/usr/bin/env bash

buildtime="$(date -u '+%Y-%m-%d %I:%M:%S%p')"
BRANCH=`git rev-parse --abbrev-ref HEAD`
COMMIT=`git rev-parse --short HEAD`
GOLDFLAGS="-s -w -X 'main.buildtime=$buildtime' -X 'main.branch=$BRANCH' -X 'main.commit=$COMMIT'"
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "$GOLDFLAGS" -o hello main.go
```

输出以下内容

```shell
➜  gobuild git:(master) ./build.sh                         
➜  gobuild git:(master) ./hello -v                         
Build Time: 2022-04-16 05:16:03PM
GitCommit: master:4dd79e7
go version: go1.17
Platform: linux:amd64
```

这样看着好像没毛病，但是这样也有一个问题：**就是在交叉编译的时候无法正确反应出 go 的版本**。比如，你是在 OSX 下编译 linux 的可执行程序，这时候你通过 `-v` 参数查看显示的也是 linux 平台，而不是期待的 darwin 平台。

我们把构建命令改为如下，在linux下编译windows的可执行文件

```shell
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags "$GOLDFLAGS" -o hello main.go
```

```shell
➜  gobuild git:(master) ./build.sh
➜  gobuild git:(master) ./hello.exe -v
Build Time: 2022-04-16 05:27:13PM
GitCommit: master:4dd79e7
go version: go1.17
Platform: windows:amd64
```

我们发现，编译的Platform并不是我们实际的linux， 由于是用runtimel来获取的信息。我们修改一下go version的获取方式。

```golang
package main

import (
	"fmt"
	"os"
)

var (
	buildtime = ""
	branch    = ""
	commit    = ""
	goversion = ""
)

func main() {
	args := os.Args
	if len(args) == 2 && (args[1] == "--version" || args[1] == "-v") {
		fmt.Printf("Build Time: %s\n", buildtime)
		fmt.Printf("GitCommit: %s:%s\n", branch, commit)
		fmt.Printf("GO Version: %s\n", goversion)
		// fmt.Printf("Platform: %s/%s\n", runtime.GOOS, runtime.GOARCH)
	}
}

```

```shell
#!/usr/bin/env bash

buildtime="$(date -u '+%Y-%m-%d %I:%M:%S%p')"
BRANCH=`git rev-parse --abbrev-ref HEAD`
COMMIT=`git rev-parse --short HEAD`
GOVERSION=`go version`
GOLDFLAGS="-s -w -X 'main.buildtime=$buildtime' -X 'main.branch=$BRANCH' -X 'main.commit=$COMMIT' -X 'main.goversion=$GOVERSION'"
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags "$GOLDFLAGS" -o hello.exe main.go
```

可以看到输出的内容是所预期的。

```shell
➜  gobuild git:(master) ./build.sh    
➜  gobuild git:(master) ./hello.exe -v
Build Time: 2022-04-16 05:31:37PM
GitCommit: master:4dd79e7
GO Version: go version go1.17 linux/amd64
```

## 2. 如何构建最小go可执行文件？

在第一阶段中我们用以下这段编译命令，如何理解以下的各个参数？这么构建有什么好处？

```bash
buildtime="$(date -u '+%Y-%m-%d %I:%M:%S%p')"
goversion="$(go version)"
goarch=amd64
flags="-s -w -extldflags '-static' -X 'main.buildtime=$buildtime' -X 'main.goversion=$goversion'"
CGO_ENABLED=0 GOOS=linux GOARCH=$goarch go build -a -ldflags "$flags" -o hello main.go
```

- `CGO_ENABLED` 

  默认情况下，Go的runtime环境变量CGO_ENABLED=1，即默认开始cgo，允许你在Go代码中调用C代码，如果标准库中是在CGO_ENABLED=1情况下编译的，那么编译出来的最终二进制文件可能是动态链接，所以建议设置 CGO_ENABLED=0以避免移植过程中出现的不必要问题。

- `-ldflags`  [sets the flags that are passed to 'go tool link'](https://golang.org/cmd/go/#hdr-Compile_packages_and_dependencies)

  - `-s` 忽略符号表和调试信息，使用该参数后可执行文件无法使用gdb进行调试
  - `-w` 忽略DWARF符号表
  - `-X` 根据指定的路径，动态注入变量值  add string value definition of the form importpath.name=value

  通过这两个参数，可以进一步减少编译的程序的尺寸，更多的参数可以参考[go link](https://golang.org/cmd/link/), 或者 `go tool link -help`(另一个有用的命令是`go tool compile -help`)

  - `-extldflags '-static'`  完全静态编译go程序，无第三方依赖库

- `-a` 它强制重新编译相关的包,一般不需要使用



## 3. 如何完全静态编译一个Go程序？

在docker化的今天，我们一般都是只要一个可执行文件，且无引用其他第三方依赖。在golang中标准库`net` 会使用静态链接库， 依赖glibc等库，如果我们设置`CGO_ENABLED=1`  则编译的可执行文件在镜像内可能就会报错 

> sh: /app: not found

所以我们需要编译一个完全静态的可执行文件，或者在基础镜像中加入glibc库。有以下几种解决方法

- 1. 设置 `CGO_ENABLED=0`
- 2. 编译是使用纯go的net: `go build -tags netgo -a -v` 
- 3. 使用基础镜像加glibc(或等价库musl、uclibc)， 比如 [busybox:glibc](https://hub.docker.com/_/busybox/)、alpine + `RUN apk add --no-cache libc6-compat`、[frolvlad/alpine-glibc](https://hub.docker.com/r/frolvlad/alpine-glibc/)

如果代码中代码中确实必须使用CGO，因为需要依赖一些C/C++的库。目前没有对应的Go库可替代， 那么可以使用`-extldflags "-static"` 来完全静态编译。 ` go tool link help`介绍了`extldflags`的功能：

> -extldflags flags
> 	Set space-separated flags to pass to the external linker.
>
> -static means do not link against shared libraries



## 4. 命令

1. 查看golang支持交叉编译的架构

> go tool dist list



## 5. reference

1. [what-do-these-go-build-flags-mean-netgo-extldflags-lm-lstdc-static](https://stackoverflow.com/questions/37630274/what-do-these-go-build-flags-mean-netgo-extldflags-lm-lstdc-static)
2. [鸟窝-创建最小的go镜像](https://colobu.com/2018/08/13/create-minimal-docker-image-for-go-applications/)
3. [Golang -ldflags 的一个技巧 go version 信息注入](https://ms2008.github.io/2018/10/08/golang-build-version/)

