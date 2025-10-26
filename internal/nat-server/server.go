package natserver

import (
	"io"
	"log/slog"

	"github.com/Mohd-Sayeedul-Hoda/tunnel/internal/shared/config"

	"github.com/hashicorp/yamux"
)

func HandleConnection(cfg *config.Config, conn io.ReadWriteCloser) error {

	session, err := yamux.Server(conn, yamux.DefaultConfig())
	if err != nil {
		return err
	}

	stream, err := session.Accept()
	if err != nil {
		return err
	}
	buf := make([]byte, 4)
	stream.Read(buf)
	slog.Info(string(buf))

	return nil
}
