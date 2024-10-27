# Minecraft User Presence Whatsapp Identifier

This service is divided in two main parts:

1. This service will poll the latest.log minecraft file via ftp and parse the content looking for joins/leaves on the minecraft server.
2. This service will send a whatsapp message to a chat group with the name of the user who joined or left.

- [x] Connect to the minecraft server FTP.
- [x] Poll latest.log file via FTP
- [x] Send a message to a whatsapp group.
- [x] Send a message to a whatsapp group when someone joins or leaves.
