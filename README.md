# Simulador de Trenes

Este proyecto es una simulación de trenes con varias funcionalidades, como la programación de trenes con distintas características, la simulación de su llegada a destino, y la detección de cruces entre trenes en base a su velocidad y tiempo de salida.

Características

- Añadir trenes: Puedes agregar trenes con sus respectivos parámetros:
  - Nombre del tren.
  - Velocidad del tren (en km/h).
  - Paradas programadas del tren (nombre y distancia).
  - Hora de salida.
  - Opción para marcar el tren como prioritario.

- Simulación: Al simular, el programa calcula:
  - La hora de llegada de cada tren.
  - Los puntos de cruce entre trenes, si es que se cruzan en el trayecto.

## Requisitos

- Go 1.18 o superior.
- Biblioteca **Fyne** para la interfaz gráfica (GUI).

## Instalación

1. Clonar el repositorio

2. Instalar dependencias: Asegúrate de tener instalada la biblioteca Fyne. Para mas información dirígase a: https://docs.fyne.io/started/

go get fyne.io/fyne/v2

Es recomendable ejecutar go mod tidy para asegurarte de que todas las dependencias estén actualizadas y las que no se usen sean eliminadas

3. Compilar y ejecutar el proyecto:

go run main.go

Esto abrirá la ventana de la interfaz gráfica donde podrás agregar trenes, simular y ver los resultados.
