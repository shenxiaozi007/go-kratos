package data

import (
	"context"
	"time"
)

type CustomerData struct {
	data *Data
}

func NewCustomerData(data *Data) *CustomerData {
	return &CustomerData{
		data: data,
	}
}

func (cd CustomerData) SetVerifyCode(telephone string, code string, expire int64) error {
	//设置值
	status := cd.data.Rdb.Set(context.Background(), "cvc:"+telephone, code, time.Second*time.Duration(expire))

	if _, err := status.Result(); err != nil {
		return err
	}
	return nil
}
