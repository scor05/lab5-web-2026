# Lab 5 Web
## Descripción
En este repositorio se creó un servidor simple de go para mostrar datos guardados en una
base de datos en formato de HTML.

El servidor se encuentra en `/go/server.go`, junto con todos los otros archivos que contienen
los handlers para todos los métodos que usa el server.

Para correr el servidor, simplemente correr en la terminal (dentro de `/go/server.go`) el comando
`go run .` y empezará a escuchar en el puerto 8080 y se mostrará ahí la tabla con la información dada.

## Challenges Implementados
De todos los challenges que estaban en el pool, implementé los siguientes: 
- Barra de progreso (ep actual vs totales)
- Marcar serie con texto especial si está terminada (muestra un checkmark al lado del nombre)
- Botón -1 (para decrementar episodios)
- Función para eliminar series
- Función para editar series
- Actualizar sin reload

## Preview
Al correr el servidor, se debería ver en `localhost:8080/` o `127.0.0.1:8080/` lo siguiente:
<img width="929" height="425" alt="image" src="https://github.com/user-attachments/assets/e489bac7-595b-4312-88fa-c89b6a486254" />
