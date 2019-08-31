package util

import (
	"fmt"
	"github.com/mhconradt/blog-api/config"
	"os"
	"strconv"
	"strings"
)

func ToInt64(i interface{}) int64 {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("panic'd in ToInt64 with val: ", i)
		}
	}()
	valStr := i.(string)
	val, _ := strconv.ParseInt(valStr, 10, 64)
	return val
}

func ToInt(i interface{}) int {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("panic'd in ToInt with val: ", i)
		}
	}()
	valStr := i.(string)
	val, _ := strconv.Atoi(valStr)
	return val
}

func ToStringSlice(is []interface{}) []string {
	numVals := len(is)
	vs := make([]string, numVals, numVals)
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

func ZipMap(keys []string, vals []interface{}) map[string]interface{} {
	m := make(map[string]interface{})
	for i, k := range keys {
		if vals[i] != nil {
			m[k] = vals[i]
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
