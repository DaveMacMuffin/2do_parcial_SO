FROM ubuntu:24.04

#solo intala lo que se necesita para correr el godot en modo headless... o ese era el plan pero sin el frontcofig no iniciaba ¯\_(ツ)_/¯
RUN apt-get update && apt-get install -y fontconfig fonts-dejavu \
   # ca-certificates es para que requests de HTTPS se marquen como que confiamos en ellas e.g. decargar godot de github 
   # wget es para descargar cosas de internet
   # unzip se explica solo pero descomprime zips
    ca-certificates wget unzip \
   # -q para que no imprima nada en la consola
   # -O para que lo guarde en un archivo en la ruta que le digamos
   # usamos el wget para descargar godot de github 
    && wget -qO /tmp/godot.zip \
         https://github.com/godotengine/godot-builds/releases/download/4.4-stable/Godot_v4.4-stable_linux.x86_64.zip \
   # usamos el unzip para descomprimir el zip en la ruta que le decimos con -d
    && unzip /tmp/godot.zip -d /usr/local/bin/ \
   # eliminamos el zip porque ya no sirve :(
    && rm /tmp/godot.zip \
   # chmod es para cambiar permisos de archivo
   # +x es el permiso para que sea ejecutable
    && chmod +x /usr/local/bin/Godot_v4.4-stable_linux.x86_64 \
   # ln -s es para crear un 'symbolic link' 
   # estamos linkeando el archivo para que se pueda ejecutar con solo 'godot' en la terminal
    && ln -s /usr/local/bin/Godot_v4.4-stable_linux.x86_64 /usr/local/bin/godot

# CMD es para ejecutar un comando cuando se inicia el contenedor
# godot --headless es para correr godot en modo headless ese tambien se explicaba solo
CMD ["godot", "--headless"]


# para iniciarr el contenedor:
# navegar al directorio con el Dockerfile
# docker build -t imagen_examen -f cEx2.DockerFile .
# docker run -d --name contenedor imagen_examen

# para correr el contenedor en modo interactivo:
# docker run -it --name contenedor imagen_examen
# o si ya iniciamos el contenedor:
# docker exec -it contenedor /bin/bash
