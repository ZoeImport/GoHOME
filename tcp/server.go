package tcp

import (
	"context"
	"go-redis/interface/tcp"
	"go-redis/lib/logger"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Config struct {
	Address string
}

func ListenAndServerWithSignal(cfg *Config, handler tcp.Handler) error {
	listener, err := net.Listen("tcp", cfg.Address)
	if err != nil {
		return err
	}

	closeChan := make(chan struct{})
	signalsChan := make(chan os.Signal)
	signal.Notify(signalsChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)
	go func() {
		sig := <-signalsChan
		switch sig {
		case syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT:
			closeChan <- struct{}{}
		}
	}()
	logger.Info("start listen on " + cfg.Address)
	ListenAndServer(listener, handler, closeChan)
	return nil
}

func ListenAndServer(listener net.Listener, handler tcp.Handler, closeChan <-chan struct{}) {
	go func() {
		<-closeChan
		logger.Info("stop listening on " + listener.Addr().String())
		_ = listener.Close()
		_ = handler.Close()
	}()
	defer func() {
		_ = listener.Close()
		_ = handler.Close()
	}()

	ctx := context.Background()
	var waitDone sync.WaitGroup
	for {
		conn, err := listener.Accept()
		waitDone.Add(1)
		if err != nil {
			break
		}
		logger.Info("accept " + conn.RemoteAddr().String())
		go func() {
			defer func() {
				waitDone.Done()
			}()
			handler.Handle(ctx, conn)
		}()
	}
	waitDone.Wait()
}
