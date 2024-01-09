# Forest App CLI (Unofficial)

This command line interface (CLI) tool is designed to interact with the Forest App API. It provides functionality for listing, adding, updating, and managing trees and tags associated with a user's account. 

## Features

- **Login**: Authenticate with the Forest App API.
- **List Trees**: Retrieve and display all trees associated with the user's account.
- **Add Tree**: Add a new tree with customizable start time, end time, tag, tree type, and notes.
- **Update Tree**: Update an existing tree's details.
- **List Tags**: Retrieve and display all user-defined tags.
- **List Unlocked Trees**: Display all the user's unlocked tree types.

## Prerequisites

- Go (Golang) environment set up on your machine.
- Access to the [Forest App](https://play.google.com/store/apps/details?id=cc.forestapp).

## Installation

1. Clone the repository:
   ```bash
   git clone git@github.com:kierenfunk/forest-cli.git 
   ```
2. Navigate to the project directory:
   ```bash
   cd forest-cli
   ```
3. Build the project:
   ```bash
   make
   ```

## Usage

Before using the CLI commands, make sure to set your Forest App credentials as environment variables or pass them as arguments.

### Setting Up Environment Variables

```bash
export FOREST_USERNAME="your_email@example.com"
export FOREST_PASSWORD="your_password"
```

### Available Commands

#### 1. List Trees

List all trees associated with the user's account.

```bash
forest list
```

#### 2. Add Tree

Add a new tree to the user's account.

```bash
forest add --start-time "2024-01-01T00:00:00+00:00" --end-time "2024-01-01T01:00:00+00:00" --tag "Work" --tree "Cedar" --note "Focus session"
```

Alternatively you can set the tree as chosen randomly from your set of unlocked trees.

```bash
forest add --start-time "2024-01-01T00:00:00+00:00" --end-time "2024-01-01T01:00:00+00:00" --tag "Work" --note "Focus session" --random
```

#### 3. Update Tree

Update an existing tree's details.

```bash
forest update --plant-id 123 --start-time "2024-01-02T00:00:00+00:00" --end-time "2024-01-02T01:00:00+00:00" --tag "Study" --note "Study session" 
```

#### 4. List Tags

List all user-defined tags.

```bash
forest tags
```

#### 5. List Unlocked Trees

List all unlocked tree types for the user.

```bash
forest trees
```

## Contributing

Contributions are welcome! If you have suggestions for improvements or encounter any issues, please feel free to open an issue or submit a pull request.

---
