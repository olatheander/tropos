FROM ubuntu

# Install SSHD
RUN apt-get update && apt-get install -y openssh-server sshfs
RUN mkdir /var/run/sshd
# Only allow key-based login
RUN passwd -d root
RUN sed -ri 's/^#?PermitRootLogin\s+.*/PermitRootLogin yes/' /etc/ssh/sshd_config
RUN sed -ri 's/UsePAM yes/#UsePAM yes/g' /etc/ssh/sshd_config
RUN mkdir /root/.ssh

COPY scripts /scripts
RUN chmod +x /scripts/*.sh

RUN apt-get clean && \
    rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

EXPOSE 22
VOLUME /workspace /root/.ssh/authorized_keys

ENTRYPOINT ["/usr/sbin/sshd", "-D"]
