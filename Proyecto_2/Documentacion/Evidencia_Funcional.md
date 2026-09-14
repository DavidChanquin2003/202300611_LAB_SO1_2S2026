# Evidencia Funcional y Capturas de Prueba
**Proyecto 2 - Sonda de Kernel y Daemon en Go**

En este documento se presentan las pruebas gráficas que certifican el correcto funcionamiento de todos los componentes de la arquitectura requerida para el proyecto.

### 1. Interacción a Nivel de Kernel (C y eBPF)
*   **Captura 1 (Archivo /proc):** 
    *`cat /proc/continfo_pr2_so1_202300611`*

![Modulo C](img/moduloc.png)

*   **Descripción:** Se evidencia la correcta carga del módulo en C, mostrando la captura directa de las estructuras de memoria del sistema y la lista de procesos activos.



*   **Captura 2 (Intercepción eBPF):**

![logs del Daemon en Go](img/DaemonGO.png)

*   **Descripción:** Muestra al Daemon recibiendo la confirmación de la llamada al sistema `sys_kill` interceptada por eBPF al momento de eliminar un proceso.

### 2. Orquestación y Toma de Decisiones (Go)
*   **Captura 3 (Loop Principal):**

![Loop](img/Loop_Principal.png)

*   **Descripción:** Ejecución del Daemon gestor operando en el ciclo exigido de 20 segundos.

*   **Captura 4 (Lógica de Eliminación):**

![eliminacion](img/eliminacion.png)


*   **Descripción:** Prueba de la toma de decisiones autónoma para mantener estrictamente 3 contenedores de bajo consumo y 2 de alto consumo.


### 3. Simulación de Carga e Infraestructura
*   **Captura 5 (Cronjob):**

![Cronjob](img/Cronjob.png)

*   **Descripción:** Verificación de la tarea programada encargada de disparar el script generador cada minuto.


*   **Captura 6 (Límites en Docker):**
    *(Inserta aquí tu imagen ejecutando `docker stats`)*

![Límites](img/Límites.png)


*   **Descripción:** Observación en vivo de los contenedores donde se aprecia el blindaje de RAM aplicado para evitar la saturación del sistema operativo.

### 4. Tablero de Observabilidad (Grafana)
*   **Captura 7 (Dashboard Completo):**
    *(Inserta aquí tu captura de pantalla final de Grafana)*

![Grafana](img/Grafana.png)

*   **Descripción:** Visualización de la interfaz final consumiendo los datos desde Valkey, incluyendo contadores, gráficas de pastel para los Tops y tablas de evolución temporal.