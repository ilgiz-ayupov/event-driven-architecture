package usecase

type UserLog struct {
	UserID string
}

func NewUserLog(userID string) UserLog {
	return UserLog{
		UserID: userID,
	}
}
