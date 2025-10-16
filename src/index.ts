import express from 'express';
import swaggerUi from 'swagger-ui-express';
// @ts-ignore
import { RegisterRoutes } from './routes';

const app = express();
const PORT = process.env.PORT || 8000;

app.use(express.json());

RegisterRoutes(app);

try {
  const swaggerDocument = require('../docs/swagger.json');
  app.use('/api-docs', swaggerUi.serve, swaggerUi.setup(swaggerDocument));
} catch (error) {
  console.error('Error: Unable to load swagger.json. Please run the docs build command.');
}

app.listen(PORT, () => {
  console.log(`Server is running on http://localhost:${PORT}`);
  console.log(`Swagger docs available at http://localhost:${PORT}/api-docs`);
});
