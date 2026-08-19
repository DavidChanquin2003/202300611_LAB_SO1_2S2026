# Proyecto 1 - Sistemas Operativos 1 (2S2026)

### Estudiante: David Arturo Chanquin Velàsquez | Carnet: 202300611

Este repositorio contiene el código fuente, los Dockerfiles y la guía de configuración para el despliegue de un entorno virtualizado basado en microservicios (APIs) escritos en Go, hechos a través de múltiples motores de contenedores (Containerd, Podman, Docker) y KVM.

## Estructura del Repositorio

202300611_LAB_SO1_2S2026/
│
Proyecto 1
│
├── api1/
│   ├── main.go          # Código fuente de la API 1
│   └── Dockerfile       # Archivo de construcción multi-stage
│
├── api2/
│   ├── main.go          # Código fuente de la API 2
│   └── Dockerfile       # Archivo de construcción multi-stage
│
├── api3/
│   ├── main.go          # Código fuente de la API 3
│   └── Dockerfile       # Archivo de construcción multi-stage
│
├── diagrama.png
│
├── manualtecnico.md
│
└── README.md            # Guía de instalación y configuración

## Requisitos Previos

Para ejecutar este proyecto, se requiere un sistema anfitrión con Linux y KVM habilitado. Se deben crear tres máquinas virtuales en la misma subred con la siguiente configuración:

#### VM1: Containerd instalado de forma nativa.

#### VM2: Podman instalado.

#### VM3: Docker Engine instalado.

## Guía de Instalación y Despliegue

#### Paso 1: Clonar el repo

En la compu anfitriona o en las máquinas virtuales destinadas a la compilación, clonar este repositorio:

** git clone https://github.com/DavidChanquin2003/202300611_LAB_SO1_2S2026.git

** cd 202300611_LAB_SO1_2S2026

#### Paso 2: Despliegue del Registro Zot (VM3)

El registro privado Zot centralizará las imágenes de los contenedores. Ingresar a la VM3 y ejecutar:

** docker run -d -p 5000:5000 --name zot ghcr.io/project-zot/zot-linux-amd64:latest

#### Paso 3: Construcción y Subida de Imágenes:

Navegar a cada uno de los directorios de las APIs (api1, api2, api3) y construir las imágenes utilizando el formato [ API#-#CARNET ]. Luego, empujarlas al registro Zot alojado en la VM3 (reemplazar <IP_VM3> por la IP correspondiente).

**Ejemplo para la API 1:**

*** cd api1

*** docker build -t <IP_VM3>:5000/api1-202300611:latest .

***docker push <IP_VM3>:5000/api1-202300611:latest

(Se repite este proceso para la API 2 y API3)

#### Paso 4: Despliegue de Contenedores en VMs Finales

Una vez que las imágenes estén en el registro Zot, se deben descargar y ejecutar en sus respectivas máquinas virtuales utilizando su motor correspondiente.

** En la VM1 (Containerd - API 1 y API 2):

- sudo ctr images pull --plain-http <IP_VM3>:5000/api1-202300611:latest
- sudo ctr run -d --net-host <IP_VM3>:5000/api1-202300611:latest api1_container

- sudo ctr images pull --plain-http <IP_VM3>:5000/api2-202300611:latest
- sudo ctr run -d --net-host <IP_VM3>:5000/api2-202300611:latest api2_container

** En la VM2 (Podman - API 3):

- podman run -d -p 8083:8083 --name api3 <IP_VM3>:5000/api3-202300611:latest

## Pruebas de Funcionamiento

Una vez levantados los contenedores, se puede validar el estado de los servicios mediante peticiones HTTP al endpoint /health.

** Desde la compu anfitriona se abre una terminal y luego, ejecutar:

- curl -s http://<IP_VM1>:8081/health | jq
- curl -s http://<IP_VM1>:8082/health | jq
- curl -s http://<IP_VM2>:8083/health | jq

** Para probar la comunicación cruzada entre las APIs, utilizar las rutas especificadas en el diseño tenemos este ejemplo: Ejemplo llamando a la API 2 desde la API 1:

- curl -s http://<IP_VM1>:8081/api1/202300611/call-api2 | jq

