import { pgTable, serial, text, varchar } from 'drizzle-orm/pg-core'

export const messagesTable = pgTable('pgTable', {
  id: serial('id').primaryKey(),
  from: text('from').notNull(),
  text: text('text').notNull(),
  timeSent: text('time_sent'),
})
