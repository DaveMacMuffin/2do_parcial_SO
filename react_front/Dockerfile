# usamos la imgen lts de node porque segun san google es la ulima version de 'Long term support' (LTS)
# especificamos que sea alpine poque segun es segura y no pesa mucho tambien porque el nginx que usamos es base alpine
FROM node:lts-alpine AS build

# ver goApi para explicacion
WORKDIR /app

# ver goApi para explicacion
COPY . .

# npm = node package manager
# install para instala las dependencias para poder usar node
RUN npm install


# copila el proyecto
RUN npm run build

# usamos la imgen de nginx porque es un servidor web y es lo que necesitamos para correr el front
FROM nginx:alpine

# copiamos el build que hicimos en la imagen anterior al contenedor
COPY --from=build /app/build /usr/share/nginx/html

# el CMD es para correr nginx en primer plano porue si no se cierra el contenedor
# -g es para pasarle un argumento a nginx
# daemon off es para que corra en primer plano
CMD ["nginx", "-g", "daemon off;"]

# no especifico un puerto porque voy a usar el default de nginx: 80