# TP0: Docker + Comunicaciones + Concurrencia

El paquete de la apuesta está compuesto por varios campos:

    Tamaño del paquete: longitud total en bytes del paquete a recibir. 2 bytes
    Longitud del nombre: 1 byte
    Nombre: longitud variable
    Longitud del apellido: 1 byte
    Apellido: longitud variable
    Longitud del documento: 1 byte
    Documento: longitud variable
    Longitud de la fecha de nacimiento: 1 byte
    Fecha de nacimiento: longitud variable
    Número: 2 bytes
    Número de agencia: 2 bytes

Podría haberse establecido una longitud fija para el documento y la fecha de nacimiento para hacer que el paquete fuera más claro. Al principio del paquete se usa el tamaño para poder leer el paquete completo y analizarlo internamente.

La respuesta consiste en un solo byte: 0 si se ha guardado correctamente y 1 si ha habido un error.

El tamaño máximo del paquete es de 1030 bytes, debido al máximo en los cuatro campos variables, más las longitudes de cada uno y los cuatro bytes de los datos numéricos. Por lo tanto, no se puede exceder el tamaño máximo de paquete indicado por la cátedra.

El protocolo y la clase Socket no son muy elegantes, pero son funcionales. Podrían ser un poco más simétricos para facilitar la programación.

----

Ahora se envian las apuestas en batches, al final de un batch hay 2 bytes que dicen si hay mas batches o no, esto es mas util para el proximo ejercicio. Pero en este caso se utilizan en el servidor para escribir las apuestas, enviar si todo salio bien y cerrar la conexion con el cliente.

Los batches son de tamanio variable y dependen del largo de las apuestas que contienen. Se puede configurar la cantidad maxima que puede haber en un batch, ya que no podemos saber a priori cuantos entran. 


----