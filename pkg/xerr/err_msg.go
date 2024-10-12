package xerr

var codeText = map[int]string{
	SERVER_COMMON_ERROR: "Server exception, try again later",
	REQUEST_PARAM_ERROR: "The request parameters are incorrect",
	DB_ERROR:            "Database is busy, try again later",
}

func ErrMsg(errcode int) string {
	if msg, ok := codeText[errcode]; ok {
		return msg
	}
	return codeText[SERVER_COMMON_ERROR]
}