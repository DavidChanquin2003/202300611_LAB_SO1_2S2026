# Guía de Instalación y Ejecución
## Proyecto 2 - Sonda de Kernel y Daemon en Go

Esta guía detalla los pasos exactos para compilar, levantar, probar y apagar de forma segura la arquitectura de monitoreo de contenedores. 

### Requisitos Previos
El entorno de ejecución requiere una distribución de Linux (se recomiendan sistemas basados en Debian como Ubuntu o Linux Mint) con los siguientes paquetes instalados:
*   **Motor de Docker y Docker Compose**
*   **Lenguaje Go** (versión 1.20 o superior)
*   **GCC y Linux Headers** (para compilar el módulo en C)
*   **Make**

---

### 1. Arranque del Sistema (Orden Estricto)

Para garantizar la correcta comunicación entre microservicios, los componentes deben iniciarse en el siguiente orden:

#### Paso 1.1: Levantar la Infraestructura (Base de Datos y Dashboard)
Abra una terminal en la carpeta raíz del proyecto (donde se encuentra el archivo `docker-compose.yml`) y ejecute:
```bash
docker compose up -d
```

#### Paso 1.2: Compilar y Cargar el Módulo del Kernel (C)

Navegue hacia la carpeta donde se encuentra el código fuente del módulo en C y compile utilizando el archivo Makefile provisto:

```
sudo insmod continfo.ko
```

(Nota: Si su archivo objeto tiene un nombre distinto, reemplace continfo.ko por el nombre compilado correspondiente).

#### Paso 1.3: Iniciar el Gestor y Orquestador (Go)

Abra una nueva pestaña en la terminal, navegue a la carpeta del daemon de Go (ej. daemon_go) y ejecute el servicio con privilegios de superusuario para permitir la inyección de la sonda eBPF y la lectura de /proc:

```
sudo go run .
```

El sistema iniciará el loop principal de 20 segundos y registrará automáticamente el cronjob en el sistema.

### 2. Comandos de Prueba y Verificación

Una vez que el sistema esté corriendo, puede utilizar los siguientes comandos desde cualquier ubicación en la terminal para verificar el correcto funcionamiento de cada capa de la arquitectura:

**Verificar lectura directa del Kernel:** 
```
cat /proc/continfo_pr2_so1_202300611
```

**Verificar la inyección del Cronjob (cada 1 minuto):**

```
sudo crontab -l
```

**Monitorear la gestión de memoria y eliminación de contenedores en vivo:**

```
docker stats
```

**Auditar la persistencia de datos en Valkey:**

```
docker exec -it valkey_db valkey-cli

KEYS *
```

### 3. Apagado

Para evitar fugas de memoria, contenedores huérfanos o kernel panics, es obligatorio detener el sistema en el siguiente orden:

**1. Detener el Daemon de Go:**
Vaya a la terminal donde se está ejecutando Go y presione Ctrl + C. Esto ejecutará la función de limpieza (cleanup()) que elimina el cronjob del sistema operativo.

**2. Eliminar contenedores residuales de carga:**
Ejecute el siguiente comando para limpiar cualquier contenedor generado por el script que haya quedado activo:

```
docker rm -f $(docker ps -aq --filter "name=bajo_consumo" --filter "name=alto_" --filter "name=intruso") 2>/dev/null || true
```

**3. Apagar la infraestructura principal:**
Navegue a la raíz del proyecto y detenga los servicios de persistencia:

```
docker compose down
```
**4. Descargar el módulo del Kernel:**
Retire de forma segura el módulo de la memoria del sistema operativo:

```
sudo rmmod continfo
```