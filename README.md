## Explore GitHub (Backend API Server)

This repository contains the source code for the backend of the "Explore GitHub" project. This project is a web application that allows viewers to explore the social graph induced by GitHub user's contributions to repositories. This project was developed to streamline the process for finding new, interesting projects based on a viewer's interest in a particular repository. Each repository is represented as a node in the graph, and the edges between nodes represent a user's contribution to multiple repositories. By double-clicking on a repository node, the viewer can pull and display the projects `README.md` file, which provides a brief overview of the project. The viewer can also click on a repository node to view the repository's GitHub page.

### Getting Started

This project was developed from a protobuf definition of the API, which was used to generate the server and client code. The server was developed using the `grpc` framework, and the client was developed using the `grpc-web` framework. The server was developed in `Go`, and the client was developed in `TypeScript`. The server was containerized using `Docker`, and the client was containerized using `Docker` and `Nginx`.

To get started, clone the repository and navigate to the root directory of the project. Create the `egh-api.env` settings file.

```bash
git clone --recursive
cd egh-api
touch egh-api.env
nano egh-api.env
```

Once the `egh-api.env` file is created, add a GitHub [personal access token](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token) to the file.

```bash
GITHUB_TOKEN=ghp_...  # access token required to leverage the GraphQL API (does not need any permissions)
```

### Making changes

The principal method for making changes is by adding new API objects in the `proto` directory. Once the changes are made, the server and client code can be regenerated using the `buf` tool. The server and client code can be regenerated using the following commands:

```bash
buf generate
```

### Running the server

To run the server, use the following command:

```bash
go run cmd/server/main.go
```

### Contributing changes

Changes, additions, and improvements are welcome. Please feel free to open an issue or submit a pull request.
