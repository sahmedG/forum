package api

import (
	"time"
)

// this map stores the users sessions. For larger scale applications, you can use a database or cache for this purpose
var (
	Sessions = make(map[string]Session, 0)
)

// each session contains the username of the user and the time at which it expires
type Session struct {
	userID      int
	SessionUUID string
	expiry      time.Time
}

// we'll use this method later to determine if the session has expired
func (s Session) IsExpired() bool {
	return s.expiry.Before(time.Now())
}

func (s Session) Get_UserID() int {
	return s.userID
}
