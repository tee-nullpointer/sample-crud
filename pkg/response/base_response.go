package response

const (
	CodeSuccess             = "00"
	CodeInternalServerError = "50"
)

type BaseResponse struct {
	ResponseCode    string      `json:"response_code"`
	ResponseMessage string      `json:"response_message"`
	Data            interface{} `json:"data,omitempty"`
}

type CreatedData struct {
	ID uint `json:"id"`
}

func Error(responseCode string, responseMessage string) BaseResponse {
	return BaseResponse{
		ResponseCode:    responseCode,
		ResponseMessage: responseMessage,
	}
}

func Success(data interface{}) BaseResponse {
	return BaseResponse{
		ResponseCode:    CodeSuccess,
		ResponseMessage: "Success",
		Data:            data,
	}
}

func NewCreatedData(id uint) CreatedData {
	return CreatedData{
		ID: id,
	}
}

func Created(id uint) BaseResponse {
	return Success(NewCreatedData(id))
}

func InternalServerError() BaseResponse {
	return Error(CodeInternalServerError, "INTERNAL_SERVER_ERROR")
}
