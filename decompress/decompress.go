package decompress

import (
	"encoding/json"
	"errors"
	"github.com/V-I-C-T-O-R/DataCompress/compress"
	"github.com/V-I-C-T-O-R/DataCompress/utils"
	"log"
	"reflect"
	"strconv"
	"strings"
)

var matchMap = map[int]interface{}{
	1: reflect.Bool,
	2: reflect.Int,
	3: reflect.Float64,
	4: reflect.String,
}

func DoDeCompress(file string) (data []byte, err error) {
	if !utils.CheckFileIsExist(file) {
		log.Println("file no exist")
		err = errors.New("file no exist")
		return
	}
	s := utils.ReadF(file)
	var content interface{}
	err = json.Unmarshal(s, &content)
	if err != nil {
		log.Println("json file unmarshal failed")
		return
	}
	data, _ = json.Marshal(baseMap(&compress.MesMark{Mark: true, Value: (content).(map[string]interface{})}))
	log.Println("data decompress complete")
	return
}
func DoDeCompressFromData(b []byte) (data []byte, err error) {
	var content interface{}
	err = json.Unmarshal(b, &content)
	if err != nil {
		log.Println("json file unmarshal failed")
		return
	}
	data, _ = json.Marshal(baseMap(&compress.MesMark{Mark: true, Value: (content).(map[string]interface{})}))
	log.Println("data decompress complete")
	return
}
func baseMap(m *compress.MesMark) (x interface{}) {
	if m.Mark {
		x = make(map[string]interface{})
		for k, v := range m.Value.(map[string]interface{}) {
			switch v.(type) {
			case bool, byte, int, int8, int16, int32, int64, uint16, uint32, uint64, float32, float64, string:
				x.(map[string]interface{})[k] = v
				continue
			case []interface{}:
				x.(map[string]interface{})[k] = baseMap(&compress.MesMark{Mark: false, Value: v, Key: k})
			default:
				m.Value.(map[string]interface{})[k] = baseMap(&compress.MesMark{Mark: true, Value: v, Key: k})
			}
		}
	} else {
		var listCount []interface{}
		var listMap []interface{}
		var flag bool
		for _, v := range m.Value.([]interface{}) {
			switch v.(type) {
			case bool, byte, int, int8, int16, int32, int64, uint16, uint32, uint64, float32, float64, string:
				listCount = append(listCount, v)
				flag = true
				continue
			case []interface{}:
				baseMap(&compress.MesMark{Mark: false, Value: v})
			default:
				listMap = append(listMap, baseMap(&compress.MesMark{Mark: true, Value: v}))
			}
		}
		if flag {
			var slice []interface{}
			for _, v := range listCount {
				p := parse(v)
				for _, j := range p.([]interface{}) {
					slice = append(slice, j)
				}
			}
			x = slice
		} else {
			x = listMap
		}
	}
	return
}
func toValue(v interface{}, value string) (s interface{}) {
	switch v {
	case reflect.Bool:
		s, _ = strconv.ParseBool(value)
	case reflect.Int:
		s, _ = strconv.Atoi(value)
	case reflect.Float64:
		s, _ = strconv.ParseFloat(value, 64)
	case reflect.String:
		s = value
	}
	return
}
func parse(s interface{}) interface{} {
	str := strings.Split(s.(string), "::")
	k, _ := matchMap[invertToInt(str[1])]
	value := toValue(k, str[0])
	var slice []interface{}
	for i := 0; i < invertToInt(str[2]); i++ {
		slice = append(slice, value)
	}
	return slice
}
func invertToInt(value string) (s int) {
	s, _ = strconv.Atoi(value)
	return
}
