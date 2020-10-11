package parsup

import (
	"encoding/json"
	"errors"
	"reflect"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ParamsSupport 参数辅助类
type ParamsSupport struct {
	IsDeep       *bool // 深度递归
	IsConvOID    *bool // 转化ObjectID
	IsConvTime   *bool // 转化时间对象
	IsDenyInject *bool // 防注入
}

// ParSup 工厂方法
func ParSup(isDeep, isConvOID, isConvTime, isDenyInject bool) *ParamsSupport {
	return &ParamsSupport{
		IsDeep:       &isDeep,
		IsConvOID:    &isConvOID,
		IsConvTime:   &isConvTime,
		IsDenyInject: &isDenyInject,
	}
}

// ConvBase base handler
func (p *ParamsSupport) ConvBase(i interface{}) (interface{}, error) {
	v := reflect.ValueOf(i)
	switch v.Kind() {
	case reflect.Map:
		if *p.IsDeep {
			return p.ConvMap(i.(map[string]interface{}))
		}
		break
	case reflect.Slice:
		if *p.IsDeep {
			return p.ConvSlice(i.([]interface{}))
		}
		break
	case reflect.Invalid:
		return nil, nil
	case reflect.String:
		return p.ConvStr(i.(string))
	}
	return v, nil

}

// ConvStr string handler
func (p *ParamsSupport) ConvStr(s string) (interface{}, error) {
	if *p.IsDenyInject {
		if strings.ContainsAny(s, "$[]{}()") {
			return nil, errors.New("不能包含$[]{}()等特殊符号")
		}
	}

	if *p.IsConvOID {
		if oid, err := primitive.ObjectIDFromHex(s); err == nil {
			return oid, nil
		}
	}
	if *p.IsConvTime {
		if t, err := time.Parse(time.RFC3339, s); err == nil {
			return t.Local(), nil
		}
	}
	return s, nil
}

// ConvMap map handler
func (p *ParamsSupport) ConvMap(m map[string]interface{}) (map[string]interface{}, error) {
	var res map[string]interface{}
	for k, v := range m {
		dv, err := p.ConvBase(v)
		if err != nil {
			return nil, err
		}
		res[k] = dv
	}
	return res, nil
}

// ConvSlice slice handler
func (p *ParamsSupport) ConvSlice(s []interface{}) ([]interface{}, error) {
	res := make([]interface{}, len(s))
	for k, v := range s {
		dv, err := p.ConvBase(v)
		if err != nil {
			return nil, err
		}
		res[k] = dv
	}
	return res, nil
}

// ConvJSON byte handler
func (p *ParamsSupport) ConvJSON(s []byte) (map[string]interface{}, error) {
	var m map[string]interface{}
	err := json.Unmarshal(s, &m)
	if err != nil {
		return nil, err
	}
	return p.ConvMap(m)
}
