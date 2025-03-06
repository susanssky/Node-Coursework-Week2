import { Request, Response } from 'express'
import type { MessagesType } from '../utils/types'
import { db } from '../db/db'
import { messagesTable } from '../db/schema'
import { eq, desc } from 'drizzle-orm'
import redisClient from '../utils/redis' // Import Redis client

const validateMessage = (message: MessagesType) => {
  if (!message.from || !message.text) {
    return false
  }
  return true
}

// Helper function to cache data in Redis
const cacheData = async (key: string, data: any, expiration: number = 3600) => {
  await redisClient.setEx(key, expiration, JSON.stringify(data))
}

// Helper function to get cached data
const getCachedData = async (key: string) => {
  const cached = await redisClient.get(key)
  return cached ? JSON.parse(cached) : null
}

export const getAllMessages = async (req: Request, res: Response) => {
  const cacheKey = 'all_messages'
  const cachedMessages = await getCachedData(cacheKey)

  if (cachedMessages) {
    return res.send(cachedMessages)
  }

  const allMessages = await db
    .select()
    .from(messagesTable)
    .orderBy(messagesTable.id)

  await cacheData(cacheKey, allMessages)
  res.send(allMessages)
}

export const getLatestMessages = async (req: Request, res: Response) => {
  const cacheKey = 'latest_messages'
  const cachedMessages = await getCachedData(cacheKey)

  if (cachedMessages) {
    return res.send(cachedMessages)
  }

  const latestMessages = await db
    .select()
    .from(messagesTable)
    .orderBy(desc(messagesTable.id))
    .limit(10)

  await cacheData(cacheKey, latestMessages)
  return res.send(latestMessages)
}

export const getSearchedMessages = async (req: Request, res: Response) => {
  if (!req.query.text) return res.status(400).send('Search text is required')
  const searchText = req.query.text as string
  const cacheKey = `searched_messages:${searchText}`
  const cachedMessages = await getCachedData(cacheKey)

  if (cachedMessages) {
    return res.send(cachedMessages)
  }

  const matchedMessages = await db
    .select()
    .from(messagesTable)
    .orderBy(messagesTable.id)
    .where(eq(messagesTable.text, searchText))

  await cacheData(cacheKey, matchedMessages)
  return res.send(matchedMessages)
}

export const getMessage = async (req: Request, res: Response) => {
  const messagesId = req.params.messagesId
  const cacheKey = `message:${messagesId}`
  const cachedMessage = await getCachedData(cacheKey)

  if (cachedMessage) {
    return res.send(cachedMessage)
  }

  const oneMessage = await db
    .select()
    .from(messagesTable)
    .where(eq(messagesTable.id, Number(messagesId)))

  await cacheData(cacheKey, oneMessage)
  return res.send(oneMessage)
}

export const createMessage = async (req: Request, res: Response) => {
  const isValid = validateMessage(req.body)
  if (!isValid) {
    return res.status(400).json('Your name or message are missing.')
  }

  const newMessage = {
    ...req.body,
    timeSent: new Date().toISOString(),
  }
  const returnData = await db
    .insert(messagesTable)
    .values(newMessage)
    .returning()

  // Invalidate relevant caches
  await redisClient.del('all_messages')
  await redisClient.del('latest_messages')

  return res.status(201).json(returnData)
}

export const updateMessage = async (req: Request, res: Response) => {
  const messagesId = req.params.messagesId

  const isValid = validateMessage(req.body)
  if (!isValid) {
    return res.status(400).json('Your name or message are missing.')
  }

  const updatedMessage = {
    ...req.body,
  }

  const updatedData = await db
    .update(messagesTable)
    .set(updatedMessage)
    .where(eq(messagesTable.id, Number(messagesId)))
    .returning()

  // Invalidate relevant caches
  await redisClient.del('all_messages')
  await redisClient.del('latest_messages')
  await redisClient.del(`message:${messagesId}`)

  return res.send(updatedData)
}

export const deleteMessage = async (req: Request, res: Response) => {
  const messagesId = req.params.messagesId
  await db.delete(messagesTable).where(eq(messagesTable.id, Number(messagesId)))

  // Invalidate relevant caches
  await redisClient.del('all_messages')
  await redisClient.del('latest_messages')
  await redisClient.del(`message:${messagesId}`)

  const allMessages = await db.select().from(messagesTable)
  await cacheData('all_messages', allMessages) // Optional: Recache updated data
  return res.send(allMessages)
}
