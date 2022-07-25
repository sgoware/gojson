package gojson

import (
	"errors"
	"gojson/internal/conv"
	"gojson/internal/type/stringx"
	"reflect"
	"strings"
)

func (j *Json) findContentPointer(pattern string) *interface{} {
	if pattern == "" {
		return nil
	}
	// "."返回全部内容
	if pattern == "." {
		return j.jsonContent
	}
	pointer := j.jsonContent
	nodes := strings.Split(pattern, ".")
	for _, n := range nodes {
		if stringx.IsIndex(n) {
			// 在数组或切片中寻找
			if arr, ok := (*pointer).([]interface{}); ok {
				i, err := stringx.GetIndex(n)
				if err != nil {
					return nil
				}
				arrLen := len(arr)
				if arrLen == 0 ||
					i > arrLen-1 {
					return nil
				}
				return &arr[i]
			}
		} else {
			// 在map中寻找

			if mp, ok := (*pointer).(map[string]interface{}); ok {
				mapValue, ok := mp[n]
				if !ok {
					return nil
				}
				pointer = &mapValue
			}

			if mp, ok := (*pointer).(map[string][]interface{}); ok {
				var result interface{}
				mapValue, ok := mp[n]
				if !ok {
					return nil
				}
				result = &mapValue
				pointer = &result
			}
		}
	}
	return pointer
}

