# Admin Service

This is the Admin service

Generated with

```
micro new --namespace=go.micro --type=service admin
```

## Getting Started

- [Configuration](#configuration)
- [Dependencies](#dependencies)
- [Usage](#usage)

## Configuration

- FQDN: go.micro.service.admin
- Type: service
- Alias: admin

## Dependencies

Micro services depend on service discovery. The default is multicast DNS, a zeroconf system.

In the event you need a resilient multi-host setup we recommend etcd.

```
# install etcd
brew install etcd

# run etcd
etcd
```

## Usage

A Makefile is included for convenience

Build the binary

```
make build
```

Run the service
```
./admin-service
```

Build a docker image
```
make docker
```
INSERT INTO `user` (`phone`,`password`,`updated_at`,`created_at`,`version`) VALUES (?, ?, ?, ?, ?) []interface {}{13711400432, "$2a$10$b1H8LmV7KbYcXhbxvG0W6.cZ8/KNYI6ou.wG2rJktw7qR74t8eQM2", "2024-10-02 02:11:10", "2024-10-02 02:11:10", 1}

insert into role (id,name,menu_ids) VALUES(1,"admin",CONVERT("[]" USING BINARY));

update `user` set role_ids=CONVERT("[1]" USING BINARY) where id = 1;
