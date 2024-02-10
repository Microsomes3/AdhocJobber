package typeos

import "database/sql/driver"

type NullType uint

func NewNullType() NullType {
	return NullType(0)
}

func (nt NullType) Value() (driver.ValueConverter, error) {
	return nil, nil
}
