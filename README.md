# Tiny url

## Todos:
- [X] Create simple API server with GET /ping (will user for heath check)     
- [X] Fix usage of context    
- [X] MongoDB:    
> - [X] URL collection:
>> * User ID
>> * URL hash
>> * Real URL
>> * Experation
>> * Last time visited
>> * Number of visits

> - [X] Uers collection: 
>> * Email
>> * Password
>> * URLs data

- [X] Users:
> - [X] Add new user with PWD and JWT token   
> - [X] Login     
> - [X] Get statistics

- [X] URLs:
> - [X] Add new URL   
> - [X] GET /{hashed} redirect to actual
> - [X] Update statistics
> - [X] Timeout URLs