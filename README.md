## StudyNotes Backend
###### Que es?
MemoryCards es una aplicación para acompañar todo tipo de aprendizaje que requiera memorizacion. <br/>

###### Tecnologías usadas:
- Golang
- Echo framework
- JWT
- MySQL GORM

#### Como bajar y correr el proyecto en local
###### Dependencias:
- Go
- Git
- MySQL

###### Pasos:
1) Clonar el repositorio con el comando: <br/> `https://github.com/CardozoCasariegoLuciano/StudyNotes-backend`
2) Pararse en la raíz del proyecto e instalar las dependencias con `go mod download`
3) Tener corriendo MySQL
4) Crear una base de datos llamada `studyNotes`
5) Crear un archivo .env en la raiz del proyecto y escribir <br/>
`DB_PASSWORD=[tu clave de conexion a la BD]` <br/>
`DB_HOST=localhost:3306`<br/>
`DB_NAME=studyNotes `<br/>

- Nota: por defecto el usuario es root, si tenes uno distinto lo podes modificar en `/database/db_config.go` <br/>
6) Levantar el proyecto con `go run main.go`

- Nota: opcionalmente se puede levantar el [Frontend](https://github.com/CardozoCasariegoLuciano/StudyNotes-Frontend/tree/develop)
     para interactuar con la aplicación como lo haría el usuario final

#### Como trabajamos
En la aplicación estamos aplicando el git flow lo que implica:

1) cuando se esta desarrollando una funcionalidad o característica nueva tenemos que
crear una rama nueva a partir de `develop`
2) Esa rama que creemos tiene que tener como nombre la incidencia de Jira (ver en Links importantes)
de la tarea que estemos tratando
- Ejemplo `SNTAB-1232`
3) Cada commit que hagamos tiene que comenzar con el nombre de la rama, seguido del
    resumen de lo realizado en el commit
4) Una vez terminada la funcionalidad la pusheamos al repositorio remoto
    y desde ahi creamos una Pull Request a `develop`
5) Las Pull Request tienen que ser revisadas y aprobadas para poder realizar el merge
    a la rama develop

#### Links importantes
- [swager](http://localhost:5000/swagger/index.html):
   documentación de la API del backend (requiere tener el backen en local corriendo)
- [Jira](https://studynotes-project.atlassian.net/jira/software/projects/SNTAB/boards/1):
    Estado de las tareas actuales
- [Figma](https://www.figma.com/file/ifSIZqKuHld2q15debAYky/StudiesNotesApp?node-id=115%3A445&t=zR3M3cv3vdtcF78P-1):
    Diseño de la aplicación
- [FigJam](https://www.figma.com/file/5JzllwcWgURAKeaQvyBXMs/StudyNotes?node-id=0%3A1&t=MGov6Z7RawYZ7i7q-1):
    Estructura y planes a futuro de la aplicación

#### Hot reload
- Recomendamos instalar [fresh](https://github.com/gravityblast/fresh), una utilidad que vuelve a levantar el proyecto cada vez que se realiza un cambio


#### Test
- Opcionalmente se puede instalar [gotestsum](https://github.com/gotestyourself/gotestsum)
para visualizar los test de una forma mas cómoda

Para correr todos los test: `go test -v ./...`

#### Mocks

###### Actualizar los mocks:
- Cuando se modifique la interfaz Istorage IMPORTANTE actualizar los mocks ejecutando el script:<br/>
`./updateMocks.sh` <br/>
tanto para tenerlos al momento de hacer los test como para que no se rompan los que ya estan


#### Actualizar swagger
- Para actualizar la documentación de swager ejecutar el script: `./updateSwager.sh`

