package zid

import (
	"github.com/rosas99/streaming/pkg/id"
)

const defaultABC = "abcdefghijklmnopqrstuvwxyz1234567890"

type ZID string

const (
	// ID for the user resource in streaming-usercenter.
	User ZID = "user"
	// ID for the order resource in streaming-fakeserver.
	Order ZID = "order"
	// ID for the template resource in streaming-sms.
	Template ZID = "template"
)

func (zid ZID) String() string {
	return string(zid)
}

func (zid ZID) New(i uint64) string {
	// use custom option
	str := id.NewCode(
		i,
		id.WithCodeChars([]rune(defaultABC)),
		id.WithCodeL(6),
		id.WithCodeSalt(Salt()),
	)
	return zid.String() + "-" + str
}
