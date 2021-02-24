package server

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gopkg.in/olahol/melody.v1"
	"github.com/googollee/go-socket.io"
	"log"
	"time"
)

type Server struct {
	HttpAddr string
	Engine *gin.Engine
	Melody *melody.Melody
	SocketIO *socketio.Server
}

func New(host string, port uint) Server {
	gin.SetMode(gin.ReleaseMode)

	socketioServer, _ := socketio.NewServer(nil)

	srv := Server{
		Engine:  gin.New(),
		HttpAddr: fmt.Sprintf("%s:%d", host, port),
		Melody: melody.New(),
		SocketIO: socketioServer,
	}

	srv.Engine.Use(cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
		AllowHeaders:     []string{"Authorization", "Content-Length", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	return srv
}


func (srv *Server) SetupRoutes() {
	srv.Engine.GET("/socket.io/*any", gin.WrapH(srv.SocketIO))
	srv.Engine.POST("/socket.io/*any", gin.WrapH(srv.SocketIO))

	//srv.Engine.GET("/socket.io/*path",  func(ctx *gin.Context) {
	//	log.Println(ctx.Param("path"))
	//	err := srv.Melody.HandleRequest(ctx.Writer, ctx.Request)
	//	log.Println("error", err.Error())
	//})

	srv.SocketIO.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		fmt.Println("connected:", s.ID())
		return nil
	})

	srv.SocketIO.OnEvent("/", "default", func(s socketio.Conn, msg string) {
		fmt.Println("notice:", msg)
		s.Emit("reply", "have "+msg)
	})

	srv.SocketIO.OnEvent("/chat", "msg", func(s socketio.Conn, msg string) string {
		s.SetContext(msg)
		return "recv " + msg
	})

	srv.SocketIO.OnEvent("/", "bye", func(s socketio.Conn) string {
		last := s.Context().(string)
		s.Emit("bye", last)
		s.Close()
		return last
	})

	srv.SocketIO.OnError("/", func(s socketio.Conn, e error) {
		fmt.Println("meet error:", e)
	})

	srv.SocketIO.OnDisconnect("/", func(s socketio.Conn, reason string) {
		fmt.Println("closed", reason)
	})

}

func (srv *Server) Run() error {
	log.Println("Server running on", srv.HttpAddr)

	srv.SetupRoutes()

	//srv.Melody.HandleMessage(func(s *melody.Session, msg []byte) {
	//	log.Println("WS Broadcast", msg)
	//	srv.Melody.Broadcast(msg)
	//})

	go srv.SocketIO.Serve()
	defer srv.SocketIO.Close()
	return srv.Engine.Run(srv.HttpAddr)
}