package clip
import (
    "github.com/gin-gonic/gin"
    "strconv"
)

type SbHTTPResponse struct {
    Code   int          `json:"-"`
    Errors *ErrorStruct `json:"errors"`
    Data   interface{}  `json:"data"`
    Meta   interface{}  `json:"meta"`
}

type FiserResponse struct {
    Errors *ErrorStruct `json:"errors"`
    Data   string       `json:"data"`
    Meta   interface{}  `json:"meta"`
}


// ErrorStruct hold entity for any error response returned to user
type ErrorStruct struct {
    Status  string `json:"status,omitempty"`
    Code    int    `json:"code,omitempty"`
    Message string `json:"detail,omitempty"`
}

func NewResponse() SbHTTPResponse {
    resp := SbHTTPResponse{}

    errResp := new(ErrorStruct)
    errResp.Code = 200
    errResp.Status = "200"
    errResp.Message = ""

    resp.Errors = errResp

    return resp
}

// ErrorResponse will create error response
func (k *SbHTTPResponse) ErrorResponse(code int, message string) {
    k.Data = nil

    err := new(ErrorStruct)
    err.Code = code
    err.Message = message
    err.Status = strconv.Itoa(code)

    k.Errors = err

    // set parent error code
    k.Code = code
}


func (k *SbHTTPResponse) SetData(data interface{}) {
    k.Data = data
}

func (k *SbHTTPResponse) GetResponse() *SbHTTPResponse {
    if k.Errors.Code == 200 {
        k.Errors = nil

    } else if k.Errors.Code != 200 {
        k.Data = nil
    }

    return k
}

func (k *SbHTTPResponse) GetErrorResponse() *ErrorStruct {
    return k.Errors
}

func (k *SbHTTPResponse) GetCode() int {
    return k.Code
}


func (k *SbHTTPResponse) JSON(c *gin.Context) {
    c.Header("Content-Type", "application/vnd.api+json; charset=UTF-8")
    c.JSON(k.GetCode(), k.GetResponse())
    return
}

