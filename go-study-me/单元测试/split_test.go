// split/split_test.go

package split

import (
	"fmt"
	"reflect"
	"testing"
)

func TestSplit(t *testing.T) {
	got := Split("a:b:c", ":")       // 程序返回的结果
	want := []string{"a", "b", "c"}  // 期望返回的结果
	if !reflect.DeepEqual(want, got) {
		// slice 不能直接比较, 借助反射包中方法比较
		t.Errorf("excepted:%v, got:%v", want, got)   //测试失败输出错误提示
	}
}

func TestMultiSplit(t *testing.T) {
	got := Split("abcdgsdgdbcsdfsdfdbc", "bc")
	want := []string{
		"a",
		"dgsdgd",
		"sdfsdfd",
	}
	if !reflect.DeepEqual(want, got) {
		t.Errorf("excepted:%v, got:%v", want, got)   //测试失败输出错误提示
	}
}


// 使用测试组合并多个测试函数+切片
func TestComSpilt(t *testing.T) {
	// 定义一个测试用例类型
	type test struct {
		input string
		sep string
		want []string
	}

	// 定义一个存储测试用例的切片
	tests := []test{
		{input: "a:b:c", sep: ":", want: []string{"a", "b", "c"}},
		{input: "a:b:c", sep: ",", want: []string{"a:b:c"}},
		{input: "abcdgsdgdbcsdfsdfdbc", sep: "bc", want: []string{"a", "dgsdgd", "sdfsdfd"}},
		{input: "可甜可咸可钱可言可", sep: "可", want: []string{"甜", "咸", "钱", "言"}},
	}

	for _, tc := range tests {
		got := Split(tc.input, tc.sep)
		if !reflect.DeepEqual(got, tc.want) {
			t.Errorf("excepted:%#v, got:%#v", tc.want, got)   //测试失败输出错误提示
		}
	}
}



// BenchmarkSplit ...
// func BenchmarkSplit(b *testing.B) {
// 	for i := 0; i < b.N; i++ {
// 		Split("abdfdsfsdfsdfdsfdfdfjkledf", "d")
// 	}
// }

/*
需要注意的是：在调用TestMain时, flag.Parse并没有被调用。
所以如果TestMain 依赖于command-line标志 (包括 testing 包的标记), 则应该显示的调用flag.Parse
*/

// TestMain   Setup and Teardown
// func TestMain(m *testing.M) {
// 	fmt.Println("write setup code here...")
// 	retCode := m.Run()   // 执行测试
// 	fmt.Println("write teardown code here...")
// 	os.Exit(retCode)     // 退出测试
// }


// 测试集的Setup Teardown
func setupTestCase(t *testing.T) func(t *testing.T) {
	t.Log("执行测试之前的coding")
	return func(t *testing.T) {
		t.Log("执行测试的coding")
	}
}


// 子测试的Setup Teardown
func setupSubTestCase(t *testing.T) func(t *testing.T) {
	t.Log("执行子测试之前的coding")
	return func(t *testing.T) {
		t.Log("执行子测试的coding")
	}
}


// map + 子测试
func TestSplitComMapWithSuntest(t *testing.T) {
	// 定义一个测试用例类型
	type test struct {
		input string
		sep string
		want []string
	}

	tests := map[string]test{ // 测试用例使用map存储
		"simple":      {input: "a:b:c", sep: ":", want: []string{"a", "b", "c"}},
		"wrong sep":   {input: "a:b:c", sep: ",", want: []string{"a:b:c"}},
		"more sep":    {input: "abcdgsdgdbcsdfsdfdbc", sep: "bc", want: []string{"a", "dgsdgd", "sdfsdfd"}},
		"leading sep": {input: "沙河有沙又有河", sep: "沙", want: []string{"", "河有", "又有河"}},
	}

	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			teardownSubTest :=  setupSubTestCase(t)
			defer teardownSubTest(t)
			// 使用t.Run执行子测试
			got := Split(tc.input, tc.sep)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("excepted:%#v, got:%#v", tc.want, got)   //测试失败输出错误提示
			}
		})
	}
}

// 示例函数
func Example_Split() {
    fmt.Println(Split("a:b:c", ":"))
	fmt.Println(Split("锅中有肉中有油中有菜", "中"))
	// Output:
	// [a b c]
	// [锅 有肉 有油 有菜]
}
