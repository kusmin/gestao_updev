// scripts/mcp-server.js
const http = require('http');
const { exec } = require('child_process');

const PORT = process.env.PORT || 8081;

const tools = [
  {
    name: "hello-world",
    description: "A simple tool that says hello.",
    parameters: [],
    execute: async () => {
      return {
        type: "text",
        content: "Hello from gestao_updev MCP Server!"
      };
    }
  },
  {
    name: "generate-api-types",
    description: "Generates TypeScript API types for the frontend from the backend's OpenAPI specification.",
    parameters: [],
    execute: async () => {
      return new Promise((resolve, reject) => {
        exec('npm run generate:api-types', { cwd: 'frontend' }, (error, stdout, stderr) => {
          if (error) {
            console.error(`exec error: ${error}`);
            return reject({ type: "text", content: `Error generating API types: ${stderr}` });
          }
          if (stderr) {
            console.warn(`stderr: ${stderr}`);
          }
          resolve({ type: "text", content: `API types generated successfully:\n${stdout}` });
        });
      });
    }
  },
  {
    name: "run-all-tests",
    description: "Executes all backend and frontend tests.",
    parameters: [],
    execute: async () => {
      return new Promise(async (resolve, reject) => {
        let output = "";
        try {
          output += "Running backend tests...\n";
          const backendOutput = await new Promise((res, rej) => {
            exec('docker compose -f docker-compose.test.yml up --build --abort-on-container-exit', (error, stdout, stderr) => {
              if (error) {
                console.error(`Backend test exec error: ${error}`);
                return rej(`Error running backend tests: ${stderr}`);
              }
              res(stdout);
            });
          });
          output += backendOutput + "\n";
          output += "Backend tests completed.\n\n";

          output += "Running frontend tests...\n";
          const frontendOutput = await new Promise((res, rej) => {
            exec('npm run test', { cwd: 'frontend' }, (error, stdout, stderr) => {
              if (error) {
                console.error(`Frontend test exec error: ${error}`);
                return rej(`Error running frontend tests: ${stderr}`);
              }
              res(stdout);
            });
          });
          output += frontendOutput + "\n";
          output += "Frontend tests completed.\n";

          resolve({ type: "text", content: output });
        } catch (error) {
          reject({ type: "text", content: `Error during tests: ${error}` });
        }
      });
    }
  },
  {
    name: "deploy-backend-staging",
    description: "Deploys the Go backend to the staging environment.",
    parameters: [],
    execute: async () => {
      return new Promise(async (resolve, reject) => {
        let output = "Starting backend deployment to staging...\n";
        try {
          // TODO: Implement actual deployment logic here.
          // This might involve:
          // 1. Building the Go application: exec('go build ./cmd/api', { cwd: 'backend' })
          // 2. Building the Docker image: exec('docker build -t gestao_updev_api .', { cwd: 'backend' })
          // 3. Pushing to a container registry
          // 4. Updating a Kubernetes deployment or similar
          output += "Simulating deployment...\n";
          await new Promise(res => setTimeout(res, 3000)); // Simulate work
          output += "Backend deployed to staging successfully (simulated).\n";
          resolve({ type: "text", content: output });
        } catch (error) {
          console.error(`Deployment error: ${error}`);
          reject({ type: "text", content: `Error deploying backend to staging: ${error}` });
        }
      });
    }
  }
];

const server = http.createServer(async (req, res) => {
  if (req.url === '/tools' && req.method === 'GET') {
    res.writeHead(200, { 'Content-Type': 'application/json' });
    res.end(JSON.stringify(tools.map(tool => ({
      name: tool.name,
      description: tool.description,
      parameters: tool.parameters
    }))));
  } else if (req.url.startsWith('/tools/') && req.method === 'POST') {
    const toolName = req.url.split('/')[2];
    const tool = tools.find(t => t.name === toolName);

    if (tool) {
      let body = '';
      req.on('data', chunk => {
        body += chunk.toString();
      });
      req.on('end', async () => {
        const params = body ? JSON.parse(body) : {};
        try {
          const result = await tool.execute(params);
          res.writeHead(200, { 'Content-Type': 'application/json' });
          res.end(JSON.stringify(result));
        } catch (error) {
          res.writeHead(500, { 'Content-Type': 'application/json' });
          res.end(JSON.stringify({ error: error.message }));
        }
      });
    } else {
      res.writeHead(404, { 'Content-Type': 'application/json' });
      res.end(JSON.stringify({ error: "Tool not found" }));
    }
  } else {
    res.writeHead(404, { 'Content-Type': 'application/json' });
    res.end(JSON.stringify({ error: "Not Found" }));
  }
});

server.listen(PORT, () => {
  console.log(`MCP Server running on port ${PORT}`);
});
