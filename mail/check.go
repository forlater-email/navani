package mail

import (
	"fmt"
	"log"
	"net"
	"net/mail"
	"strings"

	"blitiri.com.ar/go/spf"
	"golang.org/x/net/idna"
)

func verifySPF(ips []net.IP, domain string, from string) (passes []spf.Result, err error) {
	passes = []spf.Result{}
	for _, ip := range ips {
		result, _ := spf.CheckHostWithSender(ip, domain, from)
		// if err != nil {
		// 	return nil, fmt.Errorf("check spf: %w\n", err)
		// }
		if result == spf.Pass {
			passes = append(passes, result)
		}
	}
	return passes, nil
}

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

	mxPasses := []spf.Result{}
	// Check MX IPs against SPF record.
	for _, mx := range mxs {
		ips, err := net.LookupIP(mx.Host)
		if err != nil {
			return false, fmt.Errorf("ip lookup: %w\n", err)
		}
		passes, err := verifySPF(ips, domain, e.Address)
		if err != nil {
			return false, fmt.Errorf("mx spf: %w\n", err)
		}
		mxPasses = append(mxPasses, passes...)
	}

	if len(mxPasses) == 0 {
		// Check domain IP against SPF record.
		ips, err := net.LookupIP(domain)
		if err != nil {
			return false, fmt.Errorf("ip lookup: %w\n", err)
		}
		passes, err := verifySPF(ips, domain, e.Address)
		if err != nil {
			return false, fmt.Errorf("domain spf: %w\n", err)
		}
		// If both MX IP and domain IP fail SPF, we return false.
		if len(passes) == 0 {
			return false, nil
		}
	}

	log.Printf("spf passed: %s\n", email)
	return true, nil
}
