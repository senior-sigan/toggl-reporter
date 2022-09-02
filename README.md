



## How to build 

1. Copy config.example.toml into config.toml
2. Build project with make
3. 


## How to make new page

1. Copy existing template (e.g `index.toml`) from `web/templates` in same folder.  
   Name accordingly.  
2. Register route in `web/main.go` and add new location.
3. Add functions for each request type for location and other necessary things.
