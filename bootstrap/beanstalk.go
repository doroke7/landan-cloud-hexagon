package bootstrap

import (
	"fmt"
	"time"

	beanstalk "github.com/beanstalkd/go-beanstalk"
)

func NewBeanstalk() (*beanstalk.Conn, error) {
	sAddr := fmt.Sprintf("%s:%s", CONFIG.BEANSTALK.HOST, CONFIG.BEANSTALK.PORT)

	return beanstalk.DialTimeout("tcp", sAddr, time.Duration(CONFIG.BEANSTALK.TIMEOUT)*time.Millisecond)
}
