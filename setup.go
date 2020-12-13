package whoami_cname

import (
	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
)

func init() { plugin.Register("whoami_cname", setup) }

func setup(c *caddy.Controller) error {
	c.Next() // 'whoami_cname'
	if c.NextArg() {
		return plugin.Error("whoami_cname", c.ArgErr())
	}

	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		return Whoami_Cname{}
	})

	return nil
}
