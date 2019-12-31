# RedisSlowlogsParser for Redis3.x
Parser redis slow logs, you can easily to analyze slow query logs.

## What is redis slow log?
see: https://redis.io/commands/slowlog

## How do I collect logs?
`redis-cli slowlog get 100 > /user/path/slowlogs/redis-slow.log`

## How to use it?
1. `go build`
2. `RedisSlowlogsParser -slowlog /user/path/slowlogs/`

## Accepted parameters?
1. slowlog: slow log file or directory path.(required)
2. command: redis command.(optional)
3. duration(ms): slower than this value, will be display.(optional)


## Interfaces
```
// You can use this interface to customize your development.
slowlogsparser.ParserLogs(filePaths []string, durationThreshold float64, redisCommand string)
```

## Advanced
use linux crontab tool and this tool for automatic timing analysis.

