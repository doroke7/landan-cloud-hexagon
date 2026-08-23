package bootstrap

import (
	"fmt"
	"time"

	beanstalk "github.com/beanstalkd/go-beanstalk"
	"github.com/charmbracelet/log"
)

func NewBeanstalk() (*beanstalk.Conn, error) {
	sAddr := fmt.Sprintf("%s:%s", CONFIG.BEANSTALK.HOST, CONFIG.BEANSTALK.PORT)

	oConn, oErr := beanstalk.DialTimeout("tcp", sAddr, time.Duration(CONFIG.BEANSTALK.TIMEOUT)*time.Millisecond)

	log.Info("[INFO] BEANSTALK 連線完成. ", sAddr)

	return oConn, oErr
}
