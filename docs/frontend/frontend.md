## Getting Started

First, run the development server:

```bash
pnpm install
cd web
pnpm run dev
```

Open [http://localhost:3000](http://localhost:3000) with your browser to see the result.

## Generating the Frontend API Client

If you need to make changes to the api spec (creating a new response type or creating a new endpoint). This will need to be done in the `openapi.yml` file first.

To regenerate the frontend API Client you can use the following command:

```bash
make generate
```

## How the API Client works

We're using [Orval](https://orval.dev/guides/swr) to handle the client generation. It produces:

- **Types** (for the models, response/request objects)
- **SWR Hooks** (for GET requests)
- **Axios Requests** (for POST/OPTIONS/PUT/DELETE requests)
