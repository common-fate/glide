package types

import (
	"errors"
	"fmt"
	"strings"
)

func (o *Operation) Match(r *Resource) (bool, error) {
	switch o.OperationType {
	case IN:
		{
			if o.Values == nil {
				return false, errors.New("for IN operation, values field cannot be empty")
			}

			for _, v := range *o.Values {
				if v == r.Id {
					return true, nil
				}
			}
		}
	case BEGINSWITH:
		{
			if v, ok := r.Attributes[o.Attribute]; ok {
				return strings.HasPrefix(v, *o.Value), nil
			}

			return false, fmt.Errorf("attribute %s not found", o.Attribute)
		}
	}

	return false, nil
}

func (r *Resource) Match(filter ResourceFilter) (bool, error) {
	for _, operation := range filter {
		matched, err := operation.Match(r)
		return matched, err
	}

	return true, nil
}
