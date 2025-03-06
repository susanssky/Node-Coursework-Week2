import { Request, Response } from 'express'
import type { MessagesType } from '../utils/types'

const welcomeMessage: MessagesType = {
  id: (0).toString(),
  from: 'Bart',
  text: 'Welcome to CYF chat system!',
}

//This array is our "data store".
//We will start with one message in the array.
//Note: messages will be lost when Glitch restarts our server.
export let messages: MessagesType[] = [welcomeMessage]

const validateMessage = (message: MessagesType) => {
  if (!message.from || !message.text) {
    return false
  }
  return true
}

export const getAllMessages = (req: Request, res: Response) => {
  res.send(messages)
}
export const getLatestMessages = (req: Request, res: Response) => {
  res.send(
    messages.length <= 10
      ? messages.reverse()
      : messages.slice(messages.length - 10).reverse()
  )
}
export const getSearchedMessages = (req: Request, res: Response) => {
  if (!req.query.text) return
  const searchText = req?.query?.text as string
  res.send(messages.filter(({ text }) => text.includes(searchText)))
}

export const getMessage = (req: Request, res: Response) => {
  const messagesId = req.params.messagesId
  res.send(messages.filter(({ id }) => id === messagesId))
}

export const createMessage = (req: Request, res: Response) => {
  const isValid = validateMessage(req.body)
  if (!isValid) {
    return res.status(400).json(`Your name or message are missing.`)
  }

  const newMessage = {
    ...req.body,
    id: (Number(messages[messages.length - 1].id) + 1).toString(),
    timeSent: new Date().toISOString(),
  }

  messages.push(newMessage)
  return res.status(201).send(newMessage)
}

export const updateMessage = (req: Request, res: Response) => {
  const index = messages.findIndex(({ id }) => id === req.params.messagesId)
  const isValid = validateMessage(req.body)
  if (!isValid) {
    return res.status(400).json(`Your name or message are missing.`)
  }
  if (index === -1) {
    return res.status(404).json(`no message with id ${req.params.messagesId}`)
  }

  const updatedMessage = {
    ...messages[index],
    ...req.body,
  }
  messages[index] = updatedMessage
  return res.send(updatedMessage)
}
export const deleteMessage = (req: Request, res: Response) => {
  messages = messages.filter(({ id }) => id !== req.params.messagesId)
  res.send(messages)
}
