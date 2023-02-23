package main

import (
	"fmt"

	goarray "github.com/Yelphp/gopackage-goarray"
)

// 入口函数
func main() {
	// 创建一个map 对象
	var maps = make(map[string]int)
	// New 一个 goarray  ， 注意，创建maps 时已经确定了键值、值的数据类型
	array := goarray.New(maps)
	//从后面插入数据
	array.Append("2第一个", 1)
	array.Append("3第二个", 2)
	array.Append("1第三个", 3)
	array.Append("4第四个", 4)
	array.Append("3第五个", 5)
	// 从前面插入一个数据
	array.Unshift("第六个", 6)
	fmt.Print("根据插入顺序输出map的Key:")
	fmt.Println(array.GetOrderMapKey())
	fmt.Print("根据插入顺序输出map的value:")
	fmt.Println(array.GetOrderMapValue())
	fmt.Print("输出array的长度:")
	fmt.Println(array.Length())
	fmt.Print("获取map值:")
	fmt.Println(array.GetMap())
	fmt.Print("判断第六个是否存在:")
	fmt.Println(array.IsKey("第六个"))
	fmt.Print("从前面输出一个，判断第六个是否存在:")
	array.Shift()
	fmt.Println(array.IsKey("第六个"))
	fmt.Print("判断第五个是否存在:")
	fmt.Println(array.IsKey("第五个"))
	fmt.Print("从后面输出一个，判断第五个是否存在:")
	array.Pop()
	fmt.Println(array.IsKey("第五个"))
	fmt.Println(array.JsonEndoce())

	// 把json 转成 goarray
	var arr = []interface{}{1, "string", 3, 4, 5, 4, 3, 2}
	arr_goarray := goarray.NewArray(arr)
	fmt.Println(arr_goarray.GetOrderMapValue())
	json_str := arr_goarray.JsonEndoce()
	fmt.Println(json_str)
	// 把 json 转成 goarray
	goarray.JsonDecode(json_str)

	// 获取 array[0:2]
	fmt.Println(array.SliceArray(0, 2))
	// 键值排序
	fmt.Println(array)
	array.SortKeyDesc()
	fmt.Println(array)
	array.SortKeyAsc()
	fmt.Println(array)
	array.Delete("2第一个")
	fmt.Println(array)
	array.Add(4, "sad", 123)
	fmt.Println(array)
	for i := 0; i < array.Length(); i++ {
		k, v, e := array.Each()
		if e != nil {
			fmt.Println(e)
			break
		}
		fmt.Printf("key:%v => value : %v  err %v\n", k, v, e)
	}
	k, v, e := array.Next()
	fmt.Printf("key:%v => value : %v  err %v\n", k, v, e)
	// arr_goarray.SortKeyDesc()
	fmt.Println(array)
}
