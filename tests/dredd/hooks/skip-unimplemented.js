module.exports = function registerHooks(hooks) {
  hooks.beforeEach((transaction, done) => {
    const isHealthCheck =
      transaction.fullPath === '/healthz' &&
      transaction.request?.method?.toUpperCase() === 'GET';

    if (!isHealthCheck) {
      transaction.skip = true;
      transaction.skipReason =
        'Endpoint ainda não implementado no backend (placeholder).';
    }

    done();
  });
};
