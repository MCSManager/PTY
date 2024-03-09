package start

import (
	"encoding/json"
	"fmt"
	"io"

	pty "github.com/MCSManager/pty/console"
	"github.com/zijiren233/stream"
)

type connUtils struct {
	r *stream.Reader
	w *stream.Writer
}

func newConnUtils(r io.Reader, w io.Writer) *connUtils {
	return &connUtils{
		r: stream.NewReader(r, stream.BigEndian),
		w: stream.NewWriter(w, stream.BigEndian),
	}
}

func (cu *connUtils) ReadMessage() (uint8, []byte, error) {
	var (
		length  uint16
		msgType uint8
	)
	data, err := cu.r.U8(&msgType).U16(&length).ReadBytes(int(length))
	return msgType, data, err
}

func (cu *connUtils) SendMessage(msgType uint8, data any) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return cu.w.U8(msgType).U16(uint16(len(b))).Bytes(b).Error()
}

func handleConn(u *connUtils, con pty.Console) error {
	for {
		t, msg, err := u.ReadMessage()
		if err != nil {
			return fmt.Errorf("read message error: %w", err)
		}
		switch t {
		case RESIZE:
			resize := resizeMsg{}
			err := json.Unmarshal(msg, &resize)
			if err != nil {
				_ = u.SendMessage(
					ERROR,
					&errorMsg{
						Msg: fmt.Sprintf("unmarshal resize message error: %s", err),
					},
				)
				continue
			}
			err = con.SetSize(resize.Width, resize.Height)
			if err != nil {
				_ = u.SendMessage(
					ERROR,
					&errorMsg{
						Msg: fmt.Sprintf("resize error: %s", err),
					},
				)
				continue
			}
		}
	}
}

func testResize(u *connUtils) error {
	err := u.SendMessage(
		RESIZE,
		&resizeMsg{
			Width:  20,
			Height: 20,
		},
	)
	if err != nil {
		return fmt.Errorf("send resize message error: %w", err)
	}
	return nil
}
