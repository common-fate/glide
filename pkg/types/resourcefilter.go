package types

import (
	"errors"
	"strings"
)

func (o *Operation) Match(r *Resource) (bool, error) {
	switch o.OperationType {
	case SELECT:
		{
			if o.Values == nil {
				return false, errors.New("for select operation, values field cannot be empty")
			}

			for _, v := range *o.Values {
				if v == r.Id {
					return true, nil
				}
			}
		}
	case BEGINSWITH:
		{
			if r.Attributes == nil {
				return false, errors.New("this resource doesnot contain any attributes")
			}

			if v, ok := r.Attributes[o.Attribute]; ok {
				return strings.HasPrefix(v, *o.Value), nil
			}
		}
	}

	return false, nil
}

func (r *Resource) Match(filter ResourceFilter) bool {
	for _, operation := range filter {

		matched, err := operation.Match(r)
		if err != nil {
			return false
		}

		return matched
	}
	return true
}
