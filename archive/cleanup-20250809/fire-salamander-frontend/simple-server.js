const express = require('express');
const next = require('next');

const dev = process.env.NODE_ENV !== 'production';
const app = next({ dev });
const handle = app.getRequestHandler();
const port = process.env.PORT || 8084;

app.prepare().then(() => {
  const server = express();

  // Handle all requests
  server.all('*', (req, res) => {
    return handle(req, res);
  });

  server.listen(port, '0.0.0.0', (err) => {
    if (err) throw err;
    console.log(`🚀 Fire Salamander ready on http://0.0.0.0:${port}`);
    console.log(`🌐 Also try: http://127.0.0.1:${port}`);
    console.log(`🖥️  Or: http://localhost:${port}`);
  });
});