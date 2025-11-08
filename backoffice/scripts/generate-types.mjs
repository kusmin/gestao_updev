import path from 'node:path';
import { execSync } from 'node:child_process';

const cwd = process.cwd();
const defaultSpec = path.resolve(cwd, '..', 'docs', 'api.yaml');
const specPath = process.env.SPEC_PATH
  ? path.resolve(cwd, process.env.SPEC_PATH)
  : defaultSpec;
const target = path.resolve(cwd, 'src', 'types', 'api.d.ts');

execSync(
  `npx openapi-typescript "${specPath}" -o "${target}"`,
  { stdio: 'inherit', cwd }
);

console.log(`Tipos gerados a partir de ${specPath} → ${target}`);
