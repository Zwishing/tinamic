package validate

import (
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	"reflect"
	"sync"
)

var (
	validate *validator.Validate
	once     sync.Once
)

func GetValidatorInstance() *validator.Validate {
	once.Do(func() {
		validate = validator.New()
	})
	return validate
}

// ValidateRequestBody
func ValidateRequestBody(s any) error {
	v := reflect.ValueOf(s)

	// 检查参数是否是指针，如果是则获取指针指向的值
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	// 检查参数是否是结构体
	if v.Kind() != reflect.Struct {
		return errors.New("ProcessStruct: argument must be a struct or a struct pointer")
	}

	vdt := GetValidatorInstance()
	err := vdt.Struct(s)
	if err != nil {
		return err
	}
	return nil
}
