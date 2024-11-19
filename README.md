# availability-monitor

Availability Monitor is a simple system that handles the availability state of sensors.

The sensors can typically be motion detectors in a call-box.

See https://github.com/EtincelleCoworking/availability-monitor/wiki

## Usage

```
go build .
export DB_URL=postgres://admin:admin@localhost:5432/db?sslmode=disable
export ADMIN_USERNAME=admin
# md5 hash of "abcd"
export ADMIN_PASSWORD_HASH='$2a$10$NjVOiKsE6u1lw2G5GnARN.nVTJKzmKfXQgCKjb4yjy9KASbWUzaB2'
./availability-monitor
```

- Simulate a sensor "alive"

```bash
curl -X POST http://localhost:7000/api/alive?mac=abcd
```

- Open http://localhost:7000/admin/sensors/pending
- Connect (admin/abcd)
- Set a location to the new sensor
- Open http://localhost:7000
- Simulate a sensor "motion"

```bash
curl -X POST http://localhost:7000/api/motion?mac=abcd
```
