package server

type ServerOptions struct {
	Port string
}

func NewDefaultServerOptions() *ServerOptions {
	return &ServerOptions{
		Port: "8080",
	}
}
