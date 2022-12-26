



## How to build 

1. Copy config.example.toml into config.toml
2. Set variables in config.toml
3. Build project with make

## How to launch

1. Launch the binary
2. Open page on address defined in `addr` config param (e.g localhost:8000)


## How to make new page

1. Copy existing template (e.g `index.toml`) from `web/templates` in same folder.  
   Name accordingly.  
2. Register route in `web/main.go` and add new location.
3. Add functions for each request type for location and other necessary things.
