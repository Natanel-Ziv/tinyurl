# Tiny url

## Todos:
- [X] Create simple API server with GET /ping (will user for heath check)     
- [X] Fix usage of context    
- [ ] MongoDB:    
> - [ ] URL collection:   
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

- [ ] Users:
> - [X] Add new user with PWD and JWT token   
> - [X] Login     
> - [ ] Get statistics

- [ ] URLs:
> - [X] Add new URL   
> - [ ] GET /{hashed} redirect to actual  
> - [ ] Update statistics     
> - [ ] Timeout URLs