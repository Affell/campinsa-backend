package packHandler

const Service = "pack"

type (
	OpenPackForm struct {
		PackToken string `json:"pack_token" binding:"required"`
	}
)
