const hooks = require('hooks');
const crypto = require('crypto');

const SUPPORTED = new Set([
  'GET /healthz',
  'POST /auth/signup',
  'POST /auth/login',
  'POST /auth/refresh',
  'GET /companies/me',
  'GET /users',
  'POST /users',
  'PATCH /users/{id}',
  'DELETE /users/{id}',
  'GET /clients',
  'POST /clients',
  'GET /clients/{id}',
  'PUT /clients/{id}',
  'DELETE /clients/{id}',
]);

const PUBLIC_ENDPOINTS = new Set([
  'GET /healthz',
  'POST /auth/signup',
  'POST /auth/login',
  'POST /auth/refresh',
]);

const ctx = {
  email: '',
  password: 'Senha@123',
  tenantId: '',
  accessToken: '',
  refreshToken: '',
  createdUserId: '',
  createdClientId: '',
  clientEmail: '',
};

const MEMBER_PASSWORD = 'Senha@456';

function uniqueSuffix() {
  return crypto.randomUUID().split('-')[0];
}

function normalizePath(value) {
  if (!value) {
    return '';
  }
  let path = value;
  if (path.startsWith('http')) {
    try {
      path = new URL(path).pathname;
    } catch {
      // leave as is
    }
  } else if (!path.startsWith('/')) {
    path = `/${path}`;
  }
  path = path.split('?')[0];
  return path.replace(/^\/v[0-9]+(?=\/)/, '') || '/';
}

function keyOf(transaction) {
  const method = transaction.request.method.toUpperCase();
  const rawPath =
    (transaction.origin && transaction.origin.path) ||
    transaction.fullPath ||
    transaction.request.uri ||
    '';
  return `${method} ${normalizePath(rawPath)}`;
}

function setJSONBody(transaction, payload) {
  transaction.request = transaction.request || {};
  transaction.request.body = JSON.stringify(payload);
  transaction.request.headers = transaction.request.headers || {};
  transaction.request.headers['Content-Type'] = 'application/json';
}

function setAuthHeaders(transaction) {
  transaction.request = transaction.request || {};
  transaction.request.headers = transaction.request.headers || {};
  transaction.request.headers.Authorization = `Bearer ${ctx.accessToken}`;
  transaction.request.headers['X-Tenant-ID'] = ctx.tenantId;
}

function setResourceId(transaction, resource, id) {
  const matcher = new RegExp(`/${resource}/[^/?]+`);
  const replacement = `/${resource}/${id}`;

  if (transaction.fullPath) {
    transaction.fullPath = transaction.fullPath.replace(matcher, replacement);
  }

  transaction.request = transaction.request || {};
  if (transaction.request.uri) {
    transaction.request.uri = transaction.request.uri.replace(matcher, replacement);
  }

  transaction.request.params = transaction.request.params || {};
  transaction.request.params.id = id;
}

