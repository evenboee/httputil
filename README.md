# httputil

## Request Binding

bind.All is used in handler.Func to bind a request to given parameters. 

Binding is based on struct fields e.g. url query will be bound for fields of a struct with name "Query". 

The supported types are: 

|      Type      | Field Name | Tag Name |
|----------------|------------|----------|
| path variables | Path       | path     |
| url query      | Query      | query    |
| post form      | PostForm   | form     |
| headers        | Header     | header   |
| body           | Body       | *varies  |


Body will check if pointer of type implements BodyUnmarshaler. 
If it does not, it will try to parse the body according to the `Content-Type` header, defaulting to JSON. 
Supported content types are:
- JSON: application/json
- Post form: application/x-www-form-urlencoded
