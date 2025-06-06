# Chirpy

## Mux.Handle vs mux.HandleFunc
- Handle requires that the second arguement implements the http.Handler interface where handlefunc does not. Handlefunc will take the function provided and make it impelemnt the http.Handler interface. Handle is used for when you want to implement a custom router, middleware, or other custom logic. If you want to make a function into a http.Handler then you can use the http.HanlderFunc() function to do so. Ex: middlewareMetricsInc method on apiConfig.

## Middleware logic and why.
- When these handlers are registered with mux it calls the func. This is why you need to make sure the logic for the incrementing of metrics is handled inside of the anonymous function return instead of the body of the function. If the increment happend in the body then it would happen once on startup then never again. This function+ will now instead return a new http.Handler that will get registered and then when it is called will increment the counter correctly. This new handler just icnrements the counter and calls servhttp against the original app handler for the normal website.

## More on Patterns
A pattern is a string that specifies the set of URL paths that should be matched to handle HTTP requests. Go's ServeMux router uses these patterns to dispatch requests to the appropriate handler functions based on the URL path of the request. As we saw in the previous lesson, patterns help organize the handling of different routes efficiently.

As previously mentioned, patterns generally look like this: [METHOD ][HOST]/[PATH]. Note that all three parts are optional.

### Rules and Definitions
Fixed URL Paths
A pattern that exactly matches the URL path. For example, if you have a pattern /about, it will match the URL path /about and no other paths.

### Subtree Paths
If a pattern ends with a slash /, it matches all URL paths that have the same prefix. For example, a pattern /images/ matches /images/, /images/logo.png, and /images/css/style.css. As we saw with our /app/ path, this is useful for serving a directory of static files or for structuring your application into sub-sections.

### Longest Match Wins
If more than one pattern matches a request path, the longest match is chosen. This allows more specific handlers to override more general ones. For example, if you have patterns / (root) and /images/, and the request path is /images/logo.png, the /images/ handler will be used because it's the longest match.

### Host-Specific Patterns
We won't be using this but be aware that patterns can also start with a hostname (e.g., www.example.com/). This allows you to serve different content based on the Host header of the request. If both host-specific and non-host-specific patterns match, the host-specific pattern takes precedence.