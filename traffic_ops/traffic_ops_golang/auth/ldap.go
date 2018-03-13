package auth

import (
	"crypto/tls"
	"errors"
	"fmt"

	"github.com/apache/incubator-trafficcontrol/lib/go-log"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/config"

	"gopkg.in/ldap.v2"
)

func LookupUserDN(username string, cfg *config.ConfigLDAP) (string, bool, error) {
	l, err := ldap.DialTLS("tcp", cfg.Host, &tls.Config{InsecureSkipVerify: true})
	if err != nil {
		log.Errorln("error dialing tls")
		return "", false, err
	}
	defer l.Close()
	// Bind with admin user
	err = l.Bind(cfg.AdminDN, cfg.AdminPass)
	if err != nil {
		log.Errorln("error binding admin user")
		return "", false, err
	}

	// Search for the given username
	searchRequest := ldap.NewSearchRequest(
		cfg.SearchBase,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(&(objectCategory=person)(objectClass=user)(sAMAccountName=%s))", username),
		[]string{"dn"},
		nil,
	)

	sr, err := l.Search(searchRequest)
	if err != nil {
		log.Errorln("error issuing search:")
		return "", false, err
	}

	if len(sr.Entries) != 1 {
		return "", false, errors.New("User does not exist or too many entries returned")
	}
	userDN := sr.Entries[0].DN
	return userDN, true, nil
}

func AuthenticateUserDN(userDN string, password string, cfg *config.ConfigLDAP) (bool, error) {
	l, err := ldap.DialTLS("tcp", cfg.Host, &tls.Config{InsecureSkipVerify: true})
	if err != nil {
		log.Errorln("error dialing tls")
		return false, err
	}
	defer l.Close()

	// Bind as the user to verify their password
	err = l.Bind(userDN, password)
	if err != nil {
		return false, err
	}
	return true, nil
}
