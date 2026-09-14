# Manual Técnico - Sonda de Kernel en C y Daemon en Go para la Telemetría de Contenedores

**Proyecto 2 - Sistemas Operativos 1**

**Estudiante:** David Arturo Chanquin Velasquez - 202300611

## 1. Descripción General de la Arquitectura

El sistema es una solución integral de observabilidad y gestión autónoma para contenedores Docker. Está compuesto por microservicios que interactúan desde el espacio de kernel hasta el espacio de usuario, permitiendo el monitoreo en tiempo real del uso de recursos y la toma de decisiones automatizada basada en telemetría. La arquitectura se divide en cinco componentes principales:

**Módulo de Kernel (C):** Recolector de bajo nivel.

**Sonda eBPF:** Auditor de llamadas al sistema.

**Daemon (Go):** Cerebro y gestor del flujo lógico.

**Base de Datos (Valkey):** Almacenamiento persistente en memoria.

**Dashboard (Grafana):** Visualización de la evolución temporal.

## 2. Módulo de Kernel en C

Este componente actúa como un sensor de bajo nivel que expone métricas avanzadas a través del sistema de archivos virtual en la ruta /proc/continfo_pr2_so1_202300611

**Implementación:** Se hace uso de la estructura interna task_struct del núcleo de Linux para iterar sobre los procesos activos.

**Extracción de Datos:** Por cada proceso, se capturan identificadores (PID, Nombre), y parámetros de memoria (VSZ, RSS) realizando cálculos diferenciales para determinar el porcentaje de uso de memoria y CPU. Adicionalmente, se recolecta el estado global de la memoria RAM (Total, Libre y En uso).

![Modulo C](img/moduloc.png)

## 3 Observabilidad con eBPF

El programa eBPF funciona de manera asíncrona y está dedicado exclusivamente a auditar las acciones de eliminación emitidas por el sistema.

**Intercepción:** Se adjuntó la sonda a la llamada al sistema sys_kill para interceptar las señales de terminación (SIGTERM/SIGKILL).

**Flujo de confirmación:** eBPF captura el PID del proceso receptor y envía esta información al Daemon en Go, sirviendo como confirmación real a nivel de núcleo de que un contenedor ha sido exterminado exitosamente.

## 4. Daemon Orquestador en Go

Es el programa principal que procesa la telemetría y ejecuta las políticas de control operando en un loop infinito configurado a 20 segundos.

**Parseo e Ingesta:** En cada ciclo, lee y deserializa el archivo generado por el módulo en /proc.

**Evaluación de Estado:** El script generador administrado por cronjob crea 5 contenedores de carga aleatoria cada minuto. El daemon evalúa el estado del sistema mediante el consumo real de RSS y CPU para clasificar los contenedores en Alto y Bajo consumo.

**Toma de Decisiones:** El algoritmo garantiza que siempre existan exactamente 3 contenedores de bajo consumo y 2 de alto consumo. Todo contenedor que exceda este umbral máximo es asesinado sistemáticamente mediante comandos de Docker.

**Conexión DB:** Una vez confirmada la eliminación a través de eBPF, el gestor registra los eventos y los tops de consumo (Top 5 RAM, Top 5 CPU) en estructuras de datos de Valkey.

![logs del Daemon en Go](img/DaemonGO.png)

## 5. Persistencia y Visualización

La infraestructura de persistencia se levantó utilizando docker-compose, inyectando volúmenes físicos (valkey_data y grafana_data) para asegurar que la configuración y los historiales sobrevivan a la destrucción de los contenedores efímeros. El tablero en Grafana consume los datos directamente de Valkey mediante el plugin redis-datasource, mostrando en tiempo real los contadores eBPF y la evolución temporal del consumo de RAM y contenedores eliminados.

![blindaje](img/blindaje.png)