# simple api  _(golang)_

[![N|Solid](https://cldup.com/dTxpPi9lDf.thumb.png)](https://nodesource.com/products/nsolid)

[![Build Status](https://travis-ci.org/joemccann/dillinger.svg?branch=master)](https://travis-ci.org/joemccann/dillinger)

#### Testing Machine:

    - MacOs 12.6.7
    - Docker version 24.0.2, build cb74dfc
    - golang version go1.21.3 darwin/amd64
    - redis version redis:7.0.8-alpine
    - postgres postgres:13-alpine


## Documentation
- saat development `postgres` dan `redis` di jalankan di project terpisah, sehingga memerlukan `network bridge` bawaan
- migrasi, dan build docker ada di makefile, sehingga perlu penyesuaian host `postgres` unutk menjalankan migrasi
- `makefile migrate_up` untuk menjalakan migrasi
- `makefile migrate_down` untuk roll out db
- `makefile build_up` untuk membuat dan menjalankan docker
- `makefile migrate_down` untuk menghapus container docker
-  seeder di lakukan manual terdapat pada dir `/db/seeder/manual_seeder.sql` jalankan query tersebut di db
-  unit testing berada di directory `test`
-  collection postman berada di directory `postman`



## Installation

#### Instalation project Golang:

    1. clone repository   
    2. masuk ke folder hasil clone
    3. setup .env file dan makefile
    4. go mod download
    5. go run main.go
    6. jalankan `makefile migrate_up`
    7. jalankan `makefile build_up`
    8. enjoy with golang API 



## Testing

#### Testing:
    - testing berada di directory `test`
#### Note

- jika ada kendala di invite project github / postman mohon hubungi afitra - 085230010042

### http://apitoong.com
 
 
 
 
