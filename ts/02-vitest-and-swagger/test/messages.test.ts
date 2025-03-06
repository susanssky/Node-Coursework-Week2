// test/integration/messages.test.ts
import { describe, it, expect, beforeEach } from 'vitest'
import request from 'supertest'
import { app } from '../server'
import type { MessagesType } from '../utils/types'

describe('Messages API Integration Tests', () => {
  beforeEach(async () => {
    // Reset initial state via API
    await request(app).post('/api/messages').send({
      from: 'Bart',
      text: 'Welcome to CYF chat system!',
    })
  })

  it('GET /api/messages shall return all messages', async () => {
    const response = await request(app).get('/api/messages')
    expect(response.status).toBe(200)
    expect(response.body).toContainEqual(
      expect.objectContaining({
        from: 'Bart',
        text: 'Welcome to CYF chat system!',
      })
    )
  })

  it('GET /api/messages/latest shall return the latest 10 messages when more than 10 exist', async () => {
    for (let i = 0; i < 15; i++) {
      await request(app)
        .post('/api/messages')
        .send({
          from: `User${i}`,
          text: `Message ${i}`,
        })
    }
    const response = await request(app).get('/api/messages/latest')
    expect(response.status).toBe(200)
    expect(response.body.length).toBe(10)
    expect(response.body[0].text).toBe('Message 14')
    expect(response.body[9].text).toBe('Message 5')
  })

  it('GET /api/messages/search shall return messages matching the search criteria', async () => {
    await request(app)
      .post('/api/messages')
      .send({ from: 'Alice', text: 'Hello there' })
    const response = await request(app)
      .get('/api/messages/search')
      .query({ text: 'Hello' })
    expect(response.status).toBe(200)
    expect(response.body).toContainEqual(
      expect.objectContaining({ from: 'Alice', text: 'Hello there' })
    )
  })

  it('POST /api/messages shall create a new message and allow retrieval', async () => {
    const newMessage = { from: 'Charlie', text: 'New message' }
    const postResponse = await request(app)
      .post('/api/messages')
      .send(newMessage)
    expect(postResponse.status).toBe(201)
    expect(postResponse.body).toMatchObject(newMessage)

    const getResponse = await request(app).get('/api/messages')
    expect(getResponse.body).toContainEqual(postResponse.body)
  })

  it('PUT /api/messages/:messagesId shall update a message and reflect in retrieval', async () => {
    const createResponse = await request(app)
      .post('/api/messages')
      .send({ from: 'Alice', text: 'Original' })
    const messageId = createResponse.body.id

    const updatedMessage = { from: 'Alice Updated', text: 'Updated' }
    const putResponse = await request(app)
      .put(`/api/messages/${messageId}`)
      .send(updatedMessage)
    expect(putResponse.status).toBe(200)
    expect(putResponse.body).toMatchObject(updatedMessage)

    const getResponse = await request(app).get(`/api/messages/${messageId}`)
    expect(getResponse.body[0]).toMatchObject(updatedMessage)
  })

  it('DELETE /api/messages/:messagesId shall delete a message and remove it from the list', async () => {
    const createResponse = await request(app)
      .post('/api/messages')
      .send({ from: 'Bob', text: 'To delete' })
    const messageId = createResponse.body.id

    const deleteResponse = await request(app).delete(
      `/api/messages/${messageId}`
    )
    expect(deleteResponse.status).toBe(200)

    const getResponse = await request(app).get('/api/messages')
    expect(
      getResponse.body.find((msg: MessagesType) => msg.id === messageId)
    ).toBeUndefined()
  })
})
