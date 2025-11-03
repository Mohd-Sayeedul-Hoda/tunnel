package natserver

import (
	"context"
	"io"
	"log/slog"
	"net"
	"strconv"

	"github.com/Mohd-Sayeedul-Hoda/tunnel/internal/shared/config"

	"github.com/hashicorp/yamux"
)

func ListenAndServer(ctx context.Context, w io.Writer, cfg *config.Config) error {

	listner, err := net.Listen("tcp", cfg.NatTcpServer.Host+":"+strconv.Itoa(cfg.NatTcpServer.Port))
	if err != nil {
		return err
	}
	defer listner.Close()

	slog.Info("tcp server started", slog.String("addr", cfg.NatTcpServer.Host+":"+strconv.Itoa(cfg.NatTcpServer.Port)))
	for {
		conn, err := listner.Accept()
		if err != nil {
			slog.Error("failed to accept connection", slog.String("error", err.Error()))
			continue
		}

		go func() {
			defer conn.Close()
			ManageConnection(conn, w, cfg)
		}()
	}
}

func ManageConnection(conn net.Conn, w io.Writer, cfg *config.Config) {

	yamuxConfig := yamux.DefaultConfig()
	yamuxConfig.LogOutput = w
	session, err := yamux.Server(conn, yamuxConfig)
	if err != nil {
		// happen when yamux config verfication failed
		slog.Error("unable to open the yamux session closing tcp connection", slog.Any("yamux error", err.Error()))
		conn.Close()
		return
	}
}
