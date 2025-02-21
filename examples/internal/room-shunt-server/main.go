package main

import "github.com/kercylan98/minotaur/server"

func main() {
	srv := server.New(server.NetworkWebsocket,
		server.WithShunt(func(guid int64) chan *server.Message {
			return make(chan *server.Message, 1024*100)
		}, func(conn *server.Conn) (guid int64, allowToCreate bool) {
			guid, allowToCreate = conn.GetData("roomId").(int64)
			return
		}),
	)
	srv.RegConnectionReceivePacketEvent(func(srv *server.Server, conn *server.Conn, packet []byte) {
		conn.Write(packet)
	})
	if err := srv.Run(":9999"); err != nil {
		panic(err)
	}
}
