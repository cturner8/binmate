# Database

binmate uses an embedded [SQLite](https://www.sqlite.org/) database to track
installations, versions, downloads, and activity logs. The database is created and
migrated automatically the first time binmate runs.

## Location

```
~/.local/share/binmate/user.db
```

## Tables

| Table           | Description                                                |
| --------------- | ---------------------------------------------------------- |
| `binaries`      | The binary definitions known to binmate.                   |
| `installations` | Records of each installed version on disk.                 |
| `versions`      | The active version reference for each binary.              |
| `downloads`     | A cache of downloaded release assets.                      |
| `logs`          | Activity logs for operations performed by binmate.         |

## How state is tracked

binmate keeps the filesystem and the database in sync:

- When a version is **installed**, the files are extracted to a versioned directory and
  a row is added to `installations`.
- The **active version** for each binary is recorded in the `versions` table and
  reflected on disk via a symlink.
- **Downloads** are cached and tracked in the `downloads` table to avoid re-fetching
  assets unnecessarily.

::: tip Migrations
Schema migrations are versioned and run automatically when binmate connects to the
database, so upgrading binmate keeps your existing data intact.
:::
