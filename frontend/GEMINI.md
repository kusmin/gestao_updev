# GEMINI Project Context: gestao_updev/frontend

## Project Overview

This is the frontend for the `gestao_updev` project, a SaaS platform for local businesses. It's a single-page application built with React, TypeScript, and Vite. It communicates with a backend API to fetch and display data. The project uses `@tanstack/react-query` for data fetching and state management, and `openapi-typescript` to generate TypeScript types from an OpenAPI specification, ensuring a strong contract with the backend.

## Building and Running

To get the frontend up and running, you'll need to have Node.js and npm installed.

1.  **Install dependencies:**

    ```bash
    npm install
    ```

2.  **Generate API types:**

    The frontend relies on TypeScript types generated from the backend's OpenAPI specification. Before running the development server, you need to generate these types.

    ```bash
    npm run generate:api-types
    ```

    This command runs `make -C .. api-types`, so make sure you have `make` installed and that the backend project is in the parent directory.

3.  **Run the development server:**

    ```bash
    npm run dev
    ```

    This will start the Vite development server, and you can view the application at `http://localhost:5173`.

4.  **Build for production:**

    ```bash
    npm run build
    ```

    This will create a `dist` directory with the production-ready assets.

## Development Conventions

*   **Language:** The codebase is written in TypeScript and JSX.
*   **Styling:** CSS files are used for styling.
*   **Linting:** The project uses ESLint for code quality and consistency. You can run the linter with:

    ```bash
    npm run lint
    ```
*   **API Interaction:** All communication with the backend API should be done through the `apiClient.ts` module or similar modules. The API types are available in `src/types/api.d.ts`.
