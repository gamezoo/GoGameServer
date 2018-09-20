package message

import (
	. "core/libs"
	"core/libs/grpc/ipc"
	"core/libs/sessions"
	"encoding/binary"
)

func IpcClientReceive(stream ipc.Ipc_TransferClient, msg *ipc.Res) {
	frontSession := sessions.GetFrontSession(msg.SessionId)
	msgBody := msg.Data
	if frontSession != nil {
		frontSession.Send(msgBody)
	} else {
		msgId := binary.BigEndian.Uint16(msgBody[:2])
		WARN("frontSession no exists", msgId)
	}
}