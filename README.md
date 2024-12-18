# BE EVENT BOOKING
This is Embreo technical assessment.

## How to run

### Required

- Redis

### Conf

To set up your environment, create a `.env` file by copying `.env.example` and customizing the configurations accordingly.


### Run
```
$ go run main.go 
```
Now server is running at port 8080.

## Run on docker

```bash
# docker
$ docker compose up --build
```
you can run on [http://localhost:8080](http://localhost:8080)

### User can be used
currently there are 2 HR and 2 vendor users. below are credential can be used:
1. username: `HR1` password: `password`
2. username: `HR2` password: `password`
3. username: `Vendor1` password: `password`
4. username: `Vendor2` password: `password`

## Swagger

Access swagger [http://localhost:8080/swagger](http://localhost:8080/swagger)

## ERD

please read the ERD `ERD.pdf`