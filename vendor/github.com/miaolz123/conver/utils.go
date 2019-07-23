package conver

import "fmt"

func converError(val interface{}, t string) error {
	return fmt.Errorf("conver error, the %T{%v} can not conver to a %v", val, val, t)
}

func stringToBool(val string) (bool, error) {
	switch val {
	case "1", "t", "T", "true", "TRUE", "True", "ok", "OK", "yes", "YES":
		return true, nil
	case "0", "f", "F", "false", "FALSE", "False", "":
		return false, nil
	}
	return false, converError(val, "bool")
}
