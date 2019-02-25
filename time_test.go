package elgo

import (
	"testing"
	"time"

	"github.com/magiconair/properties/assert"
)

func TestStringToTime(t *testing.T) {
	dateStr := "2018-11-06T22:48:50.00275Z"
	tt, _ := StrToTime(dateStr)
	assert.Equal(t, tt.Year(), 2018)
	assert.Equal(t, tt.Month(), time.November)
	assert.Equal(t, tt.Day(), 6)
	assert.Equal(t, tt.Hour(), 22)
}

func TestTimeToString(t *testing.T) {
	expectStr := "2018-11-06T22:48:50.00275Z"
	tt, _ := time.Parse(time.RFC3339Nano, expectStr)
	assert.Equal(t, TimeToStr(tt), expectStr)
}
