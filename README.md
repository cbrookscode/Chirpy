# Chirpy

## Mux.Handle vs mux.HandleFunc
- Handle requires that the second arguement implements the http.Handler interface where handlefunc does not. Handlefunc will take the function provided and make it impelemnt the http.Handler interface. Handle is used for when you want to implement a custom router, middleware, or other custom logic. If you want to make a function into a http.Handler then you can use the http.HanlderFunc() function to do so. Ex: middlewareMetricsInc method on apiConfig.
- Using the Go standard library, you can specify a method like this: [METHOD ][HOST]/[PATH]

## Middleware logic and why.
- When these handlers are registered with mux it calls the func. This is why you need to make sure the logic for the incrementing of metrics is handled inside of the anonymous function return instead of the body of the function. If the increment happend in the body then it would happen once on startup then never again. This function will now instead return a new http.Handler that will get registered and then when it is called will increment the counter correctly. This new handler just icnrements the counter and calls servhttp against the original app handler for the normal website.