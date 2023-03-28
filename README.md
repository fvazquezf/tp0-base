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