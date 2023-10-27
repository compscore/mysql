# MYSQL

## Check Format

```yaml
- name:
  release:
    org: compscore
    repo: mysql
    tag: latest
  credentials:
    username:
    password:
  target:
  expectedOutput:
  weight:
  options:
    database:
    table:
    field:
    connect:
    table_exists:
    row_exists:
    match:
    regex_match:
    substring_match:
```

## Parameters

|     parameter     |            path            |   type   | default  | required | description                                               |
| :---------------: | :------------------------: | :------: | :------: | :------: | :-------------------------------------------------------- |
|      `name`       |          `.name`           | `string` |   `""`   |  `true`  | `name of check (must be unique)`                          |
|       `org`       |       `.release.org`       | `string` |   `""`   |  `true`  | `organization that check repository belongs to`           |
|      `repo`       |      `.release.repo`       | `string` |   `""`   |  `true`  | `repository of the check`                                 |
|       `tag`       |       `.release.tag`       | `string` | `latest` | `false`  | `tagged version of check`                                 |
|    `username`     |  `.credentials.username`   | `string` |   `""`   | `false`  | `username of mysql user`                                  |
|    `password`     |  `.credentials.password`   | `string` |   `""`   | `false`  | `default password of mysql user`                          |
|     `target`      |         `.target`          | `string` |   `""`   |  `true`  | `mysql server network location`                           |
| `expectedOutput`  |     `.expectedOutput`      | `string` |   `""`   | `false`  | `expected output of queried row`                          |
|     `weight`      |         `.weight`          |  `int`   |   `0`    |  `true`  | `amount of points a successful check is worth`            |
|    `database`     |    `.options.database`     | `string` |   `""`   |  `true`  | `database to use for check`                               |
|      `table`      |      `.options.table`      | `string` |   `""`   | `false`  | `table to use for query based checks`                     |
|      `field`      |      `.options.field`      | `string` |   `""`   | `false`  | `field to use for query based checks`                     |
|     `connect`     |     `.options.connect`     |  `bool`  | `false`  | `false`  | `check to connect to mysql on given `                     |
|  `table_exists`   |  `.options.table_exists`   |  `bool`  | `false`  | `false`  | `check if table exists in specified database`             |
|   `row_exists`    |   `.options.row_exists`    |  `bool`  | `false`  | `false`  | `check if any row exists in specified table`              |
|      `match`      |      `.options.match`      |  `bool`  | `false`  | `false`  | `check if any row has exact match of specified field`     |
|   `regex_match`   |   `.options.regex_match`   |  `bool`  | `false`  | `false`  | `check if any row has regex match of specified field`     |
| `substring_match` | `.options.substring_match` |  `bool`  | `false`  | `false`  | `check if any row has substring match of specified field` |

## Examples

```yaml
- name: host_1-mysql
  release:
    org: compscore
    repo: mysql
    tag: latest
  credentials:
    username: root
    password: changeme
  target: 10.{{ .Team }}.1.1:3306
  expectedOutput: john_doe
  weight: 1
  options:
    database: prod
    table: users
    field: name
    match:
```

```yaml
- name: host_1-mysql
  release:
    org: compscore
    repo: mysql
    tag: latest
  credentials:
    username: root
    password: changeme
  target: 10.{{ .Team }}.1.1:3306
  weight: 1
  options:
    database: prod
    table: users
    field: name
    connect:
    table_exists:
    row_exists:
```
