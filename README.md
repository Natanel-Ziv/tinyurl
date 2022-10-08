# Tiny url

## Todos:
- [X] Create simple API server with GET /ping (will user for heath check)     
- [X] Fix usage of context    
- [ ] MongoDB:    
> - [ ] URL collection:   
>> * User ID
>> * URL hash
>> * Real URL
>> * Is active
>> * Last time visited
>> * Number of visits

> - [X] Uers collection: 
>> * Email
>> * Password
>> * URLs data

- [ ] Users:
> - [X] Add new user with PWD and JWT token   
> - [X] Login     
> - [ ] Add new URL   
> - [ ] Get statistics

- [ ] URLs:
> - [ ] GET /{hashed} redirect to actual  
> - [ ] Update statistics     
> - [ ] Timeout URLs