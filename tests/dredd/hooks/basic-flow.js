const crypto = require('crypto');

const SUPPORTED = new Set([
  'GET /v1/healthz',
  'POST /v1/auth/signup',
  'POST /v1/auth/login',
  'POST /v1/auth/refresh',
  'GET /v1/companies/me',
]);

const ctx = {
  email: '',
  password: 'Senha@123',
  tenantId: '',
  accessToken: '',
  refreshToken: '',
};

function keyOf(transaction) {
  return `${transaction.request.method.toUpperCase()} ${transaction.fullPath}`;
}

function setJSONBody(transaction, payload) {
  transaction.request.body = JSON.stringify(payload);
  transaction.request.headers = transaction.request.headers || {};
  transaction.request.headers['Content-Type'] = 'application/json';
}

function setAuthHeaders(transaction) {
  transaction.request.headers = transaction.request.headers || {};
  transaction.request.headers.Authorization = `Bearer ${ctx.accessToken}`;
  transaction.request.headers['X-Tenant-ID'] = ctx.tenantId;
}

module.exports = function registerHooks(hooks) {
  hooks.beforeEach((transaction, done) => {
    const key = keyOf(transaction);
    if (!SUPPORTED.has(key)) {
      transaction.skip = true;
      transaction.skipReason = 'Endpoint ainda não coberto pelo fluxo básico de contrato.';
      return done();
    }

    if (key === 'POST /auth/signup') {
      ctx.email = `qa-${Date.now()}@gestao.dev`;
      const document = `99${Date.now()}`.slice(0, 14);
      setJSONBody(transaction, {
        company: {
          name: `QA Company ${Date.now()}`,
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
    } else if (key === 'POST /auth/login') {
      setJSONBody(transaction, {
        email: ctx.email,
        password: ctx.password,
      });
    } else if (key === 'POST /auth/refresh') {
      setJSONBody(transaction, {
        refresh_token: ctx.refreshToken,
      });
    } else if (key === 'GET /companies/me') {
      setAuthHeaders(transaction);
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
      body = JSON.parse(transaction.real.body || '{}');
    } catch (err) {
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

    done();
  });
};
