import express, { Request, Response } from 'express'
import messages from './routes/messages'
import cors from 'cors'
import dotenv from 'dotenv'
dotenv.config()

import swaggerDocs from './utils/swaggerDoc'

import { db } from './db/db'
import { messagesTable } from './db/schema'

export const app = express()
const port: number = Number(process.env.PORT || 4000)

app.use(express.json())
app.use(cors())
app.use(express.urlencoded({ extended: false }))

swaggerDocs(app)
app.use('/api/messages', messages)

// app.get('/', function (req: Request, res: Response) {
//   res.sendFile(__dirname + '/index.html')
// })

app.listen(port, async () => {
  console.log(`listening on port ${port}...`)
  await seed()
})

async function seed() {
  const defaultData: typeof messagesTable.$inferInsert = {
    id: 0,
    from: 'Bart',
    text: 'Welcome to CYF chat system!',
  }
  await db.delete(messagesTable)
  await db.insert(messagesTable).values(defaultData)
}
