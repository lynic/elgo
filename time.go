package elgo

import "time"

//TimeToStr xxx
func TimeToStr(tt time.Time) string {
	return tt.Format(time.RFC3339Nano)
}

//StrToTime xxx
func StrToTime(ss string) (time.Time, error) {
	tt, err := time.Parse(
		time.RFC3339Nano,
		ss,
	)
	if err != nil {
		return tt, err
	}
	return tt, nil
}
