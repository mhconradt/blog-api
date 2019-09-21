package util

import (
	"fmt"
	"github.com/mhconradt/blog-api/config"
	"os"
	"strconv"
	"strings"
)

func ToInt(i interface{}) int32 {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("panic'd in ToInt with val: ", i)
		}
	}()
	valStr := i.(string)
	val, _ := strconv.Atoi(valStr)
	return int32(val)
}

func ToStringSlice(is []interface{}) []string {
	count := len(is)
	vs := make([]string, count)
	for i, v := range is {
		vs[i] = v.(string)
	}
	return vs
}

func ToArray(i interface{}) []string {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("panic'd in ToArray with val: ", i)
		}
	}()
	str := i.(string)
	arr := strings.Split(str, config.TopicSeparator)
	return arr
}

func ZipMap(ks []string, vs []interface{}) map[string]interface{} {
	m := make(map[string]interface{})
	for i, k := range ks {
		if vs[i] != nil {
			m[k] = vs[i]
		}
	}
	return m
}

func LookupWithDefault(k, dv string) string {
	v, f := os.LookupEnv(k)
	if f {
		return v
	}
	return dv
}
