import express from 'express'
import {
  createMessage,
  deleteMessage,
  getAllMessages,
  getLatestMessages,
  getMessage,
  getSearchedMessages,
  updateMessage,
} from '../controllers/messages'
const router = express.Router()

/**
 * @openapi
 * /api/messages:
 *   get:
 *     summary: Get all messages
 *     tags: [Messages]
 *     responses:
 *       200:
 *         description: All messages returned successfully
 *         content:
 *           application/json:
 *             schema:
 *               type: array
 *               items:
 *                 $ref: '#/components/schemas/Message'
 */
router.route('/').get(getAllMessages)

/**
 * @openapi
 * /api/messages:
 *   post:
 *     summary: Create a message
 *     tags: [Messages]
 *     requestBody:
 *       required: true
 *       content:
 *         application/json:
 *           schema:
 *             $ref: '#/components/schemas/MessageInput'
 *           example:
 *             from: "Alice"
 *             text: "Hello, world!"
 *     responses:
 *       201:
 *         description: Message created successfully
 *         content:
 *           application/json:
 *             schema:
 *               $ref: '#/components/schemas/Message'
 *       400:
 *         description: Required fields are missing
 */
router.route('/').post(createMessage)

/**
 * @openapi
 * /api/messages/latest:
 *   get:
 *     summary: Get the latest 10 messages
 *     tags: [Messages]
 *     responses:
 *       200:
 *         description: Successfully returned the latest message
 *         content:
 *           application/json:
 *             schema:
 *               type: array
 *               items:
 *                 $ref: '#/components/schemas/Message'
 */
router.route('/latest').get(getLatestMessages)

/**
 * @openapi
 * /api/messages/search:
 *   get:
 *     summary: Search for messages containing specific text
 *     tags: [Messages]
 *     parameters:
 *       - in: query
 *         name: text
 *         schema:
 *           type: string
 *         required: true
 *         description: The text to search for
 *     responses:
 *       200:
 *         description: Successfully return the message that meets the conditions
 *         content:
 *           application/json:
 *             schema:
 *               type: array
 *               items:
 *                 $ref: '#/components/schemas/Message'
 */
router.route('/search').get(getSearchedMessages)

/**
 * @openapi
 * /api/messages/{messagesId}:
 *   get:
 *     summary: Get a specific message by ID
 *     tags: [Messages]
 *     parameters:
 *       - in: path
 *         name: messagesId
 *         schema:
 *           type: string
 *         required: true
 *         description: The ID of the message
 *     responses:
 *       200:
 *         description: The ID of the message
 *         content:
 *           application/json:
 *             schema:
 *               type: array
 *               items:
 *                 $ref: '#/components/schemas/Message'
 */
router.route('/:messagesId').get(getMessage)

/**
 * @openapi
 * /api/messages/{messagesId}:
 *   put:
 *     summary: Update a specific message
 *     tags: [Messages]
 *     parameters:
 *       - in: path
 *         name: messagesId
 *         schema:
 *           type: string
 *         required: true
 *         description: message ID
 *     requestBody:
 *       required: true
 *       content:
 *         application/json:
 *           schema:
 *             $ref: '#/components/schemas/MessageInput'
 *           example:
 *             from: "Alice update"
 *             text: "Hello, world! update"
 *     responses:
 *       200:
 *         description: Message updated successfully
 *         content:
 *           application/json:
 *             schema:
 *               $ref: '#/components/schemas/Message'
 *       400:
 *         description: Required fields are missing
 *       404:
 *         description: The message could not be found
 */
router.route('/:messagesId').put(updateMessage)

/**
 * @openapi
 * /api/messages/{messagesId}:
 *   delete:
 *     summary: Delete a specific message
 *     tags: [Messages]
 *     parameters:
 *       - in: path
 *         name: messagesId
 *         schema:
 *           type: string
 *         required: true
 *         description: The ID of the message
 *     responses:
 *       200:
 *         description: The message was deleted successfully, and the remaining messages are returned
 *         content:
 *           application/json:
 *             schema:
 *               type: array
 *               items:
 *                 $ref: '#/components/schemas/Message'
 */
router.route('/:messagesId').delete(deleteMessage)

export default router
