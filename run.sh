docker run -d -p 53:53/udp -e PROXY_HOST=pc-server  --restart=unless-stopped --name dnsserver sowisz/dnsserver

