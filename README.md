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

He encontrado un error inusual y he modificado uno de los archivos correspondientes para poder continuar avanzando. Al calcular la posición final del archivo 3, he obtenido una posición anterior a la posición real final del archivo. Este problema se ha presentado únicamente con este archivo. He sustituido dicho archivo por otro para poder continuar con la sección 8.

El protocolo utilizado ha sido prácticamente el mismo que en el ejercicio anterior. Se han agregado mensajes para preguntar si el proceso ha finalizado, en los que el cliente envía su "id + 2000" en 2 bytes en big endian. El servidor responde con un valor de 65535 en 2 bytes en big endian en caso de que el sorteo no haya finalizado. En ese caso, el cliente cierra la conexión y se va a dormir, para luego volver a intentarlo.

En caso de que el servidor haya realizado el sorteo, éste responderá con la cantidad de apuestas en las que la agencia ha ganado, seguida del largo del DNI y el DNI correspondiente a cada apuesta ganadora.