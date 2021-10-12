package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"bitbucket.org/latonaio/calendar-module-kube/cmd/calendar-server/app"
	"bitbucket.org/latonaio/calendar-module-kube/cmd/calendar-server/proto/calendarpb"
	"google.golang.org/grpc"
)

func main() {
	serverEnv, err := app.GetServerEnv()
	if err != nil {
		log.Fatalln(err)
	}

	listen, err := net.Listen("tcp", serverEnv.Addr+":"+serverEnv.Port)
	if err != nil {
		log.Fatalln(err)
	}
	grpcServer := grpc.NewServer()

	dbEnv, err := app.GetDatabaseEnv()
	if err != nil {
		log.Fatalln(err)
	}
	db, err := app.NewDB(dbEnv)
	if err != nil {
		log.Fatalln(err)
	}
	calendarService := &app.CalendarService{Db: db}

	calendarpb.RegisterCalendarServer(grpcServer, calendarService)

	log.Printf("[info] Start Calendar Server")
	grpcServer.Serve(listen)

	signalCh := make(chan os.Signal)
	signal.Notify(signalCh, syscall.SIGTERM)

loop:
	for {
		select {
		case s := <-signalCh:
			log.Printf("[info] recieved signal: %s", s.String())
			grpcServer.GracefulStop()
			break loop
		}
	}
}
