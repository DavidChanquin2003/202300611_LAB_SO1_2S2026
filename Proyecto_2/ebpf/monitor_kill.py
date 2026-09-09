#!/usr/bin/python3
from bcc import BPF

# Programa en C de eBPF
bpf_text = """
#include <uapi/linux/ptrace.h>
#include <linux/sched.h>

// Estructura para enviar los datos a Python/Go
struct data_t {
    u32 pid_emisor;
    u32 pid_receptor;
    int sig;
};

// Canal de comunicacion
BPF_PERF_OUTPUT(events);

// Tracepoint para la llamada sys_kill
TRACEPOINT_PROBE(syscalls, sys_enter_kill) {
    struct data_t data = {};
    
    // bpf_get_current_pid_tgid() devuelve el PID en los 32 bits altos
    data.pid_emisor = bpf_get_current_pid_tgid() >> 32;
    data.pid_receptor = args->pid;
    data.sig = args->sig;
    
    // Filtramos para capturar solo señales de terminación (SIGTERM = 15, SIGKILL = 9)
    if (data.sig == 9 || data.sig == 15) {
        events.perf_submit(args, &data, sizeof(data));
    }
    
    return 0;
}
"""

# Iniciar BPF
b = BPF(text=bpf_text)

print("Sonda eBPF iniciada. Interceptando sys_kill (Señales 9 y 15)... Esperando eventos.")
print("%-15s %-15s %-10s" % ("PID_EMISOR", "PID_RECEPTOR", "SEÑAL"))

# Función para procesar los eventos capturados
def print_event(cpu, data, size):
    event = b["events"].event(data)
    print("%-15d %-15d %-10d" % (
        event.pid_emisor, 
        event.pid_receptor, 
        event.sig
    ))

# Loop de escucha
b["events"].open_perf_buffer(print_event)
while 1:
    try:
        b.perf_buffer_poll()
    except KeyboardInterrupt:
        exit()