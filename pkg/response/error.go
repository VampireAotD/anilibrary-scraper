package response

type Error struct {
	Message string `json:"message"`
}

func (r Response) ErrorJSON(code int, err error) error {
	return r.JSON(code, Error{
		Message: err.Error(),
	})
}
