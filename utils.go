package sshdocker

import "errors"

var _DOCKERFILE = []byte(`
FROM ubuntu:17.10

RUN apt-get update && apt-get install -y openssh-server vim
RUN mkdir /var/run/sshd
RUN echo 'root:root' | chpasswd
RUN sed -i 's/#PermitRootLogin prohibit-password/PermitRootLogin yes/' /etc/ssh/sshd_config

# SSH login fix. Otherwise user is kicked off after login
RUN sed 's@session\s*required\s*pam_loginuid.so@session optional pam_loginuid.so@g' -i /etc/pam.d/sshd

EXPOSE 22
CMD ["/usr/sbin/sshd", "-D"]
`)

func ArgsCountCheck(argsCount, min, max int) error {
	var err error
	if min > 0 && argsCount < min {
		err = errors.New("Too few arguments")
	}
	if max > 0 && argsCount > max {
		err = errors.New("Too many arguments")
	}
	return err
}
