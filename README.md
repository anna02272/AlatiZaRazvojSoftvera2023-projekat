# Centralized Service Configuration System
This project involves the implementation of a centralized service configuration system comprising two main components:

1. **Web Service:** Accepts user requests and processes them.
2. **Database:** Stores the system's state.

Two auxiliary components maintain the system:

- **Log and Trace Management Component:** Responsible for storing and reviewing logs and traces.
- **Metrics Management Component:** Handles the storage and review of metrics.

##Technologies: 
- Programming language Golang
- Docker
- Swagger
- Consul NoSQL database
- Jaeger

### Web Service Component

The web service is implemented using the Go programming language (Golang) and provides the following operations:

- **Add Configuration:** Accepts JSON data for adding configuration to the system.
- **Add Configuration Group:** Accepts JSON data for adding a configuration group, which can contain one or more configurations.
- **View Configuration:** Retrieves configuration by identifier.
- **View Configuration Group:** Retrieves a group by identifier.
- **Delete Configuration:** Deletes a configuration by identifier.
- **Delete Configuration Group:** Deletes a group by identifier.
- **Expand Configuration Group:** Adds new configuration within a configuration group.
- **Advanced Operations using Labels:** Supports advanced operations on a configuration group using labels.

#### Label System

- Each configuration within a group should have a set of labels used for filtering and searching.
- Multiple configurations within a group can have the same set of labels.
- Labels are textual key-value pairs separated by a semicolon (e.g., l1:v1;l2:v2, ...).
- When querying configurations within a group using labels, all query labels must match those assigned to configurations.
- Supports deletion using the label system; the same rules as for searching apply.

#### Additional Requirements

- Supports immutability; configuration can only be replaced entirely.
- Idempotent requests are supported.
- Uses UUIDs as unique identifiers.
- Enables versioning, allowing configurations to be saved in different versions.
- When querying for configuration or groups, clients must specify the desired version.

### Database

- Configuration data is stored in the Consul NoSQL database.
- Configuration groups are stored in the Consul NoSQL database.

### Additional Requirements

- The service and database are containerized using Docker with a multi-stage build.
- The database is also containerized using Docker.
- Supports tracing within the service.
- Supports metric collection within the service.
- All elements run within a Docker Compose.
- Triggers the CI system (Git Actions) when changes are merged into the master (main) branch.
- Uses Git for version control, adhering to GitFlow principles.
- Stores information about request idempotence in the Consul NoSQL database.
