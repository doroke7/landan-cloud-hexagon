package pkgWebsocket

type Hub struct {
}

/*
對 sCId 的 連線 都發消息
*/
func (oSelf *Hub) To(sCId string, aByteMessage []byte) {

}

/*
對 sChannel 的 連線 都發消息
*/
func (oSelf *Hub) Emit(sChannel string, aByteMessage []byte) {

}
