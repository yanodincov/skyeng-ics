# Skyeng ICS Maker

This is a server that generates an ICS calendar file link for Skyeng student schedule.

## Usage

### Prerequisites

- [Docker](https://docs.docker.com/engine/install/)

### Configuration

Make `.env` file:

```bash
cp .env.example .env
```

Fill in the required values to the `.env` file.

### Running the server

Build the server:

```bash
docker build -t skyeng-ics-maker .
````

Run the server:

```bash
docker run -p 8080:8080 --env-file ./.env skyeng-ics-maker
```

Run the server and subscribe your calendar to the `$(SERVER_URL)/$(ROUTE_SUFFIX)/calendar.ics` link.

## Development

### Prerequisites

- go 1.23.2+
- [Taskfile](https://taskfile.dev/installation/)

### Running the server

Install dependencies:

```bash
task deps
```

Run the server:

```bash
task run
```

Show available tasks:

```bash
task -l
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.


