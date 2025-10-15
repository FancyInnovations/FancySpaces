# Permissions

## Roles

- **Everyone**: Any visitor, including unauthenticated users
- **User**: Authenticated users with a registered account
- **SOwner**: Space Owner, the user who created the space
- **SAdmin**: Space Admin, users with administrative privileges in a space
- **SMember**: Space Member, regular members of a space

## Actions

### Space Management

| Action              | Roles          |
|---------------------|----------------|
| Get all spaces      | Everyone       |
| Get one space       | Everyone       |
| Create space        | User           |
| Update space        | SOwner, SAdmin |
| Delete space        | SOwner         |
| Change space status | SOwner, SAdmin |

### Member Management (within a space)

| Action              | Roles                   |
|---------------------|-------------------------|
| Get all member      | Everyone                |
| Get one member      | Everyone                |
| Add member          | SOwner                  |
| Update member       | SOwner                  |
| Remove member       | SOwner, Member themself |

### Issue Management (within a space)

| Action         | Roles                   |
|----------------|-------------------------|
| Get all issues | Everyone                |
| Get one issue  | Everyone                |
| Create issue   | SMember                 |
| Update issue   | SOwner, SAdmin, SMember |
| Delete issue   | SOwner, SAdmin, SMember |