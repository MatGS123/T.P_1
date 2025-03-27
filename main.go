package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	// Importación de la librería Fyne para crear la interfaz gráfica
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// Estructura que representa una parada de tren
type Parada struct {
	Nombre    string  // Nombre de la parada
	Distancia float64 // Distancia hasta la parada desde el inicio
}

// Estructura que representa un tren
type Tren struct {
	Nombre      string    // Nombre del tren
	Velocidad   float64   // Velocidad del tren en km/h
	Paradas     []Parada  // Lista de paradas que realiza el tren
	Prioritario bool      // Si el tren tiene prioridad
	Salida      time.Time // Hora de salida del tren
}

// Lista global que contiene todos los trenes
var trenes []Tren

// Variable global que representa la ventana principal de la aplicación
var ventana_principal fyne.Window

func main() {
	// Creación de la aplicación Fyne
	a := app.New()

	// Creación de la ventana principal
	ventana_principal = a.NewWindow("Simulador de Trenes")
	ventana_principal.Resize(fyne.NewSize(600, 600)) // Ajuste del tamaño de la ventana

	// Elementos de entrada de texto para que el usuario ingrese los datos del tren
	entrada_nombre := widget.NewEntry()
	entrada_nombre.SetPlaceHolder("Nombre del tren")

	entrada_velocidad := widget.NewEntry()
	entrada_velocidad.SetPlaceHolder("Velocidad (km/h)")

	entrada_paradas := widget.NewEntry()
	entrada_paradas.SetPlaceHolder("Paradas (nombre:distancia, separadas por comas)")

	entrada_salida := widget.NewEntry()
	entrada_salida.SetPlaceHolder("Hora de salida (HH:MM)")

	// Casilla de verificación para saber si el tren es prioritario
	check_prioritario := widget.NewCheck("Prioritario", nil)

	// Botón que agrega un tren a la lista
	boton_agregar := widget.NewButton("Añadir Tren", func() {
		// Recuperar los valores ingresados por el usuario
		nombre := strings.TrimSpace(entrada_nombre.Text)
		velocidad_texto := strings.TrimSpace(entrada_velocidad.Text)
		paradas_texto := strings.TrimSpace(entrada_paradas.Text)
		salida_texto := strings.TrimSpace(entrada_salida.Text)

		// Validación de campos vacíos
		if nombre == "" || velocidad_texto == "" || paradas_texto == "" || salida_texto == "" {
			mostrar_popup_log("Error: Todos los campos deben estar completos")
			return
		}

		// Convertir la velocidad de texto a número flotante
		velocidad, err := strconv.ParseFloat(velocidad_texto, 64)
		if err != nil || velocidad <= 0 {
			mostrar_popup_log("Error: Velocidad inválida")
			return
		}

		// Procesar las paradas (separadas por comas)
		paradas_datos := strings.Split(paradas_texto, ",")
		var paradas []Parada
		for _, p := range paradas_datos {
			partes := strings.Split(strings.TrimSpace(p), ":")
			if len(partes) != 2 {
				mostrar_popup_log("Error: Formato de paradas incorrecto")
				return
			}
			// Convertir la distancia de la parada a número flotante
			distancia, err := strconv.ParseFloat(strings.TrimSpace(partes[1]), 64)
			if err != nil || distancia < 0 {
				mostrar_popup_log("Error: Distancia inválida")
				return
			}
			paradas = append(paradas, Parada{Nombre: partes[0], Distancia: distancia})
		}

		// Convertir la hora de salida de texto a formato time
		salida, err := time.Parse("15:04", salida_texto)
		if err != nil {
			mostrar_popup_log("Error: Hora de salida inválida")
			return
		}

		// Crear el tren y añadirlo a la lista de trenes
		trenes = append(trenes, Tren{Nombre: nombre, Velocidad: velocidad, Paradas: paradas, Prioritario: check_prioritario.Checked, Salida: salida})
		mostrar_popup_log(fmt.Sprintf("Tren agregado: %s", nombre))
	})

	// Botón para simular el recorrido de los trenes
	boton_simular := widget.NewButton("Simular", func() {
		simular_trenes() // Llamada a la función de simulación
	})

	// Imagen que representa un tren (opcional)
	imagen := canvas.NewImageFromFile("tren.jpg")
	imagen.FillMode = canvas.ImageFillContain
	imagen.SetMinSize(fyne.NewSize(250, 150))

	// Diseño de la interfaz gráfica de la ventana principal
	ventana_principal.SetContent(container.NewVBox(
		entrada_nombre, entrada_velocidad, entrada_paradas, entrada_salida, check_prioritario, boton_agregar, boton_simular, imagen,
	))

	// Mostrar la ventana y ejecutar la aplicación
	ventana_principal.ShowAndRun()
}