function ensureResource(transaction, resource, id, reason) {
  if (!id) {
    transaction.skip = true;
    transaction.skipReason = reason;
    return false;
  }
  setResourceId(transaction, resource, id);
  return true;
}

  hooks.beforeEach((transaction, done) => {
    const key = keyOf(transaction);
    if (!SUPPORTED.has(key)) {
      transaction.skip = true;
      transaction.skipReason = 'Endpoint ainda não coberto pelo fluxo básico de contrato.';
      return done();
    }

    if (transaction.expected && transaction.expected.headers && transaction.expected.headers['Content-Type'] === 'application/json') {
      transaction.expected.headers['Content-Type'] = 'application/json; charset=utf-8';
    }

    if (!PUBLIC_ENDPOINTS.has(key)) {
      if (!ctx.accessToken || !ctx.tenantId) {
        transaction.skip = true;
        transaction.skipReason = 'Tokens de autenticação indisponíveis para o endpoint solicitado.';
        return done();
      }
      setAuthHeaders(transaction);
    }

    switch (key) {
      case 'POST /auth/signup': {
        ctx.email = `qa-${uniqueSuffix()}@gestao.dev`;
        ctx.createdUserId = '';
        ctx.createdClientId = '';
        ctx.clientEmail = '';
        const document = crypto.randomUUID().replace(/-/g, '').slice(0, 14);
        setJSONBody(transaction, {
          company: {
            name: `QA Company ${uniqueSuffix()}`,
            document,
            phone: '+55 11 90000-0000',
          },
          user: {
            name: 'QA Admin',
            email: ctx.email,
            password: ctx.password,
            phone: '+55 11 91111-0000',
          },
        });
        if (transaction.expected) {
          delete transaction.expected.body;
        }
        break;
      }
      case 'POST /auth/login':
        if (!ctx.email) {
          transaction.skip = true;
          transaction.skipReason = 'Conta administrativa ainda não criada para login.';
          return done();
        }
        setJSONBody(transaction, {
          email: ctx.email,
          password: ctx.password,
        });
        break;
      case 'POST /auth/refresh':
        if (!ctx.refreshToken) {
          transaction.skip = true;
          transaction.skipReason = 'Refresh token ainda não disponível para renovação.';
          return done();
        }
        setJSONBody(transaction, {
          refresh_token: ctx.refreshToken,
        });
        break;
      case 'POST /users': {
        const email = `qa-user-${uniqueSuffix()}@gestao.dev`;
        setJSONBody(transaction, {
          name: 'QA Member',
          email,
          role: 'manager',
          phone: '+55 11 92222-0000',
          password: MEMBER_PASSWORD,
        });
        break;
      }
      case 'PATCH /users/{id}':
        if (!ensureResource(transaction, 'users', ctx.createdUserId, 'Usuário de teste indisponível para PATCH.')) {
          return done();
        }
        setJSONBody(transaction, {
          role: 'operator',
          active: true,
          phone: '+55 11 95555-0000',
        });
        break;
      case 'DELETE /users/{id}':
        if (!ensureResource(transaction, 'users', ctx.createdUserId, 'Usuário de teste indisponível para DELETE.')) {
          return done();
        }
        break;
      case 'POST /clients': {
        ctx.clientEmail = `qa-client-${uniqueSuffix()}@gestao.dev`;
        setJSONBody(transaction, {
          name: 'QA Client',
          phone: '+55 11 98888-0000',
          email: ctx.clientEmail,
          notes: 'Cliente criado nos testes de contrato.',
        });
        break;
      }
      case 'GET /clients/{id}':
        if (!ensureResource(transaction, 'clients', ctx.createdClientId, 'Cliente de teste indisponível para GET.')) {
          return done();
        }
        break;
      case 'PUT /clients/{id}':
        if (!ensureResource(transaction, 'clients', ctx.createdClientId, 'Cliente de teste indisponível para PUT.')) {
          return done();
        }
        setJSONBody(transaction, {
          id: ctx.createdClientId,
          name: 'QA Client Atualizado',
          phone: '+55 11 97777-0000',
          email: ctx.clientEmail,
          notes: 'Registro atualizado pelos testes de contrato.',
          tags: ['vip'],
        });
        break;
      case 'DELETE /clients/{id}':
        if (!ensureResource(transaction, 'clients', ctx.createdClientId, 'Cliente de teste indisponível para DELETE.')) {
          return done();
        }
        break;
      case 'GET /users':
      case 'GET /clients':
        // only auth headers needed
        break;
      case 'GET /companies/me':
        // auth headers already set outside switch
        break;
      default:
        break;
    }

    done();
  });

  hooks.afterEach((transaction, done) => {
    const key = keyOf(transaction);
    if (!SUPPORTED.has(key)) {
      return done();
    }

    let body;
    try {
    body = JSON.parse((transaction.real && transaction.real.body) || '{}');
    } catch {
      return done();
    }
    const data = body.data || {};

    if (key === 'POST /auth/signup') {
      ctx.tenantId = data.tenant_id;
      ctx.accessToken = data.access_token;
      ctx.refreshToken = data.refresh_token;
    } else if (key === 'POST /auth/login') {
      ctx.accessToken = data.access_token;
      ctx.refreshToken = data.refresh_token;
    } else if (key === 'POST /auth/refresh') {
      ctx.accessToken = data.access_token;
      ctx.refreshToken = data.refresh_token;
    }

    if (key === 'POST /users' && data.id) {
      ctx.createdUserId = data.id;
    } else if (key === 'POST /clients' && data.id) {
      ctx.createdClientId = data.id;
    }

    done();
  });
