# DEMO para presentación de gRPC

Este proyecto contiene la definición básica de cliente (`client.go`) y servidor (`server.go`) 
para lograr una comunicación usando gRPC.
Posee un ejemplo para cada una de los tipos de llamadas:
- Unario o Única
- Streaming del cliente
- Streaming del servidor
- Streaming bidireccional

Contiene además la definición de los mensajes y el contrato en el archivo `item.proto`. 
Además el comando para poder generar los códigos correspondientes (`generate.sh`).