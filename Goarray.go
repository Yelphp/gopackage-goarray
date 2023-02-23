package goarray

/**
	注意：golang 1.18 及以上 ，因为1.18及以后支持泛型
**/
import (
	"encoding/json"
	"errors"
)

// unit 类型
type uintType interface {
	uint | uint16 | uint32 | uint64 | uint8
}

// int 类型
type intType interface {
	int | int16 | int8 | int32 | int64
}

// map 类型
type mapType interface {
	map[string]string | map[string]int | map[string]uint | map[int]int | map[uint]uint
}

// 类型数据
type valuetypedata interface {
	uintType | intType | ~string | interface{} | mapType
}

type keytypedata interface {
	uintType | intType | ~string
}

// goarray 数据结构体
type goarray[k keytypedata, t valuetypedata] struct {
	// 存储map数据
	value map[k]t
	// 记录 key 的位置
	position []k // key值
	// 记录长度
	length int
	// 虚拟指针
	cursor int
}

// 创建一个为空的goarray
func initGoarray[k keytypedata, t valuetypedata]() goarray[k, t] {
	var array = goarray[k, t]{
		length:   0,             // 长度
		cursor:   0,             // 游标
		value:    make(map[k]t), // map 值
		position: []k{},         // map 的键值
	}
	return array
}

// 创建数组
func New[k keytypedata, t valuetypedata](argv map[k]t) goarray[k, t] {
	// 创建一个空的map类型
	array := initGoarray[k, t]()
	// 参数
	for ks, vs := range argv {
		array.Append(ks, vs)
	}
	return array
}

// 创建数组
func NewArray[t valuetypedata](argv []t) goarray[int, t] {
	// 创建一个空的map类型
	array := initGoarray[int, t]()
	// 参数
	for ks, vs := range argv {
		array.Append(ks, vs)
	}
	return array
}

// 从尾部插入数据
func (array *goarray[k, t]) Append(key k, value t) bool {
	// 判断是否存在
	if !array.IsKey(key) {
		// key的位置
		array.position = append(array.position, key)
		// 长度
		array.length++
		// 赋值
		array.value[key] = value
		return true
	}
	return false
}

// 从首为插入数据
func (array *goarray[k, t]) Unshift(key k, value t) bool {
	if !array.IsKey(key) {
		// key的位置
		array.position = append([]k{key}, array.position...)
		// 长度
		array.length++
		// 赋值
		array.value[key] = value
		return true
	}
	return false
}

// 在某个位置插入数据
func (array *goarray[k, t]) Add(place_num int, key k, value t) bool {
	if place_num >= array.length {
		// 位置大于长度，则插入到最后
		return array.Append(key, value)
	} else if place_num < 0 {
		// 位置小于0，则失败
		return false
	} else if place_num == 0 {
		// 位置等于0，则在第一位插入
		return array.Unshift(key, value)
	} else {
		// 把数组进行切割
		var before = append([]k{}, array.position[:place_num]...)
		var after = append([]k{}, array.position[place_num:]...)
		// 然后进行重新组合
		array.position = append(before, key)
		array.position = append(array.position, after...)
		// 长度
		array.length++
		// 赋值
		array.value[key] = value
		return true
	}
}

// 从尾部删除一个数据
func (array *goarray[k, t]) Pop() bool {
	if array.length > 0 {
		// 获取最后一个元素 KEY
		key := array.position[array.length-1:][0]
		// 删除map 中的数据
		delete(array.value, key)
		// 减一个长度
		array.length--
		// 删除最后面的数组
		var new_position = make([]k, array.length)
		// 把 goarray赋值到新的地址
		copy(new_position, array.position[:array.length])
		// 赋值
		array.position = new_position
		return true
	}
	return false
}

// 删除第一个数据
func (array *goarray[k, t]) Shift() bool {
	if array.length > 0 {
		// 获取第一个值 的 key
		key := array.position[:1][0]
		// 删除map的值
		delete(array.value, key)
		// 长度减一
		array.length--
		var new_position = make([]k, array.length)
		copy(new_position, array.position[1:array.length+1])
		array.position = new_position
		return true
	}
	return false
}