func (j *Json) setContentWithOptions(pattern string, data interface{}, options Options) error {
	var (
		err     error
		content interface{} = nil
	)
	if pattern == "" {
		return errors.New(invalidPattern)
	}

	// 先将data转换成map格式
	if data != nil {
		content, err = j.convertContent(data, options)
		if err != nil {
			return err
		}
	}

	// "."替换全部内容
	if pattern == "." {
		j.jsonContent = &content
		return nil
	}
	nodes := strings.Split(pattern, ".")
	nodesLength := len(nodes)
	// 设置内容时,如果jsonContent为空指针需要初始化
	// 防止删除jsonContent后指针为空的情况
	if *j.jsonContent == nil {
		if stringx.IsIndex(nodes[0]) {
			*j.jsonContent = make([]interface{}, 0)
		} else {
			*j.jsonContent = make(map[string]interface{})
		}
	}
	var parentPointer *interface{} = nil // 当前节点的父节点的指针
	curPointer := j.jsonContent          // 当前节点的指针
	// 开始遍历path的各个节点
	for i := 0; i < nodesLength; i++ {
		switch (*curPointer).(type) {
		// 枚举*currentPointer的种类

		case map[string]interface{}:
			// 是map[string]interface{}的情况

			if i != nodesLength-1 {
				// 不是叶子节点的情况

				if val, ok := (*curPointer).(map[string]interface{})[nodes[i]]; ok {
					// 如果map的key:nodes[i]有对应的val,直接将这个val替换
					parentPointer = curPointer
					curPointer = &val
				} else {
					// 如果map的key:nodes[i]没有对应的val,创建一个新的节点

					if content == nil {
						// 传入的内容为空,但是当前节点没有对应的内容,所以直接返回空

						return nil
					}

					if stringx.IsIndex(nodes[i+1]) {
						chileNodeIndex, _ := stringx.GetIndex(nodes[i+1])
						var v interface{} = make([]interface{}, chileNodeIndex+1)
						parentPointer = j.setPointer(curPointer, nodes[i], v)
						curPointer = &v
					} else {
						var v interface{} = make(map[string]interface{})
						parentPointer = j.setPointer(curPointer, nodes[i], v)
						curPointer = &v
					}
				}
			} else {
				// 是叶子节点的情况

				if content == nil {
					// 传入的值是空内容意思是删除

					delete((*curPointer).(map[string]interface{}), nodes[i])
				} else {
					if (*curPointer).(map[string]interface{}) == nil {
						// 如果map为空,创建一个新的map

						*curPointer = map[string]interface{}{}
					}
					// 叶子节点直接替换值

					(*curPointer).(map[string]interface{})[nodes[i]] = content
				}
			}

		case []interface{}:
			// 是空接口切片的情况
			if !stringx.IsIndex(nodes[i]) {
				// 当前节点的key不是索引的情况
				if i != nodesLength-1 {
					// 如果不是叶子节点

					var v interface{} = make(map[string]interface{})
					*curPointer = v
					parentPointer = curPointer
					curPointer = &v
				} else {
					// 如果是叶子节点

					*curPointer = map[string]interface{}{nodes[i]: content}
				}
			} else {
				// 当前节点的key是索引的情况

				curNodeIndex, _ := stringx.GetIndex(nodes[i])
				if i != nodesLength-1 {
					// 如果不是叶子节点

					if stringx.IsIndex(nodes[i+1]) {
						// 如果当前节点的下一个节点是索引
						// 那么就需要为下一个节点的遍历开辟好空间

						childCurNodeIndex, _ := stringx.GetIndex(nodes[i+1])
						if curValueIfaces, ok := (*curPointer).([]interface{}); ok {
							// 如果当前指针所指的值不是空接口切片类型

							var v interface{} = make([]interface{}, curNodeIndex+1)
							parentPointer = j.setPointer(curPointer, nodes[i], v)
							curPointer = &v
						} else {
							// 如果当前指针所指的值是空接口切片类型

							if len(curValueIfaces) < curNodeIndex {
								// 如果当前切片空间小于结点的索引,需要开辟新的空间

								if content == nil {
									// 传入的内容为空,但是当前节点没有对应的内容,所以直接返回空

									return nil
								}
								var v interface{} = make([]interface{}, curNodeIndex+1)
								parentPointer = j.setPointer(curPointer, nodes[i], v)
								curPointer = &v
							} else {

								chileValueIface := curValueIfaces[curNodeIndex]
								if childPointerIfaces, ok := chileValueIface.([]interface{}); ok {
									// 如果当前指针所指的空接口切片的对应的索引值是空接口切片

									for j := 0; j < childCurNodeIndex-len(childPointerIfaces); j++ {
										// 如果索引大于空接口切片的大小
										// 开辟新的空间
										childPointerIfaces = append(childPointerIfaces, nil)
									}
									parentPointer = curPointer
									curPointer = &curValueIfaces[curNodeIndex]
								} else {
									// 如果下一个节点是一个索引,但是下一个节点又不是空接口切片类型
									// 创建一个新的空接口切片
									// 这里相当于覆盖下一个节点所指的内容

									if content == nil {
										// 传入的内容为空,但是当前节点没有对应的内容,所以直接返回空

										return nil
									}

									var v interface{} = make([]interface{}, curNodeIndex+1)
									parentPointer = j.setPointer(curPointer, nodes[i], v)
									curPointer = &v
								}
							}
						}
					} else {
						// 如果当前节点的下一个节点不是索引

						if curValueIfaces, ok := (*curPointer).([]interface{}); !ok {
							// 如果当前节点所指的值不是空接口切片类型
							// 在当前节点创建一个空接口切片
							// 在下一个节点开一个map[string]interface{}

							s := make([]interface{}, curNodeIndex+1)
							s[curNodeIndex] = make(map[string]interface{})
							if parentPointer == nil {
								// i=0

								var v interface{} = s
								*curPointer = v
								parentPointer = curPointer
								curPointer = &s[curNodeIndex]
							} else {
								// i>0

								j.setPointer(parentPointer, nodes[i-1], s)
								parentPointer = curPointer
								curPointer = &s[curNodeIndex]

							}

						} else {
							// 如果当前节点所指的值是空接口切片类型

							if len(curValueIfaces) > curNodeIndex {
								// 若空接口切片的空间大小大于索引
								// 直接在空接口切片对应的索引值赋值

								parentPointer = curPointer
								curPointer = &(*curPointer).([]interface{})[curNodeIndex]
							} else {
								// 若空接口切片的空间大小小于索引
								// 先开辟新的内存空间,再赋值

								s := make([]interface{}, curNodeIndex+1)
								copy(s, curValueIfaces)
								s[curNodeIndex] = make(map[string]interface{})
								if parentPointer == nil {
									// i=0

									var v interface{} = s
									*curPointer = v
									parentPointer = curPointer
									curPointer = &s[curNodeIndex]
								} else {
									// i>0

									j.setPointer(parentPointer, nodes[i-1], s)
									parentPointer = curPointer
									curPointer = &s[curNodeIndex]

								}
							}
						}
					}
				} else {
					// 如果是叶子节点

					if curValueIfaces, ok := (*curPointer).([]interface{}); !ok {
						// 如果当前指针所指的值不是空接口切片
						// 开一个新的切片,然后直接覆盖

						s := make([]interface{}, curNodeIndex+1)
						s[curNodeIndex] = content
						j.setPointer(parentPointer, nodes[i-1], s)
					} else {
						// 如果当前指针所指的值是空接口切片

						if len(curValueIfaces) > curNodeIndex {
							// 空接口切片的空间大于当前节点的索引
							// 直接覆盖值
							if content == nil {
								// 如果是删除操作

								if parentPointer == nil {
									// 如果父节点不为空
									*curPointer = append(curValueIfaces[:curNodeIndex], curValueIfaces[curNodeIndex+1:]...)
								} else {
									j.setPointer(parentPointer, nodes[i-1], append(curValueIfaces[:curNodeIndex], curValueIfaces[curNodeIndex+1:]...))
								}
							} else {
								curValueIfaces[curNodeIndex] = content
							}

						} else {
							// 空接口切片的空间小于当前节点的索引
							// 开辟新的空间然后覆盖

							if content == nil {
								// 传入的内容为空,但是当前节点没有对应的内容,所以直接返回空

								return nil
							}
							if parentPointer == nil {
								// i=0
								// 直接覆盖当前节点值

								j.setPointer(curPointer, nodes[i], content)
							} else {
								// i>0
								// 先取父节点的值
								// 然后开辟新的空间再赋值

								s := make([]interface{}, curNodeIndex+1)
								copy(s, curValueIfaces)
								s[curNodeIndex] = content
								j.setPointer(parentPointer, nodes[i-1], s)
							}
						}
					}

				}
			}

		default:
			// 如果当前节点不是可索引的类型,直接替换当前值

			if content == nil {
				return errors.New(invalidPattern)
			}

			if stringx.IsIndex(nodes[i]) {
				// 如果当前节点的key是索引的情况
				//  创建一个新的切片

				curNodeIndex, _ := stringx.GetIndex(nodes[i])
				s := make([]interface{}, curNodeIndex+1)
				if i == nodesLength-1 {
					s[curNodeIndex] = content
				}
				if parentPointer == nil {
					// i=0

					*curPointer = s
					parentPointer = curPointer
				} else {
					// i>0

					parentPointer = j.setPointer(parentPointer, nodes[i-1], s)
				}
			} else {
				// 如果当前节点的key不是索引的情况
				// 创建一个新的map[string]interface{}

				var v1, v2 interface{}
				if i != nodesLength-1 {
					// 如果当前节点不是叶子节点

					v1 = map[string]interface{}{
						nodes[i]: nil,
					}
				} else {
					// 如果当前节点是叶子节点

					v1 = map[string]interface{}{
						nodes[i]: content,
					}
				}
				if parentPointer == nil {
					// i=0

					*curPointer = v1
					parentPointer = curPointer
				} else {
					// i>0

					parentPointer = j.setPointer(parentPointer, nodes[i-1], v1)
				}
				v2 = v1.(map[string]interface{})[nodes[i]]
				curPointer = &v2
			}
		}
	}
	return nil
}

