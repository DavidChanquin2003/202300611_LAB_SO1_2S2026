FROM alpine:latest
CMD sh -c "while true; do cat /etc/shadow 2>/dev/null || echo 'Intentando leer archivo sensible...'; sleep 3; done"