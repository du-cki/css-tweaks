# CSS Tweaks

A stateless, crowdsourced CSS patch service. Users can select multiple CSS snippets and receive a single URL that combines them.

There is no database. The app uses deterministic **[CRC32](https://en.wikipedia.org/wiki/Cyclic_redundancy_check)** hashing to generate stable IDs based on filenames.

## Adding New Snippets

To contribute, submit a Pull Request adding a `.css` file to the `snippets/` directory. You can also use subfolders to categorize them (e.g., `snippets/cosmetics/my-fix.css`).

## Development (Local)

You can run the server directly with Go for quick iteration.

```bash
go run .
```

The site will be available at http://localhost:8080/.

Note: When running this way, the repository metadata will default to "dev" or "unknown". To see specific version metadata locally, you must manually provide the variables or use the production script below.

### Production

For deployment or to run the application with full metadata and containerization, use the provided script. It injects the correct version info into the build.

```bash
./run.sh
```
