package oauth_flag

import "flag"

var(
	AuthPath string
	LoginPath string
)

func init() {
	flag.StringVar(&AuthPath, "auth", "./server/views/auth.html", "auth path")
	flag.StringVar(&LoginPath, "login", "./server/views/login.html", "login path")
	flag.Parse()
}
