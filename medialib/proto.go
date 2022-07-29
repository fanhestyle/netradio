package medialib

const (
	DEFAULT_PORT      = "1989"
	DEFAULT_MGROUP    = "224.2.2.2"
	DEFAULT_PLAYERCMD = "mpg123"
	DEFAULT_MEDIADIR  = "/home/hzx/tmp/media"
	DEFAULT_IF        = "bridge0"
	RUN_FOREGROUND    = false
)

type Server_conf_st struct {
	Rcvport   string
	Mgroup    string
	Media_dir string
	Runmode   bool
	Ifname    string
}

var Server_conf = &Server_conf_st{
	Rcvport:   DEFAULT_PORT,
	Mgroup:    DEFAULT_MGROUP,
	Media_dir: DEFAULT_MEDIADIR,
	Ifname:    DEFAULT_IF,
	Runmode:   RUN_FOREGROUND,
}
