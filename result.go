package ftgo

var (
	ResultOk         = map[string]int{"code": 1}
	ResultParamError = ResultError("参数有问题")
)

type resultCodeError struct {
	Code  int    `json:"code"`
	Error string `json:"error,omitempty"`
}

var resultCode = 9

func ResultError(e string) resultCodeError {
	resultCode++
	return resultCodeError{
		Code:  resultCode,
		Error: e,
	}
}

func ResultZero(e string) resultCodeError {
	return resultCodeError{
		Code:  0,
		Error: e,
	}
}

func ResultMap(m map[string]interface{}) map[string]interface{} {
	m["code"] = 1
	return m
}

func ResultData(d interface{}) map[string]interface{} {
	return map[string]interface{}{
		"code": 1,
		"data": d,
	}
}
