# gorm-gen

`gorm-gen` is a tool that can help you generate Golang code based on your existing database schema, and it utilizes the [GORM](https://gorm.io/) library for database operations.

## Installation

You can install `gorm-gen` using `go get`:

```bash
go install github.com/ydssx/gorm-gen
```

## Usage

`gorm-gen` requires a configuration file to specify the database connection information and the code generation options. The configuration file should be in the YAML format, and it must include the following fields:

```yaml
database:
  host: localhost
  port: 3306
  username: root
  password: 965213
  name: gva

output: ./model

tables:
  - sys_users
  - terminals
```

Here is the meaning of each field:

- `database`: the database connection information.
  - `host`: the database host.
  - `port`: the database port.
  - `name`: the database name.
  - `user`: the database user.
  - `password`: the database password.
- `output`: the output directory and package name.
  - `path`: the output directory.
  - `package`: the package name.
- `tables`: the list of models to be generated.
  - `table`: the name of the database table.
  

After creating the configuration file, you can run `gorm-gen` using the following command:

```bash
gorm-gen -c config.yaml
```

This will generate Golang code for the specified models in the output directory.

## Examples

Here is an example configuration file for generating Golang code for a `users` table with `id`, `name`, and `email` fields:

```yaml
database:
  host: "localhost"
  port: 3306
  name: "test"
  user: "root"
  password: "password"
output:
  path: "./models"
  package: "models"
models:
  - name: "user"
    table: "users"
    fields:
      - name: "id"
        type: "uint"
        tag: "primary_key"
        comment: "user id"
      - name: "name"
        type: "string"
        tag: "not null"
        comment: "user name"
      - name: "email"
        type: "string"
        tag: "unique_index"
        comment: "user email"
```

## License

`gorm-gen` is licensed under the MIT License. See the [LICENSE](LICENSE) file for more details.
