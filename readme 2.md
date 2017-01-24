## Clip ? ##
You can call it WEB summary.
Clip digs informations from a given URL and then deliver the summary of a WEB (given URL).

## How To Use ##
1. Golang must installed on your machine
2. Clone This Repo
3. Do `go get .` on your cloned repo
4. setup .env *see instruction below
5. run program using `./run.sh serve`
6. make POST request with URL param. see this example :
   ```
   curl -X POST -F 'url=https://youtu.be/50efl4S8VQc' http://localhost:3001
   ```
7. see the results

## How to config environment ##
- For .env example please open file named .env.example
- Edit as you need
- save as ".env" (*program will only read this name)

## Doc ##
Please open docs.md

## License ##

This package is licensed under MIT license. See LICENSE for details.