// Función que simula los trenes y muestra los resultados
func simular_trenes() {
	// Si no hay trenes, mostrar un mensaje de error
	if len(trenes) == 0 {
		mostrar_popup_log("No hay trenes para simular")
		return
	}

	// Mensaje que se mostrará con la hora de salida y llegada de cada tren
	mensaje_log := "Orden de paso y tiempos:\n"
	for _, t := range trenes {
		// Calcular el tiempo requerido para cada tren
		tiempo_requerido := calcular_tiempo(t)
		// Calcular la hora de llegada sumando el tiempo requerido a la hora de salida
		tiempo_llegada := t.Salida.Add(time.Duration(tiempo_requerido * float64(time.Hour)))
		// Añadir los detalles del tren al mensaje
		mensaje_log += fmt.Sprintf("Tren %s: Sale a las %s, llega a destino a las %s\n", t.Nombre, t.Salida.Format("15:04"), tiempo_llegada.Format("15:04"))
	}

	// Analizar los cruces entre trenes
	mensaje_log += "\nPuntos de cruce entre trenes:\n"
	encontro_cruce := false
	for i := 0; i < len(trenes); i++ {
		for j := i + 1; j < len(trenes); j++ {
			// Calcular el tiempo y la distancia de intersección
			tiempo_cruce, distancia_cruce := calcular_interseccion_trenes(trenes[i], trenes[j])
			if !tiempo_cruce.IsZero() {
				// Mostrar la información sobre los cruces
				mensaje_log += fmt.Sprintf("Trenes %s y %s se cruzan a las %s en el km %.2f\n",
					trenes[i].Nombre, trenes[j].Nombre, tiempo_cruce.Format("15:04"), distancia_cruce)
				encontro_cruce = true
			}
		}
	}

	// Si no se encontraron cruces, mostrar un mensaje indicando esto
	if !encontro_cruce {
		mensaje_log += "No se encontraron cruces entre trenes.\n"
	}

	// Mostrar el mensaje en un popup
	mostrar_popup_log(mensaje_log)
}

// Función para calcular el tiempo de recorrido de un tren basado en su distancia total y velocidad
func calcular_tiempo(t Tren) float64 {
	distancia_total := 0.0
	for _, p := range t.Paradas {
		distancia_total += p.Distancia // Sumar la distancia de todas las paradas
	}
	// Tiempo = Distancia total / Velocidad
	return distancia_total / t.Velocidad
}

// Función para calcular el punto de cruce entre dos trenes
func calcular_interseccion_trenes(t1, t2 Tren) (time.Time, float64) {
	// Asegurarse de que el primer tren sea el que sale primero
	if t1.Salida.After(t2.Salida) {
		t1, t2 = t2, t1
	}
	deltaT := t2.Salida.Sub(t1.Salida).Hours() // Calcular la diferencia de tiempo de salida
	if deltaT < 0 {
		return time.Time{}, 0
	}

	v1, v2 := t1.Velocidad, t2.Velocidad
	if v1 == v2 {
		return time.Time{}, 0 // Si las velocidades son iguales, no habrá cruce
	}

	// Calcular la distancia y el tiempo de cruce
	distancia := (v2 * deltaT) / (1 - v1/v2)
	if distancia < 0 {
		return time.Time{}, 0
	}

	// Calcular el tiempo de cruce
	tiempo_cruce := t1.Salida.Add(time.Duration(distancia/v1) * time.Hour)
	return tiempo_cruce, distancia
}

// Función para mostrar un mensaje en un popup de información
func mostrar_popup_log(mensaje string) {
	// Mostrar el mensaje solo si la ventana está activa
	if ventana_principal != nil {
		dialog.ShowInformation("Log", mensaje, ventana_principal)
	}
}
