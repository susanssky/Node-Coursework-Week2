import { Express, Request, Response } from 'express'
import swaggerJsdoc from 'swagger-jsdoc'
import swaggerUi from 'swagger-ui-express'
import { version } from '../package.json'

const options: swaggerJsdoc.Options = {
  definition: {
    openapi: '3.0.2',
    info: {
      title: 'CYF Chat application of REST API Docs',
      version,
    },
    components: {
      schemas: {
        Message: {
          type: 'object',
          properties: {
            id: { type: 'string' },
            from: { type: 'string' },
            text: { type: 'string' },
            timeSent: { type: 'string' },
          },
          required: ['id', 'from', 'text'],
        },
        MessageInput: {
          type: 'object',
          properties: {
            from: { type: 'string' },
            text: { type: 'string' },
          },
          required: ['from', 'text'],
        },
      },
      securitySchemes: {
        bearerAuth: {
          type: 'http',
          scheme: 'bearer',
          bearerFormat: 'JWT',
        },
      },
    },
    security: [
      {
        bearerAuth: [],
      },
    ],
  },
  apis: ['./routes/messages.ts'],
}

const swaggerSpec = swaggerJsdoc(options)

function swaggerDocs(app: Express) {
  // Swagger page
  app.use('/api/messages/docs', swaggerUi.serve, swaggerUi.setup(swaggerSpec))

  // Docs in JSON format
  app.get('/api/messages/docs.json', (req: Request, res: Response) => {
    res.setHeader('Content-Type', 'application/json')
    res.send(swaggerSpec)
  })

  console.log(`Docs available at http://localhost:4000/api/messages/docs`)
}

export default swaggerDocs
