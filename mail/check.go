package mail

import (
	"fmt"
	"net"
	"net/mail"
	"strings"

	"blitiri.com.ar/go/spf"
	"golang.org/x/net/idna"
)

func VerifySPF(email string) (bool, error) {
	e, err := mail.ParseAddress(email)
	if err != nil {
		return false, fmt.Errorf("parse address: %w", err)
	}
	domain := strings.Split(e.Address, "@")[1]

	domain, err = idna.ToASCII(domain)
	if err != nil {
		return false, fmt.Errorf("to ascii: %w\n", err)
	}

	mxs, err := net.LookupMX(domain)
	if err != nil {
		return false, fmt.Errorf("mx lookup: %w\n", err)
	}

	passes := []spf.Result{}
	for _, mx := range mxs {
		ips, err := net.LookupIP(mx.Host)
		if err != nil {
			return false, fmt.Errorf("ip lookup: %w\n", err)
		}

		for _, ip := range ips {
			result, _ := spf.CheckHostWithSender(ip, domain, e.Address)
			if result == spf.Pass {
				passes = append(passes, result)
			}
		}
	}

	if len(passes) > 0 {
		return true, nil
	} else {
		return false, nil
	}
}