// 删除某个值
func (array *goarray[k, t]) Delete(key k) bool {
	if array.IsKey(key) {
		// 删除map的值
		delete(array.value, key)
		// 减位置
		array.length--
		// 删除在map key 中的 array.position
		var new_position = make([]k, array.length)
		var i = 0
		for _, vs := range array.position {
			if vs != key {
				new_position[i] = vs
				i++
			}
		}
		array.position = new_position
		return true
	}
	return false
}

// 获取goarray长度
func (array *goarray[k, t]) Length() int {
	return array.length
}

// 获取某个数据的值
func (array *goarray[k, t]) Value(key k) t {
	return array.value[key]
}

// 判断某个值是否存在
func (array *goarray[k, t]) IsKey(key k) bool {
	_, ok := array.value[key]
	return ok
}

// 按照添加顺序输出Map的Key
func (array *goarray[k, t]) GetOrderMapKey() []k {
	return array.position
}

// 按照添加顺序输出map 的 value
func (array *goarray[k, t]) GetOrderMapValue() []t {
	var mapvalue = make([]t, array.length)
	for ks, vs := range array.position {
		mapvalue[ks] = array.value[vs]
	}
	return mapvalue
}

// 无顺序输出数据
func (array *goarray[k, t]) GetMap() map[k]t {
	return array.value
}

// 把 goarray 转成json
func (array *goarray[k, t]) JsonEndoce() string {
	bytes, _ := json.Marshal(array.value)
	return string(bytes)
}

// 把 json 转成 goarray
func JsonDecode(json_str string) goarray[string, interface{}] {
	// 创建一个map
	maps := make(map[string]interface{})
	// 将字符串转成map
	json.Unmarshal([]byte(json_str), &maps)
	array := New(maps)
	return array
}

// 获取 goarray 切片
// @param int start_num 开始位置
// @param int end_num   结束位置
// @return 返回一个新的 goarray
func (array *goarray[k, t]) SliceArray(start_num int, end_num int) goarray[k, t] {
	// 创建临时空间，存放切片的数据
	slice_arr := array.position[start_num:end_num]
	// 创建一个 goarray 数据
	new_array := goarray[k, t]{
		length:   0,
		value:    make(map[k]t),
		position: []k{},
	}
	for _, vs := range slice_arr {
		new_array.Append(vs, array.value[vs])
	}
	return new_array
}

// 根据键值降序排序
func (array *goarray[string, t]) SortKeyDesc() {
	for {
		swapped := false
		for i := 1; i < array.length; i++ {
			if array.position[i] > array.position[i-1] {
				array.position[i-1], array.position[i] = array.position[i], array.position[i-1]
				swapped = true
			}
		}
		if !swapped {
			return
		}
	}
}

// 根据键值升序排序
func (array *goarray[k, t]) SortKeyAsc() {
	for {
		swapped := false
		for i := 1; i < array.length; i++ {
			if array.position[i] < array.position[i-1] {
				array.position[i-1], array.position[i] = array.position[i], array.position[i-1]
				swapped = true
			}
		}
		if !swapped {
			return
		}
	}
}

// 根据游标，获取当前位置的键与值
func (array *goarray[k, t]) Current() (keys k, value t, err error) {
	if array.cursor < 0 || array.cursor >= array.length {
		array.cursor = 0
		return keys, value, errors.New("游标超出范围")
	} else {
		keys = array.position[array.cursor]
		return keys, array.value[keys], nil
	}
}

// 将游标指向第一位，并获取当前位置的键与值
func (array *goarray[k, t]) Reset() (k, t, error) {
	array.cursor = 0
	return array.Current()
}

// 将指针指向最后一个元素 ，并获取当前位置的键与值
func (array *goarray[k, t]) End() (k, t, error) {
	array.cursor = array.length - 1
	return array.Current()
}

// 将指针指向下一个元素 ， 并获取键与值
func (array *goarray[k, t]) Next() (k, t, error) {
	array.cursor++
	return array.Current()
}

// 将指针指向上一个元素 ， 并获取键与值
func (array *goarray[k, t]) Prev() (k, t, error) {
	array.cursor--
	return array.Current()
}

// 获取键与值，并且将指针向下移动一位
func (array *goarray[k, t]) Each() (k, t, error) {
	key, value, err := array.Current()
	if err != nil {
		array.cursor = 0
	} else {
		array.cursor++
	}
	return key, value, err
}
