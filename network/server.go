package network

type ServerOpts struct {
	Transports []Transport
}

type Server struct {
	ServerOpts
}

func NewServer(opts ServerOpts) *Server {
	return &Server{
		ServerOpts: opts,
	}
}
