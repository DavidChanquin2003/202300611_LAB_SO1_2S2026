#!/bin/bash
echo "Inyectando modulo continfo.ko en el kernel..."
cd modulo_kernel
# Usamos insmod. Si falla (porque ya está cargado), no detiene el script.
insmod continfo.ko 2>/dev/null
echo "Modulo cargado exitosamente."