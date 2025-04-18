# usamos la imgen de go 1.24 porue segun San-google es la ultima version
FROM golang:1.24 

# es alparcer una practica estandar poner un dirctorio para el proyecto
WORKDIR /app

# copia lo que hay en el directirio donde esta el file al contenedor 
COPY . .

# gomod tidy medio se explica solo. se asegura de que las dependecias sean las correctas y borra las que no siven por eso 'tidy'
RUN go mod tidy
# copila el proyecto en y lo guarda en main| el '.' le dice que lo haga en el directorio actual
RUN go build -o main .

# ponemos este puerto para que lo use el API, arbitrariamente decidi usar el mismo del proyecto por la misma razon no se pudo 8080 sera
EXPOSE 8080

# el CMD esta vez corre el main que creamos arriba
CMD ["./main"]