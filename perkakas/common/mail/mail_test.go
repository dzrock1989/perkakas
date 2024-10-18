package mail

import (
	"testing"

	"github.com/tigapilarmandiri/perkakas/common/util"
)

func TestMail(t *testing.T) {
	m := NewMail()

	m.SetSubject("Omama")

	m.SetRecipient("andi@gmail.com")
	m.SetRecipient("indah@gmail.com")
	m.SetRecipient("lala@gmail.com")

	m.SetCC("tralala@gmail.com")
	m.SetCC("trelele@gmail.com")

	m.SetBCC("tralala@gmail.com")
	m.SetBCC("trelele@gmail.com")

	m.SetBody("<h1>Hello</>")
	m.Print()

	m = NewMail()
	err := m.Send()
	if err != nil {
		util.Log.Error().Msg(err.Error())
	}
}
