package errno

var (
	HTTPOK = &Errno{HTTP: 200, Code: "OK", Message: "Success"}

	InternalServerError = &Errno{HTTP: 500, Code: "InternalError", Message: "Internal server error"}

	PageNotFoundError = &Errno{HTTP: 404, Code: "ResourceNotFound.PageNotFound", Message: "Page Not Found"}

	BindParameterError = &Errno{HTTP: 400, Code: "InvalidParameter.BindError", Message: "Error occurred while binding request body to the struct"}

	InvalidParameterError = &Errno{HTTP: 400, Code: "InvalidParameter", Message: "parameter verifycation failed"}

	SignTokenError = &Errno{HTTP: 401, Code: "AuthFailure.SignTokenError", Message: "Error occurred while signing the JSON web token"}

	TokenInvalidError = &Errno{HTTP: 401, Code: "AuthFailure.TokenInvalid", Message: "Token was invalid"}

	UnAuthoizedError = &Errno{HTTP: 401, Code: "AuthFailure.UnAuthoized", Message: "UnAuthoized"}
)
