package client

import (
	"google.golang.org/grpc"

	pbResourceAnnouncement "example/pb/source/announcement"
)

func NewAnnouncement(oClientConn *grpc.ClientConn) *Announcement {

	oAnnouncement := &Announcement{
		Lottery: pbResourceAnnouncement.NewLotteryControllerClient(oClientConn),
	}

	return oAnnouncement
}

type Announcement struct {
	Lottery pbResourceAnnouncement.LotteryControllerClient
}

func NewSourceClient(oClientConn *grpc.ClientConn, oAnnouncement *Announcement) *SourceClient {

	return &SourceClient{
		conn:         oClientConn,
		Announcement: oAnnouncement,
	}
}

type SourceClient struct {
	conn         *grpc.ClientConn
	Announcement *Announcement
}

func (oClient *SourceClient) Close() error {
	oErr := oClient.conn.Close()

	return oErr
}
