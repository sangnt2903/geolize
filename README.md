# Geolize

Geolize is a geolocation service that provides IP lookup and modification capabilities. It leverages MaxMind's GeoLite2 database to offer detailed geolocation data.

## Features

- IP address lookup with detailed geolocation information
- Support for updating geolocation data
- Integration with gRPC and HTTP servers
- CORS support for cross-origin requests

## Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/yourusername/geolize.git
   cd geolize

2. Setup the environment:
    
Make sure you have installed the Docker and Docker Compose.
    
    ```bash
    make bootstrap
    ```

3. Run the service:
    
    ```bash
    make run service=geolize
    ```

## How to update the api contract

1. Make changes to the protobuf files in the geolize directory `service-protos/services/geolize/service.proto`.

2. Generate the Go code:
   
    ```bash
    make gen service=geolize
    ```
## When I need to update the database to the latest version?

1. You will need to stop the service first:

2. Mount the new GeoLite2 database file in the `data/db` directory.

3. Set Geolize DB to dev.ini file with the name of the new database file.

For example: If database file name is `GeoLite2-City.mmdb`, then you need to update the `dev.ini` file with the following:
    
    ```ini
    [geolize]
    db=GeoLite2-City.mmdb
    ```

4. Empty the version file data/version

5. Run the service will affect and sync all history files to the new database file.

## Contributions
Contributions are welcome! If you find any issues or have suggestions for improvements, please open an issue or submit a pull request.

## License
This project is licensed under the MIT License. See the LICENSE file for more information.
