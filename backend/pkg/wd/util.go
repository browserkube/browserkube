package wd

import (
	"fmt"
	"net"
	"strings"

	"github.com/pkg/errors"
)

func ParseSessionPath(uPath string) (string, string, error) {
	p := RemoveBase(uPath)
	parts := strings.SplitN(p, "/", 4)
	length := len(parts)
	if length < 3 {
		return "", "", errors.Errorf("unable to parse session id from url: %s", uPath)
	}
	sessionID := parts[2]
	command := ""
	if length > 3 {
		command = "/" + parts[3]
	}
	return sessionID, command, nil
}

// RemoveBase removes WD base path from URL
func RemoveBase(uPath string) string {
	return strings.TrimPrefix(uPath, "/wd/hub")
}

// ReplaceSession replaces session in the given URL path
func ReplaceSession(path, session string) (string, error) {
	parts := strings.SplitN(path, "/", 4)
	if len(parts) < 3 {
		return "", errors.Errorf("unable to parse session from %s", path)
	}
	command := ""
	if len(parts) > 3 {
		command = "/" + parts[3]
	}

	return "/session/" + session + command, nil
}

// BuildHostPort ...
func BuildHostPort(session, service, port string) string {
	return net.JoinHostPort(fmt.Sprintf("%s.%s", session, service), port)
}
