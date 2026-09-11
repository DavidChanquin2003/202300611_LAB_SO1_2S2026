#!/bin/bash

# Script para generar los 5 contenedores aleatorios
# Tipos de contenedores: 
# 0 = Bajo Consumo (alpine sleep)
# 1 = Alto Consumo RAM (alto_ram)
# 2 = Alto Consumo CPU (alto_cpu)
# 3 = Contenedor Intruso (intruso_img)

for i in {1..5}; do
    RANDOM_TYPE=$((RANDOM % 4))
    TIMESTAMP=$(date +%s)
    
    case $RANDOM_TYPE in
    0)
        echo "Creando contenedor de Bajo Consumo..."
        docker run -d --name "bajo_consumo_${TIMESTAMP}_${i}" alpine sleep 240
        ;;
    1)
        echo "Creando contenedor de Alto Consumo RAM..."
        docker run -d --name "alto_ram_${TIMESTAMP}_${i}" --memory="300m" alto_ram
        ;;
    2)
        echo "Creando contenedor de Alto Consumo CPU..."
        docker run -d --name "alto_cpu_${TIMESTAMP}_${i}" --memory="200m" alto_cpu
        ;;
    3)
        echo "Creando Contenedor Intruso..."
        docker run -d --name "intruso_${TIMESTAMP}_${i}" --memory="200m" intruso_img
        ;;
esac
done

echo "¡se han creado 5 contenedores de forma exitosa vos!"