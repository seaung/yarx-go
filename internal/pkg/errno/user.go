package errno

var (
	UserAlreadyExistError = &Errno{HTTP: 400, Code: "FailedOperation.UserAlreadyExist", Message: "User already exist"}

	PasswordIncorrectError = &Errno{HTTP: 401, Code: "InvalidParameter.PasswordIncorrect", Message: "Password was incorrect"}

	UserNotFoundError = &Errno{HTTP: 404, Code: "ResourceNotFound.UserNotFound", Message: "User was not found"}
)