func (j *Json) setPointer(pointer *interface{}, key string, value interface{}) *interface{} {
	switch (*pointer).(type) {
	case map[string]interface{}:
		(*pointer).(map[string]interface{})[key] = value
		return &value
	case []interface{}:
		index, _ := stringx.GetIndex(key)
		if len((*pointer).([]interface{})) > index {
			(*pointer).([]interface{})[index] = value
			return &(*pointer).([]interface{})[index]
		} else {
			s := make([]interface{}, index+1)
			copy(s, (*pointer).([]interface{}))
			s[index] = value
			*pointer = s
			return &s[index]
		}
	default:
		*pointer = value
	}
	return pointer
}

func (j *Json) convertContent(data interface{}, options Options) (convertedValue interface{}, err error) {
	if data == nil {
		return nil, errors.New(emptyContest)
	}
	switch data.(type) {
	case map[string]interface{}, map[string][]interface{}:
		// 传入的已经是解码好的json数据的情况
		return data, nil
	case string, []byte:
		// 传入的是字符串或者bytes的情况:
		// 判断数据的格式(json,yaml,toml...),转化成json格式
		// 然后将数据解码成map[string]interface{}的形式
		content := conv.ToBytes(data)
		if len(content) == 0 {
			return nil, errors.New(emptyContest)
		}
		parsedContent, err := j.parseContent(content, options)
		if err != nil {
			return nil, err
		}
		return parsedContent, nil
	default:
		var content interface{}
		switch reflect.ValueOf(data).Kind() {
		case reflect.Struct, reflect.Map:
			// 传入的是可递归结构的情况:

			// 如果结构体是接口的情况:
			// 取值然后再递归下去
			// 方法①:
			//   先将结构体marshal成[]bytes,然后unmarshal成map[string]interface{}
			//   但是效率慢,这里想自己写个递归方法
			// 方法②:
			//   直接将结构体转化成map[string]interface{}
			//   利用反射层层递归
			// 这里采用方法②
			content = conv.MapSearch(data, "json")
		case reflect.Slice, reflect.Array:
			// 返回空接口切片
			content = conv.ToInterfaces(data)
		default:
			content = data
			return content, nil
		}
		return content, nil
	}
}
