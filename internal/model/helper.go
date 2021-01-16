package model

import "strconv"

type HelperInterface interface {
	GetByID(id int) (string, bool)
	GetAll() map[string]interface{}
}

type Helper struct {
	data map[int]string
}

func NewHelper() *Helper {
	return &Helper{}
}

func (h *Helper) SetData(data map[int]string) {
	h.data = data
}

func (h *Helper) GetAll() map[string]interface{} {
	m := make(map[string]interface{})
	for k, v := range h.data {
		m[strconv.Itoa(k)] = v
	}
	return m
}

func (h *Helper) GetByID(id int) (string, bool) {
	val, ok := h.data[id]
	if !ok {
		return "", false
	}
	return val, true
}
