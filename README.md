# Cameras-v2

New version of the camera project.

This version is:
 * in Go
 * using Redis:
    * for incomming commands from the server
    * to retrieve configuration files for the RPI
    * to push logs to the server
    * to push images to the server

The generated binary should be the only thing running on the RPI, and launched after boot