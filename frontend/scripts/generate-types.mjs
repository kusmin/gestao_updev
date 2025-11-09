import path from 'node:path';
import os from 'node:os';
import fs from 'node:fs';
import { execSync } from 'node:child_process';

const cwd = process.cwd();
const defaultSpec = path.resolve(cwd, '..', 'docs', 'api.yaml');
const specPath = process.env.SPEC_PATH
  ? path.resolve(cwd, process.env.SPEC_PATH)
  : defaultSpec;
const target = path.resolve(cwd, 'src', 'types', 'api.d.ts');
let specForTypes = specPath;

const specContent = fs.readFileSync(specPath, 'utf8');
const isSwagger2 = /^\s*swagger:\s*["']?2/.test(specContent);

if (isSwagger2) {
  const tempDir = fs.mkdtempSync(path.join(os.tmpdir(), 'gestao-spec-'));
  const converted = path.join(tempDir, 'openapi.yaml');
  execSync(
    `npx swagger2openapi "${specPath}" --outfile "${converted}"`,
    { stdio: 'inherit', cwd }
  );
  specForTypes = converted;
}

execSync(
  `npx openapi-typescript "${specForTypes}" -o "${target}"`,
  { stdio: 'inherit', cwd }
);

console.log(`Tipos gerados a partir de ${specPath} → ${target}`);
