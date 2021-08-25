package controllers

//基类
type JSONStruct struct {
	Status string      `json:"status"`
	Code   int         `json:"code"`
	Result interface{} `json:"result"`
	Msg    string      `json:"msg"`
}

type Page struct {
	Total       int         `json:"total"`
	PerPage     int         `json:"per_page"`
	CurrentPage int         `json:"current_page"`
	PageNum     int         `json:"page_num"`
	Data        interface{} `json:"data"`
}

func SuccessRep(result interface{}) JSONStruct {
	return JSONStruct{
		Status: "success",
		Code:   0,
		Result: result,
		Msg:    "ok",
	}
}

func ErrorRep(code int, msg string) JSONStruct {
	return JSONStruct{
		Status: "error",
		Code:   code,
		Result: "",
		Msg:    msg,
	}
}
