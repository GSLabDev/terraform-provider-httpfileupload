package main

import "net/http"

type Config struct{
        ip_address string
}

func (c *Config) GetResponse(request *http.Request) (error){

        return nil
}
